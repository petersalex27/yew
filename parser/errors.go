// =================================================================================================
// Alex Peters - January 24, 2024
// =================================================================================================

package parser

import (
	"github.com/petersalex27/yew/common"
	"github.com/petersalex27/yew/errors"
)

const (
	ExpectedEqual string = "expected equality"
	ExpectedIdentifier string = "expected identifier"
	ExpectedModule string = "expected 'module'"
	ExpectedWhere string = "expected 'where'"
	// TODO: better message??
	NoContextAdded string = "no context added"
	UnexpectedToken string = "unexpected token"
)

func (parser *Parser) calcLocation(start, end int) (line1, line2, char1, char2 int) {
	line1 = 1 + common.SearchRange(parser.PositionRanges, start)   // 1 + result = 0 or greater
	line2 = 1 + common.SearchRange(parser.PositionRanges, end) // 1 + result = 0 or greater
	if line1 > 0 {
		char1 = (parser.PositionRanges[line1-1] + 1) - start
	}
	if line2 > 0 {
		char2 = (parser.PositionRanges[line2-1] + 1) - end
	}
	return
}

// creates a syntax error from the arguments
func makeSyntaxError(msg string, path string, line, lineEnd, start, end int) errors.ErrorMessage {
	e := errors.MakeError("Syntax", msg, line, lineEnd, start, end)
	e.SourceName = path
	return e
}

// creates a syntax error from the arguments
func makeSyntaxWarning(msg string, path string, line, lineEnd, start, end int) errors.ErrorMessage {
	e := errors.MakeWarning("Syntax", msg, line, lineEnd, start, end)
	e.SourceName = path
	return e
}

// adds a warning constructed using arguments
func (parser *Parser) warning2(msg string, startPos, endPos int) {
	line1, line2, char1, char2 := parser.calcLocation(startPos, endPos)
	w := makeSyntaxWarning(msg, parser.Path, line1, line2, char1, char2)
	parser.addMessage(w)
}

// adds an error constructed using parser's data and the message string passed as an argument
func (parser *Parser) error(msg string) {
	start, end := parser.Next.Start, parser.Next.End
	line1, line2, char1, char2 := parser.calcLocation(start, end)	
	e := makeSyntaxError(msg, parser.Path, line1, line2, char1, char2)
	parser.addMessage(e)
}