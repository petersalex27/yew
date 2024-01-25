// =================================================================================================
// Alex Peters - January 22, 2024
//
// errors specific to lexical analysis step
// =================================================================================================
package lexer

import (
	"github.com/petersalex27/yew/common"
	"github.com/petersalex27/yew/errors"
)

const (
	InvalidCharacter                string = "invalid character"
	InvalidAffixId                  string = "invalid affixed id"
	InvalidCharacterAtEndOfNumConst string = "invalid character at end of numeric constant"
	UnexpectedEOF                   string = "unexpected end of file"
	ExpectedCharLiteral             string = "expected character literal"
	IllegalCharLiteral              string = "illegal character literal"
	IllegalEscapeSequence           string = "illegal escape sequence"
	IllegalStringLiteral            string = "illegal string literal"
	IllegalUnderscoreSequence       string = "illegal contiguous sequence of underscores"
	InvalidUnderscore               string = "invalid underscore"
	InvalidAnnotation               string = "invalid annotation"
)

// creates a OS error from the given message
func makeOSError(msg string) errors.ErrorMessage {
	return errors.MakeError("OS", msg)
}

// creates a lexical error from the arguments
func makeLexicalError(msg string, src pathSpec, line, lineEnd, start, end int) errors.ErrorMessage {
	e := errors.MakeError("Lexical", msg, line, lineEnd, start, end)
	e.SourceName = src.String()
	return e
}

// adds an error for some section of code located at `lex.path`, `line` number, `start` char number,
// `end` char number
func (lex *Lexer) error2(msg string, line, lineEnd, start, end int) {
	e := makeLexicalError(msg, lex.path, line, lineEnd, start, end)
	lex.addMessage(e)
}

// adds an error constructed using lexer's data and the message string passed as an argument
func (lex *Lexer) error(msg string) {
	start, _ := lex.SavedChar.Pop()
	line1 := 1 + common.SearchRange(lex.PositionRanges, start)   // 1 + result = 0 or greater
	line2 := 1 + common.SearchRange(lex.PositionRanges, lex.Pos) // 1 + result = 0 or greater
	char1, char2 := 0, 0
	if line1 > 0 {
		char1 = (lex.PositionRanges[line1-1] + 1) - start
	}
	if line2 > 0 {
		char2 = (lex.PositionRanges[line2-1] + 1) - lex.Pos
	}
	e := makeLexicalError(msg, lex.path, line1, line2, char1, char2)
	lex.addMessage(e)
}
