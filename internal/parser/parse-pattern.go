package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/api/util/fun"
	"github.com/petersalex27/yew/common/data"
)

func ParsePattern(p Parser) data.Either[data.Ers, pattern] {
	return parsePattern(p, false)
}

// rule:
//
//	```
//	pattern atom = literal | name | "[]" | hole ;
//	```
func parsePatternAtom(p Parser) data.Either[data.Ers, patternAtom] {
	es, mPatAtom := maybePatternAtom(p)
	if es != nil {
		return data.PassErs[patternAtom](*es)
	} else if r, just := mPatAtom.Break(); !just {
		return data.Fail[patternAtom](ExpectedPattern, p)
	} else {
		return data.Ok(r)
	}
}

func maybePatternAtom(p Parser) (*data.Ers, data.Maybe[patternAtom]) {
	// pattern atom as literal
	if lookahead1(p, literalL1s...) {
		lit := literalAsPatternAtom(p.current())
		p.advance()
		return nil, data.Just(lit)
	}

	// pattern atom as name, [], or hole
	var n name
	var isSomething bool
	if matchCurrent(token.EmptyBracketEnclosure)(p) {
		n = data.EOne[name](p.current())
		p.advance()
	} else if matchCurrent(token.Hole)(p) {
		h := holeAsPatternAtom(p.current())
		p.advance()
		return nil, data.Just(h)
	} else if n, isSomething = maybeParseName(p).Break(); !isSomething {
		return nil, data.Nothing[patternAtom](p)
	}
	return nil, data.Just(nameAsPatternAtom(n))
}

// rule:
//
//	```
//	pattern = pattern term, {pattern term} ;
//	enc pattern = pattern term, {{"\n"}, pattern term} ;
//	```
func parsePattern(p Parser, enclosed bool) data.Either[data.Ers, pattern] {
	es, mPat := maybeParsePattern(p, enclosed)
	if es != nil {
		return data.PassErs[pattern](*es)
	} else if pat, just := mPat.Break(); !just {
		return data.Fail[pattern](ExpectedPattern, p)
	} else {
		return data.Ok(pat)
	}
}

func maybeParsePattern(p Parser, enclosed bool) (*data.Ers, data.Maybe[pattern]) {
	var first pattern
	var just bool
	es, mFirst := maybeParsePatternTerm(p, enclosed)
	if es != nil {
		return es, data.Nothing[pattern](p)
	} else if first, just = mFirst.Break(); !just {
		return nil, data.Nothing[pattern](p)
	}

	maybeFunc := fun.BinBind1st_PairTarget(maybeParsePatternTerm, enclosed)
	es, res, has2ndTerm := parseOneOrMore[pattern](p, first, enclosed, maybeFunc)
	if es != nil {
		return es, data.Nothing[pattern](p)
	} else if !has2ndTerm {
		return nil, data.Just(res.Head())
	}
	// this must have at least 2 elements, so it's safe to call Break() w/o checking
	app, _ := data.NonEmptyToAppLikePair[patternApp](res).Break()
	return nil, data.Just[pattern](app)
}

func closeEnclosedPattern(p Parser, opener api.Token, closerType token.Type) func(patternEnclosed) data.Either[data.Ers, pattern] {
	return func(ps patternEnclosed) data.Either[data.Ers, pattern] {
		ps.Position = ps.Update(opener)
		if !parseKeywordAtCurrent(p, closerType, &ps.Position) {
			isRp := token.RightParen.Int() == closerType.Int()
			return data.Fail[pattern](ifThenElse(isRp, ExpectedRightParen, ExpectedRightBrace), p)
		}
		ps.implicit = token.LeftBrace.Match(opener)
		return data.Ok[pattern](ps)
	}
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
//	enc pattern inner = enc pattern, {{"\n"}, ",", enc pattern}, [{"\n"}, ","] ;
//	```
func maybeParsePatternTerm(p Parser, enclosed bool) (*data.Ers, data.Maybe[pattern]) {
	if enclosed && matchCurrentEqual(p) {
		// allow '=' to be used as a pattern term as long as the term(s) is/are enclosed
		eq := pattern(data.EOne[name](p.current()))
		p.advance()
		return nil, data.Just(eq)
	}

	es, mAtom := maybePatternAtom(p)
	if es != nil {
		return es, data.Nothing[pattern](p)
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
		return nil, data.Just[pattern](data.EOne[wildcard](wildcardToken))
	}

	opener, closerType, found := parseEnclosedOpener(p)
	if !found {
		return nil, data.Nothing[pattern](p)
	}

	f := func(p Parser) (*data.Ers, data.Maybe[pattern]) { return maybeParsePattern(p, true) }

	esPats, pats, isPatsRight := parseSepSequenced[struct{ data.NonEmpty[pattern] }](p, ExpectedPattern, token.Comma, f).Break()
	if !isPatsRight {
		return &esPats, data.Nothing[pattern](p)
	}
	pe := patternEnclosed{NonEmpty: pats.NonEmpty}

	p.dropNewlines()
	res := closeEnclosedPattern(p, opener, closerType)(pe)
	esOut, out, isRight := res.Break()
	if !isRight {
		return &esOut, data.Nothing[pattern](p)
	}
	return nil, data.Just(out)
}
