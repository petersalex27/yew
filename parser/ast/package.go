package ast

import (
	fmt "fmt"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
)

type Package struct {
	belongsToPackage PackageMembership
	program          ProgramTop
}

func MakePackage(id Id, program ProgramTop) Package {
	return Package{belongsToPackage: PackageMembership(id), program: program}
}

func MakePackage2(id scan.IdToken, program ProgramTop) Package {
	return Package{belongsToPackage: PackageMembership(MakeId(id)), program: program}
}

func (pack Package) Make(p *Parser) bool {
	if valid, e := p.Stack.Validate(packageRule); !valid {
		e(p.Input).Print()
		return false
	}
	// analysis

	pack.program = p.Stack.Pop().(ProgramTop)
	pack.belongsToPackage = p.Stack.Pop().(PackageMembership)
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

	pack2, ok := ast.(Package)
	return ok && pack.belongsToPackage.Equal_test(pack.belongsToPackage) &&
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
func (pack Package) ResolveNames(p *Parser) bool {
	panic("TODO")
}
func (pack Package) FindStartToken() scan.Token {
	return pack.belongsToPackage.token
}
