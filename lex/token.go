package scan

import (
	"yew/info"
	types "yew/type"
	util "yew/utils"
	"yew/value"

	err "yew/error"
)

type TokenType int

const (
	// variable
	ID TokenType = iota
	TYPE_ID
	VALUE

	// groupings
	LBRACK
	RBRACK
	LPAREN
	RPAREN
	LCURL
	RCURL

	// keywords
	STRING
	CHAR
	BOOL
	FLOAT
	INT
	TYPE
	CLASS
	WHERE
	LET
	CONST
	MUT
	LAZY

	PLUS
	PLUS_PLUS
	MINUS
	STAR
	SLASH

	HAT

	EQUALS

	COLON
	COLON_COLON
	SEMI_COLON
	QUESTION
	COMMA

	BANG
	BANG_EQUALS
	EQUALS_EQUALS
	BAR
	AMPER_AMPER
	BAR_BAR

	ARROW
	FAT_ARROW

	DOT
	DOT_DOT

	GREAT
	LESS
	GREAT_EQUALS
	LESS_EQUALS

	MOD

	BACKSLASH

	AT

	UNDERSCORE

	PACKAGE

	MODULE

	NEW_LINE

	COLON_COLON_EQUAL

	ERROR

	EOF

	BANG_POSTFIX__
	MINUS_PREFIX__
	PLUS_PREFIX__

	_ANY__
	_REPEAT__
	_START_GROUP__
	_END_GROUP__
	_ALTERNATION__
	END__
)

type Token interface {
	util.Stringable
	info.Locatable
	GetType() TokenType
	preludeExcludeTransform(bool) Token 
}

func ToLoc(t Token) info.Loc {
	loc := t.GetLocation()
	return info.MakeLocation(loc.GetLine(), loc.GetChar())
}

type ValueToken struct {
	Value value.Value
	Line  int
	Char  int
}

type IdToken struct {
	id   string
	line int
	char int
}

type TypeIdToken IdToken

type PreludeIdToken struct {
	tokenType TokenType
	id string
	line int
	char int
}

type OtherToken struct {
	tokenType TokenType
	line int
	char int
}

func (p PreludeIdToken) preludeExcludeTransform(exclude bool) Token {
	if exclude {
		return IdToken{id: p.id, line: p.line, char: p.char}
	}
	return OtherToken{tokenType: p.tokenType, line: p.line, char: p.char}
}

func (p PreludeIdToken) GetType() TokenType {
	return p.tokenType
}

func (p PreludeIdToken) ToString() string {
	return p.id
}

func (p PreludeIdToken) GetLocation() info.Location {
	return info.MakeLocation(p.line, p.char)
}

var inPrelude = []TokenType{
	PLUS, PLUS_PLUS, MINUS, STAR, HAT, SLASH, EQUALS_EQUALS,
	COLON, BANG, BANG_EQUALS, GREAT, GREAT_EQUALS, LESS, LESS_EQUALS,
	AMPER_AMPER, BAR_BAR,
}

// true iff receiver is an operator in prelude
func PreludeIncludes(tt TokenType) bool {
	for _, tokenType := range inPrelude {
		if tt == tokenType {
			return true
		}
	}
	return false
}

func (id TypeIdToken) AsType() types.Tau {
	return types.MakeTau(id.id, ToLoc(id))
}

func (o1 OtherToken) Equal_test_weak(o2 OtherToken) bool {
	return o1.tokenType == o2.tokenType
}
func (o1 OtherToken) Equal_test(o2 OtherToken) bool {
	return o1.Equal_test_weak(o2) && o1.char == o2.char && o1.line == o2.line
}
func MakeOtherToken(t TokenType, line int, char int) OtherToken {
	return OtherToken{
		tokenType: t,
		line:      line,
		char:      char,
	}
}

type ErrorToken struct {
	err err.UserMessage
}
type AnotationToken IdToken

func (in InputStream) MakeErrorLocation(loca info.Locatable) err.ErrorLocation {
	loc := loca.GetLocation()
	return err.MakeErrorLocation(
		loc.GetLine(),
		loc.GetChar(),
		in.path,
		in.source)
}

func (o OtherToken) ChangeTokenType(ty TokenType) Token {
	o.tokenType = ty
	return o
}

func (v ValueToken) GetValue() value.Value {
	return v.Value
}

func (id TypeIdToken) GetLocation() info.Location {
	return IdToken(id).GetLocation()
}
func (a AnotationToken) GetLocation() info.Location {
	return info.MakeLocation(a.line, a.char)
}
func (v ValueToken) GetLocation() info.Location {
	return info.MakeLocation(v.Line, v.Char)
}
func (id IdToken) GetLocation() info.Location {
	return info.MakeLocation(id.line, id.char)
}
func (o OtherToken) GetLocation() info.Location {
	return info.MakeLocation(o.line, o.char)
}
func (e ErrorToken) GetLocation() info.Location {
	return e.err.GetLocation()
}

func (a AnotationToken) GetType() TokenType {
	return AT
}
func (v ValueToken) GetType() TokenType {
	return VALUE
}
func (id IdToken) GetType() TokenType {
	return ID
}
func (id TypeIdToken) GetType() TokenType {
	return TYPE_ID
}
func (t OtherToken) GetType() TokenType {
	return t.tokenType
}
func (e ErrorToken) GetType() TokenType {
	return ERROR
}

func (a AnotationToken) preludeExcludeTransform(bool) Token {
	return a
}
func (v ValueToken) preludeExcludeTransform(bool) Token {
	return v
}
func (id IdToken) preludeExcludeTransform(bool) Token {
	return id
}
func (id TypeIdToken) preludeExcludeTransform(bool) Token {
	return id
}
func (t OtherToken) preludeExcludeTransform(bool) Token {
	return t
}
func (e ErrorToken) preludeExcludeTransform(bool) Token {
	return e
}

func (a AnotationToken) ToString() string {
	return "@" + a.id
}
func (v ValueToken) ToString() string {
	return v.Value.ToString()
}
func (id IdToken) ToString() string {
	return id.id
}
func (id TypeIdToken) ToString() string {
	return IdToken(id).ToString()
}
func (t OtherToken) ToString() string {
	return t.tokenType.ToString()
}
func (tt TokenType) ToString() string {
	switch tt {
	case ID:
		return "_ID_"
	case VALUE:
		return "_VALUE_"
	case LBRACK:
		return "["
	case RBRACK:
		return "]"
	case LPAREN:
		return "("
	case RPAREN:
		return ")"
	case LCURL:
		return "{"
	case RCURL:
		return "}"
	case STRING:
		return "String"
	case CHAR:
		return "Char"
	case BOOL:
		return "Bool"
	case FLOAT:
		return "Float"
	case INT:
		return "Int"
	case TYPE:
		return "type"
	case CLASS:
		return "class"
	case WHERE:
		return "where"
	case LET:
		return "let"
	case CONST:
		return "const"
	case MUT:
		return "mut"
	case LAZY:
		return "lazy"
	case PLUS:
		return "+"
	case PLUS_PLUS:
		return "++"
	case MINUS:
		return "-"
	case STAR:
		return "*"
	case SLASH:
		return "/"
	case HAT:
		return "^"
	case EQUALS:
		return "="
	case COLON:
		return ":"
	case COLON_COLON:
		return "::"
	case SEMI_COLON:
		return ";"
	case QUESTION:
		return "?"
	case COMMA:
		return ","
	case BANG:
		return "!"
	case BANG_EQUALS:
		return "!="
	case EQUALS_EQUALS:
		return "=="
	case BAR:
		return "|"
	case AMPER_AMPER:
		return "&&"
	case BAR_BAR:
		return "||"
	case ARROW:
		return "->"
	case FAT_ARROW:
		return "=>"
	case DOT:
		return "."
	case DOT_DOT:
		return ".."
	case GREAT:
		return ">"
	case LESS:
		return "<"
	case GREAT_EQUALS:
		return ">="
	case LESS_EQUALS:
		return "<="
	case UNDERSCORE:
		return "_"
	case EOF:
		return "_EOF_"
	case NEW_LINE:
		return "_NEW_LINE_"
	case ERROR:
		return "_ERROR_"
	case BACKSLASH:
		return `\`
	case AT:
		return "@"
	}
	return "_UNKNOWN_"
}
func (e ErrorToken) ToString() string {
	return e.err.ToString()
}

func MakeBlankToken() OtherToken {
	return OtherToken{tokenType: _ANY__, line: 0, char: 0}
}

func MakeIdToken(id string, line int, char int) IdToken {
	return IdToken{id: id, line: line, char: char}
}
