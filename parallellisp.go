package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/logrusorgru/aurora"
	. "github.com/parof/parallellisp/cell"
	"github.com/parof/parallellisp/parser"
)

func repl() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(aurora.BrightYellow("🙉 ⇝ "))
		source, _ := reader.ReadString('\n')
		sexpr, err := parser.Parse(source)
		if err != nil {
			fmt.Println("  ", aurora.Red(err))
		} else {
			evalAndPrint(sexpr)
		}
	}
}

func evalAndPrint(sexpr Cell) {
	ansChan := make(chan EvalResult)
	EvalService <- NewEvalRequest(sexpr, EmptyEnv(), ansChan)
	result := <-ansChan
	if result.Err != nil {
		fmt.Println("  😘 ", aurora.Red(result.Err))
	} else {
		fmt.Println("  🙈", result.Cell)
	}
}

func main() {
	InitMemory()
	repl()

	// tokens := Tokenize(" -33 +   \"ciao\" \"come -33 stai?\" io bene")
	// fmt.Println(tokens)

	// m := NewMemory()
	// ans := make(chan Cell)

	// i1 := MakeInt(3, m, ans)
	// fmt.Println(i1)

	// i2 := MakeInt(3, m, ans)
	// fmt.Println(i2)

	// if i1 == i2 {
	// 	fmt.Println("uguali numeri")
	// }

	// s := MakeString("ciao", m, ans)
	// fmt.Println(s)

	// sym := MakeSymbol("car", m, ans)
	// fmt.Println(sym)

	// sym1 := MakeSymbol("t", m, ans)
	// fmt.Println(sym1)

	// if sym == sym1 {
	// 	fmt.Println("we")
	// }

	// if f, ok := sym.IsBuiltinLambda(); ok {
	// 	f()
	// }

	// switch address := i.(type) {
	// case *IntCell:
	// 	fmt.Println(address)
	// default:
	// 	fmt.Println("boh")

	// }

	// c := MakeCons(i, i, m, ans)
	// fmt.Println(c)

	for i := 0; i < 10; i++ {
		go evalCellInRandomTime(i)
	}
	time.Sleep(time.Duration(10) * time.Second)

	// switch c := intcell.(type) {
	// case *IntCell:
	// 	fmt.Println(c)
	// }
}

func evalCellInRandomTime(i int) {
	ans := make(chan Cell)

	Mem.MakeInt <- IntRequest{i, ans}

	intCell := <-ans

	ansChan := make(chan EvalResult)
	EvalService <- NewEvalRequest(intCell, EmptyEnv(), ansChan)

	r := rand.Intn(1000)
	time.Sleep(time.Duration(r) * time.Millisecond)

	result := <-ansChan
	if result.Err != nil {
		fmt.Println(aurora.Red(result.Err))
	} else {
		fmt.Print(i)
		fmt.Print(" -> ")
		fmt.Println(result.Cell)
	}
}
