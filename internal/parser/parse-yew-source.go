package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/common/data"
)

func makeYewSource(h data.Maybe[header], b data.Maybe[body], f data.Maybe[annotations]) yewSource {
	return yewSource{
		header: h,
		body: b,
		footer: footer{f},
		Position: api.WeakenRangeOver[api.Node](h, b, f),
	}
}

// Parse yew source file.
//
//	`yew source = {"\n"}, [meta, {"\n"}], [header, {"\n"}], [body, {"\n"}], footer ;`
func parseYewSource(p Parser) Parser {
	p.dropNewlines()

	// parse the yew source file surface from the top down
	es, header, isHeader := parseHeader(p).Break()
	if !isHeader {
		return writeErrors(p, es)
	}

	res, mFooterAnnots := parseBody(p)
	es, mb, isMB := res.Break()
	if !isMB {
		return writeErrors(p, es)
	}

	p.dropNewlines()
	if es := assertEof(p); es != nil {
		return writeErrors(p, *es)
	}

	ys := makeYewSource(header, mb, mFooterAnnots)
	// record the AST in the parser state on success
	if ps, ok := p.(*ParserState); ok {
		ps.ast = ys
		p = ps
	} else if pso, ok := p.(*ParserState_optional); ok {
		// suspicious, but okay ...
		pso.ParserState.ast = ys
		p = pso
	}

	return p
}

func assertEof(p Parser) *data.Ers {
	if !matchCurrent(token.EndOfTokens)(p) {
		e := data.MkErr(ExpectedEndOfFile, p.current())
		es := data.Makes(e)
		return &es
	}
	return nil
}