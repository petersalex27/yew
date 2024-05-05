package token

import "fmt"

type Type uint

func (tokenType Type) Make() Token {
	return Token{Type: tokenType, Value: tokenStringMap[tokenType]}
}

func (tokenType Type) MakeValued(val string) Token {
	return Token{Value: val, Type: tokenType}
}

func IsKeyword(tokenType Type) bool {
	return tokenType > _keywords_start_ && tokenType < _keywords_end_
}

const (
	IntValue Type = iota
	CharValue
	FloatValue
	StringValue

	Id
	ImplicitId
	Affixed
	Hole

	// keywords start
	_keywords_start_ // ===================================
	Alias
	Derives
	Import
	In
	Let
	Module
	Use
	Trait
	Where
	As
	With
	Of
	Extension
	Case
	Public
	Open
	Automatic
	Mutual

	LeftParen
	RightParen
	LeftBracket
	RightBracket
	LeftBrace
	RightBrace
	Comma
	Dot
	DotDot
	Colon
	ThickArrow
	Arrow
	Bar
	Equal
	Backslash
	// keywords end
	_keywords_end_ // ==============================

	Comment

	Underscore

	Percent

	Indent

	Newline

	ExtensionKey

	EndOfTokens
)

func (ty Type) String() string {
	switch ty {
	case IntValue:
		return "IntValue"
	case CharValue:
		return "CharValue"
	case FloatValue:
		return "FloatValue"
	case StringValue:
		return "StringValue"
	case Id:
		return "Id"
	case ImplicitId:
		return "ImplicitId"
	case Affixed:
		return "Affixed"
	case Hole:
		return "Hole"
	case Alias:
		return "Alias"
	case Derives:
		return "Derives"
	case With:
		return "With"
	case Import:
		return "Import"
	case In:
		return "In"
	case Let:
		return "Let"
	case Module:
		return "Module"
	case Use:
		return "Use"
	case Trait:
		return "Trait"
	case Where:
		return "Where"
	case Of:
		return "Of"
	case As:
		return "As"
	case Extension:
		return "Extension"
	case Case:
		return "Case"
	case Public:
		return "Public"
	case Open:
		return "Open"
	case Automatic:
		return "Automatic"
	case Mutual:
		return "Mutual"
	case LeftParen:
		return "LeftParen"
	case RightParen:
		return "RightParen"
	case LeftBracket:
		return "LeftBracket"
	case RightBracket:
		return "RightBracket"
	case LeftBrace:
		return "LeftBrace"
	case RightBrace:
		return "RightBrace"
	case Comma:
		return "Comma"
	case Dot:
		return "Dot"
	case DotDot:
		return "DotDot"
	case Colon:
		return "Colon"
	case ThickArrow:
		return "ThickArrow"
	case Arrow:
		return "Arrow"
	case Bar:
		return "Bar"
	case Equal:
		return "Equal"
	case Backslash:
		return "Backslash"
	case Comment:
		return "Comment"
	case Underscore:
		return "Underscore"
	case Percent:
		return "Percent"
	case Indent:
		return "Indent"
	case Newline:
		return "Newline"
	case ExtensionKey:
		return "ExtensionKey"
	case EndOfTokens:
		return "EndOfTokens"
	default:
		return fmt.Sprintf("Type(%d)", ty)
	}
}

var tokenStringMap = map[Type]string{
	IntValue:    "",
	CharValue:   "",
	FloatValue:  "",
	StringValue: "",

	Alias:     "alias",
	Derives:   "derives",
	With:      "with",
	Import:    "import",
	In:        "in",
	Let:       "let",
	Module:    "module",
	Use:       "use",
	Trait:     "trait",
	Where:     "where",
	As:        "as",
	Of:        "of",
	Extension: "extension",
	Case:      "case",
	Public:    "public",
	Open:      "open",
	Automatic: "automatic",
	Mutual:    "mutual",

	LeftParen:    "(",
	RightParen:   ")",
	LeftBracket:  "[",
	RightBracket: "]",
	LeftBrace:    "{",
	RightBrace:   "}",
	Comma:        ",",
	Dot:          ".",
	DotDot:       "..",
	Colon:        ":",
	ThickArrow:   "=>",
	Arrow:        "->",
	Bar:          "|",
	Equal:        "=",
	Backslash:    "\\",

	Underscore: "_",

	Percent: "%",

	Newline: "\n",

	ExtensionKey: "",

	Hole: "",
}
