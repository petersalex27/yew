// =================================================================================================
// Alex Peters - January 25, 2024
// =================================================================================================
package parser

import "github.com/petersalex27/yew/token"

func (parser *Parser) backslashToken() (token.Token, bool) {
	return parser.getToken(token.Backslash, UnexpectedToken)
}

// parse lambda binders
func (parser *Parser) parseBinders(lambda *Lambda) bool {
	endOptional := parser.StartOptional() // start
	lambda.Binders = []Ident{}
	v, ok := parser.parseId()
	for ; ok; v, ok = parser.parseId() {
		lambda.Binders = append(lambda.Binders, v)
	}
	endOptional() // end

	if len(lambda.Binders) < 1 {
		// record error
		end := parser.StopOptional() // force error now
		// ok == false b/c input will still be the same as when it was not ok before
		_, _ = parser.parseId()
		end() // stop forcing errors
		return false
	}
	return true
}

func (parser *Parser) parseLambda() (lambda Lambda, ok bool) {
	var backslash token.Token
	if backslash, ok = parser.backslashToken(); !ok {
		return
	}
	lambda.Start = backslash.Start

	if ok = parser.parseBinders(&lambda); !ok {
		return
	}

	if _, ok = parser.bindingArrowToken(); !ok {
		return
	}

	if lambda.Bound, ok = parser.parseExpression(); !ok {
		return
	}

	_, lambda.End = lambda.Bound.Pos()
	return
}
