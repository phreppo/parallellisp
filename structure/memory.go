package structure

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
}

func NewMemory() *Memory {
	m := Memory{
		make(chan IntRequest),
		make(chan ConsRequest),
	}

	go func() {
		for {
			select {
			case request := <-m.MakeInt:
				newInt := new(IntCell)
				newInt.Val = request.Val
				request.AnswerChan <- newInt

			case request := <-m.MakeCons:
				newCons := new(ConsCell)
				newCons.Car = request.Car
				newCons.Cdr = request.Cdr
				request.AnswerChan <- newCons
			}
		}
	}()

	return &m
}

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
