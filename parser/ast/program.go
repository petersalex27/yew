package ast

import (
	"fmt"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	"yew/symbol"
	types "yew/type"
)

// can get any order of definitions and expressions by nesting Programs in a
// program's `expression` member (and possibly leaving definitions empty)
type Program struct {
	statements []Statement
	expression Expression
}

func (p Program) GetProgram() Program {
	return p
}
func (p Program) GetNameSpace() Id {
	return MakeEmptyTypedId(scan.UnderscoreIdToken)
}

func (p Program) ExpressionType() types.Types {
	return p.expression.ExpressionType()
}
func (p Program) ResolveNames(table *symbol.SymbolTable) bool {
	for _, s := range p.statements {
		if !s.ResolveNames(table) {
			return false
		}
	}
	return p.expression.ResolveNames(table)
}
func (p Program) DoTypeInference(newTypeInformation types.Types) types.Types {
	panic("TODO: implement")
}

func (p Program) GetNodeType() NodeType { return PROGRAM }

func (p Program) Make(parser *Parser) bool {
	valid, _ := parser.Stack.TryValidate([]NodeType{EXPRESSION})
	if valid {
		p.expression = parser.Stack.Pop().(Expression)
	} else {
		p.expression = EmptyExpression{}
	}

	//var e err.Error
	p.statements = []Statement{}
	valid, _ = parser.Stack.TryValidate([]NodeType{STATEMENT})
	for valid {
		p.statements = append(p.statements, parser.Stack.Pop().(Statement))
		valid, _ = parser.Stack.TryValidate([]NodeType{STATEMENT})
	}

	// reverse statements
	for i, j := 0, len(p.statements)-1; i < j; i, j = i+1, j-1 {
		p.statements[i], p.statements[j] = p.statements[j], p.statements[i]
	}
	
	parser.Stack.Push(p)
	return true
}

func MakeProgram(ss []Statement, e Expression) Program {
	return Program{statements: ss, expression: e}
}

func (p Program) Equal_test(a Ast) bool {
	equal := a.GetNodeType() == PROGRAM
	if !equal {
		return false
	}

	p2 := a.(Program)
	equal = equal &&
		len(p2.statements) == len(p.statements) &&
		p.expression.Equal_test(p2.expression)
	if !equal {
		return false
	}
	for i, d := range p2.statements {
		equal = equal && p.statements[i].Equal_test(d)
	}
	return equal
}

func (p Program) Print(lines []string) {
	lines = printLines(lines)
	fmt.Printf("Program\n")
	lines = append(lines, " ├─")
	for _, d := range p.statements {
		d.Print(lines)
	}
	lines[len(lines)-1] = " └─"
	p.expression.Print(lines)
}
