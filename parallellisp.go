package main

import (
	"fmt"
	"math/rand"
	"time"

	. "github.com/parof/parallellisp/cell"
)

func main() {
	m := NewMemory()
	ans := make(chan Cell)

	i := MakeInt(3, m, ans)
	fmt.Println(i)

	s := MakeString("ciao", m, ans)
	fmt.Println(s)

	c := MakeCons(i, i, m, ans)
	fmt.Println(c)

	for i := 0; i < 10; i++ {
		go makeAndPrintCell(i, m)
	}
	time.Sleep(time.Duration(10) * time.Second)

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
