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
func ParseType(p parser) data.Either[data.Ers, Typ] {
	return parseType(p, false)
}

// decides b/w parsing a forall type or a type tail
func parseTypeHelper(p parser, enclosed bool) data.Either[data.Ers, Typ] {
	// only allowed at the beginning of a type
	if matchCurrentForall(p) {
		binding := parseForallBinding(p)
		if es, binders, ok := binding.Break(); !ok {
			return data.PassErs[Typ](es)
		} else {
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
func parseType(p parser, enclosed bool) data.Either[data.Ers, Typ] {
	lparen, found := getKeywordAtCurrent(p, token.LeftParen, dropAfter)
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
		return data.PassErs[Typ](es)
	}

	ty = ty.updatePosTyp(lparen) // update position to include left paren

	rparen, found := getKeywordAtCurrent(p, token.RightParen, dropBefore)
	if !found {
		return data.Fail[Typ](ExpectedRightParen, p)
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
//	type tail = type term, {type term rhs}, [{"\n"}, function rhs] ;
//	```
func parseMaybeTypeTail(p parser, enclosed bool) (*data.Ers, data.Maybe[Typ]) {
	origin := getOrigin(p)

	es, mTerm := maybeParseTypeTerm(p)
	if es != nil {
		return es, data.Nothing[Typ](p)
	}

	term, isSomething := mTerm.Break()
	if !isSomething {
		return nil, mTerm
	}

	es, head := parseJustAppTypeOrJustType(p, term, enclosed)
	if es != nil {
		return es, data.Nothing[Typ](p)
	}

	// parse function rhs
	return parseMaybeFunctionRhs(p, head, origin, enclosed)
}

// rule:
//
//	```
//	arrow = {"\n"}, ("->" | "=>"), {"\n"} ;
//	```
func getArrow(p parser) (arrow api.Token, found bool) {
	arrow, found = getKeywordAtCurrent(p, token.Arrow, dropBeforeAndAfter)
	if !found {
		arrow, found = getKeywordAtCurrent(p, token.ThickArrow, dropBeforeAndAfter)
	}
	return arrow, found
}

// rule:
//
//	```
//	function rhs = ("->" | "=>"), {"\n"}, type tail ;
//	```
func parseMaybeFunctionRhs(p parser, lhs Typ, origin int, enclosed bool) (*data.Ers, data.Maybe[Typ]) {
	arrowOrigin := getOrigin(p) // location before parsing arrow

	arrow, found := getArrow(p)
	if !found {
		return nil, data.Just(lhs) // not a function type, return term
	}

	// builds function type from rhs
	var fromRhs func(Typ) Typ

	if token.ThickArrow.Match(arrow) { // this means lhs is a constraint
		arrowEndOrigin := getOrigin(p)

		// unfortunately, this is a bit inefficient, but it keeps things simple--we need to re-parse the input
		// as a constraint
		p = setOrigin(p, origin)
		result := parseConstraint(p)
		es, constraint, isConstraint := result.Break()
		if !isConstraint {
			return &es, data.Nothing[Typ](p)
		}

		// no need to parse arrow again, check that position ends at the same place as the arrow start
		if arrowOrigin != getOrigin(p) {
			// failed to parse constraint
			e := data.MkErr(UnexpectedToken, p)
			es := data.Makes(e)
			return &es, data.Nothing[Typ](p)
		}

		setOrigin(p, arrowEndOrigin) // set to after the arrow (and possible newlines)

		fromRhs = makeConstrainedType(constraint)
	} else {
		fromRhs = fun.Curry(makeFunc)(lhs)
	}

	construct := fun.Compose(data.Ok, fromRhs)

	res := parseTypeTail(p, enclosed)

	es, function, isFunction := data.Cases(res, data.PassErs[Typ], construct).Break()
	if !isFunction {
		return &es, data.Nothing[Typ](p)
	}

	function = function.updatePosTyp(arrow)

	return nil, data.Just(function)
}

// parse tail-end of `type` rule
//
// rule:
//
//	```
//	type tail = type term, {type term}, [{"\n"}, ("->" | "=>"), {"\n"}, type tail] ;
//	```
func parseTypeTail(p parser, enclosed bool) data.Either[data.Ers, Typ] {
	es, mTail := parseMaybeTypeTail(p, enclosed)
	if es != nil {
		return data.PassErs[Typ](*es)
	}
	unit, just := mTail.Break()
	if just {
		return data.Ok(unit)
	}
	return data.Fail[Typ](ExpectedType, p)
}

func makeConstrainedType(constraint constraint) func(Typ) Typ {
	return func(rhs Typ) Typ {
		return data.EMakePair[constrainedType](constraint, rhs)
	}
}

func parseForallBinders(p parser, forallKey api.Token) data.Either[data.Ers, forallBinders] {
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

func forallBindParsedType(p parser, fb forallBinders) data.Either[data.Ers, Typ] {
	if in, found := getKeywordAtCurrent(p, token.In, dropBeforeAndAfter); found {
		res := data.Cases(parseTypeTail(p, false), data.PassErs[Typ], assembleForallType(fb))
		res = res.Update(in)
		return res
	}
	return data.Fail[Typ](ExpectedForallIn, p)
}

func assembleForallType(fb forallBinders) func(Typ) data.Either[data.Ers, Typ] {
	return func(t Typ) data.Either[data.Ers, Typ] {
		return data.Ok(Typ(forallType{data.MakePair(fb, t)}))
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
func maybeParseTypeTerm(p parser) (*data.Ers, data.Maybe[Typ]) {
	return maybeParseTypeTermHelper(p, false)
}

func maybeParseTypeTermRhs(p parser) (*data.Ers, data.Maybe[Typ]) {
	return maybeParseTypeTermHelper(p, true)
}

func maybeParseTypeTermHelper(p parser, rhs bool) (*data.Ers, data.Maybe[Typ]) {
	var termRes data.Either[data.Ers, Typ]
	if rhs && lookahead1(p, token.Dot) {
		es, acc := parseAccess(p)
		if es != nil {
			return es, data.Nothing[Typ](p)
		}
		termRes = data.Ok[Typ](acc)
	} else if lookahead1(p, exprAtomLAs...) {
		termRes = data.Cases(parseExprAtom(p), data.PassErs[Typ], exprAtomAsExprRes)
	} else if lookahead1(p, typeTermExceptionLAs...) {
		termRes = parseTypeTermException(p) // "_", "()", "[]", or "="
	} else if lookahead1(p, token.LeftParen, token.LeftBrace) {
		termRes = parseEnclosedType(p)
	} else {
		return nil, data.Nothing[Typ](p)
	}

	es, term, isRight := termRes.Break()
	if !isRight {
		return &es, data.Nothing[Typ](p)
	}
	return nil, data.Just(term)
}

// ASSUMPTION: the current token is Either "_", "()", "[]", "="
func parseTypeTermException(p parser) data.Either[data.Ers, Typ] {
	tok, found := getKeywordAtCurrent(p, token.Underscore, dropNone)
	if found {
		return data.Ok(Typ(data.EOne[wildcard](tok)))
	}

	tok, found = getKeywordAtCurrent(p, token.Equal, dropNone)
	if found {
		return data.Ok[Typ](data.EOne[name](tok))
	}

	tok, found = getKeywordAtCurrent(p, token.EmptyBracketEnclosure, dropNone)
	if found {
		return data.Ok[Typ](data.EOne[name](tok))
	}

	tok, found = getKeywordAtCurrent(p, token.EmptyParenEnclosure, dropNone)
	if !found {
		panic("expected '_', '()', '[]', '='") // input was not validated before calling
	}
	return data.Ok[Typ](data.EOne[unitType](tok))
}

func parseJustAppTypeOrJustType(p parser, lhs Typ, enclosed bool) (*data.Ers, Typ) {
	if lhs == nil {
		panic("lhs cannot be nil")
	}

	es, types, has2ndTerm := parseOneOrMore(p, lhs, enclosedDependentIt(enclosed), maybeParseTypeTermRhs)
	if es != nil {
		return es, nil
	}

	if !has2ndTerm {
		return nil, types.Head()
	}

	// otherwise, construct app type
	app, _ := data.NonEmptyToAppLikePair[TypeApp](types).Break()
	return nil, app
}

// rule:
//
//	```
//	forall binding = "forall", {"\n"}, forall binders, {"\n"}, "in", {"\n"}
//	```
func parseForallBinding(p parser) data.Either[data.Ers, forallBinders] {
	forallTok, found := getKeywordAtCurrent(p, token.Forall, dropAfter)
	if !found {
		panic("expected 'forall'") // input was not validated before calling
	}
	return parseForallBinders(p, forallTok)
}

func parseOptionalModality(p parser) data.Maybe[modality] {
	mModality := data.Nothing[modality](p)
	// parse optional multiplicity modality
	mode, found := getKeywordAtCurrent(p, token.Erase, dropAfter)
	if !found {
		mode, found = getKeywordAtCurrent(p, token.Once, dropAfter)
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
func parseEnclosedType(p parser) (out data.Either[data.Ers, Typ]) {
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
		return data.PassErs[Typ](es)
	} else if !mModality.IsNothing() {
		// validate that the it's a single term and that single term is an identifier
		term := lhs.Head()
		valid := lhs.Len() == 1 && (nodeType.LowerIdent.Match(term) || nodeType.UpperIdent.Match(term))
		if !valid {
			return data.Fail[Typ](ExpectedModalId, lhs)
		}
	}

	colon, foundColon := getKeywordAtCurrent(p, token.Colon, dropBeforeAndAfter)
	if !foundColon && !mModality.IsNothing() {
		return data.Fail[Typ](ExpectedTypeSig, p)
		// multiplicity requires term to be typed
	} else if foundColon {
		return parseEnclosedRhs(p, lhs, opener, closerType, implicit, mModality, colon)
	} else if lhs.Tail().Len() == 0 {
		return closeEnclosedTyp(p, opener, closerType)(lhs.Head())
	} else {
		return closeEnclosedTyp(p, opener, closerType)(lhs)
	}
}

func parseEnclosedRhs(p parser, lhs innerTypeSeq, opener api.Token, closerType token.Type, implicit bool, mModality data.Maybe[modality], colon api.Token) data.Either[data.Ers, Typ] {
	// parse typing and assemble inner type signature
	ty := parseType(p, true)
	typing := data.Cases(ty, data.PassErs[innerTyping], assembleInnerTyping(mModality, lhs, colon))
	attached := optionalAttachDefaultExpr(p, implicit, typing)
	construct := closeEnclosedTyp(p, opener, closerType)
	return data.Cases(attached, data.PassErs[Typ], construct)
}

// panics if colonEqual is nil
func getColonEqualAtCurrent(p parser, colonEqual *api.Token) (found bool) {
	if colonEqual == nil {
		panic("nil pointer for ':=' token")
	}
	*colonEqual, found = getKeywordAtCurrent(p, token.ColonEqual, dropBeforeAndAfter)
	return found
}

func optionalAttachDefaultExpr(p parser, implicit bool, typing data.Either[data.Ers, innerTyping]) data.Either[data.Ers, Typ] {
	var colonEqual api.Token
	noDefault := !(data.IsRight(typing) && implicit && getColonEqualAtCurrent(p, &colonEqual))
	if noDefault {
		return data.Cases(typing, data.PassErs[Typ], fun.Compose(data.Ok, (innerTyping).asTyp))
	}

	es, de, isDE := ParseExpr(p).Break()
	if !isDE {
		return data.PassErs[Typ](es)
	}
	return data.Cases(typing, data.PassErs[Typ], appendDefaultExpr(de))
}

func appendDefaultExpr(de Expr) func(innerTyping) data.Either[data.Ers, Typ] {
	return func(t innerTyping) data.Either[data.Ers, Typ] {
		return data.Ok[Typ](data.EMakePair[implicitTyping](t, data.EOne[defaultExpr](de)))
	}
}

func closeEnclosedTyp(p parser, opener api.Token, closerType token.Type) func(Typ) data.Either[data.Ers, Typ] {
	return func(t Typ) data.Either[data.Ers, Typ] {
		t = t.updatePosTyp(opener)
		closer, found := getKeywordAtCurrent(p, closerType, dropBefore)
		if !found {
			isRp := closerType.Match(token.RightParen.Make())
			return data.Fail[Typ](ifThenElse(isRp, ExpectedRightParen, ExpectedRightBrace), p)
		}
		et := enclosedType{implicit: token.LeftBrace.Match(opener), Typ: t}
		// NOTE: it's very important to call update on the `typ` field--if it's called on `et`, then the
		// enclosed type gets enclosed a second time ... lol
		et.Typ = et.Typ.updatePosTyp(closer)
		return data.Ok[Typ](et)
	}
}

func getNames(mModality data.Maybe[modality], lhs innerTypeSeq) data.Either[data.Ers, innerTypingNames] {
	var terms data.NonEmpty[name]
	fst := lhs.Head()
	tail := lhs.Tail()
	if n, ok := fst.(name); ok {
		terms = data.Singleton(n)
	} else {
		return data.Fail[innerTypingNames](ExpectedName, fst)
	}

	for tail.Len() > 0 {
		ty, _ := tail.Head().Break()
		if n, ok := ty.(name); ok {
			terms = terms.Snoc(n)
		} else {
			return data.Fail[innerTypingNames](ExpectedName, tail.Head())
		}

		tail, _ = tail.Tail().Break()
	}

	var names innerTypingNames
	if mode, ok := mModality.Break(); ok {
		if terms.Len() != 1 {
			return data.Fail[innerTypingNames](IllegalMultipleModality, lhs)
		}

		names = data.Inl[data.NonEmpty[name]](data.MakePair(mode, terms.Head()))
	} else {
		names = data.Inr[data.Pair[modality, name]](terms)
	}
	return data.Ok(names)
}

func assembleInnerTyping(mModality data.Maybe[modality], lhs innerTypeSeq, colon api.Token) func(Typ) data.Either[data.Ers, innerTyping] {
	return func(t Typ) data.Either[data.Ers, innerTyping] {
		es, names, isNames := getNames(mModality, lhs).Break()
		if !isNames {
			return data.PassErs[innerTyping](es)
		}

		it := data.EMakePair[innerTyping](names, t)
		it.Position = api.WeakenRangeOver[api.Node](colon, mModality, it)
		return data.Ok(it)
	}
}

// rule:
//
//	```
//	inner type terms = type tail, [{{"\n"}, ",", {"\n"}, type tail}, [{"\n"}, ","]] ;
//	```
func parseInnerTypeTerms(p parser) data.Either[data.Ers, innerTypeSeq] {
	maybeParse := fun.BinBind1st_PairTarget(parseMaybeTypeTail, true)
	return parseSepSequenced[innerTypeSeq](p, ExpectedType, token.Comma, maybeParse)
}
