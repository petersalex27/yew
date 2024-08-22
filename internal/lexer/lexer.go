// =============================================================================
// Alex Peters - January 21, 2024
// =============================================================================
package lexer

import (
	"bufio"
	"fmt"
	"io"
	"maps"
	"os"
	"strings"
	
	"github.com/petersalex27/yew/common"
	"github.com/petersalex27/yew/common/stack"
	"github.com/petersalex27/yew/errors"
	"github.com/petersalex27/yew/source"
	"github.com/petersalex27/yew/token"
)

type Lexer struct {
	keepComments bool
	illegalWhitespace byte
	// write some amount of source to lexer
	write func(*Lexer) bool
	source.SourceCode
	// current line number
	Line int
	// current position in source
	Pos int
	// saved char number
	SavedChar *stack.Stack[int]
	// Affixed Identifiers
	Affixed      []int
	AffixIndexes []int
	// tokens created from source
	Tokens []token.Token
	// errors, warnings, and logs written during lexical analysis
	messages []errors.ErrorMessage
	// keywords
	keywords map[string]token.Type
	// index for Next calls
	nextIndex int
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
func Init(path source.PathSpec) *Lexer {
	lex := new(Lexer)

	lex.SetKeepComments(false)
	lex.SetSpaceWhitespace()

	lex.Path = path
	// beyond 8 being a small power of two, it's an arbitrary choice
	const cap uint = 8
	lex.SavedChar = stack.NewStack[int](cap)

	lex.AffixIndexes = make([]int, 0, 8)
	lex.Affixed = make([]int, 0, 8)

	// generate source code write-to-lexer function
	if _, ok := path.(source.StandardInput); ok {
		lex.write = genWriteFromStdin()
	} else if _, ok := source.AssertStringInput(path); ok {
		lex.write = lexWrite_fromString
	} else {
		lex.write = lexWrite_fromPath
	}

	lex.keywords = make(map[string]token.Type, len(keywords))
	maps.Copy(lex.keywords, keywords)

	return lex
}

func (lex *Lexer) SetKeepComments(truthy bool) {
	lex.keepComments = truthy
}

func (lex *Lexer) SetSpaceWhitespace() {
	lex.illegalWhitespace = '\t'
}

func (lex *Lexer) SetTabWhitespace() {
	lex.illegalWhitespace = ' '
}

// returns index position for given line from start (inclusive) to end (exclusive)
func (lexer *Lexer) LinePos(line int) (start, end int) {
	numLines := len(lexer.PositionRanges)
	if line > numLines || line < 1 {
		panic("bug: illegal argument, line > number of lines or < 1")
	}

	if line > 1 {
		start = lexer.PositionRanges[line-2]
	} // else start=0

	end = lexer.PositionRanges[line-1]
	return
}

// saves current char

// converts receiver to string
//
// intended for debugging and fail messages for tests
func (lex *Lexer) String() string {
	return fmt.Sprintf(
		"Lexer{Path: %v, write: nil ? %t, Source: %s, Line: %d, Pos: %d, Tokens: %v, messages: %v}",
		lex.Path, lex.write == nil, string(lex.Source), lex.Line, lex.Pos, lex.Tokens, lex.messages,
	)
}

// returns true if and only if line number is a valid line position for lex's source
func (lex *Lexer) ValidLine() bool {
	return (lex.Line > 0) && (lex.Line <= len(lex.PositionRanges))
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
	token.Start, token.End = start, lex.Pos
	lex.Tokens = append(lex.Tokens, token)
}

// adds a message to lex's message slice
func (lex *Lexer) addMessage(e errors.ErrorMessage) {
	lex.messages = append(lex.messages, e)
}

// returns current char of source
func (lex *Lexer) currentSourceChar() (char byte, ok bool) {
	if lex.Pos >= len(lex.Source) {
		return 0, false
	}
	return lex.Source[lex.Pos], true
}

// returns current line of source
func (lex *Lexer) currentSourceLine() (line string, ok bool) {
	if ok = lex.ValidLine(); ok {
		start, end := lex.LinePos(lex.Line)
		line = string(lex.Source[start:end])
	}

	return line, ok
}

// writes input from stdin
func genWriteFromStdin() func(lex *Lexer) bool {
	return genWriteFromStream(os.Stdin)
}

// writes to source slice from input stream line-by-line--this allows calling lex.Write multiple times
func genWriteFromStream(stream *os.File) func(lex *Lexer) bool {
	reader := bufio.NewReader(stream)//bufio.NewReader(os.Stdin)
	// closure on `reader`
	return func(lex *Lexer) bool {
		switch line, err := reader.ReadBytes('\n'); err {
		case nil:
			length := len(line)
			if len(lex.PositionRanges) != 0 {
				length += lex.PositionRanges[len(lex.PositionRanges)-1]
			}
			lex.PositionRanges = append(lex.PositionRanges, length)
			lex.Source = append(lex.Source, line...)
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
	lex.write = nil // prevent further writing

	f := lex.openPath()
	if f == nil {
		return false
	}

	defer f.Close()

	lex.Source, lex.PositionRanges = read(f)

	return true
}

func lexWrite_fromString(lex *Lexer) bool {
	lex.write = nil // prevent further writing
	si, _ := source.AssertStringInput(lex.Path)
	r := common.NewStringReader(si.GetInput())
	lex.Source, lex.PositionRanges = read(r)
	return true
}

// opens file located by lex.path
//
// appends an error on failure and returns nil
//
// on success, returns file opened
func (lex *Lexer) openPath() *os.File {
	path := lex.Path.Path()

	f, err := os.Open(path)
	if err == nil {
		return f
	}

	msg := err.(*os.PathError).Err.Error()
	e := makeOSError(msg)
	lex.addMessage(e)
	return nil
}

func read(r io.Reader) ([]byte, []int) {
	buf := []byte{}
	pos := []int{}
	scanner := bufio.NewScanner(r)
	scanner.Split(ScanLines)
	tot := 0
	for scanner.Scan() {
		text := scanner.Text()

		tot += len(text)
		pos = append(pos, tot)
		buf = append(buf, []byte(text)...)
	}
	return buf, pos
}

// reads entire source file, splitting input at newlines
// func readSourceFile(f *os.File) ([]byte, []int) {
// 	buf := []byte{}
// 	pos := []int{}
// 	scanner := bufio.NewScanner(f)
// 	scanner.Split(ScanLines)
// 	tot := 0
// 	for scanner.Scan() {
// 		text := scanner.Text()

// 		tot += len(text)
// 		pos = append(pos, tot)
// 		buf = append(buf, []byte(text)...)
// 	}
// 	return buf, pos
// }

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
