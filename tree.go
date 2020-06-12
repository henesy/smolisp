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

const dummyName = "dummy root"

// The AST
type Tree struct {
	Symbol
	Children []*Tree
	Eval		func(*Tree)(*Tree, error)
}


// Get number of nodes in a tree
func (t Tree) Size() int {
	if t.Symbol.Kind == NIL {
		return 0
	}

	count := 1
	for _, child := range t.Children {
		if child != nil {
			count += child.Size()
		}
	}

	return count
}

// Parse a list of tokens into an AST
func parse(ts *TokenScanner) (*Tree, error) {
	var tree *Tree
	
	tree, _ = InitTree(Symbol{Dummy, dummyName})
	
	for ts.hasNext() {
		token := ts.next()

		if chatty {
			fmt.Println("» parsing token =", token)
		}

		ingest(ts, tree, token)
	}
	
	if len(tree.Children) <= 0 {
		return nil, errors.New("no expression provided (tree's size is 0)")
	}
	
	// Remove the dummy head
	tree = tree.Children[0]
	
	return tree, nil
}

// Ingest a token into the tree
func ingest(ts *TokenScanner, tree *Tree, token Token) error { 
	if chatty {
		fmt.Println("»» ingesting tree size =", tree.Size())
		fmt.Println("»» ingesting token =", token)
	}
	
	switch token.Kind {
	case Begin:
		// Shift down into the next child and forward a token
		child := NewTree()
		tree.Children = append(tree.Children, child)
		ingest(ts, child, ts.next())
		
	case End:
		// Propagate back up
		return nil
	
	case Procedure:
		// Take advantage of begin percolating down the tree, take over current node
		symbol, err := token2symbol(token)
		if err != nil {
			return err
		}
		
		tree.Symbol = symbol
		tree.Eval, err = getHandler(symbol)
		if err != nil {
			return err
		}
		
		// Ingest our children now
		loop:
		for {
			token := ts.next()
		
			switch token.Kind {
			case End:
				break loop
				
			default:
				ingest(ts, tree, token)
			}
		}
		
	case Integral:
		// Insert ourselves as a child
		symbol, err := token2symbol(token)
		if err != nil {
			return errors.New(fmt.Sprint("token→symbol failed - ", err))
		}
		
		child, err := InitTree(symbol)
		if err != nil {
			return err
		}

		tree.Children = append(tree.Children, child)
	}
	
	return nil
}

// Create a new tree containing nothing
func NewTree() (*Tree) {
	return &Tree{Symbol{NIL, "NewTree holder"}, make([]*Tree, 0, maxChildren), nil}
}

// Create a new tree from a Token
func InitTree(symbol Symbol) (*Tree, error) {
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

				return InitTree(result)
			}, nil

	case Integral:
		return func(tree *Tree) (*Tree, error) {
				return InitTree(symbol)
			}, nil

	default:
		return func(tree *Tree) (*Tree, error) {
				return InitTree(Symbol{Kind: NIL})
			}, nil
	}
}
