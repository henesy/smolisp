package main

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
	Datum		Symbol
	Children	[]*Tree
}

// Create a new, single-node tree
func NewTree(symbol Symbol) *Tree {
	return &Tree{symbol, make([]*Tree, 0, maxChildren)}
}
