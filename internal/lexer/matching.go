// =================================================================================================
// Alex Peters - January 22, 2024
//
// Regular expression and other matching related functions
// =================================================================================================

package lexer

import (
	"regexp"

	"github.com/petersalex27/yew/internal/token"
)

const (
	idRegexClassRaw         = `[a-z][a-zA-Z0-9']`
	typeIdRegexClassRaw     = `[A-Z][a-zA-Z0-9']`
	implicitRaw             = `[a-z][0-9]*'*`
	identifierRegexRaw      = `(([a-z][a-zA-Z0-9']*)?_?((([a-z][a-zA-Z0-9']*_)+([a-z][a-zA-Z0-9']*)?)|([a-z][a-zA-Z0-9']*)))|([A-Z][a-zA-Z0-9']*)`
	symbolRegexClassRaw     = `[\(\)\[\]\{\}!@#\$%\^\&\*~,<>\.\?/;:\|\-\+=\\` + "`]"
	freeSymbolRegexClassRaw = `[!@#\$\^\&\*~,<>\?/:\|\-\+=\\` + "`]"
)

var (
	capIdRegex      = regexp.MustCompile(typeIdRegexClassRaw + `*`)
	implicitIdRegex = regexp.MustCompile(implicitRaw)
	// the following regular expression occurring between or in front of zero or one '_':
	//
	//	`\([!@#\$%\^\&\*~,<>\.\?/:\|-\+=`]+\)`
	//
	// examples:
	//
	//	`_>>=_`
	//	`_!`
	//	`!`
	symbolRegex          = regexp.MustCompile(affixedRegexGen(freeSymbolRegexClassRaw + `+`))
	identifierRegex      = regexp.MustCompile(identifierRegexRaw)
	endMultiCommentRegex = regexp.MustCompile(`\*-`)
	intRegex             = regexp.MustCompile(`[0-9](_*[0-9]+)*`)
	hexRegex             = regexp.MustCompile(`(0x|0X)[0-9a-fA-F](_*[0-9a-fA-F]+)*`)
	octRegex             = regexp.MustCompile(`(0o|0O)[0-7](_*[0-7]+)*`)
	binRegex             = regexp.MustCompile(`(0b|0B)(0|1)(_*(0|1)+)*`)
	justZeros            = regexp.MustCompile(`0+(_0+)*`)
	// the following regular expression, or, alpha-numeric id regex occurring between or in-front-of/after zero or one '_':
	//
	//	`([a-z][a-zA-Z0-9']*)?_?((([a-z][a-zA-Z0-9']*_)+([a-z][a-zA-Z0-9']*)?)|([a-z][a-zA-Z0-9']*))`
	//
	// examples:
	//
	//	`not_`
	//	`_mod_`
	//	`zipWith`
	affixedIdRegex   = regexp.MustCompile(affixedRegexGen(idRegexClassRaw + `*`))
	camelCaseIdRegex = regexp.MustCompile(idRegexClassRaw + `*`)
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

func camelCase(s string) bool {
	sLen := len(s)
	return sLen != 0 && len(matchRegex(camelCaseIdRegex, s)) == sLen
}

var keywords = map[string]token.Type{
	"alias":     token.Alias,
	"derives":   token.Derives,
	"with":      token.With,
	"import":    token.Import,
	"in":        token.In,
	"let":       token.Let,
	"module":    token.Module,
	"use":       token.Use,
	"spec":      token.Spec,
	"where":     token.Where,
	"of":        token.Of,
	"as":        token.As,
	"syntax":    token.Syntax,
	"case":      token.Case,
	"public":    token.Public,
	"open":      token.Open,
	"automatic": token.Automatic,
	"mutual":    token.Mutual,
	"erase":     token.Erase,
	"once":      token.Once,
	"@":         token.At,
	"(":         token.LeftParen,
	")":         token.RightParen,
	"[":         token.LeftBracket,
	"]":         token.RightBracket,
	"{":         token.LeftBrace,
	"}":         token.RightBrace,
	",":         token.Comma,
	".":         token.Dot,
	"..":        token.DotDot,
	":":         token.Colon,
	":=":        token.ColonEqual,
	"=>":        token.ThickArrow,
	"->":        token.Arrow,
	"|":         token.Bar,
	"=":         token.Equal,
	"\\":        token.Backslash,
}

func (lex *Lexer) isKeyword(s string) (token.Type, bool) {
	ty, found := keywords[s]
	return ty, found
}

func isImplicitId(s string) bool {
	// true iff entire string `s` matches regex
	return (len(s) != 0) && (len(matchRegex(implicitIdRegex, s)) == len(s))
}

// tries to match string `s` to a keyword token
func matchKeyword(s string, otherwise token.Type) token.Type {
	if ty, found := keywords[s]; found {
		return ty
	}
	return otherwise
}
