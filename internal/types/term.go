// =================================================================================================
// Alex Peters - February 29, 2024
// =================================================================================================
package types

type Term interface {
	// M[x:=e]
	Substitute(*Term, Variable, Term)
	CollectVariables(m map[string]Variable) map[string]Variable
	String() string
	Locate(Variable) bool
	Pos() (start, end int)
	// returns the term (with possibly updated kind [if it was nil to begin with]) and the kind
	GetKind() (Term, Type)
}

func Split(t Term) (c string, terms []Term) {
	const (
		appc string = "_$"
		pic  string = "_->"
		lamc string = "_λ"
	)
	
	if c, ok := t.(Constant); ok {
		return c.C, []Term{}
	} else if u, ok := t.(Universe); ok {
		return u.String(), []Term{}
	} else if a, ok := t.(Application); ok {
		return appc, a.terms
	} else if p, ok := t.(Pi); ok {
		return pic, []Term{p.binderVar.Kind, p.dependent}
	} else if _, ok := t.(Literal); ok {
		return t.String(), []Term{}
	} else if lam, ok := t.(Lambda); ok {
		return lamc, []Term{lam.binder, lam.bound} // TODO: need binder type??
	} else if v, ok := t.(Variable); ok {
		if v.isHole {
			panic("bug: tried to split holed variable")
		}
		return v.x, []Term{}
	} else {
		panic("bug: tried to split unknown term")
	}
	// else if _, ok := t.(*TermBind); ok {
	// 	panic("let binding occurrence inside type")
	// }

}

func reduce(a Term) (ra Term) {
	// TODO: we assume here that every function is total! this is not true!!!!
	// 	the compiler will break on non-terminating functions :(
	ra = a // TODO: lookup replacements of terms in WHNF
	return
}
