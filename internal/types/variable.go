// =================================================================================================
// Alex Peters - March 02, 2024
// =================================================================================================
package types

import (
	"fmt"
	"strings"
)

type Multiplicity byte

const (
	Unrestricted Multiplicity = iota
	Erase
	Once
)

type Variable struct {
	isHole     bool
	x          string
	demangler  uint
	mult       Multiplicity
	Kind       Type
	Start, End int
}

func (v Variable) SetMultiplicity(m Multiplicity) Variable {
	v.mult = m
	return v
}

func (v Variable) CollectVariables(m map[string]Variable) map[string]Variable {
	m[v.x] = v
	return m
}

func (v Variable) GetKind() (Term, Type) {
	if v.Kind == nil {
		v.Kind = Hole(v.x)
	}
	return v, v.Kind
}

type positioned interface {
	Pos() (start, end int)
}
type strPos = interface {
	fmt.Stringer
	positioned
}

func (v Variable) Pos() (start, end int) {
	return v.Start, v.End
}

func (env *Environment) Use(v *Variable) (ok bool) {
	// cannot use erased
	mult := (*v).mult
	// unrestricted :-> unrestricted
	// linear :-> erase
	// erase :-> erase
	(*v).mult = mult ^ 1
	if ok = mult != Erase; !ok {
		env.error(multiplicityPreventsUse(*v), *v)
		return false
	}
	return ok
}

// wildcard variable
var __ Variable

func init() {
	// init wildcard var
	name := Constant{C: "_"}
	__ = MakeVar(name, 0, Erase, nil)
}

// wildcard variable
func Wildcard() Variable {
	return Var(__)
}

func (v Variable) Locate(u Variable) bool {
	return v.Equals(u)
}

func (Variable) Known() bool { return false }

func Hole(s string) Variable {
	v := DummyVar(s)
	v.isHole = true
	v.mult = Erase
	t := v.Kind.(Variable)
	t.isHole = true
	t.mult = Erase
	v.Kind = t
	return v
}

func typingChain(s strPos, m Multiplicity, Ts ...Term) (v Variable) {
	if len(Ts) == 0 {
		return Var(s)
	}
	v = MakeVar(s, nextEnvUid(), m, nil)
	for i := len(Ts) - 1; i > 0; i-- {
		T, ok := Ts[i-1].(Type)
		if !ok {
			panic("bug: expected type")
		}
		var A Type
		if A, ok = Ts[i].(Type); !ok {
			panic("bug: expected type")
		}
		if !SetKind(&T, A) {
			panic("bug: failed to set kind")
		}
		Ts[i-1] = T
	}
	v.Kind = Ts[0].(Type)
	return v
}

func processDummyName(s string) string {
	s = strings.TrimLeft(s, "?")
	s = strings.TrimPrefix(s, "t.")
	s = strings.TrimPrefix(s, "T.")
	s = strings.TrimRight(s, "0123456789")
	return s
}

func preprocessed_dummyVarKind(s string, uid uint) Variable {
	name := "?T." + s + fmt.Sprint(uid)
	return Variable{true, name, uid, Unrestricted, nil, 0, 0}
}

func preprocessed_dummyVarType(s string) Variable {
	uid := nextEnvUid()
	name := "?t." + s + fmt.Sprint(uid)
	return Variable{true, name, uid, Unrestricted, nil, 0, 0}
}

func dummyVarType(s string, withKind bool) Variable {
	s = processDummyName(s)
	v := preprocessed_dummyVarType(s)
	if withKind {
		v.Kind = preprocessed_dummyVarKind(s, v.demangler)
	}
	return v
}

func DummyVar(s string) Variable {
	uid := nextEnvUid()
	s = processDummyName(s)
	ty := preprocessed_dummyVarType(s)
	v := Variable{true, "?" + s + fmt.Sprint(uid), uid, Unrestricted, ty, 0, 0}
	return v
}

func MakeVar(name strPos, dem uint, m Multiplicity, A Type) (v Variable) {
	v.x = name.String()
	v.demangler = dem
	v.mult = m
	v.Kind = A
	v.Start, v.End = name.Pos()
	return
}

func ErasedVar(name strPos) Variable {
	v := DummyVar(name.String())
	v.Start, v.End = name.Pos()
	v.mult = Erase
	return v
}

func LinearizedVar(name strPos) Variable {
	v := DummyVar(name.String())
	v.Start, v.End = name.Pos()
	v.mult = Once
	return v
}

func Var(name strPos) Variable {
	uid := nextEnvUid()
	v := Variable{x: name.String(), demangler: uid, mult: Unrestricted}
	v.Start, v.End = name.Pos()
	v.Kind = dummyVarType(v.x, false)
	return v
}

func (v Variable) Equals(u Variable) bool {
	return v.x == u.x //&& v.demangler == u.demangler
}

func (v Variable) Substitute(dest *Term, u Variable, s Term) {
	if v.Equals(u) {
		*dest = s
		return
	}
	*dest = v
}

func (v Variable) String() string {
	return v.x
}

// variable class
func (Variable) TypeClassification() typeClassification {
	return variableClass
}
