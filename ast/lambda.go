package ast

import (
	"fmt"
	"yew/symbol"
	"yew/type"
)

type Lambda struct {
	binder Id
	bound Expression
}

func (lambda Lambda) GetNodeType() NodeType { return LAMBDA }
func (lambda Lambda) Make(stack *AstStack) bool {
	if !stack.Validate([]NodeType{IDENTIFIER, EXPRESSION}) {
		return false
	}

	lambda.bound = stack.Pop().(Expression)
	lambda.binder = stack.Pop().(Id)
	stack.Push(lambda)
	return true
}
func (lambda Lambda) ExpressionType() types.Types {
	return types.Function{
		Domain: lambda.binder.id.GetType(), 
		Codomain: lambda.bound.ExpressionType(),
	}
}
func (lambda Lambda) ResolveNames(*symbol.SymbolTable) {

}
func (lambda Lambda) DoTypeInference(newTypeInformation types.Types) types.Types {
	panic("") // TODO
}

func (lambda Lambda) equal_test(a Ast) bool {
	equal := a.GetNodeType() == LAMBDA
	l2 := a.(Lambda)
	return equal && 
			lambda.binder.equal_test(l2.binder) &&
			lambda.bound.equal_test(l2.bound)
}

func (l Lambda) print(n int) {
	printSpaces(n)
	fmt.Printf("Lambda\n")
	l.binder.print(n + 1)
	l.bound.print(n + 1)
}