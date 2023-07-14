package scan

import (
	//"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	err "yew/error"
	"yew/source"
	types "yew/type"
	"yew/value"
)

// underscore id token
var UnderscoreIdToken = IdToken{id: "_", line: 0, char: 0}

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
	excludePrelude  bool
	source          source.Source

	tokens []Token
}

func CreateInputStream(path string, streamIndex int64, source []string, tokens ...Token) InputStream {
	cap := len(tokens)
	length := len(tokens)
	return InputStream{
		path: path, 
		streamIndex: streamIndex, 
		source: source, 
		streamLength: int64(length),
		streamCapacity: int64(cap),
		tokens: tokens,
	}
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
		excludePrelude: false,
		path:         path,
		streamIndex:  0,
		streamLength: 0,
		tokens:       make([]Token, 0, calculateStreamBufferSize(0, sourceLen)),
	}
}

type Input struct {
	lineNumber int
	charNumber int
	path       string
	source     source.Source
}

func (a Input) equal_test(b Input) bool {
	ok := a.lineNumber == b.lineNumber &&
		a.charNumber == b.charNumber &&
		a.path == b.path &&
		len(a.source) == len(b.source)

	if !ok {
		return false
	}

	for i, source := range a.source {
		if source != b.source[i] {
			return false
		}
	}
	return true
}

func (a Input) toString_test() string {
	sourceStringify := func(source []string) string {
		var builder strings.Builder
		for _, s := range source {
			builder.WriteString("\t\t")
			builder.WriteString(s)
		}
		return builder.String()
	}

	return "Input{ " +
		"\n\tlineNumber: " + strconv.Itoa(a.lineNumber) +
		",\n\tcharNumber: " + strconv.Itoa(a.charNumber) +
		",\n\tpath: " + a.path +
		",\n\tsource: ```\n" + sourceStringify(a.source) +
		"```\n}"
}

func (in InputStream) GetPath() string     { return in.path }
func (in Input) GetPath() string           { return in.path }
func (in InputStream) GetSource() []string { return in.source }
func (in Input) GetSource() []string       { return in.source }

func Init(path string) (Input, error) {
	input := Input{path: path}
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
	input.source = make([]string, 0, 1)
	for i, j := 0, 0; i < len(tmp); i++ {
		if tmp[i] == '\n' || i+1 == len(tmp) {
			input.source = append(input.source, string(tmp[j:i+1]))
			j = i + 1
		}
	}
	input.charNumber = 0 // 0 marks no chars read on line (after new line or before reading anything)
	input.lineNumber = 1
	return input, nil
}

func (in *Input) getCurrentLine() string {
	return in.source[in.lineNumber-1]
}

func (in *Input) getNextIndexes() (line int, char int, eof bool) {
	eof = false
	line = in.lineNumber
	char = in.charNumber

	//println("len(\"\n\") =", len("\n"))

	if in.lineNumber == len(in.source) {
		if in.charNumber == len(in.getCurrentLine()) {
			eof = true
		} else {
			char = char + 1
		}
	} else if in.charNumber == len(in.getCurrentLine()) {
		line = line + 1
		char = 1
	} else {
		char = char + 1
	}

	return
}

func (in *Input) isAtEndOfLine() bool { return in.charNumber == len(in.getCurrentLine()) }
func (in *Input) isInFinalLine() bool { return in.lineNumber == len(in.source) }

func (in *Input) isAtEof() bool {
	return in.isInFinalLine() && in.isAtEndOfLine()
}

func (in *Input) ungetNextIndexes() (line int, char int) {
	line = in.lineNumber
	char = in.charNumber

	if in.charNumber <= 1 {
		if in.lineNumber > 1 {
			line = line - 1
			char = len(in.source.GetLine(line))
		}
	} else {
		char = char - 1
	}

	return
}

func (in *Input) unsetNextIndexes() (line int, char int, c byte) {
	line, char = in.ungetNextIndexes()
	in.charNumber = char
	in.lineNumber = line
	c = in.source[line-1][char-1]
	return
}

func (in *Input) setNextIndexes() (line int, char int, c byte) {
	var eof bool
	line, char, eof = in.getNextIndexes()
	in.lineNumber = line
	in.charNumber = char
	if eof {
		c = 0
	} else {
		c = in.source[line-1][char-1]
	}
	return
}

func (in *Input) peek() byte {
	line, char, eof := in.getNextIndexes()
	if eof {
		return 0
	}
	return in.source[line-1][char-1]
}

func (in *Input) nextChar() byte {
	_, _, next := in.setNextIndexes()
	return next
}

func (in *Input) readUntil(until byte) (c byte) {
	c = 0
	if until == 0 {
		return
	}

	for {
		c = in.nextChar()
		if c == until || c == 0 {
			return
		}
	}
}

func (in *Input) skipWhitespace() (c byte, errorToken *ErrorToken) {
	errorToken = nil

	for c = in.nextChar(); c != 0; c = in.nextChar() {
		if c == '-' && in.peek() == '-' {
			// skip single-line comments
			c = in.readUntil('\n')
			return
		} else if c == '-' && in.peek() == '*' {
			// skip multi-line comments
			for c = in.readUntil('*'); ; c = in.readUntil('*') {
				if c == 0 {
					errorToken = new(ErrorToken)
					*errorToken = inputErrors[E_UNTERMINATED_COMMENT](in)
					return
				}
				if in.peek() == '-' {
					c = in.nextChar()
					break
				}
			}
			if c == '-' {
				continue // end of comment, continue eating whitespace
			}
		} else if c == ' ' || c == '\t' {
			continue
		}
		break
	}
	return
}

func (in *Input) tokenUnitN(tokenType TokenType, n int) Token {
	return OtherToken{
		tokenType: tokenType,
		line: in.lineNumber,
		char: in.charNumber - (n - 1),
	}
}

func (in *Input) tokenUnit(tokenType TokenType) Token {
	return in.tokenUnitN(tokenType, 1)
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

func (in *Input) getSourceSlice() string {
	return in.source.GetLineSlice(in.lineNumber, in.charNumber)
}

func (in *Input) getSourceSliceN(n int) string {
	return in.source.GetLineSliceN(in.lineNumber, in.charNumber, n)
}

// regex for unicode escape sequence
var UNICODE_REGEX = regexp.MustCompile("[0-9a-fA-F]{1,4}")

// returns unicode code point converted to a string
func (in *Input) getUnicode() (string, *ErrorToken) {
	loc := UNICODE_REGEX.FindStringIndex(in.getSourceSlice())
	if loc == nil || loc[0] != 0 {
		// failed to find code point || failed to find code point next to escape char
		et := new(ErrorToken)
		*et = inputErrors[E_BAD_UNICODE](in)
		return "", et
	}

	// get unicode as int/code-point
	res, _ := strconv.ParseInt(in.getSourceSliceN(loc[1]), 16, 32)

	var builder strings.Builder  // here just for converting from code point to string
	builder.WriteRune(rune(res)) // write code point as string

	// update location information
	in.charNumber += loc[1]

	return builder.String(), nil
}

// get a string value token
func (in *Input) getString() Token {
	start := in.charNumber - 1
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
			return ValueToken{out, in.lineNumber, start}
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
	} else if c == '\'' {
		return inputErrors[E_ILLEGAL_CHAR](in)
	} else {
		// ignore
	}

	if in.nextChar() != '\'' {
		in.ungetChar()
		return inputErrors[E_EXPECTED_CHAR_CLOSE](in)
	}

	return ValueToken{value.Char(c), in.lineNumber, in.charNumber}
}

type _errorGenFn (func(in *Input) ErrorToken)

func inputErrorGen(message string) _errorGenFn {
	return (func(in *Input) ErrorToken {
		return ErrorToken{err.CompileMessage(message, err.ERROR, err.INPUT, in.path, in.lineNumber, in.charNumber, in.source)}
	})
}

type inputErrorType int

const (
	E_END_OF_FILE = iota
	E_ILLEGAL_ESCAPE
	E_UNEXPECTED_SYMBOL
	E_EXPECTED_CHAR_CLOSE
	E_EXPECTED_RAW_ID_CLOSE
	E_UNEXPECTED_CONTROL
	E_TRAILING_UNDERSCORE
	E_STRING_ONLY_ESCAPE
	E_BAD_UNICODE
	E_UNTERMINATED_COMMENT
	E_ILLEGAL_CHAR
)

var inputErrors = map[inputErrorType]_errorGenFn{
	E_END_OF_FILE:           inputErrorGen("unexpected end of file"),
	E_ILLEGAL_ESCAPE:        inputErrorGen("illegal escape sequence"),
	E_UNEXPECTED_SYMBOL:      inputErrorGen("unexpected symbol"),
	E_EXPECTED_CHAR_CLOSE:   inputErrorGen("expected end of character literal"),
	E_EXPECTED_RAW_ID_CLOSE: inputErrorGen("expected closing backtick"),
	E_UNEXPECTED_CONTROL:    inputErrorGen("unexpected control character"),
	E_TRAILING_UNDERSCORE:   inputErrorGen("trailing digit seperator"),
	E_STRING_ONLY_ESCAPE:    inputErrorGen("escape sequence is not allowed for char literals"),
	E_BAD_UNICODE:           inputErrorGen("malformed unicode code point escape sequence"),
	E_UNTERMINATED_COMMENT:  inputErrorGen("unterminated multi-line comment"),
	E_ILLEGAL_CHAR:          inputErrorGen("illegal character literal"),
}

var ID_REGEX = regexp.MustCompile("[A-Za-z][A-Za-z0-9_]*(')*")
var EXT_ID_REGEX = regexp.MustCompile(`[^0-9\{\}\[\]\(\)\s\x60[[:cntrl:]]@]+`)
// - ! # $ % ^ & * / ? . : > < = + ~
var SYMBOL_ID_REGEX_0 = regexp.MustCompile(`[!|#\$%\^\&\*/\?\.:><=\+~-]`)
var SYMBOL_ID_REGEX_1 = regexp.MustCompile(`[!|#\$%\^\&\/\?\.:><=\+~]`) // prev == -
var SYMBOL_ID_REGEX_2 = regexp.MustCompile(`[!|#\$%\^\&\*/\?\.:><=\+~]`) // prev == *
func nextRegexForSymb(previous byte) *regexp.Regexp {
	if previous == '-' {
		return SYMBOL_ID_REGEX_1
	} else if previous == '*' {
		return SYMBOL_ID_REGEX_2
	}
	return SYMBOL_ID_REGEX_0
}

var HEX_REGEX = regexp.MustCompile("0(x|X)([0-9A-Fa-f]_?)*[0-9A-Fa-f]+")
var OCT_REGEX = regexp.MustCompile("0(o|O)([0-7]_?)*[0-7]+")
var BIN_REGEX = regexp.MustCompile("0(b|B)([01]_?)*[01]+")
var INT_REGEX = regexp.MustCompile("([0-9]_?)*[0-9]+")
var FLOAT_REGEX = regexp.MustCompile(`[0-9]+(((\.[0-9]*)?(e|E)(\+|-)?([0-9]_?)*[0-9]+)|(\.[0-9]*))`)

var keywords = map[string]Token{
	"Int":     OtherToken{INT, 0, 0},
	"Bool":    OtherToken{BOOL, 0, 0},
	"Char":    OtherToken{CHAR, 0, 0},
	"Float":   OtherToken{FLOAT, 0, 0},
	"String":  OtherToken{STRING, 0, 0},
	"class":   OtherToken{CLASS, 0, 0},
	"where":   OtherToken{WHERE, 0, 0},
	"lazy":    OtherToken{LAZY, 0, 0},
	"let":     OtherToken{LET, 0, 0},
	"mut":     OtherToken{MUT, 0, 0},
	"const":   OtherToken{CONST, 0, 0},
	"mod":     OtherToken{MOD, 0, 0},
	"True":    ValueToken{value.Bool(true), 0, 0},
	"False":   ValueToken{value.Bool(false), 0, 0},
	"package": OtherToken{PACKAGE, 0, 0},
	"module":  OtherToken{MODULE, 0, 0},
}

func (in *Input) getIdOrKeyword2(forceId bool) Token {
	loc := ID_REGEX.FindStringIndex(in.getSourceSlice())
	if loc == nil || loc[0] != 0 {
		// this should be impossible
		err.PrintBug()
	}

	str := in.getSourceSliceN(loc[1])
	//fmt.Printf("string = `%s`\n", str)
	var key Token
	var found bool = false // keyword is never found when id is forced
	if !forceId {
		key, found = keywords[str]
	}
	start := in.charNumber
	in.charNumber += (loc[1] - 1)
	if !found { // found is false when token should be id
		idToken := IdToken{id: str, line: in.lineNumber, char: start}
		if unicode.IsUpper(rune(str[0])) {
			return TypeIdToken(idToken) // return type id token
		}
		// always returns here when id is forced
		return idToken
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
	input := in.getSourceSlice()
	loc := FLOAT_REGEX.FindStringIndex(input)
	if nil != loc && loc[0] == 0 {
		res, _ := strconv.ParseFloat(removeUnderscore(input[:loc[1]]), 64)
		start := in.charNumber
		in.charNumber += (loc[1] - 1)
		if in.hasTrailingUnderscore() {
			return inputErrors[E_TRAILING_UNDERSCORE](in)
		}
		return ValueToken{Value: value.Float(res), Line: in.lineNumber, Char: start}
	}

	loc = HEX_REGEX.FindStringIndex(input)
	if nil != loc && loc[0] == 0 {
		res, _ := strconv.ParseInt(removeUnderscore(input[:loc[1]]), 0, 64)
		start := in.charNumber
		in.charNumber += (loc[1] - 1)
		if in.hasTrailingUnderscore() {
			return inputErrors[E_TRAILING_UNDERSCORE](in)
		}
		return ValueToken{Value: value.Int(res), Line: in.lineNumber, Char: start}
	}

	loc = OCT_REGEX.FindStringIndex(input)
	if nil != loc && loc[0] == 0 {
		res, _ := strconv.ParseInt(removeUnderscore(input[:loc[1]]), 0, 64)
		start := in.charNumber
		in.charNumber += (loc[1] - 1)
		if in.hasTrailingUnderscore() {
			return inputErrors[E_TRAILING_UNDERSCORE](in)
		}
		return ValueToken{Value: value.Int(res), Line: in.lineNumber, Char: start}
	}

	loc = BIN_REGEX.FindStringIndex(input)
	if nil != loc && loc[0] == 0 {
		res, _ := strconv.ParseInt(removeUnderscore(input[:loc[1]]), 0, 64)
		start := in.charNumber
		in.charNumber += (loc[1] - 1)
		if in.hasTrailingUnderscore() {
			return inputErrors[E_TRAILING_UNDERSCORE](in)
		}
		return ValueToken{Value: value.Int(res), Line: in.lineNumber, Char: start}
	}

	loc = INT_REGEX.FindStringIndex(input)
	if nil != loc && loc[0] == 0 {
		res, _ := strconv.ParseInt(removeUnderscore(input[:loc[1]]), 10, 64)
		start := in.charNumber
		in.charNumber += (loc[1] - 1)
		if in.hasTrailingUnderscore() {
			return inputErrors[E_TRAILING_UNDERSCORE](in)
		}
		return ValueToken{Value: value.Int(res), Line: in.lineNumber, Char: start}
	}

	err.PrintBug()
	panic("")
}

func (in *Input) ungetChar() {
	in.unsetNextIndexes()
}


func (in *Input) getStructuralToken(c byte) (token Token, found bool) {
	found = true
	switch c {
	case '"':
		token = in.getString()
	case '\'':
		token = in.getChar()
	case '\\':
		token = in.tokenUnit(BACKSLASH)
	case '@':
		token = in.tokenUnit(AT)
		//token = AnotationToken(in.getIdOrKeyword2(true).(IdToken))
	case ';':
		token = in.tokenUnit(SEMI_COLON)
	case '\n':
		token = in.tokenUnit(NEW_LINE)
	case '(':
		token = in.tokenUnit(LPAREN)
	case ')':
		token = in.tokenUnit(RPAREN)
	case '[':
		token = in.tokenUnit(LBRACK)
	case ']':
		token = in.tokenUnit(RBRACK)
	case '{':
		token = in.tokenUnit(LCURL)
	case '}':
		token = in.tokenUnit(RCURL)
	case '_':
		token = in.tokenUnit(UNDERSCORE)
	case ',':
		token = in.tokenUnit(COMMA)
	case 0:
		token = in.tokenUnit(EOF)
	default:
		found = false
	}

	return
}

var builtinOps = map[string]TokenType{
	// operators that cannot be excluded (these are operators that are not part of prelude)
	"->": ARROW,
	"=>": FAT_ARROW,
	"=": EQUALS,
	"::": COLON_COLON,
	"|": BAR,
	".": DOT,
	"..": DOT_DOT,
	"?": QUESTION,

	// operators that can be excluded (these are operators that are part of prelude)
	"+": PLUS,
	"++": PLUS_PLUS,
	"-": MINUS,
	"*": STAR,
	"^": HAT,
	"/": SLASH,
	"==": EQUALS_EQUALS,
	":": COLON,
	"!": BANG,
	"!=": BANG_EQUALS,
	">": GREAT,
	">=": GREAT_EQUALS,
	"<": LESS,
	"<=": LESS_EQUALS,
	"&&": AMPER_AMPER,
	"||": BAR_BAR,
}

func (in *Input) getSymbolId() Token {
	start := in.charNumber // current place in stream
	var source = in.getSourceSlice() // get line sliced from char number - 1 to end of line/input
	idx := 0
	var prev byte = ' '
	sourceLen := len(source)
	for idx < sourceLen {
		if !nextRegexForSymb(prev).MatchString(source[idx:idx+1]) {
			break
		}
		prev = source[idx]
		idx++
	}

	if idx == 0 { 
		// illegal char
		return inputErrors[E_UNEXPECTED_SYMBOL](in) // generate error at start position
	}

	in.charNumber += (idx - 1) // update char number
	symb := source[:idx] // token regex matched
	symbLength := len(symb) // length of matched token
	
	if tokenType, found := builtinOps[symb]; found {
		if PreludeIncludes(tokenType) {
			return PreludeIdToken{tokenType: tokenType, id: symb, line: in.lineNumber, char: start}
		}
		return in.tokenUnitN(tokenType, symbLength) // return builtin token
	}
	return IdToken{id: symb, line: in.lineNumber, char: start} // return user defined token
}

func (in *Input) Next() Token {
	c, e := in.skipWhitespace()
	if nil != e {
		return *e
	}

	if unicode.IsDigit(rune(c)) {
		return in.getNumber()
	} else if nil != ID_REGEX.FindStringIndex(string(c)) {
		res := in.getIdOrKeyword()
		return res
	}

	if token, shouldReturn := in.getStructuralToken(c); shouldReturn {
		return token
	}

	return in.getSymbolId()
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

func (in *InputStream) ExcludePrelude() {
	in.excludePrelude = true
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
	return in.tokens[in.streamIndex-1].preludeExcludeTransform(in.excludePrelude)
}

// grows stream's buffer with respect to length of unread input; always grows buffer by at least one
// if there are no issues allocating memory for the new buffer
func (stream *InputStream) grow(in *Input) {
	// number of chars left to read
	remainingSourceLength := len(in.getSourceSlice())
	if in.lineNumber < len(in.source) {
		for _, s := range in.source[in.lineNumber:] {
			remainingSourceLength = remainingSourceLength + len(s)
		}
	}

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
		prev := stream.tokens[stream.streamLength-1].GetType()
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

	sourceLength := 0
	for _, s := range in.source {
		sourceLength = sourceLength + len(s)
	}

	stream := InitStream(in.path, int64(sourceLength))
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
