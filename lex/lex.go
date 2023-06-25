package scan

import (
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	err "yew/error"
	types "yew/type"
	"yew/value"
)

// underscore id token
var UnderscoreIdToken = IdToken{id: "_", index: 0, line: 0, char: 0}

type Lexer interface {
	Next() Token
	Match(TokenPattern) int
	GetPath() string
	GetSource() string
}

type InputStream struct {
	path            string
	streamIndex     int64
	streamLength    int64
	streamCapacity  int64
	asStringPattern string
	source          string

	tokens []Token
}

/*
GetTokenAtOffset returns token at current index + offset (offset can be negative).
The parser holds the token at offset -1 in its `Next` field and the token at offset -2
in its `Current` field; thus, to get the most recently forgoten token (the one that was
in `Current` before the current token in `Current`) use an offset of -3.
*/
func (stream *InputStream) GetTokenAtOffset(offset int64) Token {
	index := stream.streamIndex + offset
	if index < 0 || index >= stream.streamLength {
		err.PrintBug()
		panic("")
	}
	return stream.tokens[index]
}

// TODO: find average number of tokens based on length of stream, use this for init of `tokens`
const streamBufferFactor = 0.25

func calculateStreamBufferSize(offset int64, sourceLen int64) int64 {
	return offset + int64(float64(sourceLen)*streamBufferFactor) + 1
}
func InitStream(path string, sourceLen int64) InputStream {
	return InputStream{
		path:         path,
		streamIndex:  0,
		streamLength: 0,
		tokens:       make([]Token, 0, calculateStreamBufferSize(0, sourceLen)),
	}
}

type Input struct {
	//stored []Token
	lineNumber                 int
	prevLineLength             int
	charNumber                 int
	sourceIndex                int
	sourceLength               int
	path                       string
	source                     string
	scanningInterpolatedString int
	endingInterpolatedString   int
}

func (in InputStream) GetPath() string   { return in.path }
func (in Input) GetPath() string         { return in.path }
func (in InputStream) GetSource() string { return in.source }
func (in Input) GetSource() string       { return in.source }

func Init(path string) (Input, error) {
	input := Input{sourceIndex: 0, path: path}
	//fmt.Println(os.Getwd())
	f, err := os.Open(path)
	if nil != err {
		return input, err
	}
	tmp, err2 := io.ReadAll(f)
	f.Close()
	if nil != err2 {
		return input, err2
	}
	input.source = string(tmp)
	input.charNumber = 0 // 0 marks no chars read on line (after new line or before reading anything)
	input.lineNumber = 1
	input.sourceLength = len(input.source)
	return input, nil
}

func (in *Input) peek() (c byte) {
	if in.sourceIndex >= in.sourceLength {
		return 0
	}
	return in.source[in.sourceIndex]
}

func (in *Input) nextChar() (c byte) {
	if in.sourceIndex >= len(in.source) {
		return 0
	}

	c = in.source[in.sourceIndex]
	if c == '\n' {
		in.prevLineLength = in.charNumber
		in.charNumber = 0
		in.lineNumber++
	} else {
		in.charNumber++
	}
	in.sourceIndex++
	return
}

func (in *Input) readUntil(until byte) byte {
	for {
		c := in.nextChar()
		if c == until {
			return until
		} else if until == 0 || c == 0 {
			return 0
		}
	}
}

func (in *Input) skipWhitespace() (c byte) {
	for c = in.nextChar(); c != 0 && in.sourceIndex < len(in.source); c = in.nextChar() {
		if c == '-' && in.peek() == '-' {
			// skip single-line comments
			return in.readUntil('\n')
		} else if c == '-' && in.peek() == '*' {
			// skip multi-line comments
			for c = in.readUntil('*'); c != 0; c = in.readUntil('*') {
				if in.peek() == '-' {
					in.nextChar()
					return in.nextChar()
				}
			}
		} else if c == ' ' || c == '\t' {
			continue
		}
		break
	}
	return
}

func (in *Input) tokenUnit(tokenType TokenType) Token {
	return OtherToken{tokenType: tokenType, line: in.lineNumber, char: in.charNumber}
}

// getEscape scans and returns result of escape sequence.
// The argument `isChar` controls whether certain escape sequences are allowed,
// e.g., "\uXXXX" is allowed for strings (char sequences) but not for a
// single char.
func (in *Input) getEscape(isChar bool) (byte, *ErrorToken) {
	c := in.nextChar()
	switch c {
	case 'n':
		return '\n', nil
	case 't':
		return '\t', nil
	case 'b':
		return '\b', nil
	case 'v':
		return '\v', nil
	case 'r':
		return '\r', nil
	case 'u':
		if isChar {
			et := new(ErrorToken)
			*et = inputErrors[E_STRING_ONLY_ESCAPE](in)
			return 0, et
		}
		return 1, nil // 1 represents request for unicode
	case '"':
		return '"', nil
	case '\\':
		return '\\', nil
	case '\'':
		return '\'', nil
	case '%':
		return '%', nil
	}

	et := new(ErrorToken)
	*et = inputErrors[E_ILLEGAL_ESCAPE](in)
	return 0, et
}

func isUnicodeRequest(c byte) bool {
	return c == 1
}

func isEscapeFailure(c byte) bool {
	return c == 0
}

func stringValue(s string) value.Array {
	return value.MakeArray[types.Char]([]value.Char(s))
}

// regex for unicode escape sequence
var UNICODE_REGEX = regexp.MustCompile("[0-9a-fA-F]{1,4}")

// returns unicode code point converted to a string
func (in *Input) getUnicode() (string, *ErrorToken) {
	loc := UNICODE_REGEX.FindStringIndex(in.source[in.sourceIndex:])
	if loc == nil || loc[0] != 0 {
		// failed to find code point || failed to find code point next to escape char
		et := new(ErrorToken)
		*et = inputErrors[E_BAD_UNICODE](in)
		return "", et
	}

	// get unicode as int/code-point
	res, _ := strconv.ParseInt(in.source[in.sourceIndex:in.sourceIndex+loc[1]], 16, 32)

	var builder strings.Builder  // here just for converting from code point to string
	builder.WriteRune(rune(res)) // write code point as string

	// update location information
	in.charNumber += loc[1]
	in.sourceIndex += loc[1]

	return builder.String(), nil
}

// get a string value token
func (in *Input) getString() Token {
	start := in.charNumber - 1
	index := in.sourceIndex
	var builder strings.Builder
	for c := in.nextChar(); ; c = in.nextChar() {
		if unicode.IsControl(rune(c)) && c != '\t' {
			return inputErrors[E_UNEXPECTED_CONTROL](in)
		} else if c == '\\' {
			var errorToken *ErrorToken
			c, errorToken = in.getEscape(false)
			if isEscapeFailure(c) {
				return *errorToken
			} else if isUnicodeRequest(c) {
				var uc string
				uc, errorToken = in.getUnicode()
				if nil != errorToken {
					// failure, return error
					return *errorToken
				}
				// success, write unicode
				builder.WriteString(uc)
			} else {
				builder.WriteByte(c)
			}
		} else if c == '"' {
			out := stringValue(builder.String())
			return ValueToken{out, index, in.lineNumber, start}
		} else if c == 0 {
			return inputErrors[E_END_OF_FILE](in)
		} else if c == '%' {
			// TODO
			/*if in.peek() == '{' { // do not consume character, will read later
				in.scanningInterpolatedString++
				out := stringValue(builder.String())
				return ValueToken{out, index, in.lineNumber, start}
			}*/
			// else just write '%'
			builder.WriteByte(c)
		} else {
			builder.WriteByte(c)
		}
	}
}

func (in *Input) getChar() Token {
	c := in.nextChar()
	if unicode.IsControl(rune(c)) && c != '\t' {
		return inputErrors[E_UNEXPECTED_CONTROL](in)
	} else if c == '\\' {
		var et *ErrorToken
		c, et = in.getEscape(true)
		if isEscapeFailure(c) {
			// failure, return error
			return *et
		}
		// success, continue
	} else if c == 0 {
		return inputErrors[E_END_OF_FILE](in)
	} else {
		// ignore
	}

	if in.nextChar() != '\'' {
		return inputErrors[E_EXPECTED_CHAR_CLOSE](in)
	}

	return ValueToken{value.Char(c), in.sourceIndex, in.lineNumber, in.charNumber}
}

type _errorGenFn (func(in *Input) ErrorToken)

func inputErrorGen(message string) _errorGenFn {
	return (func(in *Input) ErrorToken {
		return ErrorToken{in.sourceIndex, err.CompileMessage(message, err.ERROR, err.INPUT, in.path, in.lineNumber, in.charNumber, in.sourceIndex, in.source)}
	})
}

type inputErrorType int

const (
	E_END_OF_FILE = iota
	E_ILLEGAL_ESCAPE
	E_UNEXPECTED_TOKEN
	E_EXPECTED_CHAR_CLOSE
	E_EXPECTED_RAW_ID_CLOSE
	E_UNEXPECTED_CONTROL
	E_TRAILING_UNDERSCORE
	E_STRING_ONLY_ESCAPE
	E_BAD_UNICODE
)

var inputErrors = map[inputErrorType]_errorGenFn{
	E_END_OF_FILE:           inputErrorGen("unexpected end of file"),
	E_ILLEGAL_ESCAPE:        inputErrorGen("illegal escape sequence"),
	E_UNEXPECTED_TOKEN:      inputErrorGen("unexpected token"),
	E_EXPECTED_CHAR_CLOSE:   inputErrorGen("expected end of character literal"),
	E_EXPECTED_RAW_ID_CLOSE: inputErrorGen("expected closing backtick"),
	E_UNEXPECTED_CONTROL:    inputErrorGen("unexpected control character"),
	E_TRAILING_UNDERSCORE:   inputErrorGen("trailing digit seperator"),
	E_STRING_ONLY_ESCAPE:    inputErrorGen("escape sequence is not allowed for char literals"),
	E_BAD_UNICODE:           inputErrorGen("malformed unicode code point escape sequence"),
}

var ID_REGEX = regexp.MustCompile("[A-Za-z][A-Za-z0-9_]*(')*")
var EXT_ID_REGEX = regexp.MustCompile(`[^0-9\{\}\[\]\(\)\s\x60[[:cntrl:]]@]+`)
var HEX_REGEX = regexp.MustCompile("0(x|X)([0-9A-Fa-f]_?)*[0-9A-Fa-f]+")
var OCT_REGEX = regexp.MustCompile("0(o|O)([0-7]_?)*[0-7]+")
var BIN_REGEX = regexp.MustCompile("0(b|B)([01]_?)*[01]+")
var INT_REGEX = regexp.MustCompile("([0-9]_?)*[0-9]+")
var FLOAT_REGEX = regexp.MustCompile(`[0-9]+(((\.[0-9]*)?(e|E)(\+|-)?([0-9]_?)*[0-9]+)|(\.[0-9]*))`)

var keywords = map[string]Token{
	"Int":     OtherToken{INT, 0, 0, 0},
	"Bool":    OtherToken{BOOL, 0, 0, 0},
	"Char":    OtherToken{CHAR, 0, 0, 0},
	"Float":   OtherToken{FLOAT, 0, 0, 0},
	"String":  OtherToken{STRING, 0, 0, 0},
	"class":   OtherToken{CLASS, 0, 0, 0},
	"where":   OtherToken{WHERE, 0, 0, 0},
	"lazy":    OtherToken{LAZY, 0, 0, 0},
	"let":     OtherToken{LET, 0, 0, 0},
	"mut":     OtherToken{MUT, 0, 0, 0},
	"const":   OtherToken{CONST, 0, 0, 0},
	"mod":     OtherToken{MOD, 0, 0, 0},
	"True":    ValueToken{value.Bool(true), 0, 0, 0},
	"False":   ValueToken{value.Bool(false), 0, 0, 0},
	"package": OtherToken{PACKAGE, 0, 0, 0},
	"module":  OtherToken{MODULE, 0, 0, 0},
}

func (in *Input) getIdOrKeyword2(forceId bool) Token {
	loc := ID_REGEX.FindStringIndex(in.source[in.sourceIndex:])
	if loc == nil || loc[0] != 0 {
		// this should be impossible
		err.PrintBug()
	}

	str := in.source[in.sourceIndex : in.sourceIndex+loc[1]]
	var key Token
	var found bool = false // keyword is never found when id is forced
	if !forceId {
		key, found = keywords[str]
	}
	start := in.charNumber
	in.charNumber += loc[1]
	in.sourceIndex += loc[1]
	if !found {
		// always returns here when id is forced
		return IdToken{id: str, line: in.lineNumber, char: start}
	}
	if VALUE == key.GetType() {
		out := key.(ValueToken)
		out.Char = start
		out.Line = in.lineNumber
		return out
	} else {
		out := key.(OtherToken)
		out.char = start
		out.line = in.lineNumber
		return out
	}
}

func (in *Input) getIdOrKeyword() Token {
	return in.getIdOrKeyword2(false)
}

func (in *Input) hasTrailingUnderscore() bool {
	if in.peek() == '_' {
		in.charNumber++
		in.sourceIndex++
		return true
	}
	return false
}

func removeUnderscore(s string) string {
	var builder strings.Builder
	for _, c := range s {
		if c != '_' {
			builder.WriteRune(c)
		}
	}
	return builder.String()
}

func (in *Input) getNumber() Token {
	input := in.source[in.sourceIndex:]
	loc := FLOAT_REGEX.FindStringIndex(input)
	if nil != loc && loc[0] == 0 {
		res, _ := strconv.ParseFloat(removeUnderscore(input[:loc[1]]), 64)
		start := in.charNumber
		in.charNumber += loc[1]
		in.sourceIndex += loc[1]
		if in.hasTrailingUnderscore() {
			return inputErrors[E_TRAILING_UNDERSCORE](in)
		}
		return ValueToken{Value: value.Float(res), Line: in.lineNumber, Char: start}
	}

	loc = HEX_REGEX.FindStringIndex(input)
	if nil != loc && loc[0] == 0 {
		res, _ := strconv.ParseInt(removeUnderscore(input[:loc[1]]), 0, 64)
		start := in.charNumber
		in.charNumber += loc[1]
		in.sourceIndex += loc[1]
		if in.hasTrailingUnderscore() {
			return inputErrors[E_TRAILING_UNDERSCORE](in)
		}
		return ValueToken{Value: value.Int(res), Line: in.lineNumber, Char: start}
	}

	loc = OCT_REGEX.FindStringIndex(input)
	if nil != loc && loc[0] == 0 {
		res, _ := strconv.ParseInt(removeUnderscore(input[:loc[1]]), 0, 64)
		start := in.charNumber
		in.charNumber += loc[1]
		in.sourceIndex += loc[1]
		if in.hasTrailingUnderscore() {
			return inputErrors[E_TRAILING_UNDERSCORE](in)
		}
		return ValueToken{Value: value.Int(res), Line: in.lineNumber, Char: start}
	}

	loc = BIN_REGEX.FindStringIndex(input)
	if nil != loc && loc[0] == 0 {
		res, _ := strconv.ParseInt(removeUnderscore(input[:loc[1]]), 0, 64)
		start := in.charNumber
		in.charNumber += loc[1]
		in.sourceIndex += loc[1]
		if in.hasTrailingUnderscore() {
			return inputErrors[E_TRAILING_UNDERSCORE](in)
		}
		return ValueToken{Value: value.Int(res), Line: in.lineNumber, Char: start}
	}

	loc = INT_REGEX.FindStringIndex(input)
	if nil != loc && loc[0] == 0 {
		res, _ := strconv.ParseInt(removeUnderscore(input[:loc[1]]), 10, 64)
		start := in.charNumber
		in.charNumber += loc[1]
		in.sourceIndex += loc[1]
		if in.hasTrailingUnderscore() {
			return inputErrors[E_TRAILING_UNDERSCORE](in)
		}
		return ValueToken{Value: value.Int(res), Line: in.lineNumber, Char: start}
	}

	err.PrintBug()
	panic("")
}

func (in *Input) ungetChar() {
	if in.sourceIndex <= 0 {
		return
	}

	in.sourceIndex--
	if in.source[in.sourceIndex] == '\n' {
		in.lineNumber--
		if in.prevLineLength == 0 {
			// need to find previous location since only one is saved
			in2 := Input{lineNumber: 1, charNumber: 0, sourceIndex: 0, path: in.path, source: in.source}
			for ; in.lineNumber > in2.lineNumber; in2.readUntil('\n') {
				if in2.peek() == 0 {
					// prevents loop from spinning forever (shouldn't happen, but in case I missed something)
					err.PrintBug()
					panic("")
				}
			}
			in.lineNumber = in2.lineNumber
			in.charNumber = in2.charNumber
			in.prevLineLength = in2.prevLineLength
			in.sourceIndex = in2.sourceIndex
			return
		}

		in.charNumber = in.prevLineLength
		in.prevLineLength = 0
	} else {
		in.charNumber--
	}
}

/*func (in *Input) nonAsciiId() Token {
	loc := EXT_ID_REGEX.FindStringIndex(in.source[in.sourceIndex:])
	if loc == nil || loc[0] != 0 {
		return inputErrors[E_UNEXPECTED_TOKEN](in)
	}
	length := loc[1]
	id := in.source[in.sourceIndex:in.sourceIndex+length]
	charBefore := in.charNumber
	indexBefore := in.sourceIndex
	in.charNumber += length
	in.sourceIndex += length
	return IdToken{id: id, index: indexBefore, line: in.lineNumber, char: charBefore}
} */

func (in *Input) Next() Token {
	/*if length := len(in.stored); length > 0 {
		tok := in.stored[length - 1]
		in.stored = in.stored[:length - 1]
		return tok
	}*/

	c := in.skipWhitespace()
	if unicode.IsDigit(rune(c)) {
		in.ungetChar()
		return in.getNumber()
	} else if nil != ID_REGEX.FindStringIndex(string(c)) {
		in.ungetChar()
		return in.getIdOrKeyword()
	}

	switch c {
	case '+':
		if in.peek() == '+' {
			in.nextChar()
			return in.tokenUnit(PLUS_PLUS)
		} else if in.peek() == '=' {
			in.nextChar()
			return in.tokenUnit(PLUS_EQUALS)
		}
		return in.tokenUnit(PLUS)
	case '-':
		if in.peek() == '=' {
			in.nextChar()
			return in.tokenUnit(MINUS_EQUALS)
		} else if in.peek() == '>' {
			in.nextChar()
			return in.tokenUnit(ARROW)
		}
		return in.tokenUnit(MINUS)
	case '*':
		if in.peek() == '=' {
			in.nextChar()
			return in.tokenUnit(STAR_EQUALS)
		}
		return in.tokenUnit(STAR)
	case '/':
		if in.peek() == '=' {
			in.nextChar()
			return in.tokenUnit(SLASH_EQUALS)
		}
		return in.tokenUnit(SLASH)
	case '=':
		if in.peek() == '=' {
			in.nextChar()
			return in.tokenUnit(EQUALS_EQUALS)
		} else if in.peek() == '>' {
			in.nextChar()
			return in.tokenUnit(FAT_ARROW)
		}
		return in.tokenUnit(EQUALS)
	case ':':
		if in.peek() == ':' {
			in.nextChar()
			if in.peek() == '=' {
				in.nextChar()
				return in.tokenUnit(COLON_COLON_EQUAL)
			}
			return in.tokenUnit(COLON_COLON)
		}
		return in.tokenUnit(COLON)
	case ';':
		return in.tokenUnit(SEMI_COLON)
	case '\n':
		in.prevLineLength = in.charNumber
		in.lineNumber++
		in.charNumber = 0
		return in.tokenUnit(NEW_LINE)
	case '!':
		if in.peek() == '=' {
			in.nextChar()
			return in.tokenUnit(BANG_EQUALS)
		}
		return in.tokenUnit(BANG)
	case '>':
		if in.peek() == '=' {
			in.nextChar()
			return in.tokenUnit(GREAT_EQUALS)
		}
		return in.tokenUnit(GREAT)
	case '<':
		if in.peek() == '=' {
			in.nextChar()
			return in.tokenUnit(LESS_EQUALS)
		}
		return in.tokenUnit(LESS)
	case '?':
		return in.tokenUnit(QUESTION)
	case '(':
		return in.tokenUnit(LPAREN)
	case ')':
		return in.tokenUnit(RPAREN)
	case '[':
		return in.tokenUnit(LBRACK)
	case ']':
		return in.tokenUnit(RBRACK)
	case '{':
		return in.tokenUnit(LCURL)
	case '}':
		return in.tokenUnit(RCURL)
	case '_':
		return in.tokenUnit(UNDERSCORE)
	case ',':
		return in.tokenUnit(COMMA)
	case '.':
		if in.peek() == '.' {
			in.nextChar()
			return in.tokenUnit(DOT_DOT)
		}
		return in.tokenUnit(DOT)
	case '^':
		return in.tokenUnit(HAT)
	case '&':
		c = in.nextChar()
		if c != '&' {
			return inputErrors[E_UNEXPECTED_TOKEN](in)
		}
		return in.tokenUnit(AMPER_AMPER)
	case '|':
		if in.peek() == '|' {
			in.nextChar()
			return in.tokenUnit(BAR_BAR)
		}
		return in.tokenUnit(BAR)
	/*case '`':
	tok := in.nonAsciiId()
	if tok.GetType() == ERROR {
		return tok
	}
	if in.nextChar() != '`' {
		return inputErrors[E_EXPECTED_RAW_ID_CLOSE](in)
	}
	return tok*/
	case '"':
		return in.getString()
	case '\'':
		return in.getChar()
	case '\\':
		return in.tokenUnit(BACKSLASH)
	case '@':
		return AnotationToken(in.getIdOrKeyword2(true).(IdToken))
	case 0:
		return in.tokenUnit(EOF)
	default:
		/*in.ungetChar()
		return in.nonAsciiId()*/
		return inputErrors[E_UNEXPECTED_TOKEN](in)
	}
	//return inputErrors[E_UNEXPECTED_TOKEN](in)
}

var patternMap = map[TokenType]byte{
	ID:             '0',
	VALUE:          '1',
	LBRACK:         '2',
	RBRACK:         '3',
	LPAREN:         '4',
	RPAREN:         '5',
	LCURL:          '6',
	RCURL:          '7',
	STRING:         '8',
	CHAR:           '9',
	BOOL:           'a',
	FLOAT:          'b',
	INT:            'c',
	TYPE:           'd',
	CLASS:          'e',
	WHERE:          'f',
	LET:            'g',
	CONST:          'h',
	MUT:            'i',
	LAZY:           'j',
	PLUS:           'k',
	PLUS_PLUS:      'l',
	MINUS:          'm',
	STAR:           'n',
	SLASH:          'o',
	HAT:            'p',
	EQUALS:         'q',
	PLUS_EQUALS:    'r',
	MINUS_EQUALS:   's',
	STAR_EQUALS:    't',
	SLASH_EQUALS:   'u',
	COLON:          'v',
	COLON_COLON:    'w',
	SEMI_COLON:     'x',
	QUESTION:       'y',
	COMMA:          'z',
	BANG:           'A',
	BANG_EQUALS:    'B',
	EQUALS_EQUALS:  'C',
	BAR:            'D',
	AMPER_AMPER:    'E',
	BAR_BAR:        'F',
	ARROW:          'G',
	FAT_ARROW:      'H',
	DOT:            'I',
	DOT_DOT:        'J',
	GREAT:          'K',
	LESS:           'L',
	GREAT_EQUALS:   'M',
	LESS_EQUALS:    'N',
	BACKSLASH:      'O',
	AT:             'P',
	UNDERSCORE:     'Q',
	NEW_LINE:       'R',
	ERROR:          'S',
	EOF:            'T',
	_ALTERNATION__: '|',
	_START_GROUP__: '(',
	_END_GROUP__:   ')',
	_ANY__:         '.',
	END__:          ' ',
}

func GetTokenPattern(ts []TokenType) (pat string, newlines int) {
	newlines = 0
	var builder strings.Builder
	for i := 0; i < len(ts); i++ {
		if ts[i] == _REPEAT__ {
			if len(ts) > i+1 {
				i += 1
				if int(ts[i]) == 0 {
					builder.WriteByte('*')
				} else if int(ts[i]) == -1 {
					builder.WriteByte('+')
				} else if int(ts[i]) == -2 {
					builder.WriteByte('?')
				} else {
					if len(ts) <= i+1 {
						err.PrintBug()
						panic("")
					}

					builder.WriteByte('{')
					builder.WriteString(strconv.Itoa(int(ts[i])))
					builder.WriteByte(',')
					builder.WriteString(strconv.Itoa(int(ts[i+1])))
					builder.WriteByte('}')
					i += 1
				}
			} else {
				err.PrintBug()
				panic("")
			}
		} else {
			c, found := patternMap[ts[i]]
			if !found {
				err.PrintBug()
				panic("")
			}

			if NEW_LINE == ts[i] {
				newlines++
			}

			builder.WriteByte(c)
		}
	}
	return builder.String(), newlines
}

type TokenPattern struct {
	newLines int
	pattern  *regexp.Regexp
}

func CompileTokenPattern(ts []TokenType) TokenPattern {
	pat, nl := GetTokenPattern(ts)
	return TokenPattern{pattern: regexp.MustCompile(pat), newLines: nl}
}

func (in InputStream) Match(pattern TokenPattern) int {
	if in.asStringPattern == "" {
		ts := make([]TokenType, len(in.tokens))
		for i, t := range in.tokens {
			ts[i] = t.GetType()
		}
		in.asStringPattern, _ = GetTokenPattern(ts)
	}

	loc := pattern.pattern.FindStringIndex(in.asStringPattern[in.streamIndex:])
	if nil == loc || loc[0] != 0 {
		return 0 // not found
	}
	return loc[1]
}

func (in *InputStream) Next() Token {
	if in.streamIndex+1 >= in.streamLength {
		if in.tokens[in.streamIndex].GetType() != EOF {
			err.PrintBug()
			panic("")
		}
		return in.tokens[in.streamIndex] // should return EOF
	}
	in.streamIndex++
	return in.tokens[in.streamIndex-1]
}

// grows stream's buffer with respect to length of unread input; always grows buffer by at least one
// if there are no issues allocating memory for the new buffer
func (stream *InputStream) grow(in *Input) {
	// number of chars left to read
	remainingSourceLength := in.sourceLength - in.sourceIndex
	// new size
	size := calculateStreamBufferSize(stream.streamLength, int64(remainingSourceLength))
	// new buffer
	newBuff := make([]Token, stream.streamLength, size)
	copy(newBuff, stream.tokens)

	stream.tokens = newBuff
	stream.streamCapacity = size
}

func (stream *InputStream) addNextToken(in *Input) (addedType TokenType) {
	if stream.streamLength+1 == stream.streamCapacity {
		stream.grow(in)
	}

	tok := in.Next()
	addedType = tok.GetType()
	add := true
	if addedType == NEW_LINE && stream.streamLength > 0 { 
		// only add one new line in sequences of new lines
		prev := stream.tokens[stream.streamLength - 1].GetType()
		if prev == NEW_LINE || prev == SEMI_COLON {
			add = false
		}
	}
	
	if add {
		stream.tokens = append(stream.tokens, tok)
		stream.streamLength++
	}
	return
}

func TokenizeFromInput(in *Input) (InputStream, *err.Error) {
	errOut := new(err.Error)
	stream := InitStream(in.path, int64(in.sourceLength))
	stream.source = in.source
	for {
		// always adds a token (poss. EOF or ERROR)
		tokType := stream.addNextToken(in)

		if ERROR == tokType {
			*errOut = stream.tokens[stream.streamLength-1].(ErrorToken).err.(err.Error)
			return stream, errOut // failure, return current stream and error
		} else if EOF == tokType {
			/*for _, t := range stream.tokens {
				fmt.Println(t.ToString())
			}*/
			return stream, nil // success, return stream
		}
		// else continue
	}
}

// Tokenize takes entire input and converts it into a sequence of tokens.
//   - (stream, nil) on success, (stream, someError) on failure
func Tokenize(path string) (InputStream, *err.Error) {
	errOut := new(err.Error)
	input, e := Init(path)
	if nil != e {
		*errOut = err.SystemError(e.Error()).(err.Error)
		return InputStream{}, errOut // failure, return empty stream and error
	}

	return TokenizeFromInput(&input)
}
