package ast

import (
	fmt "fmt"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	symbol "yew/symbol"
)

type PacakgeMembership Id
func (m PacakgeMembership) Make(p *Parser) bool {
	if valid, e := p.Stack.Validate(packageMembershipRule); !valid {
		e(p.Input).Print()
		return false
	}

	m = PacakgeMembership(p.Stack.Pop().(Id))
	p.Stack.Push(m)
	return true
}
func (m PacakgeMembership) GetNodeType() NodeType {
	return PACKAGE_MEMBERSHIP
}
func (m PacakgeMembership) Equal_test(ast Ast) bool {
	if ast.GetNodeType() != PACKAGE_MEMBERSHIP {
		return false
	}
	m2, ok := ast.(PacakgeMembership)
	return ok && Id(m).Equal_test(Id(m2))
}
func (m PacakgeMembership) Print(lines []string) {
	lines = printLines(lines)
	fmt.Printf("Package Membership\n")
	lines = append(lines, " └─")
	Id(m).Print(lines)
}
func (m PacakgeMembership) ResolveNames(table *symbol.SymbolTable) bool {
	panic("TODO")
}
func (m PacakgeMembership) FindStartToken() scan.Token {
	return m.token
}