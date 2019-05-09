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

func NextToken(source string) (token, string) {
	if source == "" {
		return token{typ: tokNone}, source
	}
	nextChar, index := firstChar(source)
	switch nextChar {
	case '(':
		return token{typ: tokOpen}, source[index+1:]
	case ')':
		return token{typ: tokClose}, source[index+1:]
	case '.':
		return token{typ: tokDot}, source[index+1:]
	case '\'':
		return token{typ: tokQuote}, source[index+1:]
	}
	return token{typ: tokNone}, source
}

func firstChar(str string) (byte, int) {
	for i := 0; i < len(str); i++ {
		if str[i] != ' ' {
			return str[i], i
		}
	}
	return 0, -1
}
