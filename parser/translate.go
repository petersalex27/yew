// =================================================================================================
// Alex Peters - 2024
//
// translates nodes into a version that can be used for type checking and name analysis
// =================================================================================================
package parser

import "github.com/petersalex27/yew/types"

func (id Ident) Translate() types.Term {
	return types.Constant(id.Name)
}

func (id Ident) TranslateVar(boundSet *map[string]uint) types.Variable {
	return types.VarWith(id.Name, *boundSet)
}

func (app Application) Translate() types.Term {
	//types.Application
	panic("TODO: implement")
}

func (FunctionType) Translate() types.Term {
	panic("TODO: implement")
}

func (Lambda) Translate() types.Term {
	panic("TODO: implement")
}

func (StringConst) Translate() types.Term {
	panic("TODO: implement")
}

func (FloatConst) Translate() types.Term {
	panic("TODO: implement")
}

func (CharConst) Translate() types.Term {
	panic("TODO: implement")
}

func (IntConst) Translate() types.Term {
	panic("TODO: implement")
}

func (AmbiguousTuple) Translate() types.Term {
	panic("TODO: implement")
}

func (AmbiguousList) Translate() types.Term {
	panic("TODO: implement")
}

func (Pairs) Translate() types.Term {
	panic("TODO: implement")
}

func (List) Translate() types.Term {
	panic("TODO: implement")
}

func (Key) Translate() types.Term {
	panic("TODO: implement")
}

func (EnclosedTerm) Translate() types.Term {
	panic("TODO: implement")
}
