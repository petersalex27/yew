package parser

type position struct {
	start, end int
}

func (p position) pos() (int, int) {
	return p.start, p.end
}

type (
	node interface {
		pos() (start, end int)
	}

	enclosable[T node] interface {
		resetPos(ns, ne int) T
	}

	maybe[T any] struct{ just *T }

	yewProgram struct {
		header maybe[header]
		body   maybe[body]
		footer maybe[footer]
		position
	}

	header struct {
		meta maybe[meta]
		tail maybe[headerTail]
		position
	}

	body struct {
		bodyElement maybe[[]bodyElement]
		position
	}

	footer struct {
		annotations maybe[[]annotation]
		position
	}

	meta struct {
		env    maybe[annotation]
		embeds maybe[[]embed]
		position
	}

	headerTail struct {
		module  maybe[name]
		imports maybe[[]imprt]
		position
	}

	packageIdent = name

	publicMember = name

	imprt struct {
		// 'yew/x/reflect' in 'import yew/x/reflect.( x, y, z )'
		name []packageIdent
		// 'x, y, z' in 'import yew/x/reflect.( x, y, z )'
		members maybe[[]publicMember]
		position
	}

	embed struct {
		ident string
		lang  string
		raw   string
		position
	}

	bodyElement interface {
		node
	}

	mainElement interface {
		node
	}

	def struct {
		doc maybe[annotation]
		nameAppPattern
		defTail maybe[[]defTail]
		position
	}

	typeDef struct {
		position
	}

	typeAlias struct {
		position
	}

	specDef struct {
		position
	}

	specInst struct {
		position
	}

	mutualBody struct {
		bodyElements maybe[[]bodyElement]
		position
	}

	typ interface {
		node
		visitFunction(t function[typ], parser *Parser) maybe[function[typ]]
		visitTuple(t tuple[typ], parser *Parser) maybe[tuple[typ]]
		visitIdentTyping(ti identTyping, parser *Parser) maybe[identTyping]
		visitSpecTyping(st specTyping, parser *Parser) maybe[specTyping]
		visitTyping(t typing, parser *Parser) maybe[typing]
		visitForall(fa forall, parser *Parser) maybe[forall]
	}

	identTyping struct {
		ident ident
		typ   typ
		position
	}

	ident interface {
		node
		visitAnnotation(a annotation, parser *Parser) maybe[annotation]
		visitIdentTyping(ti identTyping, parser *Parser) maybe[identTyping]
		visitForall(fa forall, parser *Parser) maybe[forall]
		visitLambdaBinder(lb lambdaBinder, parser *Parser) maybe[lambdaBinder]
		visitBindingAssignment(ba bindingAssignment, parser *Parser) maybe[bindingAssignment]
		visitPrefixedName(pn prefixedName, parser *Parser) maybe[prefixedName]
	}

	literal interface {
		node
		visitPattern(p pattern, parser *Parser) maybe[pattern]
	}

	literalValue struct {
		literal literal
		position
	}

	wildcard struct{ position }

	prefixedName name

	lambdaBinder struct {
		ident ident
		maybe[typ]
		position
	}

	forall struct {
		binding list[ident]
		typ     typ
		position
	}

	specTyping struct {
		specIdent
		typ
	}

	specIdent interface {
		visitSpecTyping(st specTyping, parser *Parser) maybe[specTyping]
		visitDerivesElemGroup(deg []derivesElem, parser *Parser) maybe[[]derivesElem]
	}

	derivesElem struct {
		position
	}

	infix[T ident] struct {
		ident T
		position
	}

	lowerIdent name

	upperIdent name

	symbol name

	typing struct {
		name
		typ
		position
	}

	explicitTyping struct {
		name
		explicitTyp
		position
	}

	function[T node] [2]T

	explicitTyp interface {
		node
		visitExplicitFunction(f function[explicitTyp], parser *Parser) maybe[function[explicitTyp]]
		visitExplicitTuple(t tuple[explicitTyp], parser *Parser) maybe[tuple[explicitTyp]]
		visitExplicitTyping(t explicitTyping, parser *Parser) maybe[explicitTyping]
		visitImplicitType(it implicitTyp, parser *Parser) maybe[implicitTyp]
	}

	tuple[T node] struct {
		fst, snd T
		position
	}

	app[T node] struct {
		elems []T
		position
	}

	list[T node] []T

	pattern interface {
		node
		typ
		explicitTyp
		expression
		visitPatternApp(pp app[pattern], parser *Parser) maybe[app[pattern]]
		visitPatternTuple(pp tuple[pattern], parser *Parser) maybe[tuple[pattern]]
		visitPatternList(ps list[pattern], parser *Parser) maybe[list[pattern]]
		visitNameAppPattern(pp nameAppPattern, parser *Parser) maybe[nameAppPattern]
		visitTyp(t typ, parser *Parser) maybe[typ]
		visitExplicitType(t explicitTyp, parser *Parser) maybe[explicitTyp]
	}

	expression interface {
		node
		visitExpressionApp(ee app[expression], parser *Parser) maybe[app[expression]]
		visitExpressionTuple(ee tuple[expression], parser *Parser) maybe[tuple[expression]]
		visitExpressionList(es list[expression], parser *Parser) maybe[list[expression]]
		visitLambdaAbstraction(abs lambdaAbstraction, parser *Parser) maybe[lambdaAbstraction]
		visitBindingAssignment(ba bindingAssignment, parser *Parser) maybe[bindingAssignment]
		visitAssignment(a assignment, parser *Parser) maybe[assignment]
		visitCaseArm(ca caseArm, parser *Parser) maybe[caseArm]
		visitWithArm(wa withArm, parser *Parser) maybe[withArm]
		visitImplicitType(it implicitTyp, parser *Parser) maybe[implicitTyp]
	}

	lambdaAbstraction struct {
		position
	}

	bindingAssignment struct {
		position
	}

	caseArm struct {
		position
	}

	implicitTyp struct {
		position
	}

	nameAppPattern struct {
		pattern maybe[[]pattern]
		position
	}

	defTail interface {
		node
		visitDef()
	}

	withClause struct {
		view maybe[expression]
		arms maybe[[]withArm]
		position
	}

	assignment struct {
		expression
		position
	}

	whereClause struct {
		position
	}

	// rule:
	//	`with clause arm = [ lhs, keyword bar ], pattern, ( ( keyword equal, expression ) | with clause ) ;`
	withArm struct {
		// 'f' in 'f | e ...
		lhs maybe[withArmLhs]
		// 'e' in 'f | e ...'
		rhs pattern
		// '...' in 'f | e ...'
		tail maybe[withArmTail]
		position
	}

	// rule:
	//	`lhs = name app pattern | ( open paren, name app pattern, close paren ) | "_" ;`
	withArmLhs struct {
		// if nameAppPattern is nothing, then this struct represents a wildcard '_'
		nameAppPattern maybe[nameAppPattern]
		position
	}

	withArmTail struct {
		isWithClause bool
		// '= m' in 'f | e = m'
		expression
		// 'with x of ...' in 'f | e with x of ...'
		withClause
		position
	}

	name struct {
		ident string
		position
	}

	annotation struct {
		id    ident
		value string
		position
	}
)

func unit[A node](a A) maybe[A] {
	return maybe[A]{just: &a}
}

func (ma maybe[A]) bind(f func(A) maybe[A]) maybe[A] {
	if ma.just == nil {
		return ma
	}
	return f(*ma.just)
}
