// =================================================================================================
// Alex Peters - March 04, 2024
// =================================================================================================
package parser

import (
	"github.com/petersalex27/yew/token"
)

func (parser *Parser) get(ty token.Type) (tok token.Token, ok bool) {
	typ := parser.Peek().Type
	if ok = ty == typ; !ok {
		return
	}
	tok = parser.Advance()
	return
}

func (parser *Parser) getPascalCaseIdent() (ident termElem, ok bool, upperFailed bool) {
	tok, ok := parser.get(token.Id)
	if ok {
		upperFailed = !startsWithUppercase(tok.Value)
	}
	if ok = ok && !upperFailed; !ok {
		return
	}

	ident = termElem{
		Term: Ident{
			Name:  tok.Value,
			Start: tok.Start,
			End:   tok.End,
		},
		termInfo: termInfo{},
	}
	return
}