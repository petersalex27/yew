package errorgen

import (
	"strings"
	err "yew/error"
	scan "yew/lex"
)

type GenerateErrorFunction func (scan.Token, scan.InputStream) err.Error

func GenerateError(message string, errSub err.ErrorSubType) GenerateErrorFunction {
	return func(t scan.Token, i scan.InputStream) err.Error {
		loc := t.GetLocation()
		return err.CompileMessage(
			message, err.ERROR, errSub, (i).GetPath(), loc.GetLine(), loc.GetChar(),
			(i).GetSource()).(err.Error)
	}
}

type SyntaxMessage string

const (
	TypeDecExpectsEqual SyntaxMessage = "unexpected token, expected equal token after type declaration"
	ExpectedLCurl SyntaxMessage = "unexpected token, expected an open curly brace"
	ExpectedWhere SyntaxMessage = "unexpected token, expected `where`"
	ExpectedFunctionDeclaration SyntaxMessage = "expected function declaration"
	ExpectedInstanceFunctionDeclaration SyntaxMessage = "expected instance function declaration"
	ExpectedInstanceFunctionDefinition SyntaxMessage = "expected instance function definition"
	//NonClassAsClass SyntaxMessage = "cannot define an instance function for"
)

type SyntaxMessageComplex struct{argc int; msg []string}

var NonClassAsClass = SyntaxMessageComplex{
	2, []string{
		"cannot define an instance function for ", 
		"", 
		" because ",
		"",
		" is not a class",
	},
}
func (m SyntaxMessageComplex) Generate(args ...string) GenerateErrorFunction {
	argsIndex := 0
	if len(args) != m.argc {
		err.PrintBug() // arg count mismatch 
		panic("")
	}

	var builder strings.Builder
	for _, s := range m.msg {
		if s == "" {
			builder.WriteString(args[argsIndex])
			argsIndex++
		} else {
			builder.WriteString(s)
		}
	}
	return GenerateSyntaxError(builder.String())
}

type TypeMessage string
const (
	ExpectedTypeIdentifier TypeMessage = "expected type identifier"
	ExpectedTypeIdentifierNotVar TypeMessage = "expected type identifier but found type variable"
	UnexpectedType TypeMessage = "unexpected type"
	RedeclaredClassInstance TypeMessage = "cannot redeclare type class instance"
)

type NameMessage string
const (
	RedeclaredConstructor NameMessage = "cannot redeclare type constructor"
	ClassNotFound NameMessage = "intance function defined for a class that is not declared"
	RedeclaredClass NameMessage = "cannot redeclare type class"
	UndefinedClass NameMessage = "class not defined"
	FunctionNotInClass NameMessage = "function is not declared in type class being instantiated"
	FunctionInstanceRedefined NameMessage = "instance function redefined"
)

func (m NameMessage) Generate() GenerateErrorFunction {
	return GenerateNameError(string(m))
}

func (m TypeMessage) Generate() GenerateErrorFunction {
	return GenerateTypeError(string(m))
}

func (m SyntaxMessage) Generate() GenerateErrorFunction {
	return GenerateSyntaxError(string(m))
}

func GenerateSyntaxError(message string) GenerateErrorFunction {
	return GenerateError(message, err.SYNTAX)
}

func GenerateTypeError(message string) GenerateErrorFunction {
	return GenerateError(message, err.TYPE)
}

func GenerateNameError(message string) GenerateErrorFunction {
	return GenerateError(message, err.NAME)
}