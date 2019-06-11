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

func evalWithChan(req evalRequest) {
	req.ReplyChan <- eval(req.Cell, req.Env)
}

func evlisParallel(args Cell, env *environmentEntry) EvalResult {
	n := listLengt(args)

	if n == 0 {
		return newEvalPositiveResult(nil)
	}

	var replyChans []chan EvalResult
	act := args
	for act != nil && cdr(act) != nil {
		newChan := make(chan EvalResult)
		replyChans = append(replyChans, newChan)
		go evalWithChan(newEvalRequest(car(act), env, newChan))
		act = cdr(act)
	}

	lastArgResult := eval(car(act), env)
	if lastArgResult.Err != nil {
		return lastArgResult
	}

	var top Cell
	var actCons Cell
	var evaluedArg EvalResult
	for i := 0; i < n; i++ {
		if i < n-1 {
			evaluedArg = <-replyChans[i]
			if evaluedArg.Err != nil {
				return newEvalErrorResult(evaluedArg.Err)
			}
			appendCellToArgs(&top, &actCons, &(evaluedArg.Cell))
		} else {
			appendCellToArgs(&top, &actCons, &(lastArgResult.Cell))
		}
	}

	return newEvalPositiveResult(top)
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
