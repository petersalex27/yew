package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/api/util/fun"
	"github.com/petersalex27/yew/common/data"
	"github.com/petersalex27/yew/internal/common"
)

var (
	// cast token to node
	tokenAsNode = func(t api.Token) api.Node { return t }
	// returns parser's the current token as a node
	currentTokenAsNode = fun.Compose(tokenAsNode, (Parser).current)
	// given a token type, returns function that takes in a parser and then tests if the parser's
	// current token matches that type
	matchCurrent = fun.ComposeRightCurryFlip((token.Type).Match, currentTokenAsNode)

	matchCurrentWith        = matchCurrent(token.With)
	matchCurrentId          = matchCurrent(token.Id)
	matchCurrentEndOfTokens = matchCurrent(token.EndOfTokens)
	matchCurrentEqual       = matchCurrent(token.Equal)
	matchCurrentBackslash   = matchCurrent(token.Backslash)
	matchCurrentLeftParen   = matchCurrent(token.LeftParen)
	matchCurrentUnderscore  = matchCurrent(token.Underscore)
	matchCurrentImpossible  = matchCurrent(token.Impossible)
	matchCurrentForall      = matchCurrent(token.Forall)
	matchCurrentInfix       = matchCurrent(token.Infix)
	matchCurrentRawString   = matchCurrent(token.RawStringValue)

	notMatchCurrent = func(t token.Type) func(Parser) bool {
		return func(p Parser) bool { return !matchCurrent(t)(p) }
	}

	literalL1s = []token.Type{token.IntValue, token.FloatValue, token.StringValue, token.RawStringValue, token.ImportPath, token.CharValue}

	exprInTypeL1s = append(([]token.Type{token.Backslash, token.Id, token.Infix}), literalL1s...)

	boundIdentL2 = [][2]token.Type{{token.LeftBrace, token.Id}}

	typingL2 = [][2]token.Type{{token.Id, token.Colon}, {token.Infix, token.Colon}}
)

func bind[a, b api.Node](ma data.Maybe[a], f func(a) data.Maybe[b]) data.Maybe[b] {
	x, isJust := ma.Break()
	if !isJust {
		return data.Nothing[b](ma)
	}
	return f(x)
}

func currentIsUpperIdent(p Parser) bool {
	return matchCurrent(token.Id)(p) && common.Is_PascalCase(p.current().String())
}

func currentIsName(p Parser) bool {
	return matchCurrentId(p) || matchCurrentInfix(p)
}

func ifThenElse[a any](cond bool, true_, false_ a) a {
	if cond {
		return true_
	}
	return false_
}

func lookahead1Report(p Parser, types ...token.Type) (tt token.Type, found bool) {
	for _, t := range types {
		if t.Match(p.current()) {
			return t, true
		}
	}
	return tt, false
}

func lookahead1(p Parser, types ...token.Type) bool {
	for _, t := range types {
		if t.Match(p.current()) {
			return true
		}
	}
	return false
}

// performs a lookahead 2-ish. If the current token matches the first type in the pair, then it will try to
// match the second type in the pair, dropping newlines in between the two
func lookahead2(p Parser, types ...[2]token.Type) bool {
	ps, ok := p.(*ParserState)
	if !ok {
		var ps_o *ParserState_optional
		if ps_o, ok = p.(*ParserState_optional); ok {
			ps = ps_o.ParserState
		}
	}

	if !ok {
		return false
	}

	origin := ps.tokenCounter // capture the current token counter for restoration
	defer func() { ps.tokenCounter = origin }()

	for _, tt := range types {
		if tt[0].Match(p.current()) {
			ps.advance()
			ps.dropNewlines()
			if tt[1].Match(p.current()) {
				return true
			}
		}
	}
	return false
}

// if the current token's type matches `keyword`, updates the position at `pos` and returns
// true; otherwise, returns false and leaves `pos` unchanged.
//
//	```
//	rule = [keyword, {"\n"}] ;
//	```
//
// NOTE: if `pos` is nil, then everything above happens with the exception that the position is not updated
//
// SEE: `getKeywordAtCurrent` to return the token and found status instead
func parseKeywordAtCurrent(p Parser, keyword token.Type, pos *api.Position) (found bool) {
	var token api.Token
	if token, found = getKeywordAtCurrent(p, keyword); found {
		if pos != nil {
			*pos = pos.Update( /* keyword */ token)
		}
	}
	return found
}

func parseEnclosedOpener(p Parser) (opener api.Token, closerType token.Type, found bool) {
	closerType = token.RightParen
	opener, found = getKeywordAtCurrent(p, token.LeftParen)
	if !found {
		closerType = token.RightBrace
		opener, found = getKeywordAtCurrent(p, token.LeftBrace)
	}
	return opener, closerType, found
}

func getKeywordAtCurrent(p Parser, keyword token.Type) (token api.Token, found bool) {
	if found = keyword.Match(p.current()); found {
		token = p.current()
		p.advance()
		p.dropNewlines()
	}
	return token, found
}

func writeErrors(p Parser, es data.Ers) Parser {
	var out Parser
	for _, e := range es.Elements() {
		out = p.report(parseError(p, e), e.Fatal())
	}
	return out
}

func maybeParseName(p Parser) data.Maybe[name] {
	t := p.current()
	if !currentIsName(p) {
		return data.Nothing[name](t)
	}
	p.advance()
	return data.Just(data.EOne[name](t))
}

func createLowerIdent(_ Parser, t api.Token) data.Maybe[lowerIdent] {
	if name := t.String(); !common.Is_camelCase(name) {
		return data.Nothing[lowerIdent](t)
	}
	return data.Just(data.EOne[lowerIdent](t))
}

func createUpperIdent(t api.Token) func(Parser) data.Maybe[upperIdent] {
	return func(p Parser) data.Maybe[upperIdent] {
		if name := t.String(); !common.Is_PascalCase(name) {
			return data.Nothing[upperIdent](t)
		}
		return data.Just(upperIdent{data.One(t)})
	}
}

var liftGenLowerIdent = fun.Compose(data.Just, data.Inl[upperIdent, lowerIdent])
var liftGenUpperIdent = fun.Compose(data.Just, data.Inr[lowerIdent, upperIdent])

func parseIdent(p Parser) data.Maybe[ident] {
	t := p.current()
	if !token.Id.Match(t) {
		return data.Nothing[ident](t)
	}
	p.advance()
	upper := thunk(p, createUpperIdent(t))
	res := bind(createLowerIdent(p, t), liftGenLowerIdent)
	if res.IsNothing() {
		return bind(upper(), liftGenUpperIdent)
	}
	return res
}

func parseLowerIdent(p Parser) data.Maybe[lowerIdent] {
	t := p.current()
	if !token.Id.Match(t) {
		return data.Nothing[lowerIdent](t)
	}
	p.advance()
	return createLowerIdent(p, t)
}

// parses a rule pattern `group` (parameterized by `mem`):
//
//	```
//	group <mem> = mem | "(", {"\n"}, mem, {{"\n"}, mem}, {"\n"}, ")" ;
//	```
func parseGroup[ne data.EmbedsNonEmpty[a], a api.Node](p Parser, errorMsg string, maybeParse func(Parser) (*data.Ers, data.Maybe[a])) data.Either[data.Ers, ne] {
	leftParen, found := getKeywordAtCurrent(p, token.LeftParen) // parse '('
	var first a
	if es, mFirst := maybeParse(p); es != nil {
		return data.PassErs[ne](*es)
	} else if unit, just := mFirst.Break(); !just {
		return data.Fail[ne](errorMsg, p)
	} else {
		first = unit
	}

	var xs data.NonEmpty[a]
	var es *data.Ers
	if found { // if '(' was found, parse multiple elements and then ')'
		es, xs, _ = parseOneOrMore(p, first, true, maybeParse)
		if es != nil {
			return data.PassErs[ne](*es)
		}
		xs.Position = xs.Update(leftParen)
		rp, found := getKeywordAtCurrent(p, token.RightParen)
		if !found {
			return data.Fail[ne](ExpectedRightParen, p)
		}
		xs.Position = xs.Update(rp)
	} else { // otherwise, just return the first element
		xs = data.Singleton(first)
	}

	return data.Ok(ne{xs})
}

func actionThunk[a, b any](x a, f func(a) b) func() b {
	return func() b { return f(x) }
}

func thunk[a any](p Parser, f func(Parser) a) func() a {
	return actionThunk(p, f)
}

// lhs - the thing returned if there is no rhs; otherwise, the first thing in the non-empty list
// dropNewlinesEachIt - if true, calls `p.dropNewlines()` at the start of each loop iteration
func parseOneOrMore[a api.Node](p Parser, lhs a, dropNewlinesEachIt bool, f func(Parser) (*data.Ers, data.Maybe[a])) (_ *data.Ers, _ data.NonEmpty[a], has2ndTerm bool) {
	group := data.Singleton(lhs)
	has2ndTerm = false
	for {
		if dropNewlinesEachIt {
			p.dropNewlines()
		}

		es, mRhs := f(p)
		if es != nil {
			return es, group, has2ndTerm
		}

		rhs, isSomething := mRhs.Break()
		if !isSomething {
			return nil, group, has2ndTerm
		}

		has2ndTerm = true
		group = group.Snoc(rhs)
	}
}

func parseHandledSepSequenced[b data.EmbedsNonEmpty[a], a api.Node](p Parser, errHandler func(cur api.Token) string, sep token.Type, maybeParse func(Parser) (*data.Ers, data.Maybe[a])) data.Either[data.Ers, b] {
	ps := api.ZeroPosition()
	es, lhs := maybeParse(p)
	if es != nil {
		return data.PassErs[b](*es)
	}

	unit, just := lhs.Break()
	if !just {
		return data.Fail[b](errHandler(p.current()), p) // error when less than 1 element
	}
	ps = ps.Update(unit)

	var terms = data.Singleton(unit)
	for {
		p.dropNewlines()
		if !parseKeywordAtCurrent(p, sep, &ps) {
			break // no 'sep', end loop
		}

		es, mRhs := maybeParse(p)
		ps = ps.Update(mRhs)
		if es != nil {
			return data.PassErs[b](*es)
		} else if rhs, just := mRhs.Break(); !just {
			break // no rhs, trailing comma, end loop
		} else {
			terms = terms.Snoc(rhs)
		}
	}
	terms.Position = ps
	return data.Ok(b{terms})
}

func parseSepSequenced[b data.EmbedsNonEmpty[a], a api.Node](p Parser, emptyErrorMsg string, sep token.Type, maybeParse func(Parser) (*data.Ers, data.Maybe[a])) data.Either[data.Ers, b] {
	return parseHandledSepSequenced[b](p, fun.Constant[api.Token](emptyErrorMsg), sep, maybeParse)
}
