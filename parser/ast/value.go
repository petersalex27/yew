package ast

import (
	"fmt"
	err "yew/error"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	"yew/ir"
	"yew/symbol"
	"yew/type"
	"yew/value"
)

type Value struct { value value.Value }

/*func (v Value) Matches(p Pattern) bool {
	if 
	switch v.value.GetType().GetTypeType() {
		case types
	}
}*/

func (v Value) ExpressionType() types.Types {
	return (v.value).GetType()
}
func (Value) ResolveNames(*symbol.SymbolTable) bool {
	return true
}
func (v Value) DoTypeInference(newTypeInformation types.Types) types.Types {
	if types.ERROR == newTypeInformation.InferType((v.value).GetType()).GetTypeType() {
		return types.Error{}
	}
	return (v.value).GetType()
}
func (v Value) Compile(builder *ir.IrBuilder) {
	
}
func (v Value) GetNodeType() NodeType { return VALUE }
func (v Value) Make(*Parser) bool {
	err.PrintBug()
	panic("")
}
func MakeValue(v value.Value) Value {
	return Value{v}
}
func (v Value) Equal_test(a Ast) bool {
	return a.GetNodeType() == VALUE && 
			v.value.GetType().Equals(a.(Value).value.GetType()) &&
			v.value.ToString() == a.(Value).value.ToString()
}
func (v Value) Print(lines []string) {
	printLines(lines)
	fmt.Printf("Value == %s\n", v.value.ToString())
}