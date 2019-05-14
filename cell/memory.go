package cell

func MakeInt(i int, m *Memory, ans chan Cell) Cell {
	m.MakeInt <- IntRequest{i, ans}
	intCell := <-ans
	return intCell
}

func MakeString(s string, m *Memory, ans chan Cell) Cell {
	m.MakeString <- StringRequest{s, ans}
	stringCell := <-ans
	return stringCell
}

func MakeSymbol(s string, m *Memory, ans chan Cell) Cell {
	m.MakeSymbol <- SymbolRequest{s, ans}
	symbolCell := <-ans
	return symbolCell
}

func MakeCons(car Cell, cdr Cell, m *Memory, ans chan Cell) Cell {
	m.MakeCons <- ConsRequest{car, cdr, ans}
	consCell := <-ans
	return consCell
}

type IntRequest struct {
	Val        int
	AnswerChan chan<- Cell
}

type StringRequest struct {
	Str        string
	AnswerChan chan<- Cell
}

type SymbolRequest struct {
	Sym        string
	AnswerChan chan<- Cell
}

type ConsRequest struct {
	Car        Cell
	Cdr        Cell
	AnswerChan chan<- Cell
}

type Memory struct {
	MakeInt    chan IntRequest
	intFactory *intCellSupplier

	MakeString    chan StringRequest
	stringFactory *stringCellSupplier

	MakeSymbol    chan SymbolRequest
	symbolFactory *symbolCellSupplier

	MakeCons    chan ConsRequest
	consFactory *consCellSupplier

	BuiltinLambdas [3]BuiltinLambdaCell
	BuiltinMacros  [3]BuiltinMacroCell
}

func NewMemory() *Memory {
	m := Memory{
		MakeInt:       make(chan IntRequest),
		intFactory:    newIntCellSupplier(),
		MakeString:    make(chan StringRequest),
		stringFactory: newStringCellSupplier(),
		MakeSymbol:    make(chan SymbolRequest),
		symbolFactory: newSymbolCellSupplier(),
		MakeCons:      make(chan ConsRequest),
		consFactory:   newConsCellSupplier(),
	}

	go func() {
		for {
			select {
			case request := <-m.MakeInt:
				m.supplyInt(request)

			case request := <-m.MakeString:
				m.supplyString(request)

			case request := <-m.MakeSymbol:
				m.supplySymbol(request)

			case request := <-m.MakeCons:
				m.supplyCons(request)
			}
		}
	}()

	return &m
}

func (m *Memory) supplyInt(request IntRequest) {
	m.intFactory.makeInt <- request
}

func (m *Memory) supplyString(request StringRequest) {
	m.stringFactory.makeString <- request
}

func (m *Memory) supplySymbol(request SymbolRequest) {
	m.symbolFactory.makeSymbol <- request
}

func (m *Memory) supplyCons(request ConsRequest) {
	m.consFactory.makeCons <- request
}

/*******************************************************************************
 Factories
*******************************************************************************/

const intTapeSize = 100

type intCellSupplier struct {
	makeInt     chan IntRequest
	tape        *[intTapeSize]IntCell
	tapePointer int
}

func newIntCellSupplier() *intCellSupplier {
	supplier := intCellSupplier{
		makeInt:     make(chan IntRequest),
		tape:        new([intTapeSize]IntCell),
		tapePointer: 0,
	}

	go func() {
		for {
			request := <-supplier.makeInt
			if supplier.tapePointer >= intTapeSize {
				supplier.tape = new([intTapeSize]IntCell)
				supplier.tapePointer = 0
			}
			newInt := &(supplier.tape[supplier.tapePointer])
			newInt.Val = request.Val
			supplier.tapePointer++
			request.AnswerChan <- newInt
		}
	}()

	return &supplier
}

const stringTapeSize = 100

type stringCellSupplier struct {
	makeString  chan StringRequest
	tape        *[stringTapeSize]StringCell
	tapePointer int
}

func newStringCellSupplier() *stringCellSupplier {
	supplier := stringCellSupplier{
		makeString:  make(chan StringRequest),
		tape:        new([stringTapeSize]StringCell),
		tapePointer: 0,
	}

	go func() {
		for {
			request := <-supplier.makeString
			if supplier.tapePointer >= stringTapeSize {
				supplier.tape = new([stringTapeSize]StringCell)
				supplier.tapePointer = 0
			}
			newString := &(supplier.tape[supplier.tapePointer])
			newString.Str = request.Str
			supplier.tapePointer++
			request.AnswerChan <- newString
		}
	}()

	return &supplier
}

const symbolTapeSize = 100

type symbolCellSupplier struct {
	makeSymbol  chan SymbolRequest
	tape        *[symbolTapeSize]SymbolCell
	tapePointer int
}

func newSymbolCellSupplier() *symbolCellSupplier {
	supplier := symbolCellSupplier{
		makeSymbol:  make(chan SymbolRequest),
		tape:        new([symbolTapeSize]SymbolCell),
		tapePointer: 0,
	}

	go func() {
		for {
			request := <-supplier.makeSymbol
			if request.Sym == "nil" {
				request.AnswerChan <- nil
			} else if builtinLambdaCell, ok := BuiltinLambdas[request.Sym]; ok {
				request.AnswerChan <- &builtinLambdaCell
			} else if builtinMacroCell, ok := BuiltinMacros[request.Sym]; ok {
				request.AnswerChan <- &builtinMacroCell
			} else if request.Sym == "t" {
				request.AnswerChan <- &TrueSymbol
			} else {
				if supplier.tapePointer >= symbolTapeSize {
					supplier.tape = new([symbolTapeSize]SymbolCell)
					supplier.tapePointer = 0
				}
				newSymbol := &(supplier.tape[supplier.tapePointer])
				newSymbol.Sym = request.Sym
				supplier.tapePointer++
				request.AnswerChan <- newSymbol
			}
		}
	}()

	return &supplier
}

const consTapeSize = 100

type consCellSupplier struct {
	makeCons    chan ConsRequest
	tape        *[consTapeSize]ConsCell
	tapePointer int
}

func newConsCellSupplier() *consCellSupplier {
	supplier := consCellSupplier{
		makeCons:    make(chan ConsRequest),
		tape:        new([consTapeSize]ConsCell),
		tapePointer: 0,
	}

	go func() {
		for {
			request := <-supplier.makeCons
			if supplier.tapePointer >= consTapeSize {
				supplier.tape = new([consTapeSize]ConsCell)
				supplier.tapePointer = 0
			}
			newCons := &(supplier.tape[supplier.tapePointer])
			newCons.Car = request.Car
			newCons.Cdr = request.Cdr
			supplier.tapePointer++
			request.AnswerChan <- newCons
		}
	}()

	return &supplier
}
