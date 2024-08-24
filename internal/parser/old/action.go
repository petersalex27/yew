package parser

import (
	"os"

	"github.com/petersalex27/yew/internal/common/stack"
	"github.com/petersalex27/yew/internal/token"
	"github.com/petersalex27/yew/internal/types"
)

type (
	// function that performs a parsing action
	actionFunc func(parser *Parser, data *actionData) (term termElem, ok bool)

	// an identifier for an actionMapper
	// NOTE: two actionMappers are "equal" if they have the same actionMapperId
	actionMapperId uint8

	// a named map from token types to actions
	actionMapper struct {
		actionMap
		id actionMapperId
	}

	// a map from token types to actions
	actionMap map[token.Type]actionFunc

	// data required for performing actions on tokens
	actionData struct {
		// current rule/action set
		m actionMapper
		// location in `tokens`
		ptr uint
		// token input to parse
		tokens     []token.Token
		discharges *stack.SaveStack[func(*Parser)]
		// supplies the next list of tokens to parse and updates any relevant state the actionData
		supplyNext func(*Parser, *actionData) bool
		// the value the parser looks for to signal the end of a sub-section
		end string
	}
)

const (
	constraintId actionMapperId = iota
	standardId
	scrutineeId
	argPosImplicitId
	typePosImplicitId
	typingId
	dataTypeFollowId
	dataTypeInitId
)

var (
	// explicitly giving an argument to an implicit parameter, e.g.,
	//	id {Int} x = x
	argPositionImplicitAction actionMapper
	// actions that cover most syntactic sections
	standardActions actionMapper
	// actions for parsing scrutinees
	scrutineeActions actionMapper
	localActions     actionMapper
	// actions that parse constraints
	constraintActions actionMapper
	// actions for parsing implicit parameters, e.g.,
	//	id : {a : Type} -> a -> a
	typePositionImplicitAction actionMapper
	// actions for parsing typings
	typingActions actionMapper
	// data type parsing after name is known
	dataTypeFollowActions actionMapper
	// initial data type parsing
	dataTypeInitActions actionMapper
)

func init() {
	constraintActions = actionMapper{
		actionMap{
			token.Id:          idAction,
			token.Underscore:  wildcardAction, // possible, i suppose, but useless // TODO: perhaps make it an error
			token.Backslash:   abstractionAction,
			token.Arrow:       productAction,
			token.IntValue:    intAction,
			token.CharValue:   charAction,
			token.StringValue: stringAction,
			token.FloatValue:  floatAction,
			token.LeftParen:   parenAction,
			token.Comma:       listingAction,
		},
		constraintId,
	}

	standardActions = actionMapper{
		actionMap{
			token.Id:         idAction,
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
			token.Comma:       listingAction,
		},
		standardId,
	}

	scrutineeActions = actionMapper{
		actionMap{
			token.Id:         scrutineeIdAction,
			token.Underscore: wildcardAction,
			// why is this allowed?: scrutinizing types
			token.Arrow:       productAction,
			token.IntValue:    intAction,
			token.CharValue:   charAction,
			token.StringValue: stringAction,
			token.FloatValue:  floatAction,
			token.LeftParen:   parenAction,
			token.Comma:       listingAction,
		},
		scrutineeId,
	}

	// explicit value given to implicit argument
	argPositionImplicitAction = actionMapper{
		actionMap{
			token.Id:         idAction,
		},
		argPosImplicitId,
	}

	typePositionImplicitAction = actionMapper{
		actionMap{
			token.Id:         idAction,
			token.Underscore: wildcardAction,
			token.Colon:      labeledTypeAction,
			token.Erase:      modalityAction,
			token.Once:       modalityAction,
		},
		typePosImplicitId,
	}

	typingActions = actionMapper{
		actionMap{
			token.Id:          idAction,
			token.Equal:       toIdAction,
			token.Underscore:  wildcardAction,
			token.Backslash:   abstractionAction,
			token.Arrow:       productAction,
			token.ThickArrow:  errorActionGen(IllegalConstraintPosition),
			token.IntValue:    intAction,
			token.CharValue:   charAction,
			token.StringValue: stringAction,
			token.FloatValue:  floatAction,
			token.LeftParen:   parenAction,
			token.Colon:       labeledTypeAction,
			token.LeftBrace:   implicitAction,
		},
		typingId,
	}

	dataTypeFollowActions = actionMapper{
		actionMap{
			token.Id:          idAction,
			token.Backslash:   abstractionAction,
			token.Arrow:       productAction,
			token.IntValue:    intAction,
			token.CharValue:   charAction,
			token.StringValue: stringAction,
			token.FloatValue:  floatAction,
			token.LeftParen:   parenAction,
		},
		dataTypeFollowId,
	}

	dataTypeInitActions = actionMapper{
		actionMap{
			token.Id:         idDataTypeAction,
		},
		dataTypeInitId,
	}
}

var modalityMap = map[token.Type]types.Multiplicity{
	token.Erase: types.Erase,
	token.Once:  types.Once,
}

func (parser *Parser) setMultiplicity(mode token.Token, v types.Variable) (_ types.Variable, ok bool) {
	var mult types.Multiplicity
	if mult, ok = modalityMap[mode.Type]; !ok {
		// this should be impossible ...
		parser.errorOnToken("yikes ... what mode is this ... ?", mode)
		return v, false
	}
	v = v.SetMultiplicity(mult)
	return v, true
}

func modalityAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	panic("TODO: implement") // TODO: implement
	/*
		var mode token.Token
		if mode, ok = parser.nextToken(data); !ok {
			return
		}

		if ok = parser.allowModality; !ok {
			parser.errorOnToken(IllegalModalityLocation, mode)
			return
		}

		parser.allowModality = false // only allow a single modality

		if term, ok = parser.process(-1, data); !ok {
			return
		}

		var typing Typing
		if typing, ok = term.Term.(Typing); !ok {
			parser.illegalModalityError(mode, term.Term)
			return term, false
		}

		term.Term, ok = parser.setMultiplicity(mode, typing)
		return
	*/
}

// constructs a new actionData and returns a pointer to it
func newActionData(m actionMapper, toks []token.Token) *actionData {
	return &actionData{
		m:      m,
		ptr:    0,
		tokens: toks,
	}
}

// constructs a new actionData and returns a pointer to it
func newActionDataWithDischarger(m actionMapper, toks []token.Token) *actionData {
	return &actionData{
		m:          m,
		ptr:        0,
		tokens:     toks,
		discharges: stack.NewSaveStack[func(*Parser)](4),
	}
}

// constructs a new actionData with its internal token pointer offset by `offset` and returns a
// pointer to it
//
// panics if offset is greater than the length of the supplied tokens--note, it's okay to give an
// offset equal to the length of the tokens (this signifies that the actionData has reached the end
// of its input)
func newOffsetActionData(m actionMapper, toks []token.Token, offset uint) *actionData {
	if offset > uint(len(toks)) {
		panic("bug: ptr offset is greater than the length of the supplied tokens")
	}

	data := newActionData(m, toks)
	data.ptr = offset
	return data
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
		parser.errorEOI(tok)
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

// performs an action based on a terminal term (i.e., a token)
//
// actions are decided based on `data.m`
//
// on success, returns the parsed term and true; otherwise, returns ?, false
func (parser *Parser) actOnTerminal(data *actionData) (term termElem, ok bool) {
	if !parser.termMemory.Empty() {
		term, _ = parser.termMemory.Pop()
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

// creates the term info for a function arrow
func arrowInfo() termInfo {
	return termInfo{1, true, 2, true}
}

func (parser *Parser) idActionHelper(tok token.Token) (term termElem, ok bool, found bool) {
	// attempt to find declaration, converting it into a term if found
	if term, found = parser.findDeclAsTerm(tok); found {
		return term, true, true
	}

	return parser.makeVarTermElem(tok), true, false
}

// converts any token into an identifier and attempts to find a declared version of it
func toIdAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	var tok token.Token
	if tok, ok = parser.nextToken(data); !ok {
		return
	}

	res := token.Id.MakeValued(tok.Value)
	res.Start = tok.Start
	res.End = tok.End
	var found bool
	term, ok, found = parser.idActionHelper(res)
	if ok && !found {
		parser.errorOnToken(UndefinedName, tok)
		ok = false
	}
	return
}

// produces a function type constructor
func productAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	var tok token.Token
	if tok, ok = parser.nextToken(data); !ok {
		return
	}

	//f := makeFunctionType(tok)
	term = termElem{functionType, nil, arrowInfo(), tok.Start, tok.End}
	return
}

// produces a listing constructor
func listingAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	var comma token.Token
	if comma, ok = parser.nextToken(data); !ok {
		return
	}

	//listing := Listing{[]Term{}, comma.Start, comma.End}
	term = termElem{listingType, nil, termInfo{0, true, 2, true}, comma.Start, comma.End}
	return term, true
}

// generates an action function that reports an error based on the token read
//
// `msg` is the message to report
//
// the generated function will always
//
//	return termElem{}, false
func errorActionGen(msg string) func(parser *Parser, data *actionData) (term termElem, ok bool) {
	return func(parser *Parser, data *actionData) (termElem, bool) {
		tok, _ := data.peek()
		parser.errorOnToken(msg, tok)
		return termElem{}, false
	}
}

// data type initial name parsing
func idDataTypeAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	var tok token.Token
	if tok, ok = parser.nextToken(data); !ok {
		return
	}

	// see if already declared
	if _, found := parser.declarations.Find(tok); found {
		// error: redeclaration
		parser.errorOn(IllegalRedeclaration, tok)
		ok = false
		return
	}

	// create name
	name := types.Var(tok)
	term = termElem{datatypeType, name, termInfo{}, tok.Start, tok.End}
	// TODO: declare
	return term, true
}

// makes a var term element from a token
func (parser *Parser) makeVarTermElem(tok token.Token) termElem {
	v := types.Var(tok)
	term := termElem{identType, v, termInfo{}, tok.Start, tok.End}
	return term
}

// produces an identifier either from a previously declared id or a new, free identifier
//
// fails when `data` is out of input
func idAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	var tok token.Token
	if tok, ok = parser.nextToken(data); !ok {
		return term, false
	}

	term, ok, _ = parser.idActionHelper(tok)
	return term, ok
}

// produces an identifier either from a previously declared id or a new, free identifier
//
// fails when `data` is out of input
func scrutineeIdAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	var tok token.Token
	if tok, ok = parser.nextToken(data); !ok {
		return term, false
	}

	//var found bool
	term, ok, _ = parser.idActionHelper(tok)
	return term, ok
}

var _Int named = named{"Int", 0, 0}
var _Float named = named{"Float", 0, 0}
var _Char named = named{"Char", 0, 0}
var _String named = named{"String", 0, 0}

func literalAction[T types.Literal](parser *Parser, data *actionData, nt NodeType, typeName named, msg string) (term termElem, ok bool) {
	var tok token.Token
	if tok, ok = parser.nextToken(data); !ok {
		return
	}

	var v T
	var lit types.Literal
	if lit, ok = v.Parse(tok); !ok {
		// TODO: could not parse, but why?
		parser.errorOnToken(msg, tok)
		return term, false
	}
	v = lit.(T)

	ty, _, found := parser.env.Get(typeName)
	if ok = found; !ok {
		parser.errorOnToken(NoAppropriateType, tok)
		return
	}
	typ, yes := ty.(types.Type)
	if ok = yes; !ok {
		parser.errorOnToken(NoAppropriateType, tok)
		return
	}

	var tmp types.Term = v
	if ok = types.SetKind(&tmp, typ); !ok {
		parser.transferEnvErrors()
		return
	}
	v = tmp.(T)
	return termElem{nt, v, termInfo{}, tok.Start, tok.End}, true
}

// produces an integer literal
//
// fails when `data` is out of input
func intAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	return literalAction[*types.IntConst](parser, data, intConstType, _Int, CouldNotParseInt)
}

// produces an character literal
//
// fails when `data` is out of input
func charAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	return literalAction[*types.CharConst](parser, data, charConstType, _Char, CouldNotParseChar)
}

// produces a string literal
//
// fails when `data` is out of input
func stringAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	return literalAction[*types.StringConst](parser, data, stringConstType, _String, CouldNotParseString)
}

// produces an float literal
//
// fails when `data` is out of input
func floatAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	return literalAction[*types.FloatConst](parser, data, floatConstType, _Float, CouldNotParseFloat)
}

func wildcardAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	var __ token.Token
	if __, ok = parser.nextToken(data); !ok {
		return term, false
	}
	wildcard := types.Wildcard()
	wildcard.Start, wildcard.End = __.Pos()
	term = termElem{identType, wildcard, termInfo{}, wildcard.Start, wildcard.End}
	return term, true
}

// creates typing prototype--still requires a type
func createTyping(colon token.Token) termElem {
	info := termInfo{0, true, 2, true}
	return termElem{typingType, nil, info, colon.Start, colon.End}
}

// typing, e.g.,
//
//	x : A
func labeledTypeAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	// read past ':' token
	var colon token.Token
	if colon, ok = parser.nextToken(data); ok {
		term, ok = createTyping(colon), true
	}
	return
}

// advances past token, reports error and returns false when end marker isn't found, otherwise just
// returns true
func (parser *Parser) readEnd(data *actionData) (tok token.Token, ok bool) {
	tok, ok = parser.nextToken(data)
	if !ok {
		return tok, false
	}

	if ok = tok.Value == data.end; !ok {
		parser.errorOnToken(expectedSyntax(data.end), tok)
	}
	return tok, ok
}

type extensionUpdater struct {
	end   string
	m     func(parser *Parser, data *actionData) (actionMapper, bool)
	build func(parser *Parser, termPrev, termNew termElem, start, end int) (termElem, bool)
}

func _makeEnds_m_(*Parser, *actionData) (actionMapper, bool) { return standardActions, true }

func _makeEnds_apply_(parser *Parser, termPrev, termNew termElem, start, end int) (termElem, bool) {
	// TODO: cannot set start and end here :(
	return parser.reduce(termPrev, termNew)
}

// makes default ends
//
// fn is underlying function results are applied to
func makeEnds(ends ...string) []extensionUpdater {
	out := make([]extensionUpdater, len(ends))
	for i, end := range ends {
		out[i] = extensionUpdater{end, _makeEnds_m_, _makeEnds_apply_}
	}
	return out
}

func extensionAction(first string, fn termElem, ends ...extensionUpdater) func(parser *Parser, data *actionData) (term termElem, ok bool) {
	return func(parser *Parser, data *actionData) (term termElem, ok bool) {
		// this shadows captured `fn` so the `fn` captured by the closure isn't overwritten. It actually
		// doesn't matter logic-wise, but for printing debug info it does matter
		fn := fn
		old := data.m
		oldEnd := data.end

		defer func() {
			data.end = oldEnd
			data.m = old
		}()

		var open token.Token
		if open, ok = parser.nextToken(data); !ok {
			return
		}
		if open.Value != first {
			parser.errorOnToken(expectedSyntax(first), open)
			return term, false
		}

		for _, end := range ends {
			// update map according to updater
			data.m, ok = end.m(parser, data)
			if !ok {
				return
			}

			data.end = end.end

			if term, ok = parser.process(-1, data); !ok {
				return
			}

			var endTok token.Token
			if endTok, ok = parser.readEnd(data); !ok {
				return term, false
			}
			before := fn
			fn, ok = end.build(parser, fn, term, open.Start, endTok.End)
			if !ok {
				return
			}
			debug_log_reduce(os.Stderr, before, term, fn)
		}
		return fn, ok
	}
}

func implicitEndM(parser *Parser, data *actionData) (actionMapper, bool) {
	switch data.m.id {
	case typingId:
		return typePositionImplicitAction, true
	case standardId:
		return argPositionImplicitAction, true
	default:
		var tok token.Token
		var ok bool
		if tok, ok = parser.peek(data); !ok {
			return actionMapper{}, false
		}
		parser.errorOnToken(NoGrammarExtensionFound, tok)
		return actionMapper{}, false
	}
}

func implicitBuilder(parser *Parser, _ termElem, term termElem, start, end int) (termElem, bool) {
	var imp types.Term
	var info termInfo
	switch term.NodeType {
	case identType:
		if !parser.parsingTypeSig {
			// exit switch, treat just as a term since implicit braces are being used in a function
			// pattern/application
			break
		}

		// assume this is a term of 'Type'
		ok := types.SetKind(&term.Term, types.Type0)
		id := term.Term.(types.Variable)
		if !ok {
			parser.errorOn(CouldNotType, id)
			return termElem{}, false
		}
		imp = id
		// create new entry
		info = term.termInfo
	case typingType:
		imp = term.Term
	default:
		// TODO: implicit argument, ig?
		info = term.termInfo
		imp = term.Term
	}

	info = abstractInfo(info)
	term = termElem{implicitType, imp, info, start, end}
	return term, true
}

var implicitAction = extensionAction("{", termElem{}, extensionUpdater{"}", implicitEndM, implicitBuilder})

var parenAction_ = extensionAction("(", termElem{}, extensionUpdater{")", parenEndM, parenBuilder})

func parenAction(parser *Parser, data *actionData) (term termElem, ok bool) {
	if data.discharges == nil || data.discharges.Empty() {
		return parenAction_(parser, data)
	}
	// need to track discharges
	data.discharges.Save()
	term, ok = parenAction_(parser, data)
	discharge, stat := data.discharges.Pop()
	for ; stat.IsOk(); discharge, stat = data.discharges.Pop() {
		discharge(parser)
	}
	return
}

func parenEndM(_ *Parser, data *actionData) (actionMapper, bool) {
	return data.m, true
}

func parenBuilder(parser *Parser, _ termElem, term termElem, start, end int) (termElem, bool) {
	return standardParenCloser(parser, term, start, end)
}

// when a term is enclosed by parens, it's information is lost and becomes a more abstract function
func abstractInfo(info termInfo) termInfo {
	if info.Arity() > 0 && info.bp > 0 {
		// update binding power and associativity for enclosed functions
		//	- they lose whatever the inner bp and assoc is
		info.bp = 10
		info.rAssoc = false
	}
	info.infixed = false
	return info
}
