package main

import "github.com/parof/parallellisp/lisp"

func main() {
	// runtime.GOMAXPROCS(1)
	lisp.Repl()
}
