package ast

import (
	"yew/symbol"
	//err "yew/error"
	"fmt"
	err "yew/error"
	. "yew/parser/node-type"
	. "yew/parser/parser"
)

type DeclarationQualifier byte

const (
	LetDeclare DeclarationQualifier = iota
	ConstDeclare
	MutDeclare
	ParamDeclare
)

// `let` ID
type Declaration struct {
	Qualifier DeclarationQualifier
	id        Id
}

func (dec Declaration) GetId() Id { return dec.id }

func (dec Declaration) GetQualifier() DeclarationQualifier { return dec.Qualifier }

func MakeDeclaration2(qualifer DeclarationQualifier, id Id) Declaration {
	return Declaration{Qualifier: qualifer, id: id}
}
func MakeDeclaration(id Id) Declaration {
	return MakeDeclaration2(LetDeclare, id)
}
func MakeMutableDeclaration(id Id) Declaration {
	return MakeDeclaration2(MutDeclare, id)
}
func MakeConstantDeclaration(id Id) Declaration {
	return MakeDeclaration2(ConstDeclare, id)
}

func (dec Declaration) GetSymbol() symbol.Symbolic {
	return symbol.MakeSymbol(dec.id.token)
}

func (dec Declaration) ResolveNames(table *symbol.SymbolTable) bool {
	if table.IsDefined(symbol.MakeSymbol(dec.id.token)) {
		return true
	} else if IncludeBuiltin(table, dec.id.token.ToString()) {
		return true
	}

	// TODO: try to link with external source
	panic("TODO: implement linking")
}

// Declaration ::= Identifier Type-Annotation
var declarationRule = NodeRule{
	Production: DECLARATION /* ::= */, Expression: []NodeType{IDENTIFIER, TYPE_ANNOTATION},
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

	dec.id = id
	p.Stack.Push(dec)
	return true
}
func (dec Declaration) GetNodeType() NodeType { return DECLARATION }

func (dec Declaration) Equal_test(a Ast) bool {
	equal := a.GetNodeType() == DECLARATION
	dec2, ok := a.(Declaration)
	equal = equal && ok &&
		dec.Qualifier == dec2.Qualifier &&
		dec.id.Equal_test(dec2.id)
	return equal
}

func (dq DeclarationQualifier) ToString() string {
	switch dq {
	case LetDeclare:
		return "let"
	case ConstDeclare:
		return "const"
	case MutDeclare:
		return "mut"
	case ParamDeclare:
		return "param"
	}
	
	err.PrintBug()
	panic("")
}

func (dec Declaration) Print(lines []string) {
	lines = printLines(lines)
	fmt.Printf("Declaration (%s)\n", dec.Qualifier.ToString())
	lines = append(lines, " └─")
	dec.id.Print(lines)
}
func (dec Declaration) StackLogString() string {
	return fmt.Sprintf("%s; %s", dec.GetNodeType().ToString(), dec.id.token.ToString())
}
