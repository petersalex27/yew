package parser

import (
	"github.com/petersalex27/yew/token"
)

type actionFunc func(parser *Parser, data *actionData) (term termElem, ok bool)

type actionMapperId uint8
type actionMapper struct {
	actionMap
	id actionMapperId
}

type actionMap map[token.Type]actionFunc

type actionData struct {
	m           actionMapper
	resolutions resolutionMap
	ptr         uint
	tokens      []token.Token
}

// constructs a new actionData and returns a pointer to it
func newActionData(m actionMapper, r resolutionMap, toks []token.Token) *actionData {
	return &actionData{
		m:           m,
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
func newOffsetActionData(m actionMapper, r resolutionMap, toks []token.Token, offset uint) *actionData {
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
	action, found = data.m.actionMap[tokenType]
	return action, found
}

func (parser *Parser) terminalAction(data *actionData) (term termElem, ok bool) {
	if parser.termMemory != nil {
		term = *parser.termMemory
		parser.forget()
		return term, true
	}
	var tok token.Token
	if tok, ok = parser.peek(data); !ok {
		return
	}

	var action actionFunc
	if action, ok = data.findAction(tok.Type); !ok {
		parser.error2(UnexpectedToken, tok.Start, tok.End)
		return
	}
	return action(parser, data)
}

func arrowInfo() termInfo {
	return termInfo{bp: 1, rAssoc: true, arity: 2}
}

func makeFunctionType(tok token.Token) FunctionType {
	return FunctionType{nil, nil, nil, tok.Start, tok.End}
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

func listingAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	var comma token.Token
	if comma, ok = parser.nextToken(data); !ok {
		return
	}

	listing := Listing{[]Term{}, comma.Start, comma.End}
	term = termElem{listing, termInfo{bp: 0, rAssoc: true, arity: 2}}
	return term, true
}

func errorActionGen(msg string) func(parser *Parser, data *actionData) (term termElem, ok bool) {
	return func(parser *Parser, data *actionData) (term termElem, ok bool) {
		tok, _ := data.peek()
		parser.error2(msg, tok.Start, tok.End)
		return term, false
	}
}

// data type initial name parsing
func idDataTypeAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	var tok token.Token
	if tok, ok = parser.nextToken(data); !ok {
		return
	}

	// see if already declared
	if decl, found := parser.lookupTerm(tok); found {
		// error: redeclaration
		parser.error2(IllegalRedeclaration, decl.name.Start, decl.name.End)
		ok = false
		return
	}

	// create new data type declaration
	ident := makeIdent(tok)
	term = termElem{ident, termInfo{}}
	// TODO: declare
	return term, true
}

func defaultIdTermMaker(tok token.Token) termElem {
	ident := makeIdent(tok)
	term := termElem{ident, termInfo{}}
	return term
}

func idAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	var tok token.Token
	if tok, ok = parser.nextToken(data); !ok {
		return term, false
	}

	// attempt to find declaration, converting it into a term if found
	if term, ok = parser.findDeclAsTerm(tok); ok {
		return term, true
	}

	return defaultIdTermMaker(tok), true
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

func closingImplicit(parser *Parser, term termElem, start, end int) (termElem, bool) {
	switch term.Term.NodeType() {
	case identType:
		// assume this is a term of 'Type'
		id, _ := term.Term.(Ident)
		decl := new(Declaration)
		decl.implicit = true
		decl.name = id
		decl.termInfo = termInfo{}
		decl.typing = Ident{"Type", 0, 0}
		// create new entry
		parser.locals.Map(id, decl)
	case typingType:
		// this case updates existing locals, marking them as implicit

		if term.Term.NodeType() != identType {
			break
		}

		// entry should already exist
		decl, found := parser.locals.Find(term.Term)
		if !found {
			break // ? dunno when this would happen
		}
		// mark as implicit
		decl.implicit = true
		parser.locals.Map(term.Term, decl) // update
	default:
		break
	}

	info := abstractInfo(term.termInfo)
	term = termElem{Implicit{term.Term, start, end}, info}
	return term, true
}

func wildcardAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	panic("TODO: implement") // TODO: implement
}

// creates typing prototype--still requires a type
func createTyping(colon token.Token) termElem {
	typing := Typing{Start: colon.Start, End: colon.End}
	info := termInfo{bp: 0, rAssoc: true, arity: 2}
	return termElem{typing, info}
}

// typing, e.g.,
//
//	x : A
//
// But, specifically, typing in a type; introduces bound occurrence of term. Implicit types w/
// labels are available inside the body of the function they describe
func labeledTypeAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	// read past ':' token
	var colon token.Token
	if colon, ok = parser.nextToken(data); ok {
		term, ok = createTyping(colon), true
	}
	return
}

// begins parsing of an implicit argument or type
func implicitAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	var openBrace token.Token
	if openBrace, ok = parser.nextToken(data); !ok {
		return
	}

	term = makeOpener(openBrace)
	parser.shift(term)

	pushFirst := parser.terms.GetCount() > 0
	parser.terms.Save()

	// set parsing rules depending on the syntactic structure the implicit term is located in
	// there should only be two possibilities, type and function arguments
	switch data.m.id {
	case typingId:
		data.m = typePositionImplicitAction
	case standardId:
		data.m = argPositionImplicitAction
	default:
		panic("bug: unexpected action map")
	}

	term, ok = parser.terminalAction(data)
	if !ok {
		return
	} else if !pushFirst {
		return
	}

	parser.shift(term)
	return parser.terminalAction(data)
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

	pushFirst := parser.terms.GetCount() > 0

	parser.shift(term)
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

type parenClosingFunc = func(*Parser, termElem, int, int) (termElem, bool)

func closeAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	var tok token.Token
	if tok, ok = parser.nextToken(data); !ok {
		return
	}

	var m Marker = Marker{Start: tok.Start, End: tok.End}
	info := termInfo{11, false, 1}
	switch tok.Type {
	case token.RightParen:
		m.nodeType = closeParenType
	case token.RightBracket:
		m.nodeType = closeBracketType
	case token.RightBrace:
		m.nodeType = closeBraceType
	default:
		parser.errorOnToken(UnexpectedToken, tok)
		return term, false
	}
	term = termElem{m, info}
	return term, true
}

const (
	constraintId actionMapperId = iota
	standardId
	argPosImplicitId
	typePosImplicitId
	typingId
	dataTypeFollowId
	dataTypeInitId
)

var (
	argPositionImplicitAction  actionMapper
	standardActions            actionMapper
	constraintActions          actionMapper
	typePositionImplicitAction actionMapper
	typingActions              actionMapper
	// data type parsing after name is known
	dataTypeFollowActions actionMapper
	// initial data type parsing
	dataTypeInitActions actionMapper
)

func init() {
	constraintActions = actionMapper{
		actionMap{
			token.Id:          idAction,
			token.ImplicitId:  idAction,
			token.Underscore:  wildcardAction, // possible, i suppose, but useless // TODO: perhaps make it an error
			token.Backslash:   abstractionAction,
			token.Arrow:       productAction,
			token.IntValue:    intAction,
			token.CharValue:   charAction,
			token.StringValue: stringAction,
			token.FloatValue:  floatAction,
			token.LeftParen:   parenAction,
			token.RightParen:  closeAction,
			token.Comma:       listingAction,
		},
		constraintId,
	}

	standardActions = actionMapper{
		actionMap{
			token.Id:         idAction,
			token.ImplicitId: idAction,
			token.Underscore: wildcardAction,
			token.Backslash:  abstractionAction,
			// why is this allowed?: functions from terms to types and functions from types to types, e.g.,
			//	Dfun 0 = Int -> Int
			//	FromInt t = Int -> t
			token.Arrow:       productAction,
			token.IntValue:    intAction,
			token.CharValue:   charAction,
			token.StringValue: stringAction,
			token.FloatValue:  floatAction,
			token.LeftParen:   parenAction,
			token.RightParen:  closeAction,
			token.RightBrace:  closeAction,
		},
		standardId,
	}

	// explicit value given to implicit argument
	argPositionImplicitAction = actionMapper{
		actionMap{
			token.Id:         idAction,
			token.ImplicitId: idAction,
			token.RightBrace: closeAction,
		},
		argPosImplicitId,
	}

	typePositionImplicitAction = actionMapper{
		actionMap{
			token.Id:         idAction,
			token.ImplicitId: idAction,
			token.Underscore: wildcardAction,
			token.Colon:      labeledTypeAction,
			token.RightBrace: closeAction,
		},
		typePosImplicitId,
	}

	typingActions = actionMapper{
		actionMap{
			token.Id:          idAction,
			token.ImplicitId:  idAction,
			token.Underscore:  wildcardAction,
			token.Backslash:   abstractionAction,
			token.Arrow:       productAction,
			token.ThickArrow:  errorActionGen(IllegalConstraintPosition),
			token.IntValue:    intAction,
			token.CharValue:   charAction,
			token.StringValue: stringAction,
			token.FloatValue:  floatAction,
			token.LeftParen:   parenAction,
			token.RightParen:  closeAction,
			token.Colon:       labeledTypeAction,
			token.LeftBrace:   implicitAction,
			token.RightBrace:  closeAction,
		},
		typingId,
	}

	dataTypeFollowActions = actionMapper{
		actionMap{
			token.Id:          idAction,
			token.ImplicitId:  idAction,
			token.Backslash:   abstractionAction,
			token.Arrow:       productAction,
			token.IntValue:    intAction,
			token.CharValue:   charAction,
			token.StringValue: stringAction,
			token.FloatValue:  floatAction,
			token.LeftParen:   parenAction,
			token.RightParen:  closeAction,
		},
		dataTypeFollowId,
	}

	dataTypeInitActions = actionMapper{
		actionMap{
			token.Id:         idDataTypeAction,
			token.ImplicitId: errorActionGen(IllegalDataTypeName),
		},
		dataTypeInitId,
	}
}
