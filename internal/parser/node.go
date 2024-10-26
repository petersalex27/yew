package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/common/data"
	"github.com/petersalex27/yew/common/stack"
)

type TableEntry map[string]symbolic[declaration]

type Table struct {
	stack.MapStack[string, symbolic[declaration]]
}

type (
	// a dot followed by a name, e.g., `.run` or list of names separated by dots, e.g., `.run.run`
	accessList struct{ data.NonEmpty[name] }

	// an annotation, e.g., `--@infixr 0 ($)` and `[@infixr 0 ($)]`
	annotation = data.Either[flatAnnotation, enclosedAnnotation]

	// a block of annotations
	annotations struct{ data.NonEmpty[annotation] }

	// a type application, e.g., `a List n`
	TypeApp struct {
		data.Pair[Typ, data.NonEmpty[Typ]]
	}

	// a binding pattern of some sort for a lambda abstraction or let-in expression
	//
	//	```
	// 	binder = lower ident | upper ident | "(", {"\n"}, enc pattern, {"\n"}, ")" ;
	//	```
	binder = data.Either[ident, Pattern]

	// the body of a yew source file
	body struct{ data.List[bodyElement] }

	// a body element, e.g., a function definition, typing, or spec definition
	bodyElement struct {
		data.Either[Def, VisibleBodyElement]
	}

	caseArm struct {
		TableEntry
		data.Pair[Pattern, defBody]
	}

	caseArms struct{ data.NonEmpty[caseArm] }

	caseExpr struct{ data.Pair[Pattern, caseArms] }

	// a constrained type, e.g., `Eq a => a`
	constrainedType struct {
		data.Pair[constraint, Typ]
	}

	// a thing that constrains
	constrainer struct {
		data.Pair[upperIdent, constrainerParams]
	}

	constrainerParam = data.Pair[ident, data.Maybe[Typ]]

	constrainerParams struct {
		data.NonEmpty[constrainerParam]
	}

	constraintElem = data.Pair[data.List[upperIdent], constrainer]

	// a type constraint
	// 	```
	//	A a
	//	(A a, B b)
	//	(A, B, C a, D b a)
	//	```
	// etc.
	constraint struct {
		data.NonEmpty[data.Pair[data.List[upperIdent], constrainer]]
	}

	declaration = data.Either[api.Token, typing]

	// a function definition
	Def struct {
		// definition cannot have visibilities--these are applied to their declarations
		annotations data.Maybe[annotations]
		pattern     Pattern
		defBody     defBody
		api.Position
	}

	// body of a function definition: Either something that can be computed/proven or something impossible
	defBody struct {
		data.Either[impossible, defBodyPossible]
	}

	// def body that is possible to compute/prove; syntactically, this is a body not marked
	// 'impossible'
	defBodyPossible struct {
		data.Pair[data.Either[withClause, Expr], data.Maybe[whereClause]]
	}

	// a default value for an implicit argument
	defaultExpr struct{ data.Solo[Expr] }

	// a deriving clause of a data type definition
	deriving struct{ data.NonEmpty[constrainer] }

	// an annotation enclosed by `[...]`
	enclosedAnnotation struct {
		data.Pair[ident, data.List[api.Node]]
	}

	// a type that is enclosed in parentheses or braces
	enclosedType struct {
		implicit bool
		Typ
	}

	// express some kind of computation
	Expr interface {
		api.DescribableNode
		updatePosExpr(p api.Positioned) Expr
		//processExpr(YewSource, *Table, byte) (e Expr, bp byte, err error)
	}

	// a function application
	ExprApp struct {
		data.Pair[Expr, data.NonEmpty[Expr]]
	}

	// an expression that can appear in a type
	exprAtom = data.Either[patternAtom, LambdaAbstraction]

	// a line-annotation, e.g., `--@infixr 0 ($)`
	flatAnnotation struct{ data.Solo[api.Token] }

	// a footer of a yew source file
	footer struct{ data.Maybe[annotations] }

	// a List of binders for a forall type
	forallBinders struct{ data.NonEmpty[ident] }

	// a forall type, e.g., `Forall a => a -> a`
	forallType struct{ data.Pair[forallBinders, Typ] }

	// a function type, e.g., `Int -> Int`
	functionType struct{ data.Pair[Typ, Typ] }

	// a header of a yew source file
	header struct {
		data.Pair[data.Maybe[module], data.List[importStatement]]
	}

	// an alphanumeric identifier prefixed with a `?`
	hole struct{ data.Solo[api.Token] }

	// an identifier, e.g., `x` or 'X
	ident = data.Either[lowerIdent, upperIdent]

	// the `erase x : a` part of `{erase x : a}`
	implicitTyping struct {
		data.Pair[innerTyping, defaultExpr]
	}

	importStatement struct {
		data.Pair[data.Maybe[annotations], importing]
	}

	importing struct{ data.NonEmpty[packageImport] }

	impossible struct{ data.Solo[api.Token] }

	innerTypeSeq struct{ data.NonEmpty[Typ] }

	innerTypingNames = data.Either[data.Pair[modality, name], data.NonEmpty[name]]

	innerTyping struct {
		data.Pair[innerTypingNames, Typ]
	}

	LambdaAbstraction struct {
		TableEntry
		data.Pair[lambdaBinders, Expr]
	}

	// a binder specifically in the binding part of lambda abstractions
	//
	//	```
	// 	lambda binder = binder | "_" ;
	//	```
	lambdaBinder = struct{ data.Either[binder, wildcard] }

	// a sequence of lambda binders
	//
	//	```
	//	lambda binders = lambda binder, {{"\n"}, ",", {"\n"}, lambda binder}, [{"\n"}, ","] ;
	//	```
	lambdaBinders struct{ data.NonEmpty[lambdaBinder] }

	// a member of a let binding group
	//
	//	```
	//	binding group member = binder, {"\n"}, ":=", {"\n"}, expr | typing, [{"\n"}, ":=", {"\n"}, expr] ;
	//	```
	//
	// NOTE: this is just an alias, not an actual node
	bindingGroupMember = data.Either[data.Pair[binder, Expr], data.Pair[typing, data.Maybe[Expr]]]

	// a group of bindings in a let-in expression
	//
	//	```
	//	let binding =
	//		binding group member
	//		| "(", {"\n"}, binding group member, {{"\n"}, binding group member}, {"\n"}, ")" ;
	//	```
	letBinding struct {
		data.NonEmpty[bindingGroupMember]
	}

	// a let-in expression
	//
	//	```
	//	let expr = "let", {"\n"}, let binding, {"\n"}, "in", {"\n"}, expr ;
	//	```
	LetExpr struct{ data.Pair[letBinding, Expr] }

	literal interface {
		api.DescribableNode
		literal() literal
		asTyp() Typ
		asExpr() Expr
		asPattern() Pattern
	}

	aLiteral[a any] struct {
		LiteralType api.NodeType
		Value       a
		api.Position
	}

	lowerIdent struct{ data.Solo[api.Token] }

	mainElement interface {
		api.DescribableNode
		asBodyElement() bodyElement
		asMainElement() mainElement
		setAnnotation(data.Maybe[annotations]) mainElement
		//processMainElement(YewSource, *Table) (mainElement, error)
	}

	// meta struct{ annotations }

	module struct {
		annotations data.Maybe[annotations]
		name        data.Solo[lowerIdent]
		api.Position
	}

	modality struct{ data.Solo[api.Token] }

	name struct{ data.Solo[api.Token] }

	importPathIdent struct{ data.Solo[api.Token] }

	OfDef struct {
		visibility	data.Maybe[visibility]
		Def
	}

	// a package import statement, e.g.,
	//
	//		reflect/annot as ant
	// in
	//		import reflect/annot as ant
	packageImport struct {
		data.Pair[importPathIdent, data.Maybe[selections]]
	}

	// EITHER as clause OR List of exported names
	selections = data.Either[lowerIdent, data.Maybe[data.NonEmpty[name]]]

	// a representation of an expression that can be scrutinized
	Pattern interface {
		api.DescribableNode
		updatePosPattern(p api.Positioned) Pattern
		//processPattern(YewSource, *Table, byte) (pat Pattern, bp byte, err error)
		//secondPass(*Table) pattern
	}

	// a pattern representing an application
	patternApp struct {
		data.Pair[Pattern, data.NonEmpty[Pattern]]
	}

	// most primitive, non-wildcard pattern
	patternAtom = data.Either[literal, patternName]

	// a pattern enclosed in braces or parentheses
	patternEnclosed struct {
		implicit bool
		data.NonEmpty[Pattern]
	}

	patternName = data.Either[hole, name]

	// a raw string, e.g., `hello`
	rawString struct{ data.Solo[api.Token] }

	requiringClause = data.NonEmpty[Def]

	// a spec definition body
	specBody struct{ data.NonEmpty[specMember] }

	// a spec definition. Examples of various forms it can take:
	//
	// the basic form:
	//	```
	//	spec Semigroup s where (
	//		(<>) : s -> s -> s
	//	)
	//	```
	//
	// with a constraint:
	//	```
	//	spec Semigroup o => Monoid o where (
	//		neutral : o
	//	)
	//	```
	//
	// with dependency (and constraint):
	//	```
	//	spec Monad m => MonadState st m from m where (
	//		get : m st
	//		put : st -> m ()
	//		state : (st -> (st, a)) -> m a
	//	)
	//	```
	// with requiring block (and constraint):
	//	```
	//	spec Summand n => Nat n where (
	//		0 : n
	//	) requiring (
	//		-- require this case
	//		(1 +) : n -> n
	//	)
	//	```
	// note, for the last example, a require block can be used for more than this; this is just the
	// first example that comes to mind of a spec that has one
	specDef struct {
		annotations data.Maybe[annotations]
		visibility  data.Maybe[visibility]
		specHead    specHead
		dependency  data.Maybe[Pattern]
		specBody    specBody
		requiring   data.Maybe[data.NonEmpty[Def]]
		api.Position
	}

	specHead struct {
		data.Pair[data.Maybe[constraint], constrainer]
	}

	// A spec instance. Examples of various forms it can take:
	//
	// the basic form:
	//	```
	//	inst Monad Maybe where (
	//		Just x >>= f = f x
	//		_ >>= _ = Nothing
	//	)
	//	```
	//
	// with constraint(s):
	//	```
	//	inst Num a => Semigroup (Solo a) where (
	//		Solo x <> Solo y = Solo (x <> y)
	//	)
	//	```
	//
	// constrained and named:
	//	```
	//	inst (Cast a b, Cast b c)
	//	=> TransitiveCast = Cast a c where (
	//		cast x = cast x |> cast
	//	)
	//	```
	specInst struct {
		annotations data.Maybe[annotations]
		visibility  data.Maybe[visibility]
		// part following 'spec', Either:
		//	1. a constraint placed on (2.) or (3.)--if so, this will be a non-`nothing` value on the lhs of `head`
		//	2. the spec being instantiated, rhs of `head`, target will be `nothing`
		//  3. the name of the spec instance, rhs of `head` _AND_ `target` will not be `nothing`
		head specHead
		// optional: the spec being instantiated; if this is `nothing`, then the instantiated spec is
		// the one specified in the rhs of `head`
		target data.Maybe[constrainer] // used when the spec instance is assigned to a name
		// part following 'where'
		body specInstWhereClause
		api.Position
	}

	// a group of members of a spec instance
	specInstWhereClause = specBody

	// a member of a spec
	specMember = data.Either[Def, typing]

	// a syntax
	//
	// for example ...
	//	```
	//	syntax `if` test `then` true `else` false = ifThenElse test true false
	//	```
	syntax struct {
		annotations data.Maybe[annotations]
		visibility  data.Maybe[visibility]
		rule        data.Pair[syntaxRule, Expr]
		api.Position
	}

	syntaxRawKeyword struct{ data.Solo[rawString] }

	syntaxRule struct{ data.NonEmpty[syntaxSymbol] }

	syntaxRuleIdent struct {
		binding bool
		id      ident
		api.Position
	}

	syntaxSymbol = data.Either[syntaxRuleIdent, syntaxRawKeyword]

	// a type
	Typ interface {
		api.DescribableNode
		api.Positioned
		//kind.Typeable
		updatePosTyp(p api.Positioned) Typ
		//processTyp(YewSource, *Table, byte) (t Typ, bp byte, err error)
	}
	// type alias
	typeAlias struct {
		annotations data.Maybe[annotations]
		visibility  data.Maybe[visibility]
		alias       data.Pair[name, Typ]
		api.Position
	}

	// type definition with an optional deriving tuple
	typeDef struct {
		annotations data.Maybe[annotations]
		visibility  data.Maybe[visibility]
		typedef     data.Pair[typing, typeDefBody]
		deriving    data.Maybe[deriving]
		api.Position
	}

	// body of a type definition
	typeDefBody = data.Either[data.NonEmpty[typeConstructor], impossible]

	typeConstructor struct {
		annotations data.Maybe[annotations]
		constructor data.Pair[name, Typ]
		api.Position
	}

	// typing, e.g., `x: Int`
	typing struct {
		automatic   bool
		annotations data.Maybe[annotations]
		visibility  data.Maybe[visibility]
		typing      data.Pair[name, Typ]
		api.Position
	}

	unitType struct{ data.Solo[api.Token] }

	// identifier that begins with a string matching the regex `[A-Z]`
	upperIdent struct{ data.Solo[api.Token] }

	// visibility modifier
	visibility struct{ api.Token }

	VisibleBodyElement interface {
		api.DescribableNode
		mainElement
		//declare(p parser) data.Either[data.Ers, mainElement]
		setVisibility(v data.Maybe[visibility]) mainElement
		//processVbe(YewSource, *Table) (VisibleBodyElement, error)
	}

	// where clause
	whereClause struct {
		TableEntry
		data.NonEmpty[mainElement]
	}

	// '_' pattern
	wildcard struct{ data.Solo[api.Token] }

	// with clause
	withClause struct {
		TableEntry
		data.Pair[Pattern, withClauseArms]
	}

	withArmLhs = data.Either[Pattern, data.Pair[Pattern, Pattern]]

	// arm of a with clause
	withClauseArm struct {
		TableEntry
		data.Pair[withArmLhs, defBody]
	}

	// arms of a with clause
	withClauseArms struct{ data.NonEmpty[withClauseArm] }

	YewModule struct {
		name                 string // name of the module
		path                 string // absolute path to the directory of the module
		data.List[YewSource]        // sources in the module

		// non-syntax related data

		// all the modules required by this module
		required map[string]struct{}
	}

	// yew source code file
	YewSource struct {
		path string // absolute path to the source file

		// the header section contains the module declaration and any imports
		//
		// both are optional, but the module is required if the meta section is present; if the module
		// declaration is not present, any annotations one might have intended for the meta section
		// will instead target the body or footer--which will fail if those annotations cannot target
		// elements in those sections/that section.
		header data.Maybe[header]
		// (non-module) declarations and definitions
		//
		// this section is optional
		body data.Maybe[body]
		// zero or more annotations
		//
		// this is the only non-optional section; if nothing else, represents the terminal point of
		// the yew source--which is something
		footer footer
		// position range of the yew source file
		pos api.Position

		// non-syntax related data

		// this includes type constructors, data constructors, typings, and spec types
		symbols  TableEntry
		imported importedSymbols
		srcCode  api.SourceCode
	}
)
