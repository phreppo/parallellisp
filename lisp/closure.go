package lisp

func isClosure(formalParameters, actualParameters Cell) bool {
	return listLengt(formalParameters) > listLengt(actualParameters)
}

func buildClosure(lambdaBody, formalParameters, actualParameters Cell) Cell {
	// ((lambda (x y) (+ x y)) 1)
	// head
	top := makeCons(makeSymbol("lambda"), nil)
	act := top

	// unmatched parameters
	actFormal := formalParameters
	actActual := actualParameters
	found := false
	closureEnv := emptyEnv()

	for actFormal != nil && !found {
		if actActual == nil {
			// found
			found = true
		} else {
			closureEnv = newEnvironmentEntry((car(actFormal)).(*symbolCell), car(actActual), closureEnv)
			actFormal = cdr(actFormal)
			actActual = cdr(actActual)
		}
	}
	appendCellToArgs(&top, &act, &actFormal)

	closedBody := copyAndSubstituteSymbols(lambdaBody, closureEnv)
	appendCellToArgs(&top, &act, &closedBody)

	return top
}
