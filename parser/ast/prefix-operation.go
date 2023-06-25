package ast

import (
	"fmt"
	"yew/ir"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	symbol "yew/symbol"
	types "yew/type"
)

type UnaryOperation struct {
	op      UOpType
	operand Expression
}

func MakeUnaryOperation(op UOpType, operand Expression) UnaryOperation {
	return UnaryOperation{op, operand}
}

func (u UnaryOperation) ResolveNames(table *symbol.SymbolTable) bool {
	return u.operand.ResolveNames(table)
}

func (u UnaryOperation) ExpressionType() types.Types {
	opd := u.operand.ExpressionType()
	fn := u.op.GetFunctionType(nil)
	return fn.
		InferType(opd). // remove qualifier (if applicable)
		Apply(opd)      // apply operand's type
}
func (u UnaryOperation) DoTypeInference(newTypeInformation types.Types) types.Types {
	ty := u.op.GetFunctionType(nil)
	return ty.InferType(newTypeInformation)
}
func (u UnaryOperation) Compile(builder *ir.IrBuilder) {

}
func (u UnaryOperation) GetNodeType() NodeType { return UOPERATION }

var unaryOperationRule = NodeRule{
	Production: UOPERATION /* ::= */, Expression: []NodeType{UOP_, EXPRESSION},
}

func (u UnaryOperation) Make(p *Parser) bool {
	valid, e := p.Stack.Validate(unaryOperationRule)
	if !valid {
		e.Print()
		return false
	}

	u.operand = p.Stack.Pop().(Expression)
	u.op = p.Stack.Pop().(UOpType)
	p.Stack.Push(u)
	return true
}
func (u UnaryOperation) Equal_test(a Ast) bool {
	equal := a.GetNodeType() == UOPERATION
	u2 := a.(UnaryOperation)
	return equal &&
		u.op == u2.op &&
		u.operand.Equal_test(u2.operand)
}
func (u UnaryOperation) Print(lines []string) {
	lines = printLines(lines)
	fmt.Printf("UnaryOperation\n")
	lines = append(lines, " ├─")
	u.op.Print(lines)
	lines[len(lines)-1] = " └─"
	u.operand.Print(lines)
}
func (u UnaryOperation) FindStartToken() scan.Token {
	return u.operand.FindStartToken()
}
