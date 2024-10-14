package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/common/data"
)

type (
	// a dot followed by a name, e.g., `.run`
	access struct{ data.Solo[api.Token] }

	// an annotation, e.g., `--@infixr 0 ($)` and `[@infixr 0 ($)]`
	annotation = data.Either[flatAnnotation, enclosedAnnotation]

	// a block of annotations
	annotations struct{ data.NonEmpty[annotation] }

	// a type application, e.g., `a List n`
	appType struct {
		data.Pair[typ, data.NonEmpty[typ]]
	}

	// a binding pattern of some sort for a lambda abstraction or let-in expression
	//
	//	```
	// 	binder = lower ident | upper ident | "(", {"\n"}, enc pattern, {"\n"}, ")" ;
	//	```
	binder = data.Either[ident, pattern]

	// the body of a yew source file
	body struct{ data.List[bodyElement] }

	// a body element, e.g., a function definition, typing, or spec definition
	bodyElement interface {
		api.DescribableNode
		setAnnotation(data.Maybe[annotations]) bodyElement
	}

	caseArm struct{ data.Pair[pattern, defBody] }

	caseArms struct{ data.NonEmpty[caseArm] }

	caseExpr struct{ data.Pair[pattern, caseArms] }

	// a constrained type, e.g., `Eq a => a`
	constrainedType struct {
		data.Pair[constraint, typ]
	}

	// a thing that constrains
	constrainer struct {
		data.Pair[upperIdent, pattern]
	}

	// a type constraint
	// [([..], A a), ([..], B b a), ..]
	constraint interface {
		api.DescribableNode
		asConstraint() constraint
	}

	constraintElem = data.Pair[data.List[upperIdent], constrainer]

	constraintUnverified struct{ data.Solo[typ] }

	constraintVerified struct {
		data.NonEmpty[data.Pair[data.List[upperIdent], constrainer]]
	}

	// a function definition
	def struct {
		// definition cannot have visibilities--these are applied to their declarations
		annotations data.Maybe[annotations]
		pattern     pattern
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
		data.Pair[data.Either[withClause, expr], data.Maybe[whereClause]]
	}

	// a default value for an implicit argument
	defaultExpr struct{ data.Solo[expr] }

	// a deriving clause of a data type definition
	deriving struct{ data.Solo[derivingBody] }

	//
	derivingBody struct{ data.NonEmpty[constrainer] }

	// an annotation enclosed by `[...]`
	enclosedAnnotation struct {
		data.Pair[ident, data.List[api.Node]]
	}

	// a type that is enclosed in parentheses or braces
	enclosedType struct {
		implicit bool
		typ
	}

	// express some kind of computation
	expr interface {
		api.DescribableNode
		updatePosExpr(p api.Positioned) expr
	}

	// a function application
	exprApp struct {
		data.Pair[expr, data.NonEmpty[expr]]
	}

	// an expression that can appear in a type
	exprAtom = data.Either[patternAtom, lambdaAbstraction]

	// a line-annotation, e.g., `--@infixr 0 ($)`
	flatAnnotation struct{ data.Solo[api.Token] }

	// a footer of a yew source file
	footer struct{ data.Maybe[annotations] }

	// a List of binders for a forall type
	forallBinders struct{ data.NonEmpty[ident] }

	// a forall type, e.g., `Forall a => a -> a`
	forallType struct{ data.Pair[forallBinders, typ] }

	// a function type, e.g., `Int -> Int`
	functionType struct{ data.Pair[typ, typ] }

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

	innerTypeTerms struct{ data.NonEmpty[typ] }

	innerTyping struct {
		mode   data.Maybe[modality]
		typing data.Pair[innerTypeTerms, typ]
		api.Position
	}

	lambdaAbstraction struct{ data.Pair[lambdaBinders, expr] }

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
	bindingGroupMember = data.Either[data.Pair[binder, expr], data.Pair[typing, data.Maybe[expr]]]

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
	letExpr struct{ data.Pair[letBinding, expr] }

	literal struct{ data.Solo[api.Token] }

	lowerIdent struct{ data.Solo[api.Token] }

	mainElement interface {
		api.DescribableNode
		bodyElement
		pureMainElem() mainElement
	}

	// meta struct{ annotations }

	module struct{
		annotations data.Maybe[annotations]
		name data.Solo[lowerIdent]
		api.Position
	}

	modality struct{ data.Solo[api.Token] }

	name struct{ data.Solo[api.Token] }

	importPathIdent struct{ data.Solo[api.Token] }

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
	pattern interface {
		api.DescribableNode
		updatePosPattern(p api.Positioned) pattern
	}

	// a pattern representing an application
	patternApp struct {
		data.Pair[pattern, data.NonEmpty[pattern]]
	}

	// most primitive, non-wildcard pattern
	patternAtom = data.Either[literal, patternName]

	// a pattern enclosed in braces or parentheses
	patternEnclosed struct {
		implicit bool
		data.NonEmpty[pattern]
	}

	patternName = data.Either[hole, name]

	// a raw string, e.g., `hello`
	rawString struct{ data.Solo[api.Token] }

	requiringClause = data.NonEmpty[def]

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
		dependency  data.Maybe[pattern]
		specBody    specBody
		requiring   data.Maybe[data.NonEmpty[def]]
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
	specMember = data.Either[def, typing]

	// a syntax
	//
	// for example ...
	//	```
	//	syntax `if` test `then` true `else` false = ifThenElse test true false
	//	```
	syntax struct {
		annotations data.Maybe[annotations]
		visibility  data.Maybe[visibility]
		rule        data.Pair[syntaxRule, expr]
		api.Position
	}

	syntaxRawKeyword struct{ data.Solo[rawString] }

	syntaxRule struct{ data.NonEmpty[syntaxSymbol] }
	
	syntaxRuleIdent struct {
		binding bool
		id ident
		api.Position
	}

	syntaxSymbol = data.Either[syntaxRuleIdent, syntaxRawKeyword]
	

	// a type
	typ interface {
		api.DescribableNode
		api.Positioned
		updatePosTyp(p api.Positioned) typ
	}
	// type alias
	typeAlias struct {
		annotations data.Maybe[annotations]
		visibility  data.Maybe[visibility]
		alias       data.Pair[name, typ]
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
		constructor data.Pair[name, typ]
		api.Position
	}

	// typing, e.g., `x: Int`
	typing struct {
		annotations data.Maybe[annotations]
		visibility  data.Maybe[visibility]
		typing      data.Pair[name, typ]
		api.Position
	}

	unitType struct{ data.Solo[api.Token] }

	// identifier that begins with a string matching the regex `[A-Z]`
	upperIdent struct{ data.Solo[api.Token] }

	// visibility modifier
	visibility struct{ data.Solo[api.Token] }

	visibleBodyElement interface {
		api.DescribableNode
		bodyElement
		setVisibility(v data.Maybe[visibility]) bodyElement
	}

	// where clause body
	whereBody struct{ data.NonEmpty[mainElement] }

	// where clause
	whereClause struct{ data.Solo[whereBody] }

	// '_' pattern
	wildcard struct{ data.Solo[api.Token] }

	// with clause
	withClause struct {
		data.Pair[pattern, withClauseArms]
	}

	// arm of a with clause
	withClauseArm struct{ data.Pair[withArmLhs, defBody] }

	// arms of a with clause
	withClauseArms struct{ data.NonEmpty[withClauseArm] }

	// yew source code file
	yewSource struct {
		// the meta section is primarily distinguished by semantics, but it does require a module
		// declaration to indicate its *possible* presence--which is its syntactic component
		//meta data.Maybe[meta]

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
		api.Position
	}
)

type withArmLhs = data.Either[pattern, data.Pair[pattern, pattern]]
