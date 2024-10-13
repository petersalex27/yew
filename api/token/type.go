package token

import (
	"fmt"

	"github.com/petersalex27/yew/api"
)

type Type uint8

func (t Type) Int() int64 { return int64(t) }

func (ty Type) Make() Token {
	return Token{Typ: ty, Value: tokenStringMap[ty]}
}

func (ty Type) MakeCommand() Token {
	return Token{Typ: ty, Value: commandStringMap[ty]}
}

func (ty Type) MakeValued(val string) Token {
	return Token{Value: val, Typ: ty}
}

func IsKeyword(tokenType Type) bool {
	return tokenType > _variable_end_ && tokenType < _keywords_end_
}

const (
	Error Type = iota
	IntValue
	CharValue
	FloatValue
	StringValue
	RawStringValue

	ImportPath

	Id
	Infix
	Hole
	_variable_end_
	MethodSymbol = _variable_end_
)

const (
	Alias Type = iota + 1 + _variable_end_
	Deriving
	Import
	In
	Let
	Module
	Using
	Spec
	Where
	As
	With
	Of
	Syntax
	Case
	Public
	Open
	Auto
	Inst
	Erase
	Once
	Impossible
	Requiring
	From
	Forall
	Ref
	Term
	Pattern

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
	ColonEqual
	FlatAnnotation
	LeftBracketAt
	EmptyParenEnclosure
	_keywords_end_
	EmptyBracketEnclosure = _keywords_end_
)

const (
	Comment Type = iota + 1 + _keywords_end_
	Underscore
	Newline
	_token_types_proper_end_
	EndOfTokens = _token_types_proper_end_
)

const (
	/* = Symbol related commands ======================================================== */
	Import_c Type = iota + 1 + _token_types_proper_end_

	/* = Info query commands ============================================================ */
	Instances_c
	Help_c
	Type_c
	Kind_c

	/* = Control commands =============================================================== */
	Main_c
	Run_c
	Quit_c
	Set_c
	Api_c
	Expose_c

	// these allow one to effectively create isolated environments in the repl
	/* = State commands ================================================================= */
	Save_c
	Restore_c

	/* = Source commands ================================================================ */
	Begin_c
	End_c
	Include_c

	_commands_end_
	USER_DEFINED_START = _commands_end_
)

const UnknownTokenTypeString string = "?UnknownTokenType"

func ProperTypeString(ty Type) string {
	switch ty {
	case Error:
		return "Error"
	case IntValue:
		return "IntValue"
	case CharValue:
		return "CharValue"
	case FloatValue:
		return "FloatValue"
	case StringValue:
		return "StringValue"
	case RawStringValue:
		return "RawStringValue"
	case ImportPath:
		return "ImportPath"
	case Id:
		return "Id"
	case Infix:
		return "Infix"
	case Hole:
		return "Hole"
	case MethodSymbol:
		return "MethodSymbol"
	case Alias:
		return "Alias"
	case Deriving:
		return "Deriving"
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
	case Using:
		return "Using"
	case Spec:
		return "Spec"
	case Where:
		return "Where"
	case Of:
		return "Of"
	case As:
		return "As"
	case Impossible:
		return "Impossible"
	case Requiring:
		return "Requiring"
	case From:
		return "From"
	case Forall:
		return "Forall"
	case Ref:
		return "Ref"
	case Term:
		return "Term"
	case Pattern:
		return "Pattern"
	case FlatAnnotation:
		return "FlatAnnotation"
	case LeftBracketAt:
		return "LeftBracketAt"
	case EmptyParenEnclosure:
		return "EmptyParenEnclosure"
	case EmptyBracketEnclosure:
		return "EmptyBracketEnclosure"
	case Syntax:
		return "Syntax"
	case Case:
		return "Case"
	case Public:
		return "Public"
	case Open:
		return "Open"
	case Auto:
		return "Auto"
	case Inst:
		return "Inst"
	case Erase:
		return "Erase"
	case Once:
		return "Once"
	case ColonEqual:
		return "ColonEqual"
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
	case Newline:
		return "Newline"
	case EndOfTokens:
		return "EndOfTokens"
	default:
		return "?UnknownTokenType"
	}
}

func (ty Type) String() string {
	s := ProperTypeString(ty)
	if InReplMode() && s == UnknownTokenTypeString {
		s = ty.CommandString()
	}
	return s
}

func (ty Type) CommandString() string {
	switch ty {
	case Import_c:
		return "Command(Import)"
	case Type_c:
		return "Command(Type)"
	case Instances_c:
		return "Command(Instances)"
	case Main_c:
		return "Command(Main)"
	case Expose_c:
		return "Command(Expose)"
	case Help_c:
		return "Command(Help)"
	case Quit_c:
		return "Command(Quit)"
	case Run_c:
		return "Command(Run)"
	case Set_c:
		return "Command(Set)"
	default:
		return fmt.Sprintf("Command(%d)", ty)
	}
}

func (ty Type) Match(tok api.Node) bool {
	return ty.String() == tok.Type().String()
}

var tokenStringMap = map[Type]string{
	Error:       ``, //`<error>`,
	IntValue:    ``, //`<int-literal>`,
	CharValue:   ``, //`<char-literal>`,
	FloatValue:  ``, //`<float-literal>`,
	StringValue: ``, //`<string-literal>`,
	RawStringValue: ``, //`<raw-string-literal>`,
	ImportPath: ``, //`<import-path>`,

	Id:    ``, //`<identifier>`,
	Infix: ``, //`<(infix)>`,
	Hole:  ``, //`<?hole>`,

	MethodSymbol: ``, //`<method-symbol>`,

	Alias:      `alias`,
	Deriving:   `deriving`,
	With:       `with`,
	Import:     `import`,
	In:         `in`,
	Let:        `let`,
	Module:     `module`,
	Using:      `use`,
	Spec:       `spec`,
	Where:      `where`,
	As:         `as`,
	Of:         `of`,
	Syntax:     `syntax`,
	Case:       `case`,
	Public:     `public`,
	Open:       `open`,
	Auto:       `auto`,
	Inst:       `inst`,
	Erase:      `erase`,
	Once:       `once`,
	Impossible: `impossible`,
	Requiring:  `requiring`,
	From:       `from`,
	Forall:     `forall`,
	Ref:        `ref`,
	Term:       `term`,
	Pattern:    `pattern`,

	LeftParen:    `(`,
	RightParen:   `)`,
	LeftBracket:  `[`,
	RightBracket: `]`,
	LeftBrace:    `{`,
	RightBrace:   `}`,
	Comma:        `,`,
	Dot:          `.`,
	DotDot:       `..`,
	Colon:        `:`,
	ThickArrow:   `=>`,
	Arrow:        `->`,
	Bar:          `|`,
	Equal:        `=`,
	Backslash:    `\`,
	ColonEqual:   `:=`,

	Underscore: `_`,

	FlatAnnotation: ``, //`<annotation>`,

	LeftBracketAt:         `[@`,
	EmptyParenEnclosure:   `()`,
	EmptyBracketEnclosure: `[]`,

	Newline: "\n",

	Comment: ``, //`<comment>`,

	EndOfTokens: ``, //`<end-of-tokens>`,
}

var commandStringMap = map[Type]string{
	Import_c:    ":import",
	Type_c:      ":type",
	Instances_c: ":instances",
	Main_c:      ":main",
	Expose_c:    ":expose",
	Help_c:      ":help",
	Quit_c:      ":quit",
	Run_c:       ":run",
	Set_c:       ":set",
}

// return the standard form of a command literal string for a given command type
//
// if one does not exist, an empty string is returned
func (command Type) CommandLiteral() string {
	return commandStringMap[command]
}
