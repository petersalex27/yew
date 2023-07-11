package errorgen

import scan "yew/lex"
import err "yew/error"

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
)

type TypeMessage string
const (
	ExpectedTypeIdentifier TypeMessage = "expected type identifier"
	ExpectedTypeIdentifierNotVar TypeMessage = "expected type identifier but found type variable"
	UnexpectedType TypeMessage = "unexpected type"
)

type NameMessage string
const (
	RedeclaredConstructor NameMessage = "cannot redeclare type constructor"
	ClassNotFound NameMessage = "intance function defined for a class that is not declared"
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