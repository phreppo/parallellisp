package cell

// makeInt supplies a Int Cell. Blocking: use only in sequential code
func makeInt(i int) Cell {
	return &IntCell{i}
}

// makeString supplies a String Cell. Blocking: use only in sequential code
func makeString(s string) Cell {
	return &StringCell{s}
}

// makeSymbol supplies a Symbol Cell. Blocking: use only in sequential code
func makeSymbol(s string) Cell {
	if isBuiltin, builtinSymbol := lisp.isBuiltinSymbol(s); isBuiltin {
		return builtinSymbol
	}
	return &SymbolCell{s}
}

// makeCons supplies a Cons Cell. Blocking: use only in sequential code
func makeCons(car Cell, cdr Cell) Cell {
	return &ConsCell{car, cdr, evlisSequential}
}
