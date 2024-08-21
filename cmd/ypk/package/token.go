package pkg

import "strings"

type Token struct {
	TokenType
	Value string
}

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	PACKAGE
	REQUIRE

	IDENT

	STRING
	VERSION
	NUMBER
	PATH

	COLON
	SWITCH
)

// version stuff
const (
	indirect_version = `latest\+?|stable`
	version_range = `least|most`
	at_indirect = `(v@(` + indirect_version + `))`
	at_range = `(@(` + version_range + `))`
	version_number = `(v([\d+]\.)*[\d+])`
)

// number stuff
const (
	sci_notation_tail = `([eE][+-]?\d+)`
	int_head = `([+-]?\d+)`
	floating_tail = `(\.\d+)`
)

// string stuff
const (
	escape_characters = `(\\[abfnrtv'"])`
	not_quote_or_newline = `[^"\n]`
	standard_string = `("(` + escape_characters + `|` + not_quote_or_newline + `)*")`
	raw_string = "(`[^`]*`)"
)

// ident stuff
const (
	first_char = `[a-zA-Z]`
	rest_char = `((_?[a-zA-Z0-9]+)*)`
)

// path stuff
const (
	// valid protocol names:
	//		- https - https
	//		- ftp - file transfer protocol
	//		- file - file protocol
	//		- ssh - secure shell protocol
	//		- git - git protocol
	//		- svn - subversion protocol
	//		- http - http protocol
	protocol_regex = `(https|ftp|file|ssh|git|svn|http)`
)

const (
	HORIZONTAL_WHITESPACE = `(\h*)`
	IDENT_REGEX = `(` + first_char + rest_char + `)`
	STRING_REGEX = `(` + standard_string + `|` + raw_string + `)`
	VERSION_REGEX = `(` + at_indirect +`|` + version_number + at_range + `?)`
	NUMBER_REGEX = `(` + int_head + floating_tail + `?` + sci_notation_tail + `?)`
	PATH_REGEX = `((` + protocol_regex + `?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?)`
	PACKAGE_REGEX = `(package\b|package(?=:))`
	REQUIRE_REGEX = `(require\b|require(?=:))`
	ARG_LINE = `(\h*-\h*(\d|\w|[\(\)\*\+/><,\.;:'"\{\}!@\$%\^\&=-~`+ "`" +`\|\?])+)`
)

var tokens = map[string]TokenType{
	"package": PACKAGE,
	"require": REQUIRE,
}