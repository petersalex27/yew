// =================================================================================================
// Alex Peters - March 02, 2024
// =================================================================================================
package types

type Multiplicity byte

const (
	Erase Multiplicity = iota
	Once
	Unrestricted
)

type Variable struct {
	name      string
	demangler uint
	mult      Multiplicity
}

func (v Variable) Use() (_ Variable, ok bool) {
	// cannot use erased
	ok = v.mult != Erase 
	// unrestricted :-> unrestricted
	// linear :-> erase
	// erase :-> erase
	v.mult = v.mult ^ 1
	return v, ok
}

// wildcard variable
var __ Variable = Variable{}

func (v Variable) Locate(u Variable) bool {
	return v.Equals(u)
}

func (Variable) Known() bool { return false }

func demVar(name string, dem uint) Variable {
	return Variable{name, dem, Unrestricted}
}

func mkVar(name string, dem uint, m Multiplicity) Variable {
	return Variable{name, dem, m}
}

func EVar(name string) Variable {
	return mkVar(name, nextEnvUid(), Erase)
}

func LVar(name string) Variable {
	return mkVar(name, nextEnvUid(), Once)
}

func Var(name string) Variable {
	return mkVar(name, nextEnvUid(), Unrestricted)
}

func VarWith(name string, boundSet map[string]uint) Variable {
	if uid, found := boundSet[name]; found {
		return Variable{name: name, demangler: uid}
	}
	return Var(name)
}

func VarBind(name string, boundSet *map[string]uint) (v Variable, unbind func(*map[string]uint)) {
	oldUid, oldFound := (*boundSet)[name]
	if oldFound {
		unbind = func(m *map[string]uint) { (*m)[name] = oldUid }
	} else {
		unbind = func(m *map[string]uint) { delete(*m, name) }
	}

	v = Var(name)
	(*boundSet)[name] = v.demangler
	return v, unbind
}

func (v Variable) Equals(u Variable) bool {
	return v.name == u.name && v.demangler == u.demangler
}

func (v Variable) Substitute(dest *Term, u Variable, s Term) {
	if v == u {
		*dest = s
	}
	*dest = v
}

func (v Variable) String() string { return v.name }

// variable class
func (Variable) TypeClassification() typeClassification {
	return variableClass
}

func (Variable) _identifier_() {}
