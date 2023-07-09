package parsing

import (
	"fmt"
	"os"
	"testing"
	"yew/info"
	scan "yew/lex"
	"yew/parser/ast"
	"yew/parser/parser"
	types "yew/type"

	//"yew/symbol"
	//"yew/type"
	"yew/value"
	//"yew/info"
)

var defExpected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeDeclaration(ast.MakeId(scan.MakeIdToken("x", 1, 4))),
		ast.MakeDefinition(ast.MakeAssignment(
			ast.MakeId(scan.MakeIdToken("x", 1, 4)),
			ast.MakeValue(value.Int(1)),
		)),
	},
)

var def2Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeDeclaration(ast.MakeId(scan.MakeIdToken("x", 1, 4))),
		ast.MakeDefinition(ast.MakeAssignment(
			ast.MakeId(scan.MakeIdToken("x", 1, 4)),
			ast.MakeValue(value.Int(1)),
		)),
	},
)

var appExpected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeApplication(
			ast.MakeId(scan.MakeIdToken("myFunction", 1, 0)),
			ast.MakeValue(value.Int(1)),
		),
	},
)

var app2Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeApplication(
			ast.MakeApplication(
				ast.MakeId(scan.MakeIdToken("myFunction", 1, 0)),
				ast.MakeValue(value.Int(1)),
			),
			ast.MakeValue(value.Int(1)),
		),
	},
)

var fnDefExpected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeFunction(
			ast.MakeId(scan.MakeIdToken("id", 1, 0)),
			ast.MakeLambda(
				ast.MakeParameter(0,
					ast.MakeTypeAnnotation(
						ast.MakeId(scan.MakeIdToken("x", 1, 0)),
						types.Int{},
					),
				),
				ast.MakeTypeAnnotation(
					ast.MakeId(scan.MakeIdToken("x", 1, 0)),
					types.Int{},
				),
			),
		),
	},
)

var fnDef2Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeFunction(
			ast.MakeId(scan.MakeIdToken("fn", 1, 0)),
			ast.MakeLambda(
				ast.MakeParameter(1,
					ast.MakeTypeAnnotation(
						ast.MakeId(scan.MakeIdToken("x", 1, 0)),
						types.Int{},
					),
				),
				ast.MakeLambda(
					ast.MakeParameter(0,
						ast.MakeTypeAnnotation(
							ast.MakeId(scan.MakeIdToken("y", 1, 0)),
							types.Int{},
						),
					),
					ast.MakeTypeAnnotation(
						ast.MakeId(scan.MakeIdToken("x", 1, 0)),
						types.Int{},
					),
				),
			),
		),
	},
)

var fnDef3Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeFunction(
			ast.MakeId(scan.MakeIdToken("fn", 1, 0)),
			ast.MakeLambda(
				ast.MakeParameter(1,
					ast.MakeTypeAnnotation(
						ast.MakeId(scan.MakeIdToken("x", 1, 0)),
						types.Int{},
					),
				),
				ast.MakeLambda(
					ast.MakeParameter(0,
						ast.MakeTypeAnnotation(
							ast.MakeId(scan.MakeIdToken("y", 1, 0)),
							types.Char{},
						),
					),
					ast.MakeTypeAnnotation(
						ast.MakeId(scan.MakeIdToken("x", 1, 0)),
						types.Int{},
					),
				),
			),
		),
	},
)

var fnDef4Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeFunction(
			ast.MakeId(scan.MakeIdToken("id", 1, 0)),
			ast.MakeLambda(
				ast.MakeParameter(0,
					ast.MakeTypeAnnotation(
						ast.MakeId(scan.MakeIdToken("x", 1, 0)),
						types.Int{},
					),
				),
				ast.MakeTypeAnnotation(
					ast.Program{
						ast.MakeId(scan.MakeIdToken("x", 1, 0)),
					},
					types.Int{},
				),
			),
		),
	},
)

var fnDef5Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeFunction(
			ast.MakeId(scan.MakeIdToken("myFunction", 1, 0)),
			ast.MakeLambda(
				ast.MakeParameter(0,
					ast.MakeTypeAnnotation(
						ast.MakeId(scan.MakeIdToken("x", 1, 0)),
						types.Int{},
					),
				),
				ast.MakeTypeAnnotation(
					ast.MakeApplication(
						ast.MakeId(scan.MakeIdToken("Just", 1, 0)),
						ast.MakeId(scan.MakeIdToken("x", 1, 0)),
					),
					types.Application{types.Var("Maybe"), types.Int{}},
				),
			),
		),
	},
)

var fnDef6Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeFunction(
			ast.MakeId(scan.MakeIdToken("myFunction", 1, 0)),
			ast.MakeLambda(
				ast.MakeParameter(1,
					ast.MakeTypeAnnotation(
						ast.MakeId(scan.MakeIdToken("x", 1, 0)),
						types.Int{},
					),
				),
				ast.MakeLambda(
					ast.MakeParameter(0,
						ast.MakeTypeAnnotation(
							ast.MakeId(scan.MakeIdToken("y", 1, 0)),
							types.Int{},
						),
					),
					ast.MakeTypeAnnotation(
						ast.Program{
							ast.MakeFunction(
								ast.MakeId(scan.MakeIdToken("myFunction'", 2, 0)),
								ast.MakeLambda(
									ast.MakeParameter(0,
										ast.MakeTypeAnnotation(
											ast.MakeId(scan.MakeIdToken("z", 2, 0)),
											types.Int{},
										),
									),
									ast.MakeTypeAnnotation(
										ast.MakeBinaryOperation(
											ast.OpType(scan.MakeOtherToken(scan.PLUS, 0, 0)),
											ast.MakeId(scan.MakeIdToken("x", 0, 0)),
											ast.MakeBinaryOperation(
												ast.OpType(scan.MakeOtherToken(scan.STAR, 0, 0)),
												ast.MakeId(scan.MakeIdToken("y", 0, 0)),
												ast.MakeId(scan.MakeIdToken("z", 0, 0)),
											),
										),
										types.Int{},
									),
								),
							),
							ast.MakeBinaryOperation(
								ast.OpType(scan.MakeOtherToken(scan.PLUS, 0, 0)),
								ast.MakeBinaryOperation(
									ast.OpType(scan.MakeOtherToken(scan.MINUS, 0, 0)),
									ast.MakeId(scan.MakeIdToken("x", 0, 0)),
									ast.MakeId(scan.MakeIdToken("y", 0, 0)),
								),
								ast.MakeApplication(
									ast.MakeId(scan.MakeIdToken("myFunction'", 2, 0)),
									ast.MakeId(scan.MakeIdToken("x", 1, 0)),
								),
							),
						},
						types.Int{},
					),
				),
			),
		),
	},
)

var fnDef7Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeFunction(
			ast.MakeId(scan.MakeIdToken("myFunction", 1, 0)),
			ast.MakeLambda(
				ast.MakeParameter(0,
					ast.MakeTypeAnnotation(
						ast.MakeId(scan.MakeIdToken("x", 1, 0)),
						types.Int{},
					),
				),
				ast.MakeTypeAnnotation(
					ast.Program{
						ast.MakeFunction(
							ast.MakeId(scan.MakeIdToken("myFunction'", 2, 0)),
							ast.MakeLambda(
								ast.MakeParameter(0,
									ast.MakeTypeAnnotation(
										ast.MakeId(scan.MakeIdToken("y", 2, 0)),
										types.Int{},
									),
								),
								ast.MakeTypeAnnotation(
									ast.MakeId(scan.MakeIdToken("x", 0, 0)),
									types.Int{},
								),
							),
						),
						ast.MakeId(scan.MakeIdToken("myFunction'", 2, 0)),
					},
					types.Function{Domain: types.Int{}, Codomain: types.Int{}},
				),
			),
		),
	},
)

var constraintExpected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeFunction(
			ast.MakeId(scan.MakeIdToken("whatever", 1, 1)),
			ast.MakeLambda(
				ast.MakeParameter(0,
					ast.MakeTypeAnnotation(
						ast.MakeId(scan.MakeIdToken("x", 1, 0)),
						types.Constraint{
							Context: types.ConstraintContext{
								types.Context{
									ClassName: types.Var("Num"),
									TypeVariable: types.Var("a"),
								},
							},
							Constrained: types.Var("a"),
						},
					),
				),
				ast.MakeTypeAnnotation(
					ast.MakeId(scan.MakeIdToken("x", 0, 0)),
					types.Constraint{
						Context: types.ConstraintContext{
							types.Context{
								ClassName: types.Var("Num"),
								TypeVariable: types.Var("a"),
							},
						},
						Constrained: types.Var("a"),
					},
				),
			),
		),
	},
)

var constraint2Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeFunction(
			ast.MakeId(scan.MakeIdToken("whatever", 1, 1)),
			ast.MakeLambda(
				ast.MakeParameter(0,
					ast.MakeTypeAnnotation(
						ast.MakeId(scan.MakeIdToken("x", 1, 0)),
						types.Constraint{
							Context: types.ConstraintContext{
								types.Context{
									ClassName: types.Var("Num"),
									TypeVariable: types.Var("a"),
								},
								types.Context{
									ClassName: types.Var("Num"),
									TypeVariable: types.Var("b"),
								},
							},
							Constrained: types.Var("a"),
						},
					),
				),
				ast.MakeTypeAnnotation(
					ast.MakeId(scan.MakeIdToken("x", 0, 0)),
					types.Constraint{
						Context: types.ConstraintContext{
							types.Context{
								ClassName: types.Var("Num"),
								TypeVariable: types.Var("a"),
							},
							types.Context{
								ClassName: types.Var("Num"),
								TypeVariable: types.Var("b"),
							},
						},
						Constrained: types.Var("b"),
					},
				),
			),
		),
	},
)

var opExpected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeBinaryOperation(
			ast.OpType(scan.MakeOtherToken(scan.PLUS, 0, 0)),
			ast.Value(scan.ValueToken{Value: value.Int(1)}),
			ast.Value(scan.ValueToken{Value: value.Int(1)}),
		),
	},
)

var factorialExpected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakePostfixOperation(
			ast.PostOpType(scan.MakeOtherToken(scan.BANG_POSTFIX__, 0, 0)),
			ast.MakeValue(value.Int(1)),
		),
	},
)
var composeExpected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeApplication(
			ast.MakeApplication(
				ast.MakeId(scan.MakeIdToken("f", 1, 0)),
				ast.MakeId(scan.MakeIdToken("g", 1, 0)),
			),
			ast.MakeId(scan.MakeIdToken("h", 1, 0)),
		),
	},
)
var compose2Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeApplication(
			ast.MakeId(scan.MakeIdToken("f", 1, 0)),
			ast.MakeApplication(
				ast.MakeId(scan.MakeIdToken("g", 1, 0)),
				ast.MakeId(scan.MakeIdToken("h", 1, 0)),
			),
		),
	},
)
var compose3Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeApplication(
			ast.MakeId(scan.MakeIdToken("f", 1, 0)),
			ast.MakeApplication(
				ast.MakeId(scan.MakeIdToken("g", 1, 0)),
				ast.MakeApplication(
					ast.MakeApplication(
						ast.MakeId(scan.MakeIdToken("h", 1, 0)),
						ast.MakeValue(value.Int(1)),
					),
					ast.MakeId(scan.MakeIdToken("i", 1, 0)),
				),
			),
		),
	},
)

var prefixOperationExpected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeUnaryOperation(
			ast.UOpType(scan.MakeOtherToken(scan.PLUS_PREFIX__, 0, 0)),
			ast.MakeValue(value.Int(1)),
		),
	},
)

var packageExpected = ast.MakePackage2(
	scan.MakeIdToken("myPackage", 1, 0),
	ast.Program{},
)

var moduleExpected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeModule(
			scan.MakeIdToken("myModule", 1, 0),
			ast.Program{},
		),
	},
)

var assignAppExpected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeDeclaration(ast.MakeId(scan.MakeIdToken("i", 1, 4))),
		ast.MakeDefinition(ast.MakeAssignment(
			ast.MakeId(scan.MakeIdToken("i", 1, 4)),
			ast.MakeApplication(
				ast.MakeId(scan.MakeIdToken("myFunction", 1, 0)),
				ast.MakeValue(value.Int(1)),
			),
		)),
	},
)

var typeExpected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeTypeDefinition(
			scan.MakeIdToken("Color", 0, 0),
			types.MakeData(
				"Color",
				[]types.Tau{},
				[]types.Constructor{
					types.MakeConstructor("Red", types.Application{}),
					types.MakeConstructor("Blue", types.Application{}),
				},
			),
		),
	},
)

var type2Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeTypeDefinition(
			scan.MakeIdToken("Maybe", 0, 0),
			types.MakeData2(
				"Maybe", []string{"a"},
				[]types.Constructor{
					types.MakeConstructor("Just", types.Application{types.Var("a")}),
					types.MakeConstructor("Nothing", types.Application{}),
				},
			),
		),
	},
)

var type3Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeTypeDefinition(
			scan.MakeIdToken("Either", 0, 0),
			types.MakeData2(
				"Either", []string{"a", "b"},
				[]types.Constructor{
					types.MakeConstructor("Left", types.Application{types.Var("a")}),
					types.MakeConstructor("Right", types.Application{types.Var("b")}),
				},
			),
		),
	},
)

var patternExpected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.Pattern{
			Expression: ast.MakeId(scan.MakeIdToken("a", 0, 0)),
			Matchers: []ast.Lambda{
				ast.MakeLambda( // 1 -> 0
					ast.MakeParameter(0,
						ast.MakeTypeAnnotation(
							ast.MakeValue(value.Int(1)),
							types.Var(".t?"),
						),
					),
					ast.MakeValue(value.Int(0)),
				),
				ast.MakeLambda( // x -> x
					ast.MakeParameter(0,
						ast.MakeTypeAnnotation(
							ast.MakeId(scan.MakeIdToken("x", 0, 0)),
							types.Var(".t?"),
						),
					),
					ast.MakeId(scan.MakeIdToken("x", 0, 0)),
				),
			},
		},
	},
)

var pattern3Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.Pattern{
			Expression: ast.MakeBinaryOperation(
				ast.OpType(scan.MakeOtherToken(scan.PLUS, 0, 0)),
				ast.MakeId(scan.MakeIdToken("a", 0, 0)),
				ast.MakeId(scan.MakeIdToken("b", 0, 0)),
			),
			Matchers: []ast.Lambda{
				ast.MakeLambda( // 0 -> 1
					ast.MakeParameter(0,
						ast.MakeTypeAnnotation(
							ast.MakeValue(value.Int(0)),
							types.Var(".t?"),
						),
					),
					ast.MakeValue(value.Int(1)),
				),
				ast.MakeLambda( // x y -> 0
					ast.MakeParameter(0,
						ast.MakeTypeAnnotation(
							ast.MakeApplication(
								ast.MakeId(scan.MakeIdToken("x", 0, 0)),
								ast.MakeId(scan.MakeIdToken("y", 0, 0)),
							),
							types.Var(".t?"),
						),
					),
					ast.MakeValue(value.Int(0)),
				),
			},
		},
	},
)

var classExpected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeClass(
			ast.MakeId(scan.MakeIdToken("MyClass", 1, 6)),
			map[string]types.Function{
				"f": { Domain: types.Int{}, Codomain: types.Int{}, },
			},
		),
	},
)

var class2Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeClass(
			ast.MakeId(scan.MakeIdToken("MyClass", 1, 6)),
			map[string]types.Function{
				"f": { Domain: types.Int{}, Codomain: types.Int{}, },
				"g": { Domain: types.MakeTau("a", info.Loc{}), Codomain: types.Int{}, },
			},
		),
	},
)

var class3Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeClass(
			ast.MakeId(scan.MakeIdToken("MyClass", 1, 6)),
			map[string]types.Function{
				"f": { Domain: types.Int{}, Codomain: types.Int{}, },
				"g": { Domain: types.MakeTau("a", info.Loc{}), Codomain: types.Int{}, },
				"h": { 
					Domain: types.Constraint{
						Context: types.ConstraintContext{
							types.Context{
								ClassName: types.Var("Num"),
								TypeVariable: types.Var("x"),
							},
						},
						Constrained: types.MakeTau("a", info.Loc{}),
					}, 
					Codomain: types.Constraint{
						Context: types.ConstraintContext{
							types.Context{
								ClassName: types.Var("Num"),
								TypeVariable: types.Var("x"),
							},
						},
						Constrained: types.MakeTau("x", info.Loc{}), 
					},
				},
			},
		),
	},
)

var class4Expected = class2Expected

var class5Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeClass(
			ast.MakeId(scan.MakeIdToken("MyClass", 1, 18)),
			map[string]types.Function{
				"f": { 
					Domain: types.ConstrainType(types.Int{}, types.MakeContext("Ord", "a")),
					Codomain: types.ConstrainType(types.Int{}, types.MakeContext("Ord", "a")),
				},
			},
		),
	},
)

var class6Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeClass(
			ast.MakeId(scan.MakeIdToken("MyClass", 1, 18)),
			map[string]types.Function{
				"f": { 
					Domain: types.ConstrainType(types.Int{}, types.MakeContext("Ord", "a")),
					Codomain: types.ConstrainType(types.Int{}, types.MakeContext("Ord", "a")),
				},
				"g": { 
					Domain: types.ConstrainType(types.Var("a"), types.MakeContext("Ord", "a")), 
					Codomain: types.ConstrainType(types.Int{}, types.MakeContext("Ord", "a")),
				},
			},
		),
	},
)

var class7Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeClass(
			ast.MakeId(scan.MakeIdToken("MyClass", 1, 18)),
			map[string]types.Function{
				"f": { 
					Domain: types.ConstrainType(types.Int{}, types.MakeContext("Ord", "a")),
					Codomain: types.ConstrainType(types.Int{}, types.MakeContext("Ord", "a")),
				},
				"g": { 
					Domain: types.ConstrainType(types.Var("a"), types.MakeContext("Ord", "a")), 
					Codomain: types.ConstrainType(types.Int{}, types.MakeContext("Ord", "a")),
				},
				"h": { // (Num x, Ord a) => a -> x
					Domain: types.ConstrainType(
						types.Var("a"), 
						types.MakeContext("Ord", "a"), types.MakeContext("Num", "x"), 
					), 
					Codomain: types.ConstrainType(
						types.Var("x"), 
						types.MakeContext("Ord", "a"), types.MakeContext("Num", "x"), 
					),
				},
			},
		),
	},
)

var instanceExpected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeInstanceFunction("MyClass", types.Int(info.DefaultLoc()),
			ast.MakeId(scan.MakeIdToken("f", 1, 0)),
			ast.MakeLambda(
				ast.MakeParameter(0,
					ast.MakeTypeAnnotation(
						ast.MakeId(scan.MakeIdToken("a", 1, 0)),
						types.Var(".?"),
					),
				),
				ast.MakeTypeAnnotation(
					ast.MakeId(scan.MakeIdToken("a", 1, 0)),
					types.Var(".?"),
				),
			),
		),
	},
)

var instance2Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.Program{
		ast.MakeInstanceFunction("MyClass", types.Int(info.DefaultLoc()),
			ast.MakeId(scan.MakeIdToken("f", 1, 0)),
			ast.MakeLambda(
				ast.MakeParameter(0,
					ast.MakeTypeAnnotation(
						ast.MakeId(scan.MakeIdToken("a", 1, 0)),
						types.Var(".?"),
					),
				),
				ast.MakeTypeAnnotation(
					ast.MakeId(scan.MakeIdToken("a", 1, 0)),
					types.Var(".?"),
				),
			),
		),
		ast.MakeInstanceFunction("MyClass", types.Int(info.DefaultLoc()),
			ast.MakeId(scan.MakeIdToken("g", 1, 0)),
			ast.MakeLambda(
				ast.MakeParameter(1,
					ast.MakeTypeAnnotation(
						ast.MakeId(scan.MakeIdToken("a", 1, 0)),
						types.Var(".?"),
					),
				),
				ast.MakeLambda(
					ast.MakeParameter(0,
						ast.MakeTypeAnnotation(
							ast.MakeId(scan.MakeIdToken("b", 1, 0)),
							types.Var(".?"),
						),
					),
					ast.MakeTypeAnnotation(
						ast.Sequence{
							ast.MakeId(scan.MakeIdToken("a", 1, 0)),
							ast.MakeApplication(
								ast.MakeId(scan.MakeIdToken("b", 1, 0)),
								ast.MakeApplication(
									ast.MakeId(scan.MakeIdToken("a", 1, 0)),
									ast.MakeValue(value.Int(1)),
								),
							),
						},
						types.Var(".?"),
					),
				),
			),
		),
	},
)

var asts = []struct {
	path string
	ast_ parser.Ast
}{
	//*
	{"./test/def.yw", defExpected},
	{"./test/def2.yw", def2Expected},
	{"./test/app.yw", appExpected},
	{"./test/package.yw", packageExpected},
	{"./test/module.yw", moduleExpected},
	{"./test/app2.yw", app2Expected},
	{"./test/op.yw", opExpected},
	{"./test/factorial.yw", factorialExpected},
	{"./test/prefix.yw", prefixOperationExpected},
	{"./test/fnDef.yw", fnDefExpected},
	{"./test/fnDef2.yw", fnDef2Expected}, 
	{"./test/fnDef3.yw", fnDef3Expected},
	{"./test/fnDef4.yw", fnDef4Expected},
	{"./test/assignApp.yw", assignAppExpected},
	{"./test/type.yw", typeExpected},
	{"./test/type2.yw", type2Expected},
	{"./test/type3.yw", type3Expected},
	{"./test/fnDef5.yw", fnDef5Expected},
	{"./test/compose.yw", composeExpected},
	{"./test/compose2.yw", compose2Expected},
	{"./test/compose3.yw", compose3Expected},
	{"./test/pattern.yw", patternExpected},
	{"./test/pattern2.yw", patternExpected},
	{"./test/pattern3.yw", pattern3Expected},
	{"./test/fnDef6.yw", fnDef6Expected},
	{"./test/fnDef7.yw", fnDef7Expected},//*/
	{"./test/type-constraint.yw", constraintExpected},
	{"./test/type-constraint2.yw", constraint2Expected},
	{"./test/type-constraint3.yw", constraint2Expected},
	{"./test/class.yw", classExpected},
	{"./test/class2.yw", class2Expected},
	{"./test/class3.yw", class3Expected},
	{"./test/class4.yw", class4Expected},
	{"./test/class5.yw", class5Expected},
	{"./test/class6.yw", class6Expected},
	{"./test/class7.yw", class7Expected},
	{"./test/instance.yw", instanceExpected},
	{"./test/instance2.yw", instance2Expected},
}

func TestParse(t *testing.T) {
	for _, test := range asts {
		//fmt.Fprintf(os.Stderr, ">>> running: %s\n", test.path)
		in, e := scan.Init(test.path)
		if nil != e {
			fmt.Fprintf(os.Stderr, ">>> failed: %s <<<\n", test.path)
			fmt.Fprintf(os.Stderr, "%s\n", e.Error())
			t.FailNow()
		}

		ok, prog := Parse(&in)
		if !ok {
			fmt.Fprintf(os.Stderr, ">>> failed: %s <<<\n", test.path)
			t.FailNow()
		}

		/*
		if test.path == "./test/type-constraint3.yw" {
			ast.PrintAst(prog)
		}//*/
		//ast.PrintAst(prog)
		if !ast.EqualTest(prog, test.ast_) {
			fmt.Fprintf(os.Stderr, ">>> failed: %s <<<\n", test.path)
			fmt.Printf("Expected: \n")
			ast.PrintAst(test.ast_)

			fmt.Printf("Actual: \n")
			ast.PrintAst(prog)
			t.FailNow()
		}
	}
}

func TestConstraintBad(t *testing.T) {
	var inSource3 = []string{
		"class (C a) => MyClass a where\n",
		"  fn :: (W a) => a -> Int",
	}
	var in3 = scan.CreateInputStream(
		"test/ast-3", 0, inSource3,
	)
	var expectMsg3 =
		"[test/ast-3:2:12] Type Error: class type parameter cannot be constrained in a class function.\n" +
		"    2 |   fn :: (W a) => a -> Int\n" +
		"                   ^"
	var constraint = types.Constraint{
		Context: types.ConstraintContext{
			types.Context{
				ClassName: types.MakeTau("W", info.MakeLocation(2, 10)),
				TypeVariable: types.MakeTau("a", info.MakeLocation(2, 12)),
			},
		},
	}

	ok, e := validateConstraint(constraint, types.MakeTau("a", info.MakeLocation(1, 24)), in3)
	if ok {
		fmt.Fprintf(os.Stderr, "expected error but validating contraint succeded.\n")
		t.FailNow()
	}

	if expectMsg3 != e.ToString() {
		fmt.Fprintf(os.Stderr, "Expected:\n%s\nActual:\n%s\n", expectMsg3, e.ToString())
		t.FailNow()
	}
	
}