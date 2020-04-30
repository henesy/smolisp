package main

import (
	"flag"
	"fmt"
	"bufio"
	"os"
	"io"
	"math"
	"strings"
	"errors"
)

const (
	maxChildren		= 32
)

// Represents the inner type of a name
type Vtype int
const (
	Integral	Vtype = iota	// (+ »2« »3«)
	Floating					// (pow »3.14« 2)
	Procedure					// (»foo« a b)
	Begin						// »(«foo a b)
	End							// (foo a b»)«
	Value						// (let »a« 2)
	NIL							// Sentinel values are bad
)

// Represents a lexed token
type Token struct {
	Vtype
	name	string
}

// Represents a known symbol
type Symbol struct {
	Vtype
	Value	interface{}		// Assertion determined by .Type
}

var (
	prompt	= "» "
	symbols	map[string]Symbol
	ast		*Tree
)


// Parse a list of tokens into an AST
func parse(ts TokenScanner) (*Tree, error) {
	token := ts.next()
	var tree *Tree

	switch token.Vtype {
	case Begin:
		// TODO - validate the procedure token type, etc.
		symtok := ts.next()

		// Check if this is a known symbol
		symbol, ok := symbols[symtok.name]
		if !ok {
			symbol.Vtype = symtok.Vtype

			// Determine what the value of the symbol is
			;
		}

		tree = &Tree{symbol, make([]*Tree, 0, maxChildren)}

	case End:

	}

	return tree, nil
}


// Small, toy, lisp-like
func main() {
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)

	symbols = map[string]Symbol{
	    "+":		Symbol{Procedure,	procAdd},
	    "-":		Symbol{Procedure,	procSub},
	    "π":		Symbol{Floating,	math.Pi},
	    "billion":	Symbol{Integral,	1_000_000_000},
	}

	// Main loop
	repl:
	for {
		fmt.Print(prompt)

		// Just read a line, TODO - make by-rune reading later
		text, err := reader.ReadString('\n')
		if err == io.EOF {
			break repl
		}

		// Ignore empty lines
		if len(text) < 2 {
			continue repl
		}

		// Strip newline
		text = text[:len(text)-1]

		/* Tokenizing */
		tokens, err := tokenize(text)

		if err != nil {
			fmt.Println("err: tokenizing failed -", err)
			continue repl
		}

		fmt.Println(tokens)

		/* Parsing */
		ts := NewTokenScanner(tokens)

		ast, err := parse(ts)

		if err != nil {
			fmt.Println("err: could not parse -", err)
			continue repl
		}

		fmt.Println(*ast)

		/* Evaluate */
		;

		/* Output */
		;
	}
}


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

func (ts TokenScanner) next() Token {
	if ts.i >= len(ts.tokens) {
		return Token{NIL, "<nil>"}
	}

	// Hack
	defer func() { ts.i++ }()

	return ts.tokens[ts.i]
}

func (ts TokenScanner) previous() Token {
	if ts.i <= 0 {
		return ts.tokens[0]
	}

	return ts.tokens[ts.i-1]
}