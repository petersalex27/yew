// =================================================================================================
// Alex Peters - 2024
// =================================================================================================
package parser

import "os"

type resolutionFunc = func(parser *Parser) (term termElem, ok bool)

type resolutionMap map[NodeType]resolutionFunc

// attempts to lookup action given by NodeType `nt`
//
// on success, this returns the action and true; otherwise, returns _, false
func (data *actionData) findResolution(nt NodeType) (resolution resolutionFunc, found bool) {
	resolution, found = data.resolutions[nt]
	return resolution, found
}

var resolutionActions resolutionMap = resolutionMap{
	lambdaType: resolveLambdaAbstraction,
}

// given a previously processed (via normal actions or a previous resolution action) term `processed` ...
//	- pop (an empty) stack frame
//	- find and perform resolution based off of top type of node on current stack frame
//
// ASSUMPTION:
//	- stack frame being returned to is not empty 
func (parser *Parser) resolvingInner(data *actionData, processed termElem) (term termElem, ok bool) {
	var resolution resolutionFunc

	// return to save point
	parser.terms.Return()
	// use top term to decide on a resolution
	nt := parser.topTermType()
	resolution, ok = data.findResolution(nt)
	if !ok {
		// error: no resolution, but resolution required
		parser.reportUnresolved()
		return
	}
	// push processed term so action can use it
	parser.shift(processed)
	// run resolution
	term, ok = resolution(parser)
	return
}

// - if nothing to resolve: term arg, false, true is returned
// - if successful non-trivial resolution: new term, ?, true is returned
// - if not-successful resolution: _, false, false is returned
func (parser *Parser) resolving(data *actionData, term termElem) (_ termElem, again, ok bool) {
	again = true
	for again {
		if parser.terms.GetCount() != 0 {
			// current stack frame must be empty--call to `(*Parser) resolvingInner` will pop to the
			// next frame (this one is the one that gets "resolved")
			panic("bug: unexpected terms left on parse stack during resolution step")
		}
		
		if parser.terms.GetFrames() <= 1 {
			// nothing to resolve, but that's okay--actually done parsing
			return term, false, true
		}

		// resolve
		term, ok = parser.resolvingInner(data, term)

		// has the current stack frame been fully resolved?
		finishedFrame := parser.terms.GetCount() == 0
		// loop again if current frame is fully resolved and there are more frames to
		again = ok && finishedFrame // resolve again?
	}
	// process again?
	processAgain := ok && parser.terms.GetCount() != 0
	return term, processAgain, ok
}

func resolveLambdaAbstraction(parser *Parser) (term termElem, ok bool) {
	// pop [expression, lambda]
	var terms []termElem 
	if terms, ok = parser.popTerms(2); !ok {
		return
	}

	var bound termElem
	const lambdaIdx, boundIdx int = 1, 0
	term, bound = terms[lambdaIdx], terms[boundIdx]
	lambda := term.Term.(Lambda)
	
	// set lambda fields
	lambda.Bound = bound.Term
	_, lambda.End = bound.Pos()
	term.Term = lambda
	// remove local declarations
	if _, ok = parser.declarations.Decrease(); !ok {
		panic("bug: could not remove local declarations")
	}
	debug_log_reduce(os.Stderr, terms[lambdaIdx], terms[boundIdx], term)
	return term, ok
}