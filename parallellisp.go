package main

import (
	"fmt"

	. "github.com/parof/parallellisp/structure"
)

func main() {
	fmt.Println("hey man")
	c := IntCell{Val: 1}
	fmt.Println(c.IsAtom())
}
