package cell

import (
	"fmt"
	"time"
)

// Lisp is the global variable for the language
var Lisp = newLanguage()

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

func newLanguage() *language {
	lisp := language{
		builtinLambdas: map[string]BuiltinLambdaCell{
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
				Macro: unimplementedMacro},
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

func quoteMacro(args Cell, env *EnvironmentEntry) *EvalResult {
	switch cons := args.(type) {
	case *ConsCell:
		return newEvalResult(cons.Car, nil)
	default:
		return newEvalResult(nil, newEvalError("[quote] Can't quote"+fmt.Sprint(cons)))
	}
}

func timeMacro(args Cell, env *EnvironmentEntry) *EvalResult {
	now := time.Now()
	start := now.UnixNano()

	time.Sleep(time.Duration(100) * time.Millisecond)

	now = time.Now()
	afterEvalTime := now.UnixNano()
	elapsedMillis := (afterEvalTime - start) / 1000000
	fmt.Println("time:", elapsedMillis, "ms")

	return newEvalResult(nil, nil)
}

func unimplementedMacro(c Cell, env *EnvironmentEntry) *EvalResult {
	panic("unimplemented macro")
}
