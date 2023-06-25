package ast

import (
	"fmt"
	err "yew/error"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	"yew/symbol"
	types "yew/type"
)

type Id struct {
	token scan.IdToken
	ty    types.Types
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
func (id Id) ResolveNames(table *symbol.SymbolTable) bool {
	sym := table.Get(id.token.ToString())
	if sym == nil {
		// must be in global scope
		/*e*/
		_, added := table.AddSymbolToGlobal(sym) // declare symbol
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
	return Id{token: token, ty: types.GetNewTau()}
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
