package structure

type Cell interface {
	IsAtom() bool
}

type IntCell struct {
	Val int
}

func (i IntCell) IsAtom() bool {
	return true
}

type StringCell struct {
	Str string
}

func (s StringCell) IsAtom() bool {
	return true
}

type SymbolCell struct {
	Sym string
}

func (s SymbolCell) IsAtom() bool {
	return true
}

type ConsCell struct {
	Car *Cell
	Cdr *Cell
}

func (c ConsCell) IsAtom() bool {
	return false
}
