package ast

import (
	"fmt"
	"yew/symbol"
	"yew/type"
)

type Program struct {
	definitions []Definition
	expression Expression
}

func (p Program) ExpressionType() types.Types {
	return p.expression.ExpressionType()
}
func (p Program) ResolveNames(table *symbol.SymbolTable) {
	// TODO
}
func (p Program) DoTypeInference(newTypeInformation types.Types) types.Types {
	panic("TODO: implement")
}

func (p Program) GetNodeType() NodeType { return PROGRAM }

func (p Program) Make(stack *AstStack) bool {
	valid, reps := stack.Validate2([]NodeType{DEFINITION, REPEAT_OR_NONE__})
	if valid {
		p.expression = EmptyExpression{} // empty 
	} else {
		valid, reps = stack.Validate2([]NodeType{DEFINITION, REPEAT_OR_NONE__, EXPRESSION})
		if !valid {
			return false
		}
		// reps should be []int{n, 1}, where n >= 0
		// first, pop expression
		p.expression = stack.Pop().(Expression)
	}
	
	// now, pop definitions
	p.definitions = make([]Definition, reps[0])
	for i := reps[0] - 1; i >= 0; i-- {
		p.definitions[i] = stack.Pop().(Definition)
	}
	stack.Push(p)
	return true 
}

func MakeProgram(ds []Definition, e Expression) Program {
	return Program{definitions: ds, expression: e}
}

func (p Program) equal_test(a Ast) bool {
	equal := a.GetNodeType() == PROGRAM
	if !equal {
		return false
	}

	p2 := a.(Program)
	equal = equal && 
			len(p2.definitions) == len(p.definitions) && 
			p.expression.equal_test(p2.expression)
	if !equal {
		return false
	}
	for i, d := range p2.definitions {
		equal = equal && p.definitions[i].equal_test(d)
	}
	return equal
}

func (p Program) print(n int) {
	printSpaces(n)
	fmt.Printf("Program\n")
	for _, d := range p.definitions {
		d.print(n + 1)
	}
	p.expression.print(n + 1)
}