package main

import (
	"errors"
	"fmt"
)


// (+ a b)
func procAdd(symbols ...Symbol) (Symbol, error) {
	if len(symbols) != 2 {
		return Symbol{}, errors.New(fmt.Sprint(`"+" takes 2 arguments, got: `, len(symbols)))
	}
	numTypeErr := errors.New(`arguments must match and be numeric type ∈ {Integral, Floating}; got type "` + symbols[0].Kind.String() + `" and "` + symbols[1].Kind.String() + `"`)

	// TODO - some kind of unify for types to match the first argument
	t := symbols[0].Kind
	for _, s := range symbols[1:] {
		if s.Kind != t {
			return Symbol{}, numTypeErr
		}
	}

	switch t {
	case Integral:
		return Symbol{Integral, symbols[0].Contents.(int64) + symbols[1].Contents.(int64)}, nil
	case Floating:
		return Symbol{Floating, symbols[0].Contents.(float64) + symbols[1].Contents.(float64)}, nil
	default:
		// numTypeErr
		;
	}

	return Symbol{}, numTypeErr
}

// (- a b)
func procSub(symbols ...Symbol) (Symbol, error) {
	if len(symbols) != 2 {
		return Symbol{}, errors.New(`"-" takes 2 arguments`)
	}

	numTypeErr := errors.New(`arguments must match and be numeric type ∈ {Integral, Floating}; got type "` + symbols[0].Kind.String() + `" and "` + symbols[1].Kind.String() + `"`)

	// TODO - some kind of unify for types to match the first argument
	t := symbols[0].Kind
	for _, s := range symbols[1:] {
		if s.Kind != t {
			return Symbol{}, numTypeErr
		}
	}

	switch t {
	case Integral:
		return Symbol{Integral, symbols[0].Contents.(int64) - symbols[1].Contents.(int64)}, nil
	case Floating:
		return Symbol{Floating, symbols[0].Contents.(float64) - symbols[1].Contents.(float64)}, nil
	default:
		// numTypeErr
		;
	}

	return Symbol{}, numTypeErr
}

// (* a b)
func procMult(symbols ...Symbol) (Symbol, error) {
	if len(symbols) != 2 {
		return Symbol{}, errors.New(fmt.Sprint(`"*" takes 2 arguments, got: `, len(symbols)))
	}
	numTypeErr := errors.New(`arguments must match and be numeric type ∈ {Integral, Floating}; got type "` + symbols[0].Kind.String() + `" and "` + symbols[1].Kind.String() + `"`)

	// TODO - some kind of unify for types to match the first argument
	t := symbols[0].Kind
	for _, s := range symbols[1:] {
		if s.Kind != t {
			return Symbol{}, numTypeErr
		}
	}

	switch t {
	case Integral:
		return Symbol{Integral, symbols[0].Contents.(int64) * symbols[1].Contents.(int64)}, nil
	case Floating:
		return Symbol{Floating, symbols[0].Contents.(float64) * symbols[1].Contents.(float64)}, nil
	default:
		// numTypeErr
		;
	}

	return Symbol{}, numTypeErr
}

// (/ a b)
func procDiv(symbols ...Symbol) (Symbol, error) {
	if len(symbols) != 2 {
		return Symbol{}, errors.New(`"/" takes 2 arguments`)
	}

	numTypeErr := errors.New(`arguments must match and be numeric type ∈ {Integral, Floating}; got type "` + symbols[0].Kind.String() + `" and "` + symbols[1].Kind.String() + `"`)

	// TODO - some kind of unify for types to match the first argument
	t := symbols[0].Kind
	for _, s := range symbols[1:] {
		if s.Kind != t {
			return Symbol{}, numTypeErr
		}
	}

	switch t {
	case Integral:
		return Symbol{Integral, symbols[0].Contents.(int64) / symbols[1].Contents.(int64)}, nil
	case Floating:
		return Symbol{Floating, symbols[0].Contents.(float64) / symbols[1].Contents.(float64)}, nil
	default:
		// numTypeErr
		;
	}

	return Symbol{}, numTypeErr
}

