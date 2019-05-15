package eval

import (
	"fmt"

	. "github.com/parof/parallellisp/cell"
)

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
	ReplyChan chan *EvalResult
}

func NewEvalRequest(c Cell, replChan chan *EvalResult) *EvalRequest {
	r := new(EvalRequest)
	r.Cell = c
	r.ReplyChan = replChan
	return r
}

func (e EvalError) Error() string {
	return e.Err
}

func StartEvaluator() chan *EvalRequest {
	service := make(chan *EvalRequest)
	go server(service)
	return service
}

func server(service <-chan *EvalRequest) {
	for {
		req := <-service
		go serve(req)
	}
}

func serve(req *EvalRequest) {
	ansChan := req.ReplyChan
	toEval := req.Cell
	if toEval == nil {
		ansChan <- &EvalResult{nil, nil}
	}
	switch c := toEval.(type) {
	case *IntCell:
		ansChan <- newEvalResult(c, nil)
	case *StringCell:
		ansChan <- newEvalResult(c, nil)
	case *SymbolCell:
		ansChan <- newEvalResult(c, nil)
	default:
		error := newEvalError("unknown cell type: " + fmt.Sprintf("%v", toEval))
		ansChan <- newEvalResult(nil, error)
	}
}
