package cell

// Lisp is the global variable for the language
var Lisp *language

func initLanguage() {
	Lisp = newLanguage()
}

type language struct {
	builtinLambdas        map[string]BuiltinLambdaCell
	builtinMacros         map[string]BuiltinMacroCell
	builtinSpecialSymbols map[string]SymbolCell
	trueSymbol            SymbolCell
}

func (lang *language) IsBuiltinSymbol(s string) (bool, Cell) {
	if s == "NIL" || s == "nil" {
		return true, nil
	}
	isBuiltinLambda, builtinLambda := lang.IsBuiltinLambdaSymbol(s)
	if isBuiltinLambda {
		return true, builtinLambda
	}
	isBuiltinMacro, builtinMacro := lang.IsBuiltinMacroSymbol(s)
	if isBuiltinMacro {
		return true, builtinMacro
	}
	isBuiltinSpecialSymbol, builtinSpecialSymbol := lang.IsBuiltinSpecialSymbol(s)
	if isBuiltinSpecialSymbol {
		return true, builtinSpecialSymbol
	}
	return false, nil
}

// IsBuiltinLambdaSymbol returns the pointer to the concrete cell and if the symbol is a lambda.
// This is because in this manner one has not to perform two searches
func (lang *language) IsBuiltinLambdaSymbol(s string) (bool, Cell) {
	builtinLambda, isBuiltinLambda := (*lang).builtinLambdas[s]
	return isBuiltinLambda, &builtinLambda
}

// IsBuiltinMacroSymbol returns the pointer to the concrete cell and if the symbol is a macro.
// This is because in this manner one has not to perform two searches
func (lang *language) IsBuiltinMacroSymbol(s string) (bool, Cell) {
	builtinMacro, isBuiltinMacro := (*lang).builtinMacros[s]
	return isBuiltinMacro, &builtinMacro
}

// IsBuiltinSpecialSymbol returns the pointer to the concrete cell and if the symbol is a special symbol(eg: t).
// This is because in this manner one has not to perform two searches
func (lang *language) IsBuiltinSpecialSymbol(s string) (bool, Cell) {
	builtinSpecialSymbol, isBuiltinSpecialSymbol := (*lang).builtinSpecialSymbols[s]
	return isBuiltinSpecialSymbol, &builtinSpecialSymbol
}

func (lang *language) HasSideEffect(c Cell) bool {
	switch cell := c.(type) {
	case *BuiltinLambdaCell:
		return cell.Sym == "write" || cell.Sym == "load"
	case *BuiltinMacroCell:
		return cell.Sym == "defun"
	default:
		return false
	}
}

func (lang *language) GetTrueSymbol() Cell {
	return &(lang.trueSymbol)
}

func (lang *language) IsLambdaSymbol(c Cell) bool {
	switch sym := c.(type) {
	case *BuiltinMacroCell:
		return sym.Sym == "lambda"
	default:
		return false
	}
}

func newLanguage() *language {
	lisp := language{
		builtinLambdas: map[string]BuiltinLambdaCell{

			"car": BuiltinLambdaCell{
				Sym:    "car",
				Lambda: carLambda},

			"cdr": BuiltinLambdaCell{
				Sym:    "cdr",
				Lambda: cdrLambda},

			"cons": BuiltinLambdaCell{
				Sym:    "cons",
				Lambda: consLambda},

			"eq": BuiltinLambdaCell{
				Sym:    "eq",
				Lambda: eqLambda},

			"atom": BuiltinLambdaCell{
				Sym:    "atom",
				Lambda: atomLambda},

			"+": BuiltinLambdaCell{
				Sym:    "+",
				Lambda: plusLambda},

			"-": BuiltinLambdaCell{
				Sym:    "-",
				Lambda: minusLambda},

			"*": BuiltinLambdaCell{
				Sym:    "*",
				Lambda: multLambda},

			"/": BuiltinLambdaCell{
				Sym:    "/",
				Lambda: divLambda},

			">": BuiltinLambdaCell{
				Sym:    ">",
				Lambda: greaterLambda},

			">=": BuiltinLambdaCell{
				Sym:    ">=",
				Lambda: greaterEqLambda},

			"<": BuiltinLambdaCell{
				Sym:    "<",
				Lambda: lessLambda},

			"<=": BuiltinLambdaCell{
				Sym:    "<=",
				Lambda: lessEqLambda},

			"or": BuiltinLambdaCell{
				Sym:    "or",
				Lambda: orLambda},

			"and": BuiltinLambdaCell{
				Sym:    "and",
				Lambda: andLambda},

			"not": BuiltinLambdaCell{
				Sym:    "not",
				Lambda: notLambda},

			"list": BuiltinLambdaCell{
				Sym:    "list",
				Lambda: listLambda},

			"load": BuiltinLambdaCell{
				Sym:    "load",
				Lambda: loadLambda},

			"write": BuiltinLambdaCell{
				Sym:    "write",
				Lambda: writeLambda},

			// "label",
		},

		builtinMacros: map[string]BuiltinMacroCell{

			"quote": BuiltinMacroCell{
				Sym:   "quote",
				Macro: quoteMacro},

			"time": BuiltinMacroCell{
				Sym:   "time",
				Macro: timeMacro},

			"cond": BuiltinMacroCell{
				Sym:   "cond",
				Macro: condMacro},

			"lambda": BuiltinMacroCell{
				Sym:   "lambda",
				Macro: lambdaMacro},

			"defun": BuiltinMacroCell{
				Sym:   "defun",
				Macro: defunMacro},
		},

		builtinSpecialSymbols: map[string]SymbolCell{
			"t": SymbolCell{
				Sym: "t",
			},
		},

		trueSymbol: SymbolCell{Sym: "t"},
	}
	return &lisp
}
