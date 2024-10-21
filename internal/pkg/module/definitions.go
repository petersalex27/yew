package module

import (
	"github.com/llir/llvm/ir"
	"github.com/petersalex27/yew/internal/symbol"
)

type (
	Data struct {
		name     string
		internal *ir.Module
		symbols  map[string]symbol.Symbolic
	}

	symbolType byte

	precedence byte

	sym struct {
		symbol     symbol.Symbolic
		precedence precedence
		typ        symbolType
	}
)

const (
	Prefix symbolType = iota
	Infixl
	Infixr
)

const (
	Zero  precedence = iota // example: terms, infix function application
	One                     // example: monadic operators
	Two                     // example: logical or
	Three                   // example: logical and
	Four                    // example: functor operators, comparison operators
	Five                    // example: list operators
	Six                     // example: additive operators
	Seven                   // example: multiplicative operators
	Eight                   // example: exponentiation
	Nine                    // example: function composition
	Ten                     // example: function application
)

func (d Data) Name() string {
	return d.name
}

func (d Data) Lookup(name string) (symbol.Symbolic, bool) {
	s, ok := d.symbols[name]
	return s, ok
}
