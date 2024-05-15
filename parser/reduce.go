// =================================================================================================
// Alex Peters - May 05, 2024
//
// part of parser's second-pass
//
// functions and methods related to parse stack reductions
// =================================================================================================
package parser

import (
	"os"

	"github.com/petersalex27/yew/token"
)

type reduceFunction func(parser *Parser, a, b termElem) (termElem, bool)

var reduceTable map[NodeType]reduceFunction

var _reduceParens = closingReduction(standardParenCloser, token.LeftParen)

// var _reduceBrackets =  //closingReduction()
var _reduceBraces = closingReduction(closingImplicit, token.LeftBrace)

func init() {
	reduceTable = map[NodeType]reduceFunction{
		functionType:     reduceFuncType,
		listingType:      reduceListing,
		typingType:       reduceTyping,
		closeParenType:   reduceParens,
		closeBracketType: reduceBrackets,
		closeBraceType:   reduceBraces,
	}
}

func defaultReduce(parser *Parser, a, b termElem) (termElem, bool) {
	newInfo, decd := a.termInfo.decrementArity()
	start, end := getTermsPos(a, b)
	term := termElem{Application{a.Term, b.Term, start, end}, newInfo}
	
	if !decd {
		parser.error2(IllegalApplication, start, end)
		return termElem{}, false
	}
	return term, true
}

func (parser *Parser) declareIdFromTyping(typing Typing, b termElem) {
	if !parser.parsingTypeSig {
		return
	}

	id, isId := typing.Term.(Ident)
	if !isId {
		return
	}

	decl := new(Declaration)
	decl.name = id
	generate_setType(id.Name, decl)(b.Term)
	parser.locals.Map(id, decl) // will overwrite any names shadowed here
}

func reduceParens(parser *Parser, paren, term termElem) (termElem, bool) {
	return _reduceParens(parser, paren, term)
}

func reduceBrackets(parser *Parser, bracket, term termElem) (termElem, bool) {
	panic("TODO: implement")
}

func reduceBraces(parser *Parser, brace, term termElem) (termElem, bool) {
	return _reduceBraces(parser, brace, term)
}

func reduceTyping(parser *Parser, a, b termElem) (termElem, bool) {
	tyi, decd := a.termInfo.decrementArity() // always use 'a'
	if !decd {
		return termElem{}, false
	}

	typing := a.Term.(Typing)

	// if arity after decrementing it is 0, then a type is being declared for the term
	// 	- otherwise, the term is being given
	typingTerm := tyi.arity == 0
	if typingTerm {
		parser.declareIdFromTyping(typing, b)
	}

	typing.Start, typing.End = getTermsPos(typing, b)
	if typingTerm {
		typing.Type = b.Term
	} else {
		typing.Term = b.Term
	}

	return termElem{typing, tyi}, true
}

// copies available implicit bindings
func (parser *Parser) getAvailable() (available *declTable) {
	available = new(declTable)
	for _, local := range parser.locals.All() {
		if local.Value.implicit {
			available.Map(local.Key, local.Value)
		}
	}
	return available
}

func reduceFuncType(parser *Parser, a, b termElem) (termElem, bool) {
	fi, decd := a.termInfo.decrementArity()
	if !decd {
		return termElem{}, false
	}

	f := a.Term.(FunctionType)

	f.Start, f.End = getTermsPos(f, b)

	if fi.arity == 1 {
		f.Left = b.Term
		f.availableInDefs = parser.getAvailable()
	} else {
		f.Right = b.Term
	}
	return termElem{f, fi}, true
}

func reduceListing(parser *Parser, a, b termElem) (termElem, bool) {
	// whether or not to append a list
	// if false, then `a` will have a single element inserted
	isAppend := b.Term.NodeType() == listingType

	li, decd := a.termInfo.decrementArity()
	if !decd {
		return termElem{}, false
	}
	ls := a.Term.(Listing)

	start, end := getTermsPos(ls, b)
	ls.Start, ls.End = start, end

	if isAppend {
		ls2 := b.Term.(Listing)
		ls.Elements = append(ls.Elements, ls2.Elements...)
	} else {
		ls.Elements = append(ls.Elements, b.Term)
	}

	return termElem{ls, li}, true
}

func selectReduction(t Term) reduceFunction {
	rf := reduceTable[t.NodeType()]
	if rf == nil {
		rf = defaultReduce
	}
	return rf
}

// perform a reduction based off of the type of `a`
func (parser *Parser) reduce(a, b termElem) (termElem, bool) {
	reduction := selectReduction(a)
	term, ok := reduction(parser, a, b)

	debug_log_reduce(os.Stderr, a, b, term)

	return term, ok
}

// predicate, checks if reduction should happen after calling this function
func shouldHoldStack(left termElem, right termElem) bool {
	// if !right.infixed {

	// }
	if right.termInfo.AssociatesRight() {
		return right.Bp() >= left.Bp()
	}
	return right.Bp() > left.Bp()
}

// applies top element of stack to `fun`, i.e.,
//
//	push stack ((\e => fun e) stack.top)
func (parser *Parser) reduceTop(fun termElem) (ok bool) {
	term := parser.grab()
	if ok = !parser.panicking; !ok {
		return
	}
	if term, ok = parser.reduce(fun, term); ok {
		parser.shift(term)
	}
	return
}

func (parser *Parser) reduceStack() (_ termElem, ok bool) {
	// empty stack via reduction
	term := parser.grab()
	if ok = !parser.panicking; !ok {
		return termElem{}, false
	} else if parser.terms.Empty() {
		return term, true
	}

	fun := parser.grab()
	for {
		term, ok = parser.reduce(fun, term)
		if !ok || parser.terms.Empty() {
			break
		}
		fun = parser.grab()
	}
	return term, ok
}

func closingReduction(closer parenClosingFunc, opened token.Type) reduceFunction {
	return func(parser *Parser, marker, b termElem) (_ termElem, ok bool) {

		var lp, term termElem
		// loop reducing anything that needs to be reduced
		for {
			if term, ok = parser.reduceStack(); !ok {
				return
			}
			if ok = !parser.terms.FullEmpty(); !ok {
				parser.errorOn(UnexpectedRParen, marker) // TODO: this shouldn't be 'UnexpectedRParen'
				return
			}

			parser.terms.Return() // return stack
			if ok = parser.terms.GetCount() != 0; !ok {
				// error: size of stack frame is 0
				parser.errorOn(UnexpectedRParen, marker) // TODO: this shouldn't be 'UnexpectedRParen'
				return
			}

			lp, _ = parser.terms.Peek()
			if lp.NodeType() != lambdaType {
				break
			}

			// push reduced bound expression
			parser.shift(term)
			// resolve lambda abstraction
			term, ok = resolveLambdaAbstraction(parser)
			if !ok {
				return
			}
			parser.shift(term) // push result, try again
		}

		// at this point, lp should be a left paren
		start, end := lp.Pos()
		isOpened := lp.Term.NodeType() == syntaxExtensionType && lp.Term.String() == opened.Make().Value
		if ok = isOpened; !ok {
			parser.error2(expectedMessage(opened), start, end)
			return
		}

		// remove left paren
		_, _ = parser.terms.SaveStack.Pop()

		_, end = marker.Pos() // end of right paren

		// create enclosed term
		return closer(parser, term, start, end)
	}
}

func closeTuple(term termElem, start, end int) (termElem, bool) {
	ls := term.Term.(Listing)
	tuple := Tuple(ls)
	tuple.Start, tuple.End = start, end
	term = termElem{tuple, term.termInfo}
	return term, true
}

func closeClosedTuple(term termElem, start, end int) (termElem, bool) {
	// update start and end--enclosing in parens does nothing otherwise
	ps := term.Term.(Tuple)
	ps.Start = start
	ps.End = end
	term.Term = ps
	return term, true
}

func standardParenCloser(parser *Parser, term termElem, start, end int) (termElem, bool) {
	ty := term.Term.NodeType()
	if ty == listingType {
		return closeTuple(term, start, end)
	} else if ty == pairsType {
		// TODO: need a way to prevent the same warning from being thrown when a tuple is nested in more
		// than one extra pair of parens. Whatever the solution is, should be general enough for
		// arbitrary rules that can throw multiple similar warnings in succession
		parser.warning2(ExcessiveParens, start, end) // ((..(e1, e2, .., eN)..))
		return closeClosedTuple(term, start, end)
	}

	info := abstractInfo(term.termInfo)
	term = termElem{EnclosedTerm{term.Term, start, end}, info}
	return term, true
}