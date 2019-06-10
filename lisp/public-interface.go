package lisp

import (
	"bufio"
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

// Cell is the generic interface for every value in Lisp
type Cell interface {
	Eq(Cell) bool
}

// Eval evaluates one sexpression in the empty environment
func Eval(c Cell) EvalResult {
	return eval(c, emptyEnv())
}

// EvalResult result contains the result of one evaluation.
// It can be either a correct result or an error
type EvalResult struct {
	Cell Cell
	Err  error
}

// EvalError represents the error of a computation
type EvalError struct {
	Err string
}

func (e EvalError) Error() string {
	return e.Err
}

// ParseError represents one error found during the paring
type ParseError struct {
	err string
}

func (e ParseError) Error() string {
	return e.err
}

// SemanticError represents one error found during the semantic analysis
type SemanticError struct {
	errorString string
}

func (e *SemanticError) Error() string {
	return e.errorString
}

// Init initializes the needed variables. Must be called before using any lisp structure
func Init() {
	initLanguage()
	initGlobalEnv()
}

// Repl performs the read-eval-printline loop
func Repl() {
	Init()
	reader := bufio.NewReader(os.Stdin)
	for {
		// Read
		fmt.Print(aurora.BrightBlue("≃ "))
		source, _ := reader.ReadString('\n')
		if source == "\n" {
			fmt.Println("  Bye!")
			return
		}
		// Parse
		sexpr, err := Parse(source)
		if err != nil {
			printError(err)
		} else {
			// Semantic Analysis
			if ok, err := SemanticAnalysis(sexpr); !ok {
				printError(err)
			} else {
				// Eval
				result := Eval(sexpr)
				if result.Err != nil {
					printError(result.Err)
				} else {
					fmt.Println(" ", result.Cell, aurora.BrightGreen("✓"))
				}
			}
		}
	}
}

func printError(e error) {
	fmt.Println(" ", aurora.BrightRed(e), aurora.BrightRed("✗"))
}
