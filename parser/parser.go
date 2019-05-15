package parser

import (
	"fmt"

	. "github.com/parof/parallellisp/cell"
)

// Parse returns the result, if there were errors parsing and eventually one error message
func Parse(source string) (Cell, error) {
	tokens := tokenize(source)
	if len(tokens) == 1 && tokens[0].typ == tokNone {
		return nil, ParseError{"empty source"}
	}
	return ricParse(tokens)
}

func ricParse(tokens []token) (Cell, error) {
	actualToken := extractNextToken(tokens)

	switch actualToken.typ {
	case tokNum:
		newInt := MakeInt(actualToken.val)
		return newInt, nil
	case tokStr:
		newStr := MakeString(actualToken.str)
		return newStr, nil
	case tokSym:
		newSym := MakeSymbol(actualToken.str)
		return newSym, nil
	case tokQuote:
		return buildQuote(tokens)
	case tokOpen:
		return buildCons(tokens)
	default:
		return nil, ParseError{"parse error near token " + fmt.Sprintf("%v", actualToken)}
	}
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

func buildQuote(tokens []token) (Cell, error) {
	quoteSym := MakeSymbol("quote")
	quotedSexpression, err := ricParse(tokens)
	if err != nil {
		return nil, err
	}
	firstConsArg := MakeCons(quotedSexpression, nil)
	topCons := MakeCons(quoteSym, firstConsArg)
	return topCons, nil
}

func buildCons(tokens []token) (Cell, error) {
	left, err := ricParse(tokens)
	if err != nil {
		return nil, err
	}
	top := MakeCons(left, nil)
	actCons := top
	if readNextToken(tokens).typ == tokClose {
		return top, nil
	}
	for {
		actualToken := readNextToken(tokens)
		if actualToken.typ == tokDot {
			extractNextToken(tokens) // extract the dot
			// last element
			right, err := ricParse(tokens)

			if err != nil {
				return nil, err
			}
			closePar := extractNextToken(tokens)

			if closePar.typ != tokClose {
				return nil, ParseError{"parenthesis not closed near " + fmt.Sprintf("%v", right)}
			}
			switch cons := actCons.(type) {
			case *ConsCell:
				(*cons).Cdr = right
			}
			return top, nil
		} else {
			right, err := ricParse(tokens)
			if err != nil {
				return nil, err
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
				return top, nil
			}
		}
	}
	return top, nil
}

type ParseError struct {
	err string
}

func (e ParseError) Error() string {
	return e.err
}
