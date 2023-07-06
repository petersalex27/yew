package ast

import (
	"fmt"
	"os"
	"testing"
	err "yew/error"
	"yew/info"
	scan "yew/lex"
	"yew/parser/parser"
	types "yew/type"
	"yew/value"
)

var inSource = []string{"test1 test2 Int\n", "someId + 1"}
var in = scan.CreateInputStream(
	"test/ast", 0, inSource,
	idTok1, idTok2, scan.MakeOtherToken(scan.INT, 1, 13),
	scan.MakeOtherToken(scan.NEW_LINE, 1, 16),
	idTok3, plusTok, valTok1,
)
var idTok1 = scan.MakeIdToken("test1", 1, 1)
var idTok2 = scan.MakeIdToken("test2", 1, 7)
var idTok3 = scan.MakeIdToken("someId", 2, 1)
var id1 = MakeIdWithType(idTok1, types.Int{})
var id2 = MakeIdWithType(idTok2, types.Int{})
var id3 = MakeIdWithType(idTok3, types.Int{})
var valTok1 = scan.ValueToken{Value: value.Int(1), Line: 2, Char: 10}
var plusTok = scan.MakeOtherToken(scan.PLUS, 2, 8)
var val1 = Value(valTok1)
var op = OpType(plusTok)
var dec1 = Declaration{Qualifier: LetDeclare, id: id1}
var dec2 = Declaration{Qualifier: LetDeclare, id: id2}
var dec3 = Declaration{Qualifier: LetDeclare, id: id3}
var stack1 = parser.AstStack{dec1, dec2}

func genExpected(tok scan.Token, found string, expected string, source []string) string {
	lc := tok.GetLocation()
	msg := "found " + found + " but expected " + expected
	return err.CompileMessage(msg, err.ERROR, err.SYNTAX, "test/ast", lc.GetLine(), lc.GetChar(), source).ToString()
}

func TestApplication(t *testing.T) {
	{
		expected := genExpected(idTok1, "a declaration", "a function", inSource)
		ok, e := stack1.Validate(appRule1)
		if ok {
			fmt.Fprintf(os.Stderr, "expected validation to fail\n")
			t.FailNow()
		}
		es := e(in).ToString()
		if es != expected {
			fmt.Fprintf(os.Stderr, "Expected (len=%d):\n%s\n", len(expected), expected)
			fmt.Fprintf(os.Stderr, "Actual (len=%d):\n%s\n", len(es), es)
			t.FailNow()
		}
	}

	{
		expected := genExpected(idTok1, "a declaration", "an expression", inSource)
		ok, e := stack1.Validate(appRule2)
		if ok {
			fmt.Fprintf(os.Stderr, "expected validation to fail\n")
			t.FailNow()
		}
		es := e(in).ToString()
		if es != expected {
			fmt.Fprintf(os.Stderr, "Expected (len=%d):\n%s\n", len(expected), expected)
			fmt.Fprintf(os.Stderr, "Actual (len=%d):\n%s\n", len(es), es)
			t.FailNow()
		}
	}
}

func TestAssignment(t *testing.T) {
	{
		expected := genExpected(idTok1, "a declaration", "an identifier", inSource)
		ok, e := stack1.Validate(assignmentRule)
		if ok {
			fmt.Fprintf(os.Stderr, "expected validation to fail\n")
			t.FailNow()
		}
		es := e(in).ToString()
		if es != expected {
			fmt.Fprintf(os.Stderr, "Expected (len=%d):\n%s\n", len(expected), expected)
			fmt.Fprintf(os.Stderr, "Actual (len=%d):\n%s\n", len(es), es)
			t.FailNow()
		}
	}

	{
		expected := genExpected(idTok1, "a declaration", "an expression", inSource)
		ok, e := (&parser.AstStack{id1, dec1}).Validate(assignmentRule)
		if ok {
			fmt.Fprintf(os.Stderr, "expected validation to fail\n")
			t.FailNow()
		}
		es := e(in).ToString()
		if es != expected {
			fmt.Fprintf(os.Stderr, "Expected (len=%d):\n%s\n", len(expected), expected)
			fmt.Fprintf(os.Stderr, "Actual (len=%d):\n%s\n", len(es), es)
			t.FailNow()
		}
	}
}

func TestBinaryOp(t *testing.T) {
	{
		expected := genExpected(idTok3, "a declaration", "an expression", inSource)
		ok, e := (&parser.AstStack{dec3, op, val1}).Validate(binaryOperationRule)
		if ok {
			fmt.Fprintf(os.Stderr, "expected validation to fail\n")
			t.FailNow()
		}
		es := e(in).ToString()
		if es != expected {
			fmt.Fprintf(os.Stderr, "Expected (len=%d):\n%s\n", len(expected), expected)
			fmt.Fprintf(os.Stderr, "Actual (len=%d):\n%s\n", len(es), es)
			t.FailNow()
		}
	}
	{
		expected := genExpected(valTok1, "a literal value", "an infix operator", inSource)
		ok, e := (&parser.AstStack{val1, val1, val1}).Validate(binaryOperationRule)
		if ok {
			fmt.Fprintf(os.Stderr, "expected validation to fail\n")
			t.FailNow()
		}
		es := e(in).ToString()
		if es != expected {
			fmt.Fprintf(os.Stderr, "Expected (len=%d):\n%s\n", len(expected), expected)
			fmt.Fprintf(os.Stderr, "Actual (len=%d):\n%s\n", len(es), es)
			t.FailNow()
		}
	}
	{
		expected := genExpected(plusTok, "an infix operator", "an expression", inSource)
		ok, e := (&parser.AstStack{val1, op, op}).Validate(binaryOperationRule)
		if ok {
			fmt.Fprintf(os.Stderr, "expected validation to fail\n")
			t.FailNow()
		}
		es := e(in).ToString()
		if es != expected {
			fmt.Fprintf(os.Stderr, "Expected (len=%d):\n%s\n", len(expected), expected)
			fmt.Fprintf(os.Stderr, "Actual (len=%d):\n%s\n", len(es), es)
			t.FailNow()
		}
	}
}

func TestClass(t *testing.T) {
	var inSource = []string {
		"class MyClass a where\n",
		"  fn :: a -> Int;\n",
		"  gn :: a",
	}
	var inSource2 = []string{
		"class MyClass a where\n",
		"  fn :: a -> Int;\n",
		"  fn :: a -> Int",
	}
	var in = scan.CreateInputStream(
		"test/ast-1", 0, inSource,
	)
	var in2 = scan.CreateInputStream(
		"test/ast-2", 0, inSource2,
	)
	var expectMsg = 
		"[test/ast-1:3:9] Type Error: unexpected type, expected a function type.\n" +
		"    3 |   gn :: a\n" +
		"                ^"
	var expectMsg2 = 
		"[test/ast-2:3:3] Name Error: illegal redefinition of fn in the MyClass class.\n" +
		"    3 |   fn :: a -> Int\n" +
		"          ^"
	
	var myClassToken = scan.MakeIdToken("MyClass", 1, 6)
	var fnIdToken1 = scan.MakeIdToken("fn", 2, 3)
	var fnIdToken2 = scan.MakeIdToken("fn", 3, 3)
	var gnIdToken = scan.MakeIdToken("gn", 3, 3)
	var classParamType1 = types.MakeTau("a", info.MakeLocation(2, 9))
	var intType1 = types.Int(info.MakeLocation(2, 14))
	var classParamType2 = types.MakeTau("a", info.MakeLocation(3, 9))
	var intType2 = types.Int(info.MakeLocation(3, 14))

	class := Class{
		name: MakeId(myClassToken),
		functions: make(map[string]types.Function),
	}
	ty1 := types.Function{
		Domain: classParamType1,
		Codomain: intType1,
	}
	ty2 := types.Function{
		Domain: classParamType2,
		Codomain: intType2,
	}
	annotBad := ExpressionTypeAnnotation{
		expression: MakeId(gnIdToken),
		expressionType: classParamType2,
	}
	annotGood := ExpressionTypeAnnotation{
		expression: MakeId(fnIdToken2),
		expressionType: ty2,
	}
	class.functions[fnIdToken1.ToString()] = ty1
	
	{
		stack := parser.AstStack{class, annotBad}
		p := parser.Parser{Input: in, Stack: &stack}
		ok, e := constructClass(&p)
		if ok {
			fmt.Fprintf(os.Stderr, "expected validation to fail.\n")
			t.FailNow()
		}
		actual := e.ToString()
		if actual != expectMsg {
			fmt.Fprintf(os.Stderr, "Expected:\n%s\nActual:\n%s\n", expectMsg, actual)
			t.FailNow()
		}
	}

	{
		stack := parser.AstStack{class, annotGood}
		p := parser.Parser{Input: in2, Stack: &stack}
		ok, e := constructClass(&p)
		if ok {
			fmt.Fprintf(os.Stderr, "expected validation to fail.\n")
			t.FailNow()
		}
		actual := e.ToString()
		if actual != expectMsg2 {
			fmt.Fprintf(os.Stderr, "Expected:\n%s\nActual:\n%s\n", expectMsg2, actual)
			t.FailNow()
		}
	}
}

/*func TestDeclaration(t *testing.T) {

}

func TestDefinition(t *testing.T) {

}

func TestEmpty(t *testing.T) {

}

func TestExpression(t *testing.T) {

}

func TestFunction(t *testing.T) {

}

func TestId(t *testing.T) {

}

func TestLambda(t *testing.T) {

}

func TestList(t *testing.T) {

}

func TestModule(t *testing.T) {

}

func TestPackage(t *testing.T) {

}

func TestParam(t *testing.T) {

}

func TestPattern(t *testing.T) {

}

func TestPostfix(t *testing.T) {

}

func TestPrefix(t *testing.T) {

}

func TestProgram(t *testing.T) {

}

func TestSequence(t *testing.T) {

}

func TestStatement(t *testing.T) {

}

func TestTuple(t *testing.T) {

}

func TestTypeAnnotation(t *testing.T) {

}

func TestTypeDef(t *testing.T) {

}

func TestType(t *testing.T) {

}

func TestValue(t *testing.T) {

}*/
