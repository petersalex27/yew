package ast

import (
	fmt "fmt"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	symbol "yew/symbol"
)

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

func (mod Module) Make(p *Parser) bool {
	if valid, e := p.Stack.Validate(moduleRule); !valid {
		e(p.Input).Print()
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

	mod2, ok := ast.(Module)
	return ok && mod.belongsToModule.Equal_test(mod2.belongsToModule) &&
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
func (mod Module) FindStartToken() scan.Token {
	return mod.belongsToModule.token
}
