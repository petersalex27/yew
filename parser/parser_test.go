package parsing

import (
	"fmt"
	"os"
	"testing"
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
	ast.MakeProgram(
		[]ast.Statement{
			ast.MakeDeclaration(ast.MakeId(scan.MakeIdToken("x", 1, 4))),
			ast.MakeDefinition(ast.MakeAssignment(
				ast.MakeId(scan.MakeIdToken("x", 1, 4)),
				ast.MakeValue(value.Int(1)),
			)),
		},
		ast.EmptyExpression{},
	),
)

var def2Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.MakeProgram(
		[]ast.Statement{
			ast.MakeDeclaration(ast.MakeId(scan.MakeIdToken("x", 1, 4))),
			ast.MakeDefinition(ast.MakeAssignment(
				ast.MakeId(scan.MakeIdToken("x", 1, 4)),
				ast.MakeValue(value.Int(1)),
			)),
		},
		ast.EmptyExpression{},
	),
)

var appExpected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.MakeProgram(
		[]ast.Statement{},
		ast.MakeApplication(
			ast.MakeId(scan.MakeIdToken("myFunction", 1, 0)),
			ast.MakeValue(value.Int(1)),
		),
	),
)

var app2Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.MakeProgram(
		[]ast.Statement{},
		ast.MakeApplication(
			ast.MakeApplication(
				ast.MakeId(scan.MakeIdToken("myFunction", 1, 0)),
				ast.MakeValue(value.Int(1)),
			),
			ast.MakeValue(value.Int(1)),
		),
	),
)

var fnDefExpected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.MakeProgram(
		[]ast.Statement{
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
						ast.MakeProgram(
							[]ast.Statement{},
							ast.MakeId(scan.MakeIdToken("x", 1, 0)),
						),
						types.Int{},
					),
				),
			),
		},
		ast.EmptyExpression{},
	),
)

var fnDef2Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.MakeProgram(
		[]ast.Statement{
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
							ast.MakeProgram(
								[]ast.Statement{},
								ast.MakeId(scan.MakeIdToken("x", 1, 0)),
							),
							types.Int{},
						),
					),
				),
			),
		},
		ast.EmptyExpression{},
	),
)

var fnDef3Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.MakeProgram(
		[]ast.Statement{
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
							ast.MakeProgram(
								[]ast.Statement{},
								ast.MakeId(scan.MakeIdToken("x", 1, 0)),
							),
							types.Int{},
						),
					),
				),
			),
		},
		ast.EmptyExpression{},
	),
)

var fnDef4Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.MakeProgram(
		[]ast.Statement{
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
						ast.MakeProgram(
							[]ast.Statement{},
							ast.Sequence{
								ast.MakeId(scan.MakeIdToken("x", 1, 0)),
							},
						),
						types.Int{},
					),
				),
			),
		},
		ast.EmptyExpression{},
	),
)

var fnDef5Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.MakeProgram(
		[]ast.Statement{
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
						ast.MakeProgram(
							[]ast.Statement{},
							ast.MakeApplication(
								ast.MakeId(scan.MakeIdToken("Just", 1, 0)),
								ast.MakeId(scan.MakeIdToken("x", 1, 0)),
							),
						),
						types.Application{types.Tau("Maybe"), types.Int{}},
					),
				),
			),
		},
		ast.EmptyExpression{},
	),
)

var opExpected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.MakeProgram(
		[]ast.Statement{},
		ast.MakeBinaryOperation(
			ast.ADD,
			ast.MakeValue(value.Int(1)),
			ast.MakeValue(value.Int(1)),
		),
	),
)

var factorialExpected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.MakeProgram(
		[]ast.Statement{},
		ast.MakePostfixOperation(
			ast.FACTORIAL,
			ast.MakeValue(value.Int(1)),
		),
	),
)
var composeExpected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.MakeProgram(
		[]ast.Statement{},
		ast.MakeApplication(
			ast.MakeApplication(
				ast.MakeId(scan.MakeIdToken("f", 1, 0)),
				ast.MakeId(scan.MakeIdToken("g", 1, 0)),
			),
			ast.MakeId(scan.MakeIdToken("h", 1, 0)),
		),
	),
)
var compose2Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.MakeProgram(
		[]ast.Statement{},
		ast.MakeApplication(
			ast.MakeId(scan.MakeIdToken("f", 1, 0)),
			ast.MakeApplication(
				ast.MakeId(scan.MakeIdToken("g", 1, 0)),
				ast.MakeId(scan.MakeIdToken("h", 1, 0)),
			),
		),
	),
)
var compose3Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.MakeProgram(
		[]ast.Statement{},
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
	),
)

var prefixOperationExpected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.MakeProgram(
		[]ast.Statement{},
		ast.MakeUnaryOperation(
			ast.POSITIVE,
			ast.MakeValue(value.Int(1)),
		),
	),
)

var packageExpected = ast.MakePackage2(
	scan.MakeIdToken("myPackage", 1, 0),
	ast.MakeProgram( // empty program
		[]ast.Statement{},
		ast.EmptyExpression{},
	),
)

var moduleExpected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.MakeProgram(
		[]ast.Statement{
			ast.MakeModule(
				scan.MakeIdToken("myModule", 1, 0),
				ast.MakeProgram( //empty program
					[]ast.Statement{},
					ast.EmptyExpression{},
				),
			),
		},
		ast.EmptyExpression{},
	),
)

var assignAppExpected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.MakeProgram(
		[]ast.Statement{
			ast.MakeDeclaration(ast.MakeId(scan.MakeIdToken("i", 1, 4))),
			ast.MakeDefinition(ast.MakeAssignment(
				ast.MakeId(scan.MakeIdToken("i", 1, 4)),
				ast.MakeApplication(
					ast.MakeId(scan.MakeIdToken("myFunction", 1, 0)),
					ast.MakeValue(value.Int(1)),
				),
			)),
		},
		ast.EmptyExpression{},
	),
)

var typeExpected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.MakeProgram(
		[]ast.Statement{
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
		ast.EmptyExpression{},
	),
)

var type2Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.MakeProgram(
		[]ast.Statement{
			ast.MakeTypeDefinition(
				scan.MakeIdToken("Maybe", 0, 0),
				types.MakeData2(
					"Maybe", []string{"a"}, 
					[]types.Constructor{
						types.MakeConstructor("Just", types.Application{types.Tau("a")}),
						types.MakeConstructor("Nothing", types.Application{}),
					},
				),
			),
		},
		ast.EmptyExpression{},
	),
)

var type3Expected = ast.MakePackage(
	DefaultNameSpaceId,
	ast.MakeProgram(
		[]ast.Statement{
			ast.MakeTypeDefinition(
				scan.MakeIdToken("Either", 0, 0),
				types.MakeData2(
					"Either", []string{"a", "b"}, 
					[]types.Constructor{
						types.MakeConstructor("Left", types.Application{types.Tau("a")}),
						types.MakeConstructor("Right", types.Application{types.Tau("b")}),
					},
				),
			),
		},
		ast.EmptyExpression{},
	),
)

var asts = []struct {
	path string
	ast_ parser.Ast
}{
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
	{"./test/fnDef2.yw", fnDef2Expected},//*/
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
	
}

func TestParse(t *testing.T) {
	for _, test := range asts {
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

		/*if test.path == "./test/compose3.yw" {
			ast.PrintAst(prog)
		}//*/
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
