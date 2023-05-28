package parsing

import "yew/lex"
import "yew/ast"
import err "yew/error"
import "yew/type"
import "yew/symbol"

type parser struct {
	input scan.Lexer

	next scan.Token
	current scan.Token

	table symbol.SymbolTable
	
	stack ast.AstStack
	functions []ast.Function
}

type bindingPower int

const (
	none bindingPower = iota
	assign
	ordered
	equated
	logical
	sum
	product
	unary
	power
	call
	compose
)

func printSyntaxError(p *parser, token scan.Token) {
	unexpectedToken(token, p.input).Print()
}

func error(p *parser, token scan.Token) {
	printSyntaxError(p, token)
}

func printError(p *parser, errToken scan.Token) {
	error(p, errToken)
}

func generateSyntaxError(message string) (func (scan.Token, scan.Lexer) err.Error) {
	return func(t scan.Token, i scan.Lexer) err.Error {
		loc := t.GetLocation()
		return err.CompileMessage(
			message, err.ERROR, err.SYNTAX, (i).GetPath(), loc.GetLine(), loc.GetChar(), 
			t.GetSourceIndex(), (i).GetSource()).(err.Error)
	}
}

var expectedID = generateSyntaxError("expected identifier")
var expectedASSIGN = generateSyntaxError("expected assignment")
var unmatchedPAREN = generateSyntaxError("unmatched left-paren")
var unmatchedBRACK = generateSyntaxError("unmatched left-bracket")
var unexpectedToken = generateSyntaxError("unexpected token")

func parseTypePrefix(p *parser) (types.Types, bool) {
	current := p.current
	p.advance()

	switch current.GetType() {
	case scan.INT:
		return types.Int{}, true
	case scan.BOOL:
		return types.Bool{}, true
	case scan.CHAR:
		return types.Char{}, true
	case scan.STRING:
		return types.Array{ElemType: types.Char{}}, true
	case scan.FLOAT:
		return types.Float{}, true
	case scan.LPAREN:
		ty, ok := parseTypeAnotation(p)
		if scan.RPAREN != p.next.GetType() {
			if ok {
				unmatchedPAREN(current, p.input).Print()
			}
			return types.Error{}, false
		}
		p.advance()
		return ty, ok
	case scan.LBRACK:
		ty, ok := parseTypeAnotation(p)
		if scan.RBRACK != p.next.GetType() {
			if ok {
				unmatchedBRACK(current, p.input).Print()
			}
			return types.Error{}, false
		}
		p.advance()
		return types.Array{ElemType: ty}, ok
	case scan.ID:
		sym, added := p.table.GetElseAdd(current.(scan.IdToken))
		if nil == sym {
			return types.Error{}, false
		}
		if added {
			sym.SetType(types.GetNewTau())
		}
		return sym.GetType(), true
	default:
		unexpectedToken(current, p.input).Print()
		return types.Error{}, false
	}
}

func parseTypeRight(p *parser, left types.Types) (types.Types, bool) {
	switch p.current.GetType() {
	case scan.ARROW:
		right, ok := parseTypePrefix(p)
		if !ok {
			return right, ok
		}
		return types.Function{Domain: left, Codomain: right}, true
	case scan.COMMA:
		ts := make([]types.Types, 0, 1)
		ts = append(ts, left)
		for ; p.current.GetType() == scan.COMMA; {
			p.advance()
			if !isTypePrefixTokenType(p.current.GetType()) {
				break
			}

			elem, ok := parseTypeAnotation(p)
			if !ok {
				return elem, false
			}
			ts = append(ts, elem)
		}
		return types.Tuple(ts), true
	default:
		return left, true
	}
}

func parseTypeAnotation(p *parser) (types.Types, bool) {
	var left types.Types
	var ok bool
	left, ok = parseTypePrefix(p)
	if !ok {
		return left, false
	}
	left, ok = parseTypeRight(p, left)
	if !ok {
		return left, false 
	}
	return left, true
}

func isTypePrefixTokenType(t scan.TokenType) bool {
	switch t {
	case scan.ID: fallthrough
	case scan.LPAREN: fallthrough
	case scan.LBRACK: fallthrough
	case scan.INT: fallthrough
	case scan.BOOL: fallthrough
	case scan.FLOAT: fallthrough
	case scan.STRING: fallthrough
	case scan.CHAR:
		return true
	}
	return false
}

func initialParseTypeAnotation(p *parser) (types.Types, bool) {
	if !isTypePrefixTokenType(p.current.GetType()) {
		return types.GetNewTau(), true
	}

	return parseTypeAnotation(p)
	//return ty, ok
}

func getBindingPower(t scan.Token) bindingPower {
	return parseTable[t.GetType()].bp
}

func parseExpression(p *parser, bp bindingPower) bool {
	ok := parseTable[p.current.GetType()].unaryRule(p, p.current)
	for ; getBindingPower(p.next) > bp && ok; {
		p.advance()
		ok = parseTable[p.current.GetType()].infixRule(p, p.current)
	}
	return ok
}

// for decls
func parseId(p *parser) bool {
	if scan.ID != p.current.GetType() {
		expectedID(p.current, p.input).Print()
		return false
	}

	sym, added := p.table.GetElseAdd(p.current.(scan.IdToken))
	if !added && nil != sym {
		if sym.IsDefined() {
			s2 := symbol.MakeSymbol(p.current.(scan.IdToken))
			symbol.MakeRedefinedError(sym, s2).Print()
			return false
		}

		// update def loc
		loc := p.current.GetLocation()
		sym.SetLocation(loc.GetLine(), loc.GetChar(), loc.GetPath())
	} else if !added && nil == sym {
		// Something went wrong
		err.PrintBug()
		panic("")
	} // else added
	
	p.advance()

	p.stack.Push(ast.MakeId(sym))
	return true
}

func parseDeclaration(p *parser) bool {
	res := parseId(p)
	// whether it finds a type or not, a type anotation is pushed; if it doesn't 
	// find one, a tau type is pushed
	ty, ok := initialParseTypeAnotation(p)
	if !ok || !res {
		return false
	}
	t := ast.MakeTypeAnotation(ast.EmptyExpression{}, ty)
	(&p.stack).Push(t)

	return ast.Declaration{}.Make(&p.stack)
}

func parseInit(p *parser) bool {
	if p.current.GetType() != scan.EQUALS {
		expectedASSIGN(p.current, p.input).Print()
		return false
	}
	p.advance()
	return parseExpression(p, none) && 
		ast.Assignment{}.Make(&p.stack)
}

func letParse(p *parser) bool {
	// definition ::= `let` ID typeAnotation = expression
	return parseDeclaration(p) &&
			parseInit(p) &&
			ast.Definition{}.Make(&p.stack)
}

// for expressions
func identifier(p *parser, token scan.Token) bool {
	idToken := token.(scan.IdToken)
	sym, _ := p.table.GetElseAdd(idToken)
	sym.UpdateUse(idToken.GetLocation())
	p.stack.Push(ast.MakeId(sym))
	return true
}
func value_(p *parser, token scan.Token) bool {
	val := token.(scan.ValueToken).GetValue()
	p.stack.Push(ast.MakeValue(val))
	return true
}
func parseError(p *parser, token scan.Token) bool {
	unexpectedToken(token, p.input).Print()
	return false 
}

func bracket(p *parser, token scan.Token) bool {
	err.PrintBug()
	panic("TODO")
}

func paren(p *parser, token scan.Token) bool {
	err.PrintBug()
	panic("TODO")
}

func block(p *parser, token scan.Token) bool {
	err.PrintBug()
	panic("TODO")
}

func unaryOp(p *parser, token scan.Token) bool {
	err.PrintBug()
	panic("TODO")
}

func arithmetic(p *parser, token scan.Token) bool {
	err.PrintBug()
	panic("TODO")
}

func logic(p *parser, token scan.Token) bool {
	err.PrintBug()
	panic("TODO")
}

func lists(p *parser, token scan.Token) bool {
	err.PrintBug()
	panic("TODO")
}

func typeAnotation(p *parser, token scan.Token) bool {
	err.PrintBug()
	panic("TODO")
}

func ignore(*parser, scan.Token) bool {
	return true
}

func pattern(p *parser, token scan.Token) bool {
	err.PrintBug()
	panic("TODO")
}

func equals(p *parser, token scan.Token) bool {
	err.PrintBug()
	panic("TODO")
}

func composition(p *parser, token scan.Token) bool {
	err.PrintBug()
	panic("TODO")
}

func ordering(p *parser, token scan.Token) bool {
	err.PrintBug()
	panic("TODO")
}

func anotation(p *parser, token scan.Token) bool {
	err.PrintBug()
	panic("TODO")
}

var parseTable = map[scan.TokenType]struct{
	bp bindingPower
	unaryRule func(*parser, scan.Token) bool
	infixRule func(*parser, scan.Token) bool
} {
	scan.ID: {none, identifier, parseError},
	scan.VALUE: {none, value_, parseError},
	scan.LBRACK: {none, bracket, parseError},
	scan.RBRACK: {none, parseError, parseError},
	scan.LPAREN: {none, paren, parseError},
	scan.RPAREN: {none, parseError, parseError},
	scan.LCURL: {none, block, parseError},
	scan.RCURL: {none, parseError, parseError},
	scan.STRING: {none, parseError, parseError},
	scan.CHAR: {none, parseError, parseError},
	scan.BOOL: {none, parseError, parseError},
	scan.FLOAT: {none, parseError, parseError},
	scan.INT: {none, parseError, parseError},
	scan.TYPE: {none, parseError, parseError},
	scan.CLASS: {none, parseError, parseError},
	scan.WHERE: {none, parseError, parseError},
	scan.LET: {none, parseError, parseError},
	scan.CONST: {none, parseError, parseError},
	scan.MUT: {none, parseError, parseError},
	scan.LAZY: {none, parseError, parseError},
	scan.PLUS: {sum, unaryOp, arithmetic},
	scan.PLUS_PLUS: {sum, parseError, lists},
	scan.MINUS: {sum, unaryOp, arithmetic},
	scan.STAR: {product, parseError, arithmetic},
	scan.SLASH: {product, parseError, arithmetic},
	scan.HAT: {power, parseError, arithmetic},
	scan.EQUALS: {none, parseError, parseError},
	scan.PLUS_EQUALS: {none, parseError, parseError},
	scan.MINUS_EQUALS: {none, parseError, parseError},
	scan.STAR_EQUALS: {none, parseError, parseError},
	scan.SLASH_EQUALS: {none, parseError, parseError},
	scan.COLON: {sum, parseError, lists},
	scan.COLON_COLON: {assign, parseError, typeAnotation},
	scan.SEMI_COLON: {none, ignore, parseError},
	scan.QUESTION: {assign, parseError, pattern},
	scan.COMMA: {none, parseError, parseError},
	scan.BANG: {unary, unaryOp, parseError},
	scan.BANG_EQUALS: {equated, parseError, equals},
	scan.EQUALS_EQUALS: {equated, parseError, equals},
	scan.BAR: {none, parseError, parseError},
	scan.AMPER_AMPER: {logical, parseError, logic},
	scan.BAR_BAR: {logical, parseError, logic},
	scan.ARROW: {none, parseError, parseError},
	scan.FAT_ARROW: {none, parseError, parseError},
	scan.DOT: {compose, parseError, composition},
	scan.DOT_DOT: {none, parseError, parseError},
	scan.GREAT: {ordered, parseError, ordering},
	scan.LESS: {ordered, parseError, ordering},
	scan.GREAT_EQUALS: {ordered, parseError, ordering},
	scan.LESS_EQUALS: {ordered, parseError, ordering},
	scan.TRUE: {none, parseError, parseError},
	scan.FALSE: {none, parseError, parseError},
	scan.BACKSLASH: {none, parseError, parseError},
	scan.AT: {none, anotation, parseError},
	scan.UNDERSCORE: {none, parseError, parseError},
	scan.NEW_LINE: {none, ignore, parseError},
	scan.ERROR: {none, parseError, parseError},
	scan.EOF: {none, ignore, ignore},
	scan.END__: {none, parseError, parseError},
}

func (p *parser) advance() {
	p.current = p.next
	p.next = p.input.Next()
}

func mutParse(p *parser) bool {
	panic("TODO: implement")
}

func anotationParse(p *parser) bool {
	panic("TODO: implement")
}

func constParse(p *parser) bool {
	panic("TODO: implement")
}

func classParse(p *parser) bool {
	panic("TODO: implement")
}

func parseFunctionDeclaration(p *parser) bool {
	// TODO
	return false
}

// parses either a function defintion or function application
func functionParse(p *parser) bool {
	identifier(p, p.current) // updates usage information even if being declared
	if p.next.GetType() == scan.ID {
		// either application or defintion
		p.advance()
		identifier(p, p.current)

	} else if p.next.GetType() == scan.COLON_COLON {
		// just declaration
		p.advance()
		p.advance()
		return parseFunctionDeclaration(p)
	}
	//functionId := p.current.(scan.IdToken)
	return false
}

// for blocks, top-level scope, and sequences
func onSemi_any(p *parser, end *bool) bool {
	*end = false
	return ast.Sequence{}.Make(&p.stack)
}

// for blocks and top-level scope
func onNewLine_block_top(p *parser, end *bool) bool {
	*end = false
	return ast.Program{}.Make(&p.stack)
}
// for sequences
func onNewLine_seq(p *parser, end *bool) bool {
	*end = true
	return ast.Program{}.Make(&p.stack)
}

// for blocks
func onRBrack_block(p *parser, end *bool) bool {
	*end = true
	return ast.Program{}.Make(&p.stack)
}
// for sequences and top-level
func onRBrack_top_seq(p *parser, end *bool) bool {
	*end = true
	printError(p, p.current)
	return false
}

var parseProgramBlock = parseProgramGenerator(
		onSemi_any,
		onNewLine_block_top,
		onRBrack_block)
var parseProgramSequence = parseProgramGenerator(
		onSemi_any,
		onNewLine_seq,
		onRBrack_top_seq)
var parseProgramTop = parseProgramGenerator(
		onSemi_any,
		onNewLine_block_top,
		onRBrack_top_seq)

func parseProgramGenerator(
		onSemi func(*parser, *bool) bool, onNewLine func(*parser, *bool) bool,
		onRBrack func(*parser, *bool) bool) (func (*parser) bool) {
	return func(p *parser) bool {
		ok := true
		end := false
		for ; ok && !end; {
			p.advance()

			switch p.current.GetType() {
			case scan.LET:
				p.advance()
				ok = letParse(p)
			case scan.MUT:
				ok = mutParse(p)
			case scan.AT:
				ok = anotationParse(p)
			case scan.CONST:
				ok = constParse(p)
			case scan.CLASS:
				ok = classParse(p)
			case scan.ID:
				ok = functionParse(p)
			case scan.SEMI_COLON:
				ok = onSemi(p, &end)
			case scan.NEW_LINE:
				ok = onNewLine(p, &end)
			case scan.RBRACK:
				ok = onRBrack(p, &end)
				// will either be an error or end of program
				// so this function call ends either way
				return ok
			case scan.EOF:
				ok = ok && onNewLine(p, &end)
				return ok
			default:
				printError(p, p.current)
				return false
			}
		}

		return ok && ast.Program{}.Make(&p.stack)
	}
}

func doParse(p *parser) (bool, ast.Program) {
	p.advance()
	res := parseProgramTop(p)
	if !res {
		return false, ast.Program{}
	}

	// should be a single `Program` node on the stack
	if len(p.stack) != 1 {
		err.CompileMessage(
			"could not parse program", 
			err.ERROR, 
			err.SYNTAX, 
			"", 0, 0, 0, "",
		).Print()
		return false, ast.Program{}
	}

	prog := p.stack[0].(ast.Program)
	return true, prog
}

func initParser(in scan.Lexer) *parser {
	p := parser{
		input: in,
		table: symbol.InitSymbolTable(in.GetPath()),
		stack: make(ast.AstStack, 0, 0x40),
	}
	out := new(parser)
	*out = p
	return out
}

func Parse(in *scan.Input) (bool, ast.Program) {
	in2, e := scan.TokenizeFromInput(in)
	if nil != e {
		e.Print()
		return false, ast.Program{}
	}

	p := initParser(&in2)
	return doParse(p)
}

/*func ParseLazy(in *scan.Input) (bool, ast.Program) {
	p := initParser(in)
	return doParse(p)
}*/