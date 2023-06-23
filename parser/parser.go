package parsing

import (
	sync "sync"
	err "yew/error"
	scan "yew/lex"
	ast "yew/parser/ast"
	bp "yew/parser/bp"
	nodetype "yew/parser/node-type"
	. "yew/parser/parser"
	symbol "yew/symbol"
	types "yew/type"
)

var DefaultNameSpaceId ast.Id = ast.MakeId(scan.UnderscoreIdToken)

func pushId(p *Parser) bool {
	// check for id
	validPackageDeclaration := p.Current.GetType() == scan.ID
	if !validPackageDeclaration {
		expectedID(p.Next, p.Input).Print()
		return false
	}

	p.Stack.Push(ast.MakeId(p.Current.(scan.IdToken)))
	return true
}

func printSyntaxError(p *Parser, token scan.Token) {
	unexpectedToken(token, p.Input).Print()
}

func error(p *Parser, token scan.Token) {
	printSyntaxError(p, token)
}

func printError(p *Parser, errToken scan.Token) {
	error(p, errToken)
}

func generateSyntaxError(message string) func(scan.Token, scan.InputStream) err.Error {
	return func(t scan.Token, i scan.InputStream) err.Error {
		loc := t.GetLocation()
		return err.CompileMessage(
			message, err.ERROR, err.SYNTAX, (i).GetPath(), loc.GetLine(), loc.GetChar(),
			t.GetSourceIndex(), (i).GetSource()).(err.Error)
	}
}

var expectedID = generateSyntaxError("expected identifier")
var expectedEndOfStatement = generateSyntaxError("expected end of statement")
var expectedASSIGN = generateSyntaxError("expected assignment")
var expectedLCURL = generateSyntaxError("expected left curly brace")
var unmatchedPAREN = generateSyntaxError("unmatched left-paren")
var unmatchedBRACK = generateSyntaxError("unmatched left-bracket")
var unexpectedToken = generateSyntaxError("unexpected token")

func isPrefixTypeToken(t scan.TokenType) bool {
	return t == scan.INT ||
		t == scan.BOOL ||
		t == scan.CHAR ||
		t == scan.STRING ||
		t == scan.FLOAT ||
		t == scan.LPAREN ||
		t == scan.LBRACK ||
		t == scan.ID
}

func parseTypePrefix(p *Parser) (types.Types, bool) {
	p.Stack.Push(ast.StackMarker{})

	i := 0
	for {
		current := p.Current

		switch current.GetType() {
		case scan.INT:
			p.Stack.Push(ast.MakeType(types.Int{}))
		case scan.BOOL:
			p.Stack.Push(ast.MakeType(types.Bool{}))
		case scan.CHAR:
			p.Stack.Push(ast.MakeType(types.Char{}))
		case scan.STRING:
			p.Stack.Push(ast.MakeType(types.Array{ElemType: types.Char{}}))
		case scan.FLOAT:
			p.Stack.Push(ast.MakeType(types.Float{}))
		case scan.LPAREN:
			ty, ok := parseTypeAnnotation(p, bp.None)
			if scan.RPAREN != p.Next.GetType() {
				if ok {
					unmatchedPAREN(current, p.Input).Print()
				}
				return types.Error{}, false
			}
			p.Advance()
			p.Stack.Push(ast.MakeType(ty))
		case scan.LBRACK:
			ty, ok := parseTypeAnnotation(p, bp.None)
			if scan.RBRACK != p.Next.GetType() {
				if ok {
					unmatchedBRACK(current, p.Input).Print()
				}
				return types.Error{}, false
			}
			p.Advance()
			p.Stack.Push(ast.MakeType(types.Array{ElemType: ty}))
		case scan.ID:
			//s := symbol.MakeSymbol(current.(scan.IdToken))
			//sym, added := p.Table.GetElseAdd(s)
			/*if !added {
				return types.Error{}, false
			}//*/
			/*if added {
				newType := types.GetNewTau()
				sym = sym.SetType(newType)
				sym = p.Table.Update(sym)
			}//*/
			//return sym.GetType(), true
			str := current.(scan.IdToken).ToString()
			p.Stack.Push(ast.MakeType(types.Tau(str)))
		default:
			if i == 0 {
				unexpectedToken(current, p.Input).Print()
				return types.Error{}, false
			}
		}

		i++
		if i > 1 {
			second := p.Stack.Pop().(ast.Type)
			var ty types.Application
			if i == 2 {
				first := p.Stack.Pop().(ast.Type)
				ty = types.Application{first.GetType(), second.GetType()}
			} else {
				first := p.Stack.Pop().(ast.Type).GetType().(types.Application)
				ty = append(first, second.GetType())
			}
			p.Stack.Push(ast.MakeType(ty))
		}
		// whenever this point is reached there should only be one type on the stack
		// 	from this function call!
		if isPrefixTypeToken(p.Next.GetType()) {
			p.Advance()
		} else {
			break
		}
	}

	ty := p.Stack.Pop().(ast.Type).GetType()
	if p.Stack.GetTopNodeType() != nodetype.STACK_MARKER {
		err.PrintBug()
		panic("")
	}
	p.Stack.Pop() // remove stack marker

	return ty, true
}

func parseTypeArrow(p *Parser, left types.Types) (types.Types, bool) {
	bp := getBindingPower(p.Current)
	p.Advance()
	// along with scan.ARROW being the only thing with its binding power,
	// 	the "- 1" makes this type-operation right-to-left
	right, ok := parseTypeAnnotation(p, bp-1)
	if !ok {
		return right, ok
	}
	return types.Function{Domain: left, Codomain: right}, true
}

func parseTypeList(p *Parser, left types.Types) (types.Types, bool) {
	listBp := getBindingPower(p.Current)
	ts := make([]types.Types, 0, 1)
	ts = append(ts, left)
	for {
		p.Advance() // discard comma and put next token into current

		// check for end of parse (current is not a type token)
		if !isTypePrefixTokenType(p.Current.GetType()) {
			break // trailing comma, end parse and return tuple type
		}

		// parse element type
		elem, ok := parseTypeAnnotation(p, listBp)
		if !ok {
			return elem, false
		}

		ts = append(ts, elem)

		if p.Next.GetType() != scan.COMMA {
			break // end parse if next token is not comma
		}

		p.Advance() // move comma to current
	}
	return types.Tuple(ts), true
}

func parseTypeRight(p *Parser, left types.Types, bp bp.BindingPower) (types.Types, bool) {
	// need to use next token since the default option should not grab any tokens
	next := p.Next.GetType()
	nextBp := getBindingPower(p.Next)
	parseNext := nextBp > bp && (next == scan.ARROW || next == scan.COMMA)

	if !parseNext {
		return left, true
	}

	p.Advance()

	switch next {
	case scan.ARROW:
		return parseTypeArrow(p, left)
	case scan.COMMA:
		return parseTypeList(p, left)
	default:
		return left, true
	}
}

func parseDataType(p *Parser) bool {
	// loop until no more `|`
	for {
		ty, ok := parseTypeAnnotation(p, 0)
		if !ok {
			return false
		}
		p.Stack.Push(ast.MakeType(ty))
		if !(ast.Type{}.AddConstructor(p)) {
			return false
		}
		if p.Next.GetType() == scan.NEW_LINE {
			p.Advance()
		}

		if p.Next.GetType() != scan.BAR {
			return true // end parse
		}

		// else `|` is next, move past it
		p.Advance()
		p.Advance()
	}
}

func parseTypeAnnotation(p *Parser, bp bp.BindingPower) (types.Types, bool) {
	var left types.Types
	var ok bool
	left, ok = parseTypePrefix(p)
	if !ok {
		return left, false
	}
	left, ok = parseTypeRight(p, left, bp)
	if !ok {
		return left, false
	}
	return left, true
}

func isTypePrefixTokenType(t scan.TokenType) bool {
	switch t {
	case scan.ID:
		fallthrough
	case scan.LPAREN:
		fallthrough
	case scan.LBRACK:
		fallthrough
	case scan.INT:
		fallthrough
	case scan.BOOL:
		fallthrough
	case scan.FLOAT:
		fallthrough
	case scan.STRING:
		fallthrough
	case scan.CHAR:
		return true
	}
	return false
}

func initialParseTypeAnnotation(p *Parser) bool {
	if p.Current.GetType() == scan.COLON_COLON {
		p.Stack.Push(ast.EmptyExpression{})
		res := typeAnnotation(p, p.Current)
		if res {
			p.Advance()
		}
		return res
	}

	annot := ast.MakeTypeAnnotation(ast.EmptyExpression{}, types.GetNewTau())
	p.Stack.Push(annot)
	return true
}

func getBindingPower(t scan.Token) bp.BindingPower {
	return parseTable[t.GetType()].bp
}

func parseExpressionExpectingInitial[T Ast](p *Parser, bp bp.BindingPower, initialParse func(p *Parser) (T, bool)) (T, bool) {
	res, ok := initialParse(p)

	for getBindingPower(p.Next) > bp && ok {
		p.Advance()
		ok = parseTable[p.Current.GetType()].infixRule(p, p.Current)
	}
	return res, ok
}

func parseExpression(p *Parser, bp bp.BindingPower) bool {
	//fmt.Fprintf(os.Stderr, "%s\n", p.Current.GetType().ToString())
	elem, found := parseTable[p.Current.GetType()]
	if !found {
		err.PrintBug()
		panic("")
	}
	ok := elem.prefixRule(p, p.Current)
	for getBindingPower(p.Next) > bp && ok {
		p.Advance()
		ok = parseTable[p.Current.GetType()].infixRule(p, p.Current)
	}
	return ok
}

// for decls
func parseId(p *Parser) bool {
	id := p.Current
	if scan.ID != id.GetType() {
		expectedID(p.Current, p.Input).Print()
		return false
	}

	p.Advance()

	p.Stack.Push(ast.MakeId(id.(scan.IdToken)))
	return true
}

func parseDeclaration(p *Parser) bool {
	// whether it finds a type or not, a type annotation is pushed; if it doesn't
	// find one, a tau type is pushed
	return parseId(p) &&
		initialParseTypeAnnotation(p) &&
		ast.Declaration{}.Make(p)
}

func parseInit(p *Parser) bool {
	if p.Current.GetType() != scan.EQUALS {
		expectedASSIGN(p.Current, p.Input).Print()
		return false
	}
	p.Advance()
	return parseExpression(p, bp.None) //&& ast.Assignment{}.Make(p)
}

func letParse(p *Parser) bool {
	// definition ::= `let` ID `::` typeAnotation = expression
	return parseDeclaration(p) &&
		parseInit(p) &&
		ast.Definition{}.Make(p)
}

// for expressions
func identifier(p *Parser, token scan.Token) bool {
	idToken := token.(scan.IdToken)
	var sym symbol.Symbolic = symbol.MakeSymbol(idToken)
	sym, _ = p.Table.GetElseAdd(sym)
	//sym.UpdateUse(idToken.GetLocation())
	p.Stack.Push(ast.MakeId(sym.GetIdToken()))
	return true
}
func value_(p *Parser, token scan.Token) bool {
	val := token.(scan.ValueToken).GetValue()
	p.Stack.Push(ast.MakeValue(val))
	return true
}
func parseError(p *Parser, token scan.Token) bool {
	unexpectedToken(token, p.Input).Print()
	return false
}

func bracket(p *Parser, token scan.Token) bool {
	err.PrintBug()
	panic("TODO")
}

func paren(p *Parser, token scan.Token) bool {
	err.PrintBug()
	panic("TODO")
}

func block(p *Parser, token scan.Token) bool {
	err.PrintBug()
	panic("TODO")
}

func sequence(p *Parser, delimiter scan.TokenType, action func(*Parser, scan.Token) bool) bool {
	for p.Current.GetType() == delimiter {
		p.Advance()
		if !action(p, p.Current) {
			return false
		}
	}
	return false
}

func pushOperation(p *Parser, tokenType scan.TokenType) {
	switch tokenType {
	case scan.PLUS:
		p.Stack.Push(ast.ADD)
	case scan.STAR:
		p.Stack.Push(ast.MULTIPLY)
	case scan.SLASH:
		p.Stack.Push(ast.DIVIDE)
	case scan.MINUS:
		p.Stack.Push(ast.SUBTRACT)
	case scan.MOD:
		p.Stack.Push(ast.MOD)
	case scan.HAT:
		p.Stack.Push(ast.POWER)
	case scan.PLUS_PLUS:
		p.Stack.Push(ast.APPEND)
	case scan.COLON:
		p.Stack.Push(ast.CONSTRUCT)
	case scan.ARROW:
		p.Stack.Push(ast.MAPS_TO)
	case scan.AMPER_AMPER:
		p.Stack.Push(ast.AND)
	case scan.BAR_BAR:
		p.Stack.Push(ast.OR)
	case scan.GREAT:
		p.Stack.Push(ast.GREAT)
	case scan.LESS:
		p.Stack.Push(ast.LESS)
	case scan.GREAT_EQUALS:
		p.Stack.Push(ast.GREAT_EQUALS)
	case scan.LESS_EQUALS:
		p.Stack.Push(ast.LESS_EQUALS)
	case scan.EQUALS_EQUALS:
		p.Stack.Push(ast.EQUALS)
	case scan.BANG_EQUALS:
		p.Stack.Push(ast.NOT_EQUALS)
	case scan.BANG:
		p.Stack.Push(ast.NOT)
	case scan.BANG_POSTFIX__:
		p.Stack.Push(ast.FACTORIAL)
	case scan.PLUS_PREFIX__:
		p.Stack.Push(ast.POSITIVE)
	case scan.MINUS_PREFIX__:
		p.Stack.Push(ast.NEGATIVE)
	default:
		err.PrintBug()
		panic("")
	}
}

// returns proper token for operation now that the context (i.e., part of arithmetic expression)
// is known
func fixArithmeticToken(token scan.Token) scan.Token {
	ty := token.GetType()
	if ty == scan.BANG {
		return token.(scan.OtherToken).ChangeTokenType(scan.BANG_POSTFIX__)
	}
	return token
}

func fixUnaryToken(token scan.Token) scan.Token {
	ty := token.GetType()
	if ty == scan.PLUS {
		return token.(scan.OtherToken).ChangeTokenType(scan.PLUS_PREFIX__)
	} else if ty == scan.MINUS {
		return token.(scan.OtherToken).ChangeTokenType(scan.MINUS_PREFIX__)
	}
	return token
}

func prefix(p *Parser, token scan.Token) bool {
	p.Advance() // move past operation

	token = fixUnaryToken(token)
	operationsBindingPower := getBindingPower(token)
	tokenType := token.GetType()
	pushOperation(p, tokenType)

	if !parseExpression(p, operationsBindingPower) {
		return false
	}

	return ast.UnaryOperation{}.Make(p)
}

func binary(p *Parser, token scan.Token) bool {
	p.Advance() // move past operation

	token = fixArithmeticToken(token)
	operationsBindingPower := getBindingPower(token)
	tokenType := token.GetType()
	pushOperation(p, tokenType)

	if tokenType == scan.BANG_POSTFIX__ { // factorial
		return ast.PostfixOperation{}.Make(p)
	}

	if !parseExpression(p, operationsBindingPower) {
		return false
	}

	return ast.BinaryOperation{}.Make(p)
}

var postfixOp = binary

func lambda(p *Parser, token scan.Token) bool {
	p.Advance()
	// do not use binding power of -> here, else -> would become left-to-right; but, it's
	// right-to-left
	if !parseExpression(p, bp.None) {
		return false
	}
	return ast.Lambda{}.Make(p)
}

func createBinder(p *Parser, token scan.Token) bool {
	(p.Table).NewScope()

	action := func(par *Parser, tok scan.Token) bool {
		if tok.GetType() != scan.ID {
			return true // end of action
		}
		if !parseDeclaration(par) {
			return false
		}
		dec := par.Stack.Peek().(ast.Declaration)
		par.Stack.Push(ast.Id(dec))
		if !(ast.Binder{}.Make(p)) {
			return false
		}
		par.Advance()
		return true
	}

	if p.Next.GetType() != scan.ID {
		parseError(p, p.Next)
		return false
	}
	p.Advance()
	ok := true
	for {
		ok = action(p, p.Current)
		if !ok {
			break
		}

		ok = sequence(p, scan.COMMA, action)
		if !ok {
			break
		}
		if p.Next.GetType() != scan.ARROW {
			parseError(p, p.Next)
			ok = false
			break
		}

		p.Advance()
		ok = lambda(p, p.Current)
		break
	}
	p.Table.RemoveScope()
	return ok
}

func lists(p *Parser, token scan.Token) bool {
	err.PrintBug()
	panic("TODO")
}

// Type-Annotation ::= Expression
var typeAnnotationRule = nodetype.NodeRule{
	Production: nodetype.TYPE_ANNOTATION,
	Expression: []nodetype.NodeType{
		nodetype.EXPRESSION,
	},
}

func typeAnnotation(p *Parser, token scan.Token) bool {
	p.Advance() // move past `::`
	ty, ok := parseTypeAnnotation(p, bp.None)
	if !ok {
		error(p, p.Current)
		return false
	}

	valid, e := p.Stack.Validate(typeAnnotationRule)
	if !valid {
		e.Print()
		return false
	}

	expr := p.Stack.Pop().(ast.Expression)
	annot := ast.MakeTypeAnnotation(expr, ty)
	p.Stack.Push(annot)
	return true
}

func ignore(*Parser, scan.Token) bool {
	return true
}

func pattern(p *Parser, token scan.Token) bool {
	err.PrintBug()
	panic("TODO")
}

func annotation(p *Parser, token scan.Token) bool {
	//id := token.(scan.AnnotationToken).ToString()
	err.PrintBug()
	panic("TODO")
}

func comma(p *Parser, token scan.Token) bool {
	err.PrintBug()
	panic("TODO")
}

func parseBug(*Parser, scan.Token) bool {
	err.PrintBug()
	panic("")
}

var parseTableMutex sync.Mutex
var isParseTableInitialized = false

func compose(p *Parser, token scan.Token) bool {
	// use the binding power just below composition; this makes
	// 	composition right-to-left
	p.Advance()
	ok :=
		parseExpression(p, bp.Compose - 1) &&
			ast.Application{}.Make(p)
	return ok
}

func isApplicable(token scan.Token) bool {
	rule, found := parseTable[token.GetType()]
	if !found {
		err.PrintBug()
		panic("")
	}

	return rule.callRule != nil
}

func parseApplicationParen(p *Parser, parenToken scan.Token) bool {
	p.Advance() // move past paren

	ok := true
	for ok {
		ok = parseExpression(p, 0)
		if !ok {
			break
		}

		doBreak := false
		switch p.Next.GetType() {
		case scan.COMMA:
			fallthrough // TODO
		case scan.SEMI_COLON:
			panic("TODO") // TODO
		case scan.RPAREN:
			fallthrough
		default:
			doBreak = true
		}

		if doBreak || !ok {
			break
		}
	}

	if ok && p.Next.GetType() != scan.RPAREN {
		ok = false
		unmatchedPAREN(parenToken, p.Input).Print()
	}
	return ok
}

func parseApplication(p *Parser, _ scan.Token) bool {
	if p.Stack.Peek().GetNodeType() != nodetype.IDENTIFIER {
		// check that top of the stack is an id
		unexpectedToken(p.Input.GetTokenAtOffset(-3), p.Input).Print()
		return false
	}
	if !isApplicable(p.Current) {
		unexpectedToken(p.Current, p.Input).Print()
		return false
	}

	var ok bool = true
	for ok {
		currentType := p.Current.GetType()
		rule := parseTable[currentType]
		if currentType == scan.COLON_COLON {
			ok = rule.callRule(p, p.Current) // first, parse type, then ...
			break                            // ... end application parse
		}

		ok = rule.callRule(p, p.Current)
		if currentType != scan.DOT { // if dot, then has already been applied
			ok = ok && ast.Application{}.Make(p)
		}

		if !isApplicable(p.Next) {
			break
		}
		// else advance and continue
		p.Advance()
	}
	return ok
}

// initializes parse table (safe to call from multiple threads)
func InitParseTable() {
	// this function is needed to avoid circular initializations
	parseTableMutex.Lock()
	if !isParseTableInitialized {
		// binding power should always be binding power for infix or postfix
		// consequently, there should never be a token that can be both infix and postfix
		// (prefix and postfix is fine though)
		tmp := map[scan.TokenType]parseTableElement{
			scan.ID:        {bp.Applicative, identifier, parseApplication, identifier},
			scan.VALUE:     {bp.Applicative, value_, parseApplication, value_},
			scan.LBRACK:    {bp.Applicative, bracket, parseApplication, bracket},
			scan.RBRACK:    {bp.None, parseError, parseError, bracket},
			scan.LPAREN:    {bp.Group, paren, parseApplicationParen, paren},
			scan.RPAREN:    {bp.None, parseError, parseError, nil},
			scan.LCURL:     {bp.None, block, parseError, block},
			scan.RCURL:     {bp.None, parseError, parseError, nil},
			scan.STRING:    {bp.None, parseError, parseError, nil},
			scan.CHAR:      {bp.None, parseError, parseError, nil},
			scan.BOOL:      {bp.None, parseError, parseError, nil},
			scan.FLOAT:     {bp.None, parseError, parseError, nil},
			scan.INT:       {bp.None, parseError, parseError, nil},
			scan.TYPE:      {bp.None, parseError, parseError, nil},
			scan.CLASS:     {bp.None, parseError, parseError, nil},
			scan.WHERE:     {bp.None, parseError, parseError, nil},
			scan.LET:       {bp.None, parseError, parseError, nil},
			scan.CONST:     {bp.None, parseError, parseError, nil},
			scan.MUT:       {bp.None, parseError, parseError, nil},
			scan.LAZY:      {bp.None, parseError, parseError, nil},
			scan.PLUS:      {bp.Additive, prefix, binary, operation},
			scan.PLUS_PLUS: {bp.Additive, parseError, binary, operation},
			scan.MINUS:     {bp.Additive, prefix, binary, operation},
			scan.STAR:      {bp.Multiplicative, parseError, binary, operation},
			scan.SLASH:     {bp.Multiplicative, parseError, binary, operation},
			scan.MOD:       {bp.Multiplicative, parseError, binary, operation},
			scan.HAT:       {bp.Power, parseError, binary, operation},
			scan.EQUALS:    {bp.None, parseError, parseError, nil},
			//scan.PLUS_EQUALS:    {bp.None, parseError, parseError, nil},
			//scan.MINUS_EQUALS:   {bp.None, parseError, parseError, nil},
			//scan.STAR_EQUALS:    {bp.None, parseError, parseError, nil},
			//scan.SLASH_EQUALS:   {bp.None, parseError, parseError, nil},
			scan.COLON:             {bp.Additive, parseError, binary, operation},
			scan.COLON_COLON:       {bp.ExpressionAnotation, parseError, typeAnnotation, typeAnnotation},
			scan.COLON_COLON_EQUAL: {bp.None, parseError, parseError, nil},
			scan.SEMI_COLON:        {bp.None, ignore, parseError, nil},
			scan.QUESTION:          {bp.PatternMatch, parseError, pattern, nil},
			scan.COMMA:             {bp.None, parseError, parseError, comma},
			scan.BANG:              {bp.Postfix, prefix, postfixOp, operation},
			scan.BANG_EQUALS:       {bp.Equitable, parseError, binary, operation},
			scan.EQUALS_EQUALS:     {bp.Equitable, parseError, binary, operation},
			scan.BAR:               {bp.None, parseError, parseError, nil},
			scan.AMPER_AMPER:       {bp.Conjunctive, parseError, binary, operation},
			scan.BAR_BAR:           {bp.Disjunctive, parseError, binary, operation},
			scan.ARROW:             {bp.Mapping, parseError, lambda, nil},
			scan.FAT_ARROW:         {bp.None, parseError, parseError, nil},
			scan.DOT:               {bp.Compose, parseError, compose, compose},
			scan.DOT_DOT:           {bp.None, parseError, parseError, nil},
			scan.GREAT:             {bp.Ordered, parseError, binary, operation},
			scan.LESS:              {bp.Ordered, parseError, binary, operation},
			scan.GREAT_EQUALS:      {bp.Ordered, parseError, binary, operation},
			scan.LESS_EQUALS:       {bp.Ordered, parseError, binary, operation},
			//scan.TRUE:           {bp.None, parseError, parseError, nil}, // TODO: remove
			//scan.FALSE:          {bp.None, parseError, parseError, nil}, // TODO: remove
			scan.BACKSLASH:      {bp.None, parseError, createBinder, nil},
			scan.AT:             {bp.None, annotation, parseError, nil},
			scan.UNDERSCORE:     {bp.None, parseError, parseError, nil},
			scan.PACKAGE:        {bp.None, parseError, parseError, nil},
			scan.MODULE:         {bp.None, parseError, parseError, nil},
			scan.NEW_LINE:       {bp.None, ignore, parseError, nil},
			scan.ERROR:          {bp.None, parseError, parseError, nil},
			scan.EOF:            {bp.None, ignore, ignore, nil},
			scan.BANG_POSTFIX__: {bp.Postfix, parseBug, parseBug, parseBug},
			scan.PLUS_PREFIX__:  {bp.Unary, parseBug, parseBug, parseBug},
			scan.MINUS_PREFIX__: {bp.Unary, parseBug, parseBug, parseBug},
			scan.END__:          {bp.None, parseError, parseError, nil},
		}
		parseTable = tmp
		isParseTableInitialized = true
	}
	parseTableMutex.Unlock()
}

type parseTableElement struct {
	bp         bp.BindingPower
	prefixRule func(*Parser, scan.Token) bool
	infixRule  func(*Parser, scan.Token) bool
	callRule   func(*Parser, scan.Token) bool
}

var parseTable map[scan.TokenType]parseTableElement

func operation(p *Parser, token scan.Token) bool {
	// TODO
	return false
}

func mutParse(p *Parser) bool {
	panic("TODO: implement")
}

func annotationParse(p *Parser) bool {
	panic("TODO: implement")
}

func constParse(p *Parser) bool {
	panic("TODO: implement")
}

func classParse(p *Parser) bool {
	panic("TODO: implement")
}

func moduleParse(p *Parser) bool {
	validModuleDeclaration :=
		pushId(p) && // push module id
			ast.ModuleMembership{}.Make(p) // apply rule
	if !validModuleDeclaration {
		return false
	}

	p.Advance()
	if p.Current.GetType() != scan.LCURL {
		expectedLCURL(p.Current, p.Input).Print()
		return false
	}

	//p.Advance()
	validModuleDefinition :=
		parseProgram(p, parseProgramBlock) && // attempt to parse module definition
			ast.Module{}.Make(p) // attempt to create module
	return validModuleDefinition
}

func parseFunctionDeclaration(p *Parser, functionName ast.Id) bool {
	return ast.DeclareFunction(p, functionName, func(par *Parser) bool {
		// parses function body
		if par.Next.GetType() == scan.LCURL {
			par.Advance()
			// skip leading newlines
			for ty := par.Next.GetType(); ty == scan.RCURL || ty == scan.NEW_LINE; ty = par.Next.GetType() {
				if ty == scan.RCURL {
					// end of function, push empty valued expression to stack
					par.Stack.Push(ast.EmptyExpression{})
					break
				}
				par.Advance()
			}
			return parseProgram(p, parseProgramBlock)
		} else {
			return parseProgram(p, parseProgramSequence)
		}
	})
}

func parseInitialId(p *Parser) (ast.Id, bool) {
	if p.Current.GetType() != scan.ID {
		expectedID(p.Current, p.Input).Print()
		return ast.Id{}, false
	}

	ok := identifier(p, p.Current)
	if !ok {
		return ast.Id{}, false
	}
	id := p.Stack.Peek().(ast.Id)
	return id, true
}

func parseTypeDef(p *Parser, id ast.Id) bool {
	ok := ast.Type{}.MakeData(p) && parseDataType(p)
	if !ok {
		return false
	}
	ty := p.Stack.Pop().(ast.Type).GetType().(types.Data)
	id = id.SetType(ty)
	//println("TypeDef Name:", id.GetName())

	def := ast.TypeDefinition(id)

	e, ok := p.Table.DeclareLocal(def, ty)
	if !ok {
		e.ToError().Print()
	}
	p.Stack.Push(def)
	// TODO: declare/define constructors in p.Table
	return ok
}

// parses either a function definition or function application
func functionParse(p *Parser) bool {
	id, ok := parseExpressionExpectingInitial(p, 0, parseInitialId)
	//println("Name:", id.GetName())
	//ok := identifier(p, p.Current) // function name
	if !ok {
		return false
	}

	if p.Next.GetType() == scan.COLON_COLON_EQUAL {
		p.Advance()
		p.Advance()
		return parseTypeDef(p, id)
	}

	nodeType := p.Stack.GetTopNodeType()
	isNotFunction := nodeType != nodetype.APPLICATION
	if isNotFunction {
		// check if nested inside type annot
		if nodeType == nodetype.TYPE_ANNOTATION {
			annot := p.Stack.Peek().(ast.ExpressionTypeAnnotation)
			nodeType = annot.GetExpression().GetNodeType()
			isNotFunction = nodeType != nodetype.APPLICATION
		}

		if isNotFunction {
			return true // not a function def
		}
	}

	// possibly function def
	if p.Next.GetType() == scan.EQUALS { // is function def
		//p.Advance()
		p.Advance()
		return parseFunctionDeclaration(p, id)
	}
	return ok
}

// for blocks, top-level scope, and sequences
func onSemi_any(p *Parser, end *bool) bool {
	*end = false
	return ast.Sequence{}.Make(p)
}

// for blocks and top-level scope
func onNewLine_block_top(p *Parser, end *bool) bool {
	*end = false
	if p.Stack.GetTopNodeType() != nodetype.STACK_MARKER {
		//fmt.Printf("**** %s ****\n", p.Stack.GetTopNodeType().ToString())
		return ast.Sequence{}.Make(p)
	}
	return true
	//return ast.Sequence{}.Make(p)
}

func onEOF_block(p *Parser, end *bool) bool {
	*end = true
	printError(p, p.Current)
	return false
}

// for
func onEOF_seq_top(p *Parser, end *bool) bool {
	*end = true
	return true
}

// for sequences
func onNewLine_seq(p *Parser, end *bool) bool {
	*end = true
	return true
}

// for blocks
func onRCurl_block(p *Parser, end *bool) bool {
	*end = true
	return true
}

// for sequences and top-level
func onRCurl_top_seq(p *Parser, end *bool) bool {
	*end = true
	printError(p, p.Current)
	return false
}

var parseProgramBlock = parseProgramAction{
	onSemi_any, onNewLine_block_top, onRCurl_block, onEOF_block}
var parseProgramSequence = parseProgramAction{
	onSemi_any, onNewLine_seq, onRCurl_top_seq, onEOF_seq_top}
var parseProgramTop = parseProgramAction{
	onSemi_any, onNewLine_block_top, onRCurl_top_seq, onEOF_seq_top}

type parseProgramAction struct {
	onSemi    func(*Parser, *bool) bool
	onNewLine func(*Parser, *bool) bool
	onRCurl   func(*Parser, *bool) bool
	onEOF     func(*Parser, *bool) bool
}

func parseFromKey(tokenType scan.TokenType, p *Parser) bool {
	switch tokenType {
	case scan.LET:
		return letParse(p)
	case scan.MUT:
		return mutParse(p)
	case scan.AT:
		return annotationParse(p)
	case scan.CONST:
		return constParse(p)
	case scan.CLASS:
		return classParse(p)
	case scan.MODULE:
		return moduleParse(p)
	default:
		err.PrintBug()
		panic("")
	}
}

func isEndOfStatementToken(t scan.TokenType) bool {
	return t == scan.SEMI_COLON || t == scan.NEW_LINE || t == scan.EOF
}

func parseProgram(p *Parser, action parseProgramAction) bool {
	ok := true
	end := false

	p.Stack.Push(ast.StackMarker{})

	for ok && !end {
		p.Advance()

		tokenType := p.Current.GetType()
		switch tokenType {
		case scan.LET:
			fallthrough
		case scan.MUT:
			fallthrough
		case scan.AT:
			fallthrough
		case scan.CONST:
			fallthrough
		case scan.CLASS:
			fallthrough
		case scan.MODULE:
			p.Advance()
			ok = parseFromKey(tokenType, p)
		case scan.ID:
			ok = functionParse(p)
		case scan.SEMI_COLON:
			ok = action.onSemi(p, &end)
		case scan.NEW_LINE:
			ok = action.onNewLine(p, &end)
		case scan.RCURL:
			ok = action.onRCurl(p, &end)
		case scan.EOF:
			ok = action.onEOF(p, &end)
		default:
			ok = parseExpression(p, 0)
		}
	}

	ok = ok && ast.Program{}.Make(p)
	if ok {
		top := p.Stack.Pop()
		if p.Stack.GetTopNodeType() != nodetype.STACK_MARKER {
			err.PrintBug()
			panic("")
		}
		p.Stack.Pop() // remove marker
		p.Stack.Push(top)
	}
	return ok
}

func makeDefaultPackage(prog ast.Program) ast.Package {
	return ast.MakePackage(DefaultNameSpaceId, prog)
}

func endParseCallback_noPackage(p *Parser) (bool, ast.Package) {
	// should be a single `Program` node on the stack
	valid, _ := p.Stack.TryValidate([]nodetype.NodeType{nodetype.PROGRAM})
	valid = valid && len(*p.Stack) == 1
	if !valid {
		err.CompileMessage(
			"could not parse program",
			err.ERROR,
			err.SYNTAX,
			"", 0, 0, 0, "",
		).Print()
		return false, ast.Package{}
	}

	prog := p.Stack.Pop().(ast.Program)
	pack := makeDefaultPackage(prog)
	return true, pack
}

func endParseCallback_package(p *Parser) (bool, ast.Package) {
	valid := ast.Package{}.Make(p)
	if !valid {
		return false, ast.Package{}
	}

	valid = len(*p.Stack) == 1
	if !valid {
		err.CompileMessage(
			"could not parse program",
			err.ERROR,
			err.SYNTAX,
			"", 0, 0, 0, "",
		).Print()
		return false, ast.Package{}
	}

	pack := p.Stack.Pop().(ast.Package)
	return true, pack
}

// returns nil on error, else callback function to be used once program is
// finished being parsed
func maybePackage(p *Parser) (endParseCallback func(*Parser) (bool, ast.Package)) {
	isPackage := p.Next.GetType() == scan.PACKAGE

	if !isPackage {
		return endParseCallback_noPackage
	}

	p.Advance()

	p.Advance() // move (what should be an) id to p.Current
	validPackageDeclaration :=
		pushId(p) && // put id on top of stack
			ast.PacakgeMembership{}.Make(p) // make package membership
	if !validPackageDeclaration {
		return nil
	}

	nextType := p.Next.GetType()
	validPackageDeclaration = isEndOfStatementToken(nextType)
	if !validPackageDeclaration {
		expectedEndOfStatement(p.Next, p.Input).Print()
		return nil
	}

	// move end-of-package-statement token (whatever token that is) to current
	p.Advance()

	return endParseCallback_package
}

func doParse(p *Parser) (bool, ast.Package) {
	p.Advance()

	var parsedSuccesfully bool
	endParseCallback := maybePackage(p)
	parsedSuccesfully =
		endParseCallback != nil &&
			parseProgram(p, parseProgramTop)
	if !parsedSuccesfully {
		return false, ast.Package{}
	}

	return endParseCallback(p)
}

var initParser = InitParser

func Parse(in *scan.Input) (bool, ast.Package) {
	InitParseTable()
	in2, e := scan.TokenizeFromInput(in)
	if nil != e {
		e.Print()
		return false, ast.Package{}
	}

	p := initParser(in2)
	return doParse(p)
}

/*func ParseLazy(in *scan.Input) (bool, ast.Program) {
	p := initParser(in)
	return doParse(p)
}*/
