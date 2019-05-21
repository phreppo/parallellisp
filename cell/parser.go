package cell

import (
	"fmt"
)

var tokensIndex = 0

// Parse returns the result, if there were errors parsing and eventually one error message
func Parse(source string) (Cell, error) {
	tokensIndex = 0
	tokens := tokenize(source)
	if len(tokens) == 1 && tokens[0].typ == tokNone {
		return nil, ParseError{"empty source"}
	}
	return ricParse(tokens)
}

func ricParse(tokens []token) (Cell, error) {
	actualToken, err := extractNextToken(tokens)
	if err != nil {
		return nil, err
	}

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
		return buildCons(tokens, tokOpen, tokClose)
	case tokOpenParallel:
		cons, err := buildCons(tokens, tokCloseParallel, tokCloseParallel)
		if err != nil {
			return nil, err
		}
		(*(cons.(*ConsCell))).Evlis = evlisParallel
		return cons, nil
	default:
		return nil, ParseError{"parse error near token " + fmt.Sprintf("%v", actualToken)}
	}
}

func extractNextToken(tokens []token) (token, error) {
	if !enoughTokens(tokens) {
		return token{typ: tokNone}, ParseError{"tokens ended"}
	}
	tok := tokens[tokensIndex]
	tokensIndex++
	return tok, nil
}

func readNextToken(tokens []token) (token, error) {
	if !enoughTokens(tokens) {
		return token{typ: tokNone}, ParseError{"parenthesis not closed"}
	}
	return tokens[tokensIndex], nil
}

func enoughTokens(tokens []token) bool {
	return tokensIndex < len(tokens)
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

func buildCons(tokens []token, openParToken, closeParToken tokenType) (Cell, error) {
	nextToken, err := readNextToken(tokens)
	if err != nil {
		return nil, err
	}
	if nextToken.typ == closeParToken {
		return nil, nil
	}
	left, err := ricParse(tokens)
	if err != nil {
		return nil, err
	}
	top := MakeCons(left, nil)
	actCons := top

	nextToken, err = readNextToken(tokens)
	if err != nil {
		return nil, err
	}
	if nextToken.typ == closeParToken {
		extractNextToken(tokens)
		return top, nil
	}

	for {
		actualToken, err := readNextToken(tokens)
		if err != nil {
			return nil, err
		}
		if actualToken.typ == tokDot {
			extractNextToken(tokens) // extract the dot
			// last element
			right, err := ricParse(tokens)

			if err != nil {
				return nil, err
			}
			closePar, err := extractNextToken(tokens)
			if err != nil {
				return nil, err
			}

			if closePar.typ != closeParToken {
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

			maybeClosePar, err := readNextToken(tokens)
			if err != nil {
				return nil, err
			}

			if maybeClosePar.typ == closeParToken {
				extractNextToken(tokens)
				return top, nil
			}
		}
	}
}

type ParseError struct {
	err string
}

func (e ParseError) Error() string {
	return e.err
}
