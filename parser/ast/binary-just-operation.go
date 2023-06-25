package ast

import (
	"fmt"
	err "yew/error"
	"yew/lex"
	nodetype "yew/parser/node-type"
	"yew/parser/parser"
	"yew/symbol"
	types "yew/type"
)

type OpType scan.TokenType

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
			expression: id1,
			expressionType: taus[0],
		},
	}
	param2 := Parameter{
		paramIndex: 0,
		pattern: ExpressionTypeAnnotation{
			expression: id2,
			expressionType: taus[1],
		},
	}

	lam := Lambda{
		binder: param1,
		bound: Lambda{
			binder: param2,
			bound: expr,
		},
	}
	return Function{MakeIdWithType(idToken, ty), lam}
}

// operations
const (
	ADD = OpType(scan.PLUS)

	APPEND   = OpType(scan.PLUS_PLUS)
	SUBTRACT = OpType(scan.MINUS)
	MULTIPLY = OpType(scan.STAR)
	DIVIDE   = OpType(scan.SLASH)
	POWER    = OpType(scan.HAT)

	CONSTRUCT = OpType(scan.COLON)

	NOT_EQUALS = OpType(scan.BANG_EQUALS)
	EQUALS     = OpType(scan.EQUALS_EQUALS)
	AND        = OpType(scan.AMPER_AMPER)
	OR         = OpType(scan.BAR_BAR)

	DOT = OpType(scan.DOT)

	GREAT        = OpType(scan.GREAT)
	LESS         = OpType(scan.LESS)
	GREAT_EQUALS = OpType(scan.GREAT_EQUALS)
	LESS_EQUALS  = OpType(scan.LESS_EQUALS)

	MAPS_TO = OpType(scan.ARROW)

	MOD = OpType(scan.MOD)
)

func (o OpType) StackLogString() string {
	return fmt.Sprintf("%s; %s", o.GetNodeType().ToString(), o.ToString())
}

func (o OpType) Print(lines []string) {
	printLines(lines)
	fmt.Printf("OpType == %s\n", o.ToString())
}

func (ot OpType) ToString() string {
	switch ot {
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
	o2 := a.(OpType)
	return equal && o2 == o
}

func (b OpType) GetFunctionType(*symbol.SymbolTable) types.Types {
	switch b {
	case ADD:
		return arith(types.Tau("+"))
	case SUBTRACT:
		return arith(types.Tau("-"))
	case MULTIPLY:
		return arith(types.Tau("*"))
	case DIVIDE:
		return arith(types.Tau("/"))
	case POWER:
		return arith(types.Tau("^"))
	case EQUALS:
		return relate(types.Tau("=="))
	case NOT_EQUALS:
		return relate(types.Tau("!="))
	case GREAT:
		return order(types.Tau(">"))
	case LESS:
		return order(types.Tau("<"))
	case GREAT_EQUALS:
		return order(types.Tau(">="))
	case LESS_EQUALS:
		return order(types.Tau("<="))
	case AND:
	case OR:
		return types.Function{
			Domain: types.Bool{},
			Codomain: types.Function{
				Domain:   types.Bool{},
				Codomain: types.Bool{},
			},
		}
	case APPEND:
		return list(types.Tau("++"))
	case CONSTRUCT:
		return list(types.Tau(":"))
	}

	err.PrintBug()
	panic("")
}