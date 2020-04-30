package main

import (
	"fmt"
)

/* Magic and garbage */

func (t Token) String() string {
	return "{" + fmt.Sprint(t.Vtype) + ", " + t.name + "}"
}

func (vt Vtype) String() string {
	switch vt {
	case Integral:	return "Integral"
	case Floating:	return "Floating"
	case Procedure:	return "Procedure"
	case Begin:		return "Begin"
	case End:		return "End"
	case Value:		return "Value"
	default:		return "UNKNOWN TYPE"
	}
}
