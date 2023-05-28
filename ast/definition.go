package ast

import (
	"fmt"
	err "yew/error"
	"yew/ir"
	"yew/symbol"
	"yew/type"
)

// `let` ID
type Declaration Id

func MakeDeclaration(sym *symbol.Symbol) Declaration {
	return Declaration(MakeId(sym))
}

type Id struct {
	id *symbol.Symbol
}
func (id Id) Make(*AstStack) bool {
	err.PrintBug()
	panic("")
}
func (id Id) GetNodeType() NodeType { return IDENTIFIER }
func (id Id) ExpressionType() types.Types {
	return id.id.GetType()
}
func (id Id) ResolveNames(*symbol.SymbolTable) {
	// TODO
}
func (id Id) DoTypeInference(newTypeInformation types.Types) types.Types {
	panic("") // TODO
}
func MakeId(s *symbol.Symbol) Id {
	return Id{s}
}

func (id Id) equal_test(a Ast) bool {
	equal := a.GetNodeType() == IDENTIFIER
	id2 := a.(Id)
	equal = equal && 
			id2.id.GetFullName() == id.id.GetFullName() &&
			id2.id.GetType().Equals(id.id.GetType())
	return equal
}

func (id Id) print(n int) {
	printSpaces(n)
	fmt.Printf("Id\n")
	printSpaces(n + 1)
	fmt.Printf("Symbol%s\n", id.id.ToString())
} 

// rule: pop, transform, push
func (dec Declaration) Make(stack *AstStack) bool {
	ok := stack.Validate([]NodeType{IDENTIFIER, TYPE_ANOTATION})
	if !ok {
		return false
	}
	anot := stack.Pop().(ExpressionTypeAnotation)
	id := stack.Pop().(Id)
	if anot.expression.GetNodeType() != EMPTY__ {
		return false
	}
	id.id.SetType(anot.expressionType)
	dec = Declaration(id)
	stack.Push(dec)
	return true
}
func (dec Declaration) GetNodeType() NodeType { return DECLARATION }

func (dec Declaration) equal_test(a Ast) bool {
	equal := a.GetNodeType() == DECLARATION
	dec2 := a.(Declaration)
	equal = equal && 
			dec.id.GetFullName() == dec2.id.GetFullName() &&
			dec.id.GetType().Equals(dec2.id.GetType())
	return equal
}

func (dec Declaration) print(n int) {
	printSpaces(n)
	fmt.Printf("Declaration\n")
	printSpaces(n + 1)
	fmt.Printf("Symbol%s\n", dec.id.ToString())
}

// definition ::= declaration assignment
type Definition struct {
	declaration Declaration
	assignment Assignment
}

func (def Definition) Make(stack *AstStack) bool {
	ok := stack.Validate([]NodeType{DECLARATION, ASSIGNMENT})
	if !ok {
		return false
	}
	def.assignment = stack.Pop().(Assignment)
	def.declaration = stack.Pop().(Declaration)
	stack.Push(def)
	return true
}
func (def Definition) GetNodeType() NodeType { return DEFINITION }


func (def *Definition) Compile(builder *ir.IrBuilder) {
	
}

func MakeDefinition(dec Declaration, a Assignment) Definition {
	return Definition{declaration: dec, assignment: a}
}

func (def Definition) equal_test(a Ast) bool {
	equal := a.GetNodeType() == DEFINITION
	if !equal {
		return false
	}

	def2 := a.(Definition)
	equal = equal && 
			def2.declaration.equal_test(def.declaration) &&
			def2.assignment.equal_test(def.assignment)
	return equal
}

func (def Definition) print(n int) {
	printSpaces(n)
	fmt.Printf("Definition\n")
	def.declaration.print(n + 1)
	def.assignment.print(n + 1)
}