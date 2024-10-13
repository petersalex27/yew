package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/common/data"
)

// rule:
//
//	```
//	deriving clause = "deriving", {"\n"}, deriving body ;
//	```
func parseOptionalDerivingClause(p Parser) data.Either[data.Ers, data.Maybe[deriving]] {
	derivingToken, found := getKeywordAtCurrent(p, token.Deriving)
	if !found {
		return data.Ok(data.Nothing[deriving](p))
	}

	es, db, isDB := parseDerivingBody(p).Break()
	if !isDB {
		return data.Inl[data.Maybe[deriving]](es)
	}

	d := data.EOne[deriving](db)
	d.Position = d.Update(derivingToken)
	return data.Ok(data.Just[deriving](d))
}

// rule:
//
//	```
//	deriving body = constrainer | "(", {"\n"}, constrainer, {{"\n"}, ",", {"\n"}, constrainer}, [{"\n"}, ","], {"\n"}, ")" ;
//	```
func parseDerivingBody(p Parser) data.Either[data.Ers, derivingBody] {
	return parseSepSequenced[derivingBody](p, ExpectedDerivingBody, token.Comma, nonEnclosedMaybeParseConstrainer)
}

func nonEnclosedMaybeParseConstrainer(p Parser) (*data.Ers, data.Maybe[constrainer]) {
	return maybeParseConstrainer(p, false)
}

// rule:
//
//	```
//	constrainer = upper ident, pattern | "(", {"\n"}, enc constrainer {"\n"}, ")" ;
//	enc constrainer = upper ident, {"\n"}, pattern ;
//	```
func parseConstrainer(p Parser, enclosed bool) data.Either[data.Ers, constrainer] {
	es, mC := maybeParseConstrainer(p, enclosed)
	if es != nil {
		return data.Inl[constrainer](*es)
	}

	if c, just := mC.Break(); just {
		return data.Ok(c)
	}
	return data.Fail[constrainer](ExpectedConstrainer, p)
}

// rule:
//
//	```
//	constrainer = upper ident, pattern ;
//	```
func maybeParseConstrainer(p Parser, enclosed bool) (*data.Ers, data.Maybe[constrainer]) {
	currentIsLeftParen := matchCurrentLeftParen(p)
	if !enclosed && currentIsLeftParen {
		origin := getOrigin(p)
		lparen := p.current()
		p.advance()
		p.dropNewlines()
		if !currentIsUpperIdent(p) {
			p = resetOrigin(p, origin)
			return nil, data.Nothing[constrainer](p)
		}

		es, c, isC := parseConstrainer(p, true).Break()
		if !isC {
			return &es, data.Nothing[constrainer](p)
		}

		rparen, found := getKeywordAtCurrent(p, token.RightParen)
		if !found {
			e := data.Nil[data.Err](1).Snoc(mkErr(ExpectedRightParen, p))
			return &e, data.Nothing[constrainer](p)
		}

		c.Position = c.Update(lparen)
		c.Position = c.Update(rparen)
		return nil, data.Just(c)
	} else if enclosed && currentIsLeftParen {
		// if enclosed already, this is an error--extraneous parens
		e := data.Nil[data.Err](1).Snoc(mkErr(IllegalMultipleEnclosure, p))
		return &e, data.Nothing[constrainer](p)
	}

	// get upper ident
	upper, isSomething := createUpperIdent(p.current())(p).Break()
	if !isSomething {
		return nil, data.Nothing[constrainer](p)
	}

	// get pattern, dropping newlines if enclosed
	if enclosed {
		p.dropNewlines()
	}
	es, pat, isPat := ParsePattern(p).Break()
	if !isPat {
		return &es, data.Nothing[constrainer](p) // return error
	}

	return nil, data.Just(data.EMakePair[constrainer](upper, pat))
}

func constraintElemLA(p Parser) bool {
	return lookahead2(p, [2]token.Type{token.Id, token.Comma})
}

// this is not an actual rule in the grammar, data.Just a helper function for `parseConstraintElem`;
// though, it can be represented as the following rule:
//
//	```
//	upper id sequence = {upper ident, {"\n"}, ",", {"\n"}} ;
//	```
func parseUpperIdSequence(p Parser) data.List[upperIdent] {
	upperIds := data.Nil[upperIdent]()

	for currentIsUpperIdent(p) && constraintElemLA(p) {
		upper := p.current()
		p.advance()
		p.dropNewlines()

		comma, found := getKeywordAtCurrent(p, token.Comma)
		if !found {
			panic("verification was incorrect")
		}

		upperIds = upperIds.Snoc(data.EOne[upperIdent](upper))
		upperIds.Position = upperIds.Update(comma)
	}

	return upperIds
}

type constraintElem = data.Pair[data.List[upperIdent], constrainer]

// rule:
//
//	```
//	constraint elem = {upper ident, {"\n"}, ",", {"\n"}}, constrainer ;
//	enc constraint elem = {upper ident, {"\n"}, ",", {"\n"}}, enc constrainer ;
//	```
func parseConstraintElem(p Parser, enclosed bool) data.Either[data.Ers, constraintElem] {
	upperIds := parseUpperIdSequence(p)
	es, c, isC := parseConstrainer(p, enclosed).Break()
	if !isC {
		return data.Inl[constraintElem](es)
	}
	return data.Ok(data.MakePair(upperIds, c))
}

func maybeParseConstraintElem(p Parser) (*data.Ers, data.Maybe[constraintElem]) {
	// this will return an empty data.List if no upper idents followed by ',' are found
	//
	// so, this is safe to call b/c on a non-zero length the next call is
	// required to be a constrainer. If it's not a zero-length data.List, the
	// next call is optional.
	upperIds := parseUpperIdSequence(p)
	if upperIds.Len() != 0 {
		es, c, isC := parseConstrainer(p, false).Break()
		if !isC { // requirement not met
			return &es, data.Nothing[constraintElem](p)
		}

		return nil, data.Just(data.MakePair(upperIds, c))
	}

	// else, constrainer is optional--if 'data.Nothing', then return data.Nothing
	// 	- NOTE: upperIds is empty

	es, mC := maybeParseConstrainer(p, false)
	if es != nil {
		return es, data.Nothing[constraintElem](p)
	}

	if c, just := mC.Break(); just {
		return nil, data.Just(data.MakePair(upperIds, c))
	}
	return nil, data.Nothing[constraintElem](p)
}

// rule:
//
//	```
//	constraint = "(", {"\n"}, constraint group, {"\n"}, ")" | constraint elem ;
//		constraint group = enc constraint elem, {{"\n"}, ",", {"\n"}, enc constraint elem}, [{"\n"}, ",", {"\n"}] ;
//		constraint elem = {upper ident, {"\n"}, ",", {"\n"}}, constrainer ;
//		enc constraint elem = {upper ident, {"\n"}, ",", {"\n"}}, enc constrainer ;
//	```
//
// inlining rules:
//
//	```
//	constraint =
//		constraint elem
//		| "(", {"\n"}, constraint elem, {{"\n"}, ",", {"\n"}, constraint elem}, [{"\n"}, ",", {"\n"}], {"\n"}, ")" ;
//	```
func parseConstraint(p Parser) data.Either[data.Ers, constraintVerified] {
	if matchCurrentLeftParen(p) {
		lparen, _ := getKeywordAtCurrent(p, token.LeftParen)

		res := parseSepSequenced[constraintVerified](p, ExpectedConstraint, token.Comma, maybeParseConstraintElem)
		if res.IsLeft() {
			return res
		}
		res = res.Update(lparen) // even if it's an error, update

		rparen, found := getKeywordAtCurrent(p, token.RightParen)
		if !found {
			return data.Fail[constraintVerified](ExpectedRightParen, lparen)
		}
		res = res.Update(rparen)

		return res
	}

	es, c, isC := parseConstraintElem(p, false).Break()
	if !isC {
		return data.Inl[constraintVerified](es)
	}
	return data.Ok(data.EConstruct[constraintVerified](c))
}

// rule:
//
//	```
//	constrainer = upper ident, pattern ;
//	enc constrainer = upper ident, {"\n"}, pattern ;
//	spec head = [constraint, {"\n"}, "=>", {"\n"}], constrainer ;
//	```
func parseSpecHead(p Parser) data.Either[data.Ers, specHead] {
	sh := specHead{}
	// strategy: parse as constraint, this can data.Either return data.Just a (case A.) constrainer or (case B.)
	// a full constraint (w/o  the "=>")
	//
	// case A.
	//	- this is the case when `len(c.rest) == 0 && len(c.first.first.elements) == 0`
	//	- if a "=>" does NOT follow, this is the non-optional `constrainer` in `c.first.second`
	//	- else, fallthrough to case B.
	// case B.
	//	- this is the case otherwise (it must have more than a single constraint b/c it didn't data.Fail
	//	  and isn't case A.).
	es, c, isC := parseConstraint(p).Break()
	if !isC {
		return data.Inl[specHead](es)
	}

	var thickArrow api.Token
	// this is okay, it will be used here, by `spec def` before `spec dependency`, or by spec def
	// before `where`
	p.dropNewlines()

	if c.Tail().Len() == 0 && c.Head().Fst().Len() == 0 { // case A.
		var found bool
		if thickArrow, found = getKeywordAtCurrent(p, token.Arrow); !found {
			sh = data.EMakePair[specHead](data.Nothing[constraint](p), c.Head().Snd())
			return data.Ok(sh)
		} // else found, fall out of branch
	} // else, case B.

	// parse a constrainer for the rhs of "=>"
	esRhs, rhs, isRhs := parseConstrainer(p, false).Break()
	if !isRhs {
		return data.Inl[specHead](esRhs)
	}

	// assemble: c => rhs
	sh = data.EMakePair[specHead](data.Just[constraint](c), rhs)
	sh.Position = sh.Update(thickArrow)
	return data.Ok(sh)
}

type requiringClause = data.NonEmpty[def]

// rule:
//
//	```
//	spec def = "spec", {"\n"}, spec head, [{"\n"}, spec dependency], {"\n"}, "where", {"\n"}, spec body, [{"\n"}, requiring clause] ;
//	```
func parseSpecDef(p Parser) data.Either[data.Ers, specDef] {
	var sd specDef
	if specToken, found := getKeywordAtCurrent(p, token.Spec); !found {
		return data.Fail[specDef](ExpectedSpecDef, p)
	} else { // use of `else` allows use of 'found' later, keeping it out of scope
		sd.Position = sd.Update(specToken)
	}

	es, sh, isSH := parseSpecHead(p).Break()
	if !isSH {
		return data.Inl[specDef](es)
	}
	sd.Position = sd.Update(sh)
	sd.specHead = sh

	p.dropNewlines()
	esDep, dep, isDep := parseOptionalSpecDependency(p).Break()
	if !isDep {
		return data.Inl[specDef](esDep)
	}
	sd.Position = sd.Update(dep)
	sd.dependency = dep

	// redundant if dep.IsNothing(), but dropNewlines is idempotent (when 'p' is used only in a single
	// thread)--so this keeps the code simpler
	p.dropNewlines()
	if whereToken, found := getKeywordAtCurrent(p, token.Where); !found {
		return data.Fail[specDef](ExpectedSpecWhere, dep)
	} else {
		sd.Position = sd.Update(whereToken)
	}

	esBody, body, isBody := parseSpecBody(p).Break()
	if !isBody {
		return data.Inl[specDef](esBody)
	}
	sd.Position = sd.Update(body)
	sd.specBody = body

	p.dropNewlines()
	esReq, req, isReq := parseOptionalRequiringClause(p).Break()
	if !isReq {
		return data.Inl[specDef](esReq)
	}
	sd.Position = sd.Update(req)
	sd.requiring = req

	return data.Ok(sd)
}

// rule:
//
//	```
//	spec dependency = "from", {"\n"}, pattern ;
//	```
func parseOptionalSpecDependency(p Parser) data.Either[data.Ers, data.Maybe[pattern]] {
	fromToken, found := getKeywordAtCurrent(p, token.From)
	if !found {
		return data.Ok(data.Nothing[pattern](p))
	}

	es, pat, isPat := ParsePattern(p).Break()
	if !isPat {
		return data.Inl[data.Maybe[pattern]](es)
	}

	return data.Ok(data.Just[pattern](pat.updatePosPattern(fromToken)))
}

// rule:
//
//	```
//	spec body =
//		spec member
//		| "(", {"\n"}, spec member, {{"\n"}, spec member}, {"\n"}, ")" ;
//	```
func parseSpecBody(p Parser) data.Either[data.Ers, specBody] {
	return parseGroup[specBody](p, ExpectedTypingOrDef, parseMaybeSpecMember)
}

// rule:
//
//	```
//	spec member = [annotations_], def | [annotations_], typing ;
//	```
//
// NOTE: typing and def can be, as elsewhere, higher-order terms. It's worth noting this since
// the style and semantics are similar to type classes in Haskell but Haskell does not, in its
// standard form, have the same support for higher-order types.
func parseMaybeSpecMember(p Parser) (*data.Ers, data.Maybe[specMember]) {
	es, mAnnots, isMAnnots := parseAnnotations(p).Break()
	if !isMAnnots {
		return &es, data.Nothing[specMember](p)
	}

	// check for typing
	isTyping := lookahead2(p, [2]token.Type{token.Id, token.Colon}, [2]token.Type{token.Infix, token.Colon})
	if isTyping {
		es, t, isT := parseTypeSig(p).Break()
		if !isT {
			return &es, data.Nothing[specMember](p)
		}
		t = t.setAnnotation(mAnnots).(typing)
		return nil, data.Just(data.Inr[def](t))
	}

	// try to parse a def: see notes in `maybeParseDef` for a description of the strategy
	pEs, mDef := maybeParseDef(p)
	if pEs != nil {
		return pEs, data.Nothing[specMember](p)
	} else if unit, just := mDef.Break(); !just {
		return nil, data.Nothing[specMember](p)
	} else {
		d := unit.setAnnotation(mAnnots).(def)
		return nil, data.Just(data.Inl[typing](d))
	}
}

// Note that this does something quite subtle: if an annotation block is found, it must be followed
// by a def. This means this function cannot be used more generally where an annotation can annotate
// other nodes.
//
// In short, only use this function when `def` is the only valid annotation target
func parseMaybeAnnotatedDef(p Parser) (*data.Ers, data.Maybe[def]) {
	es, mAnnots, isMAnnots := parseAnnotations(p).Break()
	if !isMAnnots {
		return &es, data.Nothing[def](p)
	}

	// if an annotation block is found, it must be followed by a def (annotations must have targets
	// and this function specifies that the target is a def)
	if mAnnots.IsNothing() { // no annotations found, def is optional
		return maybeParseDef(p)
	}

	es2, d, isDef := parseDef(p).Break()
	if !isDef {
		return &es2, data.Nothing[def](p)
	}
	d = d.setAnnotation(mAnnots).(def)
	return nil, data.Just(d)
}

// rule:
//
//	```
//	requiring clause = "requiring", {"\n"},
//		( [annotations_], def
//		| "(", {"\n"}, [annotations_], def, {{"\n"}, [annotations_], def}, {"\n"}, ")"
//		) ;
//	```
func parseOptionalRequiringClause(p Parser) data.Either[data.Ers, data.Maybe[requiringClause]] {
	req, found := getKeywordAtCurrent(p, token.Requiring)
	if !found {
		return data.Ok(data.Nothing[requiringClause](p)) // no requiring clause, return data.Nothing
	}

	es, reqBody, isReqBody := parseGroup[struct{ data.NonEmpty[def] }](p, ExpectedDef, parseMaybeAnnotatedDef).Break()
	if !isReqBody {
		return data.Inl[data.Maybe[requiringClause]](es)
	}
	reqBody.Position = reqBody.Update(req)
	return data.Ok(data.Just(reqBody.NonEmpty))
}

// rule:
//
//	```
//	spec inst = "inst", {"\n"}, spec head, [{"\n"}, spec inst target], {"\n"}, spec inst where clause ;
//	spec inst where clause = "where", {"\n"}, spec inst member group ;
//	spec inst member group = spec member | "(", {"\n"}, spec member, {{"\n"}, spec member}, {"\n"}, ")" ;
//	```
func parseSpecInst(p Parser) data.Either[data.Ers, specInst] {
	var si specInst

	inst, found := getKeywordAtCurrent(p, token.Inst)
	if !found {
		return data.Fail[specInst](ExpectedSpecInst, p)
	}
	si.Position = si.Update(inst)

	es, sh, isSH := parseSpecHead(p).Break()
	if !isSH {
		return data.Inl[specInst](es)
	}
	si.Position = si.Update(sh)
	si.head = sh

	p.dropNewlines()
	esTarget, target, isTarget := parseOptionalSpecInstTarget(p).Break()
	if !isTarget {
		return data.Inl[specInst](esTarget)
	}
	si.Position = si.Update(target)
	si.target = target

	where, foundWhere := getKeywordAtCurrent(p, token.Where)
	if !foundWhere {
		return data.Fail[specInst](ExpectedInstWhere, target)
	}
	si.Position = si.Update(where)

	esBody, specBod, isBody := parseSpecBody(p).Break()
	if !isBody {
		return data.Inl[specInst](esBody)
	}
	si.Position = si.Update(specBod)
	si.body = specBod

	return data.Ok(si)
}

// rule:
//
//	```
//	spec inst target = "=", {"\n"}, constrainer ;
//	```
func parseOptionalSpecInstTarget(p Parser) data.Either[data.Ers, data.Maybe[constrainer]] {
	equal, found := getKeywordAtCurrent(p, token.Equal)
	if !found {
		return data.Ok(data.Nothing[constrainer](p))
	}

	es, c, isC := parseConstrainer(p, false).Break()
	if !isC {
		return data.Inl[data.Maybe[constrainer]](es)
	}
	c.Position = c.Update(equal)
	return data.Ok(data.Just(c))
}
