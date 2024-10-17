package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/api/util/fun"
	"github.com/petersalex27/yew/common/data"
	"github.com/petersalex27/yew/internal/common"
	t "github.com/petersalex27/yew/internal/parser/typ"
)

// rule:
//
//	```
//	body = [annotations_], body elem, {then, [annotations_], body elem} ;
//	```
func parseBody(p Parser) (theBody data.Either[data.Ers, data.Maybe[body]]) {
	const smallBodyCap int = 16
	sourceBody := body{data.Nil[bodyElement](smallBodyCap)}

	has1stTerm := false

	origin := getOrigin(p)
	for {
		es, mFooterAnnots, isMAnnots := parseAnnotations_(p).Break()
		if !isMAnnots { // not just annotations & not nothing -> void
			return data.PassErs[data.Maybe[body]](es)
		} else if isMAnnots && lookahead1(p, token.EndOfTokens) {
			// at footer
			resetOrigin(p, origin)
			theBody = data.Ok(data.Just(sourceBody))
			break
		}

		esBe, be, isBe := parseBodyElement(p).Break()
		if !isBe {
			return data.PassErs[data.Maybe[body]](esBe)
		}

		// attach annotations to the body element
		if d, vbe, isVbe := be.Break(); isVbe {
			be = vbe.setAnnotation(mFooterAnnots).asBodyElement()
		} else {
			be = d.setAnnotation(mFooterAnnots).asBodyElement()
		}

		sourceBody.List = sourceBody.Snoc(be)
		has1stTerm = true
		origin = getOrigin(p)
		if !then(p) {
			theBody = data.Ok(data.Just(sourceBody))
			break
		}
	}

	if !has1stTerm { // no body elements parsed, return Nothing instead of empty list
		theBody = data.Ok(data.Nothing[body](sourceBody))
	}
	return theBody
}

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
	n, isN := maybeParseName(p).Break()
	if !isN {
		return data.Fail[typing](ExpectedName, p)
	}

	colon, found := getKeywordAtCurrent(p, token.Colon, dropBeforeAndAfter)
	if !found {
		return data.Fail[typing](ExpectedTypeJudgment, n)
	}

	es, ty, isTy := ParseType(p).Break()
	if !isTy {
		return data.PassErs[typing](es)
	}

	t := makeTyping(n, ty)
	t.Position = t.Position.Update(colon)
	return data.Ok(t)
}

// rule:
//
//	```
//	auto typing = "auto", {"\n"}, typing ;
//	```
func parseAutoDefTypeSig(p Parser) data.Either[data.Ers, mainElement] {
	auto, found := getKeywordAtCurrent(p, token.Auto, dropAfter)
	if !found {
		// should be impossible to reach this point; if it is, something is wrong
		//
		// could be a lot of things: tokens changed between before calling this function (when 'auto'
		// should flag this function for calling) and actually getting the token, or another possibility
		// is this function is called from the wrong caller
		//
		// regardless, this is almost certainly a bug in the caller
		return data.Fail[mainElement](ExpectedAuto, p)
	}

	es, ty, isTy := parseTypeSig(p).Break()
	if !isTy {
		return data.PassErs[mainElement](es)
	}
	me := ty.markAuto(auto).asMainElement()
	return data.Ok(me)
}

// rule:
//
//	```
//	type def body =
//		"impossible"
//		| "(", {"\n"}, [annotations_], type cons, {then, [annotations_], type cons}, {"\n"}, ")"
//		| [annotations_], type cons ;
//	```
func parseTypeDefBody(p Parser) data.Either[data.Ers, typeDefBody] {
	if matchCurrentImpossible(p) {
		// constant "impossible" case
		td := data.Inr[data.NonEmpty[typeConstructor]](impossible{data.One(p.current())})
		p.advance()

		return data.Ok(td)
	}

	lparen, found := getKeywordAtCurrent(p, token.LeftParen, dropAfter)

	// parse first type constructor
	var bod data.NonEmpty[typeConstructor]
	if es, tcs, isTC := parseTypeDefBodyTypeCons(p).Break(); !isTC {
		return data.PassErs[typeDefBody](es)
	} else if !found {
		return data.Ok(data.Inl[impossible](tcs))
	} else {
		bod = tcs
	}

	for then(p) && !token.RightParen.Match(p.current()) {
		if es, tcs, isTC := parseTypeDefBodyTypeCons(p).Break(); !isTC {
			return data.PassErs[typeDefBody](es)
		} else {
			bod = bod.Append(tcs.Elements()...)
		}
	}

	rparen, found := getKeywordAtCurrent(p, token.RightParen, dropBefore)
	if !found {
		return data.Fail[typeDefBody](ExpectedRightParen, lparen)
	}

	bod.Position = bod.Update(lparen).Update(rparen)
	return data.Ok(data.Inl[impossible](bod))
}

// rule:
//
//	```
//	type def body cons = [annotations_], type constructor ;
//	```
func parseTypeDefBodyTypeCons(p Parser) data.Either[data.Ers, data.NonEmpty[typeConstructor]] {
	es, mAs, isMas := parseAnnotations_(p).Break()
	if !isMas {
		return data.PassErs[data.NonEmpty[typeConstructor]](es)
	}

	es, tcs, isTcs := parseTypeConstructor(p, mAs).Break()
	if !isTcs {
		return data.PassErs[data.NonEmpty[typeConstructor]](es)
	}

	return data.Ok(tcs)
}

// rule:
//
//	```
//	constructor name = infix upper ident | upper ident | symbol | infix symbol ;
//	```
func maybeParseConstructorName(p Parser) (*data.Ers, data.Maybe[name]) {
	isMethod := lookahead1(p, token.MethodSymbol)
	if isMethod { // type constructor cannot have a method name
		return nil, data.Nothing[name](p)
	}

	tok := p.current()
	// infix ids are not stored w/ parens, so this will work for both infix and non-infix
	if common.Is_camelCase2(tok) {
		return nil, data.Nothing[name](p)
	}

	return nil, maybeParseName(p)
}

type typeConstructorSeq = data.NonEmpty[typeConstructor]

// rule:
//
//	```
//	type constructor = constructor name seq, {"\n"}, ":", {"\n"}, type ;
//	constructor name seq = constructor name, {{"\n"}, ",", {"\n"}, constructor name}, [{"\n"}, ","] ;
//	```
func parseTypeConstructor(p Parser, as data.Maybe[annotations]) data.Either[data.Ers, data.NonEmpty[typeConstructor]] {
	type group struct{ data.NonEmpty[name] }
	es, names, isNames := parseHandledSepSequenced[group](p, typeConstructorNameError, token.Comma, maybeParseConstructorName).Break()
	if !isNames {
		return data.PassErs[typeConstructorSeq](es)
	}

	colon, found := getKeywordAtCurrent(p, token.Colon, dropBeforeAndAfter)
	if !found {
		return data.Fail[typeConstructorSeq](ExpectedTypeJudgment, p)
	}

	es, ty, isTy := ParseType(p).Break()
	if !isTy {
		return data.PassErs[typeConstructorSeq](es)
	}

	tcs := data.MapNonEmpty(constructConstructor(as, colon, ty))(names.NonEmpty)
	return data.Ok(tcs)
}

// rule:
//
//	```
//	optional visibility = [("open" | "public"), {"\n"}] ;
//	```
func parseOptionalVisibility(p Parser) (mv data.Maybe[visibility]) {
	if vis, found := getKeywordAtCurrent(p, token.Open, dropAfter); found {
		return data.Just(visibility{vis})
	} else if vis, found = getKeywordAtCurrent(p, token.Public, dropAfter); found {
		return data.Just(visibility{vis})
	} else {
		return data.Nothing[visibility](p)
	}
}

func attachVisibility(vis data.Maybe[visibility]) func(be mainElement) data.Either[data.Ers, mainElement] {
	return func(be mainElement) data.Either[data.Ers, mainElement] {
		vbe, ok := be.(visibleBodyElement)
		visibilityExists := !vis.IsNothing()
		if !ok && visibilityExists {
			// trying to target a non-visibility modifiable body element w/ a visibility
			return data.Fail[mainElement](IllegalVisibilityTarget, be)
		} else if ok && visibilityExists {
			// attach visibility to the body element
			return data.Ok(vbe.setVisibility(vis))
		}
		// no visibility to attach, return the body element as is
		return data.Ok(be)
	}
}

func parseOneMainElement(p Parser) data.Either[data.Ers, mainElement] {
	vis := parseOptionalVisibility(p)
	return data.Cases(parseBasicBodyStructure(p, vis), data.PassErs[mainElement], attachVisibility(vis))
}

// rule:
//
//	```
//	body element = def | visible body element ;
//	```
func parseBodyElement(p Parser) data.Either[data.Ers, bodyElement] {
	es, me, isME := parseOneMainElement(p).Break()
	if !isME {
		return data.PassErs[bodyElement](es)
	}
	be := me.asBodyElement()
	return data.Ok(be)
}

func assembleDef(pat pattern) func(defBody) data.Either[data.Ers, def] {
	return func(db defBody) data.Either[data.Ers, def] {
		return data.Ok(makeDef(pat, db))
	}
}

// here for consistency b/w 'parseDef' and 'maybeParseDef' for parsing patterns
func _parseDef_ParsePattern(p Parser) (es *data.Ers, pat pattern) {
	e, pat, isPat := ParsePattern(p).Break()
	if !isPat {
		return &e, nil
	}
	return nil, pat
}

func _finishParseDef(p Parser, pat pattern) data.Either[data.Ers, def] {
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
	p.dropNewlines()
	return _finishParseDef(p, pat)
}

func maybeParseDef(p Parser) (*data.Ers, data.Maybe[def]) {
	// try to parse def: this is a bit of a doozy since patterns can be arbitrary large and appear
	// in the lhs of multiple mutually exclusive production rules--meaning we can't look ahead to
	// determine if it's a valid def. We will record the Position of the current token and then try
	origin := getOrigin(p)
	es, pat := _parseDef_ParsePattern(p)
	if es != nil {
		// reset the origin and return Nothing (no error)
		resetOrigin(p, origin)
		return nil, data.Nothing[def](p)
	} // else, keep origin and enforce non-Nothing result. This must be a def

	p.dropNewlines()
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
	// check for "impossible" keyword
	if imp, found := getKeywordAtCurrent(p, token.Impossible, dropNone); found {
		impossible := data.EOne[impossible](imp)
		body := data.EInl[defBody](impossible)
		return data.Ok(body)
	}

	var possibleLeft data.Either[data.Ers, data.Either[withClause, expr]]
	if bindingToken, found := getKeywordAtCurrent(p, bindingTokenType, dropAfter); found {
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
	whereToken, found := getKeywordAtCurrent(p, token.Where, dropBeforeAndAfter)
	if !found {
		return data.Ok(data.Nothing[whereClause](p)) // no where clause, return Nothing
	}

	es, where, isWhereBody := parseWhereBody(p).Break()
	if !isWhereBody {
		return data.PassErs[data.Maybe[whereClause]](es)
	}

	where.Position = where.Position.Update(whereToken)
	return data.Ok(data.Just(where))
}

// rule:
//
//	```
//	where body = main elem | "(", {"\n"}, main elem, {then, main elem}, {"\n"}, ")" ;
//	```
func parseWhereBody(p Parser) data.Either[data.Ers, whereClause] {
	es, wb, isWB := parseGroup[whereClause](p, ExpectedMainElement, maybeParseMainElement).Break()
	if !isWB {
		return data.PassErs[whereClause](es)
	}
	return data.Ok(wb)
}

// helper function for `parseMainElement` that transforms an `Either[Ers, a]` to an
// `Either[Ers, mainElement]` where `a` implements `mainElement`
func knownCase[a mainElement](elemRes data.Either[data.Ers, a]) data.Either[data.Ers, mainElement] {
	return data.Cases(elemRes, data.Inl[mainElement, data.Ers], fun.Compose(data.Ok, (a).asMainElement))
}

// parses a spec, alias, inst, or syntax definition when given the corresponding token type `tt`
//
// if `tt` is not One of the following this function will panic:
//   - token.Spec
//   - token.Alias
//   - token.Inst
//   - token.Syntax
func parseKnownMainElement(p Parser, tt token.Type) data.Either[data.Ers, mainElement] {
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
//	type alias = "alias", {"\n"}, name, {"\n"}, "=", {"\n"}, type ;
//	```
func parseTypeAlias(p Parser) data.Either[data.Ers, typeAlias] {
	aliasToken, found := getKeywordAtCurrent(p, token.Alias, dropAfter)
	if !found {
		return data.Fail[typeAlias](ExpectedTypeAlias, p)
	}

	n, isSomething := maybeParseName(p).Break()
	if !isSomething {
		return data.Fail[typeAlias](ExpectedTypeAliasName, aliasToken)
	}

	equalToken, found := getKeywordAtCurrent(p, token.Equal, dropBeforeAndAfter)
	if !found {
		return data.Fail[typeAlias](ExpectedAliasBinding, n)
	}

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
	syntaxToken, found := getKeywordAtCurrent(p, token.Syntax, dropAfter)
	if !found {
		return data.Fail[syntax](ExpectedSyntax, p)
	}

	es, rule, isRule := parseSyntaxRule(p).Break()
	if !isRule {
		return data.PassErs[syntax](es)
	}

	equalToken, found := getKeywordAtCurrent(p, token.Equal, dropBeforeAndAfter)
	if !found {
		return data.Fail[syntax](ExpectedSyntaxBinding, rule)
	}

	esE, e, isE := ParseExpr(p).Break()
	if !isE {
		return data.PassErs[syntax](esE)
	}

	syn := makeSyntax(rule, e)
	syn.Position = syn.Update(syntaxToken).Update(equalToken)
	return data.Ok(syn)
}

// This production rule looks a bit different than the others, but, intuitively, it just
// accepts a sequence of syntax symbols where _at least_ one of them is a raw keyword; newlines
// are permitted between each syntax symbol
//
// rule:
//
//	```
//	syntax rule = {syntax symbol, {"\n"}}, raw keyword, {{"\n"}, syntax symbol} ;
//	```
func parseSyntaxRule(p Parser) data.Either[data.Ers, syntaxRule] {
	const smallCap int = 4
	var sym syntaxSymbol
	hasSymbol := true
	hasRawKeyword := false

	ruleInsides := data.Nil[syntaxSymbol](smallCap)

	origin := getOrigin(p)
	for { // loop until there aren't any syntax symbols in view
		es, mSym := maybeParseSyntaxSymbol(p)
		if es != nil {
			return data.PassErs[syntaxRule](*es)
		}

		if sym, hasSymbol = mSym.Break(); hasSymbol {
			// check if a raw keyword has been found
			hasRawKeyword = hasRawKeyword || sym.Type().Match(syntaxRawKeyword{})

			ruleInsides = ruleInsides.Snoc(sym)

			origin = getOrigin(p) // update origin
			p.dropNewlines()
		} else {
			resetOrigin(p, origin)
			break
		}
	}

	// attempt to strengthen list -> non-empty list
	rule, just := data.EStrengthen[syntaxRule](ruleInsides).Break()
	if !just {
		return data.Fail[syntaxRule](ExpectedSyntaxRule, p)
	}

	// check if a raw keyword was found
	if !hasRawKeyword {
		return data.Fail[syntaxRule](ExpectedRawKeyword, rule)
	}
	return data.Inr[data.Ers](rule)
}

// assumes the token is the correct type (i.e., `token.RawStringValue`)
func isValidRawSyntaxSymbol(t api.Token) bool {
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
	if matchCurrentRawString(p) && isValidRawSyntaxSymbol(p.current()) {
		return nil, data.Just(parseRawSyntaxSymbolFromCurrent(p))
	} else if lookahead2(p, boundSyntaxIdentLAs...) {
		es, sym, isSym := parseBindingSyntaxIdent(p).Break()
		if !isSym {
			return &es, data.Nothing[syntaxSymbol](p)
		}
		return nil, data.Just(sym)
	} else if unit, just := parseIdent(p).Break(); just {
		id := makeStdSyntaxRuleIdent(unit)
		return nil, data.Just(data.Inl[syntaxRawKeyword](id))
	}
	return nil, data.Nothing[syntaxSymbol](p)
}

// parses a raw syntax symbol (a non-infix name inside of '`' characters w/o anything else)
func parseRawSyntaxSymbolFromCurrent(p Parser) syntaxSymbol {
	rawKey := p.current()
	p.advance()
	key := data.EOne[syntaxRawKeyword](data.EOne[rawString](rawKey))
	return data.Inr[syntaxRuleIdent](key)
}

// rule:
//
//	```
//	binding syntax ident = "{", {"\n"}, ident, {"\n"}, "}" ;
//	```
func parseBindingSyntaxIdent(p Parser) data.Either[data.Ers, syntaxSymbol] {
	lb, _ := getKeywordAtCurrent(p, token.LeftBrace, dropAfter)

	id, isSomething := parseIdent(p).Break()
	if !isSomething {
		return data.Fail[syntaxSymbol](ExpectedSyntaxBindingId, lb)
	}

	id = id.Update(lb)

	rb, found := getKeywordAtCurrent(p, token.RightBrace, dropBefore)
	if !found {
		return data.Fail[syntaxSymbol](ExpectedRightBrace, id)
	}

	id = id.Update(rb)
	sri := makeBindingSyntaxRuleIdent(id)
	symbol := data.Inl[syntaxRawKeyword](sri)
	return data.Inr[data.Ers](symbol)
}

type structureParseAttempt byte

const (
	attemptedDef structureParseAttempt = iota
	attemptedSpecDef
	attemptedSpecInst
	attemptedTypeDef
	attemptedTypeAlias
	attemptedTyping
	attemptedSyntax
	attemptedGenericStructure
)

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

func (attempted structureParseAttempt) getErrorMessageForAttempted(vis data.Maybe[visibility]) string {
	v, just := vis.Break()
	if !just {
		return "" // no visibility modifier given, no additional error can be given
	}

	switch v.Token.Type() {
	case token.Open:
		switch attempted {
		case attemptedTyping:
			return IllegalOpenModifierTyping
		case attemptedDef:
			return IllegalVisibleDef
		case attemptedTypeDef:
			return ""
		default:
			return IllegalOpenModifier
		}
	case token.Public:
		if attempted == attemptedDef {
			return IllegalVisibleDef
		}
		return ""
	}

	panic("unknown visibility modifier")
}

// Tries to parse any of the top-level body structures, returning the result and a value encoding
// the closest type of structure that an attempt was made to parse.
//
// The top-level body structures are:
//   - function definition
//   - spec definition
//   - spec instance
//   - data type definition
//   - type alias
//   - typing
//   - syntax
//
// the only one that can return a `Nothing` value is when parsing a def. All other structures will return
// a `Just` value or an error.
func optionalParseBasicStructureHelper(p Parser) (res data.Either[data.Ers, data.Maybe[mainElement]], attempted structureParseAttempt) {
	/*
	 * `spec def`, `spec inst`, `type alias`, and `syntax` can all be distinguished with lookahead of 1.
	 * `typing` or `type def` can be distinguished from `def` with lookahead of 2.
	 * `typing` and `type def` cannot be distinguished with a fixed lookahead; however, `type def`'s lhs
	 * is the `typing` rule, so we can parse `typing` and then check for `where` to determine if it is
	 * a `type def`
	 *
	 * strategy:
	 * 	1. lookahead 1 for `spec def`, `spec inst`, `type alias`, and `syntax`
	 *	2. otherwise, lookahead 2 for `typing` and `type def`, parse, then check for `where` to determine if it is a `type def`
	 *  3. lookahead 1 for `auto`
	 *	4. otherwise, if 1., 2. and 3. fail, try to parse `def`
	 */
	// 1. lookahead 1
	if tt, found := lookahead1Report(p, bodyKeywordsLAs...); found {
		result := parseKnownMainElement(p, tt)
		return data.EitherMap(data.IdErs, data.Just[mainElement])(result), attemptedWhat(tt)
	}

	// 2. lookahead 2
	if found := lookahead2(p, typingLAs...); found {
		// put as attempted typing since typing is the only One that actually matters. `attempted` is
		// used for visibility related error messages, but a type def can have any visibility modifier
		result := parseTypeDefOrTyping(p)
		return data.EitherMap(data.IdErs, data.Just[mainElement])(result), attemptedTyping
	}

	// 3. lookahead 1 for `auto`
	if found := lookahead1(p, token.Auto); found {
		result := parseAutoDefTypeSig(p)
		return data.EitherMap(data.IdErs, data.Just[mainElement])(result), attemptedTyping
	}

	// else, 3. try to parse `def`
	es, mDef := maybeParseDef(p)
	if es != nil {
		return data.PassErs[data.Maybe[mainElement]](*es), attemptedDef
	}
	return data.Ok(data.MaybeMap((def).asMainElement)(mDef)), attemptedDef
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

func strengthenStructureType(b mainElement, weaker structureParseAttempt) structureParseAttempt {
	stronger, found := typeToStructureTypeMap[b.Type()]
	if !found {
		return weaker
	}
	return stronger
}

func parseBasicBodyStructure(p Parser, vis data.Maybe[visibility]) data.Either[data.Ers, mainElement] {
	startPosition := p.current().GetPos()
	res, structureGuess := optionalParseBasicStructureHelper(p)
	es, mStructure, isRight := res.Break()
	// if res has no errors, we can get a more specific error message if needed
	var structure mainElement
	if isRight {
		structure, _ = mStructure.Break()
		if mStructure.IsNothing() {
			return data.Fail[mainElement](ExpectedMainElement, p)
		}

		structureAttempted := strengthenStructureType(structure, structureGuess)
		errorMsg := structureAttempted.getErrorMessageForAttempted(vis)
		if isRight = errorMsg == ""; !isRight {
			rng := api.WeakenRangeOver[api.Positioned](res, startPosition)
			es = es.Snoc(data.MkErr(errorMsg, rng))
		}
	}

	if !isRight {
		return data.PassErs[mainElement](es)
	}
	return data.Ok(structure)
}

// rule:
//
//	```
//	main elem = [annotations_],
//		( def
//		| spec def
//		| spec inst
//		| type def
//		| type alias
//		| typing
//		| syntax
//		) ;
//	```
func maybeParseMainElement(p Parser) (*data.Ers, data.Maybe[mainElement]) {
	es, mAnnots, isAnnots := parseAnnotations_(p).Break()
	if !isAnnots {
		return &es, data.Nothing[mainElement](p)
	}

	var mme data.Maybe[mainElement]
	res, _ := optionalParseBasicStructureHelper(p)
	sEs, mStructure, isStructure := res.Break()
	if !isStructure {
		return &sEs, data.Nothing[mainElement](p)
	} else if structure, just := mStructure.Break(); just {
		me := structure.setAnnotation(mAnnots)
		mme = data.Just(me)
	} else {
		mme = data.Nothing[mainElement](p)
	}

	return nil, mme
}

// rule:
//
//	```
//	deriving clause = "deriving", {"\n"}, deriving body ;
//	```
func parseOptionalDerivingClause(p Parser) data.Either[data.Ers, data.Maybe[deriving]] {
	derivingToken, found := getKeywordAtCurrent(p, token.Deriving, dropAfter)
	if !found {
		return data.Ok(data.Nothing[deriving](p))
	}

	es, db, isDB := parseDerivingBody(p).Break()
	if !isDB {
		return data.Inl[data.Maybe[deriving]](es)
	}

	d := db
	d.Position = d.Update(derivingToken)
	return data.Ok(data.Just(d))
}

// rule:
//
//	```
//	deriving body = constrainer | "(", {"\n"}, constrainer, {{"\n"}, ",", {"\n"}, constrainer}, [{"\n"}, ","], {"\n"}, ")" ;
//	```
func parseDerivingBody(p Parser) data.Either[data.Ers, deriving] {
	if !matchCurrentLeftParen(p) {
		es, res := nonEnclosedMaybeParseConstrainer(p)
		if es != nil {
			return data.PassErs[deriving](*es)
		} else if res.IsNothing() {
			return data.Fail[deriving](ExpectedDerivingBody, p)
		}
		bod, _ := res.Break()
		return data.Ok(data.EConstruct[deriving](bod))
	}
	es, res, ok := parseSepSequencedGroup(p, ExpectedDerivingBody, token.Comma, nonEnclosedMaybeParseConstrainer).Break()
	if !ok {
		return data.PassErs[deriving](es)
	}
	return data.Ok(deriving{res})
}

// rule:
//
//	```
//	type def or typing = typing, [{"\n"}, "where", {"\n"}, type def body, [{"\n"}, deriving clause]] ;
//	```
func parseTypeDefOrTyping(p Parser) data.Either[data.Ers, mainElement] {
	sigEs, sig, isSig := parseTypeSig(p).Break()
	if !isSig {
		return data.PassErs[mainElement](sigEs)
	}

	where, found := getKeywordAtCurrent(p, token.Where, dropBeforeAndAfter)
	if !found {
		return data.Inr[data.Ers, mainElement](sig)
	}
	es, tdb, isTbd := parseTypeDefBody(p).Break()
	if !isTbd {
		return data.PassErs[mainElement](es)
	}

	origin := getOrigin(p)
	p.dropNewlines()

	esDeriving, mDeriving, isDeriving := parseOptionalDerivingClause(p).Break()
	if !isDeriving {
		resetOrigin(p, origin)
		return data.PassErs[mainElement](esDeriving)
	} else if mDeriving.IsNothing() {
		resetOrigin(p, origin)
	}

	td := makeTypeDef(sig, tdb, mDeriving)
	td.Position = td.Update(where)
	return data.Ok[mainElement](td)
}

// rule:
//
//	```
//	with clause = "with", {"\n"}, pattern, {"\n"}, "of", {"\n"}, with clause arms ;
//	```
func parseWithClause(p Parser) data.Either[data.Ers, withClause] {
	withToken, found := getKeywordAtCurrent(p, token.With, dropAfter)
	if !found {
		return data.Fail[withClause](ExpectedWithClause, p)
	}

	es, pat, isPat := ParsePattern(p).Break()
	if !isPat {
		return data.PassErs[withClause](es)
	}

	ofToken, found := getKeywordAtCurrent(p, token.Of, dropBeforeAndAfter)
	if !found {
		return data.Fail[withClause](ExpectedOf, pat)
	}

	esArms, arms, isArms := parseWithClauseArms(p).Break()
	if !isArms {
		return data.PassErs[withClause](esArms)
	}

	with := makeWithClause(pat, arms)
	with.Position = with.Update(withToken).Update(ofToken)
	return data.Ok(with)
}

// rule:
//
//	```
//	with clause arms =
//		"(", {"\n"}, with clause arm, {then, with clause arm}, {"\n"}, ")"
//		| with clause arm ;
//	```
func parseWithClauseArms(p Parser) data.Either[data.Ers, withClauseArms] {
	return parseGroup[withClauseArms](p, ExpectedWithClauseArm, maybeParseWithClauseArm)
}

func maybeParseWithArmLhs(p Parser) (*data.Ers, data.Maybe[withArmLhs]) {
	es, mViewRefined := maybeParsePattern(p, false)
	if es != nil { // pattern found, but error while parsing it
		return es, data.Nothing[withArmLhs](p)
	}
	if mViewRefined.IsNothing() { // no pattern found, return Nothing
		return nil, data.Nothing[withArmLhs](p)
	}

	// if "|" is found, parse the intermediate pattern scrutinee
	if bar, found := getKeywordAtCurrent(p, token.Bar, dropBeforeAndAfter); found {
		pEs, scrutinee, isScrutinee := ParsePattern(p).Break()
		if !isScrutinee {
			return &pEs, data.Nothing[withArmLhs](p)
		}
		viewRefined, _ := mViewRefined.Break() // guaranteed by earlier check
		wLhs := makeWithArmLhsRefined(viewRefined, scrutinee)
		wLhs = wLhs.Update(bar)
		return nil, data.Just(wLhs)
	}
	// first pattern was not the view refined pattern, so it must be intermediate pattern scrutinee
	scrutinee, _ := mViewRefined.Break() // guaranteed by earlier check
	return nil, data.Just(makeWithArmLhs(scrutinee))
}

// rule:
//
//	```
//	with clause arm = [view refined pattern, {"\n"}], pattern, {"\n"}, def body thick arrow ;
//	view refined pattern = pattern, {"\n"}, "|" ;
//	```
func maybeParseWithClauseArm(p Parser) (*data.Ers, data.Maybe[withClauseArm]) {
	es, mLhs := maybeParseWithArmLhs(p)
	if es != nil {
		return es, data.Nothing[withClauseArm](p)
	}

	lhs, justLhs := mLhs.Break()
	if !justLhs {
		return nil, data.Nothing[withClauseArm](p)
	}

	p.dropNewlines()
	esDb, db, isDb := parsePatternBoundBody(p, token.ThickArrow).Break()
	if !isDb {
		return &esDb, data.Nothing[withClauseArm](p)
	}

	wc := makeWithClauseArm(lhs, db)
	return nil, data.Just(wc)
}
