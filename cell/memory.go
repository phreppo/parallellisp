package cell

// makeInt supplies a Int Cell. Blocking: use only in sequential code
func makeInt(i int) Cell {
	return &intCell{i}
}

// makeString supplies a String Cell. Blocking: use only in sequential code
func makeString(s string) Cell {
	return &stringCell{s}
}

// makeSymbol supplies a Symbol Cell. Blocking: use only in sequential code
func makeSymbol(s string) Cell {
	if isBuiltin, builtinSymbol := lisp.isBuiltinSymbol(s); isBuiltin {
		return builtinSymbol
	}
	return &symbolCell{s}
}

// makeCons supplies a Cons Cell. Blocking: use only in sequential code
func makeCons(car Cell, cdr Cell) Cell {
	return &consCell{car, cdr, evlisSequential}
}
