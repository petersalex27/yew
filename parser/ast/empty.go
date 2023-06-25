package ast

import (
	"fmt"
	err "yew/error"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	"yew/symbol"
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
	valid, er := p.Stack.Validate(NodeRule{EMPTY__, []NodeType{STATEMENT}})
	if !valid {
		er.Print()
		return false
	}
	e.statement = p.Stack.Pop().(Statement)
	p.Stack.Push(e)
	return true 
}

func (e EmptyExpression) ExpressionType() types.Types {
	return types.Tuple{}
}
func (e EmptyExpression) ResolveNames(table *symbol.SymbolTable) bool { 
	if e.statement == nil {
		return true
	}
	return e.statement.ResolveNames(table)
}
func (e EmptyExpression) DoTypeInference(newTypeInformation types.Types) types.Types {
	return e.ExpressionType().InferType(newTypeInformation)
}
func (e EmptyExpression) Equal_test(a Ast) bool {
	return a.GetNodeType() == EMPTY__
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
		Index: err.BuiltinErrorLocation.GetSourceIndex(),
		Char: err.BuiltinErrorLocation.GetChar(),
		Line: err.BuiltinErrorLocation.GetLine(),
	}
}
