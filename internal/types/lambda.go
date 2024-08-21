// =================================================================================================
// Alex Peters - May 22, 2024
//
// Type and term level lambda abstractions
// =================================================================================================
package types

import (
	"fmt"
)

type functionType interface {
	Type
	Specialize(env *Environment) Type
	betaReduce(e Term) Type
}

// Lambda represents a lambda abstraction of a term.
//
// For example, given a free variable `x`, a term `t`, and a product `A -> B`
// where `x : A` and `t : B`, x abstracts t creating ...
//
//	\x => t : A -> B
type Lambda struct {
	binder     Variable
	bound      Term
	Type       functionType
	implicit   bool
	Start, End int
}

func (lambda Lambda) ImplicitParam() bool {
	return lambda.implicit
}

func (lambda Lambda) GetKind() (Term, Type) {
	return lambda, lambda.Type
}

// AutoAbstract creates a lambda term from a list of binders and a body term.
//
// For example, given the binders [a, b, c] and the body term t, AutoAbstract
// returns the lambda term
//
//	(\a, b, c => t) : ?a0 -> ?a1 -> ?a2 -> TypeOf t
func AutoAbstract(binders []Variable, body Term) Lambda {
	if len(binders) == 0 {
		panic("bug: AutoAbstract called with no binders")
	}

	var lam Lambda
	var ty Type
	body, ty = body.GetKind()
	for i := len(binders) - 1; i >= 0; i-- {
		t := DummyVar("?t" + fmt.Sprint(i))
		//vt := &VarTyping{Term: binders[i], Kind: s}
		pi := Pi{binderVar: binders[i], dependent: ty, kind: t}
		lam = Lambda{binders[i], body, pi, false, 0, 0}
		body = lam
		ty = pi
	}
	return lam
}

// assumes body has the correct type
//
// for example, given `pi` is
//
//	a -> b -> c
//
// `body` should be
//
//	body : c
//
// panics when body has the wrong type
func (env *Environment) AutoAbstract2(body Term, pi Pi, vs []Variable) (Lambda, bool) {
	//v := DummyVar(fmt.Sprintf("x%d", len(ts)))
	//ts = append(ts, v)
	if !env.Unify(vs[0].Kind, pi.binderVar.Kind) {
		s, e := calcStartEnd(vs[0], pi.binderVar)
		env.unifyingError(vs[0].Kind, pi.binderVar.Kind, s, e)
		return Lambda{}, false
	}

	vs[0].Kind = env.getUnified(vs[0].Kind)
	
	lam := Lambda{
		binder:   vs[0],
		Type:     pi,
		implicit: pi.implicit,
		Start:    pi.Start,
		End:      pi.End,
	}
	if piInner, isPi := pi.dependent.(Pi); isPi {
		body, ok := env.AutoAbstract2(body, piInner, vs[1:])
		if !ok {
			return Lambda{}, false
		}
		lam.bound = body
	} else {
		// must exactly match
		var A Type
		body, A = body.GetKind()
		if !Equals(pi.dependent, A) {
			panic("bug: AutoAbstract2 called with invalid type")
		}
		lam.bound = body
	}
	return lam, true
}

func (lambda Lambda) Pos() (start, end int) {
	return lambda.Start, lambda.End
}

func (lambda Lambda) Substitute(dest *Term, x Variable, e Term) {
	if lambda.binder.x == x.x {
		// do not substitute bound variables
		*dest = lambda
		return
	}

	term := new(Term)
	lambda.bound.Substitute(term, x, e)
	*dest = Lambda{
		binder:   lambda.binder,
		bound:    *term,
		Type:     lambda.Type,
		implicit: lambda.implicit,
		Start:    lambda.Start,
		End:      lambda.End,
	}
}

func (lambda Lambda) CollectVariables(m map[string]Variable) map[string]Variable {
	return lambda.bound.CollectVariables(m)
}

func (lambda Lambda) Locate(v Variable) bool {
	return lambda.bound.Locate(v)
}

func (lambda Lambda) String() string {
	var end string
	if lam, ok := lambda.bound.(Lambda); ok {
		end = lam.continuedString()
	} else {
		end = " => " + lambda.bound.String()
	}
	return fmt.Sprintf("\\%v%v", lambda.binder, end)
}

func (lambda Lambda) continuedString() string {
	var end string
	if lam, ok := lambda.bound.(Lambda); ok {
		end = lam.continuedString()
	} else {
		end = " => " + lambda.bound.String()
	}
	return fmt.Sprintf(", %v%v", lambda.binder, end)
}

func (lambda Lambda) betaReduce(e Term) Term {
	dest := new(Term)
	lambda.bound.Substitute(dest, lambda.binder, e)
	return *dest
}
