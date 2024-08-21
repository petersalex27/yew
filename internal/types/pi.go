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
	// x : A
	binderVar Variable
	// B
	dependent Type
	// s
	kind Sort

	Start, End int
}

func (env *Environment) specialize(t Type) Type {
	if fn, ok := t.(functionType); ok {
		return fn.Specialize(env)
	}
	return t
}

func (pi Pi) Specialize(env *Environment) Type {
	if !pi.implicit {
		return pi
	}

	v := env.NextTermHole()
	v.Kind = pi.binderVar.Kind
	term := new(Term)
	pi.dependent.Substitute(term, pi.binderVar, v)
	return env.specialize((*term).(Type))
} 

func (pi Pi) Pos() (start, end int) {
	return pi.Start, pi.End
}

func (pi Pi) GetKind() (Term, Type) {
	if pi.kind == nil {
		pi.kind = Hole(pi.binderVar.x)
	}
	return pi, pi.kind
}

func (pi Pi) Locate(v Variable) bool {
	return pi.binderVar.Kind.Locate(v) || pi.dependent.Locate(v)
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
	mult := ""
	switch pi.binderVar.mult {
	case Erase:
		mult = "erase "
	case Once:
		mult = "once "
	}

 	domain := fmt.Sprintf("%v", pi.binderVar.Kind)
	if pi.binderVar.x == "" || pi.binderVar.x == "_" {
		if pi.implicit {
			if pi.binderVar.x == "_" {
				domain = mult + "_ : " + domain
			}
			domain = "{" + domain + "}"
		} // else, domain is unchanged
	} else {
		domain = fmt.Sprintf("%s%v : ", mult, pi.binderVar) + domain
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
	if pi.binderVar.x == v.x {
		// do not substitute bound variables
		*dest = pi
		return
	}
	term := new(Term)

	//*term = pi.binderType
	pi.binderVar.Kind.Substitute(term, v, s)
	pi.binderVar.Kind = (*term).(Type)

	//*term = pi.dependent
	pi.dependent.Substitute(term, v, s)
	pi.dependent = (*term).(Type)

	*dest = pi
}

func (pi Pi) CollectVariables(m map[string]Variable) map[string]Variable {
	m = pi.binderVar.Kind.CollectVariables(m)
	m = pi.dependent.CollectVariables(m)
	return m
}

// product class
func (Pi) TypeClassification() typeClassification {
	return productClass
}

// returns the terminal type of the product type
//
// for example, given the product type
//		a -> b -> c
// the terminal type is `c`
func (pi Pi) GetTerminal() Type {
	t := pi.dependent
	p, ok := t.(Pi)
	for ok {
		t = p.dependent
		p, ok = t.(Pi)
	}
	return t
}

// performs beta reduction on product type
func (pi Pi) betaReduce(term Term) Type {
	// no replacements if binder is "_" or ""
	if pi.binderVar.x == "" || pi.binderVar.x == "_" {
		return pi.dependent
	}
	
	t := new(Term)
	*t = pi.dependent
	pi.dependent.Substitute(t, pi.binderVar, term)
	return (*t).(Type)
}

type prePi Pi

func ImplicitBind(x Variable, A Type) prePi {
	_ = SetKind(&x, A)
	return prePi{implicit: true, binderVar: x}
}

func Bind(x Variable, A Type) prePi {
	_ = SetKind(&x, A)
	return prePi{binderVar: x}
}

func (pi prePi) To(B Type) func(u Sort) Pi {
	return func(u Sort) Pi {
		pi.dependent = B
		pi.kind = u
		return Pi(pi)
	}
}