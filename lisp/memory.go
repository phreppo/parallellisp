package lisp

func makeInt(i int) Cell {
	return &intCell{i}
}

func makeString(s string) Cell {
	return &stringCell{s}
}

func makeSymbol(s string) Cell {
	if isBuiltin, builtinSymbol := lisp.isBuiltinSymbol(s); isBuiltin {
		return builtinSymbol
	}
	return &symbolCell{s}
}

func makeCons(car Cell, cdr Cell) Cell {
	return &consCell{car, cdr, evlisSequential, false}
}
