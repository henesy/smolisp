package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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
		token := Token{
			name: word,
		}

		switch {
		case word == "(":
			token.Kind = Begin

		case word == ")":
			token.Kind = End

		// TODO - case out Floating and negatives
		case word[0] >= '0' && word[0] <= '9':
			token.Kind = Integral

		default:
			if i == 0 {
				return nil, errors.New(`first rune must be "(", got "` + word + `"`)
			}

			if tokens[i-1].Kind == Begin {
				token.Kind = Procedure
			} else {
				token.Kind = Value
			}
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

// Convert a token into a full symbol
func token2symbol(token Token) (Symbol, error) {
	var symbol Symbol

	// Check if this is a known symbol
	symbol, ok := symbols[token.name]
	if !ok {
		// Unknown symbol, so we build it
		symbol.Kind = token.Kind

		// Determine what the value of the symbol is
		value, err := findValue(token)

		if err != nil {
			return symbol, err
		}

		symbol.Contents = value
	}

	return symbol, nil
}

// For scanning across tokens
type TokenScanner struct {
	tokens []Token
	i      int
}

// Create a new token scanner from a slice of tokens
func NewTokenScanner(tokens []Token) *TokenScanner {
	return &TokenScanner{
		tokens: tokens,
		i:      0,
	}
}

// Does the TokenScanner have another token?
func (ts *TokenScanner) hasNext() bool {
	if ts.i >= len(ts.tokens) {
		return false
	}
	
	return true
}

// Return current token, then shift forwards
func (ts *TokenScanner) next() Token {
	if ts.i >= len(ts.tokens) {
		return Token{NIL, "<nil>"}
	}

	// Hack due to lack of ++i
	defer func() { ts.i++ }()

	return ts.tokens[ts.i]
}

// Get the previous token read, no shifting is done
func (ts *TokenScanner) previous() Token {
	if ts.i <= 0 {
		return ts.tokens[0]
	}

	return ts.tokens[ts.i-1]
}

// Shift back 1 token and return said token
func (ts *TokenScanner) rewind() Token {
	if ts.i <= 0 {
		return ts.tokens[0]
	}

	ts.i--

	return ts.tokens[ts.i]
}

// Look up and build value of a given token
func findValue(token Token) (interface{}, error) {
	// TODO - more types
	switch token.Kind {
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
		return nil, errors.New(fmt.Sprintf(`unknown type, cannot determine value of "%v"`, token.Kind))
	}

	// Unreachable

}
