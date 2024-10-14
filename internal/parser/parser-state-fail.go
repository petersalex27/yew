package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
)

type ParserStateFail struct {
	bad *ParserState
}

// noop
func (p *ParserStateFail) acceptRoot(yewSource) Parser {
	return p
}

// noop
func (p *ParserStateFail) markOptional() Parser { return p }

// noop
func (p *ParserStateFail) demarkOptional() Parser { return p }

func (p *ParserStateFail) dropNewlines() {}

func (p *ParserStateFail) srcCode() api.SourceCode {
	return p.bad.scanner.SrcCode()
}

func (p *ParserStateFail) Pos() (int, int) {
	return p.bad.Pos()
}

func (p *ParserStateFail) GetPos() api.Position {
	return p.bad.GetPos()
}

func (p *ParserStateFail) current() api.Token { return token.EndOfTokens.Make() }

func (p *ParserStateFail) advance() { /* noop */ }

// noop
func (p *ParserStateFail) bind(func(Parser) Parser) Parser { return p }

// just adds the error
func (p *ParserStateFail) report(e error, _ bool) Parser {
	p.bad.AddError(e)
	return p
}
