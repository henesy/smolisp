package main

import (
	"errors"
	"fmt"
)

/*	Organization:

(+ (+ 2 3) 5)

		+
	  /   \
	+		5
  /	  \
2		3

	Evaluation:

A necessarily many to one function:

	+
2		3

The `+` node fans in 3→1 by consuming its children as arguments as well as itself returning its evaluated self as a single-node tree.

The `2` node fans in, incidentally 1→1 as it consumes only itself and returns itself as a single-node tree.

*/

// The AST
type Tree struct {
	Symbol
	Children []*Tree
	Eval		func(*Tree)(*Tree, error)
}


// Get number of nodes in a tree
func (t Tree) Length() int {
	if t.Symbol.Kind == NIL {
		return 0
	}

	count := 1
	for _, child := range t.Children {
		if child != nil {
			count += child.Length()
		}
	}

	return count
}

// Parse a list of tokens into an AST
func parse(ts *TokenScanner) (*Tree, error) {
	token := ts.next()
	var tree *Tree

	fmt.Println("» parsing token =", token)

	switch token.Kind {
	case Begin:
		// Beginning of an expression
		token = ts.next()

	case End:
		// Closing of an expression
		token = ts.next()
		return tree, nil
	}

	symbol, err := token2symbol(token)
	if err != nil {
		return nil, errors.New(fmt.Sprint("token→symbol failed - ", err))
	}

	switch symbol.Kind {
	case Procedure:
		// Only thing not a leaf
		;

	default:
		// Leaf
		return NewTree(symbol)
	}

	return tree, nil
}

// Create a new, single-node tree
func NewTree(symbol Symbol) (*Tree, error) {
	handler, err := getHandler(symbol)
	if err != nil {
		return nil, err
	}

	return &Tree{symbol, make([]*Tree, 0, maxChildren), handler}, nil
}

// Return a function to handle a given symbol
func getHandler(symbol Symbol) (func(*Tree)(*Tree, error), error) {
	switch symbol.Kind {
	case Procedure:
		// Note that we'll receive the `Procedure` node as the root
		return func(tree *Tree) (*Tree, error) {
				operation := symbol.Contents.(func(...Symbol) (Symbol, error))
				childSyms := make([]Symbol, 0, len(tree.Children))

				for _, child := range tree.Children {
					// Reduce the child tree to one result node
					childTree, err := child.Eval(child)
					if err != nil {
						return nil, errors.New(fmt.Sprint("child eval failed - ", err))
					}

					symbol := childTree.Symbol

					childSyms = append(childSyms, symbol)
				}

				result, err := operation(childSyms ...)

				if err != nil {
					return nil, errors.New(fmt.Sprint("procedure evaluation failed to consume children - ", err))
				}

				return NewTree(result)
			}, nil

	case Integral:
		return func(tree *Tree) (*Tree, error) {
				return NewTree(symbol)
			}, nil

	default:
		return func(tree *Tree) (*Tree, error) {
				return NewTree(Symbol{Kind: NIL})
			}, nil
	}
}
