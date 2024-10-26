package parser

import (
	"strconv"
	"strings"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/api/util/fun"
	"github.com/petersalex27/yew/common/data"
	nodeType "github.com/petersalex27/yew/internal/parser/typ"
)

func ParsePattern(p parser) data.Either[data.Ers, Pattern] {
	return parsePattern(p, false)
}

// rule:
//
//	```
//	pattern atom = literal | name | "[]" | hole ;
//	```
func parsePatternAtom(p parser) data.Either[data.Ers, patternAtom] {
	es, mPatAtom := maybePatternAtom(p)
	if es != nil {
		return data.PassErs[patternAtom](*es)
	} else if r, just := mPatAtom.Break(); !just {
		return data.Fail[patternAtom](ExpectedPattern, p)
	} else {
		return data.Ok(r)
	}
}

func parseIntLiteral(lit api.Token) (*data.Ers, aLiteral[int64]) {
	strVal := lit.String()
	strVal = strings.Replace(strVal, "_", "", -1)
	strVal = strings.ToLower(strVal) // convert any prefix 0X, etc. to lowercase 0x, etc.
	val, err := strconv.ParseInt(strVal, 0, 64)
	if err != nil {
		e := data.MkErr(err.Error(), lit)
		es := data.Makes(e)
		return &es, aLiteral[int64]{}
	}
	literal := makeLiteral(nodeType.LiteralInt, val)
	literal.Position = literal.Update(lit)
	return nil, literal
}

func parseFloatLiteral(lit api.Token) (*data.Ers, aLiteral[float64]) {
	strVal := lit.String()
	strVal = strings.Replace(strVal, "_", "", -1)
	strVal = strings.ToLower(strVal)
	val, err := strconv.ParseFloat(strVal, 64)
	if err != nil {
		e := data.MkErr(err.Error(), lit)
		es := data.Makes(e)
		return &es, aLiteral[float64]{}
	}

	literal := makeLiteral(nodeType.LiteralFloat, val)
	literal.Position = literal.Update(lit)
	return nil, literal
}

func parseStringLiteral(lit api.Token) (*data.Ers, aLiteral[[]byte]) {
	str := lit.String()
	literal := makeLiteral(nodeType.LiteralString, []byte(str))
	literal.Position = literal.Update(lit)
	return nil, literal
}

func parseCharLiteral(lit api.Token) (*data.Ers, aLiteral[byte]) {
	str := lit.String()
	if len(str) > 1 {
		panic("unexpected token") // scanner should have caught this, or wrong token passed as arg
	}

	literal := makeLiteral(nodeType.LiteralChar, str[0])
	literal.Position = literal.Update(lit)
	return nil, literal
}

func parseLiteral(p parser) (es *data.Ers, literal literal) {
	lit := p.current()
	if token.IntValue.Match(lit) {
		es, literal = parseIntLiteral(lit)
	} else if token.FloatValue.Match(lit) {
		es, literal = parseFloatLiteral(lit)
	} else if token.StringValue.Match(lit) {
		es, literal = parseStringLiteral(lit)
	} else if token.CharValue.Match(lit) {
		es, literal = parseCharLiteral(lit)
	} else {
		// bug: parser's state was not validated before calling this function
		panic("unexpected token type")
	}

	p.advance()
	return es, literal
}

func maybePatternAtom(p parser) (*data.Ers, data.Maybe[patternAtom]) {
	// pattern atom as literal
	if lookahead1(p, literalLAs...) {
		es, lit := parseLiteral(p)
		if es != nil {
			return es, data.Nothing[patternAtom](p)
		}
		return nil, data.Just(data.Inl[patternName](lit))
	}

	// name, [], (), or hole as pattern atom
	var n name
	var isSomething bool
	if matchCurrent(token.EmptyBracketEnclosure)(p) {
		nilList := p.current()
		n = data.EOne[name](nilList)
		p.advance()
	} else if matchCurrent(token.EmptyParenEnclosure)(p) {
		unit := p.current()
		n = data.EOne[name](unit)
		p.advance()
	} else if matchCurrent(token.Hole)(p) {
		hole := p.current()
		h := holeAsPatternAtom(hole)
		p.advance()
		return nil, data.Just(h)
	} else if n, isSomething = maybeParseName(p).Break(); !isSomething {
		return nil, data.Nothing[patternAtom](p)
	}
	
	return declareFromName(p, n.Break()), data.Just(nameAsPatternAtom(n))
}

// rule:
//
//	```
//	pattern = pattern term, {pattern term} ;
//	enc pattern = pattern term, {{"\n"}, pattern term} ;
//	```
func parsePattern(p parser, enclosed bool) data.Either[data.Ers, Pattern] {
	es, mPat := maybeParsePattern(p, enclosed)
	if es != nil {
		return data.PassErs[Pattern](*es)
	} else if pat, just := mPat.Break(); !just {
		return data.Fail[Pattern](ExpectedPattern, p)
	} else {
		return data.Ok(pat)
	}
}

func maybeParsePattern(p parser, enclosed bool) (*data.Ers, data.Maybe[Pattern]) {
	var first Pattern
	var just bool
	es, mFirst := maybeParsePatternTerm(p, enclosed)
	if es != nil {
		return es, data.Nothing[Pattern](p)
	} else if first, just = mFirst.Break(); !just {
		return nil, data.Nothing[Pattern](p)
	}

	maybeFunc := fun.BinBind1st_PairTarget( /* note RHS here! */ maybeParsePatternTermRhs, enclosed)
	es, res, has2ndTerm := parseOneOrMore(p, first, enclosedDependentIt(enclosed), maybeFunc)
	if es != nil {
		return es, data.Nothing[Pattern](p)
	} else if !has2ndTerm {
		return nil, data.Just(res.Head())
	}
	// this must have at least 2 elements, so it's safe to call Break() w/o checking
	app, _ := data.NonEmptyToAppLikePair[patternApp](res).Break()
	return nil, data.Just[Pattern](app)
}

func closeEnclosedPattern(p parser, opener api.Token, closerType token.Type) func(patternEnclosed) data.Either[data.Ers, Pattern] {
	return func(ps patternEnclosed) data.Either[data.Ers, Pattern] {
		ps.Position = ps.Update(opener)
		closer, found := getKeywordAtCurrent(p, closerType, dropBefore)
		if !found {
			isRp := token.RightParen.Int() == closerType.Int()
			return data.Fail[Pattern](ifThenElse(isRp, ExpectedRightParen, ExpectedRightBrace), p)
		}
		ps.implicit = token.LeftBrace.Match(opener)
		ps.Position = ps.Update(closer)
		return data.Ok[Pattern](ps)
	}
}

// rule:
//
//	```
//	access = ".", {"\n"}, name ;
//	```
func parseAccess(p parser) (_ *data.Ers, acc accessList) {
	dot, found := getKeywordAtCurrent(p, token.Dot, dropAfter)
	if !found {
		e := data.Nil[data.Err](1).Snoc(data.MkErr(ExpectedAccessDot, p))
		return &e, acc
	}

	res := parseSepSequenced[accessList](
		p, ExpectedName, token.Dot,
		func(p parser) (*data.Ers, data.Maybe[name]) {
			return nil, maybeParseName(p)
		},
	)

	es, names, isNames := res.Break()
	if !isNames {
		return &es, acc
	}

	// declare all names in the access list
	for _, n := range names.Elements() {
		declareFromName(p, n.Break())
	}

	acc = names
	acc.Position = acc.Update(dot)
	return nil, acc
}

func maybeParsePatternTermHelper(p parser, enclosed bool, rhs bool) (*data.Ers, data.Maybe[Pattern]) {
	if rhs && token.Dot.Match(p.current()) {
		es, res := parseAccess(p)
		if es != nil {
			return es, data.Nothing[Pattern](p)
		}
		return nil, data.Just[Pattern](res)
	} else if enclosed && matchCurrentEqual(p) {
		// allow '=' to be used as a pattern term as long as the term(s) is/are enclosed
		equalTok := p.current()
		eqName := data.EOne[name](equalTok)
		eq := Pattern(eqName)
		p.advance()
		return declareFromName(p, eqName.Break()), data.Just(eq)
	}

	es, mAtom := maybePatternAtom(p)
	if es != nil {
		return es, data.Nothing[Pattern](p)
	} else if unit, just := mAtom.Break(); just {
		return nil, data.Just(
			data.Cases(unit,
				(literal).asPattern,
				patternNameAsPattern,
			),
		)
	}

	if matchCurrentUnderscore(p) {
		wildcardToken := p.current()
		p.advance()
		return nil, data.Just[Pattern](data.EOne[wildcard](wildcardToken))
	}

	opener, closerType, found := parseEnclosedOpener(p)
	if !found {
		return nil, data.Nothing[Pattern](p)
	}

	f := func(p parser) (*data.Ers, data.Maybe[Pattern]) { return maybeParsePattern(p, true) }

	esPats, pats, isPatsRight := parseSepSequenced[struct{ data.NonEmpty[Pattern] }](p, ExpectedPattern, token.Comma, f).Break()
	if !isPatsRight {
		return &esPats, data.Nothing[Pattern](p)
	}
	pe := patternEnclosed{NonEmpty: pats.NonEmpty}

	res := closeEnclosedPattern(p, opener, closerType)(pe)
	esOut, out, isRight := res.Break()
	if !isRight {
		return &esOut, data.Nothing[Pattern](p)
	}
	return nil, data.Just(out)
}

// rule:
//
//	```
//	enc pattern term = "=" | pattern term ;
//	pattern term =
//		pattern atom
//		| "_"
//		| "(", {"\n"}, enc pattern inner, {"\n"}, ")"
//		| "{", {"\n"}, enc pattern inner, {"\n"}, "}" ;
//	enc pattern inner = enc pattern, {{"\n"}, ",", {"\n"}, enc pattern}, [{"\n"}, ","] ;
//	```
func maybeParsePatternTerm(p parser, enclosed bool) (*data.Ers, data.Maybe[Pattern]) {
	return maybeParsePatternTermHelper(p, enclosed, false)
}

func maybeParsePatternTermRhs(p parser, enclosed bool) (*data.Ers, data.Maybe[Pattern]) {
	return maybeParsePatternTermHelper(p, enclosed, true)
}
