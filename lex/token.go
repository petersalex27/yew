package scan

import (
	"yew/info"
	"yew/utils"
	"yew/value"

	err "yew/error"
)

type TokenType int

const (
	// variable 
	ID TokenType = iota
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
	PLUS_EQUALS
	MINUS_EQUALS
	STAR_EQUALS
	SLASH_EQUALS

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

	TRUE
	FALSE
	
	BACKSLASH

	AT

	UNDERSCORE

	NEW_LINE

	ERROR

	EOF

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
	GetSourceIndex() int
	//GetLineAndChar() (int, int)
	//GetLocation() info.Location
}
type ValueToken struct {
	value value.Value
	index int
	line int
	char int
}
type IdToken struct {
	id string
	index int
	line int
	char int
}
type OtherToken struct {
	tokenType TokenType
	index int
	line int
	char int
}
type ErrorToken struct {
	index int
	err err.UserMessage
}
type AnotationToken IdToken

func (v ValueToken) GetValue() value.Value {
	return v.value
}

func (v ValueToken) GetLocation() info.Location {
	return info.MakeLocation(v.line, v.char)
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

func (v ValueToken) GetSourceIndex() int {
	return v.index
}
func (id IdToken) GetSourceIndex() int {
	return id.index
}
func (o OtherToken) GetSourceIndex() int {
	return o.index
}
func (e ErrorToken) GetSourceIndex() int {
	return e.index
}

func (v ValueToken) GetType() TokenType {
	return VALUE
}
func (id IdToken) GetType() TokenType {
	return ID
} 
func (t OtherToken) GetType() TokenType {
	return t.tokenType
}
func (e ErrorToken) GetType() TokenType {
	return ERROR
}

func (v ValueToken) ToString() string {
	return v.value.ToString()
}
func (id IdToken) ToString() string {
	return id.id
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
		case PLUS_EQUALS:
			return "+="
		case MINUS_EQUALS:
			return "-="
		case STAR_EQUALS:
			return "*="
		case SLASH_EQUALS:
			return "/="
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
		case TRUE:
			return "True"
		case FALSE:
			return "False"
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

func MakeIdToken(id string, line int, char int) IdToken {
	return IdToken{id: id, line: line, char: char}
}