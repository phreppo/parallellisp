package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	. "github.com/parof/parallellisp/cell"
)

func repl() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(aurora.BrightGreen("😴 "))
		source, _ := reader.ReadString('\n')
		sexpr, err := Parse(source)
		if err != nil {
			fmt.Println("  ", aurora.Red(err))
		} else {
			evalAndPrint(sexpr)
		}
	}
}

func evalAndPrint(sexpr Cell) {
	ansChan := make(chan EvalResult)
	EvalService <- NewEvalRequest(sexpr, SimpleEnv(), ansChan)
	result := <-ansChan
	if result.Err != nil {
		fmt.Println("  ", aurora.Red(result.Err), "😘")
	} else {
		fmt.Println("  ", result.Cell, aurora.BrightGreen("✓"))
	}
}

func main() {
	// runtime.GOMAXPROCS(1)
	Init()
	repl()
}
