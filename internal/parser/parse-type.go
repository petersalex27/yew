package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/api/util/fun"
	"github.com/petersalex27/yew/common/data"
	nodeType "github.com/petersalex27/yew/internal/parser/typ"
)

// rule:
//
//	```
//	type =
//		["forall", {"\n"}, forall binders, {"\n"}, "in", {"\n"}], type tail
//		| "(", {"\n"}, enc type, {"\n"}, ")" ;
//	```
func ParseType(p Parser) data.Either[data.Ers, typ] {
	return parseType(p, false)
}

// decides b/w parsing a forall type or a type tail
func parseTypeHelper(p Parser, enclosed bool) data.Either[data.Ers, typ] {
	// only allowed at the beginning of a type
	if matchCurrentForall(p) {
		binding := parseForallBinding(p)
		if es, binders, ok := binding.Break(); !ok {
			return data.PassErs[typ](es)
		} else {
			p.dropNewlines()
			return forallBindParsedType(p, binders)
		}
	}
	return parseTypeTail(p, enclosed)
}

// IMPORTANT: enclosed types (ones found to be such b/c a left-paren is found--and not b/c of the
// argument) are NOT returned as `enclosedType` unless the type parsed in the rule `enc type` is
// an enclosedType. Intuitively, you can think of this as just returning whatever the rule `enc type`
// would return, otherwise whatever the `"forall" ...` production would return.
//
// rule:
//
//	```
//	type =
//		["forall", {"\n"}, forall binders, {"\n"}, "in", {"\n"}], type tail
//		| "(", {"\n"}, enc type, {"\n"}, ")" ;
//	```
//
// typeHelper handles deciding b/w a forall type and a type tail
func parseType(p Parser, enclosed bool) data.Either[data.Ers, typ] {
	lparen, found := getKeywordAtCurrent(p, token.LeftParen)
	if !found {
		// not enclosed in this call, parse with inherited `enclosed` value
		//
		// this call is the only one that moves the parsing of the type forward
		// in this function
		return parseTypeHelper(p, enclosed)
	}

	// parse enclosed type, recursively calling this function
	es, ty, isTy := parseType(p, true).Break()
	if !isTy {
		return data.PassErs[typ](es)
	}

	ty = ty.updatePosTyp(lparen) // update position to include left paren

	p.dropNewlines()
	rparen, found := getKeywordAtCurrent(p, token.RightParen)
	if !found {
		return data.Fail[typ](ExpectedRightParen, p)
	}

	ty = ty.updatePosTyp(rparen) // update position to include right paren

	return data.Ok(ty)
}

// an error will not be thrown if no `type term` is found on the lhs; otherwise, an error will be
// thrown as normal
//
// rule:
//
//	```
//	type tail = type term, {type term rhs}, [{"\n"}, ("->" | "=>"), {"\n"}, type tail] ;
//	```
func parseMaybeTypeTail(p Parser, enclosed bool) (*data.Ers, data.Maybe[typ]) {
	es, mTerm := maybeParseTypeTerm(p)
	if es != nil {
		return es, data.Nothing[typ](p)
	}

	term, isSomething := mTerm.Break()
	if !isSomething {
		return nil, mTerm
	}

	es, head := parseJustAppTypeOrJustType(p, term, enclosed)
	if es != nil {
		return es, data.Nothing[typ](p)
	}

	p.dropNewlines()
	pos := api.ZeroPosition()
	functionTyped := parseKeywordAtCurrent(p, token.Arrow, &pos)                  // is this a function type?
	someFunc := functionTyped || parseKeywordAtCurrent(p, token.ThickArrow, &pos) // if not, is it a constraint?
	if !someFunc {
		// not a function type, return term
		return nil, data.Just(head)
	}

	// otherwise, create function
	functionRes := runCases(p,
		fun.Bind1stOf2(parseTypeTail, enclosed),
		passParseErs[typ],
		constructFunction(head, functionTyped),
	)
	esFull, function, isFunction := functionRes.Break()
	if !isFunction {
		return &esFull, data.Nothing[typ](p)
	}
	// update position to include whichever arrow was parsed ("->" or "=>")
	function = function.updatePosTyp(pos)
	return nil, data.Just(function)
}

// parse tail-end of `type` rule
//
// rule:
//
//	```
//	type tail = type term, {type term}, [{"\n"}, ("->" | "=>"), {"\n"}, type tail] ;
//	```
func parseTypeTail(p Parser, enclosed bool) data.Either[data.Ers, typ] {
	es, mTail := parseMaybeTypeTail(p, enclosed)
	if es != nil {
		return data.PassErs[typ](*es)
	}
	unit, just := mTail.Break()
	if just {
		return data.Ok(unit)
	}
	return data.Fail[typ](ExpectedType, p)
}

// constructs Either a function or constrained type (constrained by an unverified constraint)
//
// the constraint will be verified during type checking
func constructFunction(lhs typ, isFunctionType bool) func(Parser, typ) data.Either[data.Ers, typ] {
	if lhs == nil {
		panic("lhs cannot be nil")
	}
	return func(p Parser, rhs typ) data.Either[data.Ers, typ] {
		if isFunctionType {
			return data.Ok(makeFunc(lhs, rhs))
		} // else, constraint
		return data.Ok(makeUnverifiedConstrainedType(lhs, rhs))
	}
}

func makeUnverifiedConstrainedType(lhs, rhs typ) typ {
	var constraint constraint = data.EOne[constraintUnverified](lhs)
	return data.EMakePair[constrainedType](constraint, rhs)
}

func parseForallBinders(p Parser, forallKey api.Token) data.Either[data.Ers, forallBinders] {
	ids := data.Nil[ident]()
	id, just := parseIdent(p).Break()
	for ; just; id, just = parseIdent(p).Break() {
		ids = ids.Snoc(id)
	}

	if forall, ok := data.EStrengthen[forallBinders](ids).Break(); ok {
		forall.Position = forall.Update(forallKey)
		return data.Ok(forall)
	} else {
		return data.Fail[forallBinders](ExpectedId, p) // use `p` b/c `ids` must be empty
	}
}

func forallBindParsedType(p Parser, fb forallBinders) data.Either[data.Ers, typ] {
	p.dropNewlines()
	pos := api.ZeroPosition()
	if parseKeywordAtCurrent(p, token.In, &pos) {
		res := data.Cases(parseTypeTail(p, false), data.PassErs[typ], assembleForallType(fb))
		res = res.Update(pos)
		return res
	}
	return data.Fail[typ](ExpectedForallIn, pos)
}

func assembleForallType(fb forallBinders) func(typ) data.Either[data.Ers, typ] {
	return func(t typ) data.Either[data.Ers, typ] {
		return data.Ok(typ(forallType{data.MakePair(fb, t)}))
	}
}

// parses term in a type
//
// rule:
//
//	```
//	type term =
//		expr root
//		| "_" | "()" | "="
//		| "(", {"\n"}, enc type inner, [{"\n"}, enc typing end], {"\n"}, ")"
//		| "{", {"\n"}, enc type inner, [{"\n"}, enc typing end, [{"\n"}, default expr]], {"\n"}, "}" ;
//	```
func maybeParseTypeTerm(p Parser) (*data.Ers, data.Maybe[typ]) {
	return maybeParseTypeTermHelper(p, false)
}

func maybeParseTypeTermRhs(p Parser) (*data.Ers, data.Maybe[typ]) {
	return maybeParseTypeTermHelper(p, true)
}

func maybeParseTypeTermHelper(p Parser, rhs bool) (*data.Ers, data.Maybe[typ]) {
	var lhsTerm data.Either[data.Ers, typ]
	if rhs && lookahead1(p, token.Dot) {
		es, acc := parseAccess(p)
		if es != nil {
			return es, data.Nothing[typ](p)
		}
		lhsTerm = data.Ok[typ](acc)
	} else if lookahead1(p, exprAtomLAs...) {
		lhsTerm = data.Cases(parseExprAtom(p), data.PassErs[typ], exprAtomAsExprRes)
	} else if lookahead1(p, typeTermExceptionLAs...) {
		lhsTerm = parseTypeTermException(p) // "_", "()", or "="
	} else if lookahead1(p, token.LeftParen, token.LeftBrace) {
		lhsTerm = parseEnclosedType(p)
	} else {
		return nil, data.Nothing[typ](p)
	}

	es, lhs, isRight := lhsTerm.Break()
	if !isRight {
		return &es, data.Nothing[typ](p)
	}
	return nil, data.Just(lhs)
}

// ASSUMPTION: the current token is Either "_", "()", "="
func parseTypeTermException(p Parser) data.Either[data.Ers, typ] {
	tok, found := getKeywordAtCurrent(p, token.Underscore)
	if found {
		return data.Ok(typ(data.EOne[wildcard](tok)))
	}

	tok, found = getKeywordAtCurrent(p, token.Equal)
	if found {
		return data.Ok[typ](data.EOne[name](tok))
	}

	tok, found = getKeywordAtCurrent(p, token.EmptyParenEnclosure)
	if !found {
		panic("expected '_', '()', '='") // input was not validated before calling
	}
	return data.Ok[typ](data.EOne[unitType](tok))
}

func parseJustAppTypeOrJustType(p Parser, lhs typ, enclosed bool) (*data.Ers, typ) {
	if lhs == nil {
		panic("lhs cannot be nil")
	}
	es, types, has2ndTerm := parseOneOrMore(p, lhs, enclosed, maybeParseTypeTermRhs)
	if es != nil {
		return es, nil
	}

	if !has2ndTerm {
		return nil, types.Head()
	}

	// otherwise, construct app type
	app, _ := data.NonEmptyToAppLikePair[appType](types).Break()
	return nil, app
}

// rule:
//
//	```
//	forall binding = "forall", {"\n"}, forall binders, {"\n"}, "in", {"\n"}
//	```
func parseForallBinding(p Parser) data.Either[data.Ers, forallBinders] {
	forallTok, found := getKeywordAtCurrent(p, token.Forall)
	if !found {
		panic("expected 'forall'") // input was not validated before calling
	}
	return parseForallBinders(p, forallTok)
}

func parseOptionalModality(p Parser) data.Maybe[modality] {
	mModality := data.Nothing[modality](p)
	// parse optional multiplicity modality
	mode, found := getKeywordAtCurrent(p, token.Erase)
	if !found {
		mode, found = getKeywordAtCurrent(p, token.Once)
	}

	// set modality to non-'Nothing' value if found
	if found {
		mModality = data.Just(data.EOne[modality](mode))
	}
	return mModality
}

// rule:
//
//	```
//	enclosed type =
//		"(", {"\n"}, enc type inner, [{"\n"}, enc typing end], {"\n"}, ")"
//		| "{", {"\n"}, enc type inner, [{"\n"}, enc typing end, [{"\n"}, default expr]], {"\n"}, "}" ;
//	enc type inner = multiplicity, {"\n"}, ident | inner type terms ;
//	enc typing end = ":", {"\n"}, enc type ;
//	inner type terms = enc type tail, [{{"\n"}, ",", {"\n"}, enc type tail}, [{"\n"}, ","]] ;
//	```
func parseEnclosedType(p Parser) (out data.Either[data.Ers, typ]) {
	opener, closerType, found := parseEnclosedOpener(p)
	if !found {
		panic("expected left paren or left brace") // input was not validated before calling
	}

	// enclosed type represents the kind of an implicit parameter
	implicit := token.LeftBrace.Match(opener)

	mModality := parseOptionalModality(p)

	// parse enclosed type inner head
	es, lhs, ok := parseInnerTypeTerms(p).Break()
	if !ok {
		return data.PassErs[typ](es)
	} else if !mModality.IsNothing() {
		// validate that the it's a single term and that single term is an identifier
		term := lhs.Head()
		valid := lhs.Len() == 1 && (nodeType.LowerIdent.Match(term) || nodeType.UpperIdent.Match(term))
		if !valid {
			return data.Fail[typ](ExpectedModalId, lhs)
		}
	}

	colon, foundColon := getKeywordAtCurrent(p, token.Colon)
	if !foundColon && !mModality.IsNothing() {
		return data.Fail[typ](ExpectedTypeSig, p)
		// multiplicity requires term to be typed
	} else if foundColon {
		return parseEnclosedRhs(p, lhs, opener, closerType, implicit, mModality, colon)
	} else if lhs.Tail().Len() == 0 {
		return closeEnclosedTyp(p, opener, closerType)(lhs.Head())
	} else {
		return closeEnclosedTyp(p, opener, closerType)(lhs)
	}
}

func parseEnclosedRhs(p Parser, lhs innerTypeTerms, opener api.Token, closerType token.Type, implicit bool, mModality data.Maybe[modality], colon api.Token) data.Either[data.Ers, typ] {
	// parse typing and assemble inner type signature
	ty := parseType(p, true)
	typing := data.Cases(ty, data.PassErs[innerTyping], assembleInnerTyping(mModality, lhs, colon))
	attached := optionalAttachDefaultExpr(p, implicit, typing)
	construct := closeEnclosedTyp(p, opener, closerType)
	return data.Cases(attached, data.PassErs[typ], construct)
}

// panics if colonEqual is nil
func getColonEqualAtCurrent(p Parser, colonEqual *api.Token) (found bool) {
	if colonEqual == nil {
		panic("nil pointer for ':=' token")
	}
	*colonEqual, found = getKeywordAtCurrent(p, token.ColonEqual)
	return found
}

func optionalAttachDefaultExpr(p Parser, implicit bool, typing data.Either[data.Ers, innerTyping]) data.Either[data.Ers, typ] {
	var colonEqual api.Token
	noDefault := !(data.IsRight(typing) && implicit && getColonEqualAtCurrent(p, &colonEqual))
	if noDefault {
		return data.Cases(typing, data.PassErs[typ], fun.Compose(data.Ok, (innerTyping).asTyp))
	}

	es, de, isDE := ParseExpr(p).Break()
	if !isDE {
		return data.PassErs[typ](es)
	}
	return data.Cases(typing, data.PassErs[typ], appendDefaultExpr(de))
}

func appendDefaultExpr(de expr) func(innerTyping) data.Either[data.Ers, typ] {
	return func(t innerTyping) data.Either[data.Ers, typ] {
		return data.Ok[typ](data.EMakePair[implicitTyping](t, data.EOne[defaultExpr](de)))
	}
}

func closeEnclosedTyp(p Parser, opener api.Token, closerType token.Type) func(typ) data.Either[data.Ers, typ] {
	return func(t typ) data.Either[data.Ers, typ] {
		t = t.updatePosTyp(opener)
		p.dropNewlines()
		closer, found := getKeywordAtCurrent(p, closerType)
		if !found {
			isRp := closerType.Match(token.RightParen.Make())
			return data.Fail[typ](ifThenElse(isRp, ExpectedRightParen, ExpectedRightBrace), p)
		}
		et := enclosedType{implicit: token.LeftBrace.Match(opener), typ: t}
		et.typ = et.updatePosTyp(closer)
		return data.Ok[typ](et)
	}
}

func assembleInnerTyping(modality data.Maybe[modality], lhs innerTypeTerms, colon api.Token) func(typ) data.Either[data.Ers, innerTyping] {
	return func(t typ) data.Either[data.Ers, innerTyping] {
		typing := data.MakePair(lhs, t)
		it := innerTyping{mode: modality, typing: typing}
		it.Position = api.WeakenRangeOver[api.Node](colon, modality, typing)
		return data.Ok(it)
	}
}

// rule:
//
//	```
//	inner type terms = type tail, [{{"\n"}, ",", {"\n"}, type tail}, [{"\n"}, ","]] ;
//	```
func parseInnerTypeTerms(p Parser) data.Either[data.Ers, innerTypeTerms] {
	return parseSepSequenced[innerTypeTerms](p, ExpectedType, token.Comma, fun.BinBind1st_PairTarget(parseMaybeTypeTail, true))
}
