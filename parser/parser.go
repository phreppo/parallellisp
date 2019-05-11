package parser

import (
	"fmt"

	. "github.com/parof/parallellisp/cell"
)

// Parse returns the result, if there were errors parsing and eventually one error message
func Parse(source string, m *Memory) (Cell, bool, string) {
	tokens := tokenize(source)
	i := new(int)
	*i = 0
	return ricParse(tokens, i, m)
}

func ricParse(tokens []token, i *int, m *Memory) (Cell, bool, string) {
	actualToken := nextToken(tokens, i)
	ansChan := make(chan Cell)

	switch actualToken.typ {
	case tokNone:
		break
	case tokNum:
		newInt := MakeInt(actualToken.val, m, ansChan)
		return newInt, false, ""
	case tokStr:
		newStr := MakeString(actualToken.str, m, ansChan)
		return newStr, false, ""
	case tokSym:
		newSym := MakeSymbol(actualToken.str, m, ansChan)
		return newSym, false, ""
	case tokOpen:
		return nil, false, ""
	default:
		return nil, true, ("parse error near token " + fmt.Sprintf("%v", actualToken))
	}
	return nil, true, ""
}

func nextToken(tokens []token, i *int) token {
	if *i > len(tokens) {
		return token{typ: tokNone}
	}
	tok := tokens[*i]
	(*i)++
	return tok
}
