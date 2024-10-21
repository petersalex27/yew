package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/common/data"
	"github.com/petersalex27/yew/internal/kind"
)

// parser is the interface for parsing source code
type parser interface {
	// must keep track of the current position in the source code
	api.Positioned
	// return the source code
	srcCode() api.SourceCode
	// add the error, and if fatal, then return a fail state
	report(error, bool) parser
	// return token at the current position
	current() api.Token
	// advance the parser to the next token
	advance()
	// drop zero or more newlines
	dropNewlines()
	// declare a symbol
	declare(names api.Token, fix_prec ...byte)
	// create a new, free variable
	newVar() kind.Type
	newScope()
	exitScope() map[string]kind.Symbol
}

// initialize the parser, providing it with a scanner to read tokens from
func Init(scanner api.ScannerPlus) parser {
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
func Run(p parser) api.Node {
	p = parseYewSource(p)
	if ps, ok := p.(*ParserState); ok {
		return ps.ast
	}
	return nil
}

func Errors(p parser) (errors []error, warnings []error) {
	var ps *ParserState
	if ps_, ok := p.(*ParserState); ok {
		ps = ps_
	} else {
		ps = p.(*ParserStateFail).bad
	}
	return ps.errors, ps.warnings
}

func then(p parser) bool {
	origin := getOrigin(p)
	p.dropNewlines()
	return getOrigin(p) > origin
}

// returns the current token counter in the case of a ParserState or ParserState_optional instance,
// otherwise, returns -1
func getOrigin(p parser) int {
	if ps, ok := p.(*ParserState); ok {
		return ps.tokenCounter
	}
	return -1
}

// noop in the case of a ParserStateFail as Parser instance
func resetOrigin(p parser, origin int) parser {
	if ps, ok := p.(*ParserState); ok {
		ps.tokenCounter = origin
		p = ps
	}
	return p
}

func passParseErs[b api.Node](_ parser, x data.Ers) data.Either[data.Ers, b] {
	return data.PassErs[b](x)
}

func runCases[a, b api.Node, c any](p parser, disjointAct func(parser) data.Either[a, b], left func(parser, a) c, right func(parser, b) c) c {
	l, r, isR := disjointAct(p).Break()
	if isR {
		return right(p, r)
	}
	return left(p, l)
}
