package lexer

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/internal/common"
)

type symbolClass byte

const (
	symbol_class symbolClass = iota
	number_class
	raw_string_class
	string_class
	identifier_class
	infix_class
	method_class
	underscore_class
	comment_class
	char_class
	hole_class
	end_class
	error_class
)

//var freeSymbolRegex = regexp.MustCompile(freeSymbolRegexClassRaw)

func isSymbol(c byte) bool {
	s := string([]byte{c})
	return common.Is_symbolCase(s) || common.IsStandaloneSymbol(s)
}

// classifies the next token based on the already seen '-' and whatever the next char is
//
// next is ...
//   - '-' or '*': comment_class
//   - otherwise, symbol_class
//
// error will be handled when tokenizing the symbol if its not a symbol
func (lex *Lexer) classifyMinus() (class symbolClass) {
	c, eof := lex.peek()

	classifyAsSymbol := eof || !(c == '-' || c == '*')
	if classifyAsSymbol {
		return symbol_class
	}

	return comment_class
}

// regex for infix identifier w/o leading '('
var infixRegex_wo_lparen = regexp.MustCompile(`([a-zA-Z][a-zA-Z0-9']*|[!@#$^&*~<>?/:|\-+=\\]+)\)`)

func (lex *Lexer) classifySymbol(r rune) (class symbolClass) {
	if r != '(' {
		return symbol_class
	}
	// check if infix identifier of some form
	line, ok := lex.remainingLine()
	if !ok {
		return symbol_class
	}

	// class to return if the infix regex matches
	matchedClass := infix_class

	// set matchedClass to method_class?
	if len(line) > 1 && line[1] == '.' {
		line = line[1:]
		matchedClass = method_class
	}

	loc := infixRegex_wo_lparen.FindStringIndex(line)
	if loc != nil && loc[0] == 0 {
		return matchedClass
	}

	return symbol_class
}

var letterRegex = regexp.MustCompile(`[a-zA-Z]`)

// Determines the class of some input section based on some byte `c` of the
// input. Unless there's a good reason to do otherwise, `c` is the first
// character of that input section.
func (lex *Lexer) classify(c byte) (class symbolClass, errorToken token.Token) {
	r := rune(c)
	if unicode.IsLetter(r) {
		class = identifier_class
	} else if unicode.IsDigit(r) {
		class = number_class
	} else if c == '\'' {
		class = char_class
	} else if c == '`' {
		class = raw_string_class
	} else if c == '"' {
		class = string_class
	} else if c == '_' {
		class = underscore_class
	} else if c == '-' {
		class = lex.classifyMinus()
	} else if c2, _ := lex.peek(); c == '?' && common.MatchRegex(letterRegex, string(c2)) != nil {
		class = hole_class
	} else if isSymbol(c) {
		class = lex.classifySymbol(r)
	} else {
		class = error_class
		errorToken = lex.error(UnexpectedSymbol)
	}
	return
}

// creates and pushes token for single-line comment with Value=`lineAfterDashes`
func (lex *Lexer) comment(lineAfterDashes string) token.Token {
	offs := len(lineAfterDashes)
	// remove newline (comment right before eof won't have newline)
	if offs > 0 && lineAfterDashes[offs-1] == '\n' {
		offs-- // don't take newline
	}
	lex.Pos += offs
	tok := token.Comment.MakeValued(lineAfterDashes)
	annot := strings.TrimLeft(lineAfterDashes, " \t")
	if len(annot) > 0 && annot[0] == '@' {
		tok.Typ = token.FlatAnnotation
		annot = strings.TrimLeft(annot[1:], " \t")
		if len(annot) == 0 {
			return lex.error(ExpectedAnnotationId)
		}
		tok.Value = annot
	}

	return lex.output(tok)
}

func (lex *Lexer) analyzeComment() (token token.Token) {
	lex.SavedChar.Push(lex.Pos)

	var c byte
	_, eof := lex.nextChar() // remove initial '-'
	if eof {
		panic("bug: source not validated before calling analyzeComment")
	}

	// remove second thing ('-' or '*'), the value of `c` will determine the
	// branch to take in the condition below the next one
	c, eof = lex.nextChar()
	if eof {
		return lex.error(UnexpectedEOF)
	}

	// given a line
	//		-- abc ..
	// or
	//		-* abc ..
	// `line` is
	//		line = " abc .."
	var line string
	line, _ = lex.remainingLine()
	if c == '-' { // single line comment
		return lex.comment(line)
	}
	panic("bug in analyzeComment: else branch reached")
}

func isNumEndCharValid(line string, numEnd int) bool {
	if len(line) <= numEnd {
		return true
	}

	return line[numEnd] != '_' && (line[numEnd] == '_' || line[numEnd] == '\t' || isSymbol(line[numEnd]))
}

// removes `strip` from `s` and returns result
func stripChar(s string, strip byte) string {
	var builder strings.Builder
	c := rune(strip)
	for _, r := range s {
		if r != c {
			builder.WriteByte(byte(r))
		}
	}
	return builder.String()
}

// read non-base ten integer token: hexadecimal, octal, or binary
func analyzeNonBase10(num, line string) (tok token.Token, numChars int, errorMessage string) {
	numChars, errorMessage = len(num), ""
	if !isNumEndCharValid(line, len(num)) {
		errorMessage = UnexpectedSymbol
	} else {
		num = stripChar(num, '_')
		tok = token.IntValue.MakeValued(num)
	}
	return
}

// returns true iff the char of string `s` at index `i` is 'e' or 'E'
func isE(s string, i int) bool {
	return s[i] == 'e' || s[i] == 'E'
}

// returns true iff the char of string `s` at index `i` has a sign (i.e., has '+' or '-')
func isSign(s string, i int) bool {
	return s[i] == '+' || s[i] == '-'
}

// assumes numChars is the correct value and that it corresponds to the length
// of the token
func returnInt(num string, numChars int) (token.Token, int, string) {
	num = stripChar(num, '_')
	// remove leading 0s so translation to llvm ir is not confused thinking it's octal
	num = strings.TrimLeft(num, "0")
	if num == "" { // number was all 0s?
		num = "0" // yes => set to a single 0
	}
	return token.IntValue.MakeValued(num), numChars, ""
}

// parses a floating point number that contains a '.' as a string
//
// returns
//   - num: (partially, in the case of there being an 'e'/'E') parsed number `num`
//   - numChars: the number of characters parsed
//   - hasE: whether parsed number is followed immediately by 'e' or 'E'
//   - errorMessage: an empty string if the function is successful, otherwise an error message
func analyzeDotNum(numOrigin, line string, numCharsOrigin int, hasEOrigin bool) (num string, numChars int, hasE bool, errorMessage string) {
	num, numChars, hasE = numOrigin, numCharsOrigin, hasEOrigin // init

	numChars = numChars + 1
	if len(line) <= numChars { // <integer>.EOL
		errorMessage = UnexpectedSymbol
		return
	}

	frac := common.IntLiteral.Match(line[numChars:])
	if frac == nil { // <integer>.<non-integer>
		errorMessage = UnexpectedSymbol
		return
	}

	numChars = numChars + len(*frac)
	num = num + "." + *frac

	if len(line) <= numChars {
		return
	}

	hasE = isE(line, numChars)
	if !hasE && !isNumEndCharValid(line, numChars) { // <integer>.<integer><illegal-char>
		errorMessage = UnexpectedSymbol
	}
	return
}

// read sign if one exists and return (possibly) empty sign and new total
// number of chars read
func analyzePossibleSign(line string, numCharsOrigin int) (sign string, numChars int) {
	// init
	numChars = numCharsOrigin

	// if signed, read sign and return it
	signed := isSign(line, numChars)
	sign = ""
	if signed {
		sign = string(line[numChars])
		numChars = numChars + 1
	}

	return sign, numChars
}

// ASSUMPTION: line[numChars] == 'e' or 'E'
//
// reads number from input at exponent marker (i.e., 'e' or 'E') to end of number
func analyzeExponentNum(numOrigin, line string, numCharsOrigin int) (num string, numChars int, errorMessage string) {
	num, numChars = numOrigin, numCharsOrigin // init

	e := line[numChars] // 'e' or 'E'
	numChars = numChars + 1
	if len(line) <= numChars { // <float>eEOL
		errorMessage = UnexpectedSymbol
		return
	}

	var sign string
	sign, numChars = analyzePossibleSign(line, numChars)

	if len(line) <= numChars { // <float>e<sign>EOL
		errorMessage = UnexpectedSymbol
		return
	}

	// read integer value that follows 'e'/'E'
	frac := common.IntLiteral.Match(line[numChars:])
	if frac == nil { // <float>e[sign]<illegal-char>
		errorMessage = UnexpectedSymbol
		return
	}

	numChars = numChars + len(*frac)
	if !isNumEndCharValid(line, numChars) { // <float>e[sign]<integer><illegal-char>
		errorMessage = UnexpectedSymbol
		return
	}

	// build value as string
	num = num + string(e) + sign + *frac
	return
}

// return a number token. Could be either floating point number or integer
func maybeFractional(num, line string) (tok token.Token, numChars int, errorMessage string) {
	numChars, errorMessage = len(num), ""

	if len(line) <= numChars { // just an integer at the end of the line
		return returnInt(num, numChars)
	}

	// because of above branch, line[numChars] must exist
	hasE := isE(line, numChars)
	hasDot := line[numChars] == '.'

	if !hasDot && !hasE {
		return returnInt(num, numChars)
	}

	// dotNum must be handled first to account for numbers like '123.123e123'
	if hasDot {
		num, numChars, hasE, errorMessage = analyzeDotNum(num, line, numChars, hasE)
	}

	// read 'e' or 'E' and exponent
	if hasE {
		num, numChars, errorMessage = analyzeExponentNum(num, line, numChars)
	}

	num = stripChar(num, '_')
	tok = token.FloatValue.MakeValued(num)
	return
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
	}
	return string(lex.Source[start:end]), true
}

// read number from input
func (lex *Lexer) number() token.Token {
	lex.SavedChar.Push(lex.Pos)
	line, ok := lex.remainingLine()
	if !ok {
		panic("bug: function called without verifying readable source exists")
	}

	var tok token.Token // token result
	var numChars int    // total number of chars read
	var errorMessage string = ""

	// 0x, 0b, and 0o must be checked first, else the lexer might falsely think
	// '0' is the number
	if num := common.HexLiteral.Match(line); num != nil { // hex
		tok, numChars, errorMessage = analyzeNonBase10(*num, line)
	} else if num := common.OctLiteral.Match(line); num != nil { // oct
		tok, numChars, errorMessage = analyzeNonBase10(*num, line)
	} else if num := common.BinLiteral.Match(line); num != nil { // bin
		tok, numChars, errorMessage = analyzeNonBase10(*num, line)
	} else if num := common.IntLiteral.Match(line); num != nil { // int or float
		tok, numChars, errorMessage = maybeFractional(*num, line)
	} else {
		numChars = 0
		errorMessage = UnexpectedSymbol
	}

	lex.Pos += numChars
	if errorMessage != "" {
		return lex.error(errorMessage)
	}

	return lex.output(tok)
}

// map of standalone symbols--these cannot be used within other identifiers
//
// in addition to these, there is also the standalone symbol '[@'
var standaloneMap = map[byte]token.Type{
	'(': token.LeftParen,
	')': token.RightParen,
	'[': token.LeftBracket,
	']': token.RightBracket,
	'{': token.LeftBrace,
	'}': token.RightBrace,
	',': token.Comma,
}

func (lex *Lexer) tryStandalone(first byte) (tok token.Token, ok bool) {
	var tokenType token.Type
	tokenType, ok = standaloneMap[first]
	if !ok {
		return tok, false
	}

	lex.SavedChar.Push(lex.Pos)
	lex.nextChar()
	if c, _ := lex.peek(); first == '[' && c == '@' {
		lex.nextChar()
		tokenType = token.LeftBracketAt
	} else if first == '(' && c == ')' {
		lex.nextChar()
		tokenType = token.EmptyParenEnclosure
	} else if first == '[' && c == ']' {
		lex.nextChar()
		tokenType = token.EmptyBracketEnclosure
	}
	return tokenType.Make(), true
}

func (lex *Lexer) standalone(fallback func(*Lexer) token.Token) token.Token {
	c, eof := lex.peek()
	if eof {
		panic("bug: function called without source validation")
	}

	if tok, ok := lex.tryStandalone(c); ok {
		return lex.output(tok)
	}
	return fallback(lex)
}

func (lex *Lexer) hole() token.Token {
	if c, _ := lex.peek(); c != '?' {
		panic("next character in input should've been validated before calling")
	}

	_, _ = lex.nextChar() // move past '?'
	tok := lex.identifier()
	if tok.Error() != nil {
		return tok
	}

	// update token
	tok.Start-- // start should be one backward to account for leading '?'
	tok.Typ = token.Hole
	tok.Value = "?" + tok.Value
	return tok
}

// reads and creates a token for some sequence of symbols
//
// this includes the standalone symbols, holes, identifiers, infixes,
func (lex *Lexer) symbol() token.Token {
	if c, _ := lex.peek(); c == '?' {
		return lex.hole()
	}

	return lex.standalone((*Lexer).identifier)
}

func (lex *Lexer) infix() token.Token {
	lex.SavedChar.Push(lex.Pos)
	
	line, ok := lex.remainingLine()
	if !ok {
		panic("bug: function called without verifying readable source exists")
	}

	var pStr *string
	if pStr = common.MethodId.Match(line); pStr != nil {
		ln := len(*pStr)
		lex.Pos += len(*pStr)
		tok := token.MethodSymbol.MakeValued((*pStr)[2:ln-1]) // remove leading "(." and trailing ")"
		return lex.output(tok)
	} else if pStr = common.InfixId.Match(line); pStr != nil {
		ln := len(*pStr)
		lex.Pos += len(*pStr)
		tok := token.Infix.MakeValued((*pStr)[1:ln-1]) // remove leading "(" and trailing ")"
		return lex.output(tok)
	} else {
		panic("bug: source was not validated before calling infix")
	}
}

var escapeMap = map[rune]byte{
	'n': '\n', 't': '\t', 'r': '\r', 'v': '\v', 'b': '\b', 'a': '\a', 'f': '\f', '\\': '\\',
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

// read escape sequence
func readEscapable(line string, end byte) (string, int, bool) {
	index := 0
	escaped := false
	for _, c := range line {
		if escaped {
			escaped = false
		} else if byte(c) == end {
			return line[:index], index, true
		} else if byte(c) == '\\' {
			escaped = true
		}
		index = index + 1
	}
	// `end` not found
	return "", index, false
}

func writeEscape(builder *strings.Builder, next bool, r rune, escapeString bool) (again, ok bool) {
	var c byte
	if next {
		again = false
		c, ok = getEscape(r, escapeString)
	} else if r == '\\' {
		again, ok = true, true
	} else {
		again, ok, c = false, true, byte(r)
	}

	if ok && !again {
		builder.WriteByte(c)
	}
	return again, ok
}

func updateEscape(s string, escapeString bool) (string, bool, int) {
	var builder strings.Builder
	var next, ok bool = false, true
	out := len(s) - 1
	for i, r := range s {
		next, ok = writeEscape(&builder, next, r, escapeString)
		if !ok {
			return "", false, i
		}
	}

	return builder.String(), true, out
}

func (lex *Lexer) char() token.Token {
	lex.SavedChar.Push(lex.Pos)

	c, eof := lex.nextChar() // should be leading '
	if ok := !(eof || c != '\''); !ok {
		if c != '\'' {
			panic("bug: called before validating next token is a character token")
		}
		return lex.error(UnexpectedEOF)
	}

	line, ok := lex.remainingLine()
	if !ok {
		return lex.error(ExpectedCharLiteral)
	}

	var escaped string
	var length int
	escaped, length, ok = readEscapable(line, '\'')
	if !ok {
		return lex.error(IllegalEscapeSequence)
	}

	lex.Pos += length
	if c, eof = lex.nextChar(); c != '\'' || eof { // remove closing `'`
		if eof {
			return lex.error(UnexpectedEOF)
		}
		return lex.error(IllegalCharLiteral)
	}

	var index int
	escaped, ok, index = updateEscape(escaped, false)
	if !ok {
		lex.Pos += index + 1
		return lex.error(IllegalEscapeSequence)
	}
	if ok = len(escaped) == 1; !ok {
		return lex.error(IllegalCharLiteral)
	}

	tok := token.CharValue.MakeValued(escaped)
	return lex.output(tok)
}

// This counts number of contiguous `c`s at the end of `s`.
//
// Examples:
//
//	countTrailing("employee", 'e') = 2
//	countTrailing("employee", 'y') = 0
//	countTrailing("", 'w') = 0
func countTrailing(s string, c byte) uint {
	trailing := uint(0)
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] != c {
			break
		}

		trailing++
	}
	return trailing
}

// returns true if `s` ends in an escaped quote
func hasFinalQuoteEscape(s string) bool {
	length := len(s)
	possibleEscapedQuote := length >= 2 && s[length-2:] == `\"`
	if !possibleEscapedQuote {
		return false
	}

	// remove final '"' so trailing backslashes can be counted
	unquoted := s[:length-1]
	// if number of escapes is 2n for some n, then there are n escaped '\\'; if there are 2n+1
	// '\\', then there are n escaped '\\' and a final escaped '"'
	isEscapedQuote := (countTrailing(unquoted, '\\') % 2) != 0

	return isEscapedQuote
}

func (lex *Lexer) getStringContent() (content string, ok bool) {
	var section string
	content = ""

	// reads string (and accounts for escaped '"')
	again := true
	for again {
		section, ok = lex.readThrough('"')
		if !ok {
			return
		}

		content = content + section
		again = hasFinalQuoteEscape(section)
	}

	charsRead := len(content)
	if charsRead > 0 {
		content = content[:charsRead-1] // remove trailing '"'
	}
	return
}

func (lex *Lexer) rawStringLiteral() token.Token {
	lex.SavedChar.Push(lex.Pos)
	c, eof := lex.nextChar() // should be first '`'
	if eof || c != '`' {
		panic("bug: source was not validated")
	}

	// writes everything except the final '`'--including control characters
	var b *strings.Builder = &strings.Builder{}
	for {
		c, eof = lex.nextChar()
		if eof {
			return lex.error(UnexpectedEOF)
		} else if c == '`' {
			break
		}
		b.WriteByte(c)
	}

	tok := token.RawStringValue.MakeValued(b.String())
	return lex.output(tok)
}

func (lex *Lexer) stringLiteral() token.Token {
	lex.SavedChar.Push(lex.Pos)
	c, eof := lex.nextChar() // should be first quotation mark
	// check for leading quotation mark
	if eof || c != '"' {
		panic("bug: source was not verified")
	}

	content, ok := lex.getStringContent()
	if !ok {
		return lex.error(IllegalStringLiteral)
	}

	updatedContent, ok, _ := updateEscape(content, true)
	if !ok {
		return lex.error(IllegalEscapeSequence)
	}

	var tok token.Token
	// check if string can be used as an import path--if so, set the type to
	// the more specific `ImportPath` instead of `StringValue`
	if common.IsImportPathString(updatedContent) {
		tok = token.ImportPath.MakeValued(updatedContent)
	} else {
		tok = token.StringValue.MakeValued(updatedContent)
	}

	return lex.output(tok)
}

// finds a substring of `line` starting at index 0 to an index > 0 that is an identifier
//
// if a substring with the above requirements cannot be found, then a non-empty `errorMessage` is
// returned
func matchId(line string) (matched string, errorMessage string, illegalArgument bool) {
	if len(line) < 1 {
		return "", "", true
	}

	var res *string = nil
	if res = common.AlphanumericId.Match(line); res != nil {
		return *res, "", false
	} else if res = common.SymbolId.Match(line); res != nil {
		return *res, "", false
	}

	return "", UnexpectedSymbol, false
}

// returns type of identifier
//
// return value `ty` is ...
//   - token.Id
//   - keyword token type
//   - token.Error (if `id` contains an underscore)
func (lex *Lexer) getIdType(id string) (ty token.Type, errorMessage string) {
	ty = token.Id

	if strings.ContainsRune(id, '_') {
		ty, errorMessage = token.Error, UnexpectedSymbol
	} else if key, yes := lex.isKeyword(id); yes {
		ty = key
	}

	return ty, errorMessage
}

func (lex *Lexer) matchIdentifier(line string) (id string, ty token.Type, errorMessage string) {
	var illegalArgument bool
	id, errorMessage, illegalArgument = matchId(line)
	if illegalArgument {
		panic("bug: illegal argument, empty string for argument `line`")
	} else if errorMessage != "" {
		return "", token.Error, errorMessage
	}

	ty, errorMessage = lex.getIdType(id)
	return id, ty, errorMessage
}

func (lex *Lexer) identifier() token.Token {
	lex.SavedChar.Push(lex.Pos)
	line, ok := lex.remainingLine()
	if !ok {
		panic("bug: function called without verifying readable source exists")
	}

	var ty token.Type
	errorMessage := ""

	var id string
	id, ty, errorMessage = lex.matchIdentifier(line)

	if errorMessage != "" {
		return lex.error(errorMessage)
	}

	// add token
	offs := len(id) // this works b/c id is copied from `line` (not modified)
	lex.Pos += offs
	tok := ty.MakeValued(id)
	return lex.output(tok)
}

func (lex *Lexer) underscore() token.Token {
	lex.SavedChar.Push(lex.Pos)
	_, eof := lex.nextChar()
	if eof {
		panic("bug: source was not validated")
	}

	tok := token.Underscore.Make()
	return lex.output(tok)
}

func (class symbolClass) analyze(lex *Lexer) (tok token.Token) {
	switch class {
	case number_class:
		return lex.number()
	case identifier_class, symbol_class, hole_class:
		return lex.symbol()
	case infix_class, method_class:
		return lex.infix()
	case char_class:
		return lex.char()
	case raw_string_class:
		return lex.rawStringLiteral()
	case string_class:
		return lex.stringLiteral()
	case underscore_class:
		return lex.underscore()
	case comment_class:
		tok = lex.analyzeComment()
		if tok.Typ == token.Comment && !lex.keepComments {
			tok = lex.analyze()
		}
		return tok
	}

	return lex.error(UnexpectedSymbol)
}

// true iff lexer is at end of source
func (lex *Lexer) eof() bool {
	return lex.Pos >= len(lex.Source)
}

// get current char
//
// panics if it fails to read a char
func (lex *Lexer) peek() (c byte, eof bool) {
	var ok bool
	if c, ok = lex.currentSourceChar(); ok {
		return
	}

	if eof = lex.eof(); eof {
		return 0, true
	}
	// bug in code: this shouldn't happen
	panic("no character at current location")
}

// advances input by a single char and returns the new current char
func (lex *Lexer) nextChar() (c byte, eof bool) {
	if c, eof = lex.peek(); !eof {
		lex.Pos++
	}
	return
}

// unadvance input by a single char and returns the new current char
func (lex *Lexer) ungetChar() (c byte) {
	char, eof := lex.charNumber()
	if char > 1 || eof {
		lex.Pos--
	} else if lex.Line > 1 { // char is 1
		lex.Line--
		lex.Pos--
	} else {
		panic("bug: cannot move input to a position before its start")
	}

	c, _ = lex.currentSourceChar()
	return
}

// reads whitespace until next non-whitespace char, then advances input and returns non-whitespace
// char (and eof==true when lexer is at end of source)
//
// whitespace is just ' ' and '\t'
func (lex *Lexer) skipWhitespace() (isEof bool) {
	// technically, condition isn't needed b/c if eof==true, then c==0 which is the default case
	var c byte

	c, isEof = lex.peek()
	for ; !isEof; c, isEof = lex.peek() {
		switch c {
		case ' ':
			lex.nextChar()
		case '\t':
			lex.nextChar()
		default: // non whitespace
			return isEof
		}
	}
	return isEof
}

func (lex *Lexer) analyze() token.Token {
	eof := lex.skipWhitespace()

	if eof {
		return lex.output(token.EndOfTokens.Make())
	}

	lex.SavedChar.Push(lex.Pos)
	c, _ := lex.nextChar()
	if c == '\n' {
		tok := token.Newline.MakeValued("\n")
		return lex.output(tok)
	}
	// use char to determine what class new token will belong to
	class, errorToken := lex.classify(c)
	if class == error_class {
		return errorToken
	}
	lex.ungetChar() // unget char gotten from lex.nextChar

	// use class information to get token
	return class.analyze(lex)
}

// prepares lexer for reading from source based on whether it's already read and whether the source is empty or not
func (lex *Lexer) fixLineChar() {
	if lex.Line == 0 {
		lex.Pos = 0
		lex.Line = min(1, lex.Lines()) // set line to zero when no source code, else set to one
	}
	lex.Pos, _ = lex.LinePos(lex.Line)
}

func (lex *Lexer) Scan() (tok api.Token) {
	if lex.Line == 0 {
		lex.fixLineChar()
	}

	return lex.analyze()
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
	c, _ := lex.nextChar()
	if c != ':' {
		lex.ungetChar()
		return ""
	}
	b.WriteByte(c)

	for {
		c, _ = lex.peek()
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
