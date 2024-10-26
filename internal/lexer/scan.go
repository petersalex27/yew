package lexer

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/internal/common"
)

var escapeMap = map[rune]byte{
	'a':  '\a',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'v':  '\v',
	'\\': '\\',
}

func getEscape(r rune, escapeString bool) (c byte, ok bool) {
	c, ok = escapeMap[r]
	if ok {
		return c, true
	} else if ok = escapeString && r == '"'; ok {
		return byte(r), true
	} else if ok = !escapeString && r == '\''; ok {
		return byte(r), true
	}
	return 0, false
}

func (lex *Lexer) prevEnd() (previousLineEnd int) {
	if lex.Line > 1 {
		previousLineEnd = lex.EndPositions()[lex.Line-2]
	}
	return
}

func (lex *Lexer) charNumber() (charNum int, eof bool) {
	charNum = lex.Pos - (lex.prevEnd() - 1) //(lex.Pos) % (lex.PositionRanges[lex.Line-1]) + 1
	eof = lex.Pos >= len(lex.Source)
	return
}

// return current line starting from current char until the end of line
func (lex *Lexer) remainingLine() (line string, ok bool) {
	endPositions := lex.EndPositions()
	if lex.Line > len(endPositions) {
		return "", false
	}
	start, end := lex.Pos, endPositions[lex.Line-1]
	if start >= end {
		return "", false
	} else if lex.Source[end-1] == '\n' {
		// this will be true except for the last line
		end--
	}
	return string(lex.Source[start:end]), true
}

// true iff lexer is at end of source
func (lex *Lexer) eof() bool {
	return lex.Pos >= len(lex.Source)
}

// get current char
//
// panics if it fails to read a char
func (lex *Lexer) peek() byte {
	if c, ok := lex.currentSourceChar(); ok {
		return c
	} else if eof := lex.eof(); eof {
		return 0
	}
	// bug in code: this shouldn't happen
	panic("no character at current location")
}

func (lex *Lexer) peek2() (first, second byte) {
	if lex.Pos >= len(lex.Source) {
		return 0, 0
	} else if lex.Pos+1 >= len(lex.Source) {
		return lex.Source[lex.Pos], 0
	}
	return lex.Source[lex.Pos], lex.Source[lex.Pos+1]
}

// advances input by a single char and returns the new current char
func (lex *Lexer) nextChar() byte {
	if c := lex.peek(); c != 0 {
		lex.Pos++
		return c
	}
	return 0
}

// reads whitespace until next non-whitespace char, then advances input and returns non-whitespace
// char (and eof==true when lexer is at end of source)
//
// whitespace is just ' ' and '\t'
func (lex *Lexer) skipWhitespace() (isEof bool) {
	c := lex.peek()
	for ; c != 0; c = lex.peek() {
		switch c {
		case ' ':
			lex.nextChar()
		case '\t':
			lex.nextChar()
		default: // non whitespace
			return c == 0
		}
	}
	return c == 0
}

// prepares lexer for reading from source based on whether it's already read and whether the source is empty or not
func (lex *Lexer) fixLineChar() {
	if lex.Line == 0 {
		lex.Pos = 0
		lex.Line = min(1, lex.Lines()) // set line to zero when no source code, else set to one
	}
	lex.Pos, _ = lex.LinePos(lex.Line)
}

func (lex *Lexer) processId(line string) token.Token {
	if tok, isKeyword := lex.isKeyword(line); isKeyword {
		return tok.Make()
	}
	return token.Id.MakeValued(line)
}

// Translates an integer string into its canonical form. An integer string's canonical form is
// defined by the following transformations:
//   - All underscore separators are removed from the string
//   - All alphabetic characters are converted to lowercase
//   - All leading 0s are removed from the string except those necessary to represent a zero value
//     or in a prefix
func processNumberString(result string) string {
	result = strings.ReplaceAll(result, "_", "") // remove separators
	result = strings.ToLower(result) // convert prefixes and scientific notation 'E' or 'e' to lowercase
	if len(result) > 1 && result[1] != 'x' && result[1] != 'b' && result[1] != 'o' {
		result = strings.TrimLeft(result, "0") // remove leading 0s
		if len(result) == 0 || result[0] == '.' || result[0] == 'e' {
			result = "0" + result // add leading 0 (to empty string or floating point tail)
		}
	}
	return result
}

// Translates a string literal into its canonical form. A string literal's canonical form is defined
// by the following transformations:
//   - All escape sequences are replaced with their corresponding characters
func processStringOrCharString(result string, isString bool) (str string, i int) {
	var builder strings.Builder
	escaped := false
	for i, r := range result {
		if escaped {
			escaped = false
			if c, ok := getEscape(r, isString); ok {
				builder.WriteByte(c)
			} else {
				return builder.String(), i
			}
		} else if r == '\\' {
			escaped = true
		} else {
			builder.WriteRune(r)
		}
	}

	if escaped { // if escaped, len(result) must be at least 1
		return builder.String(), len(result) - 1
	}
	return builder.String(), len(result)
}

func processAnnotationString(annot string) string {
	return strings.TrimSpace(annot)
}

func (lex *Lexer) fromSymbol() token.Token {
	line, ok := lex.remainingLine()
	if !ok {
		return lex.error(UnexpectedEOF)
	}

	if result := common.SymbolId.Match(line); result != nil {
		lex.Pos += len(*result)
		return lex.processId(*result)
	}

	return lex.error(CannotProcessSymbol)
}

func (lex *Lexer) fromAlpha() token.Token {
	line, ok := lex.remainingLine()
	if !ok {
		panic("expected line content")
	}

	result := common.AlphanumericId.Match(line)
	if result == nil {
		panic("expected alphanumeric id")
	}

	lex.Pos += len(*result)
	return lex.processId(*result)
}

func (lex *Lexer) fromDigit() token.Token {
	line, ok := lex.remainingLine()
	if !ok {
		panic("expected line content")
	}

	// first, try floating point number
	if result := common.FloatLiteral.Match(line); result != nil {
		lex.Pos += len(*result)
		return token.FloatValue.MakeValued(processNumberString(*result))
	}

	var result *string
	isBase10 := false

	// identify kind of integer based on prefix
	first, second := lex.peek2()
	if first == '0' {
		switch second {
		case 'x', 'X':
			second = 'x'
			result = common.HexLiteral.Match(line)
		case 'b', 'B':
			second = 'b'
			result = common.BinLiteral.Match(line)
		case 'o', 'O':
			second = 'o'
			result = common.OctLiteral.Match(line)
		default:
			isBase10 = true
		}
	} else {
		isBase10 = true
	}

	if isBase10 {
		result = common.IntLiteral.Match(line)
	}

	if result == nil {
		switch second {
		case 'x':
			return lex.error(IllegalHexLiteral)
		case 'b':
			return lex.error(IllegalBinLiteral)
		case 'o':
			return lex.error(IllegalOctLiteral)
		default:
			return lex.error(IllegalIntLiteral)
		}
	}

	lex.Pos += len(*result)
	return token.IntValue.MakeValued(processNumberString(*result))
}

func (lex *Lexer) fromLeftParen() token.Token {
	line, ok := lex.remainingLine()
	if !ok {
		panic("expected line content")
	}

	if _, c := lex.peek2(); c == ')' {
		lex.Pos += 2
		return token.EmptyParenEnclosure.Make()
	}

	if res := common.MethodId.Match(line); res != nil {
		lex.Pos += len(*res)
		return token.MethodSymbol.MakeValued((*res)[2 : len(*res)-1]) // remove leading "(." and trailing ")"
	} else if res = common.InfixId.Match(line); res != nil {
		lex.Pos += len(*res)
		return token.Infix.MakeValued((*res)[1 : len(*res)-1]) // remove leading "(" and trailing ")"
	}

	lex.Pos++
	return token.LeftParen.Make()
}

func (lex *Lexer) fromLeftBracket() token.Token {
	if _, c := lex.peek2(); c == '@' {
		lex.Pos += 2
		return token.LeftBracketAt.Make()
	} else if c == ']' {
		lex.Pos += 2
		return token.EmptyBracketEnclosure.Make()
	}

	lex.Pos++
	return token.LeftBracket.Make()
}

var annotRegex = regexp.MustCompile(`--\s*@`)

func (lex *Lexer) fromMinus() token.Token {
	_, second := lex.peek2()
	if second == '-' {
		line, ok := lex.remainingLine()
		if !ok {
			panic("expected line content")
		}

		loc := annotRegex.FindStringIndex(line)
		if loc != nil && loc[0] == 0 {
			lex.Pos += len(line)
			annot := line[loc[1]:]
			annot = processAnnotationString(annot)
			return token.FlatAnnotation.MakeValued(annot)
		}

		// comment, so skip and return either end of tokens or newline
		lex.Pos += len(line)
		return lex.fromEndOfLine()
	}

	return lex.fromSymbol()
}

func (lex *Lexer) fromApostrophe() token.Token {
	line, ok := lex.remainingLine()
	if !ok {
		panic("expected line content")
	}
	result := common.CharLiteral.Match(line)
	if result == nil {
		return lex.error(IllegalCharLiteral)
	}

	literal := (*result)[1 : len(*result)-1]   // remove leading and trailing "'"
	
	var processed int
	literal, processed = processStringOrCharString(literal, false)
	if processed != len(*result)-2 {
		lex.Pos += processed + 1
		return lex.error(IllegalCharLiteral)
	}

	lex.Pos += len(*result)
	return token.CharValue.MakeValued(literal)
}

func (lex *Lexer) stringError() token.Token {
	line, ok := lex.remainingLine()
	if !ok {
		panic("expected line content")
	}

	result := common.ILLEGAL_UnterminatedInterpretedStringLiteral.Match(line)
	if result != nil {
		lex.Pos += len(*result)
		return lex.error(UnterminatedString)
	}

	result = common.ILLEGAL_BadEscapeSequence.Match(line)
	if result != nil {
		lex.Pos += len(*result)
		return lex.error(IllegalEscapeSequence)
	}

	return lex.error(IllegalStringLiteral)
}

func (lex *Lexer) fromQuote() token.Token {
	line, ok := lex.remainingLine()
	if !ok {
		panic("expected line content")
	}

	result := common.InterpretedStringLiteral.Match(line)
	if result == nil {
		return lex.stringError() // needs additional processing to identify error
	}

	literal := (*result)[1 : len(*result)-1] // remove leading and trailing '"'
	var processed int
	literal, processed = processStringOrCharString(literal, true)
	if processed != len(*result)-2 {
		// +1 to account for leading '"' (trailing '"' is not accounted for because less than the length
		// of the string's chars were processed)
		lex.Pos += processed + 1
		return lex.error(IllegalEscapeSequence)
	} else {
		lex.Pos += len(*result)
	}
	return token.StringValue.MakeValued(literal)
}

func (lex *Lexer) fromBacktick() token.Token {
	var b strings.Builder
	_ = lex.nextChar() // move past first backtick
	for !lex.eof() {
		if c := lex.peek(); c == '`' {
			break
		} else {
			b.WriteByte(c)
			_ = lex.nextChar()
		}
	}

	if lex.eof() {
		return lex.error(UnexpectedEOF)
	}

	_ = lex.nextChar() // move past final backtick
	return token.RawStringValue.MakeValued(b.String())
}

func (lex *Lexer) fromQuestionMark() token.Token {
	_, c := lex.peek2()
	if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z') {
		return lex.fromSymbol()
	}

	line, ok := lex.remainingLine()
	if !ok {
		panic("expected line content")
	}

	result := common.HoleId.Match(line)
	if result == nil {
		return lex.error(UnexpectedSymbol)
	}

	lex.Pos += len(*result)
	return token.Hole.MakeValued(*result)
}

func (lex *Lexer) fromEndOfLine() token.Token {
	if lex.peek() == '\n' {
		lex.Pos++
		lex.Line++
		return token.Newline.Make()
	} else if lex.eof() {
		return token.EndOfTokens.Make()
	}

	panic("expected end of line")
}

func (lex *Lexer) scan() (tok token.Token) {
	if lex.Line == 0 {
		lex.fixLineChar()
	}

	if eof := lex.skipWhitespace(); eof {
		tok := token.EndOfTokens.Make()
		tok.Start, tok.End = lex.Pos, lex.Pos
		return tok
	}

	before := lex.Pos

	c, next := lex.peek2()

	// need to handle this case separately b/c comments are
	// skipped and a newline is returned instead. Because of
	// this, the start position might not be `before` like
	// all other tokens. Additionally, there's a chance
	// the token is an annotation, so the start position may
	// still be `before`.
	if c == '-' && next == '-' {
		tok := lex.fromMinus()
		if tok.Typ == token.Newline || tok.Typ == token.EndOfTokens {
			tok.Start = lex.Pos - 1 // Pos will always be > 1
		} else {
			tok.Start = before
		}
		tok.End = lex.Pos
		return tok
	}

	// set position for whatever token is read
	defer func(start int) {
		tok.Start, tok.End = start, lex.Pos
	}(before)

	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' {
		return lex.fromAlpha()
	} else if '0' <= c && c <= '9' {
		return lex.fromDigit()
	}

	switch c {
	case '(':
		return lex.fromLeftParen()
	case ')':
		lex.nextChar()
		return token.RightParen.Make()
	case '[':
		return lex.fromLeftBracket()
	case ']':
		lex.nextChar()
		return token.RightBracket.Make()
	case '{':
		lex.nextChar()
		return token.LeftBrace.Make()
	case '}':
		lex.nextChar()
		return token.RightBrace.Make()
	case ',':
		lex.nextChar()
		return token.Comma.Make()
	case '-':
		return lex.fromMinus()
	case '_':
		lex.nextChar()
		return token.Underscore.Make()
	case '\'':
		return lex.fromApostrophe()
	case '"':
		return lex.fromQuote()
	case '`':
		return lex.fromBacktick()
	case '?':
		return lex.fromQuestionMark()
	case '\n':
		return lex.fromEndOfLine()
	}

	return lex.fromSymbol()
}

func (lex *Lexer) Scan() (tok api.Token) {
	return lex.scan()
}

// NOTE: panics if not in repl mode
func (lex *Lexer) Command() string {
	if !token.InReplMode() {
		panic("illegal call: tried to lex a command in non-repl mode")
	}
	eof := lex.skipWhitespace()
	if eof {
		return ""
	}

	b := &strings.Builder{}
	c := lex.peek()
	if c != ':' {
		return ""
	}
	b.WriteByte(c)

	for {
		c = lex.peek()
		if unicode.IsLetter(rune(c)) {
			b.WriteByte(c)
			lex.nextChar()
		} else {
			break
		}
	}

	// validate command
	return commands[b.String()].CommandLiteral()
}
