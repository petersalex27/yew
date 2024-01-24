// =================================================================================================
// Alex Peters - January 22, 2024
//
// errors specific to lexical analysis step
// =================================================================================================
package lexer

import "github.com/petersalex27/yew/errors"

const (
	InvalidCharacter string = "invalid character"
	InvalidAffixId string = "invalid affixed id"
	InvalidCharacterAtEndOfNumConst string = "invalid character at end of numeric constant"
	UnexpectedEOF string = "unexpected end of file"
	ExpectedCharLiteral string = "expected character literal"
	IllegalCharLiteral string = "illegal character literal"
	IllegalEscapeSequence string = "illegal escape sequence"
	IllegalStringLiteral string = "illegal string literal"
	IllegalUnderscoreSequence string = "illegal contiguous sequence of underscores"
	InvalidUnderscore string = "invalid underscore"
)

// creates a OS error from the given message
func makeOSError(msg string) errors.ErrorMessage {
	return errors.MakeError("OS", msg)
}

// creates a lexical error from the arguments
func makeLexicalError(msg string, src pathSpec, line, start, end int) errors.ErrorMessage {
	e := errors.MakeError("Lexical", msg, line, start, end)
	e.SourceName = src.String()
	return e
}

// adds an error for some section of code located at `lex.path`, `line` number, `start` char number,
// `end` char number
func (lex *Lexer) error2(msg string, line, start, end int) {
	e := makeLexicalError(msg, lex.path, line, start, end)
	lex.addMessage(e)
}

// adds an error constructed using lexer's data and the message string passed as an argument
func (lex *Lexer) error(msg string) {
	start, _ := lex.SavedChar.Pop()
	e := makeLexicalError(msg, lex.path, lex.Line, start, lex.Char)
	lex.addMessage(e)
}