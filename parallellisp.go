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
		fmt.Print(aurora.BrightBlue("λ "))
		source, _ := reader.ReadString('\n')
		sexpr, err := Parse(source)
		if err != nil {
			printError(err)
		} else {
			if ok, err := SemanticAnalysis(sexpr); !ok {
				printError(err)
			} else {
				evalAndPrint(sexpr)
			}
		}
	}
}

func printError(e error) {
	fmt.Println(" ", aurora.BrightRed(e), aurora.BrightRed("✗"))
}

func evalAndPrint(sexpr Cell) {
	ansChan := make(chan EvalResult)
	EvalService <- NewEvalRequest(sexpr, SimpleEnv(), ansChan)
	result := <-ansChan
	if result.Err != nil {
		printError(result.Err)
	} else {
		fmt.Println(" ", result.Cell, aurora.BrightGreen("✓"))
	}
}

func main() {
	// runtime.GOMAXPROCS(1)
	Init()
	repl()
}
