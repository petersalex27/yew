// =============================================================================
// Alex Peters - January 21, 2024
// =============================================================================
package lexer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/petersalex27/yew/common/stack"
	"github.com/petersalex27/yew/errors"
	"github.com/petersalex27/yew/token"
)

type Lexer struct {
	// path to source file
	path pathSpec
	// write some amount of source to lexer
	write func(*Lexer) bool
	// source file as an array of strings for each non-empty line, does not include newline chars
	Source []string
	// current line number
	Line int
	// current char number for the given line
	Char int
	// saved char number
	SavedChar *stack.Stack[int]
	// tokens created from source
	Tokens []token.Token
	// errors, warnings, and logs during lexical analysis
	messages []errors.ErrorMessage
}

func (lex *Lexer) Messages() []errors.ErrorMessage {
	return lex.messages
}

func (lex *Lexer) FlushMessages() []errors.ErrorMessage {
	messages := lex.messages
	lex.messages = []errors.ErrorMessage{}
	return messages
}

// initialize lexer for writing to its internal source buffer `lex.Source` from the input stream
// located by path.
//
// given
//
//	lex := lexer.Init(myPath)
//
// write input to lex.Source with
//
//	lex.Write()
func Init(path pathSpec) *Lexer {
	lex := new(Lexer)

	lex.path = path
	// beyond 8 being a small power of two, it's an arbitrary choice 
	const cap uint = 8
	lex.SavedChar = stack.NewStack[int](cap) 

	// generate source code write-to-lexer function
	if _, ok := path.(standardInput); ok {
		lex.write = genWriteFromStdin()
	} else {
		lex.write = lexWrite_fromPath
	}

	return lex
}

// saves current char

// converts receiver to string
//
// intended for debugging and fail messages for tests
func (lex *Lexer) String() string {
	return fmt.Sprintf(
		"Lexer{path: %v, write: nil ? %t, Source: %v, Line: %d, Char: %d, Tokens: %v, messages: %v}", 
		lex.path, lex.write == nil, lex.Source, lex.Line, lex.Char, lex.Tokens, lex.messages,
	)
}

// returns true if and only if character number is valid for the current line
func (lex *Lexer) ValidChar() bool {
	line, ok := lex.currentSourceLine()
	return ok && validateCharNumber(line, lex.Char)
}

// returns true if and only if line number is a valid line position for lex's source
func (lex *Lexer) ValidLine() bool {
	return (lex.Line > 0) && (lex.Line <= len(lex.Source))
}

// writes contents of input to lexer's source slice
//
// returns number of lines written else -1 on error
func (lex *Lexer) Write() int {
	before := len(lex.Source)
	if !lex.write(lex) {
		return -1
	}
	after := len(lex.Source)
	return after - before
}

func (lex *Lexer) add(token token.Token) {
	start, _ := lex.SavedChar.Pop()
	token.Line, token.Start, token.End = lex.Line, start, lex.Char
	lex.Tokens = append(lex.Tokens, token)
}

// adds a message to lex's message slice
func (lex *Lexer) addMessage(e errors.ErrorMessage) {
	lex.messages = append(lex.messages, e)
}

// returns current char of source
func (lex *Lexer) currentSourceChar() (char byte, ok bool) {
	var line string
	line, ok = lex.currentSourceLine()
	if ok = ok && validateCharNumber(line, lex.Char); ok {
		char = line[lex.Char-1]
	}
	return char, ok
}

// returns current line of source
func (lex *Lexer) currentSourceLine() (line string, ok bool) {
	if ok = lex.ValidLine(); ok {
		line = lex.Source[lex.Line-1]
	}

	return line, ok
}

// writes input from stdin
func genWriteFromStdin() func(lex *Lexer) bool {
	return genWriteFromStream(os.Stdin)
}

// writes to source slice from input stream line-by-line--this allows calling lex.Write multiple times
func genWriteFromStream(stream *os.File) func(lex *Lexer) bool {
	reader := bufio.NewReader(os.Stdin)
	// closure on `reader`
	return func(lex *Lexer) bool {
		switch line, err := reader.ReadString('\n'); err {
		case nil:
			lex.Source = append(lex.Source, line)
			fallthrough
		case io.EOF:
			return true
		default:
			return false
		}
	}
}

// function writes contents of file to source slice in lexer, then prevents further writing
func lexWrite_fromPath(lex *Lexer) bool {
	f := lex.openPath()
	if f == nil {
		return false
	}

	defer f.Close()

	lex.Source = readSourceFile(f)

	lex.write = nil // prevent further writing

	return true
}

// opens file located by lex.path
//
// appends an error on failure and returns nil
//
// on success, returns file opened
func (lex *Lexer) openPath() *os.File {
	path := lex.path.String()

	f, err := os.Open(path)
	if err == nil {
		return f
	}

	msg := err.(*os.PathError).Err.Error()
	e := makeOSError(msg)
	lex.addMessage(e)
	return nil
}

// reads entire source file, splitting input at newlines
func readSourceFile(f *os.File) []string {
	buf := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		buf = append(buf, text)
	}
	return buf
}

// reads input through byte `end` and returns it and true iff successful
func (lex *Lexer) readThrough(end byte) (string, bool) {
	res, ok := lex.readUntil(end)
	if !ok {
		return res, ok
	}

	last, eof := lex.nextChar()
	if ok = !eof; !ok {
		return res, ok
	}

	return res + string(last), ok
}

// reads input to (but not including) byte `end` and returns it and true iff successful
func (lex *Lexer) readUntil(char byte) (string, bool) {
	var builder strings.Builder
	c, ok := lex.currentSourceChar()
	for ; c != char && ok; c, ok = lex.currentSourceChar() {
		in, _ := lex.nextChar()
		builder.WriteByte(in)
	}
	return builder.String(), ok
}

// returns true if and only if character number is a valid character position for `line`
func validateCharNumber(line string, charNumber int) bool {
	return charNumber > 0 && charNumber <= len(line)
}
