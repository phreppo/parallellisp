package lisp

import (
	"fmt"
	"runtime"
)

func emptyEnv() *environmentEntry {
	return nil
}

func newEnvironmentEntry(sym *symbolCell, value Cell, next *environmentEntry) *environmentEntry {
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

func initGlobalEnv() {
	// necessary
	globalEnv["id"], _ = Parse("(lambda (x) x)")
	globalEnv["t"], _ = Parse("t")
	globalEnv["null"], _ = Parse("(lambda (x) (eq x nil))")
	globalEnv["ncpu"], _ = Parse(fmt.Sprintf("%v", runtime.NumCPU()))

	globalEnv["take"], _ = Parse("(lambda (lst n) (cond ((eq n 0) nil) (t (cons (car lst) (take (cdr lst) (1- n))))))")
	globalEnv["drop"], _ = Parse("(lambda (lst n) (cond ((eq n 0) lst) (t (drop (cdr lst) (1- n)))))")
	globalEnv["first-half"], _ = Parse("(lambda (lst) (take lst (/ (length lst) 2)))")
	globalEnv["second-half"], _ = Parse("(lambda (lst) (drop lst (/ (length lst) 2)))")

	globalEnv["parallelize"], _ = Parse("(lambda (sequential-algorithm is-base-case split-left split-right combinator  generic-data) (parallelize-ric  1 sequential-algorithm is-base-case split-left split-right combinator  generic-data))")
	globalEnv["parallelize-ric"], _ = Parse("(lambda (partitions sequential-algorithm is-base-case split-left split-right combinator generic-data) (cond ((is-base-case generic-data) (sequential-algorithm generic-data)) ((< partitions ncpu) (let ((new-partitions (* partitions 2))) {combinator (parallelize-ric new-partitions sequential-algorithm is-base-case split-right split-left combinator (split-left generic-data)) (parallelize-ric new-partitions sequential-algorithm is-base-case split-right split-left combinator (split-right generic-data)) })) (t (combinator (sequential-algorithm (split-left generic-data)) (sequential-algorithm (split-right generic-data)) ))))")

	globalEnv["divide-et-impera"], _ = Parse("(lambda (sequential-algorithm combinator lst) (divide-et-impera-ric 1 sequential-algorithm combinator lst))")
	globalEnv["divide-et-impera-ric"], _ = Parse("(lambda (partitions sequential-algorithm combinator lst) (cond ((eq lst nil)        (sequential-algorithm lst)) ((eq (length lst) 1) (sequential-algorithm lst)) ((< partitions ncpu) (let ((new-partitions (* partitions 2))) {combinator (divide-et-impera-ric new-partitions sequential-algorithm combinator (first-half  lst)) (divide-et-impera-ric new-partitions sequential-algorithm combinator (second-half lst)) })) (t (combinator (sequential-algorithm (first-half  lst)) (sequential-algorithm (second-half lst)) ))))")

	// shortcuts
	globalEnv["tests"], _ = Parse("\"test.lisp\"")
	globalEnv["search"], _ = Parse("\"run-search.lisp\"")
	globalEnv["sum"], _ = Parse("\"run-sum.lisp\"")
	globalEnv["fib"], _ = Parse("\"run-fib.lisp\"")
	globalEnv["ms"], _ = Parse("\"run-mergesort.lisp\"")
	globalEnv["sorted"], _ = Parse("\"run-sorted.lisp\"")
	globalEnv["omega"], _ = Parse("(lambda (x) (x x))")
}
