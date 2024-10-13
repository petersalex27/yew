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
	// any kind of id or symbol, including infixed ids and symbols
	nameRegex = regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9']*|\([a-zA-Z][a-zA-Z0-9']*\)|[!@#$%^&*\-=+;:\\|~,<.>/?]+|\([!@#$%^&*\-=+;:\\|~,<.>/?]+\)`)
	// ids that that look like paths using camelCase ids between and after '/' and '.' respectively
	importPathRegex = regexp.MustCompile(`[a-z][a-zA-Z0-9']*(/[a-z][a-zA-Z0-9']*)*`)
	// integer literal regex
	intRegex = regexp.MustCompile(`[0-9]+(_[0-9]+)*`)
	// float literal regex
	floatRegex = regexp.MustCompile(`[0-9]+(_[0-9]+)*([.][0-9](_*[0-9]+)*)?([+-]?[eE][0-9]+(_[0-9]+)*)?`)
	// hex literal regex
	hexRegex = regexp.MustCompile(`0[xX][0-9a-fA-F](_[0-9a-fA-F]+)*`)
	// oct literal regex
	octRegex = regexp.MustCompile(`0[oO][0-7](_[0-7]+)*`)
	// bin literal regex
	binRegex = regexp.MustCompile(`0[bB][01](_[01]+)*`)
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
	AlphanumericId regexIdentifier = iota
	CamelCaseId
	PascalCaseId
	SymbolId
	NonInfixName
	Name
	ImportPath
	StandaloneSymbol
	IntLiteral
	FloatLiteral
	HexLiteral
	OctLiteral
	BinLiteral
)

var regexMap = map[regexIdentifier]*regexp.Regexp{
	AlphanumericId:   alphanumericIdRegex,
	CamelCaseId:      camelCaseIdRegex,
	PascalCaseId:     pascalCaseIdRegex,
	SymbolId:         symbolIdRegex,
	NonInfixName:     nonInfixNameRegex,
	Name:             nameRegex,
	ImportPath:       importPathRegex,
	StandaloneSymbol: standaloneRegex,
	IntLiteral:       intRegex,
	FloatLiteral:     floatRegex,
	HexLiteral:       hexRegex,
	OctLiteral:       octRegex,
	BinLiteral:       binRegex,
}

func (id regexIdentifier) Match(s string) *string {
	if re, ok := regexMap[id]; ok {
		return MatchRegex(re, s)
	}
	return nil
}
