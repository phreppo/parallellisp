package cell

import (
	"fmt"
	"runtime"
)

func EmptyEnv() *EnvironmentEntry {
	return nil
}

func SimpleEnv() *EnvironmentEntry {
	list, _ := Parse("(1 2 3 4)")
	return NewEnvironmentEntry(MakeSymbol("l").(*SymbolCell), list, nil)
}

func NewEnvironmentEntry(sym *SymbolCell, value Cell, next *EnvironmentEntry) *EnvironmentEntry {
	newEntry := new(EnvironmentEntry)
	newEntry.Pair = new(EnvironmentPair)
	newEntry.Pair.Symbol = sym
	newEntry.Pair.Value = value
	newEntry.Next = next
	return newEntry
}

type EnvironmentEntry struct {
	Pair *EnvironmentPair
	Next *EnvironmentEntry
}

func (e *EnvironmentEntry) String() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("%v -> %v\n", e.Pair.Symbol, e.Pair.Value) + fmt.Sprintf("%v", e.Next)
}

type EnvironmentPair struct {
	Symbol *SymbolCell
	Value  Cell
}

var GlobalEnv = make(map[string]Cell)

func initGlobalEnv() {
	// necessary
	GlobalEnv["id"], _ = Parse("(lambda (x) x)")
	GlobalEnv["ncpu"], _ = Parse(fmt.Sprintf("%v", runtime.NumCPU()))
	GlobalEnv["t"], _ = Parse("t")

	// shortcuts
	GlobalEnv["tests"], _ = Parse("\"test.lisp\"")
	GlobalEnv["search"], _ = Parse("\"psearch.lisp\"")
	GlobalEnv["sum"], _ = Parse("\"psum.lisp\"")
	GlobalEnv["fib"], _ = Parse("\"pfib.lisp\"")
	GlobalEnv["ms"], _ = Parse("\"pmergesort.lisp\"")
	GlobalEnv["sorted"], _ = Parse("\"psorted.lisp\"")
}
