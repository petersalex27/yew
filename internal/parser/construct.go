package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/common/data"
)

func holeAsPatternAtom(h api.Token) patternAtom {
	return data.Inr[literal](data.Inl[name](data.EOne[hole](h)))
}

func literalAsPatternAtom(lit api.Token) patternAtom {
	return data.Inl[patternName](data.EOne[literal](lit))
}

func nameAsPatternAtom(n name) patternAtom {
	return data.Inr[literal](data.Inr[hole](n))
}

func constructLambdaAbstraction(backslash api.Token, binders lambdaBinders, arrow api.Token) func(expr) data.Either[data.Ers, lambdaAbstraction] {
	return func(e expr) data.Either[data.Ers, lambdaAbstraction] {
		la := data.EMakePair[lambdaAbstraction](binders, e)
		la.Position = la.Update(backslash)
		la.Position = la.Update(arrow)
		return data.Ok(la)
	}
}

func assembleLetExpr(let api.Token, binders letBinding, in api.Token, expr expr) letExpr {
	letE := data.EMakePair[letExpr](binders, expr)
	letE.Position = letE.Update(let)
	letE.Position = letE.Update(in)
	return letE
}

func constructAlias(aliasToken api.Token, n name, equalToken api.Token, ty typ) typeAlias {
	alias := makeAlias(n, ty)
	alias.Position = alias.Position.Update(aliasToken).Update(equalToken)
	return alias
}

func makeFunc(lhs typ, rhs typ) typ {
	return typ(functionType{data.MakePair(lhs, rhs)})
}

type specRequiring = data.NonEmpty[def]

func makeSpecDef(head specHead, dep data.Maybe[pattern], body specBody, req data.Maybe[specRequiring]) specDef {
	return specDef{
		data.Nothing[annotations](),
		data.Nothing[visibility](),
		head,
		dep,
		body,
		req,
		api.WeakenRangeOver[api.Node](head, dep, body, req),
	}
}

func makeSpecInst(head specHead, target data.Maybe[constrainer], body specInstWhereClause) specInst {
	return specInst{
		data.Nothing[annotations](),
		data.Nothing[visibility](),
		head,
		target,
		body,
		api.WeakenRangeOver[api.Node](head, target, body),
	}
}

func makeSyntax(rule syntaxRule, e expr) syntax {
	return syntax{
		data.Nothing[annotations](),
		data.Nothing[visibility](),
		data.MakePair(rule, e),
		api.WeakenRangeOver[api.Node](rule, e),
	}
}

func makeAlias(n name, ty typ) typeAlias {
	return typeAlias{
		data.Nothing[annotations](),
		data.Nothing[visibility](),
		data.MakePair(n, ty),
		api.WeakenRangeOver[api.Node](n, ty),
	}
}

func makeTyping(n name, ty typ) typing {
	return typing{
		data.Nothing[annotations](),
		data.Nothing[visibility](),
		data.MakePair(n, ty),
		api.WeakenRangeOver[api.Node](n, ty),
	}
}

func makeTypeDef(head typing, body typeDefBody, der data.Maybe[deriving]) typeDef {
	return typeDef{
		data.Nothing[annotations](),
		data.Nothing[visibility](),
		data.MakePair(head, body),
		der,
		api.WeakenRangeOver[api.Node](head, body, der),
	}
}

func makeCons(n name, ty typ) typeConstructor {
	return typeConstructor{
		data.Nothing[annotations](),
		data.MakePair(n, ty),
		api.WeakenRangeOver[api.Node](n, ty),
	}
}

func makeBindingSyntaxRuleIdent(id ident) syntaxRuleIdent {
	return syntaxRuleIdent{true, id, id.GetPos()}
}

func makeStdSyntaxRuleIdent(id ident) syntaxRuleIdent {
	return syntaxRuleIdent{false, id, id.GetPos()}
}

func makeEmptyYewSource() yewSource {
	return makeYewSource(
		data.Nothing[header](),
		data.Nothing[body](),
		data.Nothing[annotations](),
	)
}

func constructConstructor(as data.Maybe[annotations], colon api.Token, ty typ) func(n name) typeConstructor {
	return func(n name) typeConstructor {
		tc := makeCons(n, ty)
		tc.Position = tc.Update(colon)
		(&tc).annotate(as)
		return tc
	}
}

var makeWithClause = data.EMakePair[withClause]

func makeWithArmLhs(pat pattern) withArmLhs {
	return data.Inl[data.Pair[pattern, pattern]](pat)
}

func makeWithArmLhsRefined(pat1, pat2 pattern) withArmLhs {
	return data.Inr[pattern](data.MakePair(pat1, pat2))
}

func makeWithClauseArm(lhs withArmLhs, db defBody) withClauseArm {
	return data.EMakePair[withClauseArm](lhs, db)
}