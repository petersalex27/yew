// =================================================================================================
// Alex Peters - 2024
//
// translates nodes into a version that can be used for type checking and name analysis
// =================================================================================================
package parser

import "github.com/petersalex27/yew/types"

func (id Ident) Translate(parser *Parser) types.Term {
	return types.Constant(id.Name)
}

func (id Ident) TranslateVar(boundSet *map[string]uint) types.Variable {
	return types.VarWith(id.Name, *boundSet)
}

func (im Implicit) Translate(parser *Parser) types.Term {
	panic("TODO: implement")
}

func (app Application) Translate(parser *Parser) types.Term {
	//types.Application
	panic("TODO: implement")
}

func (FunctionType) Translate(parser *Parser) types.Term {
	panic("TODO: implement")
}

func (Lambda) Translate(parser *Parser) types.Term {
	panic("TODO: implement")
}

func (StringConst) Translate(parser *Parser) types.Term {
	panic("TODO: implement")
}

func (FloatConst) Translate(parser *Parser) types.Term {
	panic("TODO: implement")
}

func (CharConst) Translate(parser *Parser) types.Term {
	panic("TODO: implement")
}

func (IntConst) Translate(parser *Parser) types.Term {
	panic("TODO: implement")
}

func (AmbiguousTuple) Translate(parser *Parser) types.Term {
	panic("TODO: implement")
}

func (AmbiguousList) Translate(parser *Parser) types.Term {
	panic("TODO: implement")
}

func (Tuple) Translate(parser *Parser) types.Term {
	panic("TODO: implement")
}

func (List) Translate(parser *Parser) types.Term {
	panic("TODO: implement")
}

func (Key) Translate(parser *Parser) types.Term {
	panic("TODO: implement")
}

func (EnclosedTerm) Translate(parser *Parser) types.Term {
	panic("TODO: implement")
}

func (Listing) Translate(parser *Parser) types.Term {
	panic("TODO: implement")
}

func (ConstrainedType) Translate(parser *Parser) types.Term {
	panic("TODO: implement")
}

func (Marker) Translate(parser *Parser) types.Term {
	panic("bug: found marker, but this should've been removed from the parse stack")
}