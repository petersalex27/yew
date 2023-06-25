package ast

import (
	"fmt"
	err "yew/error"
	"yew/ir"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	"yew/symbol"
	types "yew/type"
)

type ExpressionTypeAnnotation struct {
	expression     Expression
	expressionType types.Types
}

func (e ExpressionTypeAnnotation) ExpressionType() types.Types {
	return e.expressionType
}
func (e ExpressionTypeAnnotation) ResolveNames(table *symbol.SymbolTable) bool {
	return e.expression.ResolveNames(table)
}
func (e ExpressionTypeAnnotation) DoTypeInference(newTypeInformation types.Types) types.Types {
	return (e.expression).DoTypeInference(newTypeInformation)
}
func (e ExpressionTypeAnnotation) Compile(builder *ir.IrBuilder) {

}
func (ExpressionTypeAnnotation) GetNodeType() NodeType { return TYPE_ANNOTATION }
func (ExpressionTypeAnnotation) Make(*Parser) bool {
	err.PrintBug()
	panic("")
}
func (e ExpressionTypeAnnotation) Equal_test(a Ast) bool {
	equal := a.GetNodeType() == TYPE_ANNOTATION
	if !equal {
		return false
	}
	e2, ok := a.(ExpressionTypeAnnotation)
	if !ok {
		return false
	}
	return equal &&
		e.expression.Equal_test(e2.expression) &&
		checkTypeEqual(e.expressionType, e2.expressionType)
}
func (e ExpressionTypeAnnotation) Print(lines []string) {
	lines = printLines(lines)
	fmt.Printf("Expression :: %s\n", e.expressionType.ToString())
	lines = append(lines, " └─")
	e.expression.Print(lines)
}

func (e ExpressionTypeAnnotation) StackLogString() string {
	return fmt.Sprintf("%s :: %s", e.GetNodeType().ToString(), e.expressionType.ToString())
}

func MakeTypeAnnotation(e Expression, t types.Types) ExpressionTypeAnnotation {
	return ExpressionTypeAnnotation{expression: e, expressionType: t}
}

func (e ExpressionTypeAnnotation) GetExpression() Expression {
	return e.expression
}

func (e ExpressionTypeAnnotation) FindStartToken() scan.Token {
	return e.expression.FindStartToken()
}
