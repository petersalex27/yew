package pkg

import (
	"iter"
	"regexp"
	"strings"
)

type lexer struct {
	// input string
	input string
	// current position in input
	current int
}

func newLexer(input string) *lexer {
	return &lexer{input: input}
}

func (l *lexer) lexString() string {
	l.next()
	var str strings.Builder
	for {
		if r := l.next(); r == '"' {
			break
		} else {
			str.WriteRune(r)
		}
	}
	return str.String()
}

var (
	versionRegex = regexp.MustCompile(VERSION_REGEX)
	identRegex = regexp.MustCompile(IDENT_REGEX)
	pathRegex = regexp.MustCompile(PATH_REGEX)
	numberRegex = regexp.MustCompile(NUMBER_REGEX)
	stringRegex = regexp.MustCompile(STRING_REGEX)
)

func (l *lexer) readMatch(pattern *regexp.Regexp) string {
	loc := pattern.FindStringIndex(l.input[l.current:])
	if loc == nil || loc[0] != 0 {
		return ""
	}
	l.current += loc[1]
	return l.input[loc[0]:loc[1]]
}

func (l *lexer) lexVersion() string {
	iter.Seq[any](l.input[l.current:])
	l.readMatch(versionRegex)
}

func (l *lexer) lex() Token {
	for tok := (Token{}); ; {
		switch r := l.next(); {
		case r == 0:
			return Token{EOF, ""}
		case isWhitespace(r):
			l.ignore()
		case r == ':':
			return Token{COLON, string(r)}
		case r == '-':
			return Token{SWITCH, string(r)}
		case r == '"':
			l.backup()
			return Token{STRING, l.lexString()}
		case r == 'v':
			l.backup()
			return Token{VERSION, l.lexVersion()}
		case r == 'p':
			l.backup()
			return Token{PATH, l.lexPath()}
		case isLetter(r):
			l.backup()
			return Token{IDENT, l.lexIdent()}
		default:
			return Token{ILLEGAL, string(r)}
		}
	}
}


func (l *lexer) next() rune {
	if l.current >= len(l.input) {
		return 0
	}
	ch := l.input[l.current]
	l.current++
	return rune(ch)
}

func (l *lexer) peek() rune {
	if l.current >= len(l.input) {
		return 0
	}
	return rune(l.input[l.current])
}

func (l *lexer) backup() {
	l.current--
}

func (l *lexer) ignore() {
	l.current++
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}