package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/common/data"
)

func makeYewSource(h data.Maybe[header], b data.Maybe[body], f data.Maybe[annotations]) yewSource {
	return yewSource{
		header:   h,
		body:     b,
		footer:   footer{f},
		Position: api.WeakenRangeOver[api.Node](h, b, f),
	}
}

// Parse yew source file from the top down.
//
// rule:
//
//	```
//	yew source = {"\n"}, [header | body | header, then, body], {"\n"}, footer ;
//	```
func parseYewSource(p parser) parser {
	mb := data.Nothing[body]()

	// {"\n"}, [header | header, then, body | ..
	p.dropNewlines()
	es, header, isHeader := parseHeader(p).Break()
	if !isHeader {
		return writeErrors(p, es)
	}

	// the order is REALLY important here; `then` modifies state
	if header.IsNothing() || then(p) {
		// .. | body | .., then, body], {"\n"}, footer ;
		res := parseBody(p)
		es, mb_, isMB := res.Break()
		if !isMB {
			return writeErrors(p, es)
		}
		mb = mb_
	}

	// {"\n"}, footer ;
	p.dropNewlines()
	es, mFooterAnnots, ok := parseAnnotations_(p).Break()
	if !ok {
		return writeErrors(p, es)
	}

	if es := assertEof(p); es != nil {
		return writeErrors(p, *es)
	}

	ys := makeYewSource(header, mb, mFooterAnnots)
	// record the AST in the parser state on success
	if ps, ok := p.(*ParserState); ok {
		ps.ast = ys
		p = ps
	}

	return p
}

func assertEof(p parser) *data.Ers {
	if !matchCurrent(token.EndOfTokens)(p) {
		e := data.MkErr(ExpectedEndOfFile, p.current())
		es := data.Makes(e)
		return &es
	}
	return nil
}
