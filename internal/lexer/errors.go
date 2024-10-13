// =================================================================================================
// Alex Peters - January 22, 2024
//
// errors specific to lexical analysis step
// =================================================================================================
package lexer

import (
	"github.com/petersalex27/yew/internal/errors"
	"github.com/petersalex27/yew/api/token"
)

const (
	UnexpectedEOF         string = "unexpected end of file"
	ExpectedCharLiteral   string = "expected character literal"
	IllegalCharLiteral    string = "illegal character literal"
	IllegalEscapeSequence string = "illegal escape sequence"
	IllegalStringLiteral  string = "illegal string literal"
	UnexpectedUnderscore  string = "unexpected underscore"
	UnexpectedSymbol      string = "unexpected symbol"
	ExpectedAnnotationId  string = "annotation must have an identifier"
)

// adds an error constructed using lexer's data and the message string passed as an argument
func (lex *Lexer) error(msg string) token.Token {
	start, _ := lex.SavedChar.Pop()
	// endPositions := lex.EndPositions()
	// line1 := 1 + common.SearchRange(endPositions, start, false)  // 1 + result = 0 or greater
	// line2 := 1 + common.SearchRange(endPositions, lex.Pos, true) // 1 + result = 0 or greater
	// char1, char2 := 0, 0
	// if line1 > 0 {
	// 	char1 = (endPositions[line1-1] + 1) - start
	// }
	// if line2 > 0 {
	// 	char2 = (endPositions[line2-1] + 1) - lex.Pos
	// }
	value := errors.Lexical(lex.SourceCode, msg, start, lex.Pos).Error()
	return token.Token{Value: value, Typ: token.Error, Start: start, End: lex.Pos}
}
