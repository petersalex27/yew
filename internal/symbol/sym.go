package symbol

import "github.com/petersalex27/yew/api"

// represents a symbol with a kind and multiplicity
//
// examples:
//
//	`once mySym : MyType`
//	`f : MyType -> (x : MyType)`
type sym struct {
	// multiplicity of the symbol; or, in the case of undeclared symbols, the negative multiplicity 
	// offset (minus one) of the symbol to be applied to the symbol once it is declared
	//
	// number must be in the range [-3, 3)
	multiplicity int8
	// Name of the symbol
	name string
	// Type of the symbol
	typ api.Type
}

const (
	Erase     int8 = 0
	Once      int8 = 1
	Unlimited int8 = 2
)

func declareWithMult(mult int8, x api.Token, ty api.Type) sym {
	return sym{
		multiplicity: min(max(mult, 0), Unlimited), // ensure that the multiplicity is valid
		name:         x.String(),
		typ:          ty,
	}
}

func declare(x api.Token, ty api.Type) sym {
	return declareWithMult(Unlimited, x, ty)
}
