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
	// add the error, and if fatal, then return a fail state
	report(error, bool) Parser
	// return token at the current position
	current() api.Token
	// advance the parser to the next token
	advance()
	// drop zero or more newlines
	dropNewlines()
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

func then(p Parser) bool {
	origin := getOrigin(p)
	p.dropNewlines()
	return getOrigin(p) > origin
}

// returns the current token counter in the case of a ParserState or ParserState_optional instance,
// otherwise, returns -1
func getOrigin(p Parser) int {
	if ps, ok := p.(*ParserState); ok {
		return ps.tokenCounter
	}
	return -1
}

// noop in the case of a ParserStateFail as Parser instance
func resetOrigin(p Parser, origin int) Parser {
	if ps, ok := p.(*ParserState); ok {
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
