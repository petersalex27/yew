// =================================================================================================
// Alex Peters - March 02, 2024
//
// Dependent product type (function type (constant dependent product type), type constructor, etc.)
// =================================================================================================
package types

import "fmt"

// (𝚷x:A.B):s
type Pi struct {
	implicit bool
	// x
	binderVar Variable
	// A
	binderType Type
	// B
	dependent Type
	// s
	kind Sort
}

func (pi Pi) Locate(v Variable) bool {
	return pi.binderType.Locate(v) || pi.dependent.Locate(v)
}

// given binder type `A`, binding variable `x`, and bound type `B` (Pi).String() returns the string
//
//	"(x:A) -> B"
//
// when x is a non-empty string; otherwise,
//
//	"A -> B"
//
// is returned.
func (pi Pi) String() string {
	domain := fmt.Sprintf("%v", pi.binderType)
	if pi.binderVar.name == "" {
		if pi.implicit {
			domain = "{" + domain + "}"
		} // else, domain is unchanged
	} else {
		domain = fmt.Sprintf("%v : ", pi.binderVar) + domain
		if pi.implicit {
			domain = "{" + domain + "}"
		} else {
			domain = "(" + domain + ")"
		}
	}
	return fmt.Sprintf("%v -> %v", domain, pi.dependent)
}

// substitutes all free occurrences of `v` with term `s` in `pi`
func (pi Pi) Substitute(dest *Term, v Variable, s Term) {
	if pi.binderVar == v {
		return // won't contain any occurrences of `v` since pi.binderVar binds all free occurrences--which is a shadowed version
	}

	term := new(Term)

	*term = pi.binderType
	pi.binderType.Substitute(term, v, s)
	pi.binderType = (*term).(Type)

	*term = pi.dependent
	pi.dependent.Substitute(term, v, s)
	pi.dependent = (*term).(Type)

	*dest = pi
}

// product class
func (Pi) TypeClassification() typeClassification {
	return productClass
}

// performs beta reduction on product type
func (pi Pi) betaReduce(term Term) Type {
	t := new(Term)
	*t = pi.dependent
	pi.dependent.Substitute(t, pi.binderVar, term)
	return (*t).(Type)
}

type prePi Pi

func ImplicitBind(x Variable, A Type) prePi {
	return prePi{implicit: true, binderVar: x, binderType: A}
}

func Bind(x Variable, A Type) prePi {
	return prePi{binderVar: x, binderType: A}
}

func (pi prePi) To(B Type) func(u Sort) Pi {
	return func(u Sort) Pi {
		pi.dependent = B
		pi.kind = u
		return Pi(pi)
	}
}