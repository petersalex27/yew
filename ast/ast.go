package ast

//import "yew/ir"
import (
	"fmt"
	err "yew/error"
)

type NodeType int
const (
	PROGRAM NodeType = iota
	EXPRESSION
	EMPTY__
	DEFINITION
	DECLARATION
	VALUE
	ASSIGNMENT
	APPLICATION
	OPERATION
	BOP_
	UOP_
	UOPERATION
	CLASS_DEFINITION
	IDENTIFIER
	LAMBDA
	FUNCTION
	TYPE_ANOTATION
	RETURN
	SEQUENCE
	REPEAT__
	REPEAT_OR_NONE__
)

var ParseRules = map[NodeType][][]NodeType {
	EXPRESSION: {
		{IDENTIFIER},
		{EXPRESSION, OPERATION, EXPRESSION}, 
		{UOPERATION, EXPRESSION},
		{APPLICATION},
		{LAMBDA},
	},
	DECLARATION: { {IDENTIFIER} },
	DEFINITION: { {DECLARATION, ASSIGNMENT} },
	APPLICATION: { {EXPRESSION, EXPRESSION} },
	ASSIGNMENT: { {EXPRESSION} },
}

type AstStack []Ast
type Ast interface {
	//Compile(*ir.IrBuilder)
	Make(*AstStack) bool
	GetNodeType() NodeType
	equal_test(Ast) bool
	print(int)
}
func (stack *AstStack) Push(node Ast) {
	(*stack) = append((*stack), node)
}
func (stack *AstStack) Pop() Ast {
	out := (*stack)[len(*stack) - 1]
	(*stack) = (*stack)[:len(*stack) - 1]
	return out
}
func (stack *AstStack) GetTopNodeType() NodeType {
	return (*stack)[len(*stack) - 1].GetNodeType()
} 
func match(expect NodeType, actual NodeType) bool {
	switch expect {
	case EXPRESSION:
		isExpression := 
				actual == VALUE || 
				actual == EXPRESSION ||
				actual == APPLICATION ||
				actual == LAMBDA ||
				actual == IDENTIFIER ||
				actual == SEQUENCE ||
				actual == EMPTY__ ||
				actual == PROGRAM
		return isExpression
	default:
		return expect == actual
	}
}
//  - allows REPEAT__ NodeType to be used
//  - returns true if matched expected
//  - also returns number of tokens matched for each expected node type 
//	(number for each node type in range [1, ..))
func (stack *AstStack) Validate2(expectedNodes []NodeType) (bool, []int) {
	out := make([]int, 0, 1)
	j := len(expectedNodes) - 1
	if j < 0 {
		return false, out
	}

	for i := len(*stack) - 1; i >= 0 && j >= 0; {
		out = append([]int{0}, out...)
		allow0 := REPEAT_OR_NONE__ == expectedNodes[j]
		doRep := REPEAT__ == expectedNodes[j] || allow0
		if doRep {
			j-- // moves past repeat
			if j < 0 {
				// repeat is first thing in expectedNodes--this is a bug
				err.PrintBug()
				panic("")
			}
			
			// repeat until no match found
			for ; i >= 0; i-- {
				if match(expectedNodes[j], (*stack)[i].GetNodeType()) {
					out[0]++
				} else {
					break
				}
			}
			// when i is not decremented, out[0] == 0; thus, 
			// 	function will return at end of for loop
		} else {
			if match(expectedNodes[j], (*stack)[i].GetNodeType()) {
				out[0]++
			}
			i--
			// when not matched, out[0] == 0; thus, function will return at end of for loop
		}

		// check that match was found
		if out[0] == 0 && !allow0 {
			return false, out
		}
		j--
	} 
	if j >= 0 {
		return false, out // ran out of input to match
	}
	return true, out
}
func (stack *AstStack) Validate(expectedNodes []NodeType) bool {
	if len(*stack) < len(expectedNodes) {
		return false
	}
	
	j := 0
	for i := len(*stack) - len(expectedNodes); i < len(*stack); i++ {
		if !match(expectedNodes[j], (*stack)[i].GetNodeType()) {
			return false
		}
		j++
	}

	return true
}

func printSpaces(n int) {
	for i := 0; i < n; i++ {
		fmt.Printf(" | ")
	}
}

func EqualTest(a Ast, b Ast) bool {
	return a.equal_test(b)
}

func PrintAst(a Ast) {
	a.print(0)
	fmt.Printf("end.\n")
}