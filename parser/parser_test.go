package parsing

import (
	"fmt"
	"os"
	"testing"
	"yew/lex"
	"yew/ast"
	"yew/symbol"
	"yew/type"
	"yew/value"
	//"yew/info"
)

var defExpected = 
	ast.MakeProgram(
		[]ast.Definition{
			ast.MakeDefinition(
				ast.MakeDeclaration(
					symbol.MakeSymbol_testable(
						"x", 
						types.Int{}, 
						symbol.MakeLocation("./test/def.yw", 1, 4), 
						make(map[string]symbol.SymbolUse)), ),
				ast.MakeAssignment(
					ast.MakeValue(value.Int(1))), ),
		},
		ast.EmptyExpression{})

var def2Expected =
	ast.MakeProgram(
		[]ast.Definition{
			ast.MakeDefinition(
				ast.MakeDeclaration(
					symbol.MakeSymbol_testable(
						"x",
						types.Tau(".ty0"),
						symbol.MakeLocation("./test/def2.yw", 1, 4),
						make(map[string]symbol.SymbolUse)), ),
				ast.MakeAssignment(
					ast.MakeValue(value.Int(1))), ),
		},
		ast.EmptyExpression{})

/*var appExpected =
		ast.MakeProgram(
			[]ast.Definition{},
			ast.MakeApplication(
				ast.MakeId(symbol.MakeSymbol_testable(
					"myFunction",
					
				))
			)
		)*/

var asts = []struct {path string; ast_ ast.Ast} {
	{"./test/def.yw", defExpected},
	{"./test/def2.yw", def2Expected},
	//{"./test/app.yw", appExpected},
	//{"./test/app2.yw", app2Expected},
}

func TestParse(t *testing.T) {
	for _, test := range asts { 
		fmt.Fprintf(os.Stderr, "%s ==============\n", test.path)
		in, e := scan.Init(test.path)
		if nil != e {
			fmt.Fprintf(os.Stderr, "%s\n", e.Error())
			t.FailNow()
		}

		ok, prog := Parse(&in)
		if !ok {
			t.FailNow()
		}

		if !ast.EqualTest(prog, test.ast_) {
			fmt.Printf("Expected: \n")
			ast.PrintAst(test.ast_)

			fmt.Printf("Actual: \n")
			ast.PrintAst(prog)
			t.FailNow()
		}
	}
}