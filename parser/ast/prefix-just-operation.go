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

type UOpType scan.OtherToken

func (o UOpType) AsFunction(p *parser.Parser) Function {
	idToken := scan.MakeIdToken(o.ToString(), 0, 0)
	ty := o.GetFunctionType(nil)
	tau := types.GetNewTau()

	id1 := MakeIdWithType(scan.MakeIdToken("x", 0, 0), tau)

	p.Stack.Push(o)
	p.Stack.Push(id1)
	if !(UnaryOperation{}.Make(p)) {
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

// unary operations
const (
	NOT      = scan.BANG
	POSITIVE = scan.PLUS_PREFIX__
	NEGATIVE = scan.MINUS_PREFIX__
)

func (uot UOpType) ToString() string {
	switch uot.FindStartToken().GetType() {
	case NOT:
		return "!"
	case POSITIVE:
		return "+"
	case NEGATIVE:
		return "-"
	default:
		err.PrintBug()
		panic("")
	}
}
func (u UOpType) Print(lines []string) {
	printLines(lines)
	fmt.Printf("UOpType == %s\n", u.ToString())
}

func (u UOpType) StackLogString() string {
	return fmt.Sprintf("%s; %s", u.GetNodeType().ToString(), u.ToString())
}

func (u UOpType) ResolveNames(*symbol.SymbolTable) bool { return true }
func (u UOpType) GetNodeType() nodetype.NodeType { 
	return nodetype.UOP_ 
}
func (u UOpType) Make(p *parser.Parser) bool {
	err.PrintBug()
	panic("")
}
func (u UOpType) Equal_test(a parser.Ast) bool {
	equal := a.GetNodeType() == nodetype.UOP_
	if !equal {
		return false
	}
	u2, ok := a.(UOpType)
	return equal && ok && scan.OtherToken(u2).Equal_test_weak(scan.OtherToken(u))
}

var aToA = types.Function{
	Domain: types.Tau("a"),
	Codomain: types.Tau("a"),
}
var aToAToA = types.Function{
	Domain: types.Tau("a"),
	Codomain: aToA,
}

func (u UOpType) GetFunctionType(*symbol.SymbolTable) types.Types {
	switch u.FindStartToken().GetType() {
	case POSITIVE:
		return arith(aToA)
	case NEGATIVE:
		return arith(aToA)
	case NOT:
		return types.Function{Domain: types.Bool{}, Codomain: types.Bool{}}
	}

	err.PrintBug()
	panic("")
}

func (u UOpType) FindStartToken() scan.Token {
	return scan.OtherToken(u)
}