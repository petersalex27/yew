package types

import (
	"fmt"
	"strings"

	"github.com/petersalex27/yew/internal/common/math"
)

func GetKind[T Term](a *T) Type {
	tmp, t := (*a).GetKind()
	*a = tmp.(T)
	return t
}

func helperSetKind(a Term, A Type) (Term, bool) {
	switch u := (a).(type) {
	case Pi:
		s, ok := A.(Sort)
		if !ok {
			return nil, false
		}
		u.kind = s
		return u, true
	case Lambda:
		pi, ok := A.(Pi)
		if !ok {
			return nil, false
		}
		u.Type = pi
		return u, true
	case Forall:
		panic("bug: cannot set the kind of a 'forall' bound type")
	case Constant:
		panic("bug: cannot set the kind of a constant")
	case Application:
		u.kind = A
		return u, true
	case Variable:
		u.Kind = A
		return u, true
	case Literal:
		_, ku := u.GetKind()
		return u, Equals(ku, A)
	case Universe:
		// TODO: this violates cumulativity
		_, ku := u.GetKind()
		return u, Equals(ku, A)
	default:
		panic("bug: unhandled case in SetKind " + fmt.Sprintf("(%T)", a))
	}
}

func SetKind[T Term](a *T, A Type) bool {
	if a == nil {
		return false
	}
	t, ok := helperSetKind(*a, A)
	*a = t.(T)
	return ok
}

func joinStringed[T fmt.Stringer](elems []T, sep string) string {
	var b strings.Builder
	switch len(elems) {
	case 0:
		return ""
	}

	b.WriteString(elems[0].String())
	for _, elem := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(elem.String())
	}
	return b.String()
}

func calcStartEnd(elems ...positioned) (start, end int) {
	arr := make([]int, 0, len(elems)*2)
	for _, elem := range elems {
		s, e := elem.Pos()
		arr = append(arr, s, e)
	}
	return math.LowHighNon0(arr...)
}

func CalcArity(t Term) (arity uint32) {
	for {
		switch u := t.(type) {
		case Pi:
			if !u.implicit || u.binderVar.mult != Erase {
				arity++ // only count non-erased implicit arguments
			}
			t = u.dependent
		case Lambda:
			arity++
			t = u.bound
		case Forall:
			t = u.body
		default:
			return arity
		}
	}
}

func GetFinal(t Term) Term {
	for {
		switch u := t.(type) {
		case Pi:
			t = u.dependent
		case Lambda:
			t = u.bound
		case Forall:
			t = u.body
		default:
			return t
		}
	}
}

// skipVariableActions returns true if v is in vs.
func vInVs(v Variable, vs []Variable) bool {
	for _, u := range vs {
		if v.Equals(u) {
			return true
		}
	}
	return false 
}

func Locate(v Variable, t Term) bool {
	m := t.CollectVariables(make(map[string]Variable))
	_, found := m[v.String()]
	return found
}


func TypingString(t Term) string {
	if c, isConst := t.(Constant); isConst {
		return c.C
	}
	a, A := t.GetKind()
	return fmt.Sprintf("%v : %v", a, A)
}