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
		// TODO - validate the procedure token type, etc.
		token = ts.next()

		fmt.Println("» parsing operator =", token)

		if token.Kind == NIL {
			return nil, errors.New(fmt.Sprintf(`unexpected end of token stream at beginning of expression after token "%v"`, ts.previous().name))
		}

		symbol, err := token2symbol(token)
		if err != nil {
			return nil, err
		}

		tree, _ = NewTree(symbol)

		// TODO - this hack works, does this mean the logic is flawed?
		//ts.rewind()

		// Recursive descent? for consuming children
		loop:
		for {
			token = ts.next()

			fmt.Println("» next token =", token)

			switch token.Kind {
			case NIL:
				// End of token stream
				// TODO - just return?
				break loop

			case End:
				// End of expression means nothing more to consume
				// TODO - put at end of block?
				return tree, nil

			case Begin:
				// Beginning of new expression
				subtree, err := parse(ts)
				if err != nil {
					return nil, err
				}


				// TODO - more?
				if subtree != nil {
					tree.Children = append(tree.Children, subtree)
				}

				fmt.Println("»» Nested Begin length = ", tree.Length())
			}

			thisSymbol, _ := token2symbol(token)
			thisLeaf, _ := NewTree(thisSymbol)
			tree.Children = append(tree.Children, thisLeaf)

			// Recursively descend on new expression
			subtree, err := parse(ts)
			if err != nil {
				return nil, err
			}

			// TODO - more?
			if subtree != nil {
				tree.Children = append(tree.Children, subtree)
			}

			fmt.Println("»» Top Begin length = ", tree.Length())
		}

	case End:
		// TODO - wrong?
		return nil, nil
		//return nil, errors.New(fmt.Sprintf(`"unexpected end of expression after "%v"`, ts.previous()))

	default:
		// TODO - use this rather than just pushing it through?
		//return nil, errors.New(fmt.Sprintf(`could not parse for AST, unknown token type "%v"`, token.Kind))

		// If we get something that isn't for expression control, return just that symbol as a mono-tree?
		symbol, err := token2symbol(token)
		if err != nil {
			return nil, err
		}

		subtree, _ := NewTree(symbol)



		return subtree, nil

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
