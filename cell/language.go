package cell

import (
	"fmt"
	"io/ioutil"
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

func condMacro(args Cell, env *EnvironmentEntry) EvalResult {
	actBranch := args
	var condAndBody Cell
	var cond Cell
	var body Cell
	var condResult EvalResult
	for actBranch != nil {
		condAndBody = car(actBranch)
		cond = car(condAndBody)
		body = cadr(condAndBody)
		condResult = eval(cond, env)
		if condResult.Err != nil {
			return condResult
		} else if condResult.Cell != nil {
			return eval(body, env)
		}
		actBranch = cdr(actBranch)
	}
	return newEvalErrorResult(newEvalError("[cond] none condition was verified"))
}

func quoteMacro(args Cell, env *EnvironmentEntry) EvalResult {
	switch cons := args.(type) {
	case *ConsCell:
		return newEvalPositiveResult(cons.Car)
	default:
		return newEvalErrorResult(newEvalError("[quote] Can't quote" + fmt.Sprint(cons)))
	}
}

func timeMacro(args Cell, env *EnvironmentEntry) EvalResult {
	if args == nil {
		return newEvalErrorResult(newEvalError("[time] too few arguments"))
	}
	now := time.Now()
	start := now.UnixNano()

	result := eval(car(args), env)
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
	return newEvalPositiveResult(MakeCons(MakeSymbol("lambda"), args))
}

func defunMacro(args Cell, env *EnvironmentEntry) EvalResult {
	argsSlice := extractCars(args)
	if len(argsSlice) != 3 {
		return newEvalErrorResult(newEvalError("[defun] wrong number of arguments"))
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
		return newEvalErrorResult(newEvalError("[defun] the name of the lambda must be a symbol"))
	}
	return newEvalPositiveResult(ret)
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
			return newEvalPositiveResult(MakeCons(firstCons.Car, cons.Car))
		default:
			return newEvalErrorResult(newEvalError("[cons] not enough arguments"))
		}
	default:
		return newEvalErrorResult(newEvalError("[cons] not enough arguments"))
	}
}

func eqLambda(args Cell, env *EnvironmentEntry) EvalResult {
	switch firstArg := args.(type) {
	case *ConsCell:
		switch secondArg := firstArg.Cdr.(type) {
		case *ConsCell:
			if eq(firstArg.Car, secondArg.Car) {
				return newEvalPositiveResult(Lisp.GetTrueSymbol())
			} else {
				return newEvalPositiveResult(nil)
			}
		default:
			return newEvalErrorResult(newEvalError("[eq] not enough arguments"))
		}
	default:
		return newEvalErrorResult(newEvalError("[eq] not enough arguments"))
	}
}

func atomLambda(args Cell, env *EnvironmentEntry) EvalResult {
	switch firstCons := args.(type) {
	case *ConsCell:
		switch firstCons.Car.(type) {
		case *ConsCell:
			return newEvalPositiveResult(nil)
		default:
			return newEvalPositiveResult(Lisp.GetTrueSymbol())
		}
	default:
		return newEvalErrorResult(newEvalError("[atom] not enough arguments"))
	}
}

func plusLambda(args Cell, env *EnvironmentEntry) EvalResult {
	tot := 0
	act := args
	for act != nil {
		tot += (car(act).(*IntCell)).Val
		act = cdr(act)
	}
	return newEvalPositiveResult(MakeInt(tot))
}

func minusLambda(args Cell, env *EnvironmentEntry) EvalResult {
	if args == nil {
		return newEvalErrorResult(newEvalError("[-] too few arguments"))
	}
	tot := (car(args).(*IntCell)).Val
	act := cdr(args)
	for act != nil {
		tot -= (car(act).(*IntCell)).Val
		act = cdr(act)
	}
	return newEvalPositiveResult(MakeInt(tot))
}

func loadLambda(args Cell, env *EnvironmentEntry) EvalResult {
	files := extractCars(args)
	if len(files) != 1 {
		return newEvalErrorResult(newEvalError("[load] load needs exaclty one argument"))
	}
	fileName := (files[0].(*StringCell)).Str
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		return newEvalErrorResult(newEvalError("[load] error opening file " + fileName))
	}
	source := string(dat)
	sexpressions, err := ParseMultipleSexpressions(source)
	if err != nil {
		return newEvalErrorResult(err)
	}
	var lastEvalued EvalResult
	for _, sexpression := range sexpressions {
		lastEvalued = eval(sexpression, env)
		if lastEvalued.Err != nil {
			return lastEvalued
		}
	}
	return lastEvalued
}

func writeLambda(args Cell, env *EnvironmentEntry) EvalResult {
	phrases := extractCars(args)
	if len(phrases) == 0 {
		fmt.Println()
		return newEvalPositiveResult(MakeString(""))
	} else if len(phrases) > 1 {
		return newEvalErrorResult(newEvalError("[write] write needs at least one string argument"))
	}
	fmt.Println((phrases[0].(*StringCell)).Str)
	return newEvalPositiveResult(phrases[0])
}

func unimplementedMacro(c Cell, env *EnvironmentEntry) EvalResult {
	panic("unimplemented macro")
}

func unimplementedLambda(c Cell, env *EnvironmentEntry) EvalResult {
	panic("unimplemented lambda")
}
