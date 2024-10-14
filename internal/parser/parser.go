package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/common/data"
)

// Parser is the interface for parsing source code
type Parser interface {
	// must keep track of the current position in the source code
	api.Positioned
	// return the source code
	srcCode() api.SourceCode
	// only allow the parser to continue to 'f' if the parser is not in a bad state
	bind(f func(Parser) Parser) Parser
	// add the error, and if fatal, then return a fail state
	report(error, bool) Parser
	// return token at the current position
	current() api.Token
	// advance the parser to the next token
	advance()
	// drop zero or more newlines
	dropNewlines()
	// Mark whatever is parsed next as optional.
	//
	// An error is not to be reported when a parse attempt fails and the parser is in an
	// optional-marked-state. Instead, it should gracefully return the parsing to the point where
	// the optional-mark was made.
	markOptional() Parser
	// Stop parsing as if the next thing parsed was optional.
	//
	// This is *NOT* necessarily the inverse of `markOptional`; it's its own state transition.
	// Though, it works closely with `markOptional` to provide the optional parsing feature. And,
	// under ALL circumstances, it should be paired with a prior `markOptional` call.
	demarkOptional() Parser
	acceptRoot(yewSource) Parser
}

// initialize the parser, providing it with a scanner to read tokens from
func Init(scanner api.ScannerPlus) Parser {
	ps := &ParserState{
		state: createState(scanner),
		ast:   makeEmptyYewSource(),
	}
	if !ps.load() {
		return &ParserStateFail{bad: ps}
	}

	return ps
}

// Run an initialized parser
//
// SEE: `Init`
func Run(p Parser) api.Node {
	parseYewSource(p)
	return nil
}

// returns the current token counter in the case of a ParserState or ParserState_optional instance,
// otherwise, returns -1
func getOrigin(p Parser) int {
	if ps, ok := p.(*ParserState); ok {
		return ps.tokenCounter
	} else if ps, ok := p.(*ParserState_optional); ok {
		return ps.tokenCounter
	}
	return -1
}

// noop in the case of a ParserStateFail as Parser instance
func resetOrigin(p Parser, origin int) Parser {
	if ps, ok := p.(*ParserState); ok {
		ps.tokenCounter = origin
		p = ps
	} else if ps, ok := p.(*ParserState_optional); ok {
		ps.tokenCounter = origin
		p = ps
	}
	return p
}

func passParseErs[b api.Node](_ Parser, x data.Ers) data.Either[data.Ers, b] {
	return data.PassErs[b](x)
}

func runCases[a, b api.Node, c any](p Parser, disjointAct func(Parser) data.Either[a, b], left func(Parser, a) c, right func(Parser, b) c) c {
	l, r, isR := disjointAct(p).Break()
	if isR {
		return right(p, r)
	}
	return left(p, l)
}
