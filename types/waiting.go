// =================================================================================================
// Alex Peters - March 02, 2024
// =================================================================================================
package types

import "fmt"

type Waiting struct {
	head     Term
	variable Variable
	term     Term
}

func (w Waiting) String() string {
	return fmt.Sprintf("%v[%v:=%v]", w.head, w.variable, w.term)
}

func (w Waiting) Substitute(*Term, Variable, Term) {
	panic("expected term in weak head normal form")
	//return w.head.Substitute(w.variable, w.term).Substitute(v, s)
}

func (Waiting) TypeClassification() typeClassification {
	return waitingClass
}

func (w Waiting) reduce() Term {
	term := new(Term)
	*term = w.head
	w.head.Substitute(term, w.variable, w.term)
	return *term
}

func (w Waiting) Locate(v Variable) bool {
	return w.head.Locate(v) || w.variable.Locate(v) || w.term.Locate(v)
}