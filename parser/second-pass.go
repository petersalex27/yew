// =================================================================================================
// Alex Peters - 2024
// =================================================================================================
package parser

import (
	"github.com/petersalex27/yew/token"
	"github.com/petersalex27/yew/types"
)

// func (parser *Parser) process(tokens []token.Token, shiftSet map[token.Type]func(token.Token)Term) (remainder []token.Token) {
// 	for i, tok := range tokens {
// 		if f, found := shiftSet[tok.Type]; found {
// 			parser.shift(f(tok))
// 		} else {
// 			return tokens[i:]
// 		}
// 	}
// 	return nil
// }

// var typeShiftSet = map[token.Type]func(token.Token)Term{
// 	token.Id: func(t token.Token) Term {
// 		return Ident{t.Value, t.Start, t.End}
// 	},
// 	token.Affixed: func(t token.Token) Term {

// 	},
// }

//func (parser *Parser) parseConstraint() 

func makeIdent(t token.Token) Ident {
	return Ident{t.Value, t.Start, t.End}
}

func (parser *Parser) makeInt(tok token.Token) IntConst {
	i := types.IntConst{}
	ok := i.Parse(tok.Value)
	if !ok {
		panic("bug: could not parse integer")
	}
 
	return IntConst{
		int: i,
		Start: tok.Start,
		End: tok.End,
	}
}

func (parser *Parser) makeChar(tok token.Token) CharConst {
	c := types.CharConst{}
	ok := c.Parse(tok.Value)
	if !ok {
		panic("bug: could not parse char")
	}

	return CharConst{
		char: c,
		Start: tok.Start,
		End: tok.End,
	}
}

func (parser *Parser) makeFloat(tok token.Token) FloatConst {
	f := types.FloatConst{}
	ok := f.Parse(tok.Value)
	if !ok {
		panic("bug: could not parse float")
	}

	return FloatConst{
		float: f,
		Start: tok.Start,
		End: tok.End,
	}
}

func (parser *Parser) makeString(tok token.Token) StringConst {
	s := types.StringConst{}
	ok := s.Parse(tok.Value)
	if !ok {
		panic("bug: could not parse string")
	}

	return StringConst{
		string: s,
		Start: tok.Start,
		End: tok.End,
	}
}

// func (parser *Parser) shiftTypeTerm(tok token.Token) (ok bool) {
// 	switch tok.Type {
// 	case token.Id:
// 		decl, found := parser.lookupTerm(tok)
// 		if found {
// 			prec, _ := decl.precedence, decl.rAssoc
// 			parser.shift(decl.name)
// 			parser.shiftAction(ApplyAction(prec))
// 			return true
// 		}
// 		fallthrough
// 	case token.ImplicitId:
// 		parser.shift(makeIdent(tok))
// 		return true 
// 	case token.Affixed:
// 		decl, found := parser.lookupTerm(tok)
// 		if found {
// 			prec, rAssoc := decl.precedence, decl.rAssoc
// 			parser.shift(decl.name)
// 			parser.shiftAction(InfixAction(prec, rAssoc))
// 			return
// 		}
// 		ok = false
// 		parser.error2(UnknownIdent, tok.Start, tok.End)
// 		return
// 	case token.CharValue:
// 		ok = true
// 		parser.shift(parser.makeChar(tok))
// 	case token.IntValue:
// 		ok = true
// 		parser.shift(parser.makeInt(tok))
// 	case token.FloatValue:
// 		ok = true
// 		parser.shift(parser.makeFloat(tok))
// 	case token.StringValue:
// 		ok = true
// 		parser.shift(parser.makeString(tok))
// 	default:
// 		ok = false
// 		parser.error2(UnexpectedToken, tok.Start, tok.End)
// 	}
// 	return
// }

func (typ TypeElem) Parse(parser *Parser) (ok bool) {
	// var prec uint8 = 0
	// var rAssoc bool = false
	// // TODO: parse constraint
	// if len(typ.Type) == 0 {
	// 	panic("bug: type has no associated tokens")
	// }

	// ok, prec, rAssoc = parser.shiftTypeTerm(typ.Type[0])
	// if !ok {
	// 	return
	// }

	// parser.act()

	// for _, tok := range typ.Type[1:] {
	// 	ok, prec, rAssoc = parser.shiftTypeTerm(tok)
	// 	if !ok {
	// 		return
	// 	}

		
	// }
	return
}

func ToTyping(decl Declaration) (ty types.Typing) {
	return types.Typing{
		Term: types.Constant(decl.name.Name),
		Kind: decl.typing.(types.Type),
	}
}

func (dec DeclarationElem) Parse(parser *Parser) (ok bool) {
	var setType func(Term, ...uint8)
	setType, ok = parser.declare(dec.Name)
	ok = ok && dec.Typing.Parse(parser)
	if !ok {
		return
	}

	typ, _ := parser.terms.Pop()
	setType(typ)
	decl, _ := parser.lookupTerm(dec.Name)
	return parser.env.Declare(ToTyping(decl))
}

func (bind BindingElem) Parse(parser *Parser) (ok bool) {return}

func (where WhereClause) Parse(parser *Parser) (ok bool) {return}

func (let LetBindingElem) Parse(parser *Parser) (ok bool) {return}

func (tok TokenElem) Parse(parser *Parser) (ok bool) {return}

func (ee EnclosedElem) Parse(parser *Parser) (ok bool) {return}

func (cons TypeConstructorElem) Parse(parser *Parser) (ok bool) {
	// declare constructor
	return cons.DeclarationElem.Parse(parser)
}

func (data DataTypeElem) Parse(parser *Parser) (ok bool) {
	return//ok = data.TypeConstructor.Parse(parser)
}

func (trait TraitElem) Parse(parser *Parser) (ok bool) {return}

func (inst InstanceElem) Parse(parser *Parser) (ok bool) {return}

func (a AnnotationElem) Parse(parser *Parser) (ok bool) {
	panic("unimplemented")
}

func (m MutualBlockElem) Parse(parser *Parser) (ok bool) {
	panic("unimplemented")
}