package nodetype

import (
	"strings"
	err "yew/error"
)

type NodeType int

type NodeRule struct {
	Production NodeType
	Expression []NodeType
}

type nodeRuleString struct {
	production string
	expression []string
}

// turns a node rule into a structure resembling a node rule, but with its members as strings
func (rule NodeRule) ToNodeRuleString() nodeRuleString {
	out := nodeRuleString{
		production: rule.Production.ToString(),
		expression: make([]string, len(rule.Expression)),
	}

	for i := range out.expression {
		out.expression[i] = rule.Expression[i].ToString()
	}
	return out
}

func (rule NodeRule) ToString() string {
	sRule := rule.ToNodeRuleString()
	//finalRule := sRule.expression[len(sRule.expression)-1] // final rule's string
	//sRule.expression = sRule.expression[:len(sRule.expression)-1] // removes final rule's string

	var builder strings.Builder
	builder.WriteString(sRule.production)
	builder.WriteString(" ::= ")
	finalIndex := len(sRule.expression) - 1
	for i := 0; i < finalIndex; i++ {
		builder.WriteString(sRule.expression[i])
		builder.WriteByte(' ')
	}
	// write final rule (here to avoid writing extra space)
	builder.WriteString(sRule.expression[finalIndex])
	return builder.String()
}

/*
returns a string representation of a rule with a pointer to the location of the failure

negative integer for input is okay--puts pointer to rule
*/
func (rule NodeRule) RuleFailToString(failedAt int) string {
	if failedAt >= len(rule.Expression) {
		err.PrintBug()
		panic("")
	}

	var builder strings.Builder
	pointerPrefixSpaces := 0 // spaces before pointer
	tot := 0
	tmp := 0
	tmp, _ = builder.WriteString(rule.Production.ToString())
	tot = tot + tmp
	tmp, _ = builder.WriteString(" ::= ")
	tot = tot + tmp

	for i, e := range rule.Expression {
		if i == failedAt {
			pointerPrefixSpaces = tot
		}

		tmp, _ = builder.WriteString(e.ToString())
		builder.WriteByte(' ')
		tot = tot + tmp + 1
	}

	builder.WriteByte('\n')
	bs := make([]byte, pointerPrefixSpaces)
	for i := range bs {
		bs[i] = ' '
	}
	builder.WriteString(string(bs))
	builder.WriteString("^-- failed here")
	return builder.String()
}

// string representations of members of NodeType
var nodeTypeStringMap = map[NodeType]string{
	PROGRAM:            "Program",
	EMPTY__:            "Expression",
	EXPRESSION:         "Expression",
	DEFINITION:         "Definition",
	DECLARATION:        "Declaration",
	VALUE:              "Value",
	ASSIGNMENT:         "Assignment",
	APPLICATION:        "Application",
	OPERATION:          "Operation",
	BOP_:               "Infix-Operator",
	UOP_:               "Prefix-Operator",
	POP_:               "Postfix-Operator",
	UOPERATION:         "Unary-Operation",
	STATEMENT:          "Statement",
	CLASS_DEFINITION:   "Class",
	IDENTIFIER:         "Identifier",
	LAMBDA:             "Anonymous-Function",
	BINDER:             "Binder",
	FUNCTION:           "Function",
	TYPE_ANNOTATION:    "Type-Annotation",
	RETURN:             "Return",
	SEQUENCE:           "Sequence",
	PARAM:              "Parameter",
	PACKAGE:            "Package",
	MODULE:             "Module",
	PACKAGE_MEMBERSHIP: "Package-Membership",
	MODULE_MEMBERSHIP:  "Module-Membership",
	LIST:               "List",
	TUPLE:              "Tuple",
	TYPE_DEF:          "Type-Definition",
	TYPE:               "Type",
	STACK_MARKER:       "",
	REPEAT__:           "",
	REPEAT_OR_NONE__:   "",
}

func (t NodeType) ToString() string {
	if t >= IN_PROGRESS__ {
		t = t ^ IN_PROGRESS__
		return "Progress(" + t.ToString() + ")"
	}

	stringified, found := nodeTypeStringMap[t]
	if !found {
		stringified = ""
	}
	return stringified
}

const (
	PROGRAM NodeType = iota
	PROGRAM_TOP
	EXPRESSION
	EMPTY__
	DEFINITION
	DECLARATION
	VALUE
	ASSIGNMENT
	APPLICATION
	OPERATION
	BOP_
	UOP_
	POP_
	UOPERATION
	POPERATION
	STATEMENT
	CLASS_DEFINITION
	IDENTIFIER
	LAMBDA
	BINDER
	FUNCTION
	TYPE_ANNOTATION
	RETURN
	SEQUENCE
	PARAM
	STACK_MARKER
	PACKAGE
	MODULE
	PACKAGE_MEMBERSHIP
	MODULE_MEMBERSHIP
	LIST
	TUPLE
	TYPE
	TYPE_DEF
	REPEAT__
	REPEAT_OR_NONE__
	IN_PROGRESS__ NodeType = 0x80
)
