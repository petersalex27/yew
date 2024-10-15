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
//	body = {{"\n"}, [annotations_], body elem} ;
//	```
func parseBody(p Parser) (theBody data.Either[data.Ers, data.Maybe[body]], mFooterAnnots data.Maybe[annotations]) {
	const smallBodyCap int = 16
	sourceBody := body{data.Nil[bodyElement](smallBodyCap)}
	var es data.Ers
	var isMAnnots bool

	has2ndTerm := false

	for {
		p.dropNewlines()
		es, mFooterAnnots, isMAnnots = parseAnnotations(p).Break()
		if !isMAnnots { // not just annotations & not nothing -> void
			return data.PassErs[data.Maybe[body]](es), mFooterAnnots
		} else if isMAnnots && lookahead1(p, token.EndOfTokens) {
			// possibly parsed footer annotations, return body and "Maybe" footer annotations
			theBody = data.Ok(data.Just(sourceBody))
			break
		}

		esBe, be, isBe := parseBodyElement(p, mFooterAnnots).Break()
		if !isBe {
			return data.PassErs[data.Maybe[body]](esBe), mFooterAnnots
		}

		sourceBody.List = sourceBody.Snoc(be)
		has2ndTerm = true
	}

	if !has2ndTerm { // no body elements parsed, return Nothing instead of empty list
		theBody = data.Ok(data.Nothing[body](sourceBody))
	}
	return theBody, mFooterAnnots
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
//	type def body =
//		"impossible"
//		| "(",{"\n"},[annotations_],type cons,{{"\n"},[annotations_],type cons},{"\n"},")"
//		| "(", {"\n"}, "impossible", {"\n"}, ")"
//		| [annotations_], type cons ;
//	type cons = typing ;
//	```
func parseTypeDefBody(p Parser) data.Either[data.Ers, typeDefBody] {
	if matchCurrentImpossible(p) {
		// constant "impossible" case
		td := data.Inr[data.NonEmpty[typeConstructor]](impossible{data.One(p.current())})
		p.advance()

		return data.Ok(td)
	}
	// else if !matchCurrentLeftParen(p) {
	// 	res := parseTypeDefBodyTypeCons(p)
	// }

	lparen, found := getKeywordAtCurrent(p, token.LeftParen, dropAfter)
	var bod data.NonEmpty[typeConstructor]
	es, tcs, isTC := parseTypeDefBodyTypeCons(p).Break()
	if !isTC {
		return data.PassErs[typeDefBody](es)
	}
	if !found {
		return data.Ok(data.Inl[impossible](tcs))
	} else {
		bod = tcs
	}

	p.dropNewlines()
	for !token.RightParen.Match(p.current()) {
		es, tcs, isTC := parseTypeDefBodyTypeCons(p).Break()
		if !isTC {
			return data.PassErs[typeDefBody](es)
		}
		bod = bod.Append(tcs.Elements()...)
		p.dropNewlines()
	}

	rparen, found := getKeywordAtCurrent(p, token.RightParen, dropNone) // newlines already dropped
	if !found {
		return data.Fail[typeDefBody](ExpectedRightParen, lparen)
	}

	bod.Position = bod.Update(lparen).Update(rparen)
	return data.Ok(data.Inl[impossible](bod))
}

func parseTypeDefBodyTypeCons(p Parser) data.Either[data.Ers, data.NonEmpty[typeConstructor]] {
	return runCases(p, parseAnnotations, passParseErs[data.NonEmpty[typeConstructor]], parseTypeConstructor)
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
//	maybe visibility = [("open" | "public"), {"\n"}] ;
//	```
func parseOptionalVisibility(p Parser) (mv data.Maybe[visibility]) {
	if lookahead1(p, visibilityLAs...) {
		visibilityToken := p.current()
		mv = data.Just(data.EOne[visibility](visibilityToken))
		p.advance()
		p.dropNewlines()
	} else {
		mv = data.Nothing[visibility](p)
	}
	return mv
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
func parseBodyElement(p Parser, as data.Maybe[annotations]) data.Either[data.Ers, bodyElement] {
	es, me, isME := parseOneMainElement(p).Break()
	if !isME {
		return data.PassErs[bodyElement](es)
	}
	be := me.setAnnotation(as).asBodyElement()
	return data.Ok(be)
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
	return nil, pat
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
		// reset the origin and return Nothing (no error)
		p = resetOrigin(p, origin)
		return nil, data.Nothing[def](p)
	} // else, keep origin and enforce non-Nothing result. This must be a def

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
	if imp, found := getKeywordAtCurrent(p, token.Impossible, dropNone); found {
		return data.Ok(data.EInl[defBody](data.EOne[impossible](imp)))
	}

	var possibleLeft data.Either[data.Ers, data.Either[withClause, expr]]
	if bindingToken, found := getKeywordAtCurrent(p, bindingTokenType, dropBeforeAndAfter); found {
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
//	where body = main elem | "(", {"\n"}, main elem, {{"\n"}, main elem}, {"\n"}, ")" ;
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
//	"alias", {"\n"}, name, {"\n"}, "=", {"\n"}, type ;
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

	for hasSymbol { // loop until there aren't any syntax symbols in view
		es, mSym := maybeParseSyntaxSymbol(p)
		if es != nil {
			return data.PassErs[syntaxRule](*es)
		}

		if sym, hasSymbol = mSym.Break(); hasSymbol {
			// check if a raw keyword has been found
			hasRawKeyword = hasRawKeyword || sym.Type().Match(syntaxRawKeyword{})

			ruleInsides = ruleInsides.Snoc(sym)
			p.dropNewlines()
		}
	}

	// attempt to strengthen list -> non-empty list
	rule, just := ruleInsides.Strengthen().Break()
	if !just {
		return data.Fail[syntaxRule](ExpectedSyntaxRule, p)
	}

	// check if a raw keyword was found
	if !hasRawKeyword {
		return data.Fail[syntaxRule](ExpectedRawKeyword, rule)
	}
	return data.Inr[data.Ers](syntaxRule{rule})
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

	p.dropNewlines()
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

	switch v.Type() {
	case token.Open:
		switch attempted {
		case attemptedSpecInst,
			attemptedSpecDef,
			attemptedTypeAlias,
			attemptedSyntax:
			return IllegalOpenModifier
		case attemptedTypeDef:
			return IllegalOpenModifierTyping
		case attemptedTyping:
			return IllegalOpenModifierTyping
		}
	case token.Public:
		if attempted == attemptedDef {
			return IllegalVisibleDef
		}
	}
	return UnexpectedStructure
}

var okJustMainElement = fun.Compose(data.Ok, data.Just[mainElement])

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
	 *	2. otherwise, lookahead 2 for `typing` and `type def`, then check for `where` to determine if it is a `type def`
	 *	3. otherwise, if 1. and 2. data.Fail, try to parse `def`
	 */
// 1. lookahead 1
	if tt, found := lookahead1Report(p, bodyKeywordsLAs...); found {
		res := parseKnownMainElement(p, tt)
		return data.Cases(res, data.PassErs[data.Maybe[mainElement]], okJustMainElement), attemptedWhat(tt)
	}

	// 2. lookahead 2
	if found := lookahead2(p, typingLAs...); found {
		// put as attempted typing since typing is the only One that actually matters. `attempted` is
		// used for visibility related error messages, but a type def can have any visibility modifier
		res := runCases(p, parseTypeSig, passParseErs[mainElement], parseTypeDefOrTyping)
		return data.Cases(res, data.PassErs[data.Maybe[mainElement]], okJustMainElement), attemptedTyping
	}

	// else, 3. try to parse `def`
	es, mDef := maybeParseDef(p)
	if es != nil {
		return data.PassErs[data.Maybe[mainElement]](*es), attemptedDef
	}
	f := fun.Compose(data.Just[mainElement], (def).asMainElement)
	return data.Ok(data.Bind(mDef, f)), attemptedDef
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
	es, mAnnots, isAnnots := parseAnnotations(p).Break()
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

func parseTypeDefOrTyping(p Parser, t typing) data.Either[data.Ers, mainElement] {
	where, found := getKeywordAtCurrent(p, token.Where, dropAfter)
	if !found {
		return data.Inr[data.Ers, mainElement](t)
	}
	es, tdb, isTbd := parseTypeDefBody(p).Break()
	if !isTbd {
		return data.PassErs[mainElement](es)
	}

	esDeriving, mDeriving, isDeriving := parseOptionalDerivingClause(p).Break()
	if !isDeriving {
		return data.PassErs[mainElement](esDeriving)
	}

	td := makeTypeDef(t, tdb, mDeriving)
	td.Position = td.Update(where)
	return data.Ok[mainElement](td)
}

// rule:
//
//	```
//	with clause = "with", {"\n"}, pattern, {"\n"}, "of", {"\n"}, with clause arms ;
//	```
func parseWithClause(p Parser) data.Either[data.Ers, withClause] {
	with := withClause{}

	withToken, found := getKeywordAtCurrent(p, token.With, dropAfter)
	if !found {
		return data.Fail[withClause](ExpectedWithClause, p)
	}
	with.Position = with.Update(withToken)

	es, pat, isPat := ParsePattern(p).Break()
	if !isPat {
		return data.PassErs[withClause](es)
	}
	with.Position = with.Update(pat)

	ofToken, found := getKeywordAtCurrent(p, token.Of, dropBeforeAndAfter)
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
	} else if pat.IsNothing() { // no pattern found, return Nothing
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
