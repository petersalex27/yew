package parser

import "github.com/petersalex27/yew/internal/token"

func (parser *Parser) parseAnnotation() (ok bool) {
	annot := parser.Advance()
	switch annot.Value {
	case "builtin":
	case "error":
	case "warn":
	case "deprecated":
	case "todo":
	case "external":
	case "inline":
	case "noInline":
	case "specialize":
	case "noAlias":
	case "pure":
	case "noGc":
	case "infixl", "infixr":
	case "infix":
	}
	// if parser.Peek().Type != token.LeftParen {
	// 	// TODO
	// 	return ok
	// }
	panic("TODO: implement")
}

func (parser *Parser) parseEnvAnnotation(tok token.Token) (ok bool) {
	panic("TODO: implement")
}

func (parser *Parser) readEnvironment() (ok bool) {
	ok = true

	tok := parser.Peek()

	for ok && tok.Type != token.At {
		_ = parser.Advance()

		ok = parser.parseEnvAnnotation(tok)
		if ok {
			tok = parser.Peek()
		}
	}

	return ok
}