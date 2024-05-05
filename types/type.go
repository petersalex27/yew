// =================================================================================================
// Alex Peters - March 02, 2024
// =================================================================================================
package types

type (
	typeClassification byte

	Type interface {
		Term
		TypeClassification() typeClassification
	}
)

const (
	waitingClass     typeClassification = 0
	constantClass    typeClassification = 1 << iota
	productClass
	applicationClass
	variableClass
	implicitClass
	typeConstantClass
	universeClass
)

type Implicit struct {
	Term Term
	Type
}

func (im Implicit) TypeClassification() typeClassification {
	return implicitClass
}

func Equals(s, t Type) bool {
	switch s.TypeClassification() & t.TypeClassification() {
	case applicationClass:
		a := s.(Application)
		b := t.(Application)
		return testApplicationEquals(a, b)
	case productClass:
		a := s.(Pi)
		b := t.(Pi)
		return testProductEquals(a, b)
	case variableClass:
		a := s.(Variable)
		b := t.(Variable)
		return a.demangler == b.demangler && a.mult == b.mult && a.name == b.name
	case constantClass, typeConstantClass, universeClass:
		a := s.String()
		b := t.String()
		return a == b
	case implicitClass:
		a, okA := s.(Implicit)
		b, okB := t.(Implicit)
		return okA && okB && Equals(a.Type, b.Type)
	default:
		return false
	}
}

// not tech. beta equiv., but kinda close enough. Named such to follow form of rules of inference
func betaEquivalence(s, t Type) bool {
	a, b := reduceToWHNF(s), reduceToWHNF(t)
	if a.TypeClassification() != b.TypeClassification() {
		return false
	}
	return Equals(s, t)
}

// reduce to weak head normal form
func reduceToWHNF(s Type) Type {
	if w, ok := s.(Waiting); ok {
		term := new(Term)
		*term = w.head
		w.head.Substitute(term, w.variable, w.term)
		return (*term).(Type)
	}
	return s
}

func testApplicationEquals(a, b Application) bool {
	if !Equals(a.kind, b.kind) {
		return false
	}
	if len(a.terms) != len(b.terms) {
		return false
	}
	for i := range a.terms {
		if !Equals(a.terms[i].(Type), b.terms[i].(Type)) {
			return false
		}
	}
	return true
}

func testProductEquals(a, b Pi) bool {
	if a.binderVar != b.binderVar || a.kind != b.kind {
		return false
	}
	return Equals(a.binderType, b.binderType) && Equals(a.dependent, b.dependent)
}
