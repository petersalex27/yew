// =================================================================================================
// Alex Peters - January 22, 2024
//
// Regular expression and other matching related functions
// =================================================================================================

package lexer

import (
	"github.com/petersalex27/yew/api/token"
)

var keywords = map[string]token.Type{
	"alias":      token.Alias,
	"as":         token.As,
	"auto":       token.Auto,
	"case":       token.Case,
	"deriving":   token.Deriving,
	"erase":      token.Erase,
	"forall":     token.Forall,
	"from":       token.From,
	"import":     token.Import,
	"impossible": token.Impossible,
	"in":         token.In,
	"inst":       token.Inst,
	"let":        token.Let,
	"module":     token.Module,
	"of":         token.Of,
	"once":       token.Once,
	"open":       token.Open,
	"pattern":    token.Pattern,
	"public":     token.Public,
	"ref":        token.Ref,
	"requiring":  token.Requiring,
	"spec":       token.Spec,
	"syntax":     token.Syntax,
	"term":       token.Term,
	"using":      token.Using,
	"where":      token.Where,
	"with":       token.With,
	"()":         token.EmptyParenEnclosure,
	"[]":         token.EmptyBracketEnclosure,
	"(":          token.LeftParen,
	")":          token.RightParen,
	"[":          token.LeftBracket,
	"]":          token.RightBracket,
	"{":          token.LeftBrace,
	"}":          token.RightBrace,
	",":          token.Comma,
	".":          token.Dot,
	"..":         token.DotDot,
	":":          token.Colon,
	":=":         token.ColonEqual,
	"=>":         token.ThickArrow,
	"->":         token.Arrow,
	"|":          token.Bar,
	"=":          token.Equal,
	"\\":         token.Backslash,
	"[@":         token.LeftBracketAt,
}

// all valid commands and their short forms mapped to their corresponding token type
var commands = map[string]token.Type{
	":import":    token.Import_c,
	":i":         token.Import_c,
	":type":      token.Type_c,
	":t":         token.Type_c,
	":instances": token.Instances_c,
	":in":        token.Instances_c,
	":main":      token.Main_c,
	":m":         token.Main_c,
	":expose":    token.Expose_c,
	":help":      token.Help_c,
	":h":         token.Help_c,
	":quit":      token.Quit_c,
	":q":         token.Quit_c,
	":run":       token.Run_c,
	":r":         token.Run_c,
	":set":       token.Set_c,
}

func (lex *Lexer) isKeyword(s string) (token.Type, bool) {
	ty, found := keywords[s]
	return ty, found
}

// tries to match string `s` to a keyword token
func matchKeyword(s string, otherwise token.Type) token.Type {
	if ty, found := keywords[s]; found {
		return ty
	}
	return otherwise
}
