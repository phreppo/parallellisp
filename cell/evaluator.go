package cell

import (
	"fmt"
)

var EvalService = startEvalService()

type EvalError struct {
	Err string
}

func newEvalError(e string) *EvalError {
	r := new(EvalError)
	r.Err = e
	return r
}

type EvalResult struct {
	Cell Cell
	Err  error
}

func newEvalResult(c Cell, e error) *EvalResult {
	r := new(EvalResult)
	r.Cell = c
	r.Err = e
	return r
}

type EvalRequest struct {
	Cell      Cell
	Env       *EnvironmentEntry
	ReplyChan chan *EvalResult
}

func NewEvalRequest(c Cell, env *EnvironmentEntry, replChan chan *EvalResult) *EvalRequest {
	r := new(EvalRequest)
	r.Cell = c
	r.Env = env
	r.ReplyChan = replChan
	return r
}

func (e EvalError) Error() string {
	return e.Err
}

func startEvalService() chan *EvalRequest {
	service := make(chan *EvalRequest)
	go server(service)
	return service
}

func server(service <-chan *EvalRequest) {
	for {
		req := <-service
		go eval(req)
	}
}

func eval(req *EvalRequest) {
	replyChan := req.ReplyChan
	toEval := req.Cell
	env := req.Env
	if toEval == nil {
		replyChan <- newEvalResult(nil, nil)
	}
	switch c := toEval.(type) {
	case *IntCell:
		replyChan <- newEvalResult(c, nil)
	case *StringCell:
		replyChan <- newEvalResult(c, nil)
	case *SymbolCell:
		replyChan <- newEvalResult(c, nil)
	case *ConsCell:
		switch car := c.Car.(type) {
		case *BuiltinMacroCell:
			replyChan <- car.Macro(c.Cdr, env)
		default:
			argsResult := evlis(c.Cdr)
			if argsResult.Err != nil {
				replyChan <- newEvalResult(nil, argsResult.Err)
			} else {
				replyChan <- newEvalResult(argsResult.Cell, nil)
			}
		}
	default:
		error := newEvalError("[eval] Unknown cell type: " + fmt.Sprintf("%v", toEval))
		replyChan <- newEvalResult(nil, error)
	}
}

func evlis(args Cell) *EvalResult {
	unvaluedArgs := extractArgs(args)

	if len(*unvaluedArgs) == 0 {
		return newEvalResult(nil, nil)
	}

	var replyChans []chan *EvalResult
	n := len(*unvaluedArgs)
	for i := 0; i < n; i++ {
		newChan := make(chan *EvalResult)
		replyChans = append(replyChans, newChan)
		go eval(NewEvalRequest((*unvaluedArgs)[i], EmptyEnv(), newChan))
	}

	var top Cell
	var actCons Cell
	for i := 0; i < n; i++ {
		evaluedArg := <-replyChans[i]
		if evaluedArg.Err != nil {
			return newEvalResult(nil, evaluedArg.Err)
		}
		if top == nil {
			top = MakeCons(evaluedArg.Cell, nil)
			actCons = top
		} else {
			tmp := MakeCons(evaluedArg.Cell, nil)
			switch actConsCasted := actCons.(type) {
			case *ConsCell:
				actConsCasted.Cdr = tmp
				actCons = actConsCasted.Cdr
			}
		}
	}

	return newEvalResult(top, nil)
}

func extractArgs(args Cell) *[]Cell {

	act := args
	var argsArray = new([]Cell)
	if args == nil {
		return argsArray
	}
	for act != nil {
		switch actCons := act.(type) {
		case *ConsCell:
			*argsArray = append(*argsArray, actCons.Car)
			act = actCons.Cdr
		}
	}
	return argsArray
}
