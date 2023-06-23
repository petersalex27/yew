package ast

import (
	"yew/symbol"
	//err "yew/error"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	"fmt"
)

// `let` ID
type Declaration Id

func MakeDeclaration(id Id) Declaration {
	return Declaration(id)
}

func (dec Declaration) GetSymbol() symbol.Symbolic {
	return symbol.MakeSymbol(dec.token)
}

func (dec Declaration) ResolveNames(table *symbol.SymbolTable) bool {
	if table.IsDefined(symbol.MakeSymbol(dec.token)) {
		return true
	} else if IncludeBuiltin(table, dec.token.ToString()) {
		return true
	}

	// TODO: try to link with external source
	panic("TODO: implement linking")
}

// Declaration ::= Identifier Type-Annotation
var declarationRule = NodeRule{
	DECLARATION, /* ::= */ []NodeType{IDENTIFIER, TYPE_ANNOTATION},
}
// rule: pop, transform, push
func (dec Declaration) Make(p *Parser) bool {
	valid, e := p.Stack.Validate(declarationRule)
	if !valid {
		e.Print()
		return false
	}
	annot := p.Stack.Pop().(ExpressionTypeAnnotation)
	id := p.Stack.Pop().(Id)
	if annot.expression.GetNodeType() != EMPTY__ {
		return false
	}
	// try to declare
	e2, decd := p.Table.DeclareLocal(symbol.MakeSymbol(id.token), annot.expressionType)
	if !decd {
		e2.ToError().Print()
		return false
	}

	dec = Declaration(id)
	p.Stack.Push(dec)
	return true
}
func (dec Declaration) GetNodeType() NodeType { return DECLARATION }

func (dec Declaration) Equal_test(a Ast) bool {
	equal := a.GetNodeType() == DECLARATION
	dec2 := a.(Declaration)
	equal = equal && Id(dec).Equal_test(Id(dec2))
	return equal
}

func (dec Declaration) Print(lines []string) {
	lines = printLines(lines)
	fmt.Printf("Declaration\n")
	lines = append(lines, " └─")
	Id(dec).Print(lines)
}
func (dec Declaration) StackLogString() string {
	return fmt.Sprintf("%s; %s", dec.GetNodeType().ToString(), Id(dec).token.ToString())
}
