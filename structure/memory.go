package structure

func MakeInt(i int, m *Memory, ans chan Cell) Cell {
	m.MakeInt <- IntRequest{i, ans}
	intCell := <-ans
	return intCell
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

type ConsRequest struct {
	Car        Cell
	Cdr        Cell
	AnswerChan chan<- Cell
}

type Memory struct {
	MakeInt    chan IntRequest
	intFactory *intCellSupplier

	MakeCons    chan ConsRequest
	consFactory *consCellSupplier
}

func NewMemory() *Memory {
	m := Memory{
		MakeInt:     make(chan IntRequest),
		intFactory:  newIntCellSupplier(),
		MakeCons:    make(chan ConsRequest),
		consFactory: newConsCellSupplier(),
	}

	go func() {
		for {
			select {
			case request := <-m.MakeInt:
				m.supplyInt(request)

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

func (m *Memory) supplyCons(request ConsRequest) {
	m.consFactory.makeCons <- request
}

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
