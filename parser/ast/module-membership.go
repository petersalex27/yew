package ast

import (
	fmt "fmt"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
)

type ModuleMembership Id

func (m ModuleMembership) Make(p *Parser) bool {
	if valid, e := p.Stack.Validate(moduleMembershipRule); !valid {
		e(p.Input).Print()
		return false
	}

	m = ModuleMembership(p.Stack.Pop().(Id))
	p.Stack.Push(m)
	return true
}
func (m ModuleMembership) GetNodeType() NodeType {
	return MODULE_MEMBERSHIP
}
func (m ModuleMembership) Equal_test(ast Ast) bool {
	if ast.GetNodeType() != MODULE_MEMBERSHIP {
		return false
	}
	m2, ok := ast.(ModuleMembership)
	return ok && Id(m).Equal_test(Id(m2))
}
func (m ModuleMembership) Print(ls []string) {
	lines := make([]string, len(ls))
	lines = append(lines, ls...)
	lines = printLines(lines)
	fmt.Printf("Module Membership\n")
	lines = append(lines, " └─")
	Id(m).Print(lines)
}
func (m ModuleMembership) ResolveNames(p *Parser) bool {
	panic("TODO")
}
func (m ModuleMembership) FindStartToken() scan.Token {
	return m.token
}