// =================================================================================================
// Alex Peters - 2024
// =================================================================================================
package parser

import (
	"fmt"

	"github.com/petersalex27/yew/common/math"
	"github.com/petersalex27/yew/common/table"
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
		int:   i,
		Start: tok.Start,
		End:   tok.End,
	}
}

func (parser *Parser) makeChar(tok token.Token) CharConst {
	c := types.CharConst{}
	ok := c.Parse(tok.Value)
	if !ok {
		panic("bug: could not parse char")
	}

	return CharConst{
		char:  c,
		Start: tok.Start,
		End:   tok.End,
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
		End:   tok.End,
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
		Start:  tok.Start,
		End:    tok.End,
	}
}

// validator for processed constraints
var pv_constraint = processedValidation{
	[]NodeType{tupleType, applicationType, identType},
	func(t Term) string {
		if t.NodeType() == listingType {
			return UnexpectedListingMaybeEnclosed(t)
		}
		return ExpectedConstraint
	},
}

func (typ TypeElem) Parse(parser *Parser) (ok bool) {
	guess := int(math.PowerOfTwoCeil(uint(len(typ.Type))) / 8) // TODO
	parser.locals = table.MakeTable[fmt.Stringer, *Declaration](guess)
	old := parser.parsingTypeSig
	parser.parsingTypeSig = true
	defer func() { parser.parsingTypeSig = old }()

	// parse type
	var ty termElem
	if ty, ok = parser.Process(typingActions, typ.Type); !ok {
		return
	}

	var t Type
	if t, ok = ty.Term.(Type); !ok {
		parser.errorOn(ExpectedType, ty)
		return
	}

	if len(typ.Constraint) == 0 {
		// no constraint to parse, just return
		parser.shift(ty)
		return true
	}

	// parse constraint
	var constraint termElem
	constraint, ok = parser.ProcessAndValidate(constraintActions, typ.Constraint, pv_constraint)
	if !ok {
		return
	}

	var tup Tuple
	if tup, ok = constraint.Term.(Tuple); !ok {
		start, end := constraint.Pos()
		tup = Tuple{[]Term{constraint.Term}, start, end}
	}

	term := termElem{ConstrainedType{tup, t}, ty.termInfo}
	parser.shift(term)
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

	typ, _ := parser.terms.SaveStack.Pop()
	setType(typ)
	return true
}

func (bind BindingElem) Parse(parser *Parser) (ok bool) { return }

func (where WhereClause) Parse(parser *Parser) (ok bool) { return }

func (let LetBindingElem) Parse(parser *Parser) (ok bool) { return }

func (tok TokenElem) Parse(parser *Parser) (ok bool) { return }

func (ee EnclosedElem) Parse(parser *Parser) (ok bool) { return }

func (cons TypeConstructorElem) Parse(parser *Parser) (ok bool) {
	// declare constructor
	return cons.DeclarationElem.Parse(parser)
}

func (data DataTypeElem) Parse(parser *Parser) (ok bool) {
	return //ok = data.TypeConstructor.Parse(parser)
}

func (trait TraitElem) Parse(parser *Parser) (ok bool) { return }

func (inst InstanceElem) Parse(parser *Parser) (ok bool) { return }

func (a AnnotationElem) Parse(parser *Parser) (ok bool) {
	panic("unimplemented")
}

func (m MutualBlockElem) Parse(parser *Parser) (ok bool) {
	panic("unimplemented")
}
