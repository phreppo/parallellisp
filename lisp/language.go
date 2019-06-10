package lisp

// lisp is the global variable for the language
var lisp *language

func initLanguage() {
	lisp = newLanguage()
}

type language struct {
	builtinLambdas        map[string]builtinLambdaCell
	builtinMacros         map[string]builtinMacroCell
	builtinSpecialSymbols map[string]symbolCell
	trueSymbol            symbolCell
}

func (lang *language) isBuiltinSymbol(s string) (bool, Cell) {
	if s == "NIL" || s == "nil" {
		return true, nil
	}
	isBuiltinLambda, builtinLambda := lang.isBuiltinLambdaSymbol(s)
	if isBuiltinLambda {
		return true, builtinLambda
	}
	isBuiltinMacro, builtinMacro := lang.isBuiltinMacroSymbol(s)
	if isBuiltinMacro {
		return true, builtinMacro
	}
	isBuiltinSpecialSymbol, builtinSpecialSymbol := lang.isBuiltinSpecialSymbol(s)
	if isBuiltinSpecialSymbol {
		return true, builtinSpecialSymbol
	}
	return false, nil
}

// isBuiltinLambdaSymbol returns the pointer to the concrete cell and if the symbol is a lambda.
// This is because in this manner one has not to perform two searches
func (lang *language) isBuiltinLambdaSymbol(s string) (bool, Cell) {
	builtinLambda, isBuiltinLambda := (*lang).builtinLambdas[s]
	return isBuiltinLambda, &builtinLambda
}

// isBuiltinMacroSymbol returns the pointer to the concrete cell and if the symbol is a macro.
// This is because in this manner one has not to perform two searches
func (lang *language) isBuiltinMacroSymbol(s string) (bool, Cell) {
	builtinMacro, isBuiltinMacro := (*lang).builtinMacros[s]
	return isBuiltinMacro, &builtinMacro
}

// isBuiltinSpecialSymbol returns the pointer to the concrete cell and if the symbol is a special symbol(eg: t).
// This is because in this manner one has not to perform two searches
func (lang *language) isBuiltinSpecialSymbol(s string) (bool, Cell) {
	builtinSpecialSymbol, isBuiltinSpecialSymbol := (*lang).builtinSpecialSymbols[s]
	return isBuiltinSpecialSymbol, &builtinSpecialSymbol
}

func (lang *language) hasSideEffect(c Cell) bool {
	switch cell := c.(type) {
	case *builtinLambdaCell:
		return cell.Sym == "write" || cell.Sym == "load" || cell.Sym == "set"
	case *builtinMacroCell:
		return cell.Sym == "defun" || cell.Sym == "setq"
	default:
		return false
	}
}

func (lang *language) getTrueSymbol() Cell {
	return &(lang.trueSymbol)
}

func (lang *language) isLambdaSymbol(c Cell) bool {
	switch sym := c.(type) {
	case *builtinMacroCell:
		return sym.Sym == "lambda"
	default:
		return false
	}
}

func newLanguage() *language {
	lisp := language{
		builtinLambdas: map[string]builtinLambdaCell{

			"car": builtinLambdaCell{
				Sym:    "car",
				Lambda: carLambda},

			"cdr": builtinLambdaCell{
				Sym:    "cdr",
				Lambda: cdrLambda},

			"cons": builtinLambdaCell{
				Sym:    "cons",
				Lambda: consLambda},

			"eq": builtinLambdaCell{
				Sym:    "eq",
				Lambda: eqLambda},

			"atom": builtinLambdaCell{
				Sym:    "atom",
				Lambda: atomLambda},

			"+": builtinLambdaCell{
				Sym:    "+",
				Lambda: plusLambda},

			"-": builtinLambdaCell{
				Sym:    "-",
				Lambda: minusLambda},

			"*": builtinLambdaCell{
				Sym:    "*",
				Lambda: multLambda},

			"/": builtinLambdaCell{
				Sym:    "/",
				Lambda: divLambda},

			">": builtinLambdaCell{
				Sym:    ">",
				Lambda: greaterLambda},

			">=": builtinLambdaCell{
				Sym:    ">=",
				Lambda: greaterEqLambda},

			"<": builtinLambdaCell{
				Sym:    "<",
				Lambda: lessLambda},

			"<=": builtinLambdaCell{
				Sym:    "<=",
				Lambda: lessEqLambda},

			"or": builtinLambdaCell{
				Sym:    "or",
				Lambda: orLambda},

			"and": builtinLambdaCell{
				Sym:    "and",
				Lambda: andLambda},

			"not": builtinLambdaCell{
				Sym:    "not",
				Lambda: notLambda},

			"list": builtinLambdaCell{
				Sym:    "list",
				Lambda: listLambda},

			"reverse": builtinLambdaCell{
				Sym:    "reverse",
				Lambda: reverseLambda},

			"member": builtinLambdaCell{
				Sym:    "member",
				Lambda: memberLambda},

			"nth": builtinLambdaCell{
				Sym:    "nth",
				Lambda: nthLambda},

			"length": builtinLambdaCell{
				Sym:    "length",
				Lambda: lengthLambda},

			"set": builtinLambdaCell{
				Sym:    "set",
				Lambda: setLambda},

			"load": builtinLambdaCell{
				Sym:    "load",
				Lambda: loadLambda},

			"write": builtinLambdaCell{
				Sym:    "write",
				Lambda: writeLambda},

			"integerp": builtinLambdaCell{
				Sym:    "integerp",
				Lambda: integerpLambda},

			"symbolp": builtinLambdaCell{
				Sym:    "symbolp",
				Lambda: symbolpLambda},

			"1+": builtinLambdaCell{
				Sym:    "1+",
				Lambda: onePlusLambda},

			"1-": builtinLambdaCell{
				Sym:    "1-",
				Lambda: oneMinusLambda},

			// "label",
		},

		builtinMacros: map[string]builtinMacroCell{

			"quote": builtinMacroCell{
				Sym:   "quote",
				Macro: quoteMacro},

			"time": builtinMacroCell{
				Sym:   "time",
				Macro: timeMacro},

			"cond": builtinMacroCell{
				Sym:   "cond",
				Macro: condMacro},

			"lambda": builtinMacroCell{
				Sym:   "lambda",
				Macro: lambdaMacro},

			"defun": builtinMacroCell{
				Sym:   "defun",
				Macro: defunMacro},

			"setq": builtinMacroCell{
				Sym:   "setq",
				Macro: setqMacro},

			"let": builtinMacroCell{
				Sym:   "let",
				Macro: letMacro},

			"dotimes": builtinMacroCell{
				Sym:   "dotimes",
				Macro: dotimesMacro},
		},

		builtinSpecialSymbols: map[string]symbolCell{
			"t": symbolCell{
				Sym: "t",
			},
		},

		trueSymbol: symbolCell{Sym: "t"},
	}
	return &lisp
}
