package ast

import (
	"fmt"
	err "yew/error"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
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
func (b Binder) ResolveNames(p *Parser) bool {
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