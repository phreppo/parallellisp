package parser

import (
	"fmt"

	. "github.com/parof/parallellisp/cell"
)

// Parse returns the result, if there were errors parsing and eventually one error message
func Parse(source string, m *Memory) (Cell, bool, string) {
	tokens := tokenize(source)
	if len(tokens) == 1 && tokens[1].typ == tokNone {
		return nil, true, "empty source"
	}
	return ricParse(tokens, m)
}

func ricParse(tokens []token, m *Memory) (Cell, bool, string) {
	actualToken := extractNextToken(tokens)
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
		return buildCons(tokens, m, ansChan)
	default:
		return nil, true, ("parse error near token " + fmt.Sprintf("%v", actualToken))
	}
	return nil, true, ""
}

func extractNextToken(tokens []token) token {
	if len(tokens) == 0 {
		return token{typ: tokNone}
	}
	tok := tokens[0]
	tokens = append(tokens[:0], tokens[1:]...)
	return tok
}

func readNextToken(tokens []token) token {
	return tokens[0]
}

func buildCons(tokens []token, m *Memory, ansChan chan Cell) (Cell, bool, string) {
	left, errorLeft, errorText := ricParse(tokens, m)
	if errorLeft {
		return nil, true, errorText
	}
	top := MakeCons(left, nil, m, ansChan)
	actCons := top
	for {
		actualToken := readNextToken(tokens)
		if actualToken.typ == tokDot {
			extractNextToken(tokens) // extract the dot
			// last element
			right, rightError, errorText := ricParse(tokens, m)

			if rightError {
				return nil, true, errorText
			}
			closePar := extractNextToken(tokens)

			if closePar.typ != tokClose {
				return nil, true, ("parenthesis not closed near" + fmt.Sprintf("%v", right))
			}
			switch cons := actCons.(type) {
			case *ConsCell:
				(*cons).Cdr = right
			}
			return top, false, ""
		} else {
			// TODO SEI QUI
			// not last element
			// right, rightError, errorText := ricParse(tokens, m)
			// if rightError {
			// 	return nil, true, errorText
			// } else {
			// 	closePar := extractNextToken(tokens)
			// 	if closePar.typ != tokClose {
			// 		return nil, true, ("parenthesis not closed near" + fmt.Sprintf("%v", right))
			// 	}
			// 	top = MakeCons(left, right, m, ansChan)
			// 	return top, false, ""
			right, rightError, errorText := ricParse(tokens, m)
			if rightError {
				return nil, true, errorText
			}
			tmp := MakeCons(right, nil, m, ansChan)
			if top == actCons {
				// must init the top
				switch cons := top.(type) {
				case *ConsCell:
					(*cons).Cdr = tmp
				}
			}
			switch cons := actCons.(type) {
			case *ConsCell:
				(*cons).Cdr = tmp
				actCons = (*cons).Cdr
			}

			maybeClosePar := readNextToken(tokens)

			if maybeClosePar.typ == tokClose {
				return top, false, ""
			}
		}
	}
	return left, false, ""
}
