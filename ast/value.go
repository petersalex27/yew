package ast

import (
	"fmt"
	err "yew/error"
	"yew/ir"
	"yew/symbol"
	"yew/type"
	"yew/value"
)

type Value struct { value value.Value }

func (v Value) ExpressionType() types.Types {
	return (v.value).GetType()
}
func (v Value) ResolveNames(s *symbol.SymbolTable) {
	// TODO
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
func (v Value) Make(*AstStack) bool {
	err.PrintBug()
	panic("")
}
func MakeValue(v value.Value) Value {
	return Value{v}
}
func (v Value) equal_test(a Ast) bool {
	return a.GetNodeType() == VALUE && v.value.GetType().Equals(a.(Value).value.GetType())
}
func (v Value) print(n int) {
	printSpaces(n)
	fmt.Printf("Value{%s}\n", v.value.ToString())
}