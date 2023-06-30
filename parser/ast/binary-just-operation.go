package ast

import (
	"fmt"
	err "yew/error"
	scan "yew/lex"
	nodetype "yew/parser/node-type"
	"yew/parser/parser"
	"yew/symbol"
	types "yew/type"
)

type OpType scan.OtherToken

func (o OpType) AsFunction(p *parser.Parser) Function {
	idToken := scan.MakeIdToken(o.ToString(), 0, 0)
	ty := o.GetFunctionType(nil)
	taus := types.GetNewTaus(2)

	id1 := MakeIdWithType(scan.MakeIdToken("x", 0, 0), taus[0])
	id2 := MakeIdWithType(scan.MakeIdToken("y", 0, 0), taus[1])

	p.Stack.Push(id1)
	p.Stack.Push(o)
	p.Stack.Push(id2)
	if !(BinaryOperation{}.Make(p)) {
		err.PrintBug()
		panic("")
	}

	expr := p.Stack.Pop().(Expression)

	param1 := Parameter{
		paramIndex: 1,
		pattern: ExpressionTypeAnnotation{
			expression:     id1,
			expressionType: taus[0],
		},
	}
	param2 := Parameter{
		paramIndex: 0,
		pattern: ExpressionTypeAnnotation{
			expression:     id2,
			expressionType: taus[1],
		},
	}

	lam := Lambda{
		binder: param1,
		bound: Lambda{
			binder: param2,
			bound:  expr,
		},
	}
	return Function{MakeIdWithType(idToken, ty), lam}
}

// operations
const (
	ADD = scan.PLUS

	APPEND   = scan.PLUS_PLUS
	SUBTRACT = scan.MINUS
	MULTIPLY = scan.STAR
	DIVIDE   = scan.SLASH
	POWER    = scan.HAT

	CONSTRUCT = scan.COLON

	NOT_EQUALS = scan.BANG_EQUALS
	EQUALS     = scan.EQUALS_EQUALS
	AND        = scan.AMPER_AMPER
	OR         = scan.BAR_BAR

	DOT = scan.DOT

	GREAT        = scan.GREAT
	LESS         = scan.LESS
	GREAT_EQUALS = scan.GREAT_EQUALS
	LESS_EQUALS  = scan.LESS_EQUALS

	MAPS_TO = scan.ARROW

	MOD = scan.MOD
)

func (o OpType) StackLogString() string {
	return fmt.Sprintf("%s; %s", o.GetNodeType().ToString(), o.ToString())
}

func (o OpType) Print(lines []string) {
	printLines(lines)
	fmt.Printf("OpType == %s\n", o.ToString())
}

func (ot OpType) ToString() string {
	switch ot.FindStartToken().GetType() {
	case ADD:
		return "(+)"
	case APPEND:
		return "(++)"
	case SUBTRACT:
		return "(-)"
	case MULTIPLY:
		return "(*)"
	case DIVIDE:
		return "(/)"
	case POWER:
		return "(^)"
	case CONSTRUCT:
		return "(:)"
	case NOT_EQUALS:
		return "(!=)"
	case EQUALS:
		return "(==)"
	case AND:
		return "(&&)"
	case OR:
		return "(||)"
	case DOT:
		return "(.)"
	case GREAT:
		return "(>)"
	case LESS:
		return "(<)"
	case GREAT_EQUALS:
		return "(>=)"
	case LESS_EQUALS:
		return "(<=)"
	case MOD:
		return "(mod)"
	case MAPS_TO:
		return "(->)"
	default:
		err.PrintBug()
		panic("")
	}
}

func (o OpType) ResolveNames(*symbol.SymbolTable) bool { return true }

func (o OpType) GetNodeType() nodetype.NodeType { return nodetype.BOP_ }

func (o OpType) Make(p *parser.Parser) bool {
	err.PrintBug()
	panic("")
}
func (o OpType) Equal_test(a parser.Ast) bool {
	equal := a.GetNodeType() == nodetype.BOP_
	if !equal {
		return false
	}
	o2, ok := a.(OpType)
	if !ok {
		return false
	}

	tok := scan.OtherToken(o)
	tok2 := scan.OtherToken(o2)
	return equal && tok.Equal_test_weak(tok2)
}

func (b OpType) GetFunctionType(*symbol.SymbolTable) types.Types {
	switch b.FindStartToken().GetType() {
	case ADD:
		fallthrough
	case SUBTRACT:
		fallthrough
	case MULTIPLY:
		fallthrough
	case DIVIDE:
		fallthrough
	case POWER:
		return arith(aToAToA)
	case EQUALS:
		fallthrough
	case NOT_EQUALS:
		return relate(aToAToA)
	case GREAT:
		fallthrough
	case LESS:
		fallthrough
	case GREAT_EQUALS:
		fallthrough
	case LESS_EQUALS:
		return order(aToAToA)
	case AND:
		fallthrough
	case OR:
		return types.Function{
			Domain: types.Bool{},
			Codomain: types.Function{
				Domain:   types.Bool{},
				Codomain: types.Bool{},
			},
		}
	case APPEND:
		fallthrough
	case CONSTRUCT:
		return list(aToAToA)
	}

	err.PrintBug()
	panic("")
}

func (o OpType) FindStartToken() scan.Token {
	return scan.OtherToken(o)
}
