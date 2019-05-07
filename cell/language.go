package cell

import "fmt"

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
