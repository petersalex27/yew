// =============================================================================
// Alex Peters - January 21, 2024
// =============================================================================
package lexer

import (
	"maps"
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/internal/source"
)

type lexerState struct {
	source.SourceCode
	// current line number
	Line int
	// current position in source
	Pos int
}

type nextAction byte

const (
	appendSrc nextAction = iota
	restoreSrc
)

type Lexer struct {
	listening bool
	lexerState
	restore lexerState
	// keywords
	keywords  map[string]token.Type
	action    chan nextAction
	additions chan string
}

func (lex *Lexer) SrcCode() api.SourceCode {
	return lex.SourceCode
}

func (lex *Lexer) AppendSource(addition string) {
	//print("appending addition ...\n")
	lex.additions <- addition
	lex.appendAction()
}

// NOTE: there's only an effective restore history of one
func (lex *Lexer) Restore() {
	lex.action <- restoreSrc
}

func (lex *Lexer) appendAction() {
	lex.action <- appendSrc
}

func (lex *Lexer) Stop() {
	close(lex.action)
	close(lex.additions)
}

func (lex *Lexer) Eof() bool {
	return lex.Pos >= len(lex.Source)
}

func (lex *Lexer) copyState() lexerState {
	return lexerState{
		SourceCode: lex.SourceCode.Copy(),
		Line:       lex.Line,
		Pos:        lex.Pos,
	}
}

func (lex *Lexer) actionListener(addition string) {
	act := <-lex.action // block until action is received
	if act == appendSrc {
		lex.restore = lex.copyState()
		oldLength := len(lex.Source)
		lex.SourceCode.AppendSource(addition)
		if len(addition) > 0 {
			lex.Pos = oldLength // set position to end of old source (equiv. start of new source)
		}
	} else if act == restoreSrc {
		lex.lexerState = lex.restore
	}
}

func (lex *Lexer) additionListener() {
	for addition := range lex.additions {
		lex.actionListener(addition)
	}
}

// create a new, empty lexer
func New() *Lexer {
	lex := new(Lexer)

	lex.SourceCode = source.SourceCode{}
	lex.Line = 1

	lex.action = make(chan nextAction, 1)
	lex.additions = make(chan string, 1)
	lex.restore = lex.copyState()

	return lex
}

// map is copied, preventing changes to the original map from affecting the lexer
//
// do not mutate the map while the this function is running--this function is not thread-safe
func (lex *Lexer) LoadReserved(reserved map[string]token.Type) {
	if len(lex.keywords) != 0 {
		panic("reserved keywords already loaded")
	}
	lex.keywords = make(map[string]token.Type, len(reserved))
	maps.Copy(lex.keywords, reserved)
}

func (lex *Lexer) Load(src api.Source) bool {
	if len(lex.SourceCode.Source) > 0 || lex.SourceCode.Path() != "" {
		return false // already loaded
	}

	lex.SourceCode = (source.SourceCode{}).Set(src).(source.SourceCode)
	return true
} 

func Init(src api.Source) *Lexer {
	lex := New()
	lex.Load(src)
	lex.LoadReserved(keywords)
	return lex
}

func InitRepl(initial api.Source) *Lexer {
	lex := Init(initial)
	lex.addCommandKeys()
	lex.listening = true

	// start listener
	go lex.additionListener()

	return lex
}

func InitListening(src api.Source) *Lexer {
	lex := Init(src)
	lex.listening = true

	// start listener
	go lex.additionListener()

	return lex
}

// returns index position for given line from start (inclusive) to end (exclusive)
func (lexer *Lexer) LinePos(line int) (start, end int) {
	endPositions := lexer.EndPositions()
	numLines := len(endPositions)
	if line > numLines || line < 1 {
		panic("bug: illegal argument, line > number of lines or < 1")
	}

	if line > 1 {
		start = endPositions[line-2]
	} // else start=0

	end = endPositions[line-1]
	return
}

// returns true if and only if line number is a valid line position for lex's source
func (lex *Lexer) ValidLine() bool {
	return (lex.Line > 0) && (lex.Line <= lex.Lines())
}

func (lex *Lexer) addCommandKeys() {
	// add command keys
	for key, typ := range commands {
		lex.keywords[key] = typ
	}
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


// ready to load source code
func (lex *Lexer) CreateSimilar() api.ScannerPlus {
	out := New()
	out.LoadReserved(lex.keywords)
	return out
}