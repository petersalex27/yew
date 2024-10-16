package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/api/util/fun"
	"github.com/petersalex27/yew/common/data"
)

func nonEnclosedMaybeParseConstrainer(p Parser) (*data.Ers, data.Maybe[constrainer]) {
	return maybeParseConstrainer(p, false)
}

// rule:
//
//	```
//	constrainer = upper ident, pattern | "(", {"\n"}, enc constrainer {"\n"}, ")" ;
//	enc constrainer = upper ident, {"\n"}, pattern ;
//	```
func parseConstrainer(p Parser) data.Either[data.Ers, constrainer] {
	es, mC := maybeParseConstrainer(p, false)
	if es != nil {
		return data.Inl[constrainer](*es)
	}

	if c, just := mC.Break(); just {
		return data.Ok(c)
	}
	return data.Fail[constrainer](ExpectedConstrainer, p)
}

// helper function, not a rule--encodes a one-time-enclosed constrainer
func parseMaybeOneTimeEnclosedConstrainer(p Parser) (*data.Ers, data.Maybe[constrainer]) {
	f := fun.BinBind1st_PairTarget(maybeParseConstrainer, true)
	es, pos, res := maybeParseParenEnclosed(p, f)
	if es != nil {
		return es, data.Nothing[constrainer](p)
	} else if c, just := res.Break(); just {
		c.Position = c.Update(pos)
		res = data.Just(c) // re-lift
	}
	return nil, res
}

// rule:
//
//	```
//	constrainer = upper ident, pattern | "(", {"\n"}, upper ident, {"\n"}, pattern {"\n"}, ")" ;
//	```
func maybeParseConstrainer(p Parser, enclosed bool) (*data.Ers, data.Maybe[constrainer]) {
	if isLP := matchCurrentLeftParen(p); !enclosed && isLP {
		return parseMaybeOneTimeEnclosedConstrainer(p)
	} else if enclosed && isLP { // if enclosed already, this is an error--extraneous parens
		e := data.MkErr(IllegalMultipleEnclosure, p)
		es := data.Nil[data.Err](1).Snoc(e)
		return &es, data.Nothing[constrainer](p)
	}

	// get upper ident
	upper, just := parseUpperIdent(p).Break()
	if !just {
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
	return currentIsUpperIdent(p) && lookahead2(p, unverifiedConstraintLAs...)
}

// this is not an actual rule in the grammar, data.Just a helper function for `parseConstraintElem`;
// though, it can be represented as the following rule:
//
//	```
//	upper id sequence = {upper ident, {"\n"}, ",", {"\n"}} ;
//	```
func parseUpperIdSequence(p Parser) data.List[upperIdent] {
	upperIds := data.Nil[upperIdent]()

	for constraintElemLA(p) {
		upper := p.current()
		p.advance()

		comma, found := getKeywordAtCurrent(p, token.Comma, dropBeforeAndAfter)
		if !found {
			panic("verification was incorrect")
		}

		upperIds = upperIds.Snoc(data.EOne[upperIdent](upper))
		upperIds.Position = upperIds.Update(comma)
	}

	return upperIds
}

// rule:
//
//	```
//	constraint elem = {upper ident, {"\n"}, ",", {"\n"}}, enc constrainer ;
//	```
func maybeParseConstraintElem(p Parser) (*data.Ers, data.Maybe[constraintElem]) {
	// this will return an empty list if no upper idents followed by ',' are found
	upperIds := parseUpperIdSequence(p)
	es, mC := maybeParseConstrainer(p, true)
	if es != nil {
		return es, data.Nothing[constraintElem](p)
	} else if c, just := mC.Break(); !just && upperIds.Len() != 0 {
		e := data.MkErr(ExpectedConstrainer, p)
		es := data.Nil[data.Err](1).Snoc(e)
		return &es, data.Nothing[constraintElem](p)
	} else if just {
		return nil, data.Just(data.MakePair(upperIds, c)) // w/ possibly empty upperIds list
	}
	return nil, data.Nothing[constraintElem](p) // nothing found (no upper ident seq and no constrainer)
}

// rule:
//
//	```
//	constraint group = constraint elem, {{"\n"}, ",", {"\n"}, constraint elem}, [{"\n"}, ",", {"\n"}] ;
//	```
func maybeParseConstraintGroup(p Parser) (*data.Ers, data.Maybe[data.NonEmpty[constraintElem]]) {
	if !lookahead1(p, token.LeftParen) {
		return nil, data.Nothing[data.NonEmpty[constraintElem]](p) // no constraint group, return data.Nothing
	}

	// require at least one constraint elem
	res := parseGroup[constraintVerified](p, ExpectedConstraintElem, maybeParseConstraintElem)

	es, c, isC := res.Break()
	if !isC {
		return &es, data.Nothing[data.NonEmpty[constraintElem]](p)
	}

	return nil, data.Just(c.NonEmpty)
}

// rule:
//
//	```
//	constraint = "(", {"\n"}, constraint group, {"\n"}, ")" | constrainer ;
//	```
func parseConstraint(p Parser) data.Either[data.Ers, constraintVerified] {
	// if constraint group, parens are handled in call--note that this doesn't follow the rule exactly ...
	es, mCG := maybeParseConstraintGroup(p)
	if es != nil {
		return data.Inl[constraintVerified](*es)
	} else if cg, just := mCG.Break(); just {
		return data.Ok(constraintVerified{cg})
	}

	esC, c := maybeParseConstrainer(p, false)
	if esC != nil {
		return data.Inl[constraintVerified](*esC)
	} else if c, just := c.Break(); just {
		ce := data.MakePair(data.Nil[upperIdent](), c) // constraint elem
		cg := data.Singleton(ce)                       // constraint group
		return data.Ok(constraintVerified{cg})
	} else {
		return data.Fail[constraintVerified](ExpectedConstraint, p)
	}
}

// rule:
//
//	```
//	spec head = [constraint, {"\n"}, "=>", {"\n"}], constrainer ;
//	```
func parseSpecHead(p Parser) data.Either[data.Ers, specHead] {
	sh := specHead{}
	// strategy: parse as constraint, this can data.Either return data.Just a (case A.) constrainer or (case B.)
	// a full constraint (w/o  the "=>")
	//
	// case A.
	//	- this is the case when `c.Tail().Len() == 0 && c.Head().Fst().Len() == 0`
	//	- if a "=>" does NOT follow, this is the non-optional `constrainer` in `c.first.second`
	//	- else, fallthrough to case B.
	// case B.
	//	- this is the case otherwise (it must have more than a single constraint b/c it didn't fail
	//	  and isn't case A.).
	es, c, isC := parseConstraint(p).Break()
	if !isC {
		return data.Inl[specHead](es)
	}

	var thickArrow api.Token
	isCaseA := c.Tail().Len() == 0 && c.Head().Fst().Len() == 0
	if isCaseA {
		ta, found := getKeywordAtCurrent(p, token.ThickArrow, dropBeforeAndAfter)
		if !found {
			sh = data.EMakePair[specHead](data.Nothing[constraintVerified](p), c.Head().Snd())
			return data.Ok(sh)
		} // else found, fall out of branch
		thickArrow = ta
	} // else, case B.

	// parse a constrainer for the rhs of "=>"
	esRhs, rhs, isRhs := parseConstrainer(p).Break()
	if !isRhs {
		return data.Inl[specHead](esRhs)
	}

	// assemble: c => rhs
	sh = data.EMakePair[specHead](data.Just(c), rhs)
	sh.Position = sh.Update(thickArrow)
	return data.Ok(sh)
}

// rule:
//
//	```
//	spec def = "spec", {"\n"}, spec head, [{"\n"}, spec dependency], {"\n"}, "where", {"\n"}, spec body, [{"\n"}, requiring clause] ;
//	```
func parseSpecDef(p Parser) data.Either[data.Ers, specDef] {
	//var sd specDef
	position := api.ZeroPosition()
	if specToken, found := getKeywordAtCurrent(p, token.Spec, dropAfter); !found {
		return data.Fail[specDef](ExpectedSpecDef, p)
	} else { // use of `else` allows use of 'found' later, keeping it out of scope
		position = position.Update(specToken)
	}

	es, sh, isSH := parseSpecHead(p).Break()
	if !isSH {
		return data.Inl[specDef](es)
	}

	p.dropNewlines()
	esDep, dep, isDep := parseOptionalSpecDependency(p).Break()
	if !isDep {
		return data.Inl[specDef](esDep)
	}

	if whereToken, found := getKeywordAtCurrent(p, token.Where, dropBeforeAndAfter); !found {
		return data.Fail[specDef](ExpectedSpecWhere, p)
	} else {
		position = position.Update(whereToken)
	}

	esBody, body, isBody := parseSpecBody(p).Break()
	if !isBody {
		return data.Inl[specDef](esBody)
	}

	p.dropNewlines()
	esReq, req, isReq := parseOptionalRequiringClause(p).Break()
	if !isReq {
		return data.Inl[specDef](esReq)
	}

	sd := makeSpecDef(sh, dep, body, req)
	sd.Position = sd.Position.Update(position) // update w/ 'spec' and 'where' positions
	return data.Ok(sd)
}

// rule:
//
//	```
//	spec dependency = "from", {"\n"}, pattern ;
//	```
func parseOptionalSpecDependency(p Parser) data.Either[data.Ers, data.Maybe[pattern]] {
	fromToken, found := getKeywordAtCurrent(p, token.From, dropAfter)
	if !found {
		return data.Ok(data.Nothing[pattern](p))
	}

	es, pat, isPat := ParsePattern(p).Break()
	if !isPat {
		return data.PassErs[data.Maybe[pattern]](es)
	}

	return data.Ok(data.Just(pat.updatePosPattern(fromToken)))
}

// rule:
//
//	```
//	spec body =
//		spec member
//		| "(", {"\n"}, spec member, {then, spec member}, {"\n"}, ")" ;
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
	es, mAnnots, isMAnnots := parseAnnotations_(p).Break()
	if !isMAnnots {
		return &es, data.Nothing[specMember](p)
	}

	// check for typing
	isTyping := lookahead2(p, typingLAs...)
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
	es, mAnnots, isMAnnots := parseAnnotations_(p).Break()
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
//		| "(", {"\n"}, [annotations_], def, {then, [annotations_], def}, {"\n"}, ")"
//		) ;
//	```
func parseOptionalRequiringClause(p Parser) data.Either[data.Ers, data.Maybe[requiringClause]] {
	req, found := getKeywordAtCurrent(p, token.Requiring, dropAfter)
	if !found {
		return data.Ok(data.Nothing[requiringClause](p)) // no requiring clause, return data.Nothing
	}

	type group struct{ data.NonEmpty[def] }
	es, reqBody, isReqBody := parseGroup[group](p, ExpectedDef, parseMaybeAnnotatedDef).Break()
	if !isReqBody {
		return data.PassErs[data.Maybe[requiringClause]](es)
	}
	reqBody.Position = reqBody.Update(req)
	return data.Ok(data.Just(reqBody.NonEmpty))
}

// rule:
//
//	```
//	spec inst = "inst", {"\n"}, spec head, [{"\n"}, spec inst target], {"\n"}, spec inst where clause ;
//	spec inst where clause = "where", {"\n"}, spec inst member group ;
//	spec inst member group = spec member | "(", {"\n"}, spec member, {then, spec member}, {"\n"}, ")" ;
//	```
func parseSpecInst(p Parser) data.Either[data.Ers, specInst] {
	inst, found := getKeywordAtCurrent(p, token.Inst, dropAfter)
	if !found {
		return data.Fail[specInst](ExpectedSpecInst, p)
	}

	es, sh, isSH := parseSpecHead(p).Break()
	if !isSH {
		return data.PassErs[specInst](es)
	}

	p.dropNewlines()
	esTarget, target, isTarget := parseOptionalSpecInstTarget(p).Break()
	if !isTarget {
		return data.PassErs[specInst](esTarget)
	}

	where, foundWhere := getKeywordAtCurrent(p, token.Where, dropBeforeAndAfter)
	if !foundWhere {
		return data.Fail[specInst](ExpectedInstWhere, target)
	}

	esBody, specBod, isBody := parseSpecBody(p).Break()
	if !isBody {
		return data.PassErs[specInst](esBody)
	}

	si := makeSpecInst(sh, target, specBod)
	si.Position = si.Update(inst).Update(where)
	return data.Ok(si)
}

// rule:
//
//	```
//	spec inst target = "=", {"\n"}, constrainer ;
//	```
func parseOptionalSpecInstTarget(p Parser) data.Either[data.Ers, data.Maybe[constrainer]] {
	equal, found := getKeywordAtCurrent(p, token.Equal, dropAfter)
	if !found {
		return data.Ok(data.Nothing[constrainer](p))
	}

	es, c, isC := parseConstrainer(p).Break()
	if !isC {
		return data.PassErs[data.Maybe[constrainer]](es)
	}
	c.Position = c.Update(equal)
	return data.Ok(data.Just(c))
}
