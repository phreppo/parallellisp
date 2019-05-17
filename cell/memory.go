package cell

// MakeInt supplies a Int Cell. Blocking: use only in sequential code
func MakeInt(i int) Cell {
	return &IntCell{i}
}

// MakeString supplies a String Cell. Blocking: use only in sequential code
func MakeString(s string) Cell {
	return &StringCell{s}
}

// MakeSymbol supplies a Symbol Cell. Blocking: use only in sequential code
func MakeSymbol(s string) Cell {
	if isBuiltin, builtinSymbol := Lisp.IsBuiltinSymbol(s); isBuiltin {
		return builtinSymbol
	}
	return &SymbolCell{s}
}

// MakeCons supplies a Cons Cell. Blocking: use only in sequential code
func MakeCons(car Cell, cdr Cell) Cell {
	return &ConsCell{car, cdr, evlisSequential}
}
