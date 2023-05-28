package ast

import (
	"fmt"
	builtin "yew/builtin"
	//"yew/demangler"
	err "yew/error"
	"yew/ir"
	scan "yew/lex"
	symbol "yew/symbol"
	types "yew/type"
)

type BinaryOperation struct {
	op OpType
	left Expression
	right Expression
}

type UnaryOperation struct {
	op UOpType
	operand Expression
}

type OpType scan.TokenType
func (o OpType) GetNodeType() NodeType { return BOP_ }
func (o OpType) Make(stack *AstStack) bool {
	// `(op)` => (\x, y -> x `op` y)
	valid := stack.Validate([]NodeType{BOP_})
	if !valid {
		return false
	}
	/*res1, res2 := 
		demangler.CURRY.GetDemanglerPrefix(),
		demangler.CURRY.GetDemanglerPrefix()
	tok1 := scan.MakeIdToken(res1, 0, 0)
	tok2 := scan.MakeIdToken(res2, 0, 0)
	
	lam := Lambda{

	}*/
	err.PrintBug()
	panic("")
}
func (o OpType) equal_test(a Ast) bool {
	equal := a.GetNodeType() == BOP_
	if !equal {
		return false
	}
	o2 := a.(OpType)
	return equal && o2 == o 
} 

type UOpType scan.TokenType
func (u UOpType) GetNodeType() NodeType { return UOP_ }
func (u UOpType) Make(stack *AstStack) bool {
	err.PrintBug()
	panic("")
}
func (u UOpType) equal_test(a Ast) bool {
	equal := a.GetNodeType() == UOP_
	if !equal {
		return false
	}
	u2 := a.(UOpType)
	return equal && u2 == u 
} 

func (b BinaryOperation) ExpressionType() types.Types {
	left, right := b.left.ExpressionType(), b.right.ExpressionType()
	fn := b.op.GetFunctionType(nil)
	return fn.
			InferType(left).	// remove qualifier (if applicable)
			Apply(left).		// apply left type
			Apply(right)		// apply right type
}
func (BinaryOperation) ResolveNames(*symbol.SymbolTable) {
	return // nothing to resolve
}
func (b BinaryOperation) DoTypeInference(newTypeInformation types.Types) types.Types {
	ty := b.op.GetFunctionType(nil)
	return ty.InferType(newTypeInformation)
}
func (b BinaryOperation) Compile(builder *ir.IrBuilder) {
	
}
func (b BinaryOperation) GetNodeType() NodeType { return OPERATION }
func (b BinaryOperation) Make(stack *AstStack) bool {
	valid := stack.Validate([]NodeType{BOP_, EXPRESSION, EXPRESSION})
	if !valid {
		return false
	}
	b.right = stack.Pop().(Expression)
	b.left = stack.Pop().(Expression)
	b.op = stack.Pop().(OpType)
	stack.Push(b)
	return true
}
func (b BinaryOperation) equal_test(a Ast) bool {
	equal := a.GetNodeType() == OPERATION
	b2 := a.(BinaryOperation)
	return equal &&
		b.op == b2.op && 
		b.left.equal_test(b2.left) &&
		b.right.equal_test(b2.right)
}
func (b BinaryOperation) print(n int) {
	printSpaces(n)
	fmt.Printf("BinaryOperation\n")
	b.left.print(n + 1)
	b.op.print(n + 1)
	b.right.print(n + 1)
}

func (u UnaryOperation) ExpressionType() types.Types {
	opd := u.operand.ExpressionType()
	fn := u.op.GetFunctionType(nil)
	return fn.
			InferType(opd).	// remove qualifier (if applicable)
			Apply(opd)		// apply operand's type
}
func (UnaryOperation) ResolveNames(*symbol.SymbolTable) {
	return // nothing to resolve
}
func (u UnaryOperation) DoTypeInference(newTypeInformation types.Types) types.Types {
	ty := u.op.GetFunctionType(nil)
	return ty.InferType(newTypeInformation)
}
func (u UnaryOperation) Compile(builder *ir.IrBuilder) {
	
}
func (u UnaryOperation) GetNodeType() NodeType { return UOPERATION }
func (u UnaryOperation) Make(stack *AstStack) bool {
	valid := stack.Validate([]NodeType{UOP_, EXPRESSION})
	if !valid {
		return false
	}

	u.operand = stack.Pop().(Expression)
	u.op = stack.Pop().(UOpType)
	stack.Push(u)
	return true
}
func (u UnaryOperation) equal_test(a Ast) bool {
	equal := a.GetNodeType() == UOPERATION
	u2 := a.(UnaryOperation)
	return equal &&
		u.op == u2.op && 
		u.operand.equal_test(u2.operand)
}
func (u UnaryOperation) print(n int) {
	printSpaces(n)
	fmt.Printf("UnaryOperation\n")
	u.op.print(n + 1)
	u.operand.print(n + 1)
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
)

// unary operations
const (
	NOT = UOpType(scan.BANG)
	POSITIVE = UOpType(scan.END__ + scan.PLUS)
	NEGATIVE = UOpType(scan.END__ + scan.MINUS)
)

func (u UOpType) print(n int) {
	printSpaces(n)
	fmt.Printf("UOpType == %v\n", u)
}

func (o OpType) print(n int) {
	printSpaces(n)
	fmt.Printf("OpType == %v\n", o)
}


func buildQualified(class types.Class) func(types.Tau) func(types.Types) types.Qualifier {
	return func(typeVariable types.Tau) func(types.Types) types.Qualifier {
		return func(t types.Types) types.Qualifier {
			return types.Qualifier{
				Class: class, 
				TypeVariable: typeVariable, 
				Qualified: t,
			}
		}
	}
}

var arith = buildQualified(builtin.Number)("a")
var relate = buildQualified(builtin.Equalable)("a")
var order = buildQualified(builtin.Orderable)("a")
var list = buildQualified(builtin.Listable)("a")

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
				Domain: types.Bool{}, 
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

func (u UOpType) GetFunctionType(*symbol.SymbolTable) types.Types {
	switch u {
	case POSITIVE:
		return arith(types.Tau("positive"))
	case NEGATIVE:
		return arith(types.Tau("negative"))
	case NOT:
		return arith(types.Function{Domain: types.Bool{}, Codomain: types.Bool{}})
	}

	err.PrintBug()
	panic("")
}
func (u UOpType) Compile(builder *ir.IrBuilder) {
	
}