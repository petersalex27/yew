package ast

import (
	"fmt"
	scan "yew/lex"
	. "yew/parser/node-type"
	. "yew/parser/parser"
	symbol "yew/symbol"
	types "yew/type"
)

type PostfixOperation struct {
	op      PostOpType
	operand Expression
}

func MakePostfixOperation(op PostOpType, operand Expression) PostfixOperation {
	return PostfixOperation{op, operand}
}

func (post PostfixOperation) ResolveNames(table *symbol.SymbolTable) bool {
	return post.operand.ResolveNames(table)
}
func (post PostfixOperation) GetNodeType() NodeType {
	return POPERATION
}

var postOperationRule = NodeRule{Production: POPERATION, Expression: []NodeType{EXPRESSION, POP_}}

func (post PostfixOperation) Make(p *Parser) bool {
	valid, e := p.Stack.Validate(postOperationRule)
	if !valid {
		e.Print()
		return false
	}
	post.op = p.Stack.Pop().(PostOpType)
	post.operand = p.Stack.Pop().(Expression)
	p.Stack.Push(post)
	return true
}
func (post PostfixOperation) Equal_test(a Ast) bool {
	equal := a.GetNodeType() == POPERATION
	if !equal {
		return false
	}
	post2 := a.(PostfixOperation)
	return equal &&
		post.op.Equal_test(post2.op) &&
		post.operand.Equal_test(post2.operand)
}
func (post PostfixOperation) Print(ls []string) {
	lines := make([]string, len(ls))
	lines = append(lines, ls...)
	lines = printLines(lines)
	fmt.Printf("PostfixOperation\n")
	lines = append(lines, " ├─")
	post.op.Print(lines)
	lines[len(lines)-1] = " └─"
	post.operand.Print(lines)
}
func (post PostfixOperation) ExpressionType() types.Types {
	return post.op.GetFunctionType(nil)
}
func (post PostfixOperation) DoTypeInference(newTypeInformation types.Types) types.Types {
	ty := post.op.GetFunctionType(nil)
	return ty.InferType(newTypeInformation)
}
func (post PostfixOperation) FindStartToken() scan.Token {
	return post.operand.FindStartToken()
}
