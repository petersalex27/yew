package symbol

import (
	"github.com/petersalex27/yew/common/stack"
	"github.com/petersalex27/yew/internal/kind"
)

func (s *solver) Find(tp kind.Typeable) kind.Type {
	ty := tp.AsType(s)
	if v, _, isVar := breakAndReport(ty); isVar {
		if rep, found := s.unifications.Find(v); found {
			return rep
		}
	}
	return ty
}

// 'recursiveCall' controls whether we are looking at the surface of the original type or inside it
func occursHelper(v string, A kind.Type, terms *stack.Stack[kind.Type]) bool {
	_, as, isVar := breakAndReport(A)
	if isVar {
		return false
	}

	terms.PushAll(as...)
	defer terms.Reset() // clear the stack

	for terms.Len() > 0 {
		t, _ := terms.Pop() // should never fail since len > 0
		da, as, isVar := breakAndReport(t)
		if isVar && da == v {
			return true // variable occurs in the type
		}
		terms.PushAll(as...) // push the sub-terms onto the stack
	}

	return false // variable does not occur in the type
}

// checks if a variable 'v' occurs (is contained w/in) a type 'A'
func (s *solver) Occurs(v string, A kind.Type) bool {
	return occursHelper(v, A, s.occursStack)
}

func breakAndReport(a kind.Type) (da string, as []kind.Type, variable bool) {
	if da = a.Constant(); da != "" {
		return da, as, false // constants are not variables
	}

	// variables have no other structure than their name
	da, as = a.Break()
	return da, as, da != "" && len(as) == 0
}

func (s *solver) Union(v string, rep kind.Type) bool {
	// do occurs check, don't allow building of a recursive type from a variable
	if s.Occurs(v, rep) {
		return false
	}

	// map v :-> rep
	s.unifications.Map(v, rep)
	return true
}

// Type unification
//
// ... term unification ... gonna avoid that for now ...
func Unify(s kind.Solver, a, b kind.Typeable) bool {
	ta := s.Find(a)
	tb := s.Find(b)

	da, as, aIsVar := breakAndReport(ta)
	db, bs, bIsVar := breakAndReport(tb)

	if aIsVar {
		return s.Union(da, tb)
	} else if bIsVar {
		return s.Union(db, ta)
	} else if da != db || len(as) != len(bs) {
		return false
	}

	for i := range as {
		if !Unify(s, as[i], bs[i]) {
			return false
		}
	}

	return true
}
