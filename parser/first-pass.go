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
	"fmt"
	"strings"

	"github.com/petersalex27/yew/token"
)

type (
	SyntacticElem interface {
		Parse(*Parser) (ok bool)
		Pos() (start, end int)
		String() string
	}

	SymbolicElem interface {
		SyntacticElem
		GetSyntaxClass() SyntaxClass
	}

	ClauseElem interface {
		SyntacticElem
		_c_()
	}

	ExpressionElem interface {
		SyntacticElem
		_e_()
	}
)

type Visibility uint8

const (
	Private Visibility = 0b00 // default visibility
	Public  Visibility = 0b01
	Open    Visibility = 0b10
)

type ModuleElement struct {
	Visibility
	Declaration DeclarationElem
}

func (v Visibility) String() string {
	switch v {
	case Public:
		return "public"
	case Open:
		return "open"
	default:
		return ""
	}
}

type (
	ModuleElem struct {
		Name          token.Token
		Imports       ImportTable
		TopLevelDecls map[string]SyntacticElem
		Start, End    int
	}

	TypeElem struct {
		Constraint []token.Token
		Type       []token.Token
		Start, End int
	}

	DeclarationElem struct {
		annotation AnnotationElem
		Vis        Visibility
		Name       token.Token
		Typing     TypeElem
		Start, End int
	}

	BindingElem struct {
		Head       []token.Token
		Expression []ExpressionElem
		Clause     ClauseElem
		Start      int
		// End is the same as Where.End
	}

	WhereClause struct {
		// should only have declarations, bindings, and data types
		Elems      []SymbolicElem
		Start, End int
	}

	MutualBlockElem struct {
		// should only have declarations, bindings, and data types
		Traits       []TraitElem
		Types        []DataTypeElem
		Instances    []InstanceElem
		Declarations []DeclarationElem
		Bindings     []BindingElem
		atTop        bool
		Start, End   int
	}

	// just the let binding part (doesn't include 'in' and onwards). Example:
	//	let x = 3
	// or
	//	let
	//	  x = 3
	//	  y = x + 1
	LetBindingElem WhereClause

	TypeConstructorElem struct {
		DeclarationElem
		//Name   token.Token
		//Args   []token.Token
		//Typing TypeElem
		// Start = Name.Start
		// End = Typing.End
	}

	// a type constructor paired with zero or more data constructors for that type
	//
	//	List : * -> Uint -> * where
	//	  Nil : List a 0
	//	  Cons : a -> List a n -> List a (n+1)
	//
	//	Nat : * where
	//	  Zero : Nat
	//	  Succ : Nat -> Nat
	//
	// `Nat` could also be defined as ...
	//
	//	Nat : * where
	//	  Zero
	//	  Succ Nat
	DataTypeElem struct {
		TypeConstructor  TypeConstructorElem
		DataConstructors []DeclarationElem
		End              int
		// Start = TypeConstructor.Name.Start
	}

	HaskellStyleDataConstructors struct {
		Vis        Visibility
		Name       token.Token
		Params     []TypeElem
		Start, End int
	}

	HaskellStyleDataTypeElem struct {
		TypeConstructorElem TypeConstructorElem
		DataConstructors    []HaskellStyleDataConstructors
		End                 int
	}

	// `trait` { constraint } IDENT { IDENT } `where` decl { decl }
	TraitElem struct {
		// { constraint } IDENT { IDENT }
		Constraint         []token.Token
		Head               []token.Token
		MethodDeclarations []DeclarationElem
		Start, End         int
	}

	// IDENT `of` { type } `where` def { def }
	InstanceElem struct {
		TraitName   token.Token
		Instance    []token.Token
		Definitions []BindingElem
		End         int
		// Start = TraitName.Start
	}

	EnclosedElem struct {
		Expression []ExpressionElem
		Start, End int
	}

	TokenElem token.Token

	ExtensionElem struct {
		// TODO: implement
		Start, End int
	}

	AnnotationElem struct {
		Name token.Token
		Args EnclosedElem
	}
)

func (TokenElem) _e_()      {}
func (LetBindingElem) _e_() {}
func (EnclosedElem) _e_()   {}

func (WhereClause) _c_()     {}
func (MutualBlockElem) _c_() {}

func (decl DeclarationElem) GetSyntaxClass() SyntaxClass {
	return DeclClass
}

func (DataTypeElem) GetSyntaxClass() SyntaxClass {
	return TypeClass
}

func (TraitElem) GetSyntaxClass() SyntaxClass {
	return TraitClass
}

func (InstanceElem) GetSyntaxClass() SyntaxClass {
	return InstanceClass
}

func (BindingElem) GetSyntaxClass() SyntaxClass {
	return FuncClass
}

func (TypeConstructorElem) GetSyntaxClass() SyntaxClass {
	return TypeConsClass
}

func (ExtensionElem) GetSyntaxClass() SyntaxClass {
	return ExtensionClass
}

func (MutualBlockElem) GetSyntaxClass() SyntaxClass {
	return MutualClass
}

func joinTokensStringed(tokens []token.Token, sep string) string {
	var b strings.Builder
	switch len(tokens) {
	case 0:
		return ""
	}

	b.WriteString(tokens[0].Value)
	for _, tok := range tokens[1:] {
		b.WriteString(sep)
		b.WriteString(tok.Value)
	}
	return b.String()
}

func joinElemsStringed[T SyntacticElem](elems []T, sep string) string {
	var b strings.Builder
	switch len(elems) {
	case 0:
		return ""
	}

	b.WriteString(elems[0].String())
	for _, elem := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(elem.String())
	}
	return b.String()
}

func (dec DeclarationElem) Pos() (start, end int) {
	return dec.Start, dec.End
}

func (dec DeclarationElem) String() string {
	var annot, vis string
	annot = dec.annotation.String()
	if annot != "" {
		annot = annot + " "
	}
	if dec.Vis != Private {
		vis = dec.Vis.String() + " "
	}
	return fmt.Sprintf("%s%s%s : %s", vis, annot, dec.Name.Value, dec.Typing.String())
}

func (typ TypeElem) Pos() (start, end int) {
	return typ.Start, typ.End
}

func (typ TypeElem) String() string {
	if len(typ.Constraint) > 0 {
		return fmt.Sprintf("%s %s", joinTokensStringed(typ.Constraint, " "), joinTokensStringed(typ.Type, " "))
	}
	return joinTokensStringed(typ.Type, " ")
}

func (bind BindingElem) Pos() (start, end int) {
	_, end = bind.Clause.Pos()
	return bind.Start, end
}

func (bind BindingElem) String() string {
	start := fmt.Sprintf("%s = %s", joinTokensStringed(bind.Head, " "), joinElemsStringed(bind.Expression, " "))
	clause := bind.Clause.String()
	if len(clause) != 0 {
		return start + " " + clause
	}
	return start
}

func (mut MutualBlockElem) Pos() (start, end int) {
	return mut.Start, mut.End
}

func (mut MutualBlockElem) String() string {
	traits := joinElemsStringed(mut.Traits, "; ")
	types := joinElemsStringed(mut.Types, "; ")
	insts := joinElemsStringed(mut.Instances, "; ")
	decls := joinElemsStringed(mut.Declarations, "; ")
	bindings := joinElemsStringed(mut.Bindings, "; ")
	return fmt.Sprintf("mutual {%s} {%s} {%s} {%s} {%s}", traits, types, insts, decls, bindings)
}

func (where WhereClause) Pos() (start, end int) {
	return where.Start, where.End
}

func (where WhereClause) String() string {
	if len(where.Elems) == 0 {
		return ""
	}
	return fmt.Sprintf("where %s", joinElemsStringed(where.Elems, "; "))
}

func (let LetBindingElem) Pos() (start, end int) {
	return let.Start, let.End
}

func (let LetBindingElem) String() string {
	return fmt.Sprintf("let %s in", joinElemsStringed(let.Elems, "; "))
}

func (tok TokenElem) Pos() (start, end int) {
	return tok.Start, tok.End
}

func (tok TokenElem) String() string {
	return tok.Value
}

func (ee EnclosedElem) Pos() (start, end int) {
	return ee.Start, ee.End
}

func (ee EnclosedElem) String() string {
	return "(" + joinElemsStringed(ee.Expression, " ") + ")"
}

func (cons TypeConstructorElem) Pos() (start, end int) {
	return cons.DeclarationElem.Pos()
}

func (cons TypeConstructorElem) String() string {
	// if len(cons.Args) > 0 {
	// 	return fmt.Sprintf("%s %s : %s", cons.Name.Value, joinTokensStringed(cons.Args, " "), cons.Typing.String())
	// }
	return cons.DeclarationElem.String()
	//return fmt.Sprintf("%s : %s", cons.Name.Value, cons.Typing.String())
}

func (data DataTypeElem) Pos() (start, end int) {
	return data.TypeConstructor.Name.Start, data.End
}

func (data DataTypeElem) String() string {
	return fmt.Sprintf("%s where {%s}", data.TypeConstructor.String(), joinElemsStringed(data.DataConstructors, "; "))
}

func (trait TraitElem) Pos() (start, end int) {
	return trait.Start, trait.End
}

func (trait TraitElem) String() string {
	var start string
	if len(trait.Constraint) > 0 {
		start = fmt.Sprintf("trait %s %s where", joinTokensStringed(trait.Constraint, " "), joinTokensStringed(trait.Head, " "))
	} else {
		start = fmt.Sprintf("trait %s where", joinTokensStringed(trait.Head, " "))
	}
	return fmt.Sprintf("%s {%s}", start, joinElemsStringed(trait.MethodDeclarations, "; "))
}

func (inst InstanceElem) Pos() (start, end int) {
	return inst.TraitName.Start, inst.End
}

func (inst InstanceElem) String() string {
	return fmt.Sprintf("%s of %s where {%s}", inst.TraitName, joinTokensStringed(inst.Instance, " "), joinElemsStringed(inst.Definitions, "; "))
}

func (a AnnotationElem) Pos() (start int, end int) {
	return a.Name.Start, a.Args.End
}

func (a AnnotationElem) String() string {
	if a.Name.Value == "" {
		return ""
	}

	if len(a.Args.Expression) == 0 {
		return "%" + a.Name.Value
	}
	// %annotation(args..)
	return fmt.Sprintf("%%%s%s", a.Name.Value, a.Args.String())
}

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

func (parser *Parser) parseTyping() (typ TypeElem, ok bool) {
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
		return typ, false
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
		return typ, false
	}

	typ.Constraint = tokens[:lastConstraintArrowEnd]
	typ.Type = tokens[lastConstraintArrowEnd:i]
	typ.Start = tokens[0].Start
	typ.End = tokens[i-1].End
	return typ, true
}

func parseDeclarationMaybeType(parser *Parser) (ok bool) {
	// create a new context and make sure it's removed when this function returns
	declaredName, cleaner, validated := parser.openSection(declarationValidator)
	defer cleaner.Clean()

	if ok = !parser.panicking; !ok {
		return
	} else if ok = validated; !ok {
		panic("bug: declarable name not found at parser's current token")
	}

	var dec DeclarationElem
	dec, ok = parser.parseDeclaration_shared(declaredName)
	if !ok {
		return
	}

	if test := parser.dropAndAdvanceGreaterIndent(); !test {
		parser.writeDecl(dec)
		return
	}

	if parser.Peek().Type != token.Where {
		parser.error(ExpectedWhere)
		return
	}

	var typ DataTypeElem
	typ.TypeConstructor.Name = dec.Name
	typ.TypeConstructor.Typing = dec.Typing
	typ.DataConstructors, ok = parser.parseJustDecls()
	if !ok {
		return
	}
	parser.writeDataType(typ)
	return
}

// assumes some kind of identifier has already been seen
func parseDeclaration(parser *Parser) (ok bool) {
	// create a new context and make sure it's removed when this function returns
	declaredName, cleaner, validated := parser.openSection(declarationValidator)
	defer cleaner.Clean()

	if ok = !parser.panicking; !ok {
		return
	} else if ok = validated; !ok {
		panic("bug: declarable name not found at parser's current token")
	}

	var dec DeclarationElem
	dec, ok = parser.parseDeclaration_shared(declaredName)
	if ok {
		parser.writeDecl(dec)
	}
	return
}

func (parser *Parser) parseDeclaration_shared(declaredName token.Token) (dec DeclarationElem, ok bool) {
	if ok = parser.dropAndAdvanceGreaterIndent(); !ok {
		parser.error(ExpectedGreaterIndent)
		return
	}

	_, ok = parser.get(token.Colon)
	if !ok {
		parser.error(ExpectedTyping)
		return
	}

	dec.Name = declaredName
	dec.Start = declaredName.Start
	dec.Typing, ok = parser.parseTyping()
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

func (parser *Parser) parseTraitHelper() (trait TraitElem, ok bool) {
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
	var trait TraitElem
	if trait, ok = parser.parseTraitHelper(); ok {
		// add trait
		parser.writeTrait(trait)
	}
	return
}

func (parser *Parser) parseTrait() (ok bool) {
	if parser.Peek().Type != token.Trait {
		panic("bug: current token in stream not validated")
	}
	again := func(p *Parser, i int) bool { return p.equalIndent(i) }
	return parser.parseSection(parseTrait, again)
}

func (parser *Parser) parseDeclaration() (ok bool) {
	switch parser.Peek().Type {
	case token.Id, token.Affixed:
		return parseDeclaration(parser)
	default:
		panic("bug: current token in stream not validated")
	}
}

func (parser *Parser) parseDeclarationMaybeType() (ok bool) {
	switch parser.Peek().Type {
	case token.Id, token.Affixed:
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
func (parser *Parser) stackValues() (types []DataTypeElem, traits []TraitElem, insts []InstanceElem, decls []DeclarationElem, bindings []BindingElem) {
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
func (parser *Parser) stackValues_recorded() (types []DataTypeElem, traits []TraitElem, insts []InstanceElem, decls []DeclarationElem, bindings []BindingElem, totalWritten int) {
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

func (parser *Parser) whereClause() (whereToken token.Token, ok bool) {
	var tok token.Token
	tok, cleaner, validated := parser.openSection(func(p *Parser) bool { return p.Peek().Type == token.Where })
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

func (parser *Parser) parseLet() (let LetBindingElem, ok bool) {
	returns := parser.saveStacks()
	defer returns()

	letToken, cleaner, validated := parser.openSection(func(p *Parser) bool { return p.Peek().Type == token.Let })
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
		parser.writeDataType(data)
	}
	return ok
}

func (parser *Parser) writeHigherOrderFunctionDecl(typeCons TypeConstructorElem) (ok bool) {
	// higher order function, push result
	higherOrderDecl := typeCons.DeclarationElem
	parser.writeDecl(higherOrderDecl)
	return true
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
	typeCons.Name = scrutinee[0]
	typeCons.Typing, ok = parser.parseTyping()
	if !ok {
		return
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

// assumes input is at left paren
func (parser *Parser) parseEnclosed() (ee EnclosedElem, ok bool) {
	tmp := parser.Advance() // move past left paren
	ee.Start = tmp.Start
	ee.Expression, _, ok = parser.parseExpression2(endOnRParen)
	if !ok {
		return
	}

	rparen := parser.Advance()
	if ok = rparen.Type == token.RightParen; !ok {
		parser.error(ExpectedRParen)
		return
	}

	ee.End = rparen.End
	return
}

var endOnWhere = map[token.Type]bool{
	token.Where:  true,
	token.Mutual: true,
	token.In:     true,
}

var endOnRParen = map[token.Type]bool{
	token.RightParen: true,
	// 'where' tokens, and any other illegal tokens, will be caught as unexpected tokens inside of paren enclosed expressions
}

func (parser *Parser) parseExpression() (exprs []ExpressionElem, clause ClauseElem, ok bool) {
	return parser.parseExpression2(endOnWhere)
}

func (parser *Parser) parseExpression2(endMap map[token.Type]bool) (exprs []ExpressionElem, clause ClauseElem, ok bool) {
	const initialCap = 32
	exprs = make([]ExpressionElem, 0, initialCap)
	i := 0

	if !parser.dropAndAdvanceGreaterIndent() {
		parser.error(ExpectedGreaterIndent)
		return exprs, clause, false
	}

	next := parser.Peek()

	// loop until endType or a exiting-indent is found
	for !endMap[next.Type] {
		var expr ExpressionElem
		if next.Type == token.EndOfTokens {
			break // this counts as an "exiting-indent"
		} else if next.Type == token.Let {
			if expr, ok = parser.parseLet(); !ok {
				return exprs, clause, false
			}
		} else if next.Type == token.Case {
			panic("TODO: implement case expressions")
		} else if next.Type == token.LeftParen {
			expr, ok = parser.parseEnclosed()
			if !ok {
				return exprs, clause, false
			}
		} else {
			expr = TokenElem(next)
		}

		i++
		exprs = append(exprs, expr)

		var end bool
		next, end, ok = parser.prepareNextExpressionIteration()
		if !ok {
			return
		} else if end {
			break
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

func (parser *Parser) parseScrutinee(endOfScrutinee token.Type) (scrutinee []token.Token, scrutinizingDataType, ok bool) {
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

	// loop until `=` is found
	for next.Type != endOfScrutinee {
		if next.Type == token.EndOfTokens {
			parser.error(UnexpectedEOF)
			return scrutinee, false, false
		} else if scrutinizingDataType = next.Type == token.Colon; scrutinizingDataType {
			return
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

	if ok = next.Type == endOfScrutinee; !ok {
		parser.error(expectedMessage(endOfScrutinee))
		return
	} else if ok = len(scrutinee) != 0; !ok {
		parser.error(ExpectedScrutinee)
		return
	}

	_ = parser.Advance() // move past end marker
	return
}

func (parser *Parser) defineScrutinee(scrutinee []token.Token) (ok bool) {
	var binding BindingElem
	binding.Head = scrutinee
	binding.Start = binding.Head[0].Start // guaranteed to have at least one element in result
	binding.Expression, binding.Clause, ok = parser.parseExpression()
	if ok {
		parser.writeBinding(binding) // save result
	}
	return ok
}

func parseScrutinee(parser *Parser) (ok bool) {
	var scrutinee []token.Token
	var scrutinizingDataType bool
	scrutinee, scrutinizingDataType, ok = parser.parseScrutinee(token.Equal)
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

func (parser *Parser) parseDeclarationOrDefinition() (ok bool) {
	before := parser.tokenPos
	first := parser.Peek()

	_ = parser.Advance()
	if !parser.dropAndAdvanceGreaterIndent() {
		parser.error(ExpectedGreaterIndent)
		return false
	}
	next := parser.Peek()
	parser.tokenPos = before // rewind

	switch next.Type {
	case token.Colon:
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
	case token.Trait:
		ok = parser.parseTrait()
	case token.Id, token.Affixed:
		ok = parser.parseDeclarationOrDefinition()
	case token.ImplicitId:
		ok = parser.parseDefinition()
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
	case token.Id, token.Affixed:
		ok = parser.parseDeclarationOrDefinition()
	case token.ImplicitId:
		ok = parser.parseDefinition()
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
	// annot := parser.Advance()
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
	mutToken, cleaner, validated := parser.openSection(func(p *Parser) bool { return p.Peek().Type == actualOpen })
	defer cleaner.Clean()

	if ok = !parser.panicking; !ok {
		return
	} else if ok = validated; !ok {
		panic("bug: correct token not found at parser's current token")
	}
	return parser.mutualInnerHelper(!cleaner.SectionOpened(), mutToken)
}

func (parser *Parser) illegalTraitErrors(ts []TraitElem) {
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
		parser.saver.elems.Push(mut)
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
		case token.Percent:
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
		case token.Trait:
			ok = parser.parseTrait()
		case token.Id, token.Affixed:
			ok = parser.parseDeclarationOrDefinition()
		case token.ImplicitId:
			// this should parse the definition for some affixed identifier
			panic("TODO: implement implicit id")
			//ok = parser.parseDefinition()
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
