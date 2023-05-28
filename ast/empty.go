package ast

import (
	"fmt"
	err "yew/error"
	"yew/symbol"
	"yew/type"
)

type EmptyExpression struct {}

func (e EmptyExpression) GetNodeType() NodeType { return EMPTY__ }
func (e EmptyExpression) Make(*AstStack) bool {
	err.PrintBug()
	panic("")
}

func (e EmptyExpression) ExpressionType() types.Types {
	return types.Tuple{}
}
func (e EmptyExpression) ResolveNames(*symbol.SymbolTable) {
	// TODO
}
func (e EmptyExpression) DoTypeInference(newTypeInformation types.Types) types.Types {
	return e.ExpressionType().InferType(newTypeInformation)
}
func (e EmptyExpression) equal_test(a Ast) bool {
	return a.GetNodeType() == EMPTY__
}
func (e EmptyExpression) print(n int) {
	printSpaces(n)
	fmt.Printf("EmptyExpression ()\n")
}