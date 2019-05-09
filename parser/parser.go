package parser

import . "github.com/parof/parallellisp/cell"

type tokenType int

const (
	tokNone  tokenType = 0
	tokOpen  tokenType = 1
	tokClose tokenType = 2
	tokDot   tokenType = 3
	tokQuote tokenType = 4
	tokSym   tokenType = 5
	tokNum   tokenType = 6
	tokStr   tokenType = 7
)

type token struct {
	typ tokenType
	str string
	val int
}

func Parse(source string) (Cell, bool, string) {
	return nil, false, ""
}

func nextToken(source string) token {
	return token{}
}
