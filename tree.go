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

