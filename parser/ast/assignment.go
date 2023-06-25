package ast

import (
	"fmt"
	"yew/ir"
	"yew/symbol"
	. "yew/parser/node-type"
	. "yew/parser/parser"
)

type Assignment struct {
	target     Id
	expression Expression
}

func (a Assignment) Compile(builder *ir.IrBuilder) {

}

func MakeAssignment(target Id, e Expression) Assignment {
	return Assignment{target, e}
}

func (a Assignment) ResolveNames(table *symbol.SymbolTable) bool {
	return a.target.ResolveNames(table) && a.expression.ResolveNames(table)
}

func (a Assignment) GetNodeType() NodeType { return ASSIGNMENT }
// Assignment ::= Identifier Expression
var assignmentRule = NodeRule{
	Production: ASSIGNMENT,
	Expression: []NodeType{IDENTIFIER, EXPRESSION},
}
func (a Assignment) Make(p *Parser) bool {
	valid, e := p.Stack.Validate(assignmentRule)
	if !valid {
		e.Print()
		return false
	}
	a.expression = p.Stack.Pop().(Expression)
	a.target = p.Stack.Pop().(Id)
	p.Stack.Push(a)
	return true
}

func (a Assignment) Equal_test(ast Ast) bool {
	equal := ast.GetNodeType() == ASSIGNMENT
	a2 := ast.(Assignment)
	return equal &&
		a.target.Equal_test(a2.target) &&
		a.expression.Equal_test(a2.expression)
}

func (a Assignment) Print(lines []string) {
	next := make([]string, len(lines))
	next = append(next, lines...)
	next = printLines(next)
	fmt.Printf("Assignment\n")
	next = append(next, " ├─")
	a.target.Print(next)
	next[len(next)-1] = " └─"
	a.expression.Print(next)
}
