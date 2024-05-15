// =================================================================================================
// Alex Peters - 2024
// =================================================================================================
package parser

import "github.com/petersalex27/yew/common/stack"

type termStack struct {
	*stack.SaveStack[termElem]
}

// creates term array and returns number of remaining terms to pop
//
// panics if a term is not popped from the stack (happens when stack frame is empty)
func (parser *Parser) initTermsPop(numTerms uint32) (terms []termElem, nRemaining int64) {
	if numTerms == 0 {
		return []termElem{}, 0
	}

	nRemaining = int64(numTerms)
	terms = make([]termElem, 1, numTerms)
	var stat stack.StackStatus
	if terms[0], stat = parser.terms.SaveStack.Pop(); stat.NotOk() {
		panic("bug: unexpected empty stack frame")
	}
	nRemaining--
	return terms, nRemaining
}

// panics if > 0 terms are requested but the stack is empty
//
// pops terms and returns them in the order they are popped in, throws an error if numTerms cannot be popped
func (parser *Parser) popTerms(numTerms uint32) (_ []termElem, ok bool) {
	terms, n := parser.initTermsPop(numTerms)
	ok = true
	
	for ; n > 0; n-- {
		term, stat := parser.terms.SaveStack.Pop()
		if stat.NotOk() {
			previous := len(terms)-1 // will always be >= 0 because at least one terms was added
			parser.errorOn(ExpectedTermAfter, terms[previous])
			return terms, false
		}

		terms = append(terms, term)
	}

	return terms, ok
}