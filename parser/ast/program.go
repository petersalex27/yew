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

// can get any order of definitions and expressions by nesting Programs in a
// program's `expression` member (and possibly leaving definitions empty)
/*type Program struct {
	statements []Statement
	expression Expression
}//*/
type Program []Ast

func (p Program) GetProgram() Program {
	return p
}
func (p Program) GetNameSpace() Id {
	return MakeEmptyTypedId(scan.UnderscoreIdToken)
}

func (p Program) ExpressionType() types.Types {
	if nil == p || len(p) == 0 {
		return EmptyExpression{}.ExpressionType()
	}

	last := p[len(p)-1]
	if IsStatement(last.GetNodeType()) {
		return EmptyExpression{}.ExpressionType()
	} else if IsExpression(last.GetNodeType()) {
		expr := last.(Expression)
		return expr.ExpressionType()
	} else {
		err.PrintBug()
		panic("")
	}
}
func (p Program) ResolveNames(table *symbol.SymbolTable) bool {
	for _, q := range p {
		if !q.ResolveNames(table) {
			return false
		}
	}
	return true
}
func (p Program) DoTypeInference(newTypeInformation types.Types) types.Types {
	panic("TODO: implement")
}

func (p Program) GetNodeType() NodeType { return PROGRAM }

func (p Program) Make(parser *Parser) bool {
	// count len
	tmp, ok := parser.Stack.CutAtMark(parser)
	if !ok {
		return false
	}
	p = Program(tmp)
	parser.Stack.Push(p)

	return true
}

var programStatementRule = NodeRule{
	Production: PROGRAM,
	Expression: []NodeType{STATEMENT},
}

func MakeProgram(as ...Ast) Program { return Program(as) }

func (p Program) Equal_test(a Ast) bool {
	equal := a.GetNodeType() == PROGRAM
	if !equal {
		return false
	}

	p2, ok := a.(Program)
	equal = equal && ok && len(p2) == len(p)
	if !equal {
		return false
	}

	for i := range p {
		if !p[i].Equal_test(p2[i]) {
			return false
		}
	}
	return true
}

func (p Program) Print(lines []string) {
	lines = printLines(lines)
	fmt.Printf("Program\n")
	lines = append(lines, " ├─")
	if p == nil || len(p) == 0 {
		lines[len(lines)-1] = " └─"
		fmt.Printf("ø\n")
	}
	for _, q := range p[:len(p)-1] {
		q.Print(lines)
	}
	lines[len(lines)-1] = " └─"
	p[len(p)-1].Print(lines)
}

func (p Program) FindStartTokenOfPart(i int) scan.Token {
	if p == nil || len(p) <= i {
		return scan.ErrorToken{}
	}
	q := p[i]
	if IsStatement(q.GetNodeType()) {
		return q.(Statement).GetSymbol().GetIdToken()
	} else if IsExpression(q.GetNodeType()) {
		return q.(Expression).FindStartToken()
	}

	return scan.ErrorToken{}
}

func (p Program) FindStartToken() scan.Token {
	return p.FindStartTokenOfPart(0)
}
