// =================================================================================================
// Alex Peters - January 24, 2024
// =================================================================================================

package parser

import (
	"github.com/petersalex27/yew/common"
	"github.com/petersalex27/yew/errors"
	"github.com/petersalex27/yew/token"
)

const (
	// = "expected" ==================================================================================
	ExpectedBinding    string = "expected lambda binding '_->_'"
	ExpectedEqual      string = "expected assignment"
	ExpectedIdentifier string = "expected identifier"
	ExpectedModule     string = "expected 'module'"
	ExpectedName       string = "expected name"
	ExpectedRParen     string = "expected ')'"
	ExpectedType       string = "expected type"
	ExpectedWhere      string = "expected 'where'"

	// = "illegal" ===================================================================================
	IllegalWhere string = "illegal 'where'"

	// TODO: better message??
	UnusedContext string = "unreferenced context"

	// = "unexpected" ================================================================================
	UnexpectedEOF   string = "unexpected end of file"
	UnexpectedToken string = "unexpected token"
)

// (token, "expected" error message) pairs
var expectMap = map[token.Type]string{
	token.Arrow:      ExpectedBinding,
	token.Equal:      ExpectedEqual,
	token.Id:         ExpectedIdentifier,
	token.Module:     ExpectedModule,
	token.RightParen: ExpectedRParen,
	token.Where:      ExpectedWhere,
	token.Colon:      ExpectedType,
}

func (parser *Parser) calcLocation(start, end int) (line1, line2, char1, char2 int) {
	line1 = 1 + common.SearchRange(parser.PositionRanges, start) // 1 + result = 0 or greater
	line2 = 1 + common.SearchRange(parser.PositionRanges, end)   // 1 + result = 0 or greater
	if line1 > 0 {
		char1 = (parser.PositionRanges[line1-1] + 1) - start
	}
	if line2 > 0 {
		char2 = (parser.PositionRanges[line2-1] + 1) - end
	}
	return
}

// returns appropriate "expected X" where X is the expected token with type `expect`.
//
// if there's no appropriate "expected X" message for the given type, then `UnexpectedToken` is
// returned
func getExpectMessage(expect token.Type) (errorMessage string) {
	if errorMessage, found := expectMap[expect]; found {
		return errorMessage
	}

	return UnexpectedToken
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
