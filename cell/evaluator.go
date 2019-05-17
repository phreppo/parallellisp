package cell

import (
	"fmt"
)

var EvalService = startEvalService()

type EvalError struct {
	Err string
}

func newEvalError(e string) EvalError {
	r := EvalError{
		Err: e,
	}
	return r
}

type EvalResult struct {
	Cell Cell
	Err  error
}

func newEvalResult(c Cell, e error) EvalResult {
	r := EvalResult{
		Cell: c,
		Err:  e,
	}
	return r
}

type EvalRequest struct {
	Cell      Cell
	Env       *EnvironmentEntry
	ReplyChan chan EvalResult
}

func NewEvalRequest(c Cell, env *EnvironmentEntry, replChan chan EvalResult) EvalRequest {
	r := EvalRequest{
		Cell:      c,
		Env:       env,
		ReplyChan: replChan,
	}
	return r
}

func (e EvalError) Error() string {
	return e.Err
}

func startEvalService() chan EvalRequest {
	service := make(chan EvalRequest)
	go server(service)
	return service
}

func server(service <-chan EvalRequest) {
	for {
		req := <-service
		go serve(req)
	}
}

func serve(req EvalRequest) {
	req.ReplyChan <- eval(req.Cell, req.Env)
}

func eval(toEval Cell, env *EnvironmentEntry) EvalResult {
	if toEval == nil {
		return newEvalResult(nil, nil)
	}
	switch c := toEval.(type) {
	case *IntCell:
		return newEvalResult(c, nil)
	case *StringCell:
		return newEvalResult(c, nil)
	case *SymbolCell:
		return assoc(c, env)
	case *ConsCell:
		switch car := c.Car.(type) {
		case *BuiltinMacroCell:
			return car.Macro(c.Cdr, env)
		default:
			argsResult := evlis(c.Cdr, env)
			if argsResult.Err != nil {
				return newEvalResult(nil, argsResult.Err)
			} else {
				return apply(car, argsResult.Cell, env)
			}
		}
	default:
		error := newEvalError("[eval] Unknown cell type: " + fmt.Sprintf("%v", toEval))
		return newEvalResult(nil, error)
	}
}

func evlis(args Cell, env *EnvironmentEntry) EvalResult {
	unvaluedArgs := extractCars(args)

	if len(unvaluedArgs) == 0 {
		return newEvalResult(nil, nil)
	}

	var replyChans []chan EvalResult
	n := len(unvaluedArgs)
	for i := 0; i < n-1; i++ {
		newChan := make(chan EvalResult)
		replyChans = append(replyChans, newChan)
		go serve(NewEvalRequest(unvaluedArgs[i], env, newChan)) // TODO: empty env!!
	}

	lastArgResult := eval(unvaluedArgs[n-1], env)
	if lastArgResult.Err != nil {
		return lastArgResult
	}

	var top Cell
	var actCons Cell
	for i := 0; i < n; i++ {
		if i == n-1 {
			appendCellToArgs(&top, &actCons, &(lastArgResult.Cell))
		} else {
			evaluedArg := <-replyChans[i]
			if evaluedArg.Err != nil {
				return newEvalResult(nil, evaluedArg.Err)
			}
			appendCellToArgs(&top, &actCons, &(evaluedArg.Cell))
		}
	}

	return newEvalResult(top, nil)
}

func extractCars(args Cell) []Cell {
	act := args
	var argsArray []Cell
	if args == nil {
		return argsArray
	}
	var actCons *ConsCell
	for act != nil {
		actCons = act.(*ConsCell)
		argsArray = append(argsArray, actCons.Car)
		act = actCons.Cdr
	}
	return argsArray
}

// appends to append after actCell, maybe initializing top. Has side effects
func appendCellToArgs(top, actCell, toAppend *Cell) {
	if *top == nil {
		*top = MakeCons((*toAppend), nil)
		*actCell = *top
	} else {
		tmp := MakeCons((*toAppend), nil)
		actConsCasted := (*actCell).(*ConsCell)
		actConsCasted.Cdr = tmp
		*actCell = actConsCasted.Cdr
	}
}

func apply(function Cell, args Cell, env *EnvironmentEntry) EvalResult {
	switch functionCasted := function.(type) {
	case *BuiltinLambdaCell:
		return functionCasted.Lambda(args, env)
	case *ConsCell:
		newEnv, err := pairlis(unsafeCadr(function), args, env)
		if err != nil {
			return newEvalResult(nil, err)
		}
		return eval(unsafeCaddr(function), newEnv)
	case *SymbolCell:
		evaluedFunction := eval(function, env)
		if evaluedFunction.Err != nil {
			return newEvalResult(nil, evaluedFunction.Err)
		}
		return apply(evaluedFunction.Cell, args, env)
	default:
		return newEvalResult(nil, newEvalError("[apply] for now, only builtin lambads can be applied"))
	}
}

func unsafeCar(c Cell) Cell {
	return (c.(*ConsCell)).Car
}

func unsafeCdr(c Cell) Cell {
	return (c.(*ConsCell)).Cdr
}

func unsafeCaar(c Cell) Cell {
	return unsafeCar(unsafeCar(c))
}

func unsafeCadr(c Cell) Cell {
	return unsafeCar(unsafeCdr(c.(*ConsCell)))
}

func unsafeCaddr(c Cell) Cell {
	return unsafeCadr(unsafeCdr(c.(*ConsCell)))
}

// Pre: symbol != nil, env. pair != nil
func assoc(symbol *SymbolCell, env *EnvironmentEntry) EvalResult {
	if res, isInGlobalEnv := GlobalEnv[symbol.Sym]; isInGlobalEnv {
		return newEvalResult(res, nil)
	}
	if env == nil {
		return newEvalResult(nil, newEvalError("[assoc] symbol "+symbol.Sym+" not in env"))
	}
	act := env
	for act != nil {
		if *(act.Pair.Symbol) == *symbol {
			return newEvalResult(env.Pair.Value, nil)
		}
		act = act.Next
	}
	return newEvalResult(nil, newEvalError("[assoc] symbol "+symbol.Sym+" not in env"))
}

func pairlis(formalParameters, actualParameters Cell, oldEnv *EnvironmentEntry) (*EnvironmentEntry, error) {
	formalParametersSlice := extractCars(formalParameters)
	actualParametersSlice := extractCars(actualParameters)
	if len(actualParametersSlice) != len(formalParametersSlice) {
		return nil, newEvalError("[pairlis] mismatching number of formal and actual parameters ")
	}
	newEntry := oldEnv
	for i, formal := range formalParametersSlice {
		newEntry = NewEnvironmentEntry(formal.(*SymbolCell), actualParametersSlice[i], newEntry)
	}
	return newEntry, nil
}
