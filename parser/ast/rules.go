package ast

import . "yew/parser/node-type"

// Application ::= Function Expression
var appRule1 = APPLICATION.Replaces(FUNCTION, EXPRESSION)

// Application ::= Expression Expression
var appRule2 = APPLICATION.Replaces(EXPRESSION, EXPRESSION)

// Application ::= Type Type
var appRule3 = APPLICATION.Replaces(TYPE, TYPE)

// Application ::= Type Identifier
var appRule4 = APPLICATION.Replaces(TYPE, IDENTIFIER)

// Assignment ::= Identifier Expression
var assignmentRule = ASSIGNMENT.Replaces(IDENTIFIER, EXPRESSION)

// Operation ::= Expression Infix-Operator Expression
var binaryOperationRule = OPERATION.Replaces(EXPRESSION, BOP_, EXPRESSION)

// Declaration ::= Identifier Type-Annotation
var declarationRule = DECLARATION.Replaces(IDENTIFIER, TYPE_ANNOTATION)

// Definition ::= Declaration Expression
var definitionRule = DEFINITION.Replaces(DECLARATION, EXPRESSION)

// Empty ::= Statement
var emptyRule = EMPTY__.Replaces(STATEMENT)

// Function ::= Declaration Anonymous-Function
var functionRule = FUNCTION.Replaces(DECLARATION, LAMBDA)

// Progress(Function) ::= Anonymous-Function Identifier
var functionRule2 = (IN_PROGRESS__ | FUNCTION).Replaces(LAMBDA, IDENTIFIER)

// Anonymous-Function ::= Binder Expression
var lambdaRule = LAMBDA.Replaces(BINDER, EXPRESSION)

// Anonymous-Function ::= Parameter Expression
var lambdaRule2 = LAMBDA.Replaces(PARAM, EXPRESSION)

// List ::= Sequence
var listRule = LIST.Replaces(SEQUENCE)

// Module ::= Module-Membership Program
var moduleRule = MODULE.Replaces(MODULE_MEMBERSHIP, PROGRAM)

// Module-Membership ::= Identifier
var moduleMembershipRule = MODULE_MEMBERSHIP.Replaces(IDENTIFIER)

// Package-Membership ::= Identifier
var packageMembershipRule = PACKAGE_MEMBERSHIP.Replaces(IDENTIFIER)

// Package ::= Package-Membership Program-Top
var packageRule = PACKAGE.Replaces(PACKAGE_MEMBERSHIP, PROGRAM_TOP)

// Param ::= Expression
var patternParamRule = PARAM.Replaces(EXPRESSION)

// Pattern ::= Expression Anonymous-Function
var patternRule = PATTERN.Replaces(EXPRESSION, LAMBDA)

// Pattern ::= Expression Sequence
var patternRule2 = PATTERN.Replaces(EXPRESSION, SEQUENCE)

// Pattern ::= Expression Program
var patternRule3 = PATTERN.Replaces(EXPRESSION, PROGRAM)

// Postfix-Operation ::= Expression Postfix-Operator
var postOperationRule = POPERATION.Replaces(EXPRESSION, POP_)

// Prefix-Operation ::= Prefix-Operator Expression
var unaryOperationRule = UOPERATION.Replaces(UOP_, EXPRESSION)

// Sequence ::= Statement
var sequenceRuleStatement = SEQUENCE.Replaces(STATEMENT)

// Sequence ::= Sequence Expression
var sequenceRuleContinue = SEQUENCE.Replaces(SEQUENCE, EXPRESSION)

// Sequence ::= Expression
var sequenceRuleNew = SEQUENCE.Replaces(EXPRESSION)

// Tuple ::= Sequence
var tupleRule = TUPLE.Replaces(SEQUENCE)

// Type ::= Type Type
var binaryTypeRule = TYPE.Replaces(TYPE, TYPE)

// Type ::= Type
var justTypeRule = TYPE.Replaces(TYPE)

// Type ::= Type Type
var typeAppRule = TYPE.Replaces(TYPE, TYPE)

// Class ::= Identifier Type
var classRule = CLASS_DEFINITION.Replaces(CLASS_DEFINITION, TYPE_ANNOTATION)

// Type-Annotation ::= Expression Type
var typeAnnotRule = TYPE_ANNOTATION.Replaces(EXPRESSION, TYPE)

// Application ::= Function Expression
var annotationRule = ANNOTATION.Replaces(ANNOTATION, SOMETHING)
