package parser

import (
	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/common/data"
	"github.com/petersalex27/yew/internal/kind"
	"github.com/petersalex27/yew/internal/symbol"
)

func holeAsPatternAtom(h api.Token) patternAtom {
	return data.Inr[literal](data.Inl[name](data.EOne[hole](h)))
}

func makeLiteral[a any](ty api.NodeType, val a) aLiteral[a] {
	return aLiteral[a]{ty, val, api.ZeroPosition()}
}

func nameAsPatternAtom(n name) patternAtom {
	return data.Inr[literal](data.Inr[hole](n))
}

func constructLambdaAbstraction(backslash api.Token, table map[string]symbolic[declaration], binders lambdaBinders, arrow api.Token) func(Expr) data.Either[data.Ers, LambdaAbstraction] {
	return func(e Expr) data.Either[data.Ers, LambdaAbstraction] {
		la := makeAbstraction(table, binders, e)
		la.Position = la.Update(backslash)
		la.Position = la.Update(arrow)
		return data.Ok(la)
	}
}

func assembleLetExpr(let api.Token, binders letBinding, in api.Token, expr Expr) LetExpr {
	letE := data.EMakePair[LetExpr](binders, expr)
	letE.Position = letE.Update(let)
	letE.Position = letE.Update(in)
	return letE
}

func constructAlias(aliasToken api.Token, n name, equalToken api.Token, ty Typ) typeAlias {
	alias := makeAlias(n, ty)
	alias.Position = alias.Position.Update(aliasToken).Update(equalToken)
	return alias
}

func makeFunc(lhs Typ, rhs Typ) Typ {
	return Typ(functionType{data.MakePair(lhs, rhs)})
}

type specRequiring = data.NonEmpty[Def]

func makeSpecDef(head specHead, dep data.Maybe[Pattern], body specBody, req data.Maybe[specRequiring]) specDef {
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

func makeSyntax(rule syntaxRule, e Expr) syntax {
	return syntax{
		data.Nothing[annotations](),
		data.Nothing[visibility](),
		data.MakePair(rule, e),
		api.WeakenRangeOver[api.Node](rule, e),
	}
}

func makeAlias(n name, ty Typ) typeAlias {
	return typeAlias{
		data.Nothing[annotations](),
		data.Nothing[visibility](),
		data.MakePair(n, ty),
		api.WeakenRangeOver[api.Node](n, ty),
	}
}

func makeTyping(n name, ty Typ) typing {
	return typing{
		false,
		data.Nothing[annotations](),
		data.Nothing[visibility](),
		data.MakePair(n, ty),
		api.WeakenRangeOver[api.Node](n, ty),
	}
}

func makeDef(pat Pattern, db defBody) Def {
	return Def{
		data.Nothing[annotations](),
		pat,
		db,
		api.WeakenRangeOver[api.Node](pat, db),
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

func makeCons(n name, ty Typ) typeConstructor {
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

func makeEmptyYewSource() YewSource {
	return makeYewSource(
		"",
		data.Nothing[header](),
		data.Nothing[body](),
		data.Nothing[annotations](),
	)
}

func constructConstructor(as data.Maybe[annotations], colon api.Token, ty Typ) func(n name) typeConstructor {
	return func(n name) typeConstructor {
		tc := makeCons(n, ty)
		tc.Position = tc.Update(colon)
		(&tc).annotate(as)
		return tc
	}
}

func weakenVisibility(vis data.Maybe[visibility]) data.Maybe[visibility] {
	return data.Bind(vis, func(v visibility) data.Maybe[visibility] {
		if token.Open.Match(v.Token) {
			return data.Just(visibility{token.Public.Make()})
		}
		return data.Nothing[visibility]()
	})
}

func constructorSig(vis data.Maybe[visibility], cons typeConstructor) typing {
	sig := makeTyping(cons.constructor.Fst(), cons.constructor.Snd())
	sig.annotations = cons.annotations
	sig.visibility = weakenVisibility(vis) // visibility one below type def
	return sig
}

func makeWithClause(table TableEntry, pat Pattern, arms withClauseArms) withClause {
	return withClause{table, data.MakePair(pat, arms)}
}

func makeWithArmLhs(pat Pattern) withArmLhs {
	return data.Inl[data.Pair[Pattern, Pattern]](pat)
}

func makeWithArmLhsRefined(pat1, pat2 Pattern) withArmLhs {
	return data.Inr[Pattern](data.MakePair(pat1, pat2))
}

func makeWithClauseArm(table map[string]symbolic[declaration], lhs withArmLhs, db defBody) withClauseArm {
	return withClauseArm{table, data.MakePair(lhs, db)}
}

func makeModule(id lowerIdent) module {
	return module{data.Nothing[annotations](), data.One(id), id.Position}
}

func makeCaseArm(table map[string]symbolic[declaration], pat Pattern, db defBody) caseArm {
	return caseArm{table, data.MakePair(pat, db)}
}

func makeWhereClause(table map[string]symbolic[declaration], elems data.NonEmpty[mainElement]) whereClause {
	return whereClause{table, elems}
}

func makeAbstraction(table map[string]symbolic[declaration], binders lambdaBinders, e Expr) LambdaAbstraction {
	return LambdaAbstraction{table, data.MakePair(binders, e)}
}

func makeYewPackageSource(path string, name lowerIdent, lang language, body data.List[packageBodyElement], mods data.List[YewModule]) YewPackageSource {
	modMap := make(map[string]YewModule)
	for _, mod := range mods.Elements() {
		modMap[mod.path] = mod
	}

	return YewPackageSource{
		name, lang, body,

		api.WeakenRangeOver[api.Positioned](name, lang, body),

		modMap,
		make(map[string]map[string]symbol.Symbol),
		make(map[string]map[string]kind.Symbol),
		path,
	}
}

func makeOfDef(def Def) OfDef {
	return OfDef{data.Nothing[visibility](), def}
}
