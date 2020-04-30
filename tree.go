package main

import (
	"errors"
	"fmt"
)

/* Organization:

(+ (+ 2 3) 5)

		+
	  /   \
	+		5
  /	  \
2		3
*/

// The AST
type Tree struct {
	Datum    Symbol
	Children []*Tree
}

// Parse a list of tokens into an AST
func parse(ts TokenScanner) (*Tree, error) {
	token := ts.next()
	var tree *Tree

	fmt.Println("» parsing token =", token)

	switch token.Vtype {
	case Begin:
		// TODO - validate the procedure token type, etc.
		symtok := ts.next()

		if symtok.Vtype == NIL {
			return nil, errors.New(fmt.Sprintf(`unexpected end of token stream at beginning of expression after token "%v"`, ts.previous().name))
		}

		symbol, err := token2symbol(symtok)
		if err != nil {
			return nil, err
		}

		tree = NewTree(symbol)

		// TODO - this hack works, does this mean the logic is flawed?
		ts.rewind()

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

			/* TODO - feels wrong to remove this
			// Find out what symbol we have
			symbol, err := token2symbol(token)
			if err != nil {
				return nil, err
			}
			*/

			// Recursively descend on new expression
			subtree, err := parse(ts)
			if err != nil {
				return nil, err
			}

			// TODO - more?
			if subtree != nil {
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

// Create a new, single-node tree
func NewTree(symbol Symbol) *Tree {
	return &Tree{symbol, make([]*Tree, 0, maxChildren)}
}
