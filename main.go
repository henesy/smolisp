package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
)

const (
	maxChildren = 32
)

// Represents the inner type of a name
type Kind int

const (
	Integral  Kind = iota  // (+ »2« »3«)
	Floating               // (pow »3.14« 2)
	Procedure              // (»foo« a b)
	Begin                  // »(«foo a b)
	End                    // (foo a b»)«
	Value                  // (let »a« 2)
	NIL                    // Sentinel values are bad
)

// Represents a lexed token
type Token struct {
	Kind
	name string
}

// Represents a known symbol
type Symbol struct {
	Kind
	Contents interface{} // Assertion determined by .Type
}

var (
	prompt  = "» "
	symbols map[string]Symbol
	ast     *Tree
)

// Small, toy, lisp-like
func main() {
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)

	symbols = map[string]Symbol{
		"+":       Symbol{Procedure, procAdd},
		"-":       Symbol{Procedure, procSub},
		"π":       Symbol{Floating, math.Pi},
		"billion": Symbol{Integral, 1_000_000_000},
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

		fmt.Println("TOKENS =", tokens)

		/* Parsing */
		ts := NewTokenScanner(tokens)

		ast, err := parse(ts)

		if err != nil || ast == nil {
			fmt.Println("err: could not parse -", err)
			continue repl
		}

		fmt.Println("AST root =", *ast)

		/* Evaluate */

		final, err := ast.Eval(ast)

		if err != nil {
			fmt.Println("err: could not eval the AST -", err)
		}

		/* Output */

		fmt.Println(final.Symbol)

	}

	fmt.Println("\nGoodbye ☺")
}
