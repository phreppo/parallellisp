package cell

import "fmt"

type SemanticError struct {
	errorString string
}

func (e *SemanticError) Error() string {
	return e.errorString
}

func SemanticAnalysis(c Cell) (bool, error) {
	switch cell := c.(type) {
	case *ConsCell:
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
	case *ConsCell:
		if sideLeft := containsSideEffects(cell.Car); sideLeft != nil {
			return sideLeft
		} else if sideRight := containsSideEffects(cell.Cdr); sideRight != nil {
			return sideRight
		}
		return nil
	default:
		if Lisp.HasSideEffect(cell) {
			return &SemanticError{fmt.Sprintf("%v", cell)}
		}
		return nil
	}
}
