package parser

import (
	"fmt"

	. "github.com/parof/parallellisp/cell"
)

// Parse returns the result, if there were errors parsing and eventually one error message
func Parse(source string) (Cell, bool, string) {
	tokens := tokenize(source)
	if len(tokens) == 1 && tokens[0].typ == tokNone {
		return nil, true, "empty source"
	}
	return ricParse(tokens)
}

func ricParse(tokens []token) (Cell, bool, string) {
	actualToken := extractNextToken(tokens)

	switch actualToken.typ {
	case tokNone:
		break
	case tokNum:
		newInt := MakeInt(actualToken.val)
		return newInt, false, ""
	case tokStr:
		newStr := MakeString(actualToken.str)
		return newStr, false, ""
	case tokSym:
		newSym := MakeSymbol(actualToken.str)
		return newSym, false, ""
	case tokQuote:
		quoteSym := MakeSymbol("quote")
		quotedSexpression, err, errorText := ricParse(tokens)
		if err {
			return nil, true, errorText
		}
		firstConsArg := MakeCons(quotedSexpression, nil)
		topCons := MakeCons(quoteSym, firstConsArg)
		return topCons, false, ""
	case tokOpen:
		return buildCons(tokens)
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

func buildCons(tokens []token) (Cell, bool, string) {
	left, errorLeft, errorText := ricParse(tokens)
	if errorLeft {
		return nil, true, errorText
	}
	top := MakeCons(left, nil)
	actCons := top
	if readNextToken(tokens).typ == tokClose {
		return top, false, ""
	}
	for {
		actualToken := readNextToken(tokens)
		if actualToken.typ == tokDot {
			extractNextToken(tokens) // extract the dot
			// last element
			right, rightError, errorText := ricParse(tokens)

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
			right, rightError, errorText := ricParse(tokens)
			if rightError {
				return nil, true, errorText
			}
			tmp := MakeCons(right, nil)
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
				extractNextToken(tokens)
				return top, false, ""
			}
		}
	}
	return top, false, ""
}
