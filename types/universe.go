// =================================================================================================
// Alex Peters - 2024
// =================================================================================================
package types

import (
	"fmt"

	"github.com/petersalex27/yew/common/math"
)

type Universe int

const (
	Type0 Universe = iota
	Type1
)

const (
	Polytype1 Universe = -2
	Polytype0 Universe = -1
	Monotype0 Universe = 0
	Monotype1 Universe = 1
)

func (Universe) TypeClassification() typeClassification { return universeClass }

func (Universe) Locate(Variable) bool { return false }

func stringPolyUniverse(u Universe) string {
	if u == Universe(-1) {
		return "**"
	}
	return fmt.Sprintf("**{%d}", math.Abs(u)-1)
}

func (u Universe) String() string {
	if u < Universe(0) {
		return stringPolyUniverse(u)
	} else if u == Type0 {
		return "*"
	}
	return fmt.Sprintf("*{%d}", uint(u))
}

func (Universe) Known() bool { return true }

func (Universe) Substitute(*Term, Variable, Term) { return }
