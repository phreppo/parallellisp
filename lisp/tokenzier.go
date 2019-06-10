package lisp

import (
	"strconv"
	"strings"
)

type tokenType int

const (
	tokNone          tokenType = 0
	tokOpen          tokenType = 1
	tokClose         tokenType = 2
	tokDot           tokenType = 3
	tokQuote         tokenType = 4
	tokSym           tokenType = 5
	tokNum           tokenType = 6
	tokStr           tokenType = 7
	tokOpenParallel  tokenType = 8
	tokCloseParallel tokenType = 9
)

const (
	dotChar              = '.'
	openParChar          = '('
	closeParChar         = ')'
	openParParallelChar  = '{'
	closeParParallelChar = '}'
	quoteChar            = '\''
)

var atomicCharTokens = map[rune]bool{
	dotChar:              true,
	openParChar:          true,
	closeParChar:         true,
	openParParallelChar:  true,
	closeParParallelChar: true,
	quoteChar:            true,
}

func isAtmoicCharToken(r rune) bool {
	_, ok := atomicCharTokens[r]
	return ok
}

type token struct {
	typ tokenType
	str string
	val int
}

func (t token) String() string {
	switch t.typ {
	case tokNone:
		return "NONE"
	case tokOpen:
		return "("
	case tokOpenParallel:
		return "{"
	case tokClose:
		return ")"
	case tokCloseParallel:
		return "}"
	case tokDot:
		return "."
	case tokQuote:
		return "'"
	case tokSym:
		return t.str
	case tokNum:
		return strconv.Itoa(t.val)
	case tokStr:
		return "\"" + t.str + "\""
	default:
		return ""
	}
}

// tokenize produces an array fo tokens
func tokenize(source string) []token {
	tok, rest := readOneToken(removeComments(source))
	var result []token
	for (tok.typ) != tokNone {
		result = append(result, tok)
		tok, rest = readOneToken(rest)
	}
	return result
}

func removeComments(source string) string {
	if source == "" {
		return ""
	}
	if source[0] != ';' {
		return string(source[0]) + removeComments(source[1:])
	}
	for index, r := range source {
		if r == '\n' {
			return removeComments(source[index+1:])
		}
	}
	return ""
}

// returns the token and the remaining string
func readOneToken(source string) (token, string) {
	if source == "" {
		return token{typ: tokNone}, source
	}
	nextChar, index := firstChar(source)
	if index < 0 {
		return token{typ: tokNone}, source
	} else if nextChar == openParChar {
		return token{typ: tokOpen}, source[index+1:]
	} else if nextChar == closeParChar {
		return token{typ: tokClose}, source[index+1:]
	} else if nextChar == openParParallelChar {
		return token{typ: tokOpenParallel}, source[index+1:]
	} else if nextChar == closeParParallelChar {
		return token{typ: tokCloseParallel}, source[index+1:]
	} else if nextChar == dotChar {
		return token{typ: tokDot}, source[index+1:]
	} else if nextChar == quoteChar {
		return token{typ: tokQuote}, source[index+1:]
	} else if nextChar == '"' {
		rest := source[index+1:]
		stringResult, rest := readUntilDoubleQuote(rest)
		return token{typ: tokStr, str: stringResult}, rest
	} else {
		firstWord, rest := firstWordOrNumber(source)
		if num, err := strconv.Atoi(firstWord); err == nil {
			// it's a num
			return token{typ: tokNum, val: num}, rest
		}
		return token{typ: tokSym, str: firstWord}, rest
	}
}

// returns the char and the position in the string of the char
// returns -1 if the string has no first char
func firstChar(str string) (byte, int) {
	for i := 0; i < len(str); i++ {
		if str[i] != ' ' && str[i] != '\n' {
			return str[i], i
		}
	}
	return 0, -1
}

// resturns the first word or number in the string and the rest of the string
func firstWordOrNumber(str string) (string, string) {
	_, wordBeginningIndex := firstChar(str)
	stringWithoutBlanks := str[wordBeginningIndex:]
	result := ""
	for i, r := range stringWithoutBlanks {
		if r == '\n' || r == ' ' || isAtmoicCharToken(r) {
			return result, stringWithoutBlanks[i:]
		}
		result += string(r)
	}
	return result, ""
}

// reads until the first double quote in the string and resturns the rest of the string
func readUntilDoubleQuote(str string) (string, string) {
	result := ""
	for i, r := range str {
		if r == '"' {
			return result, str[i+1:]
		}
		result += string(r)
	}
	return "", ""
}

func stringContainsDoubleQuotes(str string) bool {
	return strings.Contains(str, "\"")
}
