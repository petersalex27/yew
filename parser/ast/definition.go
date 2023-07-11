package ast

import (
	"fmt"
	//err "yew/error"
	"yew/ir"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	"yew/symbol"
)

// definition <- declaration assignment
type Definition struct {
	assignment Assignment
}

func (def Definition) GetSymbol() symbol.Symbolic {
	return symbol.MakeSymbol(def.assignment.target.token)
}
func (def Definition) ResolveNames(p *Parser) bool {
	return def.assignment.ResolveNames(p)
}

// (Declaration, Expression) -> (Declaration, Definition)
func (def Definition) Make(p *Parser) bool {
	valid, e := p.Stack.Validate(definitionRule)
	if !valid {
		e(p.Input).Print()
		return false
	}
	expr := p.Stack.Pop().(Expression)
	declaration := p.Stack.Peek().(Declaration)
	def.assignment = MakeAssignment(declaration.id, expr)

	p.Stack.Push(def)
	return true
}
func (def Definition) GetNodeType() NodeType { return DEFINITION }

func (def *Definition) Compile(builder *ir.IrBuilder) {

}

func MakeDefinition(a Assignment) Definition {
	return Definition{assignment: a}
}

func (def Definition) Equal_test(a Ast) bool {
	equal := a.GetNodeType() == DEFINITION
	if !equal {
		return false
	}

	def2, ok := a.(Definition)
	equal = equal && ok &&
		def2.assignment.Equal_test(def.assignment)
	return equal
}

func (def Definition) Print(ls []string) {
	lines := make([]string, len(ls))
	lines = append(lines, ls...)
	lines = printLines(lines)
	fmt.Printf("Definition\n")
	lines = append(lines, " └─")
	def.assignment.Print(lines)
}

func (def Definition) StackLogString() string {
	return fmt.Sprintf("%s; %s", 
			def.GetNodeType().ToString(), 
			def.assignment.target.token.ToString())
}

func (def Definition) FindStartToken() scan.Token {
	return def.assignment.target.FindStartToken()
}
