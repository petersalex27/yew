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
	"strings"

	"github.com/petersalex27/yew/internal/token"
	"github.com/petersalex27/yew/internal/types"
)

type reduceFunction func(parser *Parser, a, b termElem) (termElem, bool)

var reduceTable map[NodeType]reduceFunction

func init() {
	reduceTable = map[NodeType]reduceFunction{
		functionType:        reduceFuncType,
		labeledFunctionType: reduceFuncType,
		listingType:         reduceListing,
		typingType:          reduceTyping,
	}
}

func defaultReduce(parser *Parser, a, b termElem) (termElem, bool) {
	newInfo, decd := a.termInfo.decrementArity()
	start, end := getTermsPos(a, b)
	x, A := a.Term, types.GetKind(&a.Term)
	if parser.panicking {
		return termElem{}, false
	}
	y := b.Term
	if parser.panicking {
		return termElem{}, false
	}
	z, C, ok := parser.env.Apply(x, A, y)
	if !ok {
		parser.transferEnvErrors()
		return termElem{}, false
	}
	if ok = types.SetKind(&z, C); !ok {
		parser.transferEnvErrors()
		return termElem{}, false
	}
	term := termElem{applicationType, z, newInfo, start, end}

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

func (parser *Parser) declareIdFromTyping(typing termElem) bool {
	implicit := typing.NodeType == implicitType
	id, isId := typing.Term.(stringPos)
	if !isId {
		return true // okay?
	}

	s := id.String()
	if s == "_" || s == "" {
		return true // don't declare the wildcard, but this isn't an error
	}

	// infix-id typing is not allowed
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		parser.errorOn(IllegalInfixTyping, typing)
		return false
	}

	setter, ok := parser.declare(token.Id.MakeValued(s), implicit, exports{})
	if !ok {
		return false // error already reported in declare
	}

	ty := types.GetKind(&typing.Term)
	if !setter(parser, ty, false) {
		return false // error already reported in setter
	}

	if !parser.env.Declare(id, ty) {
		parser.transferEnvErrors()
		return false
	}

	arity := types.CalcArity(ty)
	lambda := generateAssignableTerm(id, arity)
	if !parser.env.Assign(id, lambda) {
		parser.transferEnvErrors()
		return false
	}
	return true
}

func reduceTyping(parser *Parser, a, b termElem) (termElem, bool) {
	tyi, decd := a.termInfo.decrementArity() // always use 'a'
	if !decd {
		return termElem{}, false
	}

	// if arity after decrementing it is 0, then a type is being declared for the term
	// 	- otherwise, the term is being given
	typingTerm := tyi.arity == 0

	start, end := getTermsPos(a, b)
	var term termElem
	if typingTerm {
		B, ok := b.Term.(types.Type)
		if !ok {
			parser.errorOn(ExpectedType, b)
			return termElem{}, false
		}
		res := a.Term
		A := types.GetKind(&res)
		if ok = parser.env.Unify(A, B); !ok {
			parser.transferEnvErrors()
			return termElem{}, false
		}
		B = parser.env.FindUnified(B).(types.Type)
		if ok = types.SetKind(&res, B); !ok {
			parser.transferEnvErrors()
			return termElem{}, false
		}
		term = termElem{typingType, res, tyi, start, end}
		//parser.declareIdFromTyping(term)
	} else {
		// create a kind (or get one if it already exists) for the term
		term = termElem{typingType, b.Term, tyi, start, end}
	}

	return term, true
}

func (parser *Parser) splitTyping(elem termElem) (x types.Variable, ok bool) {
	var A types.Type
	if elem.NodeType == labeledFunctionType {
		if x, ok = elem.Term.(types.Variable); !ok {
			// error: expected a variable term
			parser.errorOn(ExpectedVariable, elem)
			return
		}
		// x : A, A : s
		A = types.GetKind(&x)
	} else if A, ok = elem.Term.(types.Type); ok {
		// _ : A, A : s
		x = types.AsTyping(A)
	} else {
		ok = false
		parser.errorOn(UnexpectedTerm, elem)
		return
	}

	types.GetKind(&A)
	return x, true
}

func (parser *Parser) splitReturn(elem termElem) (B types.Type, ok bool) {
	// named arguments cannot be returned (this could create confusion especially in the case of a named
	// product since functions are curried)
	if elem.NodeType == typingType {
		parser.errorOn(CannotReturnNamedArg, elem)
		return nil, false
	} else if elem.NodeType == implicitFunctionType {
		parser.errorOn(CannotReturnImplicitArg, elem)
		return nil, false
	}
	B, ok = elem.Term.(types.Type)
	if !ok {
		parser.errorOn(ExpectedType, elem)
		return nil, false
	}
	types.GetKind(&B)
	return B, true
}

// assumes fi.arity == 0
func (parser *Parser) produce(function, newTerm termElem, fi termInfo) (term termElem, ok bool) {
	var intro types.PiIntro
	// determines whether or not to use the implicit-abstracting version of the `Prod` rule
	isImplicit := function.NodeType == implicitFunctionType

	x, good := parser.splitTyping(function)
	if ok = good; !ok {
		return term, false
	}

	if isImplicit {
		intro, ok = parser.env.ImplicitProd(x)
	} else {
		intro, ok = parser.env.Prod(x)
	}

	if !ok {
		parser.transferEnvErrors()
		return term, false
	}

	var Pi types.Pi
	B, secondOk := parser.splitReturn(newTerm)
	if ok = secondOk; !ok {
		return term, false
	}
	if Pi, ok = intro(B); !ok {
		parser.transferEnvErrors()
		return term, false
	}
	start, end := getTermsPos(function, newTerm)
	term = termElem{functionType, Pi, fi, start, end}
	return term, true
}

func reduceFuncType(parser *Parser, function, newTerm termElem) (term termElem, ok bool) {
	fi, decd := function.termInfo.decrementArity()
	if !decd {
		return termElem{}, false
	}

	if fi.arity == 0 {
		return parser.produce(function, newTerm, fi)
	} else {
		// function is just "->" (nothing applied to it yet)
		start, end := getTermsPos(function, newTerm)
		term = termElem{functionType, newTerm.Term, fi, start, end}
		if newTerm.NodeType == typingType || newTerm.NodeType == implicitType {
			if newTerm.NodeType == typingType {
				term.NodeType = labeledFunctionType
			} else {
				term.NodeType = implicitFunctionType
			}
			// declare the typed variable
			parser.declareIdFromTyping(newTerm)
			if parser.panicking {
				return term, false
			}
		}
	}
	return term, true
}

func reduceListing(parser *Parser, a, b termElem) (termElem, bool) {
	panic("TODO: implement") // TODO: implement
	/*
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
	*/
}

func selectReduction(t termElem) reduceFunction {
	rf := reduceTable[t.NodeType]
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
	panic("TODO: implement") // TODO: implement
	/*
		ls := term.Term.(Listing)
		tuple := Tuple(ls)
		tuple.Start, tuple.End = start, end
		term = termElem{tuple, term.termInfo}
		return term, true
	*/
}

func closeClosedTuple(term termElem, start, end int) (termElem, bool) {
	panic("TODO: implement") // TODO: implement
	/*
		// update start and end--enclosing in parens does nothing otherwise
		ps := term.Term.(Tuple)
		ps.Start = start
		ps.End = end
		term.Term = ps
		return term, true
	*/
}

func standardParenCloser(parser *Parser, term termElem, start, end int) (termElem, bool) {
	ty := term.NodeType
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
	term.Start, term.End = start, end
	term.termInfo = info
	return term, true
}
