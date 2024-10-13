package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/api/util/fun"
	"github.com/petersalex27/yew/common/data"
	"github.com/petersalex27/yew/internal/common"
	t "github.com/petersalex27/yew/internal/parser/typ"
)

var (
	passBodyErrors = fun.BiConstant[Parser](data.Inl[bodyElement, data.Ers])
)

// parses a type signature
//
// rule:
//
//	```
//	typing = name, {"\n"}, ":", {"\n"}, typ ;
//	```
//
// NOTE: a type signature, syntactically speaking, is different from a "typing". A type signature is
// a subset of the more general "typing". A type signature can only have a "name" node appear to the
// left of ':'.
func parseTypeSig(p Parser) data.Either[data.Ers, typing] {
	var t typing
	n, isN := maybeParseName(p).Break()
	if !isN {
		return data.Fail[typing](ExpectedName, p)
	}
	t.Position = t.Update(n)

	p.dropNewlines()
	if !parseKeywordAtCurrent(p, token.Colon, &t.Position) {
		return data.Fail[typing](ExpectedTyping, n)
	}

	es, ty, isTy := ParseType(p).Break()
	if !isTy {
		return data.PassErs[typing](es)
	}

	t.Position = t.Update(ty)
	t.typing = data.MakePair(n, ty)
	t.annotations = data.Nothing[annotations](p)
	t.visibility = data.Nothing[visibility](p)
	return data.Ok(t)
}

// rule:
//
//	```
//	type def body =
//		"impossible"
//		| "(",{"\n"},[annotations_],type cons,{{"\n"},[annotations_],type cons},{"\n"},")"
//		| "(", {"\n"}, "impossible", {"\n"}, ")"
//		| [annotations_], type cons ;
//	type cons = typing ;
//	```
func parseTypeDefBody(p Parser, typ typing) data.Either[data.Ers, data.Pair[typing, typeDefBody]] {
	if matchCurrentImpossible(p) {
		// constant "impossible" case
		td := data.Inr[data.NonEmpty[typeConstructor]](impossible{data.One(p.current())})
		p.advance()

		return data.Ok(data.MakePair(typ, td))
	} else if !matchCurrentLeftParen(p) {
		// single constructor case
		return data.Cases(parseTypeDefBodyTypeCons(p),
			data.Inl[data.Pair[typing, typeDefBody], data.Ers],
			assembleSingleConstructor(typ),
		)
	}

	// enclosed case
	p.advance()
	p.dropNewlines()
	return repeat(notMatchCurrent(token.RightParen),
		parseTypeDefBodyTypeCons,
		assembleTypeConstructors(typ),
	)(p)
}

func assembleSingleConstructor(typ typing) func(typeConstructor) data.Either[data.Ers, data.Pair[typing, typeDefBody]] {
	return func(tc typeConstructor) data.Either[data.Ers, data.Pair[typing, typeDefBody]] {
		return data.Inr[data.Ers](data.MakePair(typ, data.Inl[impossible](data.Singleton(tc))))
	}
}

func assembleTypeConstructors(typ typing) func(data.List[typeConstructor]) data.Either[data.Ers, data.Pair[typing, typeDefBody]] {
	return func(l data.List[typeConstructor]) data.Either[data.Ers, data.Pair[typing, typeDefBody]] {
		if cons, just := l.Head().Break(); !just {
			return data.Fail[data.Pair[typing, typeDefBody]](ExpectedTypeConstructor, l)
		} else {
			tdb := data.Inl[impossible](data.Singleton(cons).Append(l.Elements()[1:]...))
			return data.Ok(data.MakePair(typ, tdb))
		}
	}
}

func parseTypeDefBodyTypeCons(p Parser) data.Either[data.Ers, typeConstructor] {
	return runCases(p, parseAnnotations, passParseErs[typeConstructor], parseTypeConstructor)
}

// rule:
//
//	```
//	type cons = typing ;
//	```
func parseTypeConstructor(p Parser, as data.Maybe[annotations]) data.Either[data.Ers, typeConstructor] {
	return twoCases(
		data.SeqResult[name](ExpectedTypeConstructorName)(maybeParseName(p)),
		runCases(p, ParseType, passParseErs[typ], passParseRight[data.Ers, typ]),
		constInl[typeConstructor, data.Ers, data.Ers],
		constInl[typeConstructor, data.Ers, typ],
		constInlSwap[typeConstructor, name, data.Ers],
		assembleTypeConstructor(as),
	)
}

func constInl[c, a, b api.Node](x a, _ b) data.Either[a, c] { return data.Inl[c](x) }

func constInlSwap[c, a, b api.Node](_ a, y b) data.Either[b, c] { return data.Inl[c](y) }

func assembleTypeConstructor(as data.Maybe[annotations]) func(n name, typ typ) data.Either[data.Ers, typeConstructor] {
	return func(n name, typ typ) data.Either[data.Ers, typeConstructor] {
		tc := typeConstructor{constructor: data.MakePair(n, typ)}
		(&tc).annotate(as)
		return data.Ok(tc)
	}
}

// rule:
//
//	```
//	data.Maybe visibility = [("open" | "public"), {"\n"}] ;
//	```
func parseOptionalVisibility(p Parser) (mv data.Maybe[visibility]) {
	if lookahead1(p, token.Open, token.Public) {
		visibilityToken := p.current()
		mv = data.Just(data.EOne[visibility](visibilityToken))
		p.advance()
		p.dropNewlines()
	} else {
		mv = data.Nothing[visibility](p)
	}
	return mv
}

func attachVisibility(vis data.Maybe[visibility]) func(be bodyElement) data.Either[data.Ers, bodyElement] {
	return func(be bodyElement) data.Either[data.Ers, bodyElement] {
		if vbe, ok := be.(visibleBodyElement); !ok && !vis.IsNothing() {
			return data.Fail[bodyElement](IllegalVisibilityTarget, be)
		} else {
			return data.Ok(vbe.setVisibility(vis))
		}
	}
}

func parseOneBodyElement(p Parser) data.Either[data.Ers, bodyElement] {
	vis := parseOptionalVisibility(p)
	return data.Cases(parseBasicBodyStructure(p, vis), data.PassErs[bodyElement], attachVisibility(vis))
}

func setResultAnnotation(as data.Maybe[annotations]) func(p Parser, be bodyElement) data.Either[data.Ers, bodyElement] {
	return func(p Parser, be bodyElement) data.Either[data.Ers, bodyElement] {
		return data.Ok(be.setAnnotation(as))
	}
}

// rule:
//
//	```
//	body element = def | visible body element ;
//	```
func parseBodyElement(p Parser, as data.Maybe[annotations]) data.Either[data.Ers, bodyElement] {
	return runCases(p, parseOneBodyElement, passBodyErrors, setResultAnnotation(as))
}

// rule:
//
//	```
//	body = {{"\n"}, [annotations_], body elem} ;
//	```
func parseBody(p Parser) (theBody data.Either[data.Ers, data.Maybe[body]], mFooterAnnots data.Maybe[annotations]) {
	const smallBodyCap int = 16
	sourceBody := body{data.Nil[bodyElement](smallBodyCap)}
	var es data.Ers
	var isAnnots bool

	has2ndTerm := false

	for {
		p.dropNewlines()
		es, mFooterAnnots, isAnnots = parseAnnotations(p).Break()
		if !isAnnots { // not just annotations & not nothing -> void
			return data.PassErs[data.Maybe[body]](es), mFooterAnnots
		} else if isAnnots && lookahead1(p, token.EndOfTokens) {
			// possibly parsed footer annotations, return body and "Maybe" footer annotations
			theBody = data.Ok(data.Just(sourceBody))
			break
		}

		esBE, be, isBE := parseBodyElement(p, mFooterAnnots).Break()
		if !isBE {
			return data.PassErs[data.Maybe[body]](esBE), mFooterAnnots
		}
		be = be.setAnnotation(mFooterAnnots)

		sourceBody.List = sourceBody.Snoc(be)
		has2ndTerm = true
	}

	if !has2ndTerm { // no body elements parsed, return Nothing instead of empty list
		theBody = data.Ok(data.Nothing[body](sourceBody))
	}
	return theBody, mFooterAnnots
}

func assembleDef(pat pattern) func(defBody) data.Either[data.Ers, def] {
	return func(db defBody) data.Either[data.Ers, def] {
		db.Either = db.Update(pat)
		return data.Ok(def{
			pattern:  pat,
			defBody:  db,
			Position: db.GetPos(),
		})
	}
}

// here for consistency b/w 'parseDef' and 'maybeParseDef' for parsing patterns
func _parseDef_ParsePattern(p Parser) (es *data.Ers, pat pattern) {
	e, pat, isPat := ParsePattern(p).Break()
	if !isPat {
		return &e, nil
	}
	return &e, pat
}

func _finishParseDef(p Parser, pat pattern) data.Either[data.Ers, def] {
	p.dropNewlines()
	return data.Cases(parseDefBody(p), data.Inl[def, data.Ers], assembleDef(pat))
}

// rule:
//
//	```
//	def = pattern, {"\n"}, def body ;
//	```
func parseDef(p Parser) data.Either[data.Ers, def] {
	es, pat := _parseDef_ParsePattern(p)
	if es != nil {
		return data.PassErs[def](*es)
	}

	return _finishParseDef(p, pat)
}

func maybeParseDef(p Parser) (*data.Ers, data.Maybe[def]) {
	// try to parse def: this is a bit of a doozy since patterns can be arbitrary large and appear
	// in the lhs of multiple mutually exclusive production rules--meaning we can't look ahead to
	// determine if it's a valid def. We will record the Position of the current token and then try
	origin := getOrigin(p)
	p = p.markOptional()
	es, pat := _parseDef_ParsePattern(p)
	p = p.demarkOptional()
	if es != nil {
		// reset the origin and return data.Nothing (no error)
		p = resetOrigin(p, origin)
		return nil, data.Nothing[def](p)
	} // else, keep origin and enforce non-data.Nothing result. This must be a def

	es2, d, isDef := _finishParseDef(p, pat).Break()
	if !isDef {
		return &es2, data.Nothing[def](p)
	}
	return nil, data.Just(d)
}

// rule:
//
//	```
//	def body = (with clause | "=", {"\n"}, expr), [where clause] | "impossible" ;
//	```
func parseDefBody(p Parser) data.Either[data.Ers, defBody] {
	return parsePatternBoundBody(p, token.Equal)
}

// a more general version of `parseDefBody` that allows for the binding token to be chosen
func parsePatternBoundBody(p Parser, bindingTokenType token.Type) data.Either[data.Ers, defBody] {
	if imp, found := getKeywordAtCurrent(p, token.Impossible); found {
		return data.Ok(data.EInl[defBody](data.EOne[impossible](imp)))
	}

	var possibleLeft data.Either[data.Ers, data.Either[withClause, expr]]
	if bindingToken, found := getKeywordAtCurrent(p, bindingTokenType); found {
		construct := fun.Compose(data.Ok, data.Inr[withClause, expr])
		possibleLeft = data.Cases(ParseExpr(p), data.PassErs[data.Either[withClause, expr]], construct)
		possibleLeft = possibleLeft.Update(bindingToken)
	} else if matchCurrentWith(p) {
		construct := fun.Compose(data.Ok, data.Inl[expr, withClause])
		possibleLeft = data.Cases(parseWithClause(p), data.PassErs[data.Either[withClause, expr]], construct)
	}

	return runCases(p, fun.Constant[Parser](possibleLeft), passParseErs[defBody], runDefBodyWhereClause)
}

func runDefBodyWhereClause(p Parser, possibleLeft data.Either[withClause, expr]) data.Either[data.Ers, defBody] {
	constructPossibleBody := fun.BiCompose(data.EInr[defBody], data.EMakePair[defBodyPossible])
	construct := fun.Compose(data.Ok, constructPossibleBody(possibleLeft))
	return data.Cases(parseOptionalWhereClause(p), data.PassErs[defBody], construct)
}

// rule:
//
//	```
//	where clause = {"\n"}, "where", {"\n"}, where body ;
//	```
func parseOptionalWhereClause(p Parser) data.Either[data.Ers, data.Maybe[whereClause]] {
	whereToken, found := getKeywordAtCurrent(p, token.Where)
	if !found {
		return data.Ok(data.Nothing[whereClause](p)) // no where clause, return data.Nothing
	}

	es, whereBody, isWhereBody := parseWhereBody(p).Break()
	if !isWhereBody {
		return data.PassErs[data.Maybe[whereClause]](es)
	}

	whereBody.Position = whereBody.Update(whereToken)
	return data.Ok(data.Just[whereClause](data.EOne[whereClause](whereBody)))
}

// rule:
//
//	```
//	where body = main elem | "(", {"\n"}, main elem, {"\n", main elem}, {"\n"}, ")" ;
//	```
func parseWhereBody(p Parser) data.Either[data.Ers, whereBody] {
	leftParen, found := getKeywordAtCurrent(p, token.LeftParen)
	me := parseMainElem(p)
	return data.Cases(me, data.PassErs[whereBody], func(me mainElement) data.Either[data.Ers, whereBody] {
		es, ne, _ := parseOneOrMore(p, me, true, parseMaybeMainElem)
		if es != nil {
			return data.PassErs[whereBody](*es)
		}
		wb := whereBody{ne}
		if found {
			wb.Position = wb.Update(leftParen)
			rp, found := getKeywordAtCurrent(p, token.RightParen)
			if !found {
				return data.Fail[whereBody](ExpectedRightParen, leftParen)
			}
			wb.Position = wb.Update(rp)
		}
		return data.Ok(whereBody{ne})
	})
}

// rule:
//
//	```
//	main elem = def | spec def | spec inst | type def | type alias | typing | syntax ;
//	```
func parseMainElem(p Parser) data.Either[data.Ers, mainElement] {
	es, me := parseMaybeMainElem(p)
	if es != nil {
		return data.PassErs[mainElement](*es)
	} else if unit, just := me.Break(); !just {
		return data.Fail[mainElement](ExpectedMainElement, p)
	} else {
		return data.Ok(unit)
	}
}

// helper function for `parseMainElem` that transforms an `data.Either[data.Ers, a]` to an
// `data.Either[data.Ers, mainElem]` where `a` implements `mainElem`
func knownCase[a mainElement](elemRes data.Either[data.Ers, a]) data.Either[data.Ers, mainElement] {
	return data.Cases(elemRes, data.Inl[mainElement, data.Ers], fun.Compose(data.Ok, (a).pureMainElem))
}

// parses a spec, alias, inst, or syntax definition when given the corresponding token type `tt`
//
// if `tt` is not data.One of the following this function will panic:
//   - token.Spec
//   - token.Alias
//   - token.Inst
//   - token.Syntax
func parseKnownMainElem(p Parser, tt token.Type) data.Either[data.Ers, mainElement] {
	switch tt {
	case token.Spec:
		return knownCase(parseSpecDef(p))
	case token.Alias:
		return knownCase(parseTypeAlias(p))
	case token.Inst:
		return knownCase(parseSpecInst(p))
	case token.Syntax:
		return knownCase(parseSyntax(p))
	default:
		panic("illegal argument")
	}
}

// rule:
//
//	```
//	"alias", {"\n"}, name, {"\n"}, "=", {"\n"}, type ;
//	```
func parseTypeAlias(p Parser) data.Either[data.Ers, typeAlias] {
	aliasToken, found := getKeywordAtCurrent(p, token.Alias)
	if !found {
		return data.Fail[typeAlias](ExpectedTypeAlias, p)
	}

	n, isSomething := maybeParseName(p).Break()
	if !isSomething {
		return data.Fail[typeAlias](ExpectedTypeAliasName, aliasToken)
	}

	p.dropNewlines()
	equalToken, found := getKeywordAtCurrent(p, token.Equal)
	if !found {
		return data.Fail[typeAlias](ExpectedAliasBinding, n)
	}

	p.dropNewlines()
	es, ty, isTy := ParseType(p).Break()
	if !isTy {
		return data.PassErs[typeAlias](es)
	}

	return data.Ok(constructAlias(aliasToken, n, equalToken, ty))
}

// rule:
//
//	```
//	syntax = "syntax", {"\n"}, syntax rule, {"\n"}, "=", {"\n"}, expr ;
//	```
func parseSyntax(p Parser) data.Either[data.Ers, syntax] {
	var syn syntax
	syntaxToken, found := getKeywordAtCurrent(p, token.Syntax)
	if !found {
		return data.Fail[syntax](ExpectedSyntax, p)
	}
	syn.Position = syn.Update(syntaxToken)

	es, rule, isRule := parseSyntaxRule(p).Break()
	if !isRule {
		return data.PassErs[syntax](es)
	}

	p.dropNewlines()
	equalToken, found := getKeywordAtCurrent(p, token.Equal)
	if !found {
		return data.Fail[syntax](ExpectedSyntaxBinding, rule)
	}
	syn.Position = syn.Update(equalToken)

	esE, e, isE := ParseExpr(p).Break()
	if !isE {
		return data.PassErs[syntax](esE)
	}

	syn.rule = data.MakePair(rule, e)
	syn.Position = syn.Update(syn.rule)
	return data.Ok(syn)
}

// rule:
//
//	```
//	spec def = "spec", {"\n"}, def ;
func parseSyntaxRule(p Parser) data.Either[data.Ers, syntaxRule] {
	var ruleInsides data.NonEmpty[syntaxSymbol]
	has1stTerm := false
	for { // loop until there aren't any syntax symbols in view
		es, sym := maybeParseSyntaxSymbol(p)
		if es != nil {
			return data.PassErs[syntaxRule](*es)
		} else if unit, just := sym.Break(); !just {
			break
		} else {
			ruleInsides = ruleInsides.Snoc(unit)
			has1stTerm = true
		}
	}

	if !has1stTerm {
		return data.Fail[syntaxRule](ExpectedSyntaxRule, p)
	}
	return data.Inr[data.Ers](syntaxRule{ruleInsides})
}

// assumes the token is the correct type (i.e., `token.RawStringValue`)
func validRawSyntaxSymbol(t api.Token) bool {
	s := t.String()
	res := common.NonInfixName.Match(s)
	return res != nil && len(*res) == len(s)
}

// rule:
//
//	```
//	syntax symbol = ident | "{", {"\n"}, ident, {"\n"}, "}" | raw keyword ;
//	raw keyword = ? RAW STRING OF JUST A VALID NON INFIX ident OR symbol ? ;
//	```
func maybeParseSyntaxSymbol(p Parser) (*data.Ers, data.Maybe[syntaxSymbol]) {
	if matchCurrentRawString(p) && validRawSyntaxSymbol(p.current()) {
		return nil, data.Just[syntaxSymbol](parseRawSyntaxSymbolFromCurrent(p))
	} else if lookahead2(p, boundIdentL2...) {
		es, sym, isSym := parseSyntaxBindingSymbol(p).Break()
		if !isSym {
			return &es, data.Nothing[syntaxSymbol](p)
		}
		return nil, data.Just[syntaxSymbol](sym)
	} else if unit, just := parseIdent(p).Break(); just {
		return nil, data.Just[syntaxSymbol](data.Inl[syntaxRawKeyword](unit))
	}
	return nil, data.Nothing[syntaxSymbol](p)
}

// parses a raw syntax symbol (a non-infix name inside of '`' characters w/o anything else)
func parseRawSyntaxSymbolFromCurrent(p Parser) syntaxSymbol {
	rawKey := p.current()
	p.advance()
	key := data.EOne[syntaxRawKeyword](data.EOne[rawString](rawKey))
	return data.Inr[ident](key)
}

// rule:
//
//	```
//	syntax binding symbol = "{", {"\n"}, ident, {"\n"}, "}" ;
//	```
func parseSyntaxBindingSymbol(p Parser) data.Either[data.Ers, syntaxSymbol] {
	lb, _ := getKeywordAtCurrent(p, token.LeftBrace)

	id, isSomething := parseIdent(p).Break()
	if !isSomething {
		return data.Fail[syntaxSymbol](ExpectedSyntaxBindingId, lb)
	}

	id = id.Update(lb)

	rb, found := getKeywordAtCurrent(p, token.RightBrace)
	if !found {
		return data.Fail[syntaxSymbol](ExpectedRightBrace, id)
	}

	id = id.Update(rb)

	symbol := data.Inl[syntaxRawKeyword](id)
	return data.Inr[data.Ers](symbol)
}

type structureParseAttempt uint32

const (
	attemptedNothing structureParseAttempt = 0 // always excluded
	attemptedDef                           = 1 << (iota - 1)
	attemptedSpecDef
	attemptedSpecInst
	attemptedTypeDef
	attemptedTypeAlias
	attemptedTyping
	attemptedSyntax
	attemptedGenericStructure = attemptedSyntax<<1 - 1 // all of the above
)

type exclusionMask uint32

const (
	publicVisModExclusion exclusionMask = exclusionMask(attemptedGenericStructure &^ attemptedDef)
	openVisModExclusion   exclusionMask = exclusionMask(attemptedTypeDef)
	xxx                   exclusionMask = attemptedGenericStructure &^ openVisModExclusion
	noVisModExclusion     exclusionMask = 0 // this will, literally, exclude *attempting data.Nothing* from being valid
	excludeAll            exclusionMask = exclusionMask(attemptedGenericStructure)
)

func (attempt structureParseAttempt) isNecessarilyExcludedBy(mask exclusionMask) bool {
	// generic structure is only excluded by `excludeAll`--think of this function as telling us that something is necessarily excluded
	// if it's only possible, then this will return false
	return (attempt &^ structureParseAttempt(mask)) == attempt
}

func recastMainElementAsBodyElement[a mainElement](me a) data.Either[data.Ers, bodyElement] {
	return data.Inr[data.Ers, bodyElement](me)
}

func attemptedWhat(tt token.Type) structureParseAttempt {
	switch tt {
	case token.Spec:
		return attemptedSpecDef
	case token.Inst:
		return attemptedSpecInst
	case token.Alias:
		return attemptedTypeAlias
	case token.Syntax:
		return attemptedSyntax
	default: // default to a generic structure
		return attemptedGenericStructure
	}
}

func (attempted structureParseAttempt) getErrorMessageForAttempted(exclusion exclusionMask) string {
	if !attempted.isNecessarilyExcludedBy(exclusion) { // not excluded from having the given visibility modifier
		return ""
	}
	if exclusion == openVisModExclusion {
		switch attempted {
		case attemptedSpecInst,
			attemptedSpecDef,
			attemptedTypeAlias,
			attemptedSyntax:
			return IllegalOpenModifier
		case attemptedDef:
			return IllegalVisibleDef
		case attemptedTyping:
			return IllegalOpenModifierTyping
		}
	} else if exclusion == publicVisModExclusion {
		// for 'public' visibility modifier
		switch attempted {
		case attemptedDef:
			return IllegalVisibleDef
		}
	}
	return UnexpectedStructure
}

func parseBasicBodyStructureHelper(p Parser) (res data.Either[data.Ers, bodyElement], attempted structureParseAttempt) {
	/*
	 * `spec def`, `spec inst`, `type alias`, and `syntax` can all be distinguished with lookahead of 1.
	 * `typing` or `type def` can be distinguished from `def` with lookahead of 2.
	 * `typing` and `type def` cannot be distinguished with a fixed lookahead; however, `type def`'s lhs
	 * is the `typing` rule, so we can parse `typing` and then check for `where` to determine if it is
	 * a `type def`
	 *
	 * strategy:
	 * 	1. lookahead 1 for `spec def`, `spec inst`, `type alias`, and `syntax`
	 *	2. otherwise, lookahead 2 for `typing` and `type def`, then check for `where` to determine if it is a `type def`
	 *	3. otherwise, if 1. and 2. data.Fail, try to parse `def`
	 */

	// 1. lookahead 1
	if tt, found := lookahead1Report(p, token.Spec, token.Inst, token.Alias, token.Syntax); found {
		res := parseKnownMainElem(p, tt)
		return data.Cases(res, data.Inl[bodyElement, data.Ers], recastMainElementAsBodyElement), attemptedWhat(tt)
	}

	// 2. lookahead 2
	if found := lookahead2(p, typingL2...); found {
		// put as attempted typing since typing is the only data.One that actually matters. `attempted` is
		// used for visibility related error messages, but a type def can have any visibility modifier
		res := runCases(p, parseTypeSig, passParseErs[mainElement], parseTypeDefOrTyping)
		return data.Cases(res, data.PassErs[bodyElement], recastMainElementAsBodyElement), attemptedTyping
	}

	// 3. try to parse `def`
	return data.Cases(parseDef(p), data.PassErs[bodyElement], recastMainElementAsBodyElement), attemptedDef
}

func visToExclusion(v data.Maybe[visibility]) exclusionMask {
	if v.IsNothing() {
		return noVisModExclusion
	}
	theV, _ := v.Break()
	if token.Open.Match(theV) {
		return openVisModExclusion
	} else if token.Public.Match(theV) {
		return publicVisModExclusion
	}
	return excludeAll
}

var typeToStructureTypeMap = map[api.NodeType]structureParseAttempt{
	t.Def:         attemptedDef,
	t.SpecDef:     attemptedSpecDef,
	t.SpecInst:    attemptedSpecInst,
	t.TypeDef:     attemptedTypeDef,
	t.TypeAlias:   attemptedTypeAlias,
	t.Typing:      attemptedTyping,
	t.InnerTyping: attemptedTyping,
	t.Syntax:      attemptedSyntax,
}

func strengthenStructureType(b bodyElement, weaker structureParseAttempt) structureParseAttempt {
	stronger, found := typeToStructureTypeMap[b.Type()]
	if !found {
		return weaker
	}
	return stronger
}

func parseBasicBodyStructure(p Parser, vis data.Maybe[visibility]) data.Either[data.Ers, bodyElement] {
	startPosition := p.current().GetPos()
	res, attempted := parseBasicBodyStructureHelper(p)
	lhs, rhs, isRight := res.Break()
	// if res has no errors, we can get a more specific error message if needed
	if isRight {
		attempted = strengthenStructureType(rhs, attempted)
	}

	// see if an additional error message is needed
	msg := attempted.getErrorMessageForAttempted(visToExclusion(vis))
	var e data.Err
	if msg != "" {
		e = mkErr(msg, api.WeakenRangeOver[api.Positioned](res, startPosition))
	}

	if e.Msg() == "" {
		return res
	} else if !isRight {
		res = data.PassErs[bodyElement](lhs.Snoc(e)) // add the error to the existing errors
	} else {
		// replace the result with the error since it was marked as excluded
		// due to visibility modifier attempted to be given to it
		res = data.PassErs[bodyElement](data.Nil[data.Err]().Snoc(e))
	}
	return res
}

// rule:
//
//	```
//	main elem = [annotations_], (def | spec def | spec inst | type def | type alias | typing | syntax) ;
//	```
func parseMaybeMainElem(p Parser) (*data.Ers, data.Maybe[mainElement]) {
	esAnnots, mAnnots, isAnnots := parseAnnotations(p).Break()
	if !isAnnots {
		return &esAnnots, data.Nothing[mainElement](p) // error parsing optional annotations
	}

	res, _ := parseBasicBodyStructureHelper(p)
	lhs, rhs, isRight := res.Break()
	if !isRight {
		return &lhs, data.Nothing[mainElement](p)
	}

	me, isME := rhs.setAnnotation(mAnnots).(mainElement)
	if !isME {
		panic("illegal state")
	}
	return nil, data.Just[mainElement](me)
}

func parseTypeDefOrTyping(p Parser, t typing) data.Either[data.Ers, mainElement] {
	where, found := getKeywordAtCurrent(p, token.Where)
	if !found {
		return data.Inr[data.Ers](mainElement(t))
	}
	es, tdb, isTbd := parseTypeDefBody(p, t).Break()
	if !isTbd {
		return data.PassErs[mainElement](es)
	}

	esDeriving, mDeriving, isDeriving := parseOptionalDerivingClause(p).Break()
	if !isDeriving {
		return data.PassErs[mainElement](esDeriving)
	}

	return data.Ok(mainElement(typeDef{
		// this serves two purposes: adding the value and incorporating the where Position
		annotations: data.Nothing[annotations](where),
		// constructed and b/c 'data.Nothing' needs a Position arg
		visibility: data.Nothing[visibility](where),
		typedef:    tdb,
		deriving:   mDeriving,
	}))
}

// rule:
//
//	```
//	with clause = "with", {"\n"}, pattern, {"\n"}, "of", {"\n"}, with clause arms ;
//	```
func parseWithClause(p Parser) data.Either[data.Ers, withClause] {
	with := withClause{}

	withToken, found := getKeywordAtCurrent(p, token.With)
	if !found {
		return data.Fail[withClause](ExpectedWithClause, p)
	}
	with.Position = with.Update(withToken)

	es, pat, isPat := ParsePattern(p).Break()
	if !isPat {
		return data.PassErs[withClause](es)
	}
	with.Position = with.Update(pat)

	p.dropNewlines()
	ofToken, found := getKeywordAtCurrent(p, token.Of)
	if !found {
		return data.Fail[withClause](ExpectedOf, pat)
	}
	with.Position = with.Update(ofToken)

	esArms, arms, isArms := parseWithClauseArms(p).Break()
	if !isArms {
		return data.PassErs[withClause](esArms)
	}
	with.Position = with.Update(arms)
	with.Pair = data.MakePair(pat, arms)

	return data.Ok(with)
}

// rule:
//
//	```
//	with clause arms =
//		"(", {"\n"}, with clause arm, {{"\n"}, with clause arm}, {"\n"}, ")"
//		| with clause arm ;
//	```
func parseWithClauseArms(p Parser) data.Either[data.Ers, withClauseArms] {
	return parseGroup[withClauseArms, withClauseArm](p, ExpectedWithClauseArm, maybeParseWithClauseArm)
}

// rule:
//
//	```
//	with clause arm = [pattern, {"\n"}, "|", {"\n"}], pattern, {"\n"}, def body thick arrow ;
//	```
func maybeParseWithClauseArm(p Parser) (*data.Ers, data.Maybe[withClauseArm]) {
	wca := withClauseArm{}
	es, pat := maybeParsePattern(p, false)
	if es != nil { // pattern found, but error while parsing it
		return es, data.Nothing[withClauseArm](p)
	} else if pat.IsNothing() { // no pattern found, return data.Nothing
		return nil, data.Nothing[withClauseArm](p)
	}

	p.dropNewlines()
	// check for '|': if found, parse the second pattern; otherwise, parse the def body
	found := parseKeywordAtCurrent(p, token.Bar, &wca.Position)
	var first withArmLhs
	if found {
		pEs, patRhs, isPatRhs := ParsePattern(p).Break()
		if !isPatRhs {
			return &pEs, data.Nothing[withClauseArm](p)
		}
		wca.Position = wca.Update(patRhs)
		unit, _ := pat.Break()
		first = data.Inr[pattern](data.MakePair(unit, patRhs))
		p.dropNewlines()
	}

	esDb, db, isDb := parsePatternBoundBody(p, token.ThickArrow).Break()
	if !isDb {
		return &esDb, data.Nothing[withClauseArm](p)
	}
	wca.Position = wca.Update(db)
	wca.Pair = data.MakePair(first, db)
	return nil, data.Just(wca)
}
