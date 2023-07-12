package ast

import (
	"fmt"
	"os"
	"testing"
	"yew/info"
	scan "yew/lex"
	errorgen "yew/parser/error-gen"
	"yew/parser/parser"
	"yew/type"
	"yew/value"
)

var testClass = Class{
	name: MakeId(scan.MakeIdToken("MyClass", 1, 7)),
	typeParameter: types.Var("a"),
	functions: map[string]types.Function{
		"f": types.Function{
			Domain: types.Var("a"),
			Codomain: types.Int(info.DefaultLoc()),
		},
		"g": types.Function{
			Domain: types.Var("a"),
			Codomain: types.Var("a"),
		},
	},
}

var testFunction_f = MakeFunction(
	MakeId(scan.MakeIdToken("f", 6, 5)),
	MakeLambda(
		MakeParameter(0,
			MakeTypeAnnotation(
				MakeId(scan.MakeIdToken("x", 6, 7)),
				types.Bool{},
			),
		),
		MakeTypeAnnotation(
			MakeValue(value.Int(1)),
			types.Int{},
		),
	),
)

var testFunction_g = MakeFunction(
	MakeId(scan.MakeIdToken("g", 7, 5)),
	MakeLambda(
		MakeParameter(0,
			MakeTypeAnnotation(
				MakeId(scan.MakeIdToken("x", 7, 7)),
				types.Bool{},
			),
		),
		MakeTypeAnnotation(
			MakeId(scan.MakeIdToken("x", 7, 11)),
			types.Bool{},
		),
	),
)

var testFunction_g_bad = MakeFunction(
	MakeId(scan.MakeIdToken("g", 8, 5)),
	MakeLambda(
		MakeParameter(0,
			MakeTypeAnnotation(
				MakeId(scan.MakeIdToken("x", 7, 7)),
				types.Bool{},
			),
		),
		MakeTypeAnnotation(
			MakeId(scan.MakeIdToken("x", 7, 11)),
			types.Bool{},
		),
	),
)

var testFunction_bad = MakeFunction(
	MakeId(scan.MakeIdToken("h", 9, 5)),
	MakeLambda(
		MakeParameter(0,
			MakeTypeAnnotation(
				MakeId(scan.MakeIdToken("x", 7, 7)),
				types.Bool{},
			),
		),
		MakeTypeAnnotation(
			MakeId(scan.MakeIdToken("x", 7, 11)),
			types.Bool{},
		),
	),
)

var p = parser.Parser{
	Input: scan.CreateInputStream("test/class", 0, []string{
		"class MyClass a where {\n",
		"    f :: a -> Int\n",
		"    g :: a -> a\n",
		"}\n",
		"MyClass => Bool where {\n",
		"    f x = 1\n",
		"    g x = x\n",
		"    g x = x\n",
 		"    h x = x\n",
		"}",
	}, ),
}

var boolInstance = types.Bool(info.MakeLocation(5, 12))

func TestDeclareClass(t *testing.T) {
	classTable := ClassTable{}.InitClassTable()
	
	if !classTable.DeclareClass(&p, testClass) {
		t.FailNow()
	}

	if !classTable.DeclareInstance(&p, testClass, boolInstance) {
		t.FailNow()
	}

	if !classTable.DefineInstanceFunction(&p, testClass, boolInstance, testFunction_f) {
		t.FailNow()
	}

	if !classTable.DefineInstanceFunction(&p, testClass, boolInstance, testFunction_g) {
		t.FailNow()
	}
}

var classIdTokenForInstance = scan.MakeIdToken("MyClass", 5, 1)

// tests that you cannot get a class that isn't defined and tests for correct error message
func Test_getClass(t *testing.T) {
	expected :=
		"[test/class:5:1] Name Error: class not defined.\n" +
		"    5 | MyClass => Bool where {\n\n" + 
		"        ^"
	classTable_ := ClassTable{}.InitClassTable()
	classTable := classTable_.(ClassTable)
	_, errFn := classTable.getClass("MyClass")
	if errFn == nil {
		fmt.Fprintf(os.Stderr, "expected getClass to fail\n")
		t.FailNow()
	}
	actual := errFn(classIdTokenForInstance, p.Input).ToString()
	if actual != expected {
		fmt.Fprintf(os.Stderr, "Expected:\n%s\nActual:\n%s\n",
			expected, actual)
		t.FailNow()
	}
}

func Test_checkUninstantiated(t *testing.T) {
	expected :=
		"[test/class:5:12] Type Error: " + string(errorgen.RedeclaredClassInstance) + ".\n" +
		"    5 | MyClass => Bool where {\n\n" + 
		"                   ^"
	classTable_ := ClassTable{}.InitClassTable()
	classTable := classTable_.(ClassTable)
	if !classTable.DeclareClass(&p, testClass) {
		fmt.Fprintf(os.Stderr, "failed to declare class\n")
		t.FailNow()
	}

	entry, found := classTable.GetClass(&p, testClass)
	if !found {
		t.FailNow()
	}
	if !classTable.DeclareInstance(&p, testClass, boolInstance) {
		t.FailNow()
	}

	errFn := entry.checkUninstantiated("Bool")
	if errFn == nil {
		fmt.Fprintf(os.Stderr, "expected checkUinstantiated to fail\n")
		t.FailNow()
	}

	loc := boolInstance.GetLocation()
	dummyToken := scan.MakeOtherToken(scan.TYPE_ID, loc.GetLine(), loc.GetChar())
	actual := errFn(dummyToken, p.Input).ToString()
	if actual != expected {
		fmt.Fprintf(os.Stderr, "Expected:\n%s\nActual:\n%s\n",
			expected, actual)
		t.FailNow()
	}
}

// test for correct error when trying to define a function that is not declared in the class definition
func Test_confirmFunctionDeclared(t *testing.T) {
	expected :=
		"[test/class:9:5] Name Error: " + string(errorgen.FunctionNotInClass) + ".\n" +
		"    9 |     h x = x\n\n" + 
		"            ^"
	classTable_ := ClassTable{}.InitClassTable()
	classTable := classTable_.(ClassTable)
	if !classTable.DeclareClass(&p, testClass) {
		fmt.Fprintf(os.Stderr, "failed to declare class\n")
		t.FailNow()
	}

	entry, found := classTable.GetClass(&p, testClass)
	if !found {
		t.FailNow()
	}
	if !classTable.DeclareInstance(&p, testClass, boolInstance) {
		t.FailNow()
	}

	errFn := entry.confirmFunctionDeclared(testFunction_bad)
	if errFn == nil {
		fmt.Fprintf(os.Stderr, "expected confirmFunctionDeclared to fail\n")
		t.FailNow()
	}

	actual := errFn(testFunction_bad.FindStartToken(), p.Input).ToString()
	if actual != expected {
		fmt.Fprintf(os.Stderr, "Expected:\n%s\nActual:\n%s\n",
			expected, actual)
		t.FailNow()
	}
}

func Test_confirmUniqueDefinition(t *testing.T) {
	expected :=
		"[test/class:8:5] Name Error: " + string(errorgen.FunctionInstanceRedefined) + ".\n" +
		"    8 |     g x = x\n\n" + 
		"            ^"

	inst := map[string]Function{
		"g": testFunction_g,
	}

	errFn := confirmUniqueDefinition(inst, testFunction_g_bad)
	if errFn == nil {
		fmt.Fprintf(os.Stderr, "expected confirmUniqueDefintion to fail\n")
		t.FailNow()
	}

	actual := errFn(testFunction_g_bad.FindStartToken(), p.Input).ToString()
	if actual != expected {
		fmt.Fprintf(os.Stderr, "Expected:\n%s\nActual:\n%s\n",
			expected, actual)
		t.FailNow()
	}
}