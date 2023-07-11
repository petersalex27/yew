package ast

import (
	"fmt"
	"yew/ir"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	types "yew/type"
)

type BinaryOperation struct {
	op    OpType
	left  Expression
	right Expression
}

func MakeBinaryOperation(op OpType, left Expression, right Expression) BinaryOperation {
	return BinaryOperation{op, left, right}
}

func (b BinaryOperation) ResolveNames(p *Parser) bool {
	return b.left.ResolveNames(p) && b.right.ResolveNames(p)
}

func (b BinaryOperation) ExpressionType() types.Types {
	left, right := b.left.ExpressionType(), b.right.ExpressionType()
	fn := b.op.GetFunctionType(nil)
	return fn.
		InferType(left). // remove qualifier (if applicable)
		Apply(left).     // apply left type
		Apply(right)     // apply right type
}
func (b BinaryOperation) DoTypeInference(newTypeInformation types.Types) types.Types {
	ty := b.op.GetFunctionType(nil)
	return ty.InferType(newTypeInformation)
}
func (b BinaryOperation) Compile(builder *ir.IrBuilder) {

}
func (b BinaryOperation) GetNodeType() NodeType { return OPERATION }

func (b BinaryOperation) Make(p *Parser) bool {
	valid, e := p.Stack.Validate(binaryOperationRule)
	if !valid {
		e(p.Input).Print()
		return false
	}
	b.right = p.Stack.Pop().(Expression)
	b.op = p.Stack.Pop().(OpType)
	b.left = p.Stack.Pop().(Expression)
	p.Stack.Push(b)
	return true
}
func (b BinaryOperation) Equal_test(a Ast) bool {
	equal := a.GetNodeType() == OPERATION
	b2, ok := a.(BinaryOperation)
	return equal && ok &&
		b.op.Equal_test(b2.op) &&
		b.left.Equal_test(b2.left) &&
		b.right.Equal_test(b2.right)
}
func (b BinaryOperation) Print(ls []string) {
	lines := make([]string, len(ls))
	lines = append(lines, ls...)
	lines = printLines(lines)
	fmt.Printf("BinaryOperation\n")
	lines = append(lines, " ├─")
	b.left.Print(lines)
	b.op.Print(lines)
	lines[len(lines)-1] = " └─"
	b.right.Print(lines)
}

func (b BinaryOperation) FindStartToken() scan.Token {
	return b.left.FindStartToken()
}
