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

-------

Building the tree is as so:

    nil


    nil
     |
    nil


    nil
     |
     +


    nil
     |
     +
    /
   2


    nil
     |
     +
   /   \
  2     3


    +
   / \
  2	  3

*/

const dummyName = "dummy root"

// The AST
type Tree struct {
	Symbol                              // Symbol value of this node
	Children []*Tree                    // Child nodes of this node
	Eval     func(*Tree) (*Tree, error) // Returns the "value" relative to the current node
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

	// Insert dummy head to enable simplified recursion logic
	tree, _ = InitTree(Symbol{Dummy, dummyName})

	for ts.hasNext() {
		token := ts.next()

		if chatty {
			fmt.Println("» parsing token =", token)
		}

		ingest(ts, tree, token)
	}

	if len(tree.Children) <= 0 {
		return nil, errors.New("no complete expression provided (empty AST)")
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

	// Inserts as a child into the tree (for values)
	insertChild := func() error {
		symbol, err := token2symbol(token)
		if err != nil {
			return errors.New(fmt.Sprint("token2symbol failed → ", err))
		}

		child, err := InitTree(symbol)
		if err != nil {
			return err
		}

		tree.Children = append(tree.Children, child)
		return nil
	}

	switch token.Kind {
	case Begin:
		// Shift down into the next child of the current node
		child := NewTree()
		tree.Children = append(tree.Children, child)

		// Recurse
		ingest(ts, child, ts.next())

	case End:
		// Propagate back up
		return nil


	case Procedure:
		// Take advantage of 'Begin (' percolating down the tree as a dummy
		symbol, err := token2symbol(token)
		if err != nil {
			return err
		}

		// Take over current node
		tree.Symbol = symbol
		tree.Eval, err = getHandler(symbol)
		if err != nil {
			return err
		}

		// Ingest our child nodes
	loop:
		for {
			token := ts.next()

			switch token.Kind {
			// Short circuit if we reach a NIL
			// TODO - is this grounds for better error reporting?
			case NIL:
				break loop

			// Break out when we get matching 'End )'
			case End:
				break loop

			default:
				// Recurse
				ingest(ts, tree, token)
			}
		}


	case Floating:
		fallthrough
	case Integral:
		// Insert ourselves as a child of the current node, we are but a value
		fallthrough
	case Value:
		// Look ourselves up to see if we know the symbol
		return insertChild()



	/* end type switch */
	}

	return nil
}

// Create a new tree containing one stub node
func NewTree() *Tree {
	return &Tree{Symbol{NIL, "NewTree holder"}, make([]*Tree, 0, maxChildren), nil}
}

// Create a new tree from a Symbol
func InitTree(symbol Symbol) (*Tree, error) {
	handler, err := getHandler(symbol)
	if err != nil {
		return nil, err
	}

	return &Tree{symbol, make([]*Tree, 0, maxChildren), handler}, nil
}

// Return a function to handle a given symbol
// TODO - this could probably be cleaner
func getHandler(symbol Symbol) (func(*Tree) (*Tree, error), error) {
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
					return nil, errors.New(fmt.Sprint("child eval failed → ", err))
				}

				symbol := childTree.Symbol

				childSyms = append(childSyms, symbol)
			}

			result, err := operation(childSyms...)

			if err != nil {
				return nil, errors.New(fmt.Sprint("procedure evaluation failed to consume children → ", err))
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
