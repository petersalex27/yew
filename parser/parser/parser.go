package parser

import (
	scan "yew/lex"
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
	for ; n > 0; n-- {
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
	//functions []ast.Function
}

func (p *Parser) setMarkIndex(new int) (old int) {
	old = p.markIndex
	p.markIndex = new
	return
}

func (p *Parser) getMarkIndex() int { return p.markIndex }

func InitParser(in scan.InputStream) *Parser {
	p := new(Parser)
	*p = Parser{
		Input:     in,
		Table:     symbol.InitSymbolTable(in.GetPath()),
		//ClassTable: ClassTable.InitClassTable(),
		Stack:     NewAstStack(),
		ParsingClass: false,
		markIndex: 0,
	}
	return p
}

func (p *Parser) Advance() {
	p.Current = p.Next
	p.Next = p.Input.Next()
}
