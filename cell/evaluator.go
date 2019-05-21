package cell

import (
	"fmt"

	"github.com/parof/parallellisp/scheduler"
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

func newEvalPositiveResult(c Cell) EvalResult {
	return newEvalResult(c, nil)
}

func newEvalErrorResult(e error) EvalResult {
	return newEvalResult(nil, e)
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
		scheduler.AddJob()
		go serve(req)
	}
}

func serve(req EvalRequest) {
	result := eval(req.Cell, req.Env)
	scheduler.JobEnded()
	req.ReplyChan <- result
}

func eval(toEval Cell, env *EnvironmentEntry) EvalResult {
	if toEval == nil {
		return newEvalPositiveResult(nil)
	}
	switch c := toEval.(type) {
	case *IntCell:
		return newEvalPositiveResult(c)
	case *StringCell:
		return newEvalPositiveResult(c)
	case *SymbolCell:
		return assoc(c, env)
	case *ConsCell:
		switch car := c.Car.(type) {
		case *BuiltinMacroCell:
			return car.Macro(c.Cdr, env)
		default:
			var argsResult EvalResult
			if scheduler.ShouldParallelize() {
				argsResult = c.Evlis(c.Cdr, env)
			} else {
				argsResult = evlisSequential(c.Cdr, env)
			}
			if argsResult.Err != nil {
				return newEvalErrorResult(argsResult.Err)
			} else {
				return apply(car, argsResult.Cell, env)
			}
		}
	default:
		return newEvalErrorResult(newEvalError("[eval] Unknown cell type: " + fmt.Sprintf("%v", toEval)))
	}
}

// todo: qui dovresti mettere prima di tutto l accodamento del numero di parametri, perche senno dopo il numero di goroutine va fuori controllo
func evlisParallel(args Cell, env *EnvironmentEntry) EvalResult {
	n := getNumberOfArgs(args)
	scheduler.AddJobs(int32(n - 1))

	if n == 0 {
		return newEvalPositiveResult(nil)
	}

	var replyChans []chan EvalResult
	act := args
	for act != nil && cdr(act) != nil {
		newChan := make(chan EvalResult)
		replyChans = append(replyChans, newChan)
		go serve(NewEvalRequest(car(act), env, newChan))
		act = cdr(act)
	}

	lastArgResult := eval(car(act), env)
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
			scheduler.JobEnded()
			if evaluedArg.Err != nil {
				return newEvalErrorResult(evaluedArg.Err)
			}
			appendCellToArgs(&top, &actCons, &(evaluedArg.Cell))
		}
	}

	return newEvalPositiveResult(top)
}

func getNumberOfArgs(c Cell) int {
	count := 0
	act := c
	actNotNil := (act != nil)
	for actNotNil {
		count++
		if cdr(act) == nil {
			actNotNil = false
		}
		act = cdr(act)
	}
	return count
}

func evlisSequential(args Cell, env *EnvironmentEntry) EvalResult {
	actArg := args
	var top Cell
	var actCons Cell
	var evaluedArg EvalResult

	for actArg != nil {
		evaluedArg = eval(actArg.(*ConsCell).Car, env)
		if evaluedArg.Err != nil {
			return evaluedArg
		}
		appendCellToArgs(&top, &actCons, &(evaluedArg.Cell))
		actArg = (actArg.(*ConsCell)).Cdr
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
		if Lisp.IsLambdaSymbol(functionCasted.Car) {
			newEnv, err := pairlis(cadr(function), args, env)
			if err != nil {
				return newEvalErrorResult(err)
			}
			return eval(caddr(function), newEnv)
		} else {
			// label check here
			return newEvalErrorResult(newEvalError("[apply] trying to apply a non-lambda"))
		}
	case *SymbolCell:
		evaluedFunction := eval(function, env)
		if evaluedFunction.Err != nil {
			return newEvalErrorResult(evaluedFunction.Err)
		}
		return apply(evaluedFunction.Cell, args, env)
	default:
		return newEvalErrorResult(newEvalError("[apply] trying to apply non-builtin, non-lambda, non-symbol"))
	}
}

func car(c Cell) Cell {
	return (c.(*ConsCell)).Car
}

func cdr(c Cell) Cell {
	return (c.(*ConsCell)).Cdr
}

func caar(c Cell) Cell {
	return car(car(c))
}

func cadr(c Cell) Cell {
	return car(cdr(c.(*ConsCell)))
}

func caddr(c Cell) Cell {
	return cadr(cdr(c.(*ConsCell)))
}

// Pre: symbol != nil, env. pair != nil
func assoc(symbol *SymbolCell, env *EnvironmentEntry) EvalResult {
	if res, isInGlobalEnv := GlobalEnv[symbol.Sym]; isInGlobalEnv {
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

func pairlis(formalParameters, actualParameters Cell, oldEnv *EnvironmentEntry) (*EnvironmentEntry, error) {
	actFormal := formalParameters
	actActual := actualParameters
	newEntry := oldEnv
	for actFormal != nil {
		if actActual == nil {
			return nil, newEvalError("[parilis] not enough actual parameters")
		}
		newEntry = NewEnvironmentEntry((car(actFormal)).(*SymbolCell), car(actActual), newEntry)
		actFormal = (actFormal.(*ConsCell)).Cdr
		actActual = (actActual.(*ConsCell)).Cdr
	}
	return newEntry, nil
}
