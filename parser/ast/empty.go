package ast

import (
	"fmt"
	err "yew/error"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	"yew/symbol"
	"yew/type"
)

type EmptyExpression struct {}

func (e EmptyExpression) GetNodeType() NodeType { return EMPTY__ }
func (e EmptyExpression) Make(*Parser) bool {
	err.PrintBug()
	panic("")
}

func (e EmptyExpression) ExpressionType() types.Types {
	return types.Tuple{}
}
func (e EmptyExpression) ResolveNames(*symbol.SymbolTable) bool { return true }
func (e EmptyExpression) DoTypeInference(newTypeInformation types.Types) types.Types {
	return e.ExpressionType().InferType(newTypeInformation)
}
func (e EmptyExpression) Equal_test(a Ast) bool {
	return a.GetNodeType() == EMPTY__
}
func (e EmptyExpression) Print(lines []string) {
	printLines(lines)
	fmt.Printf("EmptyExpression ()\n")
}