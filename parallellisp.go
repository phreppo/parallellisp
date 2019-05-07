package main

import (
	"fmt"
	"math/rand"
	"time"

	. "github.com/parof/parallellisp/structure"
)

func main() {
	m := NewMemory()

	m.MakeIntChan <- 3
	i := <-m.TakeIntChan
	fmt.Println(i)

	m.MakeConsChan <- CellPair{i, i}
	c := <-m.TakeConsChan
	fmt.Println(c)

	for i := 0; i < 10; i++ {
		go makAndPrintCell(i, m)
	}

	time.Sleep(time.Duration(10) * time.Second)

	// switch c := intcell.(type) {
	// case *IntCell:
	// 	fmt.Println(c)
	// }
}

func makAndPrintCell(i int, m *Memory) {
	m.MakeIntChan <- i
	r := rand.Intn(1000)
	time.Sleep(time.Duration(r) * time.Millisecond)
	intcell := <-m.TakeIntChan
	fmt.Print(i)
	fmt.Print(" > ")
	fmt.Println(intcell)
}
