package cell

import (
	"fmt"
)

// Parse returns the result, if there were errors parsing and eventually one error message
func Parse(source string) (Cell, error) {
	sexpressions, err := parseMultipleSexpressions(source)
	if len(sexpressions) > 1 {
		return nil, ParseError{"[parser] too many sexpressions"}
	}
	if err != nil {
		return nil, err
	}
	return sexpressions[0], nil
}

// parseMultipleSexpressions resturns the array of parser sexpressions
func parseMultipleSexpressions(source string) ([]Cell, error) {
	tokens := tokenize(source)
	if len(tokens) == 1 && tokens[0].typ == tokNone {
		return nil, ParseError{"empty source"}
	}

	var tokensIndex = 0
	var result []Cell

	for enoughTokens(tokens, &tokensIndex) {
		actualSexpression, err := ricParse(tokens, &tokensIndex)
		if err != nil {
			return nil, err
		}
		result = append(result, actualSexpression)
	}
	return result, nil
}

func ricParse(tokens []token, tokensIndex *int) (Cell, error) {
	actualToken, err := extractNextToken(tokens, tokensIndex)
	if err != nil {
		return nil, err
	}

	switch actualToken.typ {
	case tokNum:
		newInt := makeInt(actualToken.val)
		return newInt, nil
	case tokStr:
		newStr := makeString(actualToken.str)
		return newStr, nil
	case tokSym:
		newSym := makeSymbol(actualToken.str)
		return newSym, nil
	case tokQuote:
		return buildQuote(tokens, tokensIndex)
	case tokOpen:
		return buildCons(tokens, tokOpen, tokClose, tokensIndex)
	case tokOpenParallel:
		cons, err := buildCons(tokens, tokCloseParallel, tokCloseParallel, tokensIndex)
		if err != nil {
			return nil, err
		}
		(*(cons.(*ConsCell))).Evlis = evlisParallel
		return cons, nil
	default:
		return nil, ParseError{"parse error near token " + fmt.Sprintf("%v", actualToken)}
	}
}

func extractNextToken(tokens []token, tokensIndex *int) (token, error) {
	if !enoughTokens(tokens, tokensIndex) {
		return token{typ: tokNone}, ParseError{"tokens ended"}
	}
	tok := tokens[*tokensIndex]
	(*tokensIndex)++
	return tok, nil
}

func readNextToken(tokens []token, tokensIndex *int) (token, error) {
	if !enoughTokens(tokens, tokensIndex) {
		return token{typ: tokNone}, ParseError{"parenthesis not closed"}
	}
	return tokens[(*tokensIndex)], nil
}

func enoughTokens(tokens []token, tokensIndex *int) bool {
	return (*tokensIndex) < len(tokens)
}

func buildQuote(tokens []token, tokensIndex *int) (Cell, error) {
	quoteSym := makeSymbol("quote")
	quotedSexpression, err := ricParse(tokens, tokensIndex)
	if err != nil {
		return nil, err
	}

	firstConsArg := makeCons(quotedSexpression, nil)
	topCons := makeCons(quoteSym, firstConsArg)
	return topCons, nil
}

func buildCons(tokens []token, openParToken, closeParToken tokenType, tokensIndex *int) (Cell, error) {
	nextToken, err := readNextToken(tokens, tokensIndex)
	if err != nil {
		return nil, err
	}
	if nextToken.typ == closeParToken {
		extractNextToken(tokens, tokensIndex)
		return nil, nil
	}
	left, err := ricParse(tokens, tokensIndex)
	if err != nil {
		return nil, err
	}
	top := makeCons(left, nil)
	actCons := top

	nextToken, err = readNextToken(tokens, tokensIndex)
	if err != nil {
		return nil, err
	}
	if nextToken.typ == closeParToken {
		extractNextToken(tokens, tokensIndex)
		return top, nil
	}

	for {
		actualToken, err := readNextToken(tokens, tokensIndex)
		if err != nil {
			return nil, err
		}
		if actualToken.typ == tokDot {
			extractNextToken(tokens, tokensIndex) // extract the dot
			// last element
			right, err := ricParse(tokens, tokensIndex)

			if err != nil {
				return nil, err
			}
			closePar, err := extractNextToken(tokens, tokensIndex)
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
			right, err := ricParse(tokens, tokensIndex)
			if err != nil {
				return nil, err
			}
			tmp := makeCons(right, nil)
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

			maybeClosePar, err := readNextToken(tokens, tokensIndex)
			if err != nil {
				return nil, err
			}

			if maybeClosePar.typ == closeParToken {
				extractNextToken(tokens, tokensIndex)
				return top, nil
			}
		}
	}
}

// ParseError
type ParseError struct {
	err string
}

func (e ParseError) Error() string {
	return e.err
}
