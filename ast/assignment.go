package ast

import (
	"fmt"
	"yew/ir"
)

type Assignment struct {
	expression Expression
}

func (a Assignment) Compile(builder *ir.IrBuilder) {
	
}

func MakeAssignment(e Expression) Assignment {
	return Assignment{e}
}

func (a Assignment) GetNodeType() NodeType { return ASSIGNMENT }
func (a Assignment) Make(stack *AstStack) bool {
	ok := stack.Validate([]NodeType{EXPRESSION})
	if !ok {
		return false
	}
	a.expression = stack.Pop().(Expression)
	stack.Push(a)
	return true
}

func (a Assignment) equal_test(ast Ast) bool {
	equal := ast.GetNodeType() == ASSIGNMENT
	a2 := ast.(Assignment)
	return equal && a.expression.equal_test(a2.expression)
}

func (a Assignment) print(n int) {
	printSpaces(n)
	fmt.Printf("Assignment\n")
	a.expression.print(n + 1)
} 