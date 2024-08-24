// =================================================================================================
// Alex Peters - March 04, 2024
// =================================================================================================
package parser

import (
	"github.com/petersalex27/yew/internal/token"
)

func (parser *Parser) get(ty token.Type) (tok token.Token, ok bool) {
	tok = parser.Peek()
	if ok = ty == tok.Type; !ok {
		return
	}
	tok = parser.Advance()
	return
}

func (parser *Parser) getPascalCaseIdent() (ident Ident, ok bool, upperFailed bool) {
	tok, ok := parser.get(token.Id)
	if ok {
		upperFailed = !startsWithUppercase(tok.Value)
	}
	if ok = ok && !upperFailed; !ok {
		return Ident{Start: tok.Start, End: tok.End}, false, upperFailed
	}

	ident = Ident{
		Name:  tok.Value,
		Start: tok.Start,
		End:   tok.End,
	}
	return
}