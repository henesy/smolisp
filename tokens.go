package main

import (
	"strings"
	"strconv"
	"errors"
	"fmt"
)

/* For tokenizing, etc. */

// Tokenize text into a list of tokens
func tokenize(text string) ([]Token, error) {
	// So we can split on whitespace, TODO - be better, maybe use strings.FieldsFunc() ?
	text = strings.ReplaceAll(text, "(", " ( ")
	text = strings.ReplaceAll(text, ")", " ) ")

	words := strings.Fields(text)
	tokens := make([]Token, 0, len(words))

	// Determine the virtual type for each token
	for i, word := range words {
		token := Token {
			name: word,
		}

		switch {
		case word == "(":
			token.Vtype = Begin

		case word == ")":
			token.Vtype = End

		// TODO - case out Floating and negatives
		case word[0] >= '0' && word[0] <= '9':
			token.Vtype = Integral

		default:
			if i == 0 {
				return nil, errors.New(`first rune must be "(", got "` + word + `"`)
			}

			if tokens[i-1].Vtype == Begin {
				token.Vtype = Procedure
			} else {
				token.Vtype = Value
			}
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}


/* For scanning across tokens */
type TokenScanner struct {
	tokens	[]Token
	i		int
}

func NewTokenScanner(tokens []Token) TokenScanner {
	return TokenScanner {
		tokens: tokens,
		i:		0,
	}
}

func (ts *TokenScanner) next() Token {
	if ts.i >= len(ts.tokens) {
		return Token{NIL, "<nil>"}
	}

	// Hack
	defer func() { ts.i++ }()

	return ts.tokens[ts.i]
}

func (ts *TokenScanner) previous() Token {
	if ts.i <= 0 {
		return ts.tokens[0]
	}

	return ts.tokens[ts.i-1]
}

// Look up and build value of a given token
func findvalue(token Token) (interface{}, error) {
	var value interface{}

	// TODO - more types
	switch token.Vtype {
	case Integral:
		n, err := strconv.Atoi(token.name)
		if err != nil {
			return nil, errors.New(fmt.Sprintf(`could not convert to number - %v`, err))
		}

		return n, nil

	case Procedure:
		// TODO - Do we need to look this up twice? Might be redundant with map initialization
		switch token.name {
		case "+":
			return symbols["+"].Contents, nil
		case "-":
			return symbols["-"].Contents, nil
		default:
			return nil, errors.New(fmt.Sprintf(`unknown procedure "%v"`, token.name))
		}

	default:
		return nil, errors.New(fmt.Sprintf(`unknown type, cannot determine value of "%v"`, token.Vtype))
	}

	return value, nil
}
