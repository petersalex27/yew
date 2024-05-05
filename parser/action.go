package parser

import (
	"github.com/petersalex27/yew/token"
)

type actionFunc func(parser *Parser, data *actionData) (term termElem, ok bool)

type actionMap map[token.Type]actionFunc

type resolutionMap map[NodeType]actionFunc

type actionData struct {
	actionMap   actionMap
	resolutions resolutionMap
	ptr         uint
	tokens      []token.Token
}

// constructs a new actionData and returns a pointer to it
func newActionData(m actionMap, r resolutionMap, toks []token.Token) *actionData {
	return &actionData{
		actionMap:   m,
		resolutions: r,
		ptr:         0,
		tokens:      toks,
	}
}

// constructs a new actionData with its internal token pointer offset by `offset` and returns a
// pointer to it
//
// panics if offset is greater than the length of the supplied tokens--note, it's okay to give an
// offset equal to the length of the tokens (this signifies that the actionData has reached the end
// of its input)
func newOffsetActionData(m actionMap, r resolutionMap, toks []token.Token, offset uint) *actionData {
	if offset > uint(len(toks)) {
		panic("bug: ptr offset is greater than the length of the supplied tokens")
	}

	data := newActionData(m, r, toks)
	data.ptr = offset
	return data
}

// returns true iff receiver has more input left to process
func (data *actionData) hasMoreInput() bool {
	return data.ptr < uint(len(data.tokens))
}

// attempts to return the next token in data's token stream and, on success, advances the stream
//
// on success, peek returns the next token in the stream and true; otherwise, if the stream length
// is 0, an end-of-tokens token and false are returned
func (data *actionData) nextToken() (tok token.Token, ok bool) {
	if tok, ok = data.peek(); ok {
		data.ptr++
	}
	return
}

func (parser *Parser) nextToken(data *actionData) (tok token.Token, ok bool) {
	if tok, ok = data.nextToken(); !ok {
		parser.error2(UnexpectedFinalTok, tok.Start, tok.End)
	}
	return
}

func (parser *Parser) peek(data *actionData) (tok token.Token, ok bool) {
	if tok, ok = data.peek(); !ok {
		parser.error2(UnexpectedFinalTok, tok.Start, tok.End)
	}
	return
}

// attempts to return the next token in data's token stream but does not advance the stream
//
// on success, peek returns the next token in the stream and true; otherwise, if the stream length
// is 0, an end-of-tokens token and false are returned
func (data *actionData) peek() (tok token.Token, ok bool) {
	tokenLength := len(data.tokens)
	if ok = uint(tokenLength) > data.ptr; ok { // assumes len != 0
		tok = data.tokens[data.ptr]
		return
	}

	if tokenLength == 0 {
		tok = endOfTokensToken()
	} else {
		tok = data.tokens[tokenLength-1]
	}
	return
}

// attempts to lookup action given by `tokenType`
//
// on success, this returns the action and true; otherwise, returns _, false
func (data *actionData) findAction(tokenType token.Type) (action actionFunc, found bool) {
	action, found = data.actionMap[tokenType]
	return action, found
}

// attempts to lookup action given by NodeType `nt`
//
// on success, this returns the action and true; otherwise, returns _, false
func (data *actionData) findResolution(nt NodeType) (action actionFunc, found bool) {
	action, found = data.resolutions[nt]
	return action, found
}

func (parser *Parser) terminalAction(data *actionData) (term termElem, ok bool) {
	if parser.termMemory != nil {
		term = *parser.termMemory
		parser.forget()
		return term, true
	}
	var tok token.Token
	if tok, ok = data.peek(); !ok {
		parser.error2(UnexpectedFinalTok, tok.Start, tok.End)
		return
	}

	var action actionFunc
	if action, ok = data.findAction(tok.Type); !ok {
		parser.error2(UnexpectedToken, tok.Start, tok.End)
		return
	}
	return action(parser, data)
}

func constraintAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	panic("TODO: implement") // TODO: implement
}

func arrowInfo() termInfo {
	return termInfo{bp: 1, rAssoc: true, arity: 2}
}

func makeFunctionType(tok token.Token) FunctionType {
	return FunctionType{nil, nil, tok.Start, tok.End}
}

// creates a function type
func productAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	var tok token.Token
	if tok, ok = parser.nextToken(data); !ok {
		return
	}

	f := makeFunctionType(tok)
	term = termElem{f, arrowInfo()}
	return
}

func idAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	var tok token.Token
	if tok, ok = parser.nextToken(data); !ok {
		return
	}
	if decl, ok := parser.lookupTerm(tok); ok {
		return decl.makeTerm(), true
	}
	ident := makeIdent(tok)
	term = termElem{ident, termInfo{}}
	return term, true
}

func intAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	var tok token.Token
	if tok, ok = parser.nextToken(data); !ok {
		return
	}
	v := parser.makeInt(tok)
	return termElem{v, termInfo{}}, true
}

func charAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	var tok token.Token
	if tok, ok = parser.nextToken(data); !ok {
		return
	}
	v := parser.makeChar(tok)
	return termElem{v, termInfo{}}, true
}

func stringAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	var tok token.Token
	if tok, ok = parser.nextToken(data); !ok {
		return
	}
	v := parser.makeString(tok)
	return termElem{v, termInfo{}}, true
}

func floatAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	var tok token.Token
	if tok, ok = parser.nextToken(data); !ok {
		return
	}
	v := parser.makeFloat(tok)
	return termElem{v, termInfo{}}, true
}

func makeOpener(tok token.Token) termElem {
	return termElem{Key{tok.Value, tok.Start, tok.End}, termInfo{}}
}

func parenAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	var tok token.Token
	if tok, ok = parser.nextToken(data); !ok {
		return
	}

	term = makeOpener(tok)
	parser.shift(term)

	pushFirst := parser.terms.GetCount() > 0
	parser.terms.Save()
	term, ok = parser.terminalAction(data)
	if !ok {
		return
	} else if !pushFirst {
		return
	}

	parser.shift(term)
	return parser.terminalAction(data) // TODO: what happens on the following "fun (+) y"?
}

// when a term is enclosed by parens, it's information is lost and becomes a more abstract function
func abstractInfo(info termInfo) termInfo {
	if info.Arity() > 0 && info.bp > 0 {
		// update binding power and associativity for enclosed functions
		//	- they lose whatever the inner bp and assoc is
		info.bp = 10
		info.rAssoc = false
	}
	return info
}

func closeParenAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	var tok token.Token
	if tok, ok = parser.nextToken(data); !ok {
		return
	}

	var lp termElem
	// loop reducing anything that needs to be reduced
	for {
		if term, ok = parser.reduceStack(); !ok {
			return
		}
		if ok = !parser.terms.FullEmpty(); !ok {
			parser.error2(UnexpectedRParen, tok.Start, tok.End)
			return
		}

		parser.terms.Return() // return stack
		if ok = parser.terms.GetCount() != 0; !ok {
			// error: size of stack frame is 0
			parser.error2(UnexpectedRParen, tok.Start, tok.End)
			return
		}

		lp, _ = parser.terms.Peek()
		if lp.NodeType() != lambdaType {
			break
		}

		// push reduced bound expression
		parser.shift(term)
		// resolve lambda abstraction
		term, ok = resolveLambdaAbstraction(parser, data)
		if !ok {
			return
		}
		parser.shift(term) // push result, try again
	}

	// at this point, lp should be a left paren
	start, end := lp.Pos()
	isLParen := lp.Term.NodeType() == syntaxExtensionType && lp.Term.String() == "("
	if ok = isLParen; !ok {
		parser.error2(ExpectedLParen, start, end)
		return
	}

	// remove left paren
	_, _ = parser.terms.Pop()

	end = tok.End // end of right paren

	// create enclosed term
	info := abstractInfo(term.termInfo)
	term = termElem{EnclosedTerm{term.Term, start, end}, info}
	return
}

var resolutionActions resolutionMap = resolutionMap{
	lambdaType: resolveLambdaAbstraction,
}

var standardActions actionMap = actionMap{
	token.Id:          idAction,
	token.ImplicitId:  idAction,
	token.Backslash:   abstractionAction,
	token.ThickArrow:  constraintAction,
	token.Arrow:       productAction,
	token.IntValue:    intAction,
	token.CharValue:   charAction,
	token.StringValue: stringAction,
	token.FloatValue:  floatAction,
	token.LeftParen:   parenAction,
	token.RightParen:  closeParenAction,
}
