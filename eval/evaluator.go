package eval

import . "github.com/parof/parallellisp/cell"

type EvalError struct {
	Err string
}

type EvalResult struct {
	Cell Cell
	Err  error
}

type EvalRequest struct {
	Cell      Cell
	ReplyChan chan *EvalResult
}

func (e *EvalError) Error() string {
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
	ansChan <- &EvalResult{nil, nil}
}
