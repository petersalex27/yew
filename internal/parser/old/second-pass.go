// =================================================================================================
// Alex Peters - 2024
// =================================================================================================
package parser

import (
	"fmt"
	"strings"

	"github.com/petersalex27/yew/internal/token"
	"github.com/petersalex27/yew/internal/types"
)

func makeIdent(t token.Token) Ident {
	return Ident{t.Value, t.Start, t.End}
}

// validator for processed constraints
var pv_constraint = processedValidation{
	[]NodeType{tupleType, applicationType, identType},
	func(t termElem) string {
		if t.NodeType == listingType {
			return UnexpectedListingMaybeEnclosed(t)
		}
		return ExpectedConstraint
	},
}

// validator for processed function defs
var pv_functionDef = processedValidation{
	[]NodeType{applicationType, identType},
	// report a more specific error message
	func(termElem) string {
		return FunctionPatternExpected
	},
}

func (parser *Parser) typeSigParsing() func() {
	old := parser.parsingTypeSig
	parser.parsingTypeSig = true
	return func() {
		parser.parsingTypeSig = old
	}
}

const tupleConstName string = ","

func isTuple(t types.Term) bool {
	return types.GetConstant(t).C == tupleConstName
}

// convert constraint to tuple if it is not already
//
// for example ...
//
//	Num a => a -> b -> c
//
// is equivalent to ...
//
//	(Num a,) => a -> b -> c
//
// and ...
//
//	(Num a, Show b) => a -> b -> c
//
// is just equivalent to itself
//
// ASSUMPTION: c is a well-formed constraint
//
// this function should only be called from ParseConstraint
func helper_constraintAsTuple(c termElem) (tup Tuple) {
	terms := make([]types.Term, 0, 2)
	t := c.Term
	for isTuple(t) {
		terms := t.(types.Application).GetTerms()
		if len(terms) != 3 { // (,) x y
			panic("bug: tuple must have three terms, '(,) x y'")
		}
		left, right := terms[1], terms[2]
		terms = append(terms, left)
		t = right
	}
	terms = append(terms, t)
	tup = Tuple{terms, c.Start, c.End}
	return
}

// parseConstraint parses a constraint
func parseConstraint(parser *Parser, constraintToks []token.Token, T, u types.Type) (cT types.Type, ok bool) {
	var constraint termElem
	constraint, ok = parser.ProcessAndValidate(constraintActions, constraintToks, pv_constraint)
	if !ok {
		return
	}

	tup := helper_constraintAsTuple(constraint)
	cT = T
	if !types.SetKind(&cT, u) {
		parser.transferEnvErrors()
		return nil, false
	}
	// now, using the tuple of constraints, create a new type
	//
	// 		(Num a, Show b) => a -> b -> c
	// becomes ...
	//		{Num a} -> {Show b} -> a -> b -> c
	for i := len(tup.Elements) - 1; i >= 0; i-- {
		var app types.Application
		var intro types.PiIntro
		app, ok = tup.Elements[i].(types.Application)
		if !ok {
			parser.errorOn(ExpectedConstrainedType, tup.Elements[i])
			return nil, false
		}
		// take right constraint and abstract the current type with it
		//		Show b => a -> b -> c
		//	becomes ...
		//		{Show b} -> a -> b -> c
		A := types.AsTyping(app)
		intro, ok = parser.env.ImplicitProd(A)
		if !ok {
			parser.transferEnvErrors()
			return nil, false
		}
		cT, ok = intro(cT)
		if !ok {
			parser.transferEnvErrors()
			return nil, false
		}
	}
	return cT, true
}

// func (parser *Parser) beginLocals(initialCap int) func(*Parser) {
// 	parser.locals.IncreaseN(initialCap)
// 	return func(p *Parser) {
// 		p.locals = nil
// 	}
// }

// parses constraint (if one exists) and adds it to the type
func (parser *Parser) addConstraint(constraint []token.Token, ty termElem) (_ termElem, ok bool) {
	var T types.Type
	if T, ok = ty.Term.(types.Type); !ok {
		parser.errorOn(ExpectedType, ty)
		return
	}

	// check if type has a constraint
	if len(constraint) == 0 {
		// no constraint to parse, just return
		return ty, true
	}

	u := types.GetKind(&ty.Term) // kind of type

	// parse constraint
	if T, ok = parseConstraint(parser, constraint, T, u); !ok {
		return
	}

	ty = termElem{typingType, T, ty.termInfo, ty.Start, ty.End}
	return ty, true
}

// Parse parses a type element
func (typ TypeElem) Parse(parser *Parser) (ok bool) {
	defer parser.typeSigParsing()()

	// parse type
	var ty termElem
	if ty, ok = parser.Process(typingActions, typ.Type); !ok {
		return
	}

	// add constraint to type
	if ty, ok = parser.addConstraint(typ.Constraint, ty); !ok {
		return
	}

	parser.termPasser.Push(ty)
	return true
}

// createPatternForFunction creates a simple version of a function for pattern matching and translation
func (parser *Parser) createPatternForFunction(name types.Constant, A types.Type, arity uint32) types.Term {
	if arity == 0 {
		// 0 arg function (constant)
		return name
	}

	// create a simple version for pattern matching based on the declaration and the translation step
	//
	// for example, given the declaration ...
	//		fun : a -> b -> c
	// the simple version would be ...
	//		fun := \x, y -> `fun` x y
	ts := make([]types.Term, arity+1)
	vs := make([]types.Variable, arity)
	ts[0] = name
	for i := 0; i < int(arity); i++ {
		vs[i] = Var(fmt.Sprintf("x%d", i))
		ts[i+1] = vs[i]
	}

	// get the terminal type of the declaration
	pi, ok := A.(types.Pi)
	if !ok {
		panic("bug: expected pi type since arity > 0")
	}
	t := pi.GetTerminal()

	// create application body
	app := types.MakeApplication(t, ts...)
	// abstract application
	lam, ok := parser.env.AutoAbstract2(app, pi, vs)
	if !ok {
		panic("bug: failed to abstract application")
	}
	return lam
}

func (dec DeclarationElem) parseDeclaration_shared(parser *Parser) (d types.Constant, A, u types.Type, arity uint32, ok bool) {
	var setType setTypeFunc
	// check if declaration is an infix function, and what its precedence and associativity are
	infixed := dec.Name.Type == token.Infix
	bp := uint8(dec.bp)
	var rAssoc uint8
	if dec.rAssoc {
		rAssoc = 1
	}

	parser.env.BeginClause()
	parser.declarations.Increase()

	// parse type signature
	old := parser.parsingTypeSig
	parser.parsingTypeSig = true
	ok = dec.Typing.Parse(parser)
	parser.parsingTypeSig = old

	decs, _ := parser.declarations.Decrease()
	loc := parser.env.EndClause(parser.inParent)

	if !ok {
		return
	}

	// get result of parsing type signature
	typ, _ := parser.termPasser.Pop()
	arity = types.CalcArity(typ.Term.(types.Type))

	// prepare symbol for declaration
	export := exports{}
	export.declTable = new(declTable)
	export.Locals = new(types.Locals)
	*export.declTable = declTable(decs)
	*export.Locals = loc
	if setType, ok = parser.declare(dec.Name, false, export); !ok {
		return
	}

	// grab the type of the declaration and that type's type
	A = typ.Term.(types.Type)    // type
	u = types.GetKind(&typ.Term) // type's type

	// declare info in parser (this is for parsing)
	ok = setType(parser, A, infixed, bp, rAssoc)
	if !ok {
		return
	}
	d = types.Constant{C: dec.Name.Value, Start: dec.Start, End: dec.End}
	return
}

func (parser *Parser) declareAndAssignInEnv(f types.Variable, A types.Type, pattern types.Term) (ok bool) {
	// declare in environment (this is for type checking)
	if ok = parser.env.Declare(f, A); !ok {
		parser.transferEnvErrors()
		return false
	}

	// assign in environment (this is for translation)
	if ok = parser.env.Assign(f, pattern); !ok {
		parser.transferEnvErrors()
		return false
	}
	return true
}

// Parse parses a declaration element
func (dec DeclarationElem) Parse(parser *Parser) (ok bool) {
	var fC types.Constant
	var A types.Type
	var arity uint32
	fC, A, _, arity, ok = dec.parseDeclaration_shared(parser)
	if !ok {
		return
	}

	// convert to variable
	f := types.Var(fC)

	// f : A
	// A : u
	// s'pose
	//		A = a -> b -> c,
	// then
	//		f := (\x, y => `f` x y)

	// declare in environment (this is for type checking)
	pattern := parser.createPatternForFunction(fC, A, arity)
	if dec.Name.Type == token.Infix {
		// if infix, declare the enclosed version of the function
		ok = parser.declareAndAssignInEnv(f, A, pattern)
		if !ok {
			return
		}

		// remove enclosing parentheses from name to prepare that version for declaration
		tok2 := dec.Name
		tok2.Value = strings.TrimPrefix(tok2.Value, "(")
		tok2.Value = strings.TrimSuffix(tok2.Value, ")")
		// update f
		f = types.Var(tok2)
	}
	ok = parser.declareAndAssignInEnv(f, A, pattern)
	if !ok {
		return
	}

	// now, if local definition exists, parse it
	if dec.localDefinition == nil {
		return true
	}

	// info was just added, so it should always be found
	decl, _ := parser.declarations.Find(f)
	ok = dec.localDefinition.Parse(parser, termElem{identType, f, *decl.termInfo, dec.Start, dec.End})
	if !ok {
		return
	}
	// now remove the definition from defining so it can't be defined in another definition clause
	def, _ := parser.defining.Pop()
	if !parser.commitDefinition(def) {
		return false
	}

	return true
}

func (parser *Parser) parseLanguageConstructs(elems []SyntacticElem) (ok bool) {
	for _, elem := range elems {
		if !elem.Parse(parser) {
			return false
		}
	}
	// commit possible final remaining definition
	if parser.defining.Empty() {
		return true
	}
	def, _ := parser.defining.Pop()
	return parser.commitDefinition(def)
}

func (mod ModuleElem) Parse(parser *Parser) (ok bool) {
	ok = parser.parseLanguageConstructs(mod.Elems)
	if !ok {
		return
	}

	// TODO: finish
	return
}

func (s ExtensionElem) Parse(parser *Parser) (ok bool) {
	panic("unimplemented")
}

func (h HaskellStyleDataConstructors) Parse(parser *Parser) (ok bool) {
	panic("unimplemented")
}

func (h HaskellStyleDataTypeElem) Parse(parser *Parser) (ok bool) {
	panic("unimplemented")
}

func helper_parseClause(parser *Parser, clause ClauseElem) (ok bool) {
	if clause == nil {
		return true
	}

	ok = clause.Parse(parser)
	return
}

func (local localDefinition) Parse(parser *Parser, head termElem) (ok bool) {
	ty := types.GetKind(&head.Term)
	_, terms := types.Split(head.Term)
	var args []types.Term
	if len(terms) >= 2 {
		args = terms[1:]
	}
	name := terms[0]
	/*decl*/ _, found := parser.declarations.Find(name)
	if !found {
		parser.errorOn(UnknownIdent, name)
		return false
	}

	// add available implicit params to the environment
	//parser.declarations.IncreaseWith(*decl.available.declTable)
	//parser.env.IncreaseWith(*decl.available.Locals)
	// remove implicit params at the end
	//defer parser.env.WithDecrease()
	//defer parser.declarations.Decrease()

	// parse whatever kind of clause (where, with, ?), putting the clause's bindings in scope
	// for the binding's body to use
	parser.env.BeginClause() // this needs to manually be closed

	parser.declarations.Increase()
	defer parser.declarations.Decrease()

	if local.Clause != nil {
		ok = helper_parseClause(parser, local.Clause)
	}
	if !ok {
		_ = parser.env.KillClause()
		return
	}

	// parse the body of the binding as an expression
	var body termElem
	body, ok = parser.ProcessExpression(localActions, local.Expression)
	if !ok {
		_ = parser.env.KillClause()
		return
	}

	locals := parser.env.EndClause(parser.inParent)
	return parser.define(name, args, ty, body, locals)
}

func (bind BindingElem) Parse(parser *Parser) (ok bool) {
	ok = true
	if len(bind.Head) == 0 {
		panic("bug: binding must have a head")
	}

	// parse the head (pattern) of the binding as an expression
	// 	- it must be parsed as an expression bc the precedence and associativity of the operators
	//		are required to determine what thing is being defined
	// 	- create bindings for the head's variables

	var head termElem
	head, ok = parser.ProcessAndValidate(scrutineeActions, bind.Head, pv_functionDef)
	if !ok {
		return
	}

	return bind.localDefinition.Parse(parser, head)
}

// any scope saving and restoring should be done in the parent node's Parse method
func (where WhereClause) Parse(parser *Parser) (ok bool) {
	for _, elem := range where.Elems {
		if !elem.Parse(parser) {
			return false
		}
	}
	return true
}

func (let LetBindingElem) Parse(parser *Parser) (ok bool) {
	return WhereClause(let).Parse(parser)
}

func (tok TokensElem) Parse(parser *Parser) (ok bool) {
	panic("bug: TokensElem should never be passed to Parse. It should be handled in the parent node's Parse method")
}

func (cons TypeConstructorElem) Parse(parser *Parser) (ok bool) {
	panic("bug: TypeConstructorElem should never be passed to Parse. It should be handled in the parent node's Parse method")
}

func (data DataTypeElem) Parse(parser *Parser) (ok bool) {
	// parse type constructor
	Z, A, u, _, good := data.TypeConstructor.parseDeclaration_shared(parser)
	if ok = good; !ok {
		return
	}

	var intro types.TypeConIntro
	intro, ok = parser.env.TypeCon(u, Z, A)
	if !ok {
		parser.transferEnvErrors()
		return
	}

	// parse and introduce data constructors
	Cs := make([]types.Constant, len(data.DataConstructors))
	Ts := make([]types.Type, len(data.DataConstructors))
	for i, dataCons := range data.DataConstructors {
		Cs[i], Ts[i], _, _, ok = dataCons.parseDeclaration_shared(parser)
		if !ok {
			return false
		}
	}
	if !intro(Cs, Ts) {
		parser.transferEnvErrors()
		return false
	}
	return true
}

func (trait SpecElem) Parse(parser *Parser) (ok bool) { return }

func (inst InstanceElem) Parse(parser *Parser) (ok bool) { return }

func (a AnnotationElem) Parse(parser *Parser) (ok bool) {
	panic("unimplemented")
}

func (m MutualBlockElem) Parse(parser *Parser) (ok bool) {
	panic("unimplemented")
}
