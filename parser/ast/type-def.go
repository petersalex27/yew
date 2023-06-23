package ast

import (
	"fmt"
	err "yew/error"
	scan "yew/lex"
	nodetype "yew/parser/node-type"
	"yew/parser/parser"
	"yew/symbol"
	types "yew/type"
)

type TypeDefinition Id

func (t TypeDefinition) GetIdToken() scan.IdToken {
	return t.token
}
func (t TypeDefinition) GetType() types.Types {
	return t.ty
}
func (t TypeDefinition) SetType(types.Types) symbol.Symbolic {
	// type should already be set 
	return t
}
func (t TypeDefinition) IsDefined() bool {
	return true
}

func (t TypeDefinition) Make(*parser.Parser) bool {
	err.PrintBug()
	panic("")
}

func (t TypeDefinition) GetNodeType() nodetype.NodeType {
	return nodetype.TYPE_DEF
}
func (t TypeDefinition) Equal_test(a parser.Ast) bool {
	ok := a.GetNodeType() == nodetype.TYPE_DEF
	if !ok {
		return false
	}
	t2 := a.(TypeDefinition)
	return Id(t).Equal_test(Id(t2))
}
func (t TypeDefinition) Print(lines []string) {
	printLines(lines)
	fmt.Printf("Type-Definition == %s :: %s\n", t.token.ToString(), t.ty.ToString())
}
func (t TypeDefinition) ResolveNames(table *symbol.SymbolTable) bool {
	// TODO: is this right?
	return Id(t).ResolveNames(table)
}

func (t TypeDefinition) GetSymbol() symbol.Symbolic {
	return t
}

func MakeTypeDefinition(id scan.IdToken, ty types.Types) TypeDefinition {
	return TypeDefinition(MakeIdWithType(id, ty))
}