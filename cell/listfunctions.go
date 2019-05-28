package cell

func car(c Cell) Cell {
	return (c.(*ConsCell)).Car
}

func cdr(c Cell) Cell {
	return (c.(*ConsCell)).Cdr
}

func caar(c Cell) Cell {
	return car(car(c))
}

func cadr(c Cell) Cell {
	return car(cdr(c.(*ConsCell)))
}

func cdar(c Cell) Cell {
	return cdr(car(c.(*ConsCell)))
}

func caddr(c Cell) Cell {
	return cadr(cdr(c.(*ConsCell)))
}

func cadar(c Cell) Cell {
	return cadr(car(c.(*ConsCell)))
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
	case *IntCell:
		return cell
	case *StringCell:
		return cell
	case *BuiltinLambdaCell:
		return cell
	case *BuiltinMacroCell:
		return cell
	case *SymbolCell:
		if symbolIsInEnv(cell, env) {
			return assoc(cell, env).Cell
		}
		return cell
	case *ConsCell:
		return makeCons(copyAndSubstituteSymbols(cell.Car, env), copyAndSubstituteSymbols(cell.Cdr, env))
	default:
		return nil
	}
}
