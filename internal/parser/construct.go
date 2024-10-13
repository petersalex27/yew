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
	alias := data.MakePair(n, ty)
	return typeAlias{
		alias:       alias,
		annotations: data.Nothing[annotations](aliasToken),
		visibility:  data.Nothing[visibility](aliasToken),
		Position:    api.WeakenRangeOver[api.Node](aliasToken, ty, alias, equalToken),
	}
}

func makeFunc(lhs typ, rhs typ) typ {
	return typ(functionType{data.MakePair(lhs, rhs)})
}