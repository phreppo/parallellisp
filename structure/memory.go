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
	MakeCons   chan ConsRequest

	intTape        [100]IntCell
	intTapePointer int

	consTape        [100]ConsCell
	consTapePointer int
}

func NewMemory() *Memory {
	m := Memory{
		MakeInt:         make(chan IntRequest),
		intFactory:      new(),
		MakeCons:        make(chan ConsRequest),
		intTapePointer:  1,
		consTapePointer: 1,
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
	newCons := &(m.consTape[m.consTapePointer])
	newCons.Car = request.Car
	newCons.Cdr = request.Cdr
	m.consTapePointer++
	request.AnswerChan <- newCons
}

const INT_TAPE_SIZE = 100

type intCellSupplier struct {
	makeInt     chan IntRequest
	tape        [INT_TAPE_SIZE]IntCell
	tapePointer int
}

func new() *intCellSupplier {
	supplier := intCellSupplier{
		makeInt:     make(chan IntRequest),
		tapePointer: 1,
	}

	go func() {
		for {
			request := <-supplier.makeInt
			newInt := &(supplier.tape[supplier.tapePointer])
			newInt.Val = request.Val
			supplier.tapePointer++
			request.AnswerChan <- newInt
		}
	}()

	return &supplier
}
