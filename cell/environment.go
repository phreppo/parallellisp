package cell

import (
	"fmt"
	"runtime"
)

func emptyEnv() *environmentEntry {
	return nil
}

func newenvironmentEntry(sym *symbolCell, value Cell, next *environmentEntry) *environmentEntry {
	newEntry := new(environmentEntry)
	newEntry.Pair = new(environmentPair)
	newEntry.Pair.Symbol = sym
	newEntry.Pair.Value = value
	newEntry.Next = next
	return newEntry
}

type environmentEntry struct {
	Pair *environmentPair
	Next *environmentEntry
}

func symbolIsInEnv(c *symbolCell, env *environmentEntry) bool {
	if env == nil {
		return false
	}
	act := env
	for act != nil {
		if (act.Pair.Symbol.Sym) == c.Sym {
			return true
		}
		act = act.Next
	}
	return false
}

type environmentPair struct {
	Symbol *symbolCell
	Value  Cell
}

func (e *environmentEntry) String() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("%v -> %v\n", e.Pair.Symbol, e.Pair.Value) + fmt.Sprintf("%v", e.Next)
}

var globalEnv = make(map[string]Cell)

func initglobalEnv() {
	// necessary
	globalEnv["id"], _ = Parse("(lambda (x) x)")
	globalEnv["null"], _ = Parse("(lambda (x) (eq x nil))")
	globalEnv["ncpu"], _ = Parse(fmt.Sprintf("%v", runtime.NumCPU()))
	globalEnv["t"], _ = Parse("t")

	// shortcuts
	globalEnv["tests"], _ = Parse("\"test.lisp\"")
	globalEnv["search"], _ = Parse("\"run-search.lisp\"")
	globalEnv["sum"], _ = Parse("\"run-sum.lisp\"")
	globalEnv["fib"], _ = Parse("\"run-fib.lisp\"")
	globalEnv["ms"], _ = Parse("\"run-mergesort.lisp\"")
	globalEnv["sorted"], _ = Parse("\"run-sorted.lisp\"")
	globalEnv["omega"], _ = Parse("(lambda (x) (x x))")
}
