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
 Builtin lambda cell
*******************************************************************************/

type BuiltinLambdaCell struct {
	Sym    string
	Lambda func(Cell, *EnvironmentEntry) *EvalResult
}

func (l BuiltinLambdaCell) IsAtom() bool {
	return true
}

func (l BuiltinLambdaCell) String() string {
	return l.Sym
}

/*******************************************************************************
 Builtin macro cell
*******************************************************************************/

type BuiltinMacroCell struct {
	Sym   string
	Macro func(Cell, *EnvironmentEntry) *EvalResult
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
	// left := fmt.Sprintf("%v", c.Car)
	// right := fmt.Sprintf("%v", c.Cdr)
	// return "(" + left + " . " + right + ")"

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
