package ast

import (
	fmt "fmt"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	symbol "yew/symbol"
)

// Package-Membership ::= Identifier
var moduleMembershipRule = NodeRule{
	Production: MODULE_MEMBERSHIP /* ::= */, Expression: []NodeType{IDENTIFIER},
}

type ModuleMembership Id

func (m ModuleMembership) Make(p *Parser) bool {
	if valid, e := p.Stack.Validate(moduleMembershipRule); !valid {
		e.Print()
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
	m2 := ast.(ModuleMembership)
	return Id(m).Equal_test(Id(m2))
}
func (m ModuleMembership) Print(ls []string) {
	lines := make([]string, len(ls))
	lines = append(lines, ls...)
	lines = printLines(lines)
	fmt.Printf("Module Membership\n")
	lines = append(lines, " └─")
	Id(m).Print(lines)
}
func (m ModuleMembership) ResolveNames(table *symbol.SymbolTable) bool {
	panic("TODO")
}

type Module struct {
	belongsToModule ModuleMembership
	program         Program
}

func (m Module) GetProgram() Program { return m.program }
func (m Module) GetNameSpace() Id    { return Id(m.belongsToModule) }

func (m Module) GetSymbol() symbol.Symbolic {
	return symbol.MakeSymbol(m.belongsToModule.token)
}

func MakeModule(id scan.IdToken, program Program) Module {
	return Module{belongsToModule: ModuleMembership(MakeId(id)), program: program}
}

// Module ::= Module-Membership Program
var moduleRule = NodeRule{
	Production: MODULE /* ::= */, Expression: []NodeType{MODULE_MEMBERSHIP, PROGRAM},
}

func (mod Module) Make(p *Parser) bool {
	if valid, e := p.Stack.Validate(moduleRule); !valid {
		e.Print()
		return false
	}

	mod.program = p.Stack.Pop().(Program)
	mod.belongsToModule = p.Stack.Pop().(ModuleMembership)
	p.Stack.Push(mod)
	return true
}

func (mod Module) GetNodeType() NodeType {
	return MODULE
}
func (mod Module) Equal_test(ast Ast) bool {
	if ast.GetNodeType() != MODULE {
		return false
	}

	mod2 := ast.(Module)
	return mod.belongsToModule.Equal_test(mod2.belongsToModule) &&
		mod.program.Equal_test(mod.program)
}
func (mod Module) Print(lines []string) {
	lines = printLines(lines)
	fmt.Printf("Module\n")
	lines = append(lines, " ├─")
	mod.belongsToModule.Print(lines)
	lines[len(lines)-1] = " └─"
	mod.program.Print(lines)
}
func (mod Module) ResolveNames(table *symbol.SymbolTable) bool {
	panic("TODO")
}
