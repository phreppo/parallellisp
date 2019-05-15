package cell

import "fmt"

// Lisp is the global variable for the language
var Lisp = newLanguage()

type language struct {
	builtinLambdas map[string]BuiltinLambdaCell
	builtinMacros  map[string]BuiltinMacroCell
	trueSymbol     SymbolCell
}

func newLanguage() *language {
	lisp := language{
		builtinLambdas: BuiltinLambdas,
		builtinMacros:  BuiltinMacros,
		trueSymbol:     TrueSymbol,
	}
	return &lisp
}

func (lang *language) IsBuiltinSymbol(s string) (bool, Cell) {
	isBuiltinLambda, builtinLambda := lang.IsBuiltinLambdaSymbol(s)
	if isBuiltinLambda {
		return true, builtinLambda
	}
	isBuiltinMacro, builtinMacro := lang.IsBuiltinMacroSymbol(s)
	if isBuiltinMacro {
		return true, builtinMacro
	}
	if s == "T" || s == "t" {
		return true, &(lang.trueSymbol)
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

var BuiltinLambdas = map[string]BuiltinLambdaCell{
	"car": BuiltinLambdaCell{
		Sym:    "car",
		Lambda: func() { fmt.Println("sono car!") }},
	"cdr": BuiltinLambdaCell{
		Sym:    "cdr",
		Lambda: func() {}},
	"cons": BuiltinLambdaCell{
		Sym:    "cons",
		Lambda: func() {}},
	"eq": BuiltinLambdaCell{
		Sym:    "eq",
		Lambda: func() {}},
	"atom": BuiltinLambdaCell{
		Sym:    "atom",
		Lambda: func() {}},
	"lambda": BuiltinLambdaCell{
		Sym:    "lambda",
		Lambda: func() {}},

	// "label",
}

var BuiltinMacros = map[string]BuiltinMacroCell{
	"quote": BuiltinMacroCell{
		Sym:   "quote",
		Macro: func() {}},
	"cond": BuiltinMacroCell{
		Sym:   "cond",
		Macro: func() {}},
}

var TrueSymbol = SymbolCell{Sym: "t"}
