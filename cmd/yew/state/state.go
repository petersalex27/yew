package state

import (
	"github.com/petersalex27/yew/internal/lexer"
	"github.com/petersalex27/yew/internal/parser"
)

type Data struct {
	Parser *parser.ParserState
	Lexer  *lexer.Lexer
}
