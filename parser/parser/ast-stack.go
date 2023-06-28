package parser

//import "yew/ir"
import (
	//fmt "fmt"
	"fmt"
	"os"
	"strings"
	err "yew/error"
	nodetype "yew/parser/node-type"
	symbol "yew/symbol"
	//nodetype "yew/parser/nodetype"
)

type AstStack []Ast
type Ast interface {
	//Compile(*ir.IrBuilder)
	Make(*Parser) bool
	GetNodeType() nodetype.NodeType
	Equal_test(Ast) bool
	Print([]string)
	ResolveNames(*symbol.SymbolTable) bool
}

type stackMarker int
func (stackMarker) Make(*Parser) bool {
	err.PrintBug()
	panic("")
}
func (stackMarker) GetNodeType() nodetype.NodeType {
	return nodetype.STACK_MARKER
}
func (stackMarker) Equal_test(Ast) bool {
	return false
}
func (stackMarker) Print([]string) {}
func (stackMarker) ResolveNames(*symbol.SymbolTable) bool {
	return false 
}

type StackLoggable interface {
	StackLogString() string
}

func (stack AstStack) LogStack(file *os.File, header string) {
	fmt.Fprintf(file, "%s", header)
	for _, element := range stack {
		var str string
		if logElem, ok := element.(StackLoggable); ok {
			str = logElem.StackLogString()
		} else {
			str = element.GetNodeType().ToString()
		}
		fmt.Fprintf(file, "[%s] ", str)
	}
	if len(stack) == 0 {
		fmt.Fprintf(file, "ø")
	}
	fmt.Fprintf(file, "\n")
}

func NewAstStack() *AstStack {
	stack := new(AstStack)
	*stack = make(AstStack, 0, 0x40)
	return stack
}

// puts an element on the top of the stack
func (stack *AstStack) Push(node Ast) {
	(*stack) = append((*stack), node)
	//stack.LogStack(os.Stdout, "push: ")
}

func (stack *AstStack) Mark(p *Parser) {
	newIndex := len(*stack)
	old := p.setMarkIndex(newIndex)
	stack.Push(stackMarker(old))
}

func (stack *AstStack) Demark(p *Parser) bool {
	out := (*stack)[len(*stack)-1]
	if out.GetNodeType() != nodetype.STACK_MARKER {
		return false
	}
	(*stack) = (*stack)[:len(*stack)-1]
	old := p.setMarkIndex(int(out.(stackMarker)))
	return old == len(*stack)
}

func (stack *AstStack) CutAtMark(p *Parser) (out []Ast, ok bool) {
	idx := p.getMarkIndex()
	if idx < 0 {
		return []Ast{}, false
	}

	if idx + 1 >= len(*stack) {
		out = []Ast{}
	} else {
		tmp := (*stack)[idx + 1:]
		out = make([]Ast, 0, len(tmp))
		out = append(out, tmp...)
	}

	(*stack) = (*stack)[:idx + 1]
	ok = true
	return
}

// removes and returns top element of stack
func (stack *AstStack) Pop() Ast {
	out := (*stack)[len(*stack)-1]

	if out.GetNodeType() == nodetype.STACK_MARKER {
		err.PrintBug()
		panic("")
	}

	(*stack) = (*stack)[:len(*stack)-1]
	//stack.LogStack(os.Stdout, "pop : ")
	return out
}

// returns the top element of the stack but does not remove that element
func (stack *AstStack) Peek() Ast {
	return (*stack)[len(*stack)-1]
}

func (stack *AstStack) GetTopNodeType() nodetype.NodeType {
	return (*stack)[len(*stack)-1].GetNodeType()
}

func IsExpression(actual nodetype.NodeType) bool {
	return actual == nodetype.VALUE ||
		actual == nodetype.EXPRESSION ||
		actual == nodetype.TYPE_ANNOTATION ||
		actual == nodetype.APPLICATION ||
		actual == nodetype.LAMBDA ||
		actual == nodetype.IDENTIFIER ||
		actual == nodetype.SEQUENCE ||
		actual == nodetype.EMPTY__ ||
		actual == nodetype.OPERATION ||
		actual == nodetype.UOPERATION ||
		actual == nodetype.POPERATION ||
		actual == nodetype.PROGRAM ||
		actual == nodetype.TUPLE ||
		actual == nodetype.PATTERN ||
		actual == nodetype.LIST
}

func IsStatement(actual nodetype.NodeType) bool {
	return actual == nodetype.DEFINITION ||
		actual == nodetype.DECLARATION ||
		actual == nodetype.FUNCTION ||
		actual == nodetype.MODULE ||
		actual == nodetype.TYPE_DEF ||
		actual == nodetype.STATEMENT
}

func match(expect nodetype.NodeType, actual nodetype.NodeType) bool {
	switch expect {
	case nodetype.EXPRESSION:
		return IsExpression(actual)
	case nodetype.STATEMENT:
		return IsStatement(actual)
	case nodetype.PROGRAM_TOP:
		isProgramTop := actual == nodetype.PROGRAM || actual == nodetype.MODULE
		return isProgramTop
	default:
		return expect == actual
	}
}

//   - allows REPEAT__ NodeType to be used
//   - returns true if matched expected
//   - also returns number of tokens matched for each expected node type
//     (number for each node type in range [1, ..))
func (stack *AstStack) Validate2(expectedNodes []nodetype.NodeType) (bool, []int) {
	out := make([]int, 0, 1)
	j := len(expectedNodes) - 1
	if j < 0 {
		return false, out
	}

	for i := len(*stack) - 1; i >= 0 && j >= 0; {
		out = append([]int{0}, out...)
		allow0 := nodetype.REPEAT_OR_NONE__ == expectedNodes[j]
		doRep := nodetype.REPEAT__ == expectedNodes[j] || allow0
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

// (true, 0) on success, (false, -1) on stack does not contain enough things to replace,
// (false, n) where n is a non-negative number and n is the index where the rule failed in
// expectedNodes
func (stack *AstStack) TryValidate(expectedNodes []nodetype.NodeType) (bool, int) {
	ruleIndex := len(*stack) - len(expectedNodes)
	if ruleIndex < 0 {
		return false, -1
	}

	j := 0
	for i := ruleIndex; i < len(*stack); i++ {
		if !match(expectedNodes[j], (*stack)[i].GetNodeType()) {
			return false, j // failed at expected of j
		}
		j++
	}

	return true, 0
}

var dummyLocation = err.MakeErrorLocation(0, 0, "", []string{""})

func (stack *AstStack) Validate(rule nodetype.NodeRule) (bool, err.Error) {
	valid, index := stack.TryValidate(rule.Expression)
	if !valid {
		// TODO: need location info!!
		var builder strings.Builder
		builder.WriteString("could not apply production rule:\n")
		builder.WriteString(rule.RuleFailToString(index))
		return false, err.SyntaxError(builder.String(), dummyLocation)
	}

	return true, err.Error{}
}
