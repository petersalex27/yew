// =================================================================================================
// Alex Peters - January 22, 2024
//
// errors specific to lexical analysis step
// =================================================================================================
package lexer

import "github.com/petersalex27/yew/errors"

const (
	InvalidCharacter string = "invalid character"
)

func makeOSError(msg string) errors.ErrorMessage {
	return errors.MakeError("OS", msg)
}

func makeLexicalError(msg string, src pathSpec, line, start, end int) errors.ErrorMessage {
	e := errors.MakeError("Lexical", msg, line, start, end)
	e.SourceName = src.String()
	return e
}

// adds an error constructed using lexer's data and the message string passed as an argument
func (lex *Lexer) error(msg string) {
	start, _ := lex.SavedChar.Peek()
	e := makeLexicalError(msg, lex.path, lex.Line, start, lex.Char)
	lex.addMessage(e)
}