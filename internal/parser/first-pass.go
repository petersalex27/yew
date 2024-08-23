// =================================================================================================
// Alex Peters - 2024
//
// Code for the first pass of parser. The first pass of the parser collects information about ...
//   - affixed identifiers
//   - syntax extensions
//   - syntactic blocks
//
// =================================================================================================
package parser

import (
	"github.com/petersalex27/yew/internal/token"
)

func (parser *Parser) appendTokenAndAdvance(tokens *[]token.Token, next token.Token) (ok bool) {
	*tokens = append(*tokens, next)

	before := parser.tokenPos
	if ok = parser.dropAndAdvanceGreaterIndent(); !ok {
		parser.error(ExpectedGreaterIndent)
		return
	}

	if before == parser.tokenPos {
		_ = parser.Advance()
	}
	return
}

func (parser *Parser) endSection() bool {
	before := parser.tokenPos
	end := !parser.dropAndAdvanceGreaterIndent() || parser.Peek().Type == token.EndOfTokens //|| parser.Peek().Type == token.In
	if end {
		parser.tokenPos = before
	}
	return end
}

func (parser *Parser) parseTyping() (typ TypeElem, localDef, ok bool) {
	const initialCap int = 32

	// collect either type constraint or type
	tokens := make([]token.Token, 0, initialCap)
	// index in tokens
	i := 0
	// exclusive end index of last constraint arrow
	lastConstraintArrowEnd := 0
	var next token.Token

	end := false
	if !parser.dropAndAdvanceGreaterIndent() {
		parser.error(ExpectedType)
		return typ, false, false
	}

	for !end {
		next = parser.Peek()
		// if before == parser.tokenPos {
		// 	_ = parser.Advance()
		// }
		// before = parser.tokenPos
		// next = parser.Peek()
		switch next.Type {
		case token.EndOfTokens:
			end = true
		case token.ColonEqual:
			localDef = true
			end = true
		case token.ThickArrow:
			lastConstraintArrowEnd = i + 1
			fallthrough
		default:
			i++
			ok = parser.appendTokenAndAdvance(&tokens, next)
			if !ok {
				return
			}
			end = parser.endSection()
		}
	}

	if len(tokens) == 0 {
		parser.error(ExpectedType)
		return typ, localDef, false
	}

	typ.Constraint = tokens[:lastConstraintArrowEnd]
	typ.Type = tokens[lastConstraintArrowEnd:i]
	typ.Start = tokens[0].Start
	typ.End = tokens[i-1].End
	return typ, localDef, true
}

func parseDeclarationMaybeType(parser *Parser) (ok bool) {
	// create a new context and make sure it's removed when this function returns
	declaredName, cleaner, validated := parser.openSection((*Parser).getDeclarableId)
	defer cleaner.Clean()

	if ok = !parser.panicking; !ok {
		return
	} else if ok = validated; !ok {
		panic("bug: declarable name not found at parser's current token")
	}

	var decls []DeclarationElem
	decls, ok = parser.parseDeclaration_shared(declaredName)
	if !ok {
		return
	}

	if test := parser.dropAndAdvanceGreaterIndent(); !test {
		return parser.writeDecls(decls)
	}

	if parser.Peek().Type != token.Where {
		parser.error(ExpectedWhere)
		return
	}

	if len(decls) > 1 {
		parser.errorOn(IllegalTypeConsList, decls[1])
		return false
	}

	// get declared type constructor
	dec := decls[0]

	var typ DataTypeElem
	typ.TypeConstructor.Name = dec.Name
	typ.TypeConstructor.Typing = dec.Typing
	typ.DataConstructors, ok = parser.parseJustDecls()
	if !ok {
		return
	}
	return parser.write(typ)
}

// assumes some kind of identifier has already been seen
func parseDeclaration(parser *Parser) (ok bool) {
	// create a new context and make sure it's removed when this function returns
	declaredName, cleaner, validated := parser.openSection((*Parser).getDeclarableId)
	defer cleaner.Clean()

	if ok = !parser.panicking; !ok {
		return
	} else if ok = validated; !ok {
		panic("bug: declarable name not found at parser's current token")
	}

	var decls []DeclarationElem
	decls, ok = parser.parseDeclaration_shared(declaredName)
	if ok {
		return parser.writeDecls(decls)
	}
	return
}

func (parser *Parser) parseLocalDefinitionIteration() (local *localDefinition, ok bool) {
	local = new(localDefinition)
	*local, ok = parser.parseBindRHS(endForLocalDef)
	return
}

// parses a comma surrounded by optional newlines and greater-indents
//
// rule, with indent state [.., n]:
//
//	comma ::= { indent>=n } `,` { indent>=n }
//
// errorIfNone is true iff an error is reported when the comma grammar rule above is not satisfied
func (parser *Parser) parseComma(errorIfNone bool) (comma token.Token, found, ok bool) {
	// if no comma is found, still return okay
	okIfNone := !errorIfNone
	// original position
	before := parser.tokenPos

	// allow valid newlines and indents before ','
	if ok = parser.dropAndAdvanceGreaterIndent(); !ok {
		parser.conditionalError(errorIfNone, before, ExpectedGreaterIndent)
		return comma, false, okIfNone
	}
	// get comma
	if comma, ok = parser.get(token.Comma); !ok {
		parser.conditionalError(errorIfNone, before, ExpectedComma)
		return comma, false, okIfNone
	}
	// allow valid newlines and indents after ','
	if ok = parser.dropAndAdvanceGreaterIndent(); !ok {
		parser.conditionalError(errorIfNone, before, ExpectedGreaterIndent)
		return comma, false, okIfNone
	}
	return comma, true, true
}

func (parser *Parser) parseLocalDefinition(decls []DeclarationElem) (_ []DeclarationElem, ok bool) {
	if len(decls) == 0 {
		panic("bug: no declarations to associate with local definition")
	}

	// allow valid newlines and indents after ':='
	if ok = parser.dropAndAdvanceGreaterIndent(); !ok {
		parser.error(ExpectedGreaterIndent)
		return decls, false
	}

	decls[0].localDefinition, ok = parser.parseLocalDefinitionIteration()
	if !ok || len(decls) == 1 {
		return decls, false
	}

	for i := range decls[1:] {
		if _, _, ok = parser.parseComma(true); !ok {
			return decls, false
		}

		decls[i+1].localDefinition, ok = parser.parseLocalDefinitionIteration()
		if !ok {
			return decls, false
		}
	}

	return decls, true
}

func (parser *Parser) parseDeclHeadList(declaredName token.Token) (decls []DeclarationElem, ok bool) {
	decls = make([]DeclarationElem, 1, 4)
	decls[0] = DeclarationElem{Name: declaredName, Start: declaredName.Start}

	for parser.Peek().Type == token.Comma {
		_ = parser.Advance()
		if ok = parser.dropAndAdvanceGreaterIndent(); !ok {
			parser.error(ExpectedGreaterIndent)
			return
		}

		// get next name to declare
		declaredName, ok = parser.getDeclarableId()
		if !ok {
			parser.error(ExpectedIdentifier)
			return
		}
		decls = append(decls, DeclarationElem{Name: declaredName, Start: declaredName.Start})

		if ok = parser.dropAndAdvanceGreaterIndent(); !ok {
			parser.error(ExpectedGreaterIndent)
			return
		}
	}
	return decls, true
}

func (parser *Parser) parseDeclaration_shared(declaredName token.Token) (decls []DeclarationElem, ok bool) {
	if ok = parser.dropAndAdvanceGreaterIndent(); !ok {
		parser.error(ExpectedGreaterIndent)
		return
	}

	if decls, ok = parser.parseDeclHeadList(declaredName); !ok {
		return
	}

	_, ok = parser.get(token.Colon)
	if !ok {
		parser.error(ExpectedTyping)
		return
	}

	// now parse (shared if multiple decls) typing
	var typing TypeElem
	var localDef bool
	if typing, localDef, ok = parser.parseTyping(); !ok {
		return
	}

	// give type to decls
	for i, decl := range decls {
		decl.Typing = typing
		decls[i] = decl
	}
	// each declaration must be associated with a local definition if ':=' was found
	//
	// each definition is separated by a ','
	if localDef {
		if decls, ok = parser.parseLocalDefinition(decls); !ok {
			return
		}
	}
	return
}

// reports any errors b/c of non-declarations showing up in symbolic elems
func (parser *Parser) reportNonDeclsInElems() {
	count := int(parser.saver.elems.GetCount())
	elems, _ := parser.saver.elems.MultiCheck(count)
	for _, elem := range elems {
		if _, ok := elem.(DeclarationElem); !ok {
			start, end := elem.Pos()
			parser.error2(ExpectedDeclaration, start, end)
		}
	}
}

func (parser *Parser) parseJustDecls() (decls []DeclarationElem, ok bool) {
	if ok = parser.Peek().Type == token.Where; !ok {
		parser.error(ExpectedWhere)
		return
	}

	clean := parser.saveStacks() // create return point
	defer clean()                // ensure return to return point

	// now, parse where
	again := func(p *Parser, i int) bool { return p.equalIndent(i) }
	ok = parser.parseSection(parseDeclaration, again)
	if !ok {
		return
	}

	var recorded int
	_, _, _, decls, _, recorded = parser.stackValues_recorded()
	// all that was written was declarations?
	if ok = recorded == len(decls); !ok {
		parser.reportNonDeclsInElems()
		return nil, false
	}

	return decls, true
}

func (parser *Parser) parseTraitHelper() (trait SpecElem, ok bool) {
	returns := parser.saveStacks()
	defer returns()
	// TODO: this is an assumption, might need to change
	//
	// Usually traits will just be of the following form:
	//	trait TraitName v where ...
	// so, set initial capacity to two for best average time
	const initialCap int = 2

	// collect either type constraint or type
	tokens := make([]token.Token, 0, initialCap)
	// index in tokens
	i := 0
	// exclusive end index of last constraint arrow
	lastConstraintArrowEnd := 0

	ok = parser.dropAndAdvanceGreaterIndent()
	if !ok {
		return
	}

	next := parser.Peek()

	// loop until `where` is found
	for next.Type != token.Where {
		switch next.Type {
		case token.EndOfTokens:
			ok = false
			parser.error(UnexpectedEOF)
			return
		case token.ThickArrow:
			lastConstraintArrowEnd = i + 1
		}

		i++
		tokens = append(tokens, next)
		before := parser.tokenPos
		if ok = parser.dropAndAdvanceGreaterIndent(); !ok {
			parser.error(ExpectedGreaterIndent)
			return
		}
		if parser.tokenPos == before {
			_ = parser.Advance()
		}
		next = parser.Peek()
	}

	trait.Constraint = tokens[:lastConstraintArrowEnd]
	trait.Head = tokens[lastConstraintArrowEnd:]

	// set declarations in trait
	trait.MethodDeclarations, ok = parser.parseJustDecls()
	return trait, ok
}

func parseTrait(parser *Parser) (ok bool) {
	var trait SpecElem
	if trait, ok = parser.parseTraitHelper(); ok {
		// add trait
		return parser.write(trait)
	}
	return
}

func (parser *Parser) parseSpec() (ok bool) {
	if parser.Peek().Type != token.Spec {
		panic("bug: current token in stream not validated")
	}
	again := func(p *Parser, i int) bool { return p.equalIndent(i) }
	return parser.parseSection(parseTrait, again)
}

func (parser *Parser) parseDeclaration() (ok bool) {
	switch parser.Peek().Type {
	case token.Id, token.Infix:
		return parseDeclaration(parser)
	default:
		panic("bug: current token in stream not validated")
	}
}

func (parser *Parser) parseDeclarationMaybeType() (ok bool) {
	switch parser.Peek().Type {
	case token.Id, token.Infix:
		again := func(p *Parser, i int) bool { return p.equalIndent(i) }
		return parser.parseSection(parseDeclarationMaybeType, again)
	default:
		panic("bug: current token in stream not validated")
	}
}

func (parser *Parser) declCount() int {
	return int(parser.saver.decls.GetCount())
}

func (parser *Parser) bindingCount() int {
	return int(parser.saver.bindings.GetCount())
}

func (parser *Parser) typeCount() int {
	return int(parser.saver.types.GetCount())
}

func (parser *Parser) traitCount() int {
	return int(parser.saver.traits.GetCount())
}

func (parser *Parser) instCount() int {
	return int(parser.saver.inst.GetCount())
}

type SyntaxClass uint16

const (
	EmptyClass SyntaxClass = 0
	DeclClass  SyntaxClass = 1 << iota
	FuncClass
	TypeClass
	TraitClass
	InstanceClass
	TypeConsClass
	MutualClass
	ExtensionClass
	AnyClass SyntaxClass = ^SyntaxClass(0)
)

func checkClass(c SyntaxClass) func(se SymbolicElem) bool {
	return func(se SymbolicElem) bool { return se.GetSyntaxClass()&c != 0 }
}

// returns most recently parsed syntactic element
func (parser *Parser) lastSyntacticElemParsed(cls SyntaxClass) (elem SyntacticElem, ok bool) {
	predicate := checkClass(cls)
	ok, elem = parser.saver.elems.Search(predicate)
	return
}

// for mutual blocks
func (parser *Parser) stackValues() (types []DataTypeElem, traits []SpecElem, insts []InstanceElem, decls []DeclarationElem, bindings []BindingElem) {
	countDecls := parser.declCount()
	countBindings := parser.bindingCount()
	countTypes := parser.typeCount()
	countTraits := parser.traitCount()
	countInst := parser.instCount()

	types, _ = parser.saver.types.MultiCheck(countTypes)
	traits, _ = parser.saver.traits.MultiCheck(countTraits)
	insts, _ = parser.saver.inst.MultiCheck(countInst)
	bindings, _ = parser.saver.bindings.MultiCheck(countBindings)
	decls, _ = parser.saver.decls.MultiCheck(countDecls)
	return
}

// total written is NOT just the same of the lengths of the values returned; it is the total number
// of things written while the stacks were saved which may differ from the number of things returned
func (parser *Parser) stackValues_recorded() (types []DataTypeElem, traits []SpecElem, insts []InstanceElem, decls []DeclarationElem, bindings []BindingElem, totalWritten int) {
	totalWritten = int(parser.saver.elems.GetCount())
	types, traits, insts, decls, bindings = parser.stackValues()
	return
}

// for non-mutual sections
func (parser *Parser) stackElems() []SymbolicElem {
	count := int(parser.saver.elems.GetCount())
	syms, _ := parser.saver.elems.MultiCheck(count)
	return syms
}

func (parser *Parser) saveStacks() (clean func()) {
	// create return points
	parser.saver.decls.Save()
	parser.saver.bindings.Save()
	parser.saver.traits.Save()
	parser.saver.inst.Save()
	parser.saver.types.Save()
	parser.saver.elems.Save()
	// ensure return to return points
	return func() {
		parser.saver.decls.Return()
		parser.saver.bindings.Return()
		parser.saver.traits.Return()
		parser.saver.inst.Return()
		parser.saver.types.Return()
		parser.saver.elems.Return()
	}
}

func whereValidator(p *Parser) (token.Token, bool) {
	out := p.Advance()
	return out, out.Type == token.Where
}

func (parser *Parser) whereClause() (whereToken token.Token, ok bool) {
	var tok token.Token
	tok, cleaner, validated := parser.openSection(whereValidator)
	whereToken = tok
	defer cleaner.Clean()

	if ok = !parser.panicking; !ok {
		return
	} else if ok = validated; !ok {
		panic("bug: 'where' not found at parser's current token")
	}

	againCondition := func(p *Parser, i int) bool { return p.equalIndent(i) }
	ok = parser.runSection(!cleaner.SectionOpened(), parseClause, againCondition)
	if !ok {
		return
	}
	return
}

func (parser *Parser) parseWhere() (where WhereClause, ok bool) {
	returns := parser.saveStacks()
	defer returns()

	var whereToken token.Token
	whereToken, ok = parser.whereClause()
	if !ok {
		return
	}

	where.Start = whereToken.Start
	where.Elems = parser.stackElems()
	return where, ok
}

func letValidator(parser *Parser) (token.Token, bool) {
	out := parser.Advance()
	return out, out.Type == token.Let
}

func (parser *Parser) parseLet() (let LetBindingElem, ok bool) {
	returns := parser.saveStacks()
	defer returns()

	letToken, cleaner, validated := parser.openSection(letValidator)
	defer cleaner.Clean()

	if ok = !parser.panicking; !ok {
		return
	} else if ok = validated; !ok {
		panic("bug: 'let' not found at parser's current token")
	}

	againCondition := func(p *Parser, i int) bool { return p.equalIndent(i) }
	ok = parser.runSection(!cleaner.SectionOpened(), parseLetClause, againCondition)
	if !ok {
		return
	}

	if ok = parser.dropAndAdvanceGreaterIndent(); !ok {
		parser.error(ExpectedGreaterIndent)
		return
	}

	if ok = parser.Peek().Type == token.In; !ok {
		next := parser.Peek()
		if next.Type == token.EndOfTokens {
			parser.error(UnexpectedEOF)
		} else {
			parser.error2(ExpectedIn, next.Start, next.End)
		}
		return
	}

	let.Start = letToken.Start
	let.Elems = parser.stackElems()
	return let, ok
}

func (parser *Parser) parseHaskellDataType(typeCons TypeConstructorElem) (ok bool) {
	panic("TODO: implement")
}

// parse the data constructors for the type with type constructor `typeCons`
//
// on success: pushes the data type (composition of type and data constructors) and true is returned
//
// on failure: no data type is pushed, `parser` reports an error, and false is returned
func (parser *Parser) parseDataConstructors(typeCons TypeConstructorElem) (ok bool) {
	var data DataTypeElem
	data.TypeConstructor = typeCons
	if data.DataConstructors, ok = parser.parseJustDecls(); ok {
		return parser.write(data)
	}
	return ok
}

func (parser *Parser) writeHigherOrderFunctionDecl(typeCons TypeConstructorElem) (ok bool) {
	// higher order function, push result
	higherOrderDecl := typeCons.DeclarationElem
	return parser.write(higherOrderDecl)
}

// parses a data type (Haskell-like style or declaration style) or higher order function declaration
//
// returns true iff no errors occurred
//
// function panics if preconditions for calling this function are not met:
//   - next token is a colon (':')
//   - and, length of scrutinee < 1
//
// all callers of this function should check the above conditions
func (parser *Parser) parseDataType(scrutinee []token.Token) (ok bool) {
	if _, ok = parser.get(token.Colon); !ok {
		panic("bug: expected typing")
	}

	if len(scrutinee) < 1 {
		panic("bug: illegal scrutinee")
	}

	var typeCons TypeConstructorElem
	var localDef bool
	typeCons.Name = scrutinee[0]
	typeCons.Typing, localDef, ok = parser.parseTyping()
	if !ok {
		return
	} else if localDef {
		parser.errorOn(IllegalTypeConsLocalDef, typeCons)
		return false
	}

	ok = parser.dropAndAdvanceGreaterIndent()
	if !ok {
		parser.error(ExpectedGreaterIndent)
		return
	}

	// check if Haskell-style data type
	haskellStyle := len(scrutinee) > 1
	if haskellStyle {
		//typeCons.Args = scrutinee[1:]
		return parser.parseHaskellDataType(typeCons)
	}
	//typeCons.Args = []token.Token{} // no args, not Haskell style data type

	// test if declaration of higher order function or data type definition
	if parser.Peek().Type != token.Where {
		return parser.writeHigherOrderFunctionDecl(typeCons)
	}
	return parser.parseDataConstructors(typeCons)
}

func (parser *Parser) prepareNextExpressionIteration() (next token.Token, end bool, ok bool) {
	ok = true
	_ = parser.Advance()
	before := parser.tokenPos
	end = !parser.dropAndAdvanceGreaterIndent()
	if end {
		parser.tokenPos = before
		return
	}
	next = parser.Peek()
	return
}

var endOnWhere = map[token.Type]bool{
	token.Where:  true,
	token.Mutual: true,
	token.In:     true,
}

var endForLocalDef = map[token.Type]bool{
	token.Where:  true,
	token.Mutual: true,
	token.In:     true,
	token.Comma:  true,
}

func (parser *Parser) parseExpression() (exprs []ExpressionElem, clause ClauseElem, ok bool) {
	return parser.parseExpression2(endOnWhere)
}

func (parser *Parser) parseExpression2(endMap map[token.Type]bool) (exprs []ExpressionElem, clause ClauseElem, ok bool) {
	const initialCap = 32
	toks := make([]token.Token, 0, initialCap)
	exprs = make([]ExpressionElem, 0, 8)
	i := 0

	if !parser.dropAndAdvanceGreaterIndent() {
		parser.error(ExpectedGreaterIndent)
		return exprs, clause, false
	}

	next := parser.Peek()

	// loop until endType or a exiting-indent is found
	isLet := false
	end := false
	for !end {
		for !endMap[next.Type] {
			var tok token.Token
			if next.Type == token.EndOfTokens {
				break // this counts as an "exiting-indent"
			} else if isLet = next.Type == token.Let; isLet {
				break
			} else if next.Type == token.Case {
				panic("TODO: implement case expressions")
			} else {
				tok = next
			}

			i++
			toks = append(toks, tok)

			next, end, ok = parser.prepareNextExpressionIteration()
			if !ok {
				return
			} else if end {
				break
			}
		}

		var expr ExpressionElem
		// add expression to list
		if ok && i != 0 {
			expr = TokensElem(toks[:i])
			exprs = append(exprs, expr)
			// reset tokens
			toks = make([]token.Token, 0, initialCap)
			i = 0
		}

		// check if parsing stopped because of `let`
		if isLet {
			// parse let-binding expression
			isLet = false
			if expr, ok = parser.parseLet(); !ok {
				return exprs, clause, false
			}
			// add let-binding expression to list
			exprs = append(exprs, expr)
			next, end, ok = parser.prepareNextExpressionIteration()
			if !ok {
				return
			} else if end {
				break
			}
		}
	}

	// make sure something was added to tokens array
	if ok = len(exprs) != 0; !ok {
		parser.error(ExpectedExpression)
		return exprs, clause, false
	}

	// check if parsing stopped because of `where`
	if endMap[token.Where] && next.Type == token.Where {
		clause, ok = parser.parseWhere()
	} else if endMap[token.Mutual] && next.Type == token.Mutual {
		clause, ok = parser.parseMutualWhere()
	} else {
		// set end position of expression (use where to do so)
		// tokens has a length greater than one since it's guarded against above
		_, end := exprs[len(exprs)-1].Pos()
		clause = WhereClause{End: end}
	}

	return exprs[:i], clause, ok
}

func (parser *Parser) parseScrutinee(endOfScrutinee map[token.Type]bool) (scrutinee []token.Token, scrutinizingDataType, ok bool) {
	const initialCap int = 32 // TODO: larger initial cap?

	// collect either type constraint or type
	scrutinee = make([]token.Token, 0, initialCap)
	expressionIndex := 0

	ok = parser.dropAndAdvanceGreaterIndent()
	if !ok {
		parser.error(UnexpectedIndent)
		return
	}

	next := parser.Peek()

	secondCond := func(bool, int) bool { return false }
	eqEnds := endOfScrutinee[token.Equal]
	if eqEnds {
		secondCond = func(foundEq bool, i int) bool { return foundEq && i > 0 }
	}
	openParens := 0
	var firstOpen token.Token
	var eqEnclosed bool = false

	for !endOfScrutinee[next.Type] || secondCond(eqEnclosed, openParens) {
		if next.Type == token.EndOfTokens {
			parser.error(UnexpectedEOF)
			return scrutinee, false, false
		} else if scrutinizingDataType = next.Type == token.Colon; scrutinizingDataType {
			return
		} else if next.Type == token.LeftParen {
			if openParens == 0 {
				firstOpen = next // first not yet closed paren
			}
			openParens++
		} else if next.Type == token.RightParen {
			openParens--
			if ok = !(openParens >= 0); !ok {
				parser.error(UnexpectedRParen)
				return
			} else if openParens == 0 {
				eqEnclosed = false // reset (regardless of whether the enclosed term contained `=`)
			}
		} else if next.Type == token.Equal && eqEnds {
			eqEnclosed = openParens != 0 // `=` is the type `= : x -> y -> Type`
			if !eqEnclosed {
				panic("bug: `=` didn't break loop, but it should've")
			}
		}

		expressionIndex++
		scrutinee = append(scrutinee, next)
		before := parser.tokenPos
		if ok = parser.dropAndAdvanceGreaterIndent(); !ok {
			parser.error(ExpectedGreaterIndent)
			return
		}
		if parser.tokenPos == before {
			_ = parser.Advance()
		}
		next = parser.Peek()
	}

	if ok = openParens == 0; !ok {
		// must have at least one open paren
		parser.errorOn(UnmatchedLParen, firstOpen)
		return
	}

	if ok = endOfScrutinee[next.Type]; !ok {
		msg := expectedMessageMulti(endOfScrutinee)
		parser.error(msg)
		return
	} else if ok = len(scrutinee) != 0; !ok {
		parser.error(ExpectedScrutinee)
		return
	}

	_ = parser.Advance() // move past end marker
	return
}

func (parser *Parser) parseBindRHS(endMap map[token.Type]bool) (local localDefinition, ok bool) {
	local.Expression, local.Clause, ok = parser.parseExpression2(endMap)
	if ok && len(local.Expression) == 0 {
		ok = false
		parser.error(ExpectedExpression)
		return
	} else if ok {
		local.Start, _ = local.Expression[0].Pos()
		if local.Clause != nil {
			_, local.End = local.Clause.Pos()
		} else {
			_, local.End = local.Expression[len(local.Expression)-1].Pos()
		}
	} // else !ok, error already reported
	return
}

func (parser *Parser) defineScrutinee(scrutinee []token.Token) (ok bool) {
	var binding BindingElem
	if len(scrutinee) == 0 {
		panic("bug: no scrutinee to define")
	}
	binding.Head = scrutinee
	binding.Start = binding.Head[0].Start // guaranteed to have at least one element in result
	binding.localDefinition, ok = parser.parseBindRHS(endOnWhere)
	if ok {
		parser.write(binding) // save result
	}
	return ok
}

var scrutineeEnd = map[token.Type]bool{
	token.Equal: true,
}

func parseScrutinee(parser *Parser) (ok bool) {
	var scrutinee []token.Token
	var scrutinizingDataType bool
	scrutinee, scrutinizingDataType, ok = parser.parseScrutinee(scrutineeEnd)
	if !ok {
		return
	} else if scrutinizingDataType {
		return parser.parseDataType(scrutinee)
	}
	return parser.defineScrutinee(scrutinee)
}

func (parser *Parser) parseDefinition() (ok bool) {
	againCondition := func(*Parser, int) bool { return false }
	return parser.runSection(true, parseScrutinee, againCondition)
}

func (parser *Parser) getDeclarableId() (token.Token, bool) {
	// check if just an id or possibly an infix id
	first := parser.Peek()
	if first.Type != token.LeftParen {
		if first.Type == token.Id {
			return parser.Advance(), true
		}
		return first, false
	}

	// get remaining infix tokens

	// get id
	name := parser.Peek()
	if name.Type != token.Id {
		return first, false
	}
	_ = parser.Advance()

	// get right paren
	end := parser.Peek()
	if end.Type != token.RightParen {
		return first, false
	}
	_ = parser.Advance()

	// create infix id
	// TODO: test what happens when an infixed id is used for a term with 0 arity
	res := token.Infix.MakeValued("(" + name.Value + ")")
	res.Start = first.Start
	res.End = first.End
	return res, true
}

func (parser *Parser) nextIsEnclosed(first token.Token) (_ token.Token, ok bool) {
	numLeftParens := 1
	next := parser.Peek()
	for ; next.Type == token.LeftParen; numLeftParens++ {
		_ = parser.Advance()
		next = parser.Peek()
	}

	// check if identifier
	if next.Type != token.Id {
		// not an id => not enclosed
		return first, false
	}

	// get id
	out := next
	_ = parser.Advance()
	// if the same number of right parens are found, then it's an enclosed id
	next = parser.Peek()
	for ; numLeftParens > 0; numLeftParens-- {
		if next.Type != token.RightParen {
			// number of right parens doesn't _immediately_ match, might still match later, but not an
			// enclosed id
			return out, false
		}
		_ = parser.Advance()
		next = parser.Peek()
	}
	_ = parser.Advance()
	// number of right parens matches, so it's an enclosed id
	// 	((..(myId)..))
	return out, true
}

func (parser *Parser) parseDeclarationOrDefinition() (ok bool) {
	before := parser.tokenPos
	first := parser.Peek()

	_ = parser.Advance()

	var isEnclosed bool
	if first.Type == token.LeftParen {
		if first, isEnclosed = parser.nextIsEnclosed(first); !isEnclosed {
			parser.tokenPos = before
			return parser.parseDefinition()
		}
	}

	if !parser.dropAndAdvanceGreaterIndent() {
		parser.error(ExpectedGreaterIndent)
		return false
	}
	next := parser.Peek()
	parser.tokenPos = before // rewind

	switch next.Type {
	case token.Colon, token.Comma:
		if validTypeIdent(first.Value) {
			return parser.parseDeclarationMaybeType()
		}
		return parser.parseDeclaration()
	default:
		return parser.parseDefinition()
	}
}

func mutualClause(parser *Parser) (ok bool) {
	switch next := parser.Peek(); next.Type {
	case token.Alias:
		panic("TODO: implement 'alias'")
	case token.Spec:
		ok = parser.parseSpec()
	case token.Id, token.Infix:
		ok = parser.parseDeclarationOrDefinition()
	case token.EndOfTokens:
		ok = true // end should be signaled inside runSection when it checks
	default:
		parser.error(UnexpectedToken)
		ok = false
	}
	return
}

func clauseIteration(parser *Parser, endException token.Type) (ok bool) {
	switch next := parser.Peek(); next.Type {
	case token.LeftParen, token.Id, token.Infix:
		ok = parser.parseDeclarationOrDefinition()
	case token.EndOfTokens, endException:
		ok = true
	default:
		parser.error(UnexpectedToken)
		ok = false
	}
	return
}

func parseLetClause(parser *Parser) (ok bool) {
	return clauseIteration(parser, token.In)
}

func parseClause(parser *Parser) (ok bool) {
	return clauseIteration(parser, token.EndOfTokens)
}

func (parser *Parser) parseVisibility(visType token.Type, vis Visibility) (ok bool) {
	// verify parser is exploring the top level
	if !parser.ExploringTopLevel() {
		// bug: can't set the visibility of local declarations
		panic("bug: tried to set visibility outside of top level")
	}

	if parser.Peek().Type != visType {
		panic("bug: current token in stream not validated")
	}

	// stack new visibility on top of previous visibility
	parser.setVisibility(vis)
	defer parser.restoreVisibility()

	before := parser.saver.elems.GetFullCount()

	again := func(p *Parser, i int) bool { return p.equalIndent(i) }
	var visToken token.Token
	ok, visToken = parser.parseSection2(parseClause, again)
	if !ok {
		return
	}

	after := parser.saver.elems.GetFullCount()
	if ok = before < after; ok {
		return
	}
	// if a visibility modifier is set, it must be used
	parser.error2(UnusedVisibility, visToken.Start, parser.tokenPos)
	return
}

func (parser *Parser) parseAnnotation() (ok bool) {
	annot := parser.Advance()
	switch annot.Value {
	case "builtin":
	case "error":
	case "warn":
	case "deprecated":
	case "todo":
	case "external":
	case "inline":
	case "noInline":
	case "specialize":
	case "noAlias":
	case "pure":
	case "noGc":
	case "infixl", "infixr":
	case "infix":
	}
	// if parser.Peek().Type != token.LeftParen {
	// 	// TODO
	// 	return ok
	// }
	panic("TODO: implement")
}

// calculates positions of mutual block from the start token and parser's internal state.
//
// Fails, returning ok=false, when there are no syntactic elements parsed inside the block
func (parser *Parser) calcMutualBlockPos(mutToken token.Token) (start, end int, ok bool) {
	start = mutToken.Start
	// everything except extension and mutual
	const valid = AnyClass ^ ExtensionClass ^ MutualClass
	elem, ok := parser.lastSyntacticElemParsed(valid)
	if !ok {
		parser.error2(EmptyMutualBlock, mutToken.Start, mutToken.End)
		return
	}
	_, end = elem.Pos()
	return
}

func (parser *Parser) mutualInnerHelper(atMostOnce bool, mutToken token.Token) (mut MutualBlockElem, ok bool) {
	againCondition := func(p *Parser, i int) bool { return p.equalIndent(i) }
	ok = parser.runSection(atMostOnce, mutualClause, againCondition)
	if !ok {
		return
	}

	// set positions
	mut.Start, mut.End, ok = parser.calcMutualBlockPos(mutToken)
	return
}

func (parser *Parser) parseMutualInner(actualOpen token.Type) (mut MutualBlockElem, ok bool) {
	mutToken, cleaner, validated := parser.openSection(
		func(p *Parser) (token.Token, bool) {
			out := p.Advance()
			return out, out.Type == actualOpen
		})
	defer cleaner.Clean()

	if ok = !parser.panicking; !ok {
		return
	} else if ok = validated; !ok {
		panic("bug: correct token not found at parser's current token")
	}
	return parser.mutualInnerHelper(!cleaner.SectionOpened(), mutToken)
}

func (parser *Parser) illegalTraitErrors(ts []SpecElem) {
	for _, t := range ts {
		parser.errorOn(IllegalTrait, t)
	}
}

func (parser *Parser) illegalInstanceErrors(is []InstanceElem) {
	for _, i := range is {
		parser.errorOn(IllegalTrait, i)
	}
}

func (parser *Parser) parseMutualWhere() (mut MutualBlockElem, ok bool) {
	mutualToken, cleaner, validated := parser.openMutualWhereSection()
	if !validated {
		return
	}
	// this is the only time that 'mutual' does not need to be followed by a newline
	defer cleaner.Clean()

	returns := parser.saveStacks()
	defer returns()

	if mut, ok = parser.mutualInnerHelper(!cleaner.SectionOpened(), mutualToken); !ok {
		return
	}

	n := 0
	mut.Types, mut.Traits, mut.Instances, mut.Declarations, mut.Bindings, n = parser.stackValues_recorded()
	ok = len(mut.Traits) == 0 && len(mut.Instances) == 0
	if !ok {
		parser.illegalTraitErrors(mut.Traits)
		parser.illegalInstanceErrors(mut.Instances)
		return
	} else if ok = n != 0; !ok {
		parser.errorOn(EmptyMutualBlock, mut)
		return
	}
	return mut, ok
}

func (parser *Parser) parseTopMutual() (ok bool) {
	var mut MutualBlockElem
	if mut, ok = parser.parseTopMutualHelper(); ok {
		return parser.write(mut)
	}
	return
}

func (parser *Parser) parseTopMutualHelper() (mut MutualBlockElem, ok bool) {
	returns := parser.saveStacks()
	defer returns()

	if mut, ok = parser.parseMutualInner(token.Mutual); !ok {
		return
	}

	n := 0
	mut.Types, mut.Traits, mut.Instances, mut.Declarations, mut.Bindings, n = parser.stackValues_recorded()
	if ok = n != 0; !ok {
		parser.errorOn(EmptyMutualBlock, mut)
		return
	}
	return
}

func (parser *Parser) firstPass() (ok bool) {
	// collect all the tokens into syntactically significant parts. The collected parts should have
	// indentations and the likes removed--basically, remove stuff that is only intended to denote
	// syntactic parts or rules (key-words/-symbols, indentation, annotations, etc)

	ok = true
	for end := false; !end && ok; {
		parser.inTop = true // set back to true each time

		parser.drop() // TODO: is this necessary?

		switch next := parser.Peek(); next.Type {
		case token.At:
			ok = parser.parseAnnotation()
		case token.Public:
			ok = parser.parseVisibility(next.Type, Public)
		case token.Open:
			ok = parser.parseVisibility(next.Type, Open)
		case token.Mutual:
			ok = parser.parseTopMutual()
		case token.Automatic:
			panic("TODO: implement 'automatic'")
		case token.Alias:
			panic("TODO: implement 'alias'")
		case token.Spec:
			ok = parser.parseSpec()
		case token.Id, token.LeftParen:
			ok = parser.parseDeclarationOrDefinition()
		case token.EndOfTokens:
			end = true
		case token.Hole:
			ok, end = false, true
			parser.error(IllegalNonExprPosHole)
		default:
			parser.error(UnexpectedToken)
			ok, end = false, true
		}
	}
	return ok
}
