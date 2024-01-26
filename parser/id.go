// =================================================================================================
// Alex Peters - January 25, 2024
//
// =================================================================================================
package parser

import "github.com/petersalex27/yew/token"

func (parser *Parser) parseId() (id Ident, ok bool) {
	id.Token, ok = parser.idToken()
	return
}

func (parser *Parser) parseAffixedId() (affixed Ident, ok bool) {
	affixed.Token, ok = parser.getToken(token.Affixed, ExpectedName)
	return
}

// like
//
//	(*Parser) parseId
//
// but allows affixed names too
func (parser *Parser) parseFunctionName() (name Ident, ok bool) {
	if parser.Next.Type == token.Id {
		return parser.parseId()
	} else if parser.Next.Type == token.Affixed {
		return parser.parseAffixedId()
	}

	if !parser.optionalFlag {
		parser.error(ExpectedName)
	}
	ok = false
	return
}
