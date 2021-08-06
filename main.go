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
	Integral  Kind = iota // (+ »2« »3«)
	Floating              // (pow »3.14« 2)
	Procedure             // (»foo« a b)
	Begin                 // »(«foo a b)
	End                   // (foo a b»)«
	Value                 // (let »a« 2)
	Dummy                 // Fake node to be taken over later
	NIL                   // Sentinel values are bad
)

// Represents a lexed token
type Token struct {
	Kind        // Type of a Token
	name string // String representation of the Token
}

// Represents a known symbol
type Symbol struct {
	Kind                 // Type of a Symbol
	Contents interface{} // Type assertion determined by Kind
}

var (
	prompt  string            // Prompt in repl
	symbols map[string]Symbol // Known symbols
	ast     *Tree             // AST used for evaluation - TODO - move into REPL and discard after use
	chatty  bool              // Verbose debugging output
)

// Small, toy, lisp-like
func main() {
	flag.StringVar(&prompt, "p", "» ", "prompt to show in repl")
	flag.BoolVar(&chatty, "D", false, "enable verbose debug output")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)

	symbols = map[string]Symbol{
		"+":       Symbol{Procedure, procAdd},
		"-":       Symbol{Procedure, procSub},
		"*":       Symbol{Procedure, procMult},
		"/":       Symbol{Procedure, procDiv},
		"π":       Symbol{Floating, math.Pi},
		"billion": Symbol{Integral, int64(1_000_000_000)},
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

		if chatty {
			fmt.Println("TOKENS =", tokens)
		}

		/* Parsing */

		ts := NewTokenScanner(tokens)

		ast, err := parse(ts)

		if err != nil || ast == nil || ast.Kind == NIL {
			if err != nil {
				fmt.Println("err: parsing failed →", err)
			} else {
				fmt.Println("err: parsing failed → invalid expression; expected: (procedure args...)")
			}
			continue repl
		}

		if chatty {
			fmt.Println("AST root =", *ast)
		}

		/* Evaluate */

		final, err := ast.Eval(ast)

		if err != nil {
			fmt.Println("err: could not eval the AST →", err)
			continue repl
		}

		/* Output */

		fmt.Println(final.Symbol)

	}

	fmt.Println("\nGoodbye ☺")
}
