package cell

import (
	"fmt"
	"time"
)

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

func condMacro(args Cell, env *EnvironmentEntry) EvalResult {
	argsArray := extractCars(args)

	var condAndBody []Cell
	var cond Cell
	var body Cell
	var condResult EvalResult

	for _, arg := range argsArray {
		condAndBody = extractCars(arg)
		cond = (condAndBody)[0]
		body = (condAndBody)[1]
		condResult = eval(cond, env)
		if condResult.Err != nil {
			return condResult
		} else if condResult.Cell != nil {
			return eval(body, env)
		}
	}
	return newEvalResult(nil, newEvalError("[cond] none condition was verified"))
}

func quoteMacro(args Cell, env *EnvironmentEntry) EvalResult {
	switch cons := args.(type) {
	case *ConsCell:
		return newEvalResult(cons.Car, nil)
	default:
		return newEvalResult(nil, newEvalError("[quote] Can't quote"+fmt.Sprint(cons)))
	}
}

func timeMacro(args Cell, env *EnvironmentEntry) EvalResult {
	if args == nil {
		return newEvalResult(nil, newEvalError("[time] too few arguments"))
	}
	now := time.Now()
	start := now.UnixNano()

	result := eval(unsafeCar(args), env)
	if result.Err != nil {
		return result
	}

	now = time.Now()
	afterEvalTime := now.UnixNano()
	elapsedMillis := (afterEvalTime - start) / 1000000
	fmt.Println("time:", elapsedMillis, "ms")

	return result
}

func lambdaMacro(args Cell, env *EnvironmentEntry) EvalResult {
	// lambda autoquote
	return newEvalResult(MakeCons(MakeSymbol("lambda"), args), nil)
}

func defunMacro(args Cell, env *EnvironmentEntry) EvalResult {
	argsSlice := extractCars(args)
	if len(argsSlice) != 3 {
		return newEvalResult(nil, newEvalError("[defun] wrong number of arguments"))
	}
	name := argsSlice[0]
	formalParameters := argsSlice[1]
	lambdaBody := argsSlice[2]
	bodyCons := MakeCons(lambdaBody, nil)
	argsAndBodyCons := MakeCons(formalParameters, bodyCons)
	ret := MakeCons(MakeSymbol("lambda"), argsAndBodyCons)
	switch nameSymbolCell := name.(type) {
	case *SymbolCell:
		GlobalEnv[nameSymbolCell.Sym] = ret
	default:
		return newEvalResult(nil, newEvalError("[defun] the name of the lambda must be a symbol"))
	}
	return newEvalResult(ret, nil)
}

func carLambda(args Cell, env *EnvironmentEntry) EvalResult {
	switch topCons := args.(type) {
	case *ConsCell:
		switch cons := topCons.Car.(type) {
		case *ConsCell:
			return newEvalResult(cons.Car, nil)
		default:
			return newEvalResult(nil, newEvalError("[car] car applied to atom"))
		}
	default:
		return newEvalResult(nil, newEvalError("[car] not enough arguments"))
	}
}

func cdrLambda(args Cell, env *EnvironmentEntry) EvalResult {
	switch topCons := args.(type) {
	case *ConsCell:
		switch cons := topCons.Car.(type) {
		case *ConsCell:
			return newEvalResult(cons.Cdr, nil)
		default:
			return newEvalResult(nil, newEvalError("[cdr] cdr applied to atom"))
		}
	default:
		return newEvalResult(nil, newEvalError("[cdr] not enough arguments"))
	}
}

func consLambda(args Cell, env *EnvironmentEntry) EvalResult {
	switch firstCons := args.(type) {
	case *ConsCell:
		switch cons := firstCons.Cdr.(type) {
		case *ConsCell:
			result := MakeCons(firstCons.Car, cons.Car)
			return newEvalResult(result, nil)
		default:
			return newEvalResult(nil, newEvalError("[cons] not enough arguments"))
		}
	default:
		return newEvalResult(nil, newEvalError("[cons] not enough arguments"))
	}
}

func eqLambda(args Cell, env *EnvironmentEntry) EvalResult {
	switch firstArg := args.(type) {
	case *ConsCell:
		switch secondArg := firstArg.Cdr.(type) {
		case *ConsCell:
			if eq(firstArg.Car, secondArg.Car) {
				return newEvalResult(MakeSymbol("t"), nil)
			} else {
				return newEvalResult(nil, nil)
			}
		default:
			return newEvalResult(nil, newEvalError("[eq] not enough arguments"))
		}
	default:
		return newEvalResult(nil, newEvalError("[eq] not enough arguments"))
	}
}

func atomLambda(args Cell, env *EnvironmentEntry) EvalResult {
	switch firstCons := args.(type) {
	case *ConsCell:
		switch firstCons.Car.(type) {
		case *ConsCell:
			return newEvalResult(nil, nil)
		default:
			return newEvalResult(MakeSymbol("t"), nil)
		}
	default:
		return newEvalResult(nil, newEvalError("[atom] not enough arguments"))
	}
}

func plusLambda(args Cell, env *EnvironmentEntry) EvalResult {
	nums := extractCars(args)
	tot := 0
	for _, numCell := range nums {
		switch num := numCell.(type) {
		case *IntCell:
			tot += num.Val
		default:
			return newEvalResult(nil, newEvalError("[+] trying to add a non-number"))
		}
	}
	return newEvalResult(MakeInt(tot), nil)
}

func minusLambda(args Cell, env *EnvironmentEntry) EvalResult {
	nums := extractCars(args)
	tot := (nums[0].(*IntCell)).Val
	for _, numCell := range nums[1:] {
		switch num := numCell.(type) {
		case *IntCell:
			tot -= num.Val
		default:
			return newEvalResult(nil, newEvalError("[+] trying to add a non-number"))
		}
	}
	return newEvalResult(MakeInt(tot), nil)
}

func unimplementedMacro(c Cell, env *EnvironmentEntry) EvalResult {
	panic("unimplemented macro")
}

func unimplementedLambda(c Cell, env *EnvironmentEntry) EvalResult {
	panic("unimplemented lambda")
}
