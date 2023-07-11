package ast

import (
	"fmt"
	err "yew/error"
	"yew/ir"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	"yew/type"
	"yew/value"
)

type Value scan.ValueToken

/*func (v Value) Matches(p Pattern) bool {
	if 
	switch v.value.GetType().GetTypeType() {
		case types
	}
}*/

func (v Value) ExpressionType() types.Types {
	return (v.Value).GetType()
}
func (Value) ResolveNames(p *Parser) bool {
	return true
}
func (v Value) DoTypeInference(newTypeInformation types.Types) types.Types {
	if types.ERROR == newTypeInformation.InferType((v.Value).GetType()).GetTypeType() {
		return types.Error{}
	}
	return (v.Value).GetType()
}
func (v Value) Compile(builder *ir.IrBuilder) {
	
}
func (v Value) GetNodeType() NodeType { return VALUE }
func (v Value) Make(*Parser) bool {
	err.PrintBug()
	panic("")
}
func MakeValue(v value.Value) Value {
	return Value(scan.ValueToken{Value: v, Line: 0, Char: 0})
}
func (v Value) Equal_test(a Ast) bool {
	return a.GetNodeType() == VALUE && 
			v.Value.GetType().Equals(a.(Value).Value.GetType()) &&
			v.Value.ToString() == a.(Value).Value.ToString()
}
func (v Value) Print(lines []string) {
	printLines(lines)
	fmt.Printf("Value == %s\n", v.Value.ToString())
}

func (v Value) FindStartToken() scan.Token {
	return scan.ValueToken(v)
}