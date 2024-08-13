// =================================================================================================
// Alex Peters - 2024
//
// creates language/compiler builtins
// =================================================================================================

package parser

import (
	"os"

	"github.com/petersalex27/yew/token"
	"github.com/petersalex27/yew/types"
)

func builtin_Type_n() (Type, Type1 types.Sort) {
	return types.Type0, types.Type1
}

// creates ...
// 	(=) : a -> b -> Type
// 	(=) := \x, y => x `=` y
func (parser *Parser) builtin_Eq_a_b() {
	var eq types.Pi

	Type, _ := builtin_Type_n()
	a, b := Var("a"), Var("b")
	types.SetKind(&a, Type)
	types.SetKind(&b, Type)
	
	Eq := types.Constant{C: "="}
	eqIntro, _ := parser.env.Prod(types.AsTyping(a))
	eqInnerIntro, _ := parser.env.Prod(types.AsTyping(b))
	// {b : Type} -> b -> Type
	eqInner, _ := eqInnerIntro(Type)
	// {a : Type} -> a -> {b : Type} -> b -> Type
	eq, _ = eqIntro(eqInner)
	// (=) : a -> b -> Type
	// (=) := \x, y => x `=` y
	intro, _ := parser.env.TypeCon(Type, Eq, eq)
	eqFunc_, _, _ := parser.env.Get(Eq)
	eqFunc := eqFunc_.(types.Lambda)

	eqToken := token.Infix.MakeValued("(=)")
	setType, _ := parser.declare(eqToken, false, exports{})
	bp := uint8(2)
	rAssoc := uint8(0)
	_ = setType(parser, eq, true, bp, rAssoc)

	// Refl : {x : a} -> x = x
	Refl := types.Constant{C: "Refl"}
	x := Var("x")
	x.Kind = a

	eqFuncOnce, eqOnceApp, _ := parser.env.App(eqFunc, eq, x)
	ReflType_App, _, _ := parser.env.Apply(eqFuncOnce, eqOnceApp, x)
	xProd := Var("x")
	xProd.Kind = a
	introRefl, _ := parser.env.ImplicitProd(xProd)
	ReflPi, _ := introRefl(ReflType_App.(types.Type))
	_ = intro([]types.Constant{Refl}, []types.Type{ReflPi})

	rf, rt, _ := parser.env.Get(Refl)
	reflToken := token.Id.MakeValued("Refl")
	setType, _ = parser.declare(reflToken, false, exports{})
	_ = setType(parser, rt, false)

	// report debug info
	debug_log_builtin(os.Stderr, Eq, eqFunc, eq)
	debug_log_builtin(os.Stderr, Refl, rf, rt)
	return
}

// creates ...
// 	Int : Type
// 	Int = Int
// 	Uint : Type
// 	Uint = Uint
// 	Float : Type
// 	Float = Float
// 	Char : Type
// 	Char = Char
// 	String : Type
// 	String = String
func (parser *Parser) builtin_prims() {
	Type, Type_1 := builtin_Type_n()
	setType, _ := parser.declare(token.Id.MakeValued("Type"), false, exports{})
	_ = setType(parser, Type_1, false)

	var Int, Uint, Float, Char, String types.Variable
	{
		int := token.Id.MakeValued("Int")
		uint := token.Id.MakeValued("Uint")
		float := token.Id.MakeValued("Float")
		char := token.Id.MakeValued("Char")
		string := token.Id.MakeValued("String")
		Int = types.MakeVar(int, 0, types.Unrestricted, Type)
		Uint = types.MakeVar(uint, 0, types.Unrestricted, Type)
		Float = types.MakeVar(float, 0, types.Unrestricted, Type)
		Char = types.MakeVar(char, 0, types.Unrestricted, Type)
		String = types.MakeVar(string, 0, types.Unrestricted, Type)
	}
	// Int : Type
	parser.env.Declare(Int, Type)
	setType, _ = parser.declare(token.Id.MakeValued("Int"), false, exports{})
	_ = setType(parser, Type, false)

	// Uint : Type
	parser.env.Declare(Uint, Type)
	setType, _ = parser.declare(token.Id.MakeValued("Uint"), false, exports{})
	_ = setType(parser, Type, false)

	// Float : Type
	parser.env.Declare(Float, Type)
	setType, _ = parser.declare(token.Id.MakeValued("Float"), false, exports{})
	_ = setType(parser, Type, false)

	// Char : Type
	parser.env.Declare(Char, Type)
	setType, _ = parser.declare(token.Id.MakeValued("Char"), false, exports{})
	_ = setType(parser, Type, false)

	// String : Type
	parser.env.Declare(String, Type)
	setType, _ = parser.declare(token.Id.MakeValued("String"), false, exports{})
	_ = setType(parser, Type, false)

	// Int = Int
	parser.env.Assign(Int, Int)
	debug_log_builtin(os.Stderr, Int, Int, Type)
	// Uint = Uint
	parser.env.Assign(Uint, Uint)
	debug_log_builtin(os.Stderr, Uint, Uint, Type)
	// Float = Float
	parser.env.Assign(Float, Float)
	debug_log_builtin(os.Stderr, Float, Float, Type)
	// Char = Char
	parser.env.Assign(Char, Char)
	debug_log_builtin(os.Stderr, Char, Char, Type)
	// String = String
	parser.env.Assign(String, String)
	debug_log_builtin(os.Stderr, String, String, Type)
}

func (parser *Parser) declareBuiltin(fromPath string) {
	// TODO: actually use path
	parser.builtin_prims()
	parser.builtin_Eq_a_b()
	parser.printDeclarations(os.Stderr)
}