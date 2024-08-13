// =================================================================================================
// Alex Peters - 2024
// =================================================================================================
package types

import "strings"

// special case of Pi type where parameters are all variables of unspecified kind
type Forall struct {
	variables  []Variable
	body       Type
	Start, End int
}

func (forall Forall) betaReduce(e Term) Type {
	panic("bug: tried to beta reduce a forall-bound type; this should have been specialized first")
}

func (Forall) TypeClassification() typeClassification {
	return forallClass
}

func (forall Forall) Pos() (start, end int) {
	return forall.Start, forall.End
}

// forall a, b, ... . (T[x:=e])
func (forall Forall) Substitute(dest *Term, x Variable, e Term) {
	term := new(Term)
	*term = forall.body
	forall.body.Substitute(term, x, e)
	*dest = Forall{
		variables: forall.variables,
		body:      (*term).(Type),
		Start:     forall.Start,
		End:       forall.End,
	}
}

// "forall a, b, ... . T"
func (forall Forall) String() string {
	var b strings.Builder
	if len(forall.variables) == 0 {
		return forall.body.String()
	}
	b.WriteString("forall ")
	b.WriteString(forall.variables[0].String())
	for _, v := range forall.variables {
		b.WriteString(", ")
		b.WriteString(v.String())
	}
	b.WriteString(" . ")
	b.WriteString(forall.body.String())
	return b.String()
}

// locate variable `v` in forall.body
func (forall Forall) Locate(v Variable) bool {
	return forall.body.Locate(v)
}

func (forall Forall) CollectVariables(m map[string]Variable) map[string]Variable {
	return forall.body.CollectVariables(m)
}

func (forall Forall) Specialize(env *Environment) Type {
	vs := forall.variables
	var dest *Term
	var f Forall
	f = forall
	for _, v := range vs {
		e := env.NextTermHole()
		// all substitutions shouldn't conflict w/ bound variables
		f.Substitute(dest, v, e)
		f = (*dest).(Forall)
	}
	return env.specialize(f.body)
}

func (forall Forall) GetKind() (Term, Type) {
	var t Type
	var term Term
	term, t = forall.body.GetKind() // get the kind of the forall's body and check if it is bound by forall
	forall.body = term.(Type)
	if s, isSort := t.(Sort); isSort && s.Known() {
		return forall, s // sort does not depend on bound vars
	}

	// sort might depend on bound vars directly or indirectly through the body
	m := forall.body.CollectVariables(make(map[string]Variable))
	m = t.CollectVariables(m)
	rep := func() bool {
		for _, v := range forall.variables {
			if _, found := m[v.String()]; found {
				return true
			}
		}
		// neither the sort nor the body depend on the bound variables
		return false
	}()

	if !rep {
		return forall, t
	}

	vs := make([]Variable, len(forall.variables))
	copy(vs, forall.variables)

	// bind the kind with the forall
	return forall, Forall{
		variables: vs,
		body:      t,
		Start:     forall.Start,
		End:       forall.End,
	}
}
