package lexer

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/petersalex27/yew/common"
	"github.com/petersalex27/yew/token"
)

type symbolClass byte

const (
	symbol_class symbolClass = iota
	number_class
	string_class
	identifier_class
	underscore_class
	comment_class
	char_class
	end_class
	error_class
)

// ( ) [ ] { } ! @ # $ % ^ & * ~ , < > . ? / ; : | - + = \ `
const symbolRegexClassRaw string = `[\(\)\[\]\{\}!@#\$%\^\&\*~,<>\.\?/;:\|\-\+=\\` + "`]"

var symbolRegex = regexp.MustCompile(symbolRegexClassRaw)

const freeSymbolRegexClassRaw string = `[!#\$%\^\&\*~,<>\?/:\|\-\+=\\` + "`]"

//var freeSymbolRegex = regexp.MustCompile(freeSymbolRegexClassRaw)

func isSymbol(c byte) bool {
	return symbolRegex.Match([]byte{c})
}

// (re)classifies underscore_class to some (possibly the same) class
//
// if a lexical error occurs, `error_class` is returned and an error is added to lex's `messages`
func (lex *Lexer) reclassifyUnderscore() (class symbolClass) {
	c, eof := lex.peek()
	r := rune(c)
	classifyAsSymbol := !eof && (unicode.IsLetter(r) || unicode.IsSymbol(r))
	if classifyAsSymbol {
		return symbol_class
	}

	start := lex.Char - 1
	errorMessage, endChar := lex.validateUnderscoreNextChar()
	if errorMessage == "" {
		return underscore_class
	}

	lex.error2(errorMessage, lex.Line, start, endChar)
	return error_class
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

// Determines the class of some input section based on some byte `c` of the
// input. Unless there's a good reason to do otherwise, `c` is the first
// character of that input section.
func (lex *Lexer) classify(c byte) (class symbolClass) {
	r := rune(c)
	if unicode.IsLetter(r) {
		class = identifier_class
	} else if unicode.IsDigit(r) {
		class = number_class
	} else if c == '\'' {
		class = char_class
	} else if c == '"' {
		class = string_class
	} else if c == '_' {
		class = lex.reclassifyUnderscore()
	} else if c == '-' {
		class = lex.classifyMinus()
	} else if isSymbol(c) {
		class = symbol_class
	} else {
		class = error_class
		lex.error(InvalidCharacter)
	}
	return
}

var endMultiCommentRegex = regexp.MustCompile(`\*-`)

// creates and pushes token for single-line comment with Value=`lineAfterDashes`
func (lex *Lexer) getSingleLineComment(lineAfterDashes string) (ok, eof bool) {
	lex.Char += len(lineAfterDashes)
	tok := token.Comment.MakeValued(lineAfterDashes)
	lex.add(tok)
	return true, false
}

// reads, creates, and pushes token for multi-line comment with Value=`lineAfterDashes`
func (lex *Lexer) getMultiLineComment(line string) (ok, eof bool) {
	var comment string = "" // full comment
	var next string = line  // next line to analyze
	var loc []int = nil

	getNext := func(line string) []int {
		return endMultiCommentRegex.FindStringIndex(line)
	}

	// allow function to have ok==true for one line comment, i.e., "-**-"
	ok = true

	// append input read to comment until end of comment is reached
	for loc = getNext(line); loc == nil; loc = getNext(next) {
		ok, eof = lex.advanceLine()
		if eof { // at eof?
			ok = false
			lex.error(UnexpectedEOF)
			return
		}
		lex.SavedChar.Pop()
		lex.SavedChar.Push(lex.Char)

		comment = comment + next

		// get next line
		next, _ = lex.remainingLine()
	}

	startIndexOfCommentClose := loc[0]

	// grab last line of comment up to (but not including) '*-'
	next = next[:startIndexOfCommentClose]

	// create an push comment token
	comment = comment + next
	tok := token.Comment.MakeValued(comment)
	lex.Char += loc[1] // set just past end of comment
	lex.add(tok)
	return ok, false
}

func (lex *Lexer) analyzeComment() (ok, eof bool) {
	lex.SavedChar.Push(lex.Char)

	var c byte
	_, eof = lex.nextChar() // remove initial '-'
	if ok = !eof; !ok {
		panic("bug: source not validated before calling analyzeComment")
	}

	// remove second thing ('-' or '*'), the value of `c` will determine the
	// branch to take in the condition below the next one
	c, eof = lex.nextChar()
	if ok = !eof; !ok {
		lex.error(UnexpectedEOF)
		return
	}

	// given a line
	//		-- abc ..
	// or
	//		-* abc ..
	// `line` is
	//		line = " abc .."
	var line string
	line, ok = lex.remainingLine()
	if c == '-' { // single line comment
		return lex.getSingleLineComment(line)
	} else if c == '*' { // multi line comment
		return lex.getMultiLineComment(line)
	}
	panic("bug in analyzeComment: else branch reached")
}

var intRegex = regexp.MustCompile(`[0-9](_*[0-9]+)*`)
var hexRegex = regexp.MustCompile(`(0x|0X)[0-9a-fA-F](_*[0-9a-fA-F]+)*`)
var octRegex = regexp.MustCompile(`(0o|0O)[0-7](_*[0-7]+)*`)
var binRegex = regexp.MustCompile(`(0b|0B)(0|1)(_*(0|1)+)*`)

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
		errorMessage = InvalidCharacterAtEndOfNumConst
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
		errorMessage = InvalidCharacterAtEndOfNumConst
		return
	}

	frac, ok := locateAtStart(line[numChars:], intRegex)
	if !ok { // <integer>.<non-integer>
		errorMessage = InvalidCharacter
		return
	}

	numChars = numChars + len(frac)
	num = num + "." + frac

	if len(line) <= numChars {
		return
	}

	hasE = isE(line, numChars)
	if !hasE && !isNumEndCharValid(line, numChars) { // <integer>.<integer><illegal-char>
		errorMessage = InvalidCharacter
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
		errorMessage = InvalidCharacterAtEndOfNumConst
		return
	}

	var sign string
	sign, numChars = analyzePossibleSign(line, numChars)

	if len(line) <= numChars { // <float>e<sign>EOL
		errorMessage = InvalidCharacterAtEndOfNumConst
		return
	}

	// read integer value that follows 'e'/'E'
	frac, ok := locateAtStart(line[numChars:], intRegex)
	if !ok { // <float>e[sign]<illegal-char>
		errorMessage = InvalidCharacter
		return
	}

	numChars = numChars + len(frac)
	if !isNumEndCharValid(line, numChars) { // <float>e[sign]<integer><illegal-char>
		errorMessage = InvalidCharacter
		return
	}

	// build value as string
	num = num + string(e) + sign + frac
	return
}

// return a number token. Could be either floating point number or integer
func maybeFractional(num, line string) (tok token.Token, numChars int, errorMessage string) {
	numChars, errorMessage = len(num), ""
	// remove leading zeros (so 0[integer] isn't mistaken as an octal number by llvm or go)
	for numChars != 0 && num[0] == '0' {
		num = num[1:]
	}

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

// return current line starting from current char until the end of line
func (lex *Lexer) remainingLine() (line string, ok bool) {
	line, ok = lex.currentSourceLine()
	if !ok || (lex.Char > 0 && lex.Char > len(line)) {
		return "", false
	}

	return line[lex.Char-1:], true
}

// read number from input
func (lex *Lexer) analyzeNumber() (ok, eof bool) {
	lex.SavedChar.Push(lex.Char)
	line, ok := lex.remainingLine()
	if !ok {
		panic("bug: function called without verifying readable source exists")
	}

	var tok token.Token // token result
	var numChars int    // total number of chars read
	var errorMessage string = ""

	// 0x, 0b, and 0o must be checked first, else the lexer might falsely think
	// '0' is the number
	if num, ok := locateAtStart(line, hexRegex); ok { // hex
		tok, numChars, errorMessage = analyzeNonBase10(num, line)
	} else if num, ok := locateAtStart(line, octRegex); ok { // oct
		tok, numChars, errorMessage = analyzeNonBase10(num, line)
	} else if num, ok := locateAtStart(line, binRegex); ok { // bin
		tok, numChars, errorMessage = analyzeNonBase10(num, line)
	} else if num, ok := locateAtStart(line, intRegex); ok { // int or float
		tok, numChars, errorMessage = maybeFractional(num, line)
	} else {
		numChars = 0
		errorMessage = InvalidCharacter
	}

	lex.Char += numChars
	if errorMessage != "" {
		lex.error(errorMessage)
		return false, eof
	}

	lex.add(tok)
	return true, false
}

// generate affixed regex around regex `element`
func affixedRegexGen(element string) string {
	return fmt.Sprintf(`(%s)?_?(((%s_)+(%s)?)|(%s))`, element, element, element, element)
}

// the following regular expression, or, alpha-numeric id regex occurring between or in-front-of/after zero or one '_':
//
//	`([a-z][a-zA-Z0-9']*)?_?((([a-z][a-zA-Z0-9']*_)+([a-z][a-zA-Z0-9']*)?)|([a-z][a-zA-Z0-9']*))`
//
// examples:
//
//	`not_`
//	`_mod_`
//	`zipWith`
var affixedIdRegex = regexp.MustCompile(affixedRegexGen(idRegexClassRaw + `*`))

// the following regular expression occurring between or in front of zero or one '_':
//
//	`\([!@#\$%\^\&\*~,<>\.\?/:\|-\+=`]+\)`
//
// examples:
//
//	`_>>=_`
//	`_!`
//	`!`
var affixedSymbolRegex = regexp.MustCompile(affixedRegexGen(freeSymbolRegexClassRaw + `+`))

func locateAtStart(s string, regex *regexp.Regexp) (string, bool) {
	loc := regex.FindStringIndex(s)
	if loc != nil && loc[0] == 0 {
		return s[:loc[1]], true
	}
	return "", false
}

// map of standalone symbols--these cannot be used within other identifiers
var standaloneMap = map[byte]token.Type{
	'(': token.LeftParen,
	')': token.RightParen,
	'[': token.LeftBracket,
	']': token.RightBracket,
	'{': token.LeftBrace,
	'}': token.RightBrace,
	',': token.Comma,
}

func (lex *Lexer) analyzeStandalone() (ok, eof bool) {
	var c byte
	c, eof = lex.peek()
	if eof {
		panic("bug: function called without source validation")
	}

	var tokenType token.Type
	if tokenType, ok = standaloneMap[c]; ok {
		lex.SavedChar.Push(lex.Char)
		token := tokenType.Make()
		lex.Char = lex.Char + 1
		lex.add(token)
	}
	return
}

// reads and creates a token for some sequence of symbols
//
// this includes the standalone symbols that cannot be used in names:
//
//	`( ) [ ] { } ,`
//
// keywords:
//
//	`alias derives import in let module use trait where`
//
// alpha-numeric identifiers (lower and upper case), the following are examples:
//
//	`x MyType Data map x2 x'`
//
// affixed identifiers, the following are examples:
//
//	`_::_ if_then_else_ _>>=_ _! _mod_`
//
// non-alpha-numeric identifiers, the following are examples:
//
//	`! &`
func (lex *Lexer) analyzeSymbol() (ok, eof bool) {
	ok, eof = lex.analyzeStandalone()
	if ok {
		return
	}

	return lex.analyzeIdentifier()
}

func getEscape(r rune, escapeString bool) (c byte, ok bool) {
	ok = true
	switch r {
	case 'n':
		c = '\n'
	case 't':
		c = '\t'
	case 'r':
		c = '\r'
	case 'v':
		c = '\v'
	case 'b':
		c = '\b'
	case 'a':
		c = '\a'
	case 'f':
		c = '\f'
	case '\\':
		c = '\\'
	case '"':
		if escapeString {
			c = '"'
		} else {
			ok = false
		}
	case '\'':
		if !escapeString {
			c = '\''
		} else {
			ok = false
		}
	default:
		ok = false
	}
	return
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

func (lex *Lexer) analyzeChar() (ok, eof bool) {
	lex.SavedChar.Push(lex.Char)

	var c byte
	c, eof = lex.nextChar() // should be leading '
	if ok = !(eof || c != '\''); !ok {
		if c != '\'' {
			panic("bug: called before validating next token is a character token")
		} else {
			lex.error(UnexpectedEOF)
		}
		return
	}

	var line string
	line, ok = lex.remainingLine()
	if !ok {
		lex.error(ExpectedCharLiteral)
		return
	}

	var escaped string
	var length int
	escaped, length, ok = readEscapable(line, '\'')
	if !ok {
		lex.error(IllegalEscapeSequence)
		return
	}

	lex.Char += length
	if c, eof = lex.nextChar(); c != '\'' || eof { // remove closing `'`
		ok = false
		if eof {
			lex.error(UnexpectedEOF)
		} else {
			lex.error(IllegalCharLiteral)
		}
		return
	}

	var index int
	escaped, ok, index = updateEscape(escaped, false)
	if !ok {
		lex.Char += index + 1
		lex.error(IllegalEscapeSequence)
		return
	}
	if ok = len(escaped) == 1; !ok {
		lex.error(IllegalCharLiteral)
		return
	}

	tok := token.CharValue.MakeValued(escaped)
	lex.add(tok)

	return
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

func (lex *Lexer) analyzeString() (ok, eof bool) {
	lex.SavedChar.Push(lex.Char)
	var c byte
	c, eof = lex.nextChar() // should be first quotation mark
	// check for leading quotation mark
	if eof || c != '"' {
		panic("bug: source was not verified")
	}

	content, ok := lex.getStringContent()
	if !ok {
		lex.error(IllegalStringLiteral)
		return
	}

	updatedContent, ok, _ := updateEscape(content, true)
	if !ok {
		lex.error(IllegalEscapeSequence)
		return
	}

	token := token.StringValue.MakeValued(updatedContent)
	lex.add(token)
	return
}

var idRegexClassRaw = `[a-z][a-zA-Z0-9']`
var typeIdRegexClassRaw = `[A-Z][a-zA-Z0-9']`
var capIdRegex = regexp.MustCompile(typeIdRegexClassRaw + `*`)

// finds a substring of `line` starting at index 0 to an index > 0 that is an identifier
//
// if a substring with the above requirements cannot be found, then a non-empty `errorMessage` is
// returned
func matchId(line string) (matched string, errorMessage string, illegalArgument bool) {
	if len(line) < 1 {
		return "", "", true 
	}

	matched = matchRegex(affixedIdRegex, line)
	if len(matched) > 0 {
		return
	}

	matched = matchRegex(affixedSymbolRegex, line)
	if len(matched) < 1 {
		errorMessage = InvalidCharacter
	}

	return
}

// validates affixed id
func checkAffixed(line, id string) (errorMessage string) {
	// check if id ends with '_', if it does, check that it isn't actually more than one '_' in `line`
	if id[len(id)-1] != '_' {
		return ""
	}

	if len(id) > len(line) {
		panic("id shouldn't be longer than line")
	} else if len(id) != len(line) && line[len(id)] == '_' {
		return InvalidAffixId
	}
	return ""
}

// matches non capital letter starting identifier (lowercase alpha-numeric-symbolic identifier)
//
// returns token type for the returned identifier string `id` and an empty string for `errorMessage`
// on success; otherwise, returns garbage for `id` and `ty`, and a non-empty string `errorMessage`
func matchNonCapId(line string) (id string, ty token.Type, errorMessage string) {
	var illegalArgument bool
	id, errorMessage, illegalArgument = matchId(line)
	if illegalArgument {
		panic("bug: illegal argument, empty string for argument `line`")
	} else if errorMessage != "" {
		return
	}

	if strings.ContainsRune(id, '_') {
		ty = token.Affixed
		errorMessage = checkAffixed(line, id)
	} else {
		ty = matchKeyword(id, token.Id) // id or some keyword
	}
	return
}

// reads and creates a token for some kind of identifier:
//
// alpha-numeric identifiers (lower and upper case), the following are examples:
//
//	`x MyType Data map x2 x'`
//
// affixed identifiers, the following are examples:
//
//	`_::_ if_then_else_ _>>=_ _! _mod_`
//
// non-alpha-numeric identifiers, the following are examples:
//
//	`! &`
func (lex *Lexer) analyzeIdentifier() (ok, eof bool) {
	lex.SavedChar.Push(lex.Char)
	line, ok := lex.remainingLine()
	if !ok {
		panic("bug: function called without verifying readable source exists")
	}

	var ty token.Type
	errorMessage := ""

	id := matchRegex(capIdRegex, line)
	if len(id) > 0 {
		ty = token.CapId
	} else {
		id, ty, errorMessage = matchNonCapId(line)
	}

	if errorMessage != "" {
		lex.error(errorMessage)
		return false, false
	}

	// add token
	lex.Char += len(id) // this works b/c id is copied from `line` (not modified)
	token := ty.MakeValued(id)
	lex.add(token)
	return true, false
}

// when char after first underscore is valid, returns ("", 0); otherwise, returns non-empty string
// representing error message and the character to put for the end char number of the error
func (lex *Lexer) validateUnderscoreNextChar() (errorMessage string, errorEndChar int) {
	c, eof := lex.peek()
	if eof {
		return "", 0
	}

	r := rune(c)
	ok := r != '_'
	if !ok {
		return IllegalUnderscoreSequence, lex.Char + 1
	}

	ok = unicode.IsDigit(r)
	if !ok {
		return InvalidUnderscore, lex.Char
	}

	return "", 0
}

func (lex *Lexer) analyzeUnderscore() (ok, eof bool) {
	lex.SavedChar.Push(lex.Char)
	_, eof = lex.nextChar()
	if ok = !eof; !ok {
		panic("bug: source was not validated")
	}

	tok := token.Underscore.Make()
	lex.add(tok)
	return
}

func (class symbolClass) analyze(lex *Lexer) (ok, eof bool) {
	switch class {
	case number_class:
		return lex.analyzeNumber()
	case identifier_class:
		fallthrough
	case symbol_class:
		return lex.analyzeSymbol()
	case char_class:
		return lex.analyzeChar()
	case string_class:
		return lex.analyzeString()
	case underscore_class:
		return lex.analyzeUnderscore()
	case comment_class:
		return lex.analyzeComment()
	}

	lex.error(InvalidCharacter)
	return false, false
}

// true iff lexer is at end of source
func (lex *Lexer) eof() bool {
	if lex.Line > len(lex.Source) {
		return true
	}

	if lex.Line != len(lex.Source) {
		return false
	}

	line := lex.Source[lex.Line-1]
	return lex.Char > len(line)
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
		lex.Char++
	}
	return
}

// unadvance input by a single char and returns the new current char
func (lex *Lexer) ungetChar() (c byte) {
	if lex.Char > 1 {
		lex.Char = lex.Char - 1
	} else if lex.Line > 1 {
		lex.Line = lex.Line - 1
		lex.Char = len(lex.Source[lex.Line-1])
	} else {
		panic("bug: cannot move input to a position before its start")
	}

	if !lex.ValidChar() {
		panic("bug: illegal char unget")
	}

	c, _ = lex.currentSourceChar()
	return
}

// reads whitespace until next non-whitespace char, then advances input and returns non-whitespace
// char (and eof==true when lexer is at end of source)
//
// whitespace is just ' ' and '\t'
func (lex *Lexer) skipWhitespace() (eof bool) {
	// technically, condition isn't needed b/c if eof==true, then c==0 which is the default case
	var c byte
	c, eof = lex.peek()
	for ; !eof; c, eof = lex.peek() {
		switch c {
		case ' ':
			fallthrough // advance char counter
		case '\t':
			lex.Char++ // advance char counter
		default: // non whitespace
			return eof
		}
	}
	return eof
}

func eofLineAdvance(lex *Lexer) (ok, eof bool) {
	eof = true
	ok = lex.Line > len(lex.Source) // if already at EOF, then not ok; else ok
	lex.Line = len(lex.Source) + 1  // set to eof
	return
}

// sets line to next line and char to one
func setLineToNext(lex *Lexer) {
	lex.Line++
	lex.Char = 1
}

// moves line counter to next line and char counter to first char
//
// returns ok==true if line counter was advanced and eof==true either when already at end of source
// or if advancing the line counter puts lexer at end of input
func (lex *Lexer) advanceLine() (ok, eof bool) {
	if (lex.Line + 1) > len(lex.Source) {
		// returns ok==true when not at EOF before advancing; eof==true unconditionally
		return eofLineAdvance(lex)
	}

	setLineToNext(lex)
	return true, false
}

func (lex *Lexer) analyze() (ok bool, eof bool) {
	eof = lex.skipWhitespace()

	if eof {
		return true, eof
	}

	c, _ := lex.nextChar()
	if c == '\n' {
		return lex.advanceLine()
	}

	// use char to determine what class new token will belong to
	class := lex.classify(c)
	if class == error_class {
		return false, false
	}
	lex.ungetChar() // unget char gotten from lex.nextChar

	// use class information to get token
	return class.analyze(lex)
}

// prepares lexer for reading from source based on whether it's already read and whether the source is empty or not
func (lex *Lexer) fixLineChar() {
	if lex.Line == 0 {
		lex.Line = common.Min(1, len(lex.Source)) // set line to zero when no source code, else set to one
		lex.Char = lex.Line                       // this will make (line, char) == (1, 1) OR (line, char) == (0, 0)
	} else if lex.Char != 1 {
		lex.Char = 1
	}

	debug_validateLineAndChar(lex)
}

// tokenize lex source
func (lex *Lexer) Tokenize() (tokens []token.Token, ok bool) {
	lex.fixLineChar() // prepare for reading from source

	var eof bool
	// keep reading tokens until end of input
	for ok = true; ok && !eof; {
		ok, eof = lex.analyze()
	}

	// either one is fine b/c might try to get eof twice which will make ok==false
	ok = ok || eof
	return lex.Tokens, ok
}
