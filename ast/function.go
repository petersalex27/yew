package ast

import "fmt"

// gives a name to an anonomous function
type Function struct {
	dec Declaration // identifies function
	function Lambda // actual function
}

func (f Function) GetNodeType() NodeType { return FUNCTION }

func (f Function) equal_test(a Ast) bool {	
	equal := a.GetNodeType() == FUNCTION
	f2 := a.(Function)
	return equal &&
		f.dec.equal_test(f2.dec) &&
		f.function.equal_test(f2.function)
}

func (f Function) print(n int) {
	printSpaces(n)
	fmt.Printf("Function\n")
	f.dec.print(n + 1)
	f.function.print(n + 1)
}

func (f Function) Make(stack *AstStack) bool {
	valid := stack.Validate([]NodeType{IDENTIFIER, LAMBDA})
	if !valid {
		return false
	}

	lambda := stack.Pop().(Lambda)
	dec := stack.Pop().(Declaration)
	f.dec = dec
	f.function = lambda
	stack.Push(f)
	return true
}