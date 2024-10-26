package common

import (
	"regexp"
)

var (
	// alphanumeric ids (with optional, non-initial position `'`s)
	alphanumericIdRegex = regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9']*`)
	// alphanumeric ids that start with a lowercase letter (with optional, non-initial position `'`s)
	camelCaseIdRegex = regexp.MustCompile(`[a-z][a-zA-Z0-9']*`)
	// alphanumeric ids that start with an uppercase letter (with optional, non-initial position `'`s)
	pascalCaseIdRegex = regexp.MustCompile(`[A-Z][a-zA-Z0-9']*`)
	// standalone symbols
	standaloneRegex = regexp.MustCompile(`[(){}\[\],]`)
	// ids that use only and one or more non-alphanumeric characters
	symbolIdRegex = regexp.MustCompile(`[!@#$%^&*\-=+;:\\|~,<.>/?]+`)
	// any kind of id or symbol, excluding infixed ids and symbols
	nonInfixNameRegex = regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9']*|[!@#$%^&*\-=+;:\\|~,<.>/?]+`)
	// infixed ids and symbols
	infixIdRegex = regexp.MustCompile(`\([a-zA-Z][a-zA-Z0-9']*\)|\([!@#$%^&*\-=+;:\\|~,<.>/?]+\)`)
	// method ids
	methodIdRegex = regexp.MustCompile(`\([.][a-zA-Z][a-zA-Z0-9']*\)|\([.][!@#$%^&*\-=+;:\\|~,<.>/?]+\)`)
	// any kind of id or symbol, including infixed ids and symbols
	nameRegex = regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9']*|\([.]?[a-zA-Z][a-zA-Z0-9']*\)|[!@#$%^&*\-=+;:\\|~,<.>/?]+|\([.]?[!@#$%^&*\-=+;:\\|~,<.>/?]+\)`)
	// import path regex
	importPathRegex = regexp.MustCompile(`"[a-z][a-zA-Z0-9']*(/[a-z][a-zA-Z0-9']*)*"`)
	// integer literal regex
	intRegex = regexp.MustCompile(`[0-9]+(_[0-9]+)*`)
	// float literal regex
	floatRegex = regexp.MustCompile(`[0-9]+(_[0-9]+)*([.][0-9](_*[0-9]+)*([eE][+-]?[0-9]+(_[0-9]+)*)?|[eE][+-]?[0-9]+(_[0-9]+)*)`)
	// hex literal regex
	hexRegex = regexp.MustCompile(`0[xX][0-9a-fA-F](_[0-9a-fA-F]+)*`)
	// oct literal regex
	octRegex = regexp.MustCompile(`0[oO][0-7](_[0-7]+)*`)
	// bin literal regex
	binRegex = regexp.MustCompile(`0[bB][01](_[01]+)*`)

	// simple string regex
	simpleStringRegex = regexp.MustCompile(`"[^\n\r]*"`)
	// version regex // TODO: update to include pre-release suffixes, @latest tag, and @stable tag
	versionRegex = regexp.MustCompile(`v([0-9]+([.][0-9]+){2}(@[a-z]+[.][0-9]+)?|@latest|@stable)`)

	// string regex
	stringRegex = regexp.MustCompile(`"(\\[abfnrtv\\"]|[^\n\r\\])*"`)
	// unterminated string regex
	unterminatedStringRegex = regexp.MustCompile(`"(\\[abfnrtv\\"]|[^\n\r\\])*[\r\n]`)
	// bad escape string
	badEscapeStringRegex = regexp.MustCompile(`"(\\[abfnrtv\\"]|[^\n\r\\])*\\[^abfnrtv\\]`)
	// char regex
	charRegex = regexp.MustCompile(`'(\\[abfnrtv\\']|[^\n\r])'`)
	// hole regex
	holeRegex = regexp.MustCompile(`\?[A-Za-z][A-Za-z0-9']*`)
)

type stringer interface {
	String() string
}

func Is_camelCase(s string) bool {
	return len(s) != 0 && IsCompleteMatch(camelCaseIdRegex, s)
}

func Is_camelCase2[r stringer](s r) bool {
	return Is_camelCase(s.String())
}

func Is_PascalCase(s string) bool {
	return len(s) != 0 && IsCompleteMatch(pascalCaseIdRegex, s)
}

func Is_PascalCase2[r stringer](s r) bool {
	return Is_PascalCase(s.String())
}

func Is_symbolCase(s string) bool {
	return len(s) != 0 && IsCompleteMatch(symbolIdRegex, s)
}

func Is_symbolCase2[r stringer](s r) bool {
	return Is_symbolCase(s.String())
}

func IsConstructorId(s string) bool {
	// infix upper ident | upper ident | infix symbol | symbol
	lnS := len(s)
	ps := InfixId.Match(s)
	if ps != nil && len(*ps) == lnS {
		return true
	}

	ps = PascalCaseId.Match(s)
	if ps != nil && len(*ps) == lnS {
		return true
	}

	ps = SymbolId.Match(s)
	return ps != nil && len(*ps) == lnS
}

func IsConstructorId2[r stringer](s r) bool {
	return IsConstructorId(s.String())
}

func IsAlphanumericId(s string) bool {
	return len(s) != 0 && IsCompleteMatch(alphanumericIdRegex, s)
}

func IsAlphanumericId2[r stringer](s r) bool {
	return IsAlphanumericId(s.String())
}

func IsCompleteMatch(re *regexp.Regexp, s string) bool {
	matched := MatchRegex(re, s)
	return matched != nil && len(*matched) == len(s)
}

func IsCompleteMatch2[r stringer](re *regexp.Regexp, s r) bool {
	return IsCompleteMatch(re, s.String())
}

func IsImportPathString(s string) bool {
	return len(s) != 0 && IsCompleteMatch(importPathRegex, s)
}

func IsStandaloneSymbol(s string) bool {
	return len(s) != 0 && IsCompleteMatch(standaloneRegex, s)
}

func IsStandaloneSymbol2[r stringer](s r) bool {
	return IsStandaloneSymbol(s.String())
}

// Matches regular expression `r` from start of string `s` to the last matching character.
//
// If no match is found or the match is not found at start of `s`, nil is returned; otherwise, the
// matched string is returned.
//
// NOTE: if an empty string is matched by `re` at the start of `s`, an empty string is returned,
// distinguishing it from the case where no match is found
func MatchRegex(re *regexp.Regexp, s string) *string {
	loc := re.FindStringIndex(s)
	if loc == nil || loc[0] != 0 {
		return nil
	}
	res := s[:loc[1]]
	return &res
}

func MatchRegex2[r stringer](re *regexp.Regexp, s r) *string {
	return MatchRegex(re, s.String())
}

type regexIdentifier byte

const (
	// an alphanumeric id
	AlphanumericId regexIdentifier = iota
	// an alphanumeric id that starts with a lowercase letter
	CamelCaseId
	// an alphanumeric id that starts with an uppercase letter
	PascalCaseId
	// a non-alphanumeric id that is not an infix or method id
	SymbolId
	// a non-infix name
	NonInfixName
	// a symbol or alphanumeric id enclosed in parentheses, e.g., (>>=)
	InfixId
	// a symbol or alphanumeric id prefixed with a '.' and enclosed in parentheses, e.g., (.run)
	MethodId
	// an alphanumeric id prefixed with a '?', e.g., ?result
	HoleId
	// name regex: alphanumeric id, symbol id, infix id, or method id
	Name
	// one of the following: ( ) { } [ ] ,
	StandaloneSymbol
	// integer literal (base 10) with single '_' separators
	IntLiteral
	// float literal: optional scientific notation and single '_' separators
	FloatLiteral
	// hex integer literal with prefix `0x` or `0X` and single '_' separators
	HexLiteral
	// octal integer literal with prefix `0o` or `0O` and single '_' separators
	OctLiteral
	// binary integer literal with prefix `0b` or `0B` and single '_' separators
	BinLiteral
	// single quoted character literal w/ interpretation (i.e. escape sequences)
	CharLiteral
	// double quoted string literal w/ interpretation (i.e. escape sequences)
	InterpretedStringLiteral
	ILLEGAL_UnterminatedInterpretedStringLiteral
	ILLEGAL_BadEscapeSequence
	// double quoted string literal w/o interpretation (i.e. escape sequences)
	SimpleString
	// a version literal with optional milestone suffixes or relative version tags
	//
	// relative example:
	//		@latest
	// milestone example:
	//		v1.2.3@alpha.1
	// absolute example:
	//		v1.2.3
	Version
)

var regexMap = map[regexIdentifier]*regexp.Regexp{
	AlphanumericId:           alphanumericIdRegex,
	CamelCaseId:              camelCaseIdRegex,
	PascalCaseId:             pascalCaseIdRegex,
	SymbolId:                 symbolIdRegex,
	NonInfixName:             nonInfixNameRegex,
	Name:                     nameRegex,
	InfixId:                  infixIdRegex,
	MethodId:                 methodIdRegex,
	HoleId:                   holeRegex,
	StandaloneSymbol:         standaloneRegex,
	IntLiteral:               intRegex,
	FloatLiteral:             floatRegex,
	HexLiteral:               hexRegex,
	OctLiteral:               octRegex,
	BinLiteral:               binRegex,
	CharLiteral:              charRegex,
	InterpretedStringLiteral: stringRegex,
	ILLEGAL_UnterminatedInterpretedStringLiteral: unterminatedStringRegex,
	ILLEGAL_BadEscapeSequence:                    badEscapeStringRegex,
	SimpleString:                                 simpleStringRegex,
	Version:                                      versionRegex,
}

func (id regexIdentifier) Match(s string) *string {
	if re, ok := regexMap[id]; ok {
		return MatchRegex(re, s)
	}
	return nil
}
