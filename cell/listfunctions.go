package cell

func car(c Cell) Cell {
	return (c.(*consCell)).Car
}

func cdr(c Cell) Cell {
	return (c.(*consCell)).Cdr
}

func caar(c Cell) Cell {
	return car(car(c))
}

func cadr(c Cell) Cell {
	return car(cdr(c.(*consCell)))
}

func cdar(c Cell) Cell {
	return cdr(car(c.(*consCell)))
}

func caddr(c Cell) Cell {
	return cadr(cdr(c.(*consCell)))
}

func cadar(c Cell) Cell {
	return cadr(car(c.(*consCell)))
}

func listLengt(c Cell) int {
	n := 0
	for c != nil {
		n++
		c = cdr(c)
	}
	return n
}

func eq(c1, c2 Cell) bool {
	if c1 == nil && c2 == nil {
		return true
	}
	if c1 == nil || c2 == nil {
		return false
	}
	return c1.Eq(c2)
}

func copyAndSubstituteSymbols(c Cell, env *environmentEntry) Cell {
	switch cell := c.(type) {
	case *intCell:
		return cell
	case *stringCell:
		return cell
	case *builtinLambdaCell:
		return cell
	case *builtinMacroCell:
		return cell
	case *symbolCell:
		if symbolIsInEnv(cell, env) {
			return assoc(cell, env).Cell
		}
		return cell
	case *consCell:
		return makeCons(copyAndSubstituteSymbols(cell.Car, env), copyAndSubstituteSymbols(cell.Cdr, env))
	default:
		return nil
	}
}
