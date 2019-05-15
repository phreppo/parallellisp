package cell

// Mem is the global memory
var Mem = newMemory()

// MakeInt supplies a Int Cell. Blocking: use only in sequential code
func MakeInt(i int) Cell {
	ans := make(chan Cell)
	Mem.MakeInt <- IntRequest{i, ans}
	intCell := <-ans
	return intCell
}

// MakeString supplies a String Cell. Blocking: use only in sequential code
func MakeString(s string) Cell {
	ans := make(chan Cell)
	Mem.MakeString <- StringRequest{s, ans}
	stringCell := <-ans
	return stringCell
}

// MakeSymbol supplies a Symbol Cell. Blocking: use only in sequential code
func MakeSymbol(s string) Cell {
	ans := make(chan Cell)
	Mem.MakeSymbol <- SymbolRequest{s, ans}
	symbolCell := <-ans
	return symbolCell
}

// MakeCons supplies a Cons Cell. Blocking: use only in sequential code
func MakeCons(car Cell, cdr Cell) Cell {
	ans := make(chan Cell)
	Mem.MakeCons <- ConsRequest{car, cdr, ans}
	consCell := <-ans
	return consCell
}

// IntRequest allows to build Int Cell Requests for a memory
type IntRequest struct {
	Val        int
	AnswerChan chan<- Cell
}

// StringRequest allows to build String Cell Requests for a memory
type StringRequest struct {
	Str        string
	AnswerChan chan<- Cell
}

// SymbolRequest allows to build Symbol Cell Requests for a memory
type SymbolRequest struct {
	Sym        string
	AnswerChan chan<- Cell
}

// ConsRequest allows to build Cons Cell Requests for a memory
type ConsRequest struct {
	Car        Cell
	Cdr        Cell
	AnswerChan chan<- Cell
}

type memory struct {
	MakeInt    chan IntRequest
	intFactory *intCellSupplier

	MakeString    chan StringRequest
	stringFactory *stringCellSupplier

	MakeSymbol    chan SymbolRequest
	symbolFactory *symbolCellSupplier

	MakeCons    chan ConsRequest
	consFactory *consCellSupplier
}

func newMemory() *memory {
	m := memory{
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
				m.intFactory.makeInt <- request

			case request := <-m.MakeString:
				m.stringFactory.makeString <- request

			case request := <-m.MakeSymbol:
				m.symbolFactory.makeSymbol <- request

			case request := <-m.MakeCons:
				m.consFactory.makeCons <- request
			}
		}
	}()

	return &m
}

/*******************************************************************************
 Factories
*******************************************************************************/

const intTapeSize = 10000

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

const stringTapeSize = 10000

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

const symbolTapeSize = 10000

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
			if isBuiltin, builtinSymbol := Lisp.IsBuiltinSymbol(request.Sym); isBuiltin {
				request.AnswerChan <- builtinSymbol
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

const consTapeSize = 10000

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
