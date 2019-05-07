package structure

type CellPair struct {
	Left  Cell
	Right Cell
}

type Memory struct {
	MakeIntChan chan int
	TakeIntChan chan Cell

	MakeConsChan chan CellPair
	TakeConsChan chan Cell
}

func NewMemory() *Memory {
	m := Memory{
		make(chan int),
		make(chan Cell),
		make(chan CellPair),
		make(chan Cell),
	}

	go func() {
		for {
			select {
			case x := <-m.MakeIntChan:
				newInt := new(IntCell)
				newInt.Val = x
				m.TakeIntChan <- newInt
			case pair := <-m.MakeConsChan:
				newCons := new(ConsCell)
				newCons.Car = pair.Left
				newCons.Cdr = pair.Right
				m.TakeConsChan <- newCons
			}
		}
	}()

	return &m
}
