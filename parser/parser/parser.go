package parser

import (
	scan "yew/lex"
	"yew/symbol"
	types "yew/type"
)

type Class_ interface { GetClassName() string }
type InstanceFunction_ interface { GetInstanceFunction() InstanceFunction_ }
type ClassTable_ interface {
	InitClassTable() ClassTable_
	Lookup(className string) (class Class_, found bool)
	DeclareClass(p *Parser, newClass Class_) bool
	DeclareInstance(p *Parser, class Class_, instance types.Types) bool
	DefineInstanceFunction(p *Parser, class Class_, instance types.Types, function InstanceFunction_) bool
}

type Parser struct {
	Input     scan.InputStream
	Next      scan.Token
	Current   scan.Token
	Table     *symbol.SymbolTable
	ClassTable ClassTable_
	Stack     *AstStack
	markIndex int
	ParsingClass bool
	ClassVariable types.Tau
	HasConstraint bool
	ClassConstraint types.Constraint
	excludePrelude bool
	//functions []ast.Function
}

// TODO: is this even needed??
func (p *Parser) Reset() {
	p.markIndex = 0
	p.ParsingClass = false 
	p.HasConstraint = false
}

func (p *Parser) setMarkIndex(new int) (old int) {
	old = p.markIndex
	p.markIndex = new
	return
}

func (p *Parser) getMarkIndex() int { return p.markIndex }

func (p *Parser) ExcludePrelude() {
	p.Input.ExcludePrelude()
}

func (p *Parser) Advance() {
	p.Current = p.Next
	p.Next = p.Input.Next()
}
