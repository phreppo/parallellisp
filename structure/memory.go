package structure

type Memory struct {
}

func (m Memory) MakeInt(i int) *Cell {
	return &IntCell{val: i}
}
