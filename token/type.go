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
	return tokenType >= Alias && tokenType <= Backslash
}

const (
	IntValue Type = iota
	CharValue
	FloatValue
	StringValue

	Id
	Affixed
	CapId

	Alias
	Derives
	Import
	In
	Let
	Module
	Use
	Trait
	Where

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

	Comment

	Underscore

	At
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
	case Affixed:
		return "Affixed"
	case CapId:
		return "CapId"
	case Alias:
		return "Alias"
	case Derives:
		return "Derives"
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
	case At:
		return "At"
	default:
		return fmt.Sprintf("Type(%d)", ty)
	}
}

var tokenStringMap = map[Type]string{
	IntValue:    "",
	CharValue:   "",
	FloatValue:  "",
	StringValue: "",

	Alias:   "alias",
	Derives: "derives",
	Import:  "import",
	In:      "in",
	Let:     "let",
	Module:  "module",
	Use:     "use",
	Trait:   "trait",
	Where:   "where",

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

	At: "@",
}
