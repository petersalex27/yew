package ast

import (
	fmt "fmt"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
)

type PackageMembership Id

func (m PackageMembership) Make(p *Parser) bool {
	if valid, e := p.Stack.Validate(packageMembershipRule); !valid {
		e(p.Input).Print()
		return false
	}

	m = PackageMembership(p.Stack.Pop().(Id))
	p.Stack.Push(m)
	return true
}
func (m PackageMembership) GetNodeType() NodeType {
	return PACKAGE_MEMBERSHIP
}
func (m PackageMembership) Equal_test(ast Ast) bool {
	if ast.GetNodeType() != PACKAGE_MEMBERSHIP {
		return false
	}
	m2, ok := ast.(PackageMembership)
	return ok && Id(m).Equal_test(Id(m2))
}
func (m PackageMembership) Print(lines []string) {
	lines = printLines(lines)
	fmt.Printf("Package Membership\n")
	lines = append(lines, " └─")
	Id(m).Print(lines)
}
func (m PackageMembership) ResolveNames(p *Parser) bool {
	panic("TODO")
}
func (m PackageMembership) FindStartToken() scan.Token {
	return m.token
}
