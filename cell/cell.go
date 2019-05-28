package cell

import (
	"fmt"
	"strconv"
)

/*******************************************************************************
 Int cell
*******************************************************************************/

type intCell struct {
	Val int
}

func (i intCell) String() string {
	return strconv.Itoa(i.Val)
}

func (i intCell) Eq(c Cell) bool {
	switch castedC := c.(type) {
	case *intCell:
		return castedC.Val == i.Val
	default:
		return false
	}
}

/*******************************************************************************
 String cell
*******************************************************************************/

type stringCell struct {
	Str string
}

func (s stringCell) String() string {
	return ("\"" + s.Str + "\"")
}

func (s stringCell) Eq(c Cell) bool {
	switch castedC := c.(type) {
	case *stringCell:
		return castedC.Str == s.Str
	default:
		return false
	}
}

/*******************************************************************************
 Symbol cell
*******************************************************************************/

type symbolCell struct {
	Sym string
}

func (s symbolCell) String() string {
	return s.Sym
}

func (s symbolCell) Eq(c Cell) bool {
	switch castedC := c.(type) {
	case *symbolCell:
		return castedC.Sym == s.Sym
	default:
		return false
	}
}

/*******************************************************************************
 Builtin lambda cell
*******************************************************************************/

type builtinLambdaCell struct {
	Sym    string
	Lambda func(Cell, *environmentEntry) EvalResult
}

func (l builtinLambdaCell) String() string {
	return l.Sym
}

func (l builtinLambdaCell) Eq(c Cell) bool {
	switch castedC := c.(type) {
	case *builtinLambdaCell:
		return castedC.Sym == l.Sym
	default:
		return false
	}
}

/*******************************************************************************
 Builtin macro cell
*******************************************************************************/

type builtinMacroCell struct {
	Sym   string
	Macro func(Cell, *environmentEntry) EvalResult
}

func (m builtinMacroCell) String() string {
	if m.Sym == "quote" {
		return "'"
	}
	return m.Sym
}

func (m builtinMacroCell) Eq(c Cell) bool {
	switch castedC := c.(type) {
	case *builtinMacroCell:
		return castedC.Sym == m.Sym
	default:
		return false
	}
}

/*******************************************************************************
 Cons cell
*******************************************************************************/

type consCell struct {
	Car   Cell
	Cdr   Cell
	Evlis func(args Cell, env *environmentEntry) EvalResult
}

func (c consCell) String() string {
	left := fmt.Sprintf("%v", c.Car)
	rest := ""
	act := c.Cdr
	for act != nil {
		switch cell := act.(type) {
		case *consCell:
			rest += fmt.Sprintf(" %v", cell.Car)
			act = cell.Cdr
		default:
			rest += fmt.Sprintf(" . %v", act)
			act = nil
		}
	}
	if left == "'" {
		return "'" + rest[1:] // skip first char
	}
	return "(" + left + rest + ")"
}

func (c consCell) Eq(cons2 Cell) bool {
	switch castedCons2 := cons2.(type) {
	case *consCell:
		return eq(c.Car, castedCons2.Car) && eq(c.Cdr, castedCons2.Cdr)
	default:
		return false
	}
}
