package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/api/util/fun"
	"github.com/petersalex27/yew/common/data"
)

// rule:
//
//	```
//	expr = expr term, {expr term} ;
//	```
func ParseExpr(p Parser) data.Either[data.Ers, expr] {
	return parseExpr(p, false)
}

// rule:
//
//	```
//	expr = expr term, {expr term} ;
//	```
//
// This function should only be called from within this file "expr-parse.go"--nowhere else!
func parseExpr(p Parser, enclosed bool) data.Either[data.Ers, expr] {
	es, mExpr := parseMaybeExpr(p, enclosed)
	if es != nil {
		return data.PassErs[expr](*es)
	} else if unit, just := mExpr.Break(); just {
		return data.Ok(unit)
	} else {
		return data.Fail[expr](ExpectedExpr, p)
	}
}

// rule:
//
//	```
//	maybe expr = [expr] ;
//	```
//
// This function should only be called from within this file "expr-parse.go"--nowhere else!
func parseMaybeExpr(p Parser, enclosed bool) (*data.Ers, data.Maybe[expr]) {
	es, mFirst := parseMaybeExprTerm(p)
	if es != nil {
		return es, data.Nothing[expr](p)
	} else if mFirst.IsNothing() {
		return nil, data.Nothing[expr](p)
	}

	unit, _ := mFirst.Break()

	pEs, exp := parseJustAppExprOrJustExpr(p, unit, enclosed)
	if pEs != nil {
		return pEs, data.Nothing[expr](p)
	}

	return nil, data.Just(exp)
}

func parseJustAppExprOrJustExpr(p Parser, lhs expr, enclosed bool) (*data.Ers, expr) {
	es, exps, has2ndTerm := parseOneOrMore(p, lhs, enclosedDependentIt(enclosed), parseMaybeExprTermRhs)
	if es != nil {
		return es, nil
	}

	if !has2ndTerm {
		return nil, exps.Head() // only one term
	}
	// has2ndTerm guarantees that this will be not 'data.Nothing'
	app, _ := data.NonEmptyToAppLikePair[exprApp](exps).Break()
	return nil, app
}

// rule:
//
//	```
//	expr term = expr atom | "(", {"\n"}, enc expr, {"\n"}, ")" | let expr | case expr ;
//	```
func parseMaybeExprTerm(p Parser) (*data.Ers, data.Maybe[expr]) {
	return parseMaybeExprTermHelper(p, false)
}

func parseMaybeExprTermRhs(p Parser) (*data.Ers, data.Maybe[expr]) {
	return parseMaybeExprTermHelper(p, true)
}

func parseMaybeExprTermHelper(p Parser, rhs bool) (*data.Ers, data.Maybe[expr]) {
	if rhs && lookahead1(p, token.Dot) {
		es, acc := parseAccess(p)
		if es != nil {
			return es, data.Nothing[expr](p)
		}
		return nil, data.Just[expr](acc)
	} else if lookahead1(p, token.Let) { // "let"
		es, e := parseMaybeLetExpr(p)
		return es, data.MaybeMap((letExpr).asExpr)(e)
	} else if lookahead1(p, token.Case) { // "case"
		es, e := parseMaybeCaseExpr(p)
		return es, data.MaybeMap((caseExpr).asExpr)(e)
	} else if lookahead1(p, token.LeftParen) { // "("
		return parseMaybeEnclosedExpr(p)
	}
	es, e := parseMaybeExprAtom(p) // expr atom
	return es, data.MaybeMap(exprAtomAsExpr)(e)
}

// rule:
//
//	```
//	let expr = "let", {"\n"}, (binding group | binding assignment), {"\n"}, "in", {"\n"}, expr ;
//	```
func parseMaybeLetExpr(p Parser) (*data.Ers, data.Maybe[letExpr]) {
	let, found := getKeywordAtCurrent(p, token.Let, dropAfter)
	if !found {
		return nil, data.Nothing[letExpr](p)
	}

	es, binders, isRight := parseLetBinding(p).Break()
	if !isRight {
		return &es, data.Nothing[letExpr](p) // error, binders are required
	}

	in, found := getKeywordAtCurrent(p, token.In, dropBeforeAndAfter)
	if !found {
		// error, 'in' keyword is required
		es := data.Ers(data.Nil[data.Err](1).Snoc(data.MkErr(ExpectedIn, p)))
		return &es, data.Nothing[letExpr](p)
	}

	es2, expression, isExpression := ParseExpr(p).Break()
	if !isExpression {
		return &es2, data.Nothing[letExpr](p) // error, expr is required
	}
	return nil, data.Just(assembleLetExpr(let, binders, in, expression))
}

// rule:
//
//	```
//	let binding =
//		binding group member
//		| "(", {"\n"}, binding group member, {then, binding group member}, {"\n"}, ")" ;
//	```
func parseLetBinding(p Parser) data.Either[data.Ers, letBinding] {
	return parseGroup[letBinding](p, ExpectedBindingTerm, parseMaybeBindingGroupMember)
}

// Not a real rule, just useful helper function
//
// rule:
//
//	```
//	colon equal assignment = [":=", {"\n"}, expr] ;
//	```
func parseMaybeColonEqualAssignment(p Parser) (*data.Ers, data.Maybe[expr]) {
	colonEqual, found := getKeywordAtCurrent(p, token.ColonEqual, dropAfter)
	if !found { // this is okay, assignment is optional--at least syntactically
		return nil, data.Nothing[expr](p)
	}

	es, expression, isExpression := ParseExpr(p).Break()
	if !isExpression {
		return &es, data.Nothing[expr](p) // error, expression is required b/c of ':='
	}

	// update here (instead of the data.Maybe-version) in case the non-empty value is used directly
	//	- the data.Maybe version will inherit the position
	expression = expression.updatePosExpr(colonEqual)
	mExpression := data.Just(expression)
	return nil, mExpression
}

type binderMember = data.Pair[binder, expr]
type typingMember = data.Pair[typing, data.Maybe[expr]]

// rule:
//
//	```
//	binding group member = binder, {"\n"}, ":=", {"\n"}, expr | typing, [{"\n"}, ":=", {"\n"}, expr] ;
//	```
func parseMaybeBindingGroupMember(p Parser) (*data.Ers, data.Maybe[bindingGroupMember]) {
	var lhs data.Either[binder, typing]
	isTyping := lookahead2(p, bindingTypingLAs...)
	if isTyping {
		es, typing, isTyping := parseTypeSig(p).Break()
		if !isTyping {
			return &es, data.Nothing[bindingGroupMember](p)
		}
		lhs = data.Inr[binder](typing)
	} else {
		es, mBinder := parseMaybeBinder(p)
		if es != nil {
			return es, data.Nothing[bindingGroupMember](p) // error, maybe parse had an error
		} else if binder, just := mBinder.Break(); !just {
			return nil, data.Nothing[bindingGroupMember](p) // no binder, return 'data.Nothing'
		} else {
			lhs = data.Inl[typing](binder) // found a binder, continue parsing member
		}
	}

	p.dropNewlines()
	es, mExpression := parseMaybeColonEqualAssignment(p)
	if es != nil {
		return es, data.Nothing[bindingGroupMember](p)
	} else if isTyping {
		// expression might be 'data.Nothing' if the typing has no associated definition. This is okay,
		// assignment is optional--at least syntactically--for a let-typing
		//
		// NOTE: during type checking, if the typing has no associated definition, it will be
		// considered an error. But, for syntax analysis, this is permissible.
		_, r, _ := lhs.Break()
		typingMem := data.Inr[binderMember](data.MakePair(r, mExpression))
		return nil, data.Just(typingMem)
	} else if mExpression.IsNothing() { // error, expression is required for a let-binding
		e := data.Ers(data.Nil[data.Err](1).Snoc(data.MkErr(ExpectedBoundExpr, mExpression)))
		return &e, data.Nothing[bindingGroupMember](p)
	}

	l, _, _ := lhs.Break()
	unit, _ := mExpression.Break()
	binderMem := data.Inl[typingMember](data.MakePair(l, unit))
	return nil, data.Just(binderMem)
}

// rule:
//
//	```
//	case expr = "case", {"\n"}, pattern, {"\n"}, "of", {"\n"}, case arms ;
//	```
func parseMaybeCaseExpr(p Parser) (*data.Ers, data.Maybe[caseExpr]) {
	caseToken, found := getKeywordAtCurrent(p, token.Case, dropAfter)
	if !found {
		return nil, data.Nothing[caseExpr](p)
	}

	es, pat, isRight := ParsePattern(p).Break()
	if !isRight {
		return &es, data.Nothing[caseExpr](p)
	}

	of, found := getKeywordAtCurrent(p, token.Of, dropBeforeAndAfter)
	if !found {
		es := data.Ers(data.Nil[data.Err](1).Snoc(data.MkErr(ExpectedOf, p)))
		return &es, data.Nothing[caseExpr](p)
	}

	esArms, arms, isArmsRight := parseCaseArms(p).Break()
	if !isArmsRight {
		return &esArms, data.Nothing[caseExpr](p)
	}

	ce := data.EMakePair[caseExpr](pat, arms)
	ce.Position = ce.Update(caseToken)
	ce.Position = ce.Update(of)
	return nil, data.Just(ce)
}

// rule:
//
//	```
//	case arms = case arm | "(", {"\n"}, case arm, {then, case arm}, {"\n"}, ")" ;
//	```
func parseCaseArms(p Parser) data.Either[data.Ers, caseArms] {
	return parseGroup[caseArms](p, ExpectedCaseArm, maybeParseCaseArm)
}

// rule:
//
//	```
//	case arm = pattern, {"\n"}, def body thick arrow ;
//	```
func maybeParseCaseArm(p Parser) (*data.Ers, data.Maybe[caseArm]) {
	es, mPat := maybeParsePattern(p, false)
	if es != nil {
		return es, data.Nothing[caseArm](p)
	} else if mPat.IsNothing() {
		return nil, data.Nothing[caseArm](p)
	}
	pat, _ := mPat.Break()

	// the rest is required

	// drop newlines before '=>' (or 'impossible')
	p.dropNewlines()
	esBody, body, isBodyRight := parsePatternBoundBody(p, token.ThickArrow).Break()
	if !isBodyRight {
		return &esBody, data.Nothing[caseArm](p)
	}
	return nil, data.Just(data.EMakePair[caseArm](pat, body))
}

// Not a real rule, a helper function for parsing enclosed expressions
//
// rule:
//
//	```
//	enc expr' = ["(", {"\n"}, enc expr, {"\n"}, ")"] ;
//	```
func parseMaybeEnclosedExpr(p Parser) (*data.Ers, data.Maybe[expr]) {
	lparen, found := getKeywordAtCurrent(p, token.LeftParen, dropAfter)
	if !found {
		return nil, data.Nothing[expr](p) // fine, not enclosed, return 'data.Nothing'
	}

	es, e, isRight := parseExpr(p, true).Break()
	if !isRight {
		return &es, data.Nothing[expr](p) // error, expr is required
	}

	var rparen api.Token
	if rparen, found = getKeywordAtCurrent(p, token.RightParen, dropBefore); !found {
		es := data.Ers(data.Nil[data.Err](1).Snoc(data.MkErr(ExpectedRightParen, p)))
		return &es, data.Nothing[expr](p) // error, right paren is required
	}

	e = e.updatePosExpr(lparen).updatePosExpr(rparen)

	return nil, data.Just(e)
}

// rule:
//
//	```
//	expr atom = pattern atom | lambda abstraction ;
//	```
func parseMaybeExprAtom(p Parser) (*data.Ers, data.Maybe[exprAtom]) {
	if lookahead1(p, exprAtomLAs...) {
		es, e, isRight := parseExprAtom(p).Break()
		if !isRight {
			return &es, data.Nothing[exprAtom](p)
		}
		return nil, data.Just(e)
	}
	return nil, data.Nothing[exprAtom](p)
}

// rule:
//
//	```
//	expr atom = pattern atom | lambda abstraction ;
//	```
func parseExprAtom(p Parser) data.Either[data.Ers, exprAtom] {
	if matchCurrentBackslash(p) {
		return data.Cases(parseLambdaAbstraction(p), data.Inl[exprAtom, data.Ers], fun.Compose(data.Ok, data.Inr[patternAtom, lambdaAbstraction]))
	}
	return data.Cases(parsePatternAtom(p), data.Inl[exprAtom, data.Ers], fun.Compose(data.Ok, data.Inl[lambdaAbstraction, patternAtom]))
}

// rule:
//
//	```
//	lambda abstraction = "\\", {"\n"}, lambda binders, {"\n"}, "=>", {"\n"}, expr ;
//		lambda binders = lambda binder, {{"\n"}, ",", {"\n"}, lambda binder}, [{"\n"}, ","] ;
//		lambda binder = binder | "_" ;
//	```
func parseLambdaAbstraction(p Parser) data.Either[data.Ers, lambdaAbstraction] {
	backslash, found := getKeywordAtCurrent(p, token.Backslash, dropAfter)
	if !found {
		return data.Fail[lambdaAbstraction](ExpectedLambdaAbstraction, p)
	}

	es, binders, isRight := parseSepSequenced[lambdaBinders](p,
		ExpectedLambdaAbstraction,
		token.Comma,
		parseMaybeLambdaBinder,
	).Break()
	if !isRight {
		return data.PassErs[lambdaAbstraction](es)
	}

	arrow, found := getKeywordAtCurrent(p, token.ThickArrow, dropBeforeAndAfter)
	if !found {
		return data.Fail[lambdaAbstraction](ExpectedLambdaThickArrow, p)
	}
	return data.Cases(ParseExpr(p),
		data.Inl[lambdaAbstraction, data.Ers],
		constructLambdaAbstraction(backslash, binders, arrow),
	)
}

// rule:
//
//	```
//	lambda binder =  binder | "_" ;
//	```
func parseMaybeLambdaBinder(p Parser) (*data.Ers, data.Maybe[lambdaBinder]) {
	if matchCurrentUnderscore(p) {
		underscore := p.current()
		p.advance()
		return nil, data.Just(data.EInr[lambdaBinder](data.EOne[wildcard](underscore)))
	}
	es, mBinder := parseMaybeBinder(p)
	if es != nil {
		return es, data.Nothing[lambdaBinder](p)
	}
	// lift 'binder' into a 'lambdaBinder', then lift that into a 'data.Maybe'--or, return 'data.Nothing' if
	// 'binder' is 'data.Nothing'
	return nil, bind(mBinder, fun.Compose(data.Just, data.EInl[lambdaBinder]))
}

// rule:
//
//	```
//	binder = ident | "(", {"\n"}, enc pattern, {"\n"}, ")" ;
//	```
//
// NOTE: while this function cannot parse the invalid `{ecn pattern}`, it can parse the invalid
// `({enc pattern})`; however, this will be caught during name resolution--so, it is okay to parse
// this.
func parseMaybeBinder(p Parser) (*data.Ers, data.Maybe[binder]) {
	if matchCurrentLeftParen(p) {
		es, mPat := maybeParsePattern(p, false)
		if es != nil {
			return es, data.Nothing[binder](p)
		}
		// lift 'pattern' into a 'binder', then lift that into a 'data.Maybe'
		f := fun.Compose(data.Just, data.Inr[ident, pattern])
		return nil, bind(mPat, f)
	}

	id, isSomething := parseIdent(p).Break()
	if !isSomething {
		return nil, data.Nothing[binder](p)
	}
	return nil, data.Just(data.Inl[pattern](id))
}
