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

func (u Universe) asConstant() Constant { return Constant{C: u.String()} }

// universes are cumulative
//		x : Type : Type 1 : Type 2 : ...
//		x : Type n, n <= m, x : Type m
func (u Universe) GetKind() (Term, Type) {
	return u, Universe(u + 1)
}

func (Universe) CollectVariables(m map[string]Variable) map[string]Variable { return m }

const (
	Polytype1 Universe = -2
	Polytype0 Universe = -1
	Monotype0 Universe = 0
	Monotype1 Universe = 1
)

func (Universe) Pos() (start, end int) { return 0, 0 }

func (Universe) TypeClassification() typeClassification { return universeClass }

func (Universe) Locate(Variable) bool { return false }

func stringPolyUniverse(u Universe) string {
	if u == Universe(-1) {
		return "Polytype"
	}
	return fmt.Sprintf("Polytype {%d}", math.Abs(u)-1)
}

func (u Universe) String() string {
	if u < Universe(0) {
		return stringPolyUniverse(u)
	} else if u == Type0 {
		return "Type"
	}
	return fmt.Sprintf("Type %d", uint(u))
}

func (Universe) Known() bool { return true }

func (u Universe) Substitute(dest *Term, _ Variable, _ Term) {
	*dest = u
}
