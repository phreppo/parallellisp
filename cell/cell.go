package cell

import (
	"fmt"
	"strconv"
)

type Cell interface {
	IsAtom() bool
}

/*******************************************************************************
 Int cell
*******************************************************************************/

type IntCell struct {
	Val int
}

func (i IntCell) IsAtom() bool {
	return true
}

func (i IntCell) String() string {
	return strconv.Itoa(i.Val)
}

/*******************************************************************************
 String cell
*******************************************************************************/

type StringCell struct {
	Str string
}

func (s StringCell) IsAtom() bool {
	return true
}

func (s StringCell) String() string {
	return ("\"" + s.Str + "\"")
}

/*******************************************************************************
 Symbol cell
*******************************************************************************/

type SymbolCell struct {
	Sym string
}

func (s SymbolCell) IsAtom() bool {
	return true
}

func (s SymbolCell) String() string {
	return s.Sym
}

/*******************************************************************************
 Cons cell
*******************************************************************************/

type ConsCell struct {
	Car Cell
	Cdr Cell
}

func (c ConsCell) IsAtom() bool {
	return false
}

func (c ConsCell) String() string {
	left := fmt.Sprintf("%v", c.Car)
	right := fmt.Sprintf("%v", c.Cdr)
	return "(" + left + " . " + right + ")"
}
