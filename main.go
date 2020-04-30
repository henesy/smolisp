package main

import (
	"flag"
	"fmt"
	"bufio"
	"os"
	"io"
	"math"
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
	Contents	interface{}		// Assertion determined by .Type
}

var (
	prompt	= "» "
	symbols	map[string]Symbol
	ast		*Tree
)


// Convert a token into a full symbol
func token2symbol(token Token) (Symbol, error) {
	var symbol Symbol

	// Check if this is a known symbol
	symbol, ok := symbols[token.name]
	if !ok {
		// Unknown symbol, so we build it
		symbol.Vtype = token.Vtype

		// Determine what the value of the symbol is
		value, err := findvalue(token)

		if err != nil {
			return symbol, err
		}

		symbol.Contents = value
	}

	return symbol, nil
}

// Parse a list of tokens into an AST
func parse(ts TokenScanner) (*Tree, error) {
	token := ts.next()
	var tree *Tree

	switch token.Vtype {
	case Begin:
		// TODO - validate the procedure token type, etc.
		symtok := ts.next()

		if symtok.Vtype == NIL {
			return nil, errors.New(fmt.Sprintf(`unexpected end of token stream at beginning of expression after token ""`, ts.previous().name))
		}

		fmt.Println("symtok =", symtok)

		symbol, err := token2symbol(symtok)
		if err != nil {
			return nil, err
		}

		fmt.Println("Symbol value →", symbol.Contents)

		tree = NewTree(symbol)

		// Recursive descent?
		for {
			token := ts.next()

			// End of token stream
			if token.Vtype == NIL {
				// TODO - just return?
				break
			}

			// End of expression means nothing more to consume
			// TODO - put at end of block?
			if token.Vtype == End {
				return tree, nil
			}

			// Find out what symbol we have
			symbol, err := token2symbol(token)
			if err != nil {
				return nil, err
			}

			fmt.Println("→ symbol descended =", symbol)

			// Recursively descend on new expression
			subtree, err := parse(ts)
			if err != nil {
				return nil, err
			}

			// TODO - more?
			if subtree != nil {
				fmt.Println("→ tree appended =", *subtree)
				tree.Children = append(tree.Children, subtree)
			}
		}


	case End:
		// TODO - wrong?
		return nil, nil
		//return nil, errors.New(fmt.Sprintf(`"unexpected end of expression after "%v"`, ts.previous()))

	default:
		// TODO - use this rather than just pushing it through?
		//return nil, errors.New(fmt.Sprintf(`could not parse for AST, unknown token type "%v"`, token.Vtype))

		// If we get something that isn't for expression control, return just that symbol as a mono-tree?
		symbol, err := token2symbol(token)
		if err != nil {
			return nil, err
		}

		subtree := NewTree(symbol)

		return subtree, nil

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

		if err != nil || ast == nil {
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


