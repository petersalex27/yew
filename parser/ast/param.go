package ast

import (
	"fmt"
	err "yew/error"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	"yew/symbol"
)

type Parameter struct {
	paramIndex int
	pattern    ExpressionTypeAnnotation
}

func MakeParameter(paramIndex int, pattern ExpressionTypeAnnotation) Parameter {
	return Parameter{paramIndex, pattern}
}

func (p Parameter) ResolveNames(table *symbol.SymbolTable) bool {
	return p.pattern.ResolveNames(table)
}
func (Parameter) Make(*Parser) bool {
	err.PrintBug()
	panic("")
}
func (p Parameter) GetNodeType() NodeType { return PARAM }
func (p Parameter) Equal_test(a Ast) bool {
	equal := a.GetNodeType() == PARAM
	p2 := a.(Parameter)
	return equal &&
		p.paramIndex == p2.paramIndex &&
		p.pattern.Equal_test(p2.pattern)
}
func (p Parameter) Print(lines []string) {
	lines = printLines(lines)
	fmt.Printf("Parameter\n")
	lines = append(lines, " └─")
	p.pattern.Print(lines)
}

func (p Parameter) Accepts(e Expression) bool {
	return p.pattern.ExpressionType().Equals(e.ExpressionType())
}
