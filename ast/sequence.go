package ast

import (
	"fmt"
	"yew/symbol"
	types "yew/type"
	util "yew/utils"
)

type Sequence []Expression
func (s Sequence) GetNodeType() NodeType { return SEQUENCE }
func (s Sequence) Make(stack *AstStack) bool {
	// continue a sequence
	valid := stack.Validate([]NodeType{SEQUENCE, EXPRESSION})
	if valid {
		e := stack.Pop().(Expression)
		s2 := stack.Pop().(Sequence)
		s2 = append(s2, e)
		stack.Push(s2)
		return true
	}

	// start a new sequence
	valid = stack.Validate([]NodeType{EXPRESSION})
	if !valid {
		return false
	}
	tmp := make(Sequence, 0, 1)
	e := stack.Pop().(Expression)
	tmp = append(tmp, e)
	stack.Push(tmp)
	return true
}
func (s Sequence) ExpressionType() types.Types {
	e, found := util.Tail(s)
	if !found {
		return types.Tuple{} // empty tuple, i.e., ()
	}
	return e.ExpressionType()
}
func (s Sequence) ResolveNames(*symbol.SymbolTable) {
	// TODO
}
func (s Sequence) DoTypeInference(newTypeInformation types.Types) types.Types {
	panic("TODO: implement") // TODO
}

func (s Sequence) equal_test(a Ast) bool {
	equal := a.GetNodeType() == SEQUENCE
	s2 := a.(Sequence)
	equal = equal && len(s2) == len(s)
	if !equal {
		return false
	}

	for i, z := range s2 {
		equal = equal && s[i].equal_test(z)
	}
	return equal
}
func (s Sequence) print(n int) {
	printSpaces(n)
	fmt.Printf("Sequence\n")
	for _, e := range s {
		e.print(n + 1)
	}
}