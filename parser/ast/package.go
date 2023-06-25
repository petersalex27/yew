package ast

import (
	fmt "fmt"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	symbol "yew/symbol"
)

// Package-Membership ::= Identifier
var packageMembershipRule = NodeRule{
	Production: PACKAGE_MEMBERSHIP, /* ::= */ Expression: []NodeType{IDENTIFIER},
}

type PacakgeMembership Id
func (m PacakgeMembership) Make(p *Parser) bool {
	if valid, e := p.Stack.Validate(packageMembershipRule); !valid {
		e.Print()
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
	m2 := ast.(PacakgeMembership)
	return Id(m).Equal_test(Id(m2))
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

type Package struct {
	belongsToPackage PacakgeMembership
	program ProgramTop
}

func MakePackage(id Id, program ProgramTop) Package {
	return Package{belongsToPackage: PacakgeMembership(id), program: program}
}

func MakePackage2(id scan.IdToken, program ProgramTop) Package {
	return Package{belongsToPackage: PacakgeMembership(MakeId(id)), program: program}
}

// Package ::= Package-Membership Program
var pacakgeRule = NodeRule{
	Production: PACKAGE, /* ::= */ Expression: []NodeType{PACKAGE_MEMBERSHIP, PROGRAM_TOP},
}

func (pack Package) Make(p *Parser) bool {
	if valid, e := p.Stack.Validate(pacakgeRule); !valid {
		e.Print()
		return false
	}

	pack.program = p.Stack.Pop().(ProgramTop)
	pack.belongsToPackage = p.Stack.Pop().(PacakgeMembership)
	p.Stack.Push(pack)
	return true 
}

func (pack Package) GetNodeType() NodeType {
	return PACKAGE
}
func (pack Package) Equal_test(ast Ast) bool {
	if ast.GetNodeType() != PACKAGE {
		return false
	}

	pack2 := ast.(Package)
	return pack.belongsToPackage.Equal_test(pack.belongsToPackage) &&
			pack.program.Equal_test(pack2.program)
}
func (pack Package) Print(ls []string) {
	lines := make([]string, len(ls))
	lines = append(lines, ls...)
	lines = printLines(lines)
	fmt.Printf("Package\n")
	lines = append(lines, " ├─")
	pack.belongsToPackage.Print(lines)
	lines[len(lines)-1] = " └─"
	pack.program.Print(lines)
}
func (pack Package) ResolveNames(table *symbol.SymbolTable) bool {
	panic("TODO")
}