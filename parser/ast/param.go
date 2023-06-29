package ast

import (
	"fmt"
	err "yew/error"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	"yew/symbol"
	types "yew/type"
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

func (par Parameter) MakePatternParam(p *Parser) bool {
	valid, e := p.Stack.Validate(patternParamRule)
	if !valid {
		e(p.Input).Print()
		return false
	}

	var annot ExpressionTypeAnnotation
	expr := p.Stack.Pop().(Expression)
	if expr.GetNodeType() != TYPE_ANNOTATION {
		annot = ExpressionTypeAnnotation{
			expression:     expr,
			expressionType: types.GetNewTau(),
		}
	} else {
		annot = expr.(ExpressionTypeAnnotation)
	}
	par.paramIndex = 0
	par.pattern = annot
	p.Stack.Push(par)
	return true
}
func (p Parameter) GetNodeType() NodeType { return PARAM }
func (p Parameter) Equal_test(a Ast) bool {
	equal := a.GetNodeType() == PARAM
	p2, ok := a.(Parameter)
	return equal && ok &&
		p.paramIndex == p2.paramIndex &&
		p.pattern.Equal_test(p2.pattern)
}
func (p Parameter) Print(ls []string) {
	lines := make([]string, len(ls))
	lines = append(lines, ls...)
	lines = printLines(lines)
	fmt.Printf("Parameter (idx=%d)\n", p.paramIndex)
	lines = append(lines, " └─")
	p.pattern.Print(lines)
}

func (p Parameter) Accepts(e Expression) bool {
	return p.pattern.ExpressionType().Equals(e.ExpressionType())
}

func (p Parameter) FindStartToken() scan.Token {
	return p.pattern.FindStartToken()
}