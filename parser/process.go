package parser

import (
	"fmt"
	"os"

	"github.com/petersalex27/yew/common"
	"github.com/petersalex27/yew/token"
)

func getTermsPos(a, b Term) (start, end int) {
	// figure out position--it may be in reverse order (a later, b sooner) because of infix ops
	start, end = a.Pos()
	start2, end2 := b.Pos()
	start = common.Min(start, start2)
	end = common.Max(end, end2)
	return
}

func reduceFuncType(f FunctionType, fi termInfo, b termElem) (termElem, bool) {
	start, end := getTermsPos(f, b)
	f.Start, f.End = start, end

	switch fi.arity {
	case 1:
		f.Left = b.Term
	default:
		f.Right = b.Term
	}
	return termElem{f, fi}, true
}

func (a termElem) reduce(b termElem) (termElem, bool) {
	newInfo, decd := a.termInfo.decrementArity()
	if f, ok := a.Term.(FunctionType); ok && decd {
		return reduceFuncType(f, newInfo, b)
	}

	start, end := getTermsPos(a, b)

	return termElem{Application{a.Term, b.Term, start, end}, newInfo}, decd
}

func reduceImmediately(left termElem, right termElem) bool {
	if right.termInfo.AssociatesRight() {
		return right.Bp() >= left.Bp()
	}
	return right.Bp() > left.Bp()
}

// applies top element of stack to `fun`, i.e.,
//
//	push stack ((\e -> fun e) stack.top)
func (parser *Parser) reduceTop(fun termElem) (ok bool) {
	term := parser.grab()
	if term, ok = parser.reduceTerms(fun, term); ok {
		parser.shift(term)
	}
	return
}

func (parser *Parser) grab() termElem {
	term, stat := parser.terms.Pop()
	if stat.NotOk() {
		parser.reportUnresolved()
	}
	return term
}

func (parser *Parser) reduceTerms(a, b termElem) (term termElem, ok bool) {
	if term, ok = a.reduce(b); !ok {
		start, _ := a.Pos()
		_, end := b.Pos()
		parser.error2(IllegalApplication, start, end)
	}
	fmt.Fprintf(os.Stderr, "red: (%v) (%v) = %v\n", a, b, term)
	return
}

func (parser *Parser) reduceStack() (term termElem, ok bool) {
	// empty stack via reduction
	term = parser.grab()
	if ok = !parser.panicking; !ok {
		return
	} else if parser.terms.Empty() {
		return term, true
	}

	fun := parser.grab()
	for {
		term, ok = parser.reduceTerms(fun, term)
		if !ok || parser.terms.Empty() {
			break
		}
		fun = parser.grab()
	}
	return
}

func (parser *Parser) shift(term termElem) {
	parser.terms.Push(term)
	fmt.Fprintf(os.Stderr, "shf: %s\n", parser.terms.ElemString())
}

func (parser *Parser) shiftTerm(data *actionData) (ok bool) {
	var term termElem
	// get and push term
	if term, ok = parser.terminalAction(data); ok {
		parser.shift(term)
	}
	return
}

func (parser *Parser) remember(term termElem) {
	parser.termMemory = &term
}

func (parser *Parser) forget() {
	parser.termMemory = nil
}

func (parser *Parser) rightTakesTop(fun termElem, data *actionData) (term termElem, ok bool) {
	//var leftOp termElem
	for again := true; again; {
		// reduce by applying left term to infix function on the right
		if ok = parser.reduceTop(fun); !ok {
			return
		}

		result := parser.top()

		// exit if out of input
		if !data.hasMoreInput() {
			break
		}

		if ok = parser.shiftTerm(data); !ok {
			return
		}

		// exit if out of input
		if !data.hasMoreInput() {
			break
		}

		// get term to the right of the term gotten above (this term might be an infix function)
		fun, ok = parser.terminalAction(data)
		parser.remember(fun)
		// should right infix function take the term on the top of the stack?
		again = ok && reduceImmediately(result, fun)
		if again {
			parser.forget()
		}
	}

	if !ok {
		return
	}

	return parser.reduceStack()
}

func (parser *Parser) innerProcess(data *actionData) (term termElem, ok bool) {
	// check if nothing to parse
	if ok = parser.terms.GetCount() == 1 && !data.hasMoreInput(); ok { 
		term, _ = parser.terms.Pop()
		return
	}

	var rightOp termElem
	rightOp, ok = parser.terminalAction(data)
	if !ok {
		return
	}

	if reduceImmediately(parser.top(), rightOp) {
		// apply left to right (infix application)
		return parser.rightTakesTop(rightOp, data)
	}

	// apply right to left (normal application)
	term = parser.grab()
	if ok = !parser.panicking; !ok {
		return
	}
	return parser.reduceTerms(term, rightOp)
}

func (parser *Parser) processing(data *actionData) (term termElem, ok bool) {
	ok = true
	for ok && !parser.terms.Empty() {
		term, ok = parser.innerProcess(data)
		if ok && data.hasMoreInput() {
			parser.shift(term) // signals that loop should happen again
		}
	}
	return
}

func (parser *Parser) topTermType() NodeType {
	term, stat := parser.terms.Peek()
	if stat.NotOk() {
		panic("bug: I think? Some parse rule created an empty stack frame ...")
	}

	return term.NodeType()
}

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
}

func (parser *Parser) resolvingInner(data *actionData, processed termElem) (term termElem, ok bool) {
	var action actionFunc

	// return to save point
	parser.terms.Return()
	// use top term to decide on a resolution
	nt := parser.topTermType()
	action, ok = data.findResolution(nt)
	if !ok { 
		// error: no resolution, but resolution required
		parser.reportUnresolved()
		return
	}
	// push processed term so action can use it
	parser.terms.Push(processed)
	// run action
	term, ok = action(parser, data)
	return
} 

//	- if nothing to resolve: term arg, false, true is returned
// 	- if successful non-trivial resolution: new term, ?, true is returned
//	- if not-successful resolution: _, false, false is returned
func (parser *Parser) resolving(data *actionData, term termElem) (_ termElem, again, ok bool) {
	again = true
	for again {
		if parser.terms.GetCount() != 0 {
			panic("bug: unexpected terms left on parse stack during resolution step")
		} else if parser.terms.GetFullCount() == 0 {
			return term, false, true
		} else if parser.terms.GetFrames() <= 1 {
			// nothing to resolve, but that's okay--actually done parsing
			return term, false, true
		}

		// resolve
		term, ok = parser.resolvingInner(data, term)
	
		// resolve again?
		again = ok && parser.terms.GetCount() == 0
	}
	// process again?
	again = ok && parser.terms.GetCount() != 0
	return term, again, ok
}

// DO NOT CALL THIS FUNCTION RECURSIVELY!!! Specifically, `rightTakesTop` cannot be called recursively
func (parser *Parser) Process(a actionMap, tokens []token.Token) (term termElem, ok bool) {
	data := newActionData(a, resolutionActions, tokens)
	if ok = parser.shiftTerm(data); !ok {
		return
	}

	// loop until stack is empty and there are no resolutions
	again := true
	for again && ok {
		term, ok = parser.processing(data)
		// attempt to resolve
		if ok {
			term, again, ok = parser.resolving(data, term)
		}
	}
	return
}
