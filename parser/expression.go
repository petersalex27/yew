// =================================================================================================
// Alex Peters - January 26, 2024
// =================================================================================================

package parser

import "github.com/petersalex27/yew/token"

// one iteration of parenthesized expression parse
func (parser *Parser) iterateParenParse() (right ExprNode, end bool) {
	var ok bool
	right, ok = parser.parseExpression()
	if end = !ok; end {
		return
	}
	parser.dropNewlines()
	end = parser.Next.Type == token.RightParen || parser.Next.Type == token.End
	return
}

// parses one element and a separator in a list-like expression
//
//   - param listLike: pointer to list-like AST node, its field `Elems` will be filled in with
//     expression nodes
//   - param separatorType: type that separates expression elements
//   - param endType: type that denotes end of list-like
//
// returns true iff parsing should end. Note that this includes returning true when there is an
// error, not only when `endType` is found. To check if an error occurred, call
// `(*Parser) Panicking()`
//
// SEE:
//
//	(*Parser) Panicking()
func (parser *Parser) iterateListLikeParse(listLike *ListLike, separatorType token.Type, endType token.Type) (end bool) {
	next, ok := parser.parseExpression()
	if end = !ok; end {
		return
	}

	listLike.Elems = append(listLike.Elems, next)
	parser.dropNewlines()
	var forceEnd bool
	if forceEnd = parser.Next.Type != separatorType; !forceEnd {
		_ = parser.Advance()
		parser.dropNewlines()
	}

	// if forceEnd and not next is not a right paren, error will be thrown later when right paren is
	// checked for but not found
	end = forceEnd || parser.Next.Type == endType || parser.Next.Type == token.End
	return
}

// handles iterative logic for parsing parenthesized expressions
func (parser *Parser) loopParenParse(e ExprNode) (_ ExprNode, ok bool) {
	end := parser.Next.Type == token.RightParen || parser.Next.Type == token.End
	if end {
		return e, true
	}

	a := Application{}
	a.Elems = []ExprNode{e}
	for !end {
		e, end = parser.iterateParenParse()
		if !end {
			a.Elems = append(a.Elems, e)
		}
	}
	a.Start, _ = a.Elems[0].Pos()
	_, a.End = a.Elems[len(a.Elems)-1].Pos()
	return a, !parser.panicking
}

func (parser *Parser) parseApplication(first ExprNode) (a Application, ok bool) {
	again := true
	a.Elems = []ExprNode{first}
	for e := ExprNode(nil); again; {
		e, again = parser.parseExpressionIteration()
		if again {
			a.Elems = append(a.Elems, e)
		}
	}
	a.Start, _ = a.Elems[0].Pos()
	_, a.End = a.Elems[len(a.Elems)-1].Pos()
	return a, !parser.panicking
}

// parses either boring kind "()" or an parenthesized expression
func (parser *Parser) parseBoringOrEnclosed(start int) (ExprNode, bool) {
	if parser.Next.Type != token.RightParen {
		return parser.parseEnclosed(start)
	}

	end := parser.Advance().End
	return BoringKind{Start: start, End: end}, true
}

func (parser *Parser) parseElements(initElems []ExprNode, separatorType, endType token.Type) (listLike ListLike, ok bool) {
	listLike.Elems = initElems
	end := parser.Next.Type == endType || parser.Next.Type == token.End
	for !end {
		end = parser.iterateListLikeParse(&listLike, separatorType, endType)
	}
	ok = parser.validateEndOfListLike(endType)
	return
}

// parses an expression enclosed by parentheses
func (parser *Parser) parseEnclosed(start int) (pe ExprNode, ok bool) {
	parser.dropNewlines()

	var e ExprNode
	e, ok = parser.parseExpression()
	if !ok {
		return
	}

	parser.dropNewlines()

	// is it `(e, ...` ?
	isTuple := parser.Next.Type == token.Comma
	if isTuple {
		_ = parser.Advance()
		return parser.parseTuple(e, start)
	}

	if e, ok = parser.loopParenParse(e); !ok {
		return
	} else if ok = parser.Next.Type == token.RightParen; !ok {
		parser.error(UnexpectedEOF)
		return
	}

	end := parser.Advance().End
	pe = ParenExpr{Start: start, End: end, ExprNode: e}

	return
}

// a single expression parse iteration
//
// parses the following:
//   - integers
//   - characters
//   - floating point values
//   - strings
//   - IDs (including affixed ones)
//   - parenthesized expressions
//   - lists
//   - tuples
//   - other non-builtin data types (recognized by starting w/ upper case name)
func (parser *Parser) parseExpressionIteration() (expression ExprNode, ok bool) {
	tokenType := parser.Next.Type
	switch tokenType {
	case token.IntValue:
		fallthrough
	case token.CharValue:
		fallthrough
	case token.FloatValue:
		fallthrough
	case token.StringValue:
		expression, ok = Constant{parser.Advance()}, true
	case token.Id:
		fallthrough
	case token.Affixed:
		expression, ok = Ident{parser.Advance()}, true
	case token.LeftBracket:
		start := parser.Advance().Start
		expression, ok = parser.parseList(start)
	case token.CapId:
		expression, ok = KindIdent{parser.Advance()}, true
	case token.LeftParen:
		start := parser.Advance().Start
		expression, ok = parser.parseBoringOrEnclosed(start)
	}
	return
}

func (parser *Parser) parseExpression() (expression ExprNode, ok bool) {
	expression, ok = parser.parseExpressionIteration()

	// `again` is just for readability
	if again := ok; again {
		return parser.parseApplication(expression)
	}

	// set true value of ok (it's possible parseExpressionIteration just read nothing)
	ok = !parser.panicking
	return
}

func (parser *Parser) parseList(start int) (ls List, ok bool) {
	ls, ok = parser.parseElements([]ExprNode{}, token.Comma, token.RightParen)
	if ok {
		ls.Start = start
		ls.End = parser.Advance().End
	}
	return
}

func (parser *Parser) parseTuple(first ExprNode, start int) (tuple TupleKind, ok bool) {
	tuple, ok = parser.parseElements([]ExprNode{first}, token.Comma, token.RightParen)
	if ok {
		tuple.Start = start
		tuple.End = parser.Advance().End
	}
	return
}

func (parser *Parser) validateEndOfListLike(expect token.Type) (ok bool) {
	if ok = parser.Next.Type == expect; !ok {
		if parser.Next.Type == token.End {
			parser.error(UnexpectedEOF)
			return
		}
		errorMessage := getExpectMessage(expect)
		parser.error(errorMessage)
	}
	return
}
