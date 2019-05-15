package eval

import (
	"fmt"

	. "github.com/parof/parallellisp/cell"
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

func startEvalService() chan *EvalRequest {
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
	replyChan := req.ReplyChan
	toEval := req.Cell
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
	default:
		error := newEvalError("[eval] Unknown cell type: " + fmt.Sprintf("%v", toEval))
		replyChan <- newEvalResult(nil, error)
	}
}
