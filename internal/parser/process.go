package parser

import (
	"github.com/petersalex27/yew/internal/common"
	"github.com/petersalex27/yew/internal/token"
)

// get position span of terms `a` and `b`; the order the values are given to the function is irrelevant
func getTermsPos(a, b positioned) (start, end int) {
	// figure out position--it may be in reverse order (a later, b sooner) because of infix ops
	start, end = a.Pos()
	start2, end2 := b.Pos()
	start = common.Min(start, start2)
	end = common.Max(end, end2)
	return
}

// stores the argument in the parser's "term-memory"
func (parser *Parser) remember(term termElem) {
	parser.termMemory.Push(term)
}

// forgets the term the parser is remembering
func (parser *Parser) forget() {
	if parser.termMemory.Empty() {
		return
	}

	_, _ = parser.termMemory.Pop()
}

func (parser *Parser) terminalPeek(data *actionData) (termElem, bool) {
	t, ok := parser.actOnTerminal(data)
	if !ok {
		return termElem{}, false
	}

	parser.remember(t)
	return t, true
}

func (term termElem) strongerThan(bp int8) bool {
	if term.rAssoc {
		return bp <= term.bp
	}
	return bp < term.bp
}

func (parser *Parser) earlyExit(data *actionData) bool {
	if !parser.termMemory.Empty() {
		return false
	}
	res, _ := parser.peek(data)
	return res.String() == data.end
}

func (parser *Parser) application(bp int8, left termElem, data *actionData) (term termElem, ok bool) {
	if !left.strongerThan(bp) {
		return left, true // bp of left is too low to bind to anything
	}

	for parser.keepProcessing(data) && left.arity != 0 {
		var right termElem
		if right, ok = parser.process(left.bp, data); !ok {
			return left, false
		}

		if left, ok = parser.reduce(left, right); !ok {
			return right, false
		}
	}

	return left, true
}

func (parser *Parser) terminalsLeft(data *actionData) bool {
	some := data.ptr < uint(len(data.tokens)) || !parser.termMemory.Empty()
	if !some && data.supplyNext != nil {
		some = data.supplyNext(parser, data)
	}
	return some
}

func (parser *Parser) peekAtInfix(data *actionData) (term termElem, ok bool) {
	if term, ok = parser.terminalPeek(data); !ok {
		return
	}

	return term, term.infixed
}

func (parser *Parser) prefix(bp int8, data *actionData) (left termElem, ok bool) {
	if left, ok = parser.actOnTerminal(data); !ok {
		return
	} else if left, ok = parser.application(bp, left, data); !ok {
		return
	}

	return left, ok
}

// use parser.Panicking() to check for errors
func (parser *Parser) infix(bp int8, left termElem, data *actionData) (_ termElem, again bool) {
	// check if infix is next
	var right termElem
	if right, again = parser.peekAtInfix(data); !again {
		// okay: not an infixed function
		return left, false
	} else if parser.earlyExit(data) {
		// okay: exiting early
		return left, false
	}

	if !right.strongerThan(bp) {
		// okay: binding power of `right` is too weak to be used right now
		return left, false
	}

	// guaranteed to succeed, get remembered infix operator
	right, _ = parser.actOnTerminal(data)
	// swap order of infix operator and argument, then reduce
	if left, again = parser.reduce(right, left); !again {
		return // error: failed to reduce left and right
	}

	// now, treat infix operator as regular function (with a non-10 valued bp)
	return parser.application(bp, left, data)
}

func (parser *Parser) keepProcessing(data *actionData) bool {
	return parser.terminalsLeft(data) && !parser.earlyExit(data)
}

// processes a function (possibly w/ 0 args) then a possible sequence of expressions sequenced by 
// infix operators with enough binding power to bind stronger than `bp` 
func (parser *Parser) process(bp int8, data *actionData) (term termElem, ok bool) {
	var left termElem
	if left, ok = parser.prefix(bp, data); !ok {
		return
	}

	// TODO: what to do about implicit arguments? 

	for again := true; again && parser.keepProcessing(data); {
		left, again = parser.infix(bp, left, data)
	}
	return left, !parser.panicking
}

func (parser *Parser) reportProcessErrors(data *actionData) {
	if !parser.termMemory.Empty() {
		t, _ := parser.actOnTerminal(data)
		parser.errorOn(ExpectedEndOfSection, t)
		return
	}
	t, _ := parser.peek(data)
	parser.errorOnToken(UnexpectedToken, t)
}

func (tok TokensElem) parseExpressionPart(parser *Parser, data *actionData) (end, ok bool) {
	data.tokens = tok
	return true, true
}

func (let LetBindingElem) parseExpressionPart(parser *Parser, data *actionData) (end, ok bool) {
	if !let.Parse(parser) {
		return true, false
	}
	return false, true
}

func (parser *Parser) expressionProcessIteration(i int, expressionParts []ExpressionElem, data *actionData) (new_i int, end, ok bool) {
	if i >= len(expressionParts) {
		// this is fine, signals end of expression
		return i, true, true
	}

	end, ok = expressionParts[i].parseExpressionPart(parser, data)
	if !ok {
		return i, true, false
	}

	// if not end, then end if at end of expressionParts
	end = end || i >= len(expressionParts) - 1

	// if let binding, then must be followed by an expression
	if let, isLet := expressionParts[i].(LetBindingElem); end && isLet {
		// error: let binding is not followed by an expression
		parser.errorOn(ExpectedExpression, let)
		return i, true, false
	}

	i++
	return i, end, true
}
	

func (parser *Parser) ProcessExpression(a actionMapper, expressionParts []ExpressionElem) (term termElem, ok bool) {
	i := 0
	supplyNext := func(p *Parser, data *actionData) bool {
		if i >= len(expressionParts) {
			// not good, no more tokens to supply
			if len(expressionParts) == 0 {
				panic("illegal expressionParts: no tokens to supply")
			}

			parser.errorEOI(expressionParts[len(expressionParts)-1])
			return false
		}
		var end bool
		// loops until tokens are refilled or end of expressionParts (or error)
		for !end {
			i, end, ok = p.expressionProcessIteration(i, expressionParts, data)
		}
		return ok
	}

	data := newActionDataWithDischarger(a, nil)
	data.supplyNext = supplyNext
	// get first supply of tokens
	if !supplyNext(parser, data) {
		return termElem{}, false
	}

	// process the expression like the other Process functions
	if term, ok = parser.process(-1, data); !ok {
		return
	}
	if !parser.terminalsLeft(data) {
		return term, ok
	}

	parser.reportProcessErrors(data)
	return term, false
}

func (parser *Parser) Process(a actionMapper, tokens []token.Token) (term termElem, ok bool) {
	// create an action data struct for tokens and map
	data := newActionData(a, tokens)
	if term, ok = parser.process(-1, data); !ok {
		return
	}
	if !parser.terminalsLeft(data) {
		return term, ok
	}

	parser.reportProcessErrors(data)
	return term, false
}

type processedValidation struct {
	// valid terms at the end of processing
	validEndTerms []NodeType
	// error message when none of the terms are found
	getErrorMessage func(termElem) string
}

func (parser *Parser) ProcessAndValidate(a actionMapper, tokens []token.Token, pv processedValidation) (term termElem, ok bool) {
	if term, ok = parser.Process(a, tokens); !ok {
		return
	}

	actual := term.NodeType
	for _, ty := range pv.validEndTerms {
		if actual == ty {
			return term, true
		}
	}

	msg := UnexpectedSection
	if pv.getErrorMessage != nil {
		msg = pv.getErrorMessage(term)
	}

	parser.errorOn(msg, term)
	return term, false //
}
