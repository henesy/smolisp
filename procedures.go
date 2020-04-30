package main

import (
	"errors"
)

// (+ a b)
func procAdd(symbols ...Symbol) (Symbol, error) {
	if len(symbols) != 2 {
		return Symbol{}, errors.New("+ takes 2 arguments")
	}

	if symbols[0].Vtype == Integral && symbols[1].Vtype == Integral {
		return Symbol{Integral, symbols[0].Contents.(int) + symbols[1].Contents.(int)}, nil
	}

	return Symbol{}, errors.New("arguments must be numeric types {Integral, }")
}

// (- a b)
func procSub(symbols ...Symbol) (Symbol, error) {
	if len(symbols) != 2 {
		return Symbol{}, errors.New("- takes 2 arguments")
	}

	if symbols[0].Vtype == Integral && symbols[1].Vtype == Integral {
		return Symbol{Integral, symbols[0].Contents.(int) - symbols[1].Contents.(int)}, nil
	}

	return Symbol{}, errors.New("arguments must be numeric types {Integral, }")
}
