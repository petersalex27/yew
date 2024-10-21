package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/internal/kind"
	"github.com/petersalex27/yew/internal/symbol"
)

type ParserStateFail struct{ bad *ParserState }

func (p *ParserStateFail) declare(api.Token, ...byte) {}

func (p *ParserStateFail) newScope() {}

func (p *ParserStateFail) exitScope() map[string]symbol.Symbol { return nil }

// still create a new variable despite being in a fail state--this avoids nil pointer panics
func (p *ParserStateFail) newVar() kind.Type {
	return p.bad.newVar()
}

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
