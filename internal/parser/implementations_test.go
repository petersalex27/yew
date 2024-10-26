//go:build test
// +build test

package parser

import (
	"testing"

	"github.com/petersalex27/yew/api"
)

func Test_assert_expr(*testing.T) {
	var (
		_ Expr = caseExpr{}
		_ Expr = aLiteral[any]{}
		_ Expr = hole{}
		_ Expr = name{}
		_ Expr = accessList{}
		_ Expr = LambdaAbstraction{}
		_ Expr = LetExpr{}
		_ Expr = ExprApp{}
		// it's not a mistake that lambdaAbstraction is missing, it's included w/in exprRoot
	)
	// yippee!
}

func Test_assert_literal(*testing.T) {
	var (
		_ literal = aLiteral[any]{}
	)
	// yippee!
}

func Test_assert_mainElement(*testing.T) {
	var (
		_ mainElement = Def{}
		_ mainElement = specInst{}
		_ mainElement = specDef{}
		_ mainElement = syntax{}
		_ mainElement = typeAlias{}
		_ mainElement = typeDef{}
		_ mainElement = typing{}
	)
	// yippee!
}

func Test_assert_typ(*testing.T) {
	var (
		_ Typ = TypeApp{}
		_ Typ = constrainedType{}
		_ Typ = enclosedType{}
		_ Typ = aLiteral[any]{}
		_ Typ = hole{}
		_ Typ = name{}
		_ Typ = accessList{}
		_ Typ = LambdaAbstraction{}
		_ Typ = forallType{}
		_ Typ = functionType{}
		_ Typ = innerTypeSeq{}
		_ Typ = innerTyping{}
		_ Typ = implicitTyping{}
		_ Typ = unitType{}
		_ Typ = wildcard{}
	)
	// yippee!
}

func Test_assert_pattern(*testing.T) {
	var (
		_ Pattern = patternApp{}
		_ Pattern = patternEnclosed{}
		_ Pattern = aLiteral[any]{}
		_ Pattern = hole{}
		_ Pattern = name{}
		_ Pattern = wildcard{}
		_ Pattern = accessList{}
	)
	// yippee!
}

func Test_assert_visibleBodyElement(*testing.T) {
	var (
		_ VisibleBodyElement = OfDef{}
		_ VisibleBodyElement = specInst{}
		_ VisibleBodyElement = specDef{}
		_ VisibleBodyElement = typing{}
		_ VisibleBodyElement = typeDef{}
		_ VisibleBodyElement = typeAlias{}
		_ VisibleBodyElement = syntax{}
	)
	// yippee!
}

func Test_assert_DescribableNode(t *testing.T) {
	var (
		_ api.DescribableNode = accessList{}
		_ api.DescribableNode = annotations{}
		_ api.DescribableNode = TypeApp{}
		_ api.DescribableNode = body{}
		_ api.DescribableNode = caseArm{}
		_ api.DescribableNode = caseArms{}
		_ api.DescribableNode = caseExpr{}
		_ api.DescribableNode = constrainedType{}
		_ api.DescribableNode = constrainer{}
		_ api.DescribableNode = constraint{}
		_ api.DescribableNode = Def{}
		_ api.DescribableNode = defBody{}
		_ api.DescribableNode = defBodyPossible{}
		_ api.DescribableNode = defaultExpr{}
		_ api.DescribableNode = deriving{}
		_ api.DescribableNode = enclosedAnnotation{}
		_ api.DescribableNode = enclosedType{}
		_ api.DescribableNode = ExprApp{}
		_ api.DescribableNode = flatAnnotation{}
		_ api.DescribableNode = footer{}
		_ api.DescribableNode = forallBinders{}
		_ api.DescribableNode = forallType{}
		_ api.DescribableNode = functionType{}
		_ api.DescribableNode = header{}
		_ api.DescribableNode = hole{}
		_ api.DescribableNode = implicitTyping{}
		_ api.DescribableNode = importing{}
		_ api.DescribableNode = impossible{}
		_ api.DescribableNode = innerTypeSeq{}
		_ api.DescribableNode = innerTyping{}
		_ api.DescribableNode = LambdaAbstraction{}
		_ api.DescribableNode = lambdaBinders{}
		_ api.DescribableNode = letBinding{}
		_ api.DescribableNode = LetExpr{}
		_ api.DescribableNode = aLiteral[any]{}
		_ api.DescribableNode = lowerIdent{}
		_ api.DescribableNode = module{}
		_ api.DescribableNode = modality{}
		_ api.DescribableNode = name{}
		_ api.DescribableNode = importPathIdent{}
		_ api.DescribableNode = packageImport{}
		_ api.DescribableNode = patternApp{}
		_ api.DescribableNode = patternEnclosed{}
		_ api.DescribableNode = rawString{}
		_ api.DescribableNode = specBody{}
		_ api.DescribableNode = specDef{}
		_ api.DescribableNode = specHead{}
		_ api.DescribableNode = specInst{}
		_ api.DescribableNode = specInstWhereClause{}
		_ api.DescribableNode = syntax{}
		_ api.DescribableNode = syntaxRawKeyword{}
		_ api.DescribableNode = syntaxRule{}
		_ api.DescribableNode = typeAlias{}
		_ api.DescribableNode = typeDef{}
		_ api.DescribableNode = typeConstructor{}
		_ api.DescribableNode = typing{}
		_ api.DescribableNode = unitType{}
		_ api.DescribableNode = upperIdent{}
		_ api.DescribableNode = visibility{}
		_ api.DescribableNode = whereClause{}
		_ api.DescribableNode = wildcard{}
		_ api.DescribableNode = withClause{}
		_ api.DescribableNode = withClauseArm{}
		_ api.DescribableNode = withClauseArms{}
		_ api.DescribableNode = YewSource{}
	)
	// yippee!
}
