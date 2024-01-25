// =================================================================================================
// Alex Peters - January 24, 2024
// =================================================================================================

package parser

import "github.com/petersalex27/yew/errors"

const (
	ExpectedIdentifier string = "expected identifier"
	ExpectedModule string = "expected 'module'"
	ExpectedWhere string = "expected 'where'"
	UnexpectedToken string = "unexpected token"
)

// creates a syntax error from the arguments
func makeSyntaxError(msg string, path string, line, start, end int) errors.ErrorMessage {
	e := errors.MakeError("Syntax", msg, line, start, end)
	e.SourceName = path
	return e
}

// adds an error constructed using parser's data and the message string passed as an argument
func (p *Parser) error(msg string) {
	start, end, line := p.Next.Start, p.Next.End, p.Next.Line
	e := makeSyntaxError(msg, p.Path, line, start, end)
	p.addMessage(e)
}