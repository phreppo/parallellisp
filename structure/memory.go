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
	MakeInt  chan IntRequest
	MakeCons chan ConsRequest

	intTape        [100]IntCell
	intTapePointer int

	consTape        [100]ConsCell
	consTapePointer int
}

func NewMemory() *Memory {
	m := Memory{
		MakeInt:         make(chan IntRequest),
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
	newInt := &(m.intTape[m.intTapePointer])
	newInt.Val = request.Val
	m.intTapePointer++
	request.AnswerChan <- newInt
}

func (m *Memory) supplyCons(request ConsRequest) {
	newCons := &(m.consTape[m.consTapePointer])
	newCons.Car = request.Car
	newCons.Cdr = request.Cdr
	m.consTapePointer++
	request.AnswerChan <- newCons
}
