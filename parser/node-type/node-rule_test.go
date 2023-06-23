package nodetype

import (
	"testing"
	"yew"
)

type String string

func (s String) ToString() string {
	return string(s)
}

func (s String) Equals(z String) bool {
	return string(s) == string(z)
}

func TestNodeTypeToString(t *testing.T) {
	tests := []struct{in NodeType; expected string; actual string}{
		{PROGRAM, "Program", ""},
		{EXPRESSION, "Expression", ""},
		{EMPTY__, "Expression", ""},
		{DEFINITION, "Definition", ""},
		{DECLARATION, "Declaration", ""},
		{VALUE, "Value", ""},
		{ASSIGNMENT, "Assignment", ""},
		{APPLICATION, "Application", ""},
		{OPERATION, "Operation", ""},
		{BOP_, "Infix-Operator", ""},
		{UOP_, "Prefix-Operator", ""},
		{POP_, "Postfix-Operator", ""},
		{UOPERATION, "Unary-Operation", ""},
		{STATEMENT, "Statement", ""},
		{CLASS_DEFINITION, "Class", ""},
		{IDENTIFIER, "Identifier", ""},
		{LAMBDA, "Anonymous-Function", ""},
		{BINDER, "Binder", ""},
		{FUNCTION, "Function", ""},
		{TYPE_ANNOTATION, "Type-Annotation", ""},
		{RETURN, "Return", ""},
		{SEQUENCE, "Sequence", ""},
		{PARAM, "Parameter", ""},
		{STACK_MARKER, "", ""},
		{PACKAGE, "Package", ""},
		{MODULE, "Module", ""},
		{LIST, "List", ""},
		{TUPLE, "Tuple", ""},
		{TYPE_DEF, "Type-Definition", ""},
		{TYPE, "Type", ""},
		{REPEAT__, "", ""},
		{REPEAT_OR_NONE__, "", ""},
	}

	for i, test := range tests {
		tests[i].actual = test.in.ToString()
		if tests[i].actual != tests[i].expected {
			yew.PrintExpectedVsActual(tests[i].expected, tests[i].actual)
			t.FailNow()
		}
	}
}

func TestNodeRuleToString(t *testing.T) {
	rule := NodeRule{PROGRAM, []NodeType{STATEMENT, EXPRESSION}}
	rule2 := NodeRule{IN_PROGRESS__ | PROGRAM, []NodeType{STATEMENT, EXPRESSION}}
	expect := "Program ::= Statement Expression"
	expect2 := "Progress(Program) ::= Statement Expression"
	if rule.ToString() != expect {
		yew.PrintExpectedVsActual(expect, rule.ToString())
		t.FailNow()
	}

	if rule2.ToString() != expect2 {
		yew.PrintExpectedVsActual(expect2, rule2.ToString())
		t.FailNow()
	}
}