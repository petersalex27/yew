package errorgen

import "yew/lex"
import err "yew/error"

type generateErrorFunction func (scan.Token, scan.InputStream) err.Error

func GenerateError(message string, errSub err.ErrorSubType) generateErrorFunction {
	return func(t scan.Token, i scan.InputStream) err.Error {
		loc := t.GetLocation()
		return err.CompileMessage(
			message, err.ERROR, errSub, (i).GetPath(), loc.GetLine(), loc.GetChar(),
			(i).GetSource()).(err.Error)
	}
}

func GenerateSyntaxError(message string) generateErrorFunction {
	return GenerateError(message, err.SYNTAX)
}

func GenerateTypeError(message string) generateErrorFunction {
	return GenerateError(message, err.TYPE)
}