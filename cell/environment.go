package cell

import (
	"fmt"
	"runtime"
)

func EmptyEnv() *EnvironmentEntry {
	return nil
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
	GlobalEnv["null"], _ = Parse("(lambda (x) (eq x nil))")
	GlobalEnv["ncpu"], _ = Parse(fmt.Sprintf("%v", runtime.NumCPU()))
	GlobalEnv["t"], _ = Parse("t")

	// shortcuts
	GlobalEnv["tests"], _ = Parse("\"test.lisp\"")
	GlobalEnv["search"], _ = Parse("\"run-search.lisp\"")
	GlobalEnv["sum"], _ = Parse("\"run-sum.lisp\"")
	GlobalEnv["fib"], _ = Parse("\"run-fib.lisp\"")
	GlobalEnv["ms"], _ = Parse("\"run-mergesort.lisp\"")
	GlobalEnv["sorted"], _ = Parse("\"psorted.lisp\"")
}
