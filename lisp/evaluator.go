package lisp

import "fmt"

func eval(toEval Cell, env *environmentEntry) EvalResult {
	if toEval == nil {
		return newEvalPositiveResult(nil)
	}
	switch c := toEval.(type) {
	case *intCell:
		return newEvalPositiveResult(c)
	case *stringCell:
		return newEvalPositiveResult(c)
	case *symbolCell:
		return assoc(c, env)
	case *consCell:
		switch car := c.Car.(type) {
		case *builtinMacroCell:
			return car.Macro(c.Cdr, env)
		default:
			argsResult := c.Evlis(c.Cdr, env)
			if argsResult.Err != nil {
				return newEvalErrorResult(argsResult.Err)
			}
			return apply(car, argsResult.Cell, env)
		}
	// builtin symbols autoquote: allows higer order functions
	case *builtinMacroCell:
		return newEvalPositiveResult(c)
	case *builtinLambdaCell:
		return newEvalPositiveResult(c)
	default:
		return newEvalErrorResult(newEvalError("[eval] Unknown cell type: " + fmt.Sprintf("%v", toEval)))
	}
}

func evlisParallel(args Cell, env *environmentEntry) EvalResult {
	n := listLengt(args)

	if n == 0 {
		return newEvalPositiveResult(nil)
	}

	// send eval requests
	evaluedArgsChan := make(chan evalArgumentResult, n)
	act := args
	i := 0
	for act != nil && cdr(act) != nil {
		go evalArgumentWithChan(car(act), env, i, evaluedArgsChan)
		act = cdr(act)
		i++
	}

	// eval last arg
	lastArgResult := eval(car(act), env)
	if lastArgResult.Err != nil {
		return lastArgResult
	}

	// receive args
	valuedArgs := make([]Cell, n-1)
	var evaluedArg evalArgumentResult
	for i := 0; i < n-1; i++ {
		evaluedArg = <-evaluedArgsChan
		if evaluedArg.res.Err != nil {
			return newEvalErrorResult(evaluedArg.res.Err)
		}
		valuedArgs[evaluedArg.argIndex] = evaluedArg.res.Cell
	}

	// append in order
	var top Cell
	var actCons Cell
	for i := range valuedArgs {
		appendCellToArgs(&top, &actCons, &(valuedArgs[i]))
	}
	appendCellToArgs(&top, &actCons, &lastArgResult.Cell)

	return newEvalPositiveResult(top)
}

type evalArgumentResult struct {
	res      EvalResult
	argIndex int
}

func evalArgumentWithChan(argument Cell, env *environmentEntry, argIndex int, replyChan chan<- evalArgumentResult) {
	replyChan <- evalArgumentResult{eval(argument, env), argIndex}
}

func evlisSequential(args Cell, env *environmentEntry) EvalResult {
	actArg := args
	var top Cell
	var actCons Cell
	var evaluedArg EvalResult

	for actArg != nil {
		evaluedArg = eval(actArg.(*consCell).Car, env)
		if evaluedArg.Err != nil {
			return evaluedArg
		}
		appendCellToArgs(&top, &actCons, &(evaluedArg.Cell))
		actArg = (actArg.(*consCell)).Cdr
	}
	return newEvalResult(top, nil)
}

func apply(function Cell, args Cell, env *environmentEntry) EvalResult {
	switch functionCasted := function.(type) {
	case *builtinLambdaCell:
		return functionCasted.Lambda(args, env)
	case *consCell:
		if lisp.isLambdaSymbol(functionCasted.Car) {
			formalParameters := cadr(function)
			lambdaBody := caddr(function)
			if isClosure(formalParameters, args) {
				return newEvalPositiveResult(buildClosure(lambdaBody, formalParameters, args))
			}
			newEnv, err := pairlis(formalParameters, args, env)
			if err != nil {
				return newEvalErrorResult(err)
			}
			return eval(lambdaBody, newEnv)
		}
		// partial apply
		partiallyAppliedFunction := eval(function, env)
		if partiallyAppliedFunction.Err != nil {
			return partiallyAppliedFunction
		}
		return apply(partiallyAppliedFunction.Cell, args, env)
	case *symbolCell:
		evaluedFunction := eval(function, env)
		if evaluedFunction.Err != nil {
			return newEvalErrorResult(evaluedFunction.Err)
		}
		return apply(evaluedFunction.Cell, args, env)
	default:
		return newEvalErrorResult(newEvalError("[apply] trying to apply non-builtin, non-lambda, non-symbol"))
	}
}

func assoc(symbol *symbolCell, env *environmentEntry) EvalResult {
	if res, isInglobalEnv := globalEnv[symbol.Sym]; isInglobalEnv {
		return newEvalPositiveResult(res)
	}
	if env == nil {
		return newEvalErrorResult(newEvalError("[assoc] symbol " + symbol.Sym + " not in env"))
	}
	act := env
	for act != nil {
		if (act.Pair.Symbol.Sym) == symbol.Sym {
			return newEvalPositiveResult(act.Pair.Value)
		}
		act = act.Next
	}
	return newEvalErrorResult(newEvalError("[assoc] symbol " + symbol.Sym + " not in env"))
}

func pairlis(formalParameters, actualParameters Cell, oldEnv *environmentEntry) (*environmentEntry, error) {
	actFormal := formalParameters
	actActual := actualParameters
	newEntry := oldEnv
	for actFormal != nil {
		if actActual == nil {
			return nil, newEvalError("[parilis] not enough actual parameters")
		}
		newEntry = newEnvironmentEntry((car(actFormal)).(*symbolCell), car(actActual), newEntry)
		actFormal = (actFormal.(*consCell)).Cdr
		actActual = (actActual.(*consCell)).Cdr
	}
	return newEntry, nil
}

func newEvalError(e string) EvalError {
	r := EvalError{
		Err: e,
	}
	return r
}

func newEvalResult(c Cell, e error) EvalResult {
	r := EvalResult{
		Cell: c,
		Err:  e,
	}
	return r
}

func newEvalPositiveResult(c Cell) EvalResult {
	return newEvalResult(c, nil)
}

func newEvalErrorResult(e error) EvalResult {
	return newEvalResult(nil, e)
}

type evalRequest struct {
	Cell      Cell
	Env       *environmentEntry
	ReplyChan chan<- EvalResult
}

func newEvalRequest(c Cell, env *environmentEntry, replChan chan EvalResult) evalRequest {
	r := evalRequest{
		Cell:      c,
		Env:       env,
		ReplyChan: replChan,
	}
	return r
}
