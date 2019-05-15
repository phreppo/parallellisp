package cell

func EmptyEnv() *EnvironmentEntry {
	return nil
}

type EnvironmentEntry struct {
	Pair *EnvironmentPair
	Next *EnvironmentEntry
}

type EnvironmentPair struct {
	Symbol *SymbolCell
	Value  Cell
}
