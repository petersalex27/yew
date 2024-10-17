package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
)

type ParserStateFail struct{ bad *ParserState }

func (p *ParserStateFail) dropNewlines() { /* noop */ }

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

// just adds the error
func (p *ParserStateFail) report(e error, _ bool) parser {
	p.bad.AddError(e)
	return p
}
