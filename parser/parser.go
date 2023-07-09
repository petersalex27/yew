package parsing

import (
	sync "sync"
	err "yew/error"
	scan "yew/lex"
	ast "yew/parser/ast"
	bp "yew/parser/bp"
	errorgen "yew/parser/error-gen"
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

var generateSyntaxError = errorgen.GenerateSyntaxError

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
	p.Stack.Mark(p)

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
			p.Advance()
			isTuple := false
			for {
				ty, ok := parseTypeAnnotation(p, bp.None)
				if scan.COMMA == p.Next.GetType() {
					p.Advance()
					if !isTuple {
						p.Stack.Mark(p)
					}
					isTuple = true
				} else if p.Next.GetType() != scan.RPAREN {
					if ok {
						unmatchedPAREN(current, p.Input).Print()
					}
					return types.Error{}, false
				}

				p.Stack.Push(ast.MakeType(ty))
				p.Advance()
				if p.Current.GetType() == scan.RPAREN {
					break
				}
			}

			if isTuple {
				tupAst, ok := p.Stack.CutAtMark(p)
				if !(ok && p.Stack.Demark(p)) {
					err.PrintBug()
					panic("")
				}

				tup := make(types.Tuple, len(tupAst))
				for i := range tupAst {
					if tupAst[i].GetNodeType() != nodetype.TYPE {
						generateSyntaxError("expected a type")(
							tupAst[i].FindStartToken(), p.Input).Print()
						return types.Error{}, false
					}
					tup[i] = tupAst[i].(ast.Type).GetType()
				}

				p.Stack.Push(ast.MakeType(tup))
			} // else type is already on stack
		case scan.LBRACK:
			p.Advance()
			ty, ok := parseTypeAnnotation(p, bp.None)
			if scan.RBRACK != p.Next.GetType() {
				if ok {
					unmatchedBRACK(current, p.Input).Print()
				}
				return types.Error{}, false
			}
			p.Advance()
			p.Stack.Push(ast.MakeType(types.Array{ElemType: ty}))
		case scan.TYPE_ID:
			id := current.(scan.TypeIdToken)
			tau := types.MakeTau(id.ToString(), scan.ToLoc(id))
			p.Stack.Push(ast.MakeType(tau))
		case scan.ID:
			id := current.(scan.IdToken)
			tau := types.MakeTau(id.ToString(), scan.ToLoc(id))
			p.Stack.Push(ast.MakeType(tau))
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
	if !p.Stack.Demark(p) { // removes stack marker
		err.PrintBug()
		panic("")
	}

	return ty, true
}

func parseRightToLeftType(p *Parser) (types.Types, bool) {
	bp := getBindingPower(p.Current)
	p.Advance()
	// along with scan.ARROW being the only thing with its binding power,
	// 	the "- 1" makes this type-operation right-to-left
	right, ok := parseTypeAnnotation(p, bp-1)
	return right, ok
}

func parseTypeArrow(p *Parser, left types.Types) (types.Types, bool) {
	right, ok := parseRightToLeftType(p)
	if !ok {
		return right, false
	}

	return types.Function{Domain: left, Codomain: right}, true
}

func buildSingleContext(ty types.Types) (types.Context, bool) {
	if ty.GetTypeType() != types.APPLICATION {
		return types.Context{}, false // failed
	}

	app := ty.(types.Application)
	if len(app) < 2 {
		// TODO: error
		return types.Context{}, false
	} else if len(app) != 2 {
		// TODO: error
		return types.Context{}, false
	}

	// left most side should always be a tau
	if app[0].GetTypeType() != types.TAU {
		// TODO: error
		return types.Context{}, false
	} else if app[1].GetTypeType() != types.TAU {
		// TODO: error
		return types.Context{}, false
	}

	class := app[0].(types.Tau)
	typeVar := app[1].(types.Tau)
	return types.Context{ClassName: class, TypeVariable: typeVar}, true
}

// takes (<class_1> <var_1>, <class_2> <var_2>, ..) => <type>
// to (<class_1> <var_1> => <class_2> <var_2> => .. => <type>)
func buildConstraint(p *Parser, left types.Types) (types.Constraint, bool) {
	var contexts types.ConstraintContext

	if left.GetTypeType() == types.TUPLE {
		tup := left.(types.Tuple)
		if len(tup) < 1 {
			// TODO: error
			return types.Constraint{}, false
		}

		contexts = make(types.ConstraintContext, len(tup))
		for i, t := range tup {
			context, ok := buildSingleContext(t)
			if !ok {
				return types.Constraint{}, false
			}
			contexts[i] = context
			//println(context.ClassName.ToString(), context.TypeVariable.ToString())
		}
	} else {
		context, ok := buildSingleContext(left)
		if !ok {
			return types.Constraint{}, false
		}
		contexts = make(types.ConstraintContext, 1)
		contexts[0] = context
	}

	return types.Constraint{Context: contexts}, true
}

func getTypeName(ty types.Types) string {
	t := ty.GetTypeType()
	switch t {
	case types.INT:
		return "the Int type"
	case types.CHAR:
		return "the Char type"
	case types.FLOAT:
		return "the Float type"
	case types.ARRAY:
		return "array types"
	case types.BOOL:
		return "the Bool type"
	case types.CONSTRUCTOR:
		return "kinds"
	case types.DICTIONARY:
		fallthrough
	case types.DATA:
		return "data types"
	case types.TAU:
		return ty.ToString()
	case types.APPLICATION:
		return "type applications outside of class definitions"
	case types.CLASS:
		// should be impossible
		fallthrough
	case types.FUNCTION:
		fallthrough
	default:
		err.PrintBug()
		panic("")
	}
}

func validateConstraint(constraint types.Constraint, cVar types.Tau, in scan.InputStream) (bool, err.Error) {
	for _, cxt := range constraint.Context {
		// check if class constraints match function constraints
		if cxt.TypeVariable.ToString() == cVar.ToString() {
			loc := in.MakeErrorLocation(cxt.TypeVariable)
			return false, err.TypeError(
				"class type parameter cannot be constrained in a class function",
				loc,
			)
		}
	}
	return true, err.Error{}
}

func parseTypeConstraint(p *Parser, left types.Types) (types.Types, bool) {
	right, ok := parseRightToLeftType(p)
	if !ok {
		return right, false
	}
	// left_2 ==
	//		<class_name> <type_var>
	//		| (<class_name> <type_var>, ..)
	// right == <type>
	constraint, allGood := buildConstraint(p, left)
	if !allGood {
		return constraint, false
	}

	// attatch constraint to each part of right
	if right.GetTypeType() == types.FUNCTION {
		//println(right.ToString())
		if p.ParsingClass {
			if valid, e := validateConstraint(constraint, p.ClassVariable, p.Input); !valid {
				e.Print()
				return constraint, false
			}
		}
		c := constraint.Constrain(right.(types.Function))
		//println(constraint.ToString())
		//println(c.ToString())
		return c, true
	} else if right.GetTypeType() == types.APPLICATION && p.ParsingClass {
		app := right.(types.Application)
		class, ok, eMsg, loc := constraint.ConstrainApplication(app)
		if !ok {
			err.TypeError(eMsg, p.Input.MakeErrorLocation(loc)).Print()
			return class, false
		}
		constraint.Constrained = class
		return constraint, true // to get class, unwrap from constraint
	}
	err.TypeError(
		"cannot apply a type constraint to "+getTypeName(right),
		p.Input.MakeErrorLocation(right),
	).Print()
	return constraint, false
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
	parseNext := nextBp > bp && (next == scan.ARROW || next == scan.COMMA || next == scan.FAT_ARROW)

	if !parseNext {
		return left, true
	}

	p.Advance()

	switch next {
	case scan.ARROW:
		return parseTypeArrow(p, left)
	case scan.COMMA:
		return parseTypeList(p, left)
	case scan.FAT_ARROW:
		return parseTypeConstraint(p, left)
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

func parseTypeAnnotationControl(p *Parser, bp bp.BindingPower, allowComma bool) (types.Types, bool) {
	var left types.Types
	var ok bool

	for {
		left, ok = parseTypePrefix(p)
		if !ok {
			return left, false
		}
		if allowComma && p.Next.GetType() == scan.COMMA {
			p.Advance()
			p.Advance()
			p.Stack.Push(ast.MakeType(left))
			continue
		}
		break
	}

	left, ok = parseTypeRight(p, left, bp)
	if !ok {
		return left, false
	}
	return left, true
}

func parseTypeAnnotation(p *Parser, bp bp.BindingPower) (types.Types, bool) {
	return parseTypeAnnotationControl(p, bp, false)
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

func beginTypeAnnotationParse(p *Parser, expectTopExpression bool) bool {
	if !expectTopExpression {
		p.Stack.Push(ast.EmptyExpression{})
	}

	res := typeAnnotation(p, p.Current)
	if res {
		p.Advance()
	}
	return res
}

func initialParseTypeAnnotation(p *Parser) bool {
	if p.Current.GetType() == scan.COLON_COLON {
		return beginTypeAnnotationParse(p, false)
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

func parseDeclaration(p *Parser, qualifier ast.DeclarationQualifier) bool {
	// whether it finds a type or not, a type annotation is pushed; if it doesn't
	// find one, a tau type is pushed
	return parseId(p) &&
		initialParseTypeAnnotation(p) &&
		ast.Declaration{Qualifier: qualifier}.Make(p)
}

func parseInit(p *Parser) bool {
	if p.Current.GetType() != scan.EQUALS {
		expectedASSIGN(p.Current, p.Input).Print()
		return false
	}
	p.Advance()
	return parseExpression(p, bp.None) //&& ast.Assignment{}.Make(p)
}

func parseQualifiedDeclaration(p *Parser, qualifier ast.DeclarationQualifier) bool {
	// definition ::= `let` ID `::` typeAnotation = expression
	return parseDeclaration(p, qualifier) &&
		parseInit(p) &&
		ast.Definition{}.Make(p)
}

func typeIdentifier(p *Parser, token scan.Token) bool {
	tyIdToken := token.(scan.TypeIdToken)
	// add symbol to table if not already added (this does not define the
	// symbol, simply declares it if it has not yet been declared)
	var sym symbol.Symbolic = symbol.MakeSymbol(scan.IdToken(tyIdToken))
	sym, _ = p.Table.GetElseAdd(sym)
	// make type
	tau := tyIdToken.AsType()
	p.Stack.Push(ast.MakeType(tau))
	return true
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
	val := token.(scan.ValueToken)
	p.Stack.Push(ast.Value(val))
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
	p.Advance()
	if p.Current.GetType() == scan.RPAREN {
		p.Stack.Push(ast.Tuple{})
		return true
	}

	ok := parseExpression(p, 0)
	if !ok {
		return false
	}

	if p.Next.GetType() != scan.RPAREN {
		parseError(p, p.Next)
		return false
	}

	p.Advance()
	return true
}

func block(p *Parser, token scan.Token) bool {
	return parseProgram(p, false, true)
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

func pushOperation(p *Parser, token scan.OtherToken) {
	switch token.GetType() {
	case scan.PLUS:
		fallthrough
	case scan.STAR:
		fallthrough
	case scan.SLASH:
		fallthrough
	case scan.MINUS:
		fallthrough
	case scan.MOD:
		fallthrough
	case scan.HAT:
		fallthrough
	case scan.PLUS_PLUS:
		fallthrough
	case scan.COLON:
		fallthrough
	case scan.ARROW:
		fallthrough
	case scan.AMPER_AMPER:
		fallthrough
	case scan.BAR_BAR:
		fallthrough
	case scan.GREAT:
		fallthrough
	case scan.LESS:
		fallthrough
	case scan.GREAT_EQUALS:
		fallthrough
	case scan.LESS_EQUALS:
		fallthrough
	case scan.EQUALS_EQUALS:
		fallthrough
	case scan.BANG_EQUALS:
		p.Stack.Push(ast.OpType(token))
	case scan.BANG:
		fallthrough
	case scan.PLUS_PREFIX__:
		fallthrough
	case scan.MINUS_PREFIX__:
		p.Stack.Push(ast.UOpType(token))
	case scan.BANG_POSTFIX__:
		p.Stack.Push(ast.PostOpType(token))
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

	tok := fixUnaryToken(token).(scan.OtherToken)
	operationsBindingPower := getBindingPower(tok)
	pushOperation(p, tok)

	if !parseExpression(p, operationsBindingPower) {
		return false
	}

	return ast.UnaryOperation{}.Make(p)
}

func binary(p *Parser, token scan.Token) bool {
	p.Advance() // move past operation

	tok := fixArithmeticToken(token).(scan.OtherToken)
	operationsBindingPower := getBindingPower(tok)
	pushOperation(p, tok)

	if tok.GetType() == scan.BANG_POSTFIX__ { // factorial
		return ast.PostfixOperation{}.Make(p)
	}

	if !parseExpression(p, operationsBindingPower) {
		return false
	}

	return ast.BinaryOperation{}.Make(p)
}

var postfixOp = binary

func patternFunction(p *Parser, token scan.Token) bool {
	if !(ast.Parameter{}.MakePatternParam(p)) {
		return false
	}

	p.Advance()
	if !parseExpression(p, bp.Mapping) {
		return false
	}

	return ast.Lambda{}.Make2(p)
}

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
		if !parseDeclaration(par, ast.ParamDeclare) {
			return false
		}
		dec := par.Stack.Peek().(ast.Declaration)
		par.Stack.Push(dec)
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
		e(p.Input).Print()
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
	p.Advance()
	return parseExpression(p, 0) && ast.Pattern{}.MakePattern(p, token)
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
		parseExpression(p, bp.Compose-1) &&
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

func _parseApplicationBodyShared(p *Parser, applicableCheck func(scan.Token) bool) bool {
	if !applicableCheck(p.Current) {
		unexpectedToken(p.Current, p.Input).Print()
		return false
	}

	var ok bool = true
	for ok {
		currentType := p.Current.GetType()
		rule := parseTable[currentType]

		if currentType == scan.COLON_COLON {
			// never entered when applicableCheck does not allow scan.COLON_COLON
			ok = rule.callRule(p, p.Current) // first, parse type, then ...
			break                            // ... end application parse
		}

		ok = rule.callRule(p, p.Current)
		if currentType != scan.DOT { // if dot, then has already been applied
			// always entered when applicableCheck does not allow scan.DOT
			ok = ok && ast.Application{}.Make(p)
		}

		if !applicableCheck(p.Next) {
			break
		}
		// else advance and continue
		p.Advance()
	}
	return ok
}

func parseTypeApplication(p *Parser, _ scan.Token) bool {
	if p.Stack.Peek().GetNodeType() != nodetype.TYPE {
		unexpectedToken(p.Input.GetTokenAtOffset(-3), p.Input).Print()
		return false
	}
	// push back onto stack as an id
	ty := p.Stack.Pop().(ast.Type)
	id, errFn := ast.MakeIdFromType(ty)
	if errFn != nil {
		errFn(p.Input.GetTokenAtOffset(-3), p.Input).Print()
		return false
	}
	p.Stack.Push(id)

	return _parseApplicationBodyShared(p, func(t scan.Token) bool {
		tokenType := t.GetType()
		return tokenType == scan.ID || tokenType == scan.TYPE_ID
	})
}

func parseApplication(p *Parser, _ scan.Token) bool {
	if p.Stack.Peek().GetNodeType() != nodetype.IDENTIFIER {
		var transformToConstructor bool = false

		// check that top of the stack is a type id
		if p.Stack.Peek().GetNodeType() == nodetype.TYPE {
			ty := p.Stack.Pop().(ast.Type)
			if ty.GetType().GetTypeType() == types.TAU {
				transformToConstructor = true
			}
			id, e := ast.MakeIdFromType(ty)
			if e != nil {
				e(p.Input.GetTokenAtOffset(-2), p.Input).Print()
				return false
			}
			p.Stack.Push(id)
		}
		
		if !transformToConstructor {
			unexpectedToken(p.Input.GetTokenAtOffset(-2), p.Input).Print()
			return false
		}
	}
	return _parseApplicationBodyShared(p, isApplicable)
}

func newSequenceFromExpression(p *Parser) ast.Sequence {
	expr := p.Stack.Pop().(ast.Expression)
	return ast.Sequence{expr}
}

func embedStatementNewSequence(p *Parser) ast.Sequence {
	statement := p.Stack.Pop().(ast.Statement)
	return ast.Sequence{ast.MakeEmptyExpression(statement)}
}

func isDeclarationQualifier(t scan.TokenType) bool {
	return t == scan.LET || t == scan.CONST || t == scan.MUT
}

func getQualifier(t scan.TokenType) ast.DeclarationQualifier {
	switch t {
	case scan.LET:
		return ast.LetDeclare
	case scan.MUT:
		return ast.MutDeclare
	case scan.CONST:
		return ast.ConstDeclare
	}

	err.PrintBug()
	panic("")
}

func isIgnorable(t scan.TokenType) bool {
	return t == scan.NEW_LINE || t == scan.SEMI_COLON
}
func ignoreLeadingIgnorables(p *Parser) {
	for lead := p.Current.GetType(); isIgnorable(lead); lead = p.Current.GetType() {
		p.Advance()
	}
}

func parseSequence(p *Parser, _ scan.Token) bool {
	// create new sequence
	topType := p.Stack.GetTopNodeType()
	var seq ast.Sequence
	if IsExpression(topType) {
		seq = newSequenceFromExpression(p)
	} else if IsStatement(topType) {
		seq = embedStatementNewSequence(p)
	} else {
		err.PrintBug()
		panic("")
	}
	p.Stack.Push(seq)

	ok := false
	for p.Current.GetType() == scan.SEMI_COLON {
		ignoreLeadingIgnorables(p)

		currentType := p.Current.GetType()
		if isDeclarationQualifier(currentType) {
			qualifier := getQualifier(currentType)
			ok = parseQualifiedDeclaration(p, qualifier) &&
				ast.EmptyExpression{}.Make(p)
		} else {
			ok = parseExpression(p, bp.Sequencer)
		}

		ok = ok && ast.Sequence{}.Make(p)
		if !ok {
			break
		}

		if p.Next.GetType() == scan.SEMI_COLON {
			p.Advance()
		} else {
			break
		}
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
			scan.TYPE_ID:   {bp.Applicative, typeIdentifier, parseApplication, typeIdentifier},
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
			scan.SEMI_COLON:        {bp.Sequencer, parseError, parseSequence, nil},
			scan.QUESTION:          {bp.PatternMatch, parseError, pattern, nil},
			scan.COMMA:             {bp.None, parseError, parseError, comma},
			scan.BANG:              {bp.Postfix, prefix, postfixOp, operation},
			scan.BANG_EQUALS:       {bp.Equitable, parseError, binary, operation},
			scan.EQUALS_EQUALS:     {bp.Equitable, parseError, binary, operation},
			scan.BAR:               {bp.None, parseError, parseError, nil},
			scan.AMPER_AMPER:       {bp.Conjunctive, parseError, binary, operation},
			scan.BAR_BAR:           {bp.Disjunctive, parseError, binary, operation},
			scan.ARROW:             {bp.Mapping, parseError, patternFunction, nil},
			scan.FAT_ARROW:         {bp.Constraint, parseError, parseError, nil},
			scan.DOT:               {bp.Compose, parseError, compose, compose},
			scan.DOT_DOT:           {bp.None, parseError, parseError, nil},
			scan.GREAT:             {bp.Ordered, parseError, binary, operation},
			scan.LESS:              {bp.Ordered, parseError, binary, operation},
			scan.GREAT_EQUALS:      {bp.Ordered, parseError, binary, operation},
			scan.LESS_EQUALS:       {bp.Ordered, parseError, binary, operation},
			scan.BACKSLASH:         {bp.None, parseError, createBinder, nil},
			scan.AT:                {bp.None, annotation, parseError, nil},
			scan.UNDERSCORE:        {bp.None, parseError, parseError, nil},
			scan.PACKAGE:           {bp.None, parseError, parseError, nil},
			scan.MODULE:            {bp.None, parseError, parseError, nil},
			scan.NEW_LINE:          {bp.None, ignore, parseError, nil},
			scan.ERROR:             {bp.None, parseError, parseError, nil},
			scan.EOF:               {bp.None, ignore, ignore, nil},
			scan.BANG_POSTFIX__:    {bp.Postfix, parseBug, parseBug, parseBug},
			scan.PLUS_PREFIX__:     {bp.Unary, parseBug, parseBug, parseBug},
			scan.MINUS_PREFIX__:    {bp.Unary, parseBug, parseBug, parseBug},
			scan.END__:             {bp.None, parseError, parseError, nil},
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

func annotationParse(p *Parser) bool {
	panic("TODO: implement")
}

func parseClassBody(p *Parser, className ast.Id, block bool) bool {

	ignoreLeadingIgnorables(p)

	if block && p.Next.GetType() == scan.RCURL {
		err.SyntaxError(
			"cannot have an empty class definition",
			p.Input.MakeErrorLocation(p.Next),
		).Print()
		return false
	}

	class := ast.InitClass(className)
	p.Stack.Push(class)

	// loop until no more function defs
	for {
		if p.Current.GetType() != scan.ID {
			error(p, p.Current)
			return false
		}
		if !parseId(p) {
			return false
		}

		if p.Current.GetType() != scan.COLON_COLON {
			error(p, p.Next)
			return false
		}
		if !beginTypeAnnotationParse(p, true) {
			return false
		}
		ok := ast.Class{}.Make(p)
		if !ok {
			return false
		}

		if p.Current.GetType() == scan.SEMI_COLON {
			// ignore
		} else if !block {
			break
		}

		ignoreLeadingIgnorables(p)
		if block {
			if p.Current.GetType() == scan.RCURL {
				break
			}
		}
	}
	return true
}

func classParse(p *Parser) bool {
	p.ParsingClass = true
	ty, ok := parseTypeAnnotation(p, 0)
	p.ParsingClass = false

	if !ok {
		return false
	}

	if p.Next.GetType() != scan.WHERE {
		err.SyntaxError(
			"unexpected token, expected `where`",
			p.Input.MakeErrorLocation(p.Next),
		).Print()
		return false
	}
	p.Advance()
	p.Advance()
	block := p.Current.GetType() == scan.LCURL
	if block {
		p.Advance()
	}

	var className ast.Id

	if ty.GetTypeType() == types.QUALIFIER {
		// unwrap
		clsType := ty.(types.Constraint).Constrained.(types.Class)
		loc := clsType.Loc
		className = ast.MakeId(scan.MakeIdToken(clsType.Name, loc.GetLine(), loc.GetChar()))
		p.ClassVariable = clsType.TypeVariable
		p.HasConstraint = true
		p.ClassConstraint = ty.(types.Constraint)
	} else if ty.GetTypeType() == types.APPLICATION {
		app := ty.(types.Application)
		if valid, msg, loc := app.ValidClass(); !valid {
			err.TypeError(msg, p.Input.MakeErrorLocation(loc)).Print()
			return false
		}
		loc := app[0].GetLocation()
		className = ast.MakeId(scan.MakeIdToken(app[0].(types.Tau).ToString(), loc.GetLine(), loc.GetChar()))
		p.ClassVariable = app[1].(types.Tau)
	} else {
		err.SyntaxError(
			"expected something of the form `class ClassName a where ...`",
			p.Input.MakeErrorLocation(ty),
		).Print()
		return false
	}

	ok = parseClassBody(p, className, block)
	p.HasConstraint = false
	return ok
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
		parseProgram(p, false, true) && // attempt to parse module definition
			ast.Module{}.Make(p) // attempt to create module
	return validModuleDefinition
}

func parseFunctionDeclaration(p *Parser, functionName ast.Id) bool {
	return ast.DeclareFunction(p, functionName, func(par *Parser) bool {
		return parseExpression(p, 0)
	})
}

func parseInitialTypeId(p *Parser) (ast.Type, bool) {
	if p.Current.GetType() != scan.TYPE_ID {
		expectedID(p.Current, p.Input).Print()
		return ast.Type{}, false
	}

	ok := typeIdentifier(p, p.Current)
	if !ok {
		return ast.Type{}, false
	}
	ty := p.Stack.Peek().(ast.Type)
	return ty, true
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

func parseTypeDef(p *Parser) bool {
	ok := ast.Type{}.MakeData(p) && parseDataType(p)
	if !ok {
		return false
	}
	ty := p.Stack.Pop().(ast.Type).GetType().(types.Data)
	//println("TypeDef Name:", id.GetName())

	loc := ty.Loc
	name := ty.Name
	id := ast.MakeId(scan.MakeIdToken(name, loc.GetLine(), loc.GetChar()))
	id = id.SetType(ty)
	def := ast.TypeDefinition(id)

	e, ok := p.Table.DeclareLocal(def, ty)
	if !ok {
		e.ToError().Print()
	}
	p.Stack.Push(def)
	return ast.RegisterConstructors(p, ty)
}

func declareType(p *Parser) bool {
	_, ok := parseInitialTypeId(p)
	if !ok {
		return false
	}

	for p.Next.GetType() == scan.ID {
		p.Advance()
		id := p.Current.(scan.IdToken)
		typeParam := ast.MakeType(types.MakeTau(id.ToString(), scan.ToLoc(id)))
		p.Stack.Push(typeParam)
		ast.Type{}.MakeApplication(p)
	}

	if p.Next.GetType() != scan.EQUALS {
		errFn := errorgen.TypeDecExpectsEqual.Generate()
		errFn(p.Next, p.Input).Print()
		return false
	}

	p.Advance()
	p.Advance()
	return parseTypeDef(p)
}

// parses either a function definition or function application
func functionParse(p *Parser) bool {
	id, ok := parseExpressionExpectingInitial(p, 0, parseInitialId)
	if !ok {
		return false
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
		p.Advance()
		p.Advance()
		return parseFunctionDeclaration(p, id)
	}
	return ok
}

func parseFromKey(tokenType scan.TokenType, p *Parser) bool {
	switch tokenType {
	case scan.LET:
		return parseQualifiedDeclaration(p, ast.LetDeclare)
	case scan.MUT:
		return parseQualifiedDeclaration(p, ast.MutDeclare)
	case scan.CONST:
		return parseQualifiedDeclaration(p, ast.ConstDeclare)
	case scan.AT:
		return annotationParse(p)
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

func parseStatementNoIgnore(p *Parser) bool {
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
		return parseFromKey(tokenType, p)
	case scan.TYPE_ID:
		return declareType(p)
	case scan.ID:
		return functionParse(p)
	}
	return parseExpression(p, 0)
}

func parseStatement(p *Parser) bool {
	ignoreLeadingIgnorables(p)
	return parseStatementNoIgnore(p)
}

func parseProgram(p *Parser, allowEof bool, endAtRCurl bool) bool {
	ok := true

	p.Stack.Mark(p)

	for ok {
		p.Advance()
		ignoreLeadingIgnorables(p)
		curr := p.Current.GetType()
		if curr == scan.EOF {
			if allowEof {
				break
			}

			printError(p, p.Current) // TODO: should have specific error for this
			return false
		} else if curr == scan.RCURL {
			if !endAtRCurl {
				printError(p, p.Current)
				return false
			}
			break
		}

		ok = parseStatementNoIgnore(p)
	}

	ok = ok && ast.Program{}.Make(p)
	if ok {
		top := p.Stack.Pop()
		if !p.Stack.Demark(p) { // remove marker
			err.PrintBug()
			panic("")
		}

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
			"", 0, 0, []string{""},
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
			"", 0, 0, []string{""},
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
			parseProgram(p, true, false)
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
