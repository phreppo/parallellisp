package lisp

import "fmt"

// SemanticAnalysis performs the semantic analysis of one parsed sexpression and returns
// true if it is correct. This is the case if it does not contains side effects in
// nested sexpressions.
func SemanticAnalysis(c Cell) (bool, error) {
	switch cell := c.(type) {
	case *consCell:
		if err := containsSideEffects(cell.Cdr); err != nil {
			return false, &SemanticError{fmt.Sprintf("[semanalysis] expression %v contains side effects: %v", cell, err)}
		}
		return true, nil
	default:
		return true, nil
	}
}

func containsSideEffects(c Cell) error {
	if c == nil {
		return nil
	}
	switch cell := c.(type) {
	case *consCell:
		if sideLeft := containsSideEffects(cell.Car); sideLeft != nil {
			return sideLeft
		} else if sideRight := containsSideEffects(cell.Cdr); sideRight != nil {
			return sideRight
		}
		return nil
	default:
		if lisp.hasSideEffect(cell) {
			return &SemanticError{fmt.Sprintf("%v", cell)}
		}
		return nil
	}
}
