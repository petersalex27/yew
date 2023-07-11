package parser

//import "yew/ir"
import (
	//fmt "fmt"
	"fmt"
	"os"
	"strings"
	err "yew/error"
	scan "yew/lex"
	nodetype "yew/parser/node-type"
	//nodetype "yew/parser/nodetype"
)

type AstStack []Ast
type Ast interface {
	//Compile(*ir.IrBuilder)
	Make(*Parser) bool
	GetNodeType() nodetype.NodeType
	Equal_test(Ast) bool
	Print([]string)
	ResolveNames(*Parser) bool
	FindStartToken() scan.Token
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
func (stackMarker) ResolveNames(*Parser) bool {
	return false
}
func (stackMarker) FindStartToken() scan.Token {
	return scan.MakeBlankToken()
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

	if idx+1 >= len(*stack) {
		out = []Ast{}
	} else {
		tmp := (*stack)[idx+1:]
		out = make([]Ast, 0, len(tmp))
		out = append(out, tmp...)
	}

	(*stack) = (*stack)[:idx+1]
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

// assumes len(stack) >= n
func (stack *AstStack) typeTopN(n int) ([]nodetype.NodeType, bool) {
	ruleIndex := len(*stack) - n
	if ruleIndex < 0 {
		return []nodetype.NodeType{}, false
	}

	nts := make([]nodetype.NodeType, n)
	for i, j := n-1, len(*stack)-1; i >= 0; i, j = i-1, j-1 {
		nts[i] = (*stack)[j].GetNodeType()
	}
	return nts, true
}

func stringNodeTypes(nodetypes []nodetype.NodeType) string {
	var builder strings.Builder
	builder.WriteByte('(')
	for _, nt := range nodetypes {
		builder.WriteString(nt.ToString())
		builder.WriteString(", ")
	}
	builder.WriteByte(')')
	return builder.String()
}

// (true, 0) on success, (false, -1) on stack does not contain enough things to replace,
// (false, n) where n is a non-negative number and n is the index where the rule failed in
// expectedNodes
func (stack *AstStack) TryValidate(expectedNodes []nodetype.NodeType) (bool, int) {
	nts, ok := stack.typeTopN(len(expectedNodes))
	if !ok {
		return false, -1
	}

	/*
	fmt.Fprintf(os.Stderr, "expected: %s\n", stringNodeTypes(expectedNodes))
	fmt.Fprintf(os.Stderr, "actual: %s\n", stringNodeTypes(nts))//*/

	for i := 0; i < len(nts); i++ {
		if !match(expectedNodes[i], nts[i]) {
			return false, i 
		}
	}
	return true, 0
}

func (stack *AstStack) failedValidation(rule nodetype.NodeRule, index int) func(in scan.InputStream) err.Error {
	var line, char int
	var builder strings.Builder
	//fmt.Fprintf(os.Stderr, "index=%d\n", index)
	if len(*stack) <= index {
		line = 0
		char = 0
		builder.WriteString(rule.ExpectedFailure(index))
	} else if index >= 0 {
		ruleExpressionLen := len(rule.Expression)
		ast := (*stack)[len(*stack)-ruleExpressionLen+index]
		nt := ast.GetNodeType()
		loc := ast.FindStartToken().GetLocation()
		line = loc.GetLine()
		char = loc.GetChar()
		builder.WriteString(nodetype.FoundFailure(nt))
		builder.WriteString(" but ")
		builder.WriteString(rule.ExpectedFailure(index))
	} else { // index is negative
		builder.WriteString("tried to parse ")
		article, s := nodetype.GetErrorName(rule.Production)
		builder.WriteString(article)
		builder.WriteByte(' ')
		builder.WriteString(s)
		builder.WriteString(", expecting the sequence ")
		for i, expr := range rule.Expression {
			_, s = nodetype.GetErrorName(expr)
			builder.WriteString(s)
			if i < len(rule.Expression)-1 {
				builder.WriteString(", ")
			}
		}
		builder.WriteString(". The following sequence was found instead: ")
		ln := len(*stack)
		for i := 0; i < ln; i++ {
			_, s = nodetype.GetErrorName((*stack)[i].GetNodeType())
			builder.WriteString(s)
			if i < len(rule.Expression)-1 {
				builder.WriteString(", ")
			}
		}
	}

	str := builder.String()

	return func(in scan.InputStream) err.Error {
		loc := err.MakeErrorLocation(line, char, in.GetPath(), in.GetSource())
		return err.SyntaxError(str, loc)
	}
}

func (stack *AstStack) Validate(rule nodetype.NodeRule) (bool, func(in scan.InputStream) err.Error) {
	valid, index := stack.TryValidate(rule.Expression)
	if !valid {
		return false, stack.failedValidation(rule, index)
	}

	return true, nil
}
