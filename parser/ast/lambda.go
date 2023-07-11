package ast

import (
	"fmt"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	types "yew/type"
)

type Lambda struct {
	binder Parameter
	bound  Expression
}

func (lambda Lambda) GetNodeType() NodeType { return LAMBDA }

func (lambda Lambda) Make2(p *Parser) bool {
	if valid, e := p.Stack.Validate(lambdaRule2); !valid {
		e(p.Input).Print()
		return false
	}

	lambda.bound = p.Stack.Pop().(Expression)
	lambda.binder = p.Stack.Pop().(Parameter)
	p.Stack.Push(lambda)
	return true
}
func (lambda Lambda) Make(p *Parser) bool {
	valid, _ := p.Stack.TryValidate(lambdaRule.Expression)
	if valid {
		lambda.bound = p.Stack.Pop().(Expression)
		lambda.binder = Parameter(p.Stack.Pop().(Parameter))
		p.Stack.Push(lambda)
		return true
	}

	return lambda.Make2(p)
}
func MakeLambda(p Parameter, e Expression) Lambda {
	return Lambda{binder: p, bound: e}
}
func (lambda Lambda) ExpressionType() types.Types {
	return types.Function{
		Domain:   lambda.binder.pattern.ExpressionType(),
		Codomain: lambda.bound.ExpressionType(),
	}
}
func (lambda Lambda) ResolveNames(p *Parser) bool {
	return lambda.binder.ResolveNames(p) && lambda.bound.ResolveNames(p)
}
func (lambda Lambda) DoTypeInference(newTypeInformation types.Types) types.Types {
	panic("") // TODO
}

func (lambda Lambda) Equal_test(a Ast) bool {
	equal := a.GetNodeType() == LAMBDA
	l2, ok := a.(Lambda)
	return equal && ok && 
		lambda.binder.Equal_test(l2.binder) &&
		lambda.bound.Equal_test(l2.bound)
}

func (l Lambda) Print(ls []string) {
	lines := make([]string, len(ls))
	lines = append(lines, ls...)
	lines = printLines(lines)
	fmt.Printf("Lambda\n")
	lines = append(lines, " ├─")
	l.binder.Print(lines)
	lines[len(lines)-1] = " └─"
	l.bound.Print(lines)
}

func (l Lambda) FindStartToken() scan.Token {
	return l.binder.pattern.expression.FindStartToken()
}
