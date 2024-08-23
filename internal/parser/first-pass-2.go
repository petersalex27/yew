package parser

import "github.com/petersalex27/yew/internal/token"

func (parser *Parser) annotation() (ok bool) {
	panic("TODO: implement 'annotation'")
}

func (parser *Parser) parseVisibility(next token.Type, vis Visibility) (ok bool) {
	panic("TODO: implement 'parseVisibility'")
}

func (parser *Parser) mutual() (ok bool) {
	panic("TODO: implement 'mutual'")
}

func (parser *Parser) automatic() (ok bool) {
	panic("TODO: implement 'automatic'")
}

func (parser *Parser) alias() (ok bool) {
	panic("TODO: implement 'alias'")
}

func (parser *Parser) spec() (ok bool) {
	panic("TODO: implement 'spec'")
}

func (parser *Parser) declarationOrDefinition() (ok bool) {
	panic("TODO: implement 'declarationOrDefinition'")
}

func firstPass(parser *Parser) (ok bool) {
	// collect all the tokens into syntactically significant parts. The collected parts should have
	// indentations and the likes removed--basically, remove stuff that is only intended to denote
	// syntactic parts or rules (key-words/-symbols, indentation, annotations, etc)

	ok = true
	for end := false; !end && ok; {
		parser.inTop = true // set back to true each time

		parser.drop() // TODO: is this necessary?

		switch next := parser.Peek(); next.Type {
		case token.At:
			ok = parser.annotation()
		case token.Public:
			ok = parser.parseVisibility(next.Type, Public)
		case token.Open:
			ok = parser.parseVisibility(next.Type, Open)
		case token.Mutual:
			ok = parser.mutual()
		case token.Automatic:
			panic("TODO: implement 'automatic'")
		case token.Alias:
			panic("TODO: implement 'alias'")
		case token.Spec:
			ok = parser.spec()
		case token.Id, token.LeftParen:
			ok = parser.declarationOrDefinition()
		case token.EndOfTokens:
			end = true
		case token.Hole:
			ok, end = false, true
			parser.error(IllegalNonExprPosHole)
		default:
			parser.error(UnexpectedToken)
			ok, end = false, true
		}
	}
	return ok
}