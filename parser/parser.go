// =================================================================================================
// Alex Peters - January 24, 2024
//
// Parser struct and methods
// =================================================================================================

package parser

import (
	"github.com/petersalex27/yew/errors"
	"github.com/petersalex27/yew/token"
)

type Parser struct {
	Path           string
	Names          []token.Token
	PositionRanges []int
	module         *Module
	Tokens         []token.Token
	tokenCounter   int
	messages       []errors.ErrorMessage
	Current, Next  token.Token
	optionalFlag   bool
	panicking      bool
}

// get next token and advance input
func (p *Parser) Advance() (tok token.Token) {
	tok = p.Peek()
	p.Current, p.Next = p.Next, tok
	p.tokenCounter++
	return tok
}

// get next token then check if next token has the type `typeCondition`.
//
// If the next token has that type, then truthy == true; otherwise, truthy == false
//
// The token stream is advanced regardless of the truthiness of truthy
func (p *Parser) CheckAdvance(typeCondition token.Type) (tok token.Token, truthy bool) {
	tok = p.Peek()
	p.Current, p.Next = p.Next, tok
	p.tokenCounter++
	truthy = p.Current.Type == typeCondition
	return
}

// get next token but only advance input if next token has the type `typeCondition`.
//
// If the next token has that type, then truthy == true; otherwise, truthy == false
//
// The token stream is advanced if and only if truthy == true. Another way to think about this is
// a (*Parser).Peek followed by a condition advance where the condition is whether Peek returns a
// token with type `typeCondition`
func (p *Parser) ConditionalAdvance(typeCondition token.Type) (tok token.Token, truthy bool) {
	truthy = p.Next.Type == typeCondition
	if truthy {
		tok = p.Peek()
		p.Current, p.Next = p.Next, tok
		p.tokenCounter++
	} else {
		tok = p.Current
	}
	return
}

// returns parsed module and clears parsers internal module data
func (p *Parser) GetModule() Module {
	if p.module == nil {
		panic("bug: no module")
	}

	module := *p.module
	p.module = nil
	return module
}

// initializes a returns new parser
func Init(path string, positions []int, tokens []token.Token) *Parser {
	return &Parser{
		Path:           path,
		PositionRanges: positions,
		Tokens:         tokens,
		messages:       []errors.ErrorMessage{},
	}
}

// clears parser's message buffer and returns all messages
func (parser *Parser) FlushMessages() (messages []errors.ErrorMessage) {
	messages = parser.messages
	parser.messages = []errors.ErrorMessage{}
	return messages
}

// true iff parser has recorded an error since the last time the errors have been reset
func (p *Parser) Panicking() bool { return p.panicking }

// returns next token but does not advance past it
func (p *Parser) Peek() token.Token {
	if p.tokenCounter >= len(p.Tokens) {
		return endToken()
	}

	return p.Tokens[p.tokenCounter]
}

// start optional parsing--returns defer-able function that restores previous option flag value
func (parser *Parser) StartOptional() (deferEnd func()) {
	save := parser.optionalFlag
	parser.optionalFlag = true
	return func() { parser.optionalFlag = save }
}

// stops optional parsing--returns a defer-able function that restores previous option flag value
func (parser *Parser) StopOptional() (deferEnd func()) {
	save := parser.optionalFlag
	parser.optionalFlag = false
	return func() { parser.optionalFlag = save }
}

// adds a message to parser's internal messages slice
func (parser *Parser) addMessage(e errors.ErrorMessage) {
	parser.messages = append(parser.messages, e)
	parser.panicking = parser.panicking || e.IsFatal()
}

// returns an "End" token
func endToken() token.Token {
	return token.Token{Type: token.EndOfTokens}
}

// skips sequence of tokens of type ty
func (p *Parser) skip(ty token.Type) {
	if ty == token.EndOfTokens {
		return // do not allow this, just exit method
	}

	nextType := p.Peek().Type
	for nextType == ty {
		_ = p.Advance()
		nextType = p.Peek().Type
	}
}
