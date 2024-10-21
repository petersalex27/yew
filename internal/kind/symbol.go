package kind

import "github.com/petersalex27/yew/api"

// represents a symbol with a kind and multiplicity
//
// examples:
//
//	`once mySym : MyType`
//	`f : MyType -> (x : MyType)`
type Symbol struct {
	multiplicity int8
	// Name of the symbol
	name string
	// Type of the symbol
	typ Type
	// [0, 11)
	precedence byte
	// 0: none, 1: infixl, 2: infixr, 3: infix
	nameType byte
	api.Position
}

const (
	None   byte = 0
	Infixl byte = 1
	Infixr byte = 2
	Infix  byte = 3
)

const (
	Erase     int8 = 0
	Once      int8 = 1
	Unlimited int8 = 2
)

func declareWithMult(mult int8, x string, ty Type) Symbol {
	return Symbol{
		multiplicity: min(max(mult, 0), Unlimited), // ensure that the multiplicity is valid
		name:         x,
		typ:          ty,
	}
}

func declare(x api.Token, ty Type) Symbol {
	return declareWithMult(Unlimited, x.String(), ty)
}

func Sym2(name string, ty Type) Symbol {
	return declareWithMult(Unlimited, name, ty)
}

func Sym(name api.Token, ty Type) Symbol {
	return declare(name, ty)
}

func getMode(mode api.Token) (int8, string) {
	switch mode.String() {
	case "once":
		return Once, ""
	case "erase":
		return Erase, ""
	default:
		return Unlimited, UnknownModality
	}
}

func Modal(mode api.Token, name api.Token, ty Type) (symbol Symbol, errorMsg string) {
	var modality int8
	modality, errorMsg = getMode(mode)
	if errorMsg != "" {
		return symbol, errorMsg
	}

	return declareWithMult(modality, name.String(), ty), ""
}

func (s Symbol) Name() string { return s.name }
func (s Symbol) Typing() Type { return s.typ }

func (s Symbol) AsType(Solver) Type { return s.typ }

func (s Symbol) Pipe(into Type) Type {
	return into.Apply(s)
}

func (s Symbol) Apply(a Type) App {
	return Apply(s, a)
}

// if the symbol's kind is a constant, return the name of the constant
func (s Symbol) Constant() string {
	return s.typ.Constant()
}

func (s Symbol) Break() (head string, tail []Type) {
	return s.typ.Break()
}

func (s Symbol) String() string {
	name := s.name
	if name == "" {
		name = "_"
	}

	return "(" + name + " : " + s.typ.String() + ")"
}

func (sym Symbol) UseRepType(s Solver) Symbol {
	sym.typ = s.Find(sym.typ)
	return sym
}

func (sym Symbol) AccountUse() (Symbol, bool) {
	switch sym.multiplicity {
	case Erase:
		return sym, false
	case Once:
		sym.multiplicity = Erase
	}
	return sym, true
}