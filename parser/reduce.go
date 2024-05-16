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
)

type reduceFunction func(parser *Parser, a, b termElem) (termElem, bool)

var reduceTable map[NodeType]reduceFunction

func init() {
	reduceTable = map[NodeType]reduceFunction{
		functionType:     reduceFuncType,
		listingType:      reduceListing,
		typingType:       reduceTyping,
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

func rAssoc(t termElem) uint8 {
	if t.rAssoc {
		return 1
	}
	return 0
}

func (parser *Parser) declareIdFromTyping(typing Typing, b termElem) {
	id, isId := typing.Term.(Ident)
	if !isId {
		return
	}

	decl := new(Declaration)
	decl.name = id
	generate_setType(decl)(b.Term, b.infixed)
	parser.locals.Map(id, decl) // will overwrite any names shadowed here
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

	typing.Start, typing.End = getTermsPos(typing, b)
	if typingTerm {
		typing.Type = b.Term
		parser.declareIdFromTyping(typing, b)
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
