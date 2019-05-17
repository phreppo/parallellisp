package cell

import (
	"fmt"
	"strconv"
)

type Cell interface {
	IsAtom() bool
	Eq(Cell) bool
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

func (i IntCell) Eq(c Cell) bool {
	switch castedC := c.(type) {
	case *IntCell:
		return castedC.Val == i.Val
	default:
		return false
	}
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

func (s StringCell) Eq(c Cell) bool {
	switch castedC := c.(type) {
	case *StringCell:
		return castedC.Str == s.Str
	default:
		return false
	}
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

func (s SymbolCell) Eq(c Cell) bool {
	switch castedC := c.(type) {
	case *SymbolCell:
		return castedC.Sym == s.Sym
	default:
		return false
	}
}

/*******************************************************************************
 Builtin lambda cell
*******************************************************************************/

type BuiltinLambdaCell struct {
	Sym    string
	Lambda func(Cell, *EnvironmentEntry) EvalResult
}

func (l BuiltinLambdaCell) IsAtom() bool {
	return true
}

func (l BuiltinLambdaCell) String() string {
	return l.Sym
}

func (l BuiltinLambdaCell) Eq(c Cell) bool {
	switch castedC := c.(type) {
	case *BuiltinLambdaCell:
		return castedC.Sym == l.Sym
	default:
		return false
	}
}

/*******************************************************************************
 Builtin macro cell
*******************************************************************************/

type BuiltinMacroCell struct {
	Sym   string
	Macro func(Cell, *EnvironmentEntry) EvalResult
}

func (m BuiltinMacroCell) IsAtom() bool {
	return true
}

func (m BuiltinMacroCell) String() string {
	if m.Sym == "quote" {
		return "'"
	}
	return m.Sym
}

func (l BuiltinMacroCell) Eq(c Cell) bool {
	switch castedC := c.(type) {
	case *BuiltinMacroCell:
		return castedC.Sym == l.Sym
	default:
		return false
	}
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
	rest := ""
	act := c.Cdr
	for act != nil {
		switch cell := act.(type) {
		case *ConsCell:
			rest += fmt.Sprintf(" %v", cell.Car)
			act = cell.Cdr
		default:
			rest += fmt.Sprintf(" . %v", act)
			act = nil
		}
	}
	if left == "'" {
		return "'" + rest[1:] // skip first char
	} else {
		return "(" + left + rest + ")"
	}
}

func (cons1 ConsCell) Eq(cons2 Cell) bool {
	switch castedCons2 := cons2.(type) {
	case *ConsCell:
		return eq(cons1.Car, castedCons2.Car) && eq(cons1.Cdr, castedCons2.Cdr)
	default:
		return false
	}
}

func eq(c1, c2 Cell) bool {
	if c1 == nil && c2 == nil {
		return true
	} else {
		return c1.Eq(c2)
	}
}
