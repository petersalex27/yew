package parser

import (
	"math"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/api/util"
	//"github.com/petersalex27/yew/internal/symbol"
)

type state struct {
	scanner      api.ScannerPlus
	tokens       []api.Token
	tokenCounter int
	errors       []error
	warnings     []error
	//names        *symbol.Table
}

// State for parsing
type ParserState struct {
	state
	ast *yewSource
}

func createState(scanner api.ScannerPlus) state {
	return state{
		scanner:  scanner,
		tokens:   nil,
		errors:   make([]error, 0),
		warnings: make([]error, 0),
		//names:    symbol.New(),
	}
}

type ParserState_optional struct {
	timesMarked int64
	*ParserState
}

// kinda noop--panics if called when timesMarked <= 0
func (p *ParserState_optional) report(e error, fatal bool) Parser {
	if p.timesMarked <= 0 {
		// this shouldn't have been called if <= 0
		//
		// this could have, but should not have, been done by accessing the parser's
		// fields directly--but, also, calling demarkOptional() more times than markOptional() (or a
		// combination of the two .-.)
		panic("optional parser was de-marked more times than it was marked")
	}
	return p // ignore the error, as the thing that was being parsed was optional
}

func (p *ParserState) markOptional() Parser {
	return &ParserState_optional{
		timesMarked: 1,
		ParserState: p,
	}
}

// noop
func (p *ParserState) demarkOptional() Parser {
	return p
}

func (p *ParserState_optional) markOptional() Parser {
	if p.timesMarked == math.MaxInt64 {
		// this is insane and should never happen
		panic("cannot mark optional more than the value of 'math.MaxInt64' times")
	}
	p.timesMarked++
	return p
}

func (p *ParserState_optional) demarkOptional() Parser {
	if p.timesMarked > 1 {
		p.timesMarked--
	} else {
		return p.ParserState
	}
	return p
}

func (p *ParserState) dropNewlines() {
	for token.Newline.Match(p.current()) {
		p.advance()
	}
}

func (p *ParserState) srcCode() api.SourceCode {
	return p.scanner.SrcCode()
}

func (p *ParserState) Pos() (int, int) {
	return p.current().Pos()
}

func (p *ParserState) GetPos() api.Position {
	return p.current().GetPos()
}

func (p *ParserState) bind(f func(Parser) Parser) Parser { return f(p) }

// adds the error, and if fatal, then returns a *ParserStateFail
func (p *ParserState) report(e error, fatal bool) Parser {
	p.AddError(e)
	if fatal {
		return &ParserStateFail{bad: *p}
	}
	return p
}

func (parser *ParserState) AppendTokens(tokens ...api.Token) {
	if parser.tokens == nil {
		cap := int64(float64(len(tokens)) * 1.25)
		parser.tokens = make([]api.Token, 0, int64(cap))
	}
	parser.tokens = append(parser.tokens, tokens...)
}

func (parser *ParserState) Ast() api.Node {
	return parser.ast
}

func (parser *ParserState) load() bool {
	tokens, errorToken := util.Tokenize(parser.scanner, nil)
	if errorToken != nil {
		parser.AddError((*errorToken).Error())
		return false
	}
	parser.tokens = tokens
	return true
}

func (parser *ParserState) replStatement() bool {
	panic("not implemented")
}

func (parser *ParserState) ReplParse() bool {
	if !parser.load() {
		return false
	}

	return parser.replStatement()
}

func (parser *ParserState) Parse() bool {
	if !parser.load() {
		return false
	}

	parser.Run()
	return len(parser.Errors()) == 0
}

func (parser *ParserState) AddError(err error) {
	parser.errors = append(parser.errors, err)
}

// ensure that the errors slice is never nil when needed
func (parser *ParserState) Errors() []error {
	if parser.errors == nil {
		parser.errors = make([]error, 0)
	}
	return parser.errors
}

func (parser *ParserState) ReferenceScanner() *api.Scanner {
	s := api.Scanner(parser.scanner)
	return &s
}

func (*ParserState) Clear() api.Parser {
	return &ParserState{}
}

func (parser *ParserState) current() api.Token {
	if parser.tokenCounter >= len(parser.tokens) {
		return token.EndOfTokens.Make()
	}
	return parser.tokens[parser.tokenCounter]
}

func (parser *ParserState) advance() {
	if parser.tokenCounter < len(parser.tokens) {
		parser.tokenCounter++
	}
}

func (p *ParserState) Run() {
	// TODO: Implement
	panic("not implemented")
}
