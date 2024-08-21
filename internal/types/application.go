// =================================================================================================
// Alex Peters - March 02, 2024
// =================================================================================================
package types

import "strings"

type Application struct {
	terms []Term
	kind  Type
	Start, End int
}

func (app Application) GetTerms() []Term {
	return app.terms
}

func (app Application) CollectVariables(m map[string]Variable) map[string]Variable {
	for _, term := range app.terms {
		m = term.CollectVariables(m)
	}
	return app.kind.CollectVariables(m)
}

// MakeCApplication creates an application with a constant as the first term.
func MakeCApplication(kind Type, C Constant, terms ...Term) Application {
	return Application{terms: append([]Term{C}, terms...), kind: kind}
}

// MakeApplication creates an application of terms
func MakeApplication(kind Type, terms ...Term) Application {
	return Application{terms: terms, kind: kind}
}

// Pos returns the start and end position of the application.
func (app Application) Pos() (start, end int) {
	return app.Start, app.End
}

func (app Application) GetKind() (Term, Type) {
	if app.kind == nil {
		app.kind = Hole("")
	}
	return app, app.kind
}

func (a Application) Locate(v Variable) bool {
	for _, term := range a.terms {
		if term.Locate(v) {
			return true
		}
	}
	return false
}

func (a Application) String() string {
	var b strings.Builder
	for _, t := range a.terms {
		b.WriteString(t.String())
		b.WriteByte(' ')
	}
	return strings.TrimRight(b.String(), " ")
}

func (app Application) Substitute(dest *Term, v Variable, s Term) {
	//*dest = 
	a2 := Application{
		terms: make([]Term, len(app.terms)),
		kind: app.kind,
		Start: app.Start,
		End: app.End,
	}
	for i := range a2.terms {
		term := new(Term)
		app.terms[i].Substitute(term, v, s)
		a2.terms[i] = *term
	}
	*dest = a2
}

// application class
func (Application) TypeClassification() typeClassification {
	return applicationClass
}