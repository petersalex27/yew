//go:build test
// +build test

package parser

import (
	"testing"

	"github.com/petersalex27/yew/api"
)

// ensure all these implement bodyElement

// func Test_assert_bodyElement(*testing.T) {
// 	var (
// 		_ bodyElement = typing{}
// 		_ bodyElement = typeDef{}
// 		_ bodyElement = specDef{}
// 		_ bodyElement = specInst{}
// 		_ bodyElement = typeAlias{}
// 		_ bodyElement = syntax{}
// 		_ bodyElement = def{}
// 	)
// 	// yippee!
// }

func Test_assert_constraint(*testing.T) {
	var (
		_ constraint = constraintUnverified{}
		_ constraint = constraintVerified{}
	)
	// yippee!
}

func Test_assert_expr(*testing.T) {
	var (
		_ expr = caseExpr{}
		_ expr = literal{}
		_ expr = hole{}
		_ expr = name{}
		_ expr = access{}
		_ expr = lambdaAbstraction{}
		_ expr = letExpr{}
		// it's not a mistake that lambdaAbstraction is missing, it's included w/in exprRoot
	)
	// yippee!
}

func Test_assert_mainElement(*testing.T) {
	var (
		_ mainElement = def{}
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
		_ typ = appType{}
		_ typ = constrainedType{}
		_ typ = enclosedType{}
		_ typ = literal{}
		_ typ = hole{}
		_ typ = name{}
		_ typ = access{}
		_ typ = lambdaAbstraction{}
		_ typ = forallType{}
		_ typ = functionType{}
		_ typ = innerTypeTerms{}
		_ typ = innerTyping{}
		_ typ = implicitTyping{}
		_ typ = unitType{}
		_ typ = wildcard{}
	)
	// yippee!
}

func Test_assert_visibleBodyElement(*testing.T) {
	var (
		_ visibleBodyElement = specInst{}
		_ visibleBodyElement = specDef{}
		_ visibleBodyElement = typing{}
		_ visibleBodyElement = typeDef{}
		_ visibleBodyElement = typeAlias{}
		_ visibleBodyElement = syntax{}
	)
	// yippee!
}

func Test_assert_DescribableNode(t *testing.T) {
	var (
		_ api.DescribableNode = access{}
		_ api.DescribableNode = annotations{}
		_ api.DescribableNode = appType{}
		_ api.DescribableNode = body{}
		_ api.DescribableNode = caseArm{}
		_ api.DescribableNode = caseArms{}
		_ api.DescribableNode = caseExpr{}
		_ api.DescribableNode = constrainedType{}
		_ api.DescribableNode = constrainer{}
		_ api.DescribableNode = constraintUnverified{}
		_ api.DescribableNode = constraintVerified{}
		_ api.DescribableNode = def{}
		_ api.DescribableNode = defBody{}
		_ api.DescribableNode = defBodyPossible{}
		_ api.DescribableNode = defaultExpr{}
		_ api.DescribableNode = deriving{}
		_ api.DescribableNode = enclosedAnnotation{}
		_ api.DescribableNode = enclosedType{}
		_ api.DescribableNode = exprApp{}
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
		_ api.DescribableNode = innerTypeTerms{}
		_ api.DescribableNode = innerTyping{}
		_ api.DescribableNode = lambdaAbstraction{}
		_ api.DescribableNode = lambdaBinders{}
		_ api.DescribableNode = letBinding{}
		_ api.DescribableNode = letExpr{}
		_ api.DescribableNode = literal{}
		_ api.DescribableNode = lowerIdent{}
		//_ api.DescribableNode = meta{}
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
		_ api.DescribableNode = yewSource{}
	)
	// yippee!
}
