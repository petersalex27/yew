package ast

import (
	"fmt"
	//err "yew/error"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	"yew/ir"
	"yew/symbol"
)

// definition <- declaration assignment
type Definition struct {
	assignment Assignment
}

func (def Definition) GetSymbol() symbol.Symbolic {
	return symbol.MakeSymbol(def.assignment.target.token)
}
func (def Definition) ResolveNames(table *symbol.SymbolTable) bool {
	return def.assignment.ResolveNames(table)
}

// Definition ::= Declaration Expression
var definitionRule = NodeRule{
	DEFINITION, /* ::= */ []NodeType{DECLARATION, EXPRESSION},
}
// (Declaration, Expression) -> (Declaration, Definition)
func (def Definition) Make(p *Parser) bool {
	valid, e := p.Stack.Validate(definitionRule)
	if !valid {
		e.Print()
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

	def2 := a.(Definition)
	equal = equal &&
		def2.assignment.Equal_test(def.assignment)
	return equal
}

func (def Definition) Print(lines []string) {
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
