package ast

import (
	"fmt"
	scan "yew/lex"
	errorgen "yew/parser/error-gen"
	nodetype "yew/parser/node-type"
	"yew/parser/parser"
	"yew/symbol"
	err "yew/error"
	types "yew/type"
)

type Pattern struct {
	Expression Expression
	Matchers []Lambda
}

var patternRule = nodetype.NodeRule{
	Production: nodetype.PATTERN,
	Expression: []nodetype.NodeType{nodetype.EXPRESSION, nodetype.SEQUENCE},
}
var patternRule2 = nodetype.NodeRule{
	Production: nodetype.PATTERN,
	Expression: []nodetype.NodeType{nodetype.EXPRESSION, nodetype.PROGRAM},
}

var statementError = errorgen.GenerateSyntaxError("unexpected statement inside pattern")
var expressionError = errorgen.GenerateSyntaxError("expected anonymous function")
var emptyPatternError = errorgen.GenerateSyntaxError("cannot have an empty pattern")

func printStatementErrors(p *parser.Parser, statements []Statement) {
	for _, statement := range statements {
		tok := statement.GetSymbol().GetIdToken()
		statementError(tok, p.Input).Print()
	}
}

func makeFromSequence(p *parser.Parser, patStartToken scan.Token, seq Sequence) ([]Lambda, bool) {
	if seq == nil || len(seq) == 0 {
		emptyPatternError(patStartToken, p.Input).Print()
		return []Lambda{}, false
	}

	mat := make([]Lambda, len(seq))
	for i, s := range seq {
		if s.GetNodeType() != nodetype.LAMBDA {
			expressionError(s.FindStartToken(), p.Input).Print()
			return []Lambda{}, false
		}
		
		mat[i] = s.(Lambda)
	}
	return mat, true
}

func (pat Pattern) MakePattern(p *parser.Parser, patStartToken scan.Token) bool {
	if valid, _ := p.Stack.TryValidate(patternRule.Expression); valid {
		seq := p.Stack.Pop().(Sequence)
		expr := p.Stack.Pop().(Expression)
		mat, ok := makeFromSequence(p, patStartToken, seq)
		if !ok {
			return false
		}
		p.Stack.Push(Pattern{Expression: expr, Matchers: mat})
		return true
	}

	valid, e := p.Stack.Validate(patternRule2)
	if !valid {
		e.Print()
		return false
	}

	prog := p.Stack.Pop().(Program)
	expr := p.Stack.Pop().(Expression)
	mat := make([]Lambda, len(prog))
	if prog == nil || len(prog) == 0 {
		emptyPatternError(patStartToken, p.Input).Print()
		return false
	}

	for i, q := range prog {
		if q.GetNodeType() != nodetype.LAMBDA {
			expressionError(prog.FindStartTokenOfPart(i), p.Input).Print()
			return false
		}
		
		mat[i] = q.(Lambda)
	}
	p.Stack.Push(Pattern{Expression: expr, Matchers: mat})
	return true
}

func (pat Pattern) Make(p *parser.Parser) bool {
	err.PrintBug()
	panic("")
}
func (pat Pattern) GetNodeType() nodetype.NodeType {
	return nodetype.PATTERN
}
func (pat Pattern) Equal_test(a parser.Ast) bool {
	if a.GetNodeType() != nodetype.PATTERN {
		return false
	}
	pat2 := a.(Pattern)
	for i := range pat.Matchers {
		if !pat.Matchers[i].Equal_test(pat2.Matchers[i]) {
			return false
		}
	}
	return true
}
func (pat Pattern) Print(lines []string) {
	lines = printLines(lines)
	fmt.Printf("Pattern\n")
	lines = append(lines, " ├─")
	pat.Expression.Print(lines)
	for i := 0; i < len(pat.Matchers)-1; i++ {
		pat.Matchers[i].Print(lines)
	}
	if len(pat.Matchers) > 0 {
		lines[len(lines)-1] = " └─"
		pat.Matchers[len(pat.Matchers)-1].Print(lines)
	}
}
func (pat Pattern) ResolveNames(table *symbol.SymbolTable) bool {
	// TODO
	panic("TODO")
}
func (pat Pattern) ExpressionType() types.Types {
	if nil == pat.Matchers || len(pat.Matchers) == 0 {
		return types.Tuple{} // empty type
	}
	return pat.Matchers[0].ExpressionType()
}
func (pat Pattern) DoTypeInference(newTypeInformation types.Types) types.Types {
	panic("TODO") // TODO
}
func (pat Pattern) FindStartToken() scan.Token {
	return pat.Matchers[0].FindStartToken()
}
