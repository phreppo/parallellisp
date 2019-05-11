package parser

import (
	"strconv"
	"strings"
)

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

// Tokenize produces an array fo tokens
func Tokenize(source string) []token {
	tok, rest := nextToken(source)
	var result []token
	for (tok.typ) != tokNone {
		result = append(result, tok)
		tok, rest = nextToken(rest)
	}
	return result
}

// returns the token and the remaining string
func nextToken(source string) (token, string) {
	if source == "" {
		return token{typ: tokNone}, source
	}
	nextChar, index := firstChar(source)
	if index < 0 {
		return token{typ: tokNone}, source
	} else if nextChar == '(' {
		return token{typ: tokOpen}, source[index+1:]
	} else if nextChar == ')' {
		return token{typ: tokClose}, source[index+1:]
	} else if nextChar == '.' {
		return token{typ: tokDot}, source[index+1:]
	} else if nextChar == '\'' {
		return token{typ: tokQuote}, source[index+1:]
	} else if nextChar == '"' {
		rest := source[index+1:]
		stringResult, rest := readUntilDoubleQuote(rest)
		return token{typ: tokStr, str: stringResult}, rest
	} else {
		firstWord, newIndex := firstWordOrNumber(source)
		if num, err := strconv.Atoi(firstWord); err == nil {
			// it's a num
			return token{typ: tokNum, val: num}, source[newIndex+1:]
		}
		return token{typ: tokSym, str: firstWord}, source[newIndex+1:]
	}
	return token{typ: tokNone}, source
}

// returns the char and the position in the string of the char
// returns -1 if the string has no first char
func firstChar(str string) (byte, int) {
	for i := 0; i < len(str); i++ {
		if str[i] != ' ' {
			return str[i], i
		}
	}
	return 0, -1
}

// resturns the first word or number in the string and the index in which it ends
func firstWordOrNumber(str string) (string, int) {
	_, wordBeginningIndex := firstChar(str)
	stringWithoutBlanks := str[wordBeginningIndex:]
	result := ""
	for i, r := range stringWithoutBlanks {
		if r == ' ' || r == '.' || r == '(' || r == ')' || r == '\'' {
			return result, i - 1 + wordBeginningIndex
		}
		result += string(r)
	}
	return result, len(str) - 1
}

func readUntilDoubleQuote(str string) (string, string) {
	// reads until the first double quote in the string and resturns the rest of the string
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
