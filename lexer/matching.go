// =================================================================================================
// Alex Peters - January 22, 2024
//
// Regular expression and other matching related functions
// =================================================================================================

package lexer

import (
	"regexp"

	"github.com/petersalex27/yew/token"
)

// Matches regular expression `r` from start of string `s` to the last matching character.
//
// If no match is found, the match is an empty string, or the match is not found at start of `s`, an
// empty string is returned; otherwise, the matched string is returned
func matchRegex(r *regexp.Regexp, s string) string {
	loc := r.FindStringIndex(s)
	if loc == nil || loc[0] != 0 {
		return ""
	}
	return s[:loc[1]]
}

var keywords = map[string]token.Type{
	"alias":   token.Alias,
	"derives": token.Derives,
	"end":     token.End,
	"import":  token.Import,
	"in":      token.In,
	"let":     token.Let,
	"module":  token.Module,
	"use":     token.Use,
	"trait":   token.Trait,
	"where":   token.Where,
	"(":       token.LeftParen,
	")":       token.RightParen,
	"[":       token.LeftBracket,
	"]":       token.RightBracket,
	"{":       token.LeftBrace,
	"}":       token.RightBrace,
	",":       token.Comma,
	".":       token.Dot,
	"..":      token.DotDot,
	":":       token.Colon,
	"=>":      token.ThickArrow,
	"->":      token.Arrow,
	"|":       token.Bar,
	"=":       token.Equal,
	"\\":      token.Backslash,
}

// tries to match string `s` to a keyword token
func matchKeyword(s string, otherwise token.Type) token.Type {
	if ty, found := keywords[s]; found {
		return ty
	}
	return otherwise
}
