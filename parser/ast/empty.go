package ast

import (
	"fmt"
	err "yew/error"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	types "yew/type"
	"yew/value"
)

type EmptyExpression struct{
	statement Statement
}

func MakeEmptyExpression(s Statement) EmptyExpression {
	return EmptyExpression{statement: s}
}

func (e EmptyExpression) GetNodeType() NodeType { return EMPTY__ }
func (e EmptyExpression) Make(p *Parser) bool {
	valid, er := p.Stack.Validate(emptyRule)
	if !valid {
		er(p.Input).Print()
		return false
	}
	e.statement = p.Stack.Pop().(Statement)
	p.Stack.Push(e)
	return true 
}

func (e EmptyExpression) ExpressionType() types.Types {
	return types.Tuple{}
}
func (e EmptyExpression) ResolveNames(p *Parser) bool { 
	if e.statement == nil {
		return true
	}
	return e.statement.ResolveNames(p)
}
func (e EmptyExpression) DoTypeInference(newTypeInformation types.Types) types.Types {
	return e.ExpressionType().InferType(newTypeInformation)
}
func (e EmptyExpression) Equal_test(a Ast) bool {
	equal := a.GetNodeType() == EMPTY__
	if !equal {
		return false
	}
	e2, ok := a.(EmptyExpression)
	if !ok {
		return false
	}

	if e.statement == nil {
		return e2.statement == nil
	} else if e2.statement == nil {
		return false
	}

	return e.statement.Equal_test(e2.statement)
}
func (e EmptyExpression) Print(ls []string) {
	name := "EmptyExpression"
	lines := make([]string, len(ls))
	lines = append(lines, ls...)
	lines = printLines(lines)
	if e.statement == nil {
		fmt.Printf("%s ()\n", name)
		return
	}
	lines = append(lines, " └─")
	e.statement.Print(lines)
}
func (e EmptyExpression) FindStartToken() scan.Token {
	if e.statement != nil {
		return e.statement.GetSymbol().GetIdToken()
	}

	return scan.ValueToken{
		Value: value.Tuple{}, 
		Char: err.BuiltinErrorLocation.GetChar(),
		Line: err.BuiltinErrorLocation.GetLine(),
	}
}
