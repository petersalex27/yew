package ast

import (
	"fmt"
	err "yew/error"
	errorgen "yew/parser/error-gen"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	types "yew/type"
)

type Id struct {
	token scan.IdToken
	isInfix bool
	ty    types.Types
}

func MakeIdFromType(ty Type) (Id, errorgen.GenerateErrorFunction) {
	if ty.GetType().GetTypeType() != types.TAU {
		return Id{}, errorgen.ExpectedTypeIdentifier.Generate()
	}
	tau := ty.ty.(types.Tau)
	loc := tau.Loc
	idToken := scan.MakeIdToken(tau.ToString(), loc.GetLine(), loc.GetChar())
	return MakeId(idToken), nil
}

func (id Id) SetType(ty types.Types) Id {
	id.ty = ty
	return id
}
func (id Id) GetName() string {
	return id.token.ToString()
}
func (id Id) Make(p *Parser) bool {
	err.PrintBug()
	panic("")
}
func (id Id) GetNodeType() NodeType { return IDENTIFIER }
func (id Id) ExpressionType() types.Types {
	return id.ty
}
func (id Id) ResolveNames(p *Parser) bool {
	sym := p.Table.Get(id.token.ToString())
	if sym == nil {
		// must be in global scope
		/*e*/
		_, added := p.Table.AddSymbolToGlobal(sym) // declare symbol
		if !added {
			// TODO: print? e.ToError().Print()
			return false
		}
	} // else do nothing
	return true
}
func (id Id) DoTypeInference(newTypeInformation types.Types) types.Types {
	panic("") // TODO
}
func MakeId(token scan.IdToken) Id {
	return Id{token: token, ty: types.GetNewTau(), isInfix: false}
}
func (id Id) SetInfix() Id {
	id.isInfix = true
	return id
}
func MakeIdWithType(token scan.IdToken, ty types.Types) Id {
	return Id{token: token, ty: ty}
}
func MakeEmptyTypedId(token scan.IdToken) Id {
	return MakeIdWithType(token, types.Tuple{})
}

func (id Id) Equal_test(a Ast) bool {
	equal := a.GetNodeType() == IDENTIFIER
	id2, ok := a.(Id)
	if !ok {
		return false
	}
	equal = equal &&
		id2.token.ToString() == id.token.ToString()
	//println(id.token.ToString(), id2.token.ToString())
	if equal {
		return checkTypeEqual(id.ty, id2.ty)
	}
	return false
}

func (id Id) Print(lines []string) {
	printLines(lines)
	fmt.Printf("Id == %s :: %s\n", id.token.ToString(), id.ty.ToString())
}

func (id Id) StackLogString() string {
	return fmt.Sprintf("%s; %s", id.GetNodeType().ToString(), id.token.ToString())
}

func (id Id) FindStartToken() scan.Token {
	return id.token
}
