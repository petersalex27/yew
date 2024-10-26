// =================================================================================================
// Alex Peters - January 22, 2024
//
// errors specific to lexical analysis step
// =================================================================================================
package lexer

import (
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/internal/errors"
)

const (
	UnexpectedEOF         string = "unexpected end of file"
	ExpectedCharLiteral   string = "expected character literal"
	IllegalCharLiteral    string = "illegal character literal"
	IllegalEscapeSequence string = "illegal escape sequence"
	IllegalStringLiteral  string = "illegal string literal"
	UnterminatedString    string = "unterminated string literal"
	UnexpectedUnderscore  string = "unexpected underscore"
	UnexpectedSymbol      string = "unexpected symbol"
	ExpectedAnnotationId  string = "annotation must have an identifier"
	UnexpectedEndOfLine   string = "unexpected end of line"
	IllegalBinLiteral     string = "illegal binary literal"
	IllegalOctLiteral     string = "illegal octal literal"
	IllegalHexLiteral     string = "illegal hexadecimal literal"
	IllegalIntLiteral     string = "illegal integer literal"
	CannotProcessSymbol   string = "cannot process symbol"
)

func (lex *Lexer) error(msg string) token.Token {
	var value string
	var start, end int
	if lex.Pos <= 0 {
		start, end = 0, max(0, lex.Pos)
		value = errors.Lexical(lex.SourceCode, msg, start, end).Error()
	} else {
		start, end = lex.Pos-1, lex.Pos
		value = errors.Lexical(lex.SourceCode, msg, start, end).Error()
	}
	return token.Token{Value: value, Typ: token.Error, Start: start, End: end}
}
