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

func (nt NodeType) Replaces(nts ...NodeType) NodeRule {
	return NodeRule{Production: nt, Expression: nts}
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

func GetErrorName(nt NodeType) (article string, name string) {
	switch nt {
	case PROGRAM:
		return "a", "program"
	case EXPRESSION:
		return "an", "expression"
	case DEFINITION:
		return "a", "definition"
	case DECLARATION:
		return "a", "declaration"
	case VALUE:
		return "a", "literal value"
	case ASSIGNMENT:
		return "an", "assignment statement"
	case APPLICATION:
		return "an", "application"
	case OPERATION:
		return "an", "infix operation"
	case POPERATION:
		return "a", "postfix operation"
	case BOP_:
		return "an", "infix operator"
	case UOP_:
		return "a", "prefix operator"
	case POP_:
		return "a", "postfix operator"
	case UOPERATION:
		return "a", "unary operation"
	case STATEMENT:
		return "a", "statment"
	case CLASS_DEFINITION:
		return "a", "class definition"
	case IDENTIFIER:
		return "an", "identifier"
	case LAMBDA:
		return "an", "anonymous function"
	case BINDER:
		return "a", "binder"
	case FUNCTION:
		return "a", "function"
	case TYPE_ANNOTATION:
		return "a", "type annotation"
	case RETURN:
		return "a", "return statement"
	case SEQUENCE:
		return "a", "sequence"
	case PARAM:
		return "a", "parameter"
	case PACKAGE:
		return "a", "package declaration"
	case MODULE:
		return "a", "module declaration"
	case PACKAGE_MEMBERSHIP:
		return "a", "package membership declaration"
	case MODULE_MEMBERSHIP:
		return "a", "module membership declaration"
	case LIST:
		return "a", "list"
	case TUPLE:
		return "a", "tuple"
	case TYPE_DEF:
		return "a", "type definition"
	case TYPE:
		return "a", "type"
	case PATTERN:
		return "a", "pattern"
	case ANNOTATION:
		return "an", "annotation"
	case SOMETHING:
		return "", "something"
	default:
		return "", ""
	}
}

func (rule NodeRule) ExpectedFailure(failedAt int) string {
	article, name := GetErrorName(rule.Expression[failedAt])
	if article == "" {
		err.PrintBug()
		panic("")
	}

	res := "expected " + article + " " + name
	return res
}

func FoundFailure(nt NodeType) string {
	article, name := GetErrorName(nt)
	if article == "" {
		return ""
	}

	return "found " + article + " " + name
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
	TYPE_DEF:           "Type-Definition",
	TYPE:               "Type",
	PATTERN:            "Pattern",
	ANNOTATION:         "Annotation",
	SOMETHING:          "*",
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
	PATTERN
	ANNOTATION
	SOMETHING
	REPEAT__
	REPEAT_OR_NONE__
	IN_PROGRESS__ NodeType = 0x80
)
