// =================================================================================================
// Alex Peters - March 02, 2024
// =================================================================================================
package types

//import "fmt"

type (
	typeClassification uint16

	Type interface {
		Term
		TypeClassification() typeClassification
	}
)

const (
	waitingClass  typeClassification = 0
	constantClass typeClassification = 1 << iota
	productClass
	applicationClass
	variableClass
	implicitClass
	typeConstantClass
	universeClass
	forallClass
)

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
		return a.x == b.x && a.isHole == b.isHole //&& a.demangler == b.demangler && a.mult == b.mult
	case constantClass, universeClass:
		a := s.String()
		b := t.String()
		return a == b
	case forallClass:
		a := s.(Forall)
		b := t.(Forall)
		return testForallEquals(a, b)
	default:
		return false
	}
}

// not tech. beta equiv., but kinda close enough. Named such to follow form of rules of inference
func (env *Environment) betaEquivalence(s, t Type) bool {
	a, b := env.reduceToWHNF(s), env.reduceToWHNF(t)
	if a == nil || b == nil {
		return false
	}

	if a.(Type).TypeClassification() != b.(Type).TypeClassification() {
		return false
	}
	return Equals(s, t)
}

func (env *Environment) lambdaLeadingApp_reduceToWHNF(l Lambda, app Application) Term {
	if len(app.terms) == 1 {
		return l // in whnf
	}
	// beta reduce
	t := l.betaReduce(app.terms[1])
	T := l.Type.betaReduce(app.terms[1])
	if !SetKind(&t, T) {
		panic("bug: failed to set kind")
	}

	if len(app.terms) == 2 {
		// just a lambda and an argument, now we need to reduce the result
		// NOTE: might cause infinite loop
		return env.reduceToWHNF(t)
	}
	// we have more arguments to apply
	terms := make([]Term, len(app.terms)-2)
	copy(terms, app.terms[2:])
	newApp := MakeApplication(app.kind, terms...)
	return env.reduceToWHNF(newApp)
}

// reduce to weak head normal form, basically, reduce to a form where the head is a constant or
// lambda abstraction
//
// return nil on error
func (env *Environment) reduceToWHNF(s Term) Term {
	if app, ok := s.(Application); ok {
		if len(app.terms) == 0 {
			return app // TODO: i dunno--this is not right, probably
		}

		if _, ok := app.terms[0].(Constant); ok {
			return app // in whnf
		}
		if l, ok := app.terms[0].(Lambda); ok {
			return env.lambdaLeadingApp_reduceToWHNF(l, app)
		}

		// if we have a variable at the head, we can't reduce further
		// if x, ok := app.terms[0].(Variable); ok {
		// 	env.error(VarCannotReduceToWHNF, x)
		// 	return nil
		// }

	}
	// everything else should be in whnf already
	return s
}

func testForallEquals(a, b Forall) bool {
	if len(a.variables) != len(b.variables) {
		return false
	}
	// compare the specialized versions
	// - specialize a with its own variables
	// - specialize b with a's variables

	// A = forall a, b, c . a -> b a -> c
	// B = forall x, y, z . x -> y x -> z
	// C = forall x, y, z . x -> y z -> z
	// B[x:=a, y:=b, z:=c] = a -> b a -> c = A[a:=a, b:=b, c:=c] = a -> b a -> c
	// C[x:=a, y:=b, z:=c] = a -> b c -> c != A[a:=a, b:=b, c:=c] = a -> b a -> c

	for i, v := range a.variables {
		dest := new(Term)
		b.Substitute(dest, b.variables[i], v)
		b = (*dest).(Forall)
	}
	return Equals(a.body, b.body)
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
	if !Equals(a.binderVar, b.binderVar) || !Equals(a.kind, b.kind) {
		return false
	}
	return Equals(a.binderVar.Kind, b.binderVar.Kind) && Equals(a.dependent, b.dependent)
}
