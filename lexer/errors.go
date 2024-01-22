// =================================================================================================
// Alex Peters - January 22, 2024
//
// errors specific to lexical analysis step
// =================================================================================================
package lexer

import "github.com/petersalex27/yew/errors"

func makeOSError(msg string) errors.ErrorMessage {
	return errors.MakeError("OS", msg)
}

func makeLexicalError(msg string, src pathSpec, line, start, end int) errors.ErrorMessage {
	e := errors.MakeError("Lexical", msg, line, start, end)
	e.SourceName = src.String()
	return e
}