package main

import (
	"fmt"
)

/* Magic and garbage */

func (s Symbol) String() string {
	switch s.Kind {
	case Integral:
		return fmt.Sprint(s.Contents.(int))
	case Floating:
		return fmt.Sprint(s.Contents.(float64))
	case Procedure:
		return fmt.Sprintf("The procedure ", s.Contents.(func (symbols ...Symbol) (Symbol, error)))
	case Begin:
		return "( (should never happen)"
	case End:
		return ") (should never happen)"
	case Value:
		return "<not implemented>"
	case NIL:
		return "NIL (probably an interpreter bug)"
	default:
		return "UNKNOWN TYPE"
	}

}

func (t Token) String() string {
	return fmt.Sprintf(`{%v, %v}`, t.Kind, t.name)
}

func (k Kind) String() string {
	switch k {
	case Integral:
		return "Integral"
	case Floating:
		return "Floating"
	case Procedure:
		return "Procedure"
	case Begin:
		return "Begin"
	case End:
		return "End"
	case Value:
		return "Value"
	case NIL:
		return "NIL (probably an interpreter bug)"
	default:
		return "UNKNOWN TYPE"
	}
}
