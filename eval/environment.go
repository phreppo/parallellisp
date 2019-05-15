package eval

import . "github.com/parof/parallellisp/cell"

type EnvironmentEntry struct {
	Pair EnvironmentPair
	Next *EnvironmentEntry
}

type EnvironmentPair struct {
	Symbol Cell
	Value  Cell
}
