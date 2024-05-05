// =================================================================================================
// Alex Peters - March 02, 2024
// =================================================================================================
package types

import "strings"

type Application struct {
	terms []Term
	kind  Type
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
	for _, typ := range app.terms {
		typ.Substitute(&typ, v, s)
	}
}

// application class
func (Application) TypeClassification() typeClassification {
	return applicationClass
}