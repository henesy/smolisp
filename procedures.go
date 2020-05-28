package main

import (
	"errors"
	"fmt"
)

// (+ a b)
func procAdd(symbols ...Symbol) (Symbol, error) {
	if len(symbols) != 2 {
		return Symbol{}, errors.New(fmt.Sprint("+ takes 2 arguments got ", len(symbols)))
	}

	if symbols[0].Kind == Integral && symbols[1].Kind == Integral {
		return Symbol{Integral, symbols[0].Contents.(int) + symbols[1].Contents.(int)}, nil
	}

	return Symbol{}, errors.New("arguments must be numeric types {Integral, }")
}

// (- a b)
func procSub(symbols ...Symbol) (Symbol, error) {
	if len(symbols) != 2 {
		return Symbol{}, errors.New("- takes 2 arguments")
	}

	if symbols[0].Kind == Integral && symbols[1].Kind == Integral {
		return Symbol{Integral, symbols[0].Contents.(int) - symbols[1].Contents.(int)}, nil
	}

	return Symbol{}, errors.New("arguments must be numeric types {Integral, }")
}
