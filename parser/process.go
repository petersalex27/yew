package parser

import (
	"os"

	"github.com/petersalex27/yew/common"
	"github.com/petersalex27/yew/token"
)

// get position span of terms `a` and `b`; the order the values are given to the function is irrelevant
func getTermsPos(a, b Term) (start, end int) {
	// figure out position--it may be in reverse order (a later, b sooner) because of infix ops
	start, end = a.Pos()
	start2, end2 := b.Pos()
	start = common.Min(start, start2)
	end = common.Max(end, end2)
	return
}

// like `(*Parser).grabNext()` but also reports an error with when the grabbed term
// does not have the type `ty`
//
// on success, returns grabbed term of type `ty` and true
//
// on failure, returns ?, false
func (parser *Parser) grabTermOf(ty NodeType) (term termElem, ok bool) {
	if term, ok = parser.grabNext(); !ok {
		return
	}

	// test type
	if ok = term.NodeType() == ty; !ok {
		parser.expectedErrorOn(ty, term)
	}
	return
}

// grabs the top term of the term stack, reports errors if they occur and returns the grabbed term
// and whether this was successful (ok=true iff successful)
func (parser *Parser) grabNext() (term termElem, ok bool) {
	term = parser.grab()
	ok = !parser.panicking
	return
}

// grabs the top term of the term stack and returns it, reports errors if they occur
//
// An error occurred if parser.Panicking() returns false before calling `grab` but true after
// calling
func (parser *Parser) grab() termElem {
	term, stat := parser.terms.SaveStack.Pop()
	if stat.NotOk() {
		parser.reportUnresolved()
	}
	return term
}

// pushes a term onto the term stack--uses the name 'shift' to align w/ 'shift-reduce' parser
func (parser *Parser) shift(term termElem) {
	parser.terms.Push(term)
	debug_log_shift(os.Stderr, parser) // noop if not -tags=debug
}

// performs a terminal action and immediately shifts the result if the terminal action was
// successful
//
// returns true iff terminal action was successful
func (parser *Parser) shiftTerm(data *actionData) (ok bool) {
	var term termElem
	// get and push term
	if term, ok = parser.terminalAction(data); ok {
		parser.shift(term)
	}
	return
}

func (parser *Parser) postReductionShift(data *actionData) (reductionResult termElem, ok bool) {
	reductionResult = parser.top()

	var term termElem
	// get and push term
	if term, ok = parser.terminalAction(data); ok {
		parser.shift(term)
	}
	return
}

// stores the argument in the parser's "term-memory"
func (parser *Parser) remember(term termElem) { parser.termMemory = &term }

// forgets the term the parser is remembering
func (parser *Parser) forget() { parser.termMemory = nil }

func (parser *Parser) rightTakesTopIteration(fun termElem, data *actionData) (_ termElem, again bool) {
	var ok bool
	// reduce by applying left term to infix function on the right
	if ok = parser.reduceTop(fun); !ok {
		return fun, false
	}

	// return top (just reduced) term if out of input
	if !data.hasMoreInput() {
		term, _ := parser.reduceStack()
		return term, false
	}

	// result of reduction
	result := parser.top()

	var nextTerm termElem
	if nextTerm, ok = parser.terminalAction(data); !ok {
		return fun, false
	}

	if shouldHoldStack(fun, nextTerm) {
		
	} 

	if !data.hasMoreInput() {
		// reduce entire stack
		term, _ := parser.reduceStack() 
		return term, false
	}

	// get term to the right of the term gotten above (this term might be an infix function)
	fun, ok = parser.terminalAction(data)
	parser.remember(fun)
	// should right infix function take the term on the top of the stack?
	//	- decision made based on result of previous reduction
	again = ok && shouldHoldStack(result, fun)
	if !again {
		// TODO?
		res, _ := parser.reduceStack()
		return res, false
	}
	parser.forget()
	return fun, again
}

func (parser *Parser) rightTakesTop(fun termElem, data *actionData) (_ termElem, ok bool) {
	for again := true; again; {
		fun, again = parser.rightTakesTopIteration(fun, data)
	}
	return fun, !parser.Panicking()
}

func (parser *Parser) innerProcess(data *actionData) (term termElem, ok bool) {
	// check if nothing to parse
	if ok = parser.terms.GetCount() == 1 && !data.hasMoreInput(); ok {
		term = parser.grab()
		return
	}

	var rightOp termElem
	rightOp, ok = parser.terminalAction(data)
	if !ok {
		return
	}

	if parser.terms.GetCount() == 0 {
		return rightOp, true
	}

	if shouldHoldStack(parser.top(), rightOp) {
		// apply left to right (infix application)
		return parser.rightTakesTop(rightOp, data)
	}

	// apply right to left (normal application)
	term = parser.grab()
	if ok = !parser.panicking; !ok {
		return
	}
	return parser.reduce(term, rightOp)
}

func (parser *Parser) processing(data *actionData) (term termElem, ok bool) {
	// require at least one iteration through processing loop
	if parser.terms.Empty() {
		panic("bug: term must be on the parse stack before calling `(*Parser) process`")
	}

	// processing loop
	for again := true; again && !parser.terms.Empty(); {
		term, ok = parser.innerProcess(data)
		again = ok
		if ok && data.hasMoreInput() {
			parser.shift(term) // signals that loop should happen again
		}
	}
	return term, ok
}

func (parser *Parser) topTermType() NodeType {
	term, stat := parser.terms.Peek()
	if stat.NotOk() {
		panic("bug: an unknown parse rule created an empty stack frame")
	}

	return term.NodeType()
}

// reports errors based on whatever is left on the term stack
func (parser *Parser) reportUnresolved() {
	reported := 0
	for {
		n := parser.terms.GetCount()
		elems, stat := parser.terms.MultiPop(n)
		if !stat.IsOk() {
			break // done reporting
		}

		for _, elem := range elems {
			start, end := elem.Pos()
			parser.error2(ReductionFailure, start, end)
			reported++
		}

		if !parser.terms.Return().IsOk() {
			break // done reporting
		}
	}

	if reported == 0 {
		panic("unknown error caused parser to fail")
	}

	parser.panicking = true
}

// DO NOT CALL THIS FUNCTION RECURSIVELY!!! Specifically, `rightTakesTop` cannot be called recursively
func (parser *Parser) Process(a actionMapper, tokens []token.Token) (term termElem, ok bool) {
	// create an action data struct for tokens and map
	data := newActionData(a, resolutionActions, tokens)
	// shift, requiring at least one term
	if ok = parser.shiftTerm(data); !ok {
		return
	}

	// loop until stack is empty and there are no resolutions
	for again := true; again && ok; {
		// process data
		term, ok = parser.processing(data)
		if ok {
			// resolve anything that needs resolution
			if term, again, ok = parser.resolving(data, term); !ok {
				break
			}
		}
	}
	return
}

type processedValidation struct {
	// valid terms at the end of processing
	validEndTerms []NodeType
	// error message when none of the terms are found
	getErrorMessage func(Term) string
}

func (parser *Parser) ProcessAndValidate(a actionMapper, tokens []token.Token, pv processedValidation) (term termElem, ok bool) {
	if term, ok = parser.Process(a, tokens); !ok {
		return
	}

	actual := term.Term.NodeType()
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
