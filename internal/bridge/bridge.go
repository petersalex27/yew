// =================================================================================================
// Alex Peters - January 29, 2024
//
// Bridge between lexer and parser
// =================================================================================================
package bridge

import (
	"github.com/petersalex27/yew/internal/errors"
	"github.com/petersalex27/yew/internal/lexer"
	"github.com/petersalex27/yew/internal/parser"
	"github.com/petersalex27/yew/internal/source"
)

// transfers the source code at the given path to a newly initialized lexer
func TransferSourceCode(path source.PathSpec) (*lexer.Lexer, []errors.ErrorMessage) {
	lex := lexer.Init(path)
	result := lex.Write()
	if result < 0 {
		return lex, lex.Messages()
	}
	return lex, nil
}

// creates and initializes a new parser, transferring lex's data to the parser
//
// returns new, initialized parser
func TransferLexerData(lex *lexer.Lexer) *parser.Parser {
	p := parser.Init(lex.SourceCode)
	p.Load(lex.Tokens)
	return p
}
