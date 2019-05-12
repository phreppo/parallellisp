package cell

import (
	"fmt"
	"strconv"
)

type Cell interface {
	IsAtom() bool
	IsBuiltinLambda() (func(), bool)
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

func (i IntCell) IsBuiltinLambda() (func(), bool) {
	return nil, false
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

func (s StringCell) IsBuiltinLambda() (func(), bool) {
	return nil, false
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

func (s SymbolCell) IsBuiltinLambda() (func(), bool) {
	return nil, false
}

/*******************************************************************************
 Builtin lambda cell
*******************************************************************************/

type BuiltinLambdaCell struct {
	Sym    string
	Lambda func()
}

func (l BuiltinLambdaCell) IsAtom() bool {
	return true
}

func (l BuiltinLambdaCell) String() string {
	return l.Sym
}

func (l BuiltinLambdaCell) IsBuiltinLambda() (func(), bool) {
	return l.Lambda, true
}

/*******************************************************************************
 Builtin macro cell
*******************************************************************************/

type BuiltinMacroCell struct {
	Sym   string
	Macro func()
}

func (m BuiltinMacroCell) IsAtom() bool {
	return true
}

func (m BuiltinMacroCell) String() string {
	return m.Sym
}

func (m BuiltinMacroCell) IsBuiltinLambda() (func(), bool) {
	return nil, false
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

func (c ConsCell) IsBuiltinLambda() (func(), bool) {
	return nil, false
}
