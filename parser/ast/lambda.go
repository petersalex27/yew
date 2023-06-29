package ast

import (
	"fmt"
	err "yew/error"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	"yew/symbol"
	types "yew/type"
)

type Binder Parameter

func (b Binder) GetNodeType() NodeType { return BINDER }

// Binder ::= Declaration Expression
var binderRule = NodeRule{
	DECLARATION /* ::= */, []NodeType{DECLARATION, EXPRESSION},
}

func (b Binder) Make(p *Parser) bool {
	if valid, e := p.Stack.Validate(binderRule); !valid {
		e(p.Input).Print()
		return false
	}
	exp := p.Stack.Pop().(Expression)
	dec := p.Stack.Pop().(Declaration)

	b = Binder(Parameter{
		pattern: ExpressionTypeAnnotation{
			expression:     exp,
			expressionType: dec.id.ty,
		},
	})
	p.Stack.Push(b)
	return true
}
func (b Binder) ResolveNames(*symbol.SymbolTable) bool {
	err.PrintBug()
	panic("")
}
func (b Binder) DoTypeInference(newTypeInformation types.Types) types.Types {
	panic("TODO") // TODO
}
func (b Binder) Equal_test(a Ast) bool {
	equal := a.GetNodeType() == BINDER
	b2, ok := a.(Binder)
	return equal && ok && Parameter(b).Equal_test(Parameter(b2))
}
func (b Binder) Print(ls []string) {
	lines := make([]string, len(ls))
	lines = append(lines, ls...)
	lines = printLines(lines)
	fmt.Printf("Binder\n")
	lines = append(lines, " └─")
	Parameter(b).Print(lines)
}
func (b Binder) FindStartToken() scan.Token {
	return Parameter(b).FindStartToken()
}

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
func (lambda Lambda) ResolveNames(table *symbol.SymbolTable) bool {
	return lambda.binder.ResolveNames(table) && lambda.bound.ResolveNames(table)
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
