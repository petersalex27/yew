package ast

import (
	"fmt"
	err "yew/error"
	"yew/ir"
	"yew/symbol"
	"yew/type"
)

type ExpressionTypeAnotation struct {
	expression Expression
	expressionType types.Types
}
func (e ExpressionTypeAnotation) ExpressionType() types.Types {
	return e.expressionType
}
func (e ExpressionTypeAnotation) ResolveNames(*symbol.SymbolTable) {
	return // TODO
}
func (e ExpressionTypeAnotation) DoTypeInference(newTypeInformation types.Types) types.Types {
	return (e.expression).DoTypeInference(newTypeInformation)
}
func (e ExpressionTypeAnotation) Compile(builder *ir.IrBuilder) {
	
}
func (ExpressionTypeAnotation) GetNodeType() NodeType { return TYPE_ANOTATION }
func (ExpressionTypeAnotation) Make(*AstStack) bool {
	err.PrintBug()
	panic("")
}
func (e ExpressionTypeAnotation) equal_test(a Ast) bool {
	equal := a.GetNodeType() == TYPE_ANOTATION
	e2 := a.(ExpressionTypeAnotation)
	return equal && 
			e.expression.equal_test(e2.expression) && 
			e.expressionType.Equals(e2.expressionType)
}
func (e ExpressionTypeAnotation) print(n int) {
	printSpaces(n)
	fmt.Printf("Expression\n")
	e.expression.print(n + 1)
	printSpaces(n + 1)
	fmt.Printf("Type == %s\n", e.expressionType.ToString())
}

func MakeTypeAnotation(e Expression, t types.Types) ExpressionTypeAnotation {
	return ExpressionTypeAnotation{expression: e, expressionType: t}
}