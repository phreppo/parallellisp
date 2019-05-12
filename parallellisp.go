package main

import (
	"fmt"
	"math/rand"
	"time"

	. "github.com/parof/parallellisp/cell"
	. "github.com/parof/parallellisp/parser"
)

func main() {
	source := " (car (1 2) 2 3)  "
	m := NewMemory()
	sexpr, _, _ := Parse(source, m)
	fmt.Println(sexpr)

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

	// for i := 0; i < 10; i++ {
	// 	go makeAndPrintCell(i, m)
	// }
	// time.Sleep(time.Duration(10) * time.Second)

	// switch c := intcell.(type) {
	// case *IntCell:
	// 	fmt.Println(c)
	// }
}

func makeAndPrintCell(i int, m *Memory) {
	ans := make(chan Cell)

	m.MakeInt <- IntRequest{i, ans}

	r := rand.Intn(1000)
	time.Sleep(time.Duration(r) * time.Millisecond)
	intCell := <-ans

	fmt.Print(i)
	fmt.Print(" > ")
	fmt.Println(intCell)
}
