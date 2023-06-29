package ast

import (
	"fmt"
	types "yew/type"
	err "yew/error"
	"yew/lex"
	nodetype "yew/parser/node-type"
	"yew/symbol"
	"yew/parser/parser"
)

type PostOpType scan.OtherToken

func (o PostOpType) AsFunction(p *parser.Parser) Function {
	idToken := scan.MakeIdToken(o.ToString(), 0, 0)
	ty := o.GetFunctionType(nil)
	tau := types.GetNewTau()

	id1 := MakeIdWithType(scan.MakeIdToken("x", 0, 0), tau)

	p.Stack.Push(id1)
	p.Stack.Push(o)
	if !(PostfixOperation{}.Make(p)) {
		err.PrintBug()
		panic("")
	}

	expr := p.Stack.Pop().(Expression)

	param1 := Parameter{
		paramIndex: 0,
		pattern: ExpressionTypeAnnotation{
			expression: id1,
			expressionType: tau,
		},
	}

	lam := Lambda{
		binder: param1,
		bound: expr,
	}
	return Function{MakeIdWithType(idToken, ty), lam}
}

// postfix operations
const (
	FACTORIAL = scan.BANG_POSTFIX__
)

func (o PostOpType) Print(lines []string) {
	printLines(lines)
	fmt.Printf("PostOpType == %s\n", o.ToString())
}

func (o PostOpType) StackLogString() string {
	return fmt.Sprintf("%s; %s", o.GetNodeType().ToString(), o.ToString())
}

func (pot PostOpType) ToString() string {
	if scan.OtherToken(pot).GetType() != FACTORIAL {
		err.PrintBug()
		panic("")
	}
	return "_postfix_(!)"
}

func (o PostOpType) ResolveNames(*symbol.SymbolTable) bool { return true }
func (o PostOpType) GetNodeType() nodetype.NodeType { return nodetype.POP_ }
func (o PostOpType) Make(p *parser.Parser) bool {
	err.PrintBug()
	panic("")
}
func (o PostOpType) Equal_test(a parser.Ast) bool {
	equal := a.GetNodeType() == nodetype.POP_
	if !equal {
		return false
	}
	o2, ok := a.(PostOpType)
	return equal && ok && scan.OtherToken(o2).Equal_test_weak(scan.OtherToken(o))
}

func (p PostOpType) GetFunctionType(*symbol.SymbolTable) types.Types {
	switch p.FindStartToken().GetType() {
	case FACTORIAL:
		return factorial
	}

	err.PrintBug()
	panic("")
}
func (o PostOpType) FindStartToken() scan.Token {
	return scan.OtherToken(o)
}
