package parser

import (
	"yew/lex"
	"yew/symbol"
	types "yew/type"
	//util "yew/utils"
)


func (stack *TypeStack) CreateFunction() bool {
	if len(*stack) < 2 {
		return false
	}
	codomain := stack.Pop()
	domain := stack.Pop()
	stack.Push(types.Function{Domain: domain, Codomain: codomain})
	return true
}

func (stack *TypeStack) CreateTuple(n int) bool {
	if len(*stack) < n {
		return false
	}
	tup := make(types.Tuple, n)
	for ; n > 0 ; n-- {
		tup[n-1] = stack.Pop()
	}
	stack.Push(tup)
	return true
}

type TypeStack []types.Types
func (stack *TypeStack) Push(t types.Types) {
	(*stack) = append((*stack), t)
}
func (stack *TypeStack) Pop() types.Types {
	out := (*stack)[len(*stack)-1]
	(*stack) = (*stack)[:len(*stack)-1]
	return out
}
func (stack *TypeStack) Peek() types.Types {
	return (*stack)[len(*stack)-1]
}
func NewTypeStack() *TypeStack {
	stack := new(TypeStack)
	*stack = make(TypeStack, 0, 0x40)
	return stack
}

type Parser struct {
	Input   scan.InputStream
	Next    scan.Token
	Current scan.Token
	Table   *symbol.SymbolTable
	Stack   *AstStack
	TypeStack *TypeStack
	//functions []ast.Function
}

func InitParser(in scan.InputStream) *Parser {
	p := new(Parser)
	*p = Parser{
		Input: in,
		Table: symbol.InitSymbolTable(in.GetPath()),
		Stack: NewAstStack(),
		TypeStack: NewTypeStack(),
	}
	return p
}

func (p *Parser) Advance() {
	p.Current = p.Next
	p.Next = p.Input.Next()
}