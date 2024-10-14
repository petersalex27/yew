package parser

import (
	"github.com/petersalex27/yew/api/token"
	//"github.com/petersalex27/yew/api/util/fun"
	"github.com/petersalex27/yew/common/data"
)

func (ys *yewSource) acceptHeader(h data.Maybe[header]) {
	ys.header = h
	ys.Position = ys.Update(h)
}

func (ys *yewSource) acceptBody(b data.Maybe[body]) {
	ys.body = b
	ys.Position = ys.Update(b)
}

func (ys *yewSource) acceptFooter(f footer) {
	ys.footer = f
	ys.Position = ys.Update(f)
}

// Parse yew source file.
//
//	`yew source = {"\n"}, [meta, {"\n"}], [header, {"\n"}], [body, {"\n"}], footer ;`
func parseYewSource(p Parser) Parser {
	p.dropNewlines()
	ys := &yewSource{}

	// parse the yew source file surface from the top down
	p = p.
		bind(ys.parseHeader).
		bind(ys.parseBodyAndFooter).bind(assertEOF)

	// record the AST in the parser state on success
	if ps, ok := p.(*ParserState); ok {
		ps.ast = ys
	} else if pso, ok := p.(*ParserState_optional); ok {
		// suspicious, but okay ...
		pso.ParserState.ast = ys
	}

	return p
}

func assertEOF(p Parser) Parser {
	if !matchCurrent(token.EndOfTokens)(p) {
		return p.report(makeError(p, ExpectedEndOfFile, p), true)
	}
	return p
}

func (ys *yewSource) writeBody(p Parser, mb data.Maybe[body]) Parser {
	p.dropNewlines()
	ys.acceptBody(mb)
	return p
}

func (ys *yewSource) writeFooter(p Parser, m data.Maybe[annotations]) Parser {
	p.dropNewlines()
	ys.acceptFooter(footer{m})
	return p
}

func parseFooterAnnotsPostBody(ys *yewSource, p Parser, mFooterAnnots data.Maybe[annotations]) Parser {
	var es data.Ers
	valid := true

	// check if annotations were already parsed (non-nothing value for mFooterAnnots) or if EOF was
	// reached. In both cases, the footer has already been parsed; it's an annotation block or EOF
	// respectively.
	if mFooterAnnots.IsNothing() && !matchCurrentEndOfTokens(p) {
		// footer has not yet been parsed--this shouldn't happen in the way the function calls work
		// atm--but if newlines are not read before trying to parse the annotations while parsing
		// the body, then it's possible to miss the footer annotations and EOF.
		//
		// so, basically, this ensures if the logic changes for parsing the body, then the footer
		// annotations will still be parsed.
		//
		// TODO: remove this dependency without completely overhauling the parsing logic?
		p.dropNewlines()
		es, mFooterAnnots, valid = parseAnnotations(p).Break()
	}

	if !valid {
		return writeErrors(p, es)
	}

	return ys.writeFooter(p, mFooterAnnots)
}

func (ys *yewSource) parseBodyAndFooter(p Parser) Parser {
	res, mFooterAnnots := parseBody(p)
	es, b, isB := res.Break()
	if !isB {
		return writeErrors(p, es)
	}
	if b.IsNothing() {

	}
	p = ys.writeBody(p, b)

	return parseFooterAnnotsPostBody(ys, p, mFooterAnnots)
}
