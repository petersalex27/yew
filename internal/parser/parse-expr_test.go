//go:build test
// +build test

package parser

import (
	"testing"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/common/data"
)

// rule:
//
//	```
//	expr = expr term, {expr term rhs} ;
//	```
func TestParseExpr(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  expr
	}{
		{
			"1 expr term",
			[]api.Token{id_x_tok},
			exprNode,
		},
		{
			"2 expr terms",
			[]api.Token{id_x_tok, id_x_tok},
			exprAppNode,
		},
		{
			"2 + k expr terms (tested w/ k=1)",
			[]api.Token{id_x_tok, id_x_tok, id_x_tok},
			exprAppNode2,
		},
		{
			"expr w/ access - 0",
			[]api.Token{id_x_tok, dot, id_x_tok},
			exprAppAccess,
		},
		{
			"expr w/ access - 1",
			[]api.Token{id_x_tok, dot, newline, id_x_tok},
			exprAppAccess,
		},
		{
			"expr w/ double access",
			[]api.Token{id_x_tok, dot, id_x_tok, dot, id_x_tok},
			exprAppAccessDouble,
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, ParseExpr, -1))
	}
}

// Not a real rule, just useful helper function
//
// rule:
//
//	```
//	colon equal assignment = [{"\n"}, ":=", {"\n"}, expr] ;
//	```
func TestParseMaybeColonEqualAssignment(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  data.Maybe[expr]
	}{
		{
			"assignment - 0",
			[]api.Token{colonEqual, id_x_tok},
			data.Just(exprNode),
		},
		{
			"assignment - 1",
			[]api.Token{colonEqual, newline, id_x_tok},
			data.Just(exprNode),
		},
		{
			"assignment - nothing",
			[]api.Token{},
			data.Nothing[expr](),
		},
	}

	for _, test := range tests {
		t.Run(test.name, maybeOutputFUT_endCheck(test.input, test.want, parseMaybeColonEqualAssignment, -1))
	}
}

// rule:
//
//	```
//	case arm = pattern, {"\n"}, def body thick arrow ;
//	```
func TestMaybeParseCaseArm(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  caseArm
	}{
		{
			"case arm - 0",
			[]api.Token{id_x_tok, thickArrow, id_x_tok},
			caseArmNode,
		},
		{
			"case arm - 1",
			[]api.Token{id_x_tok, newline, thickArrow, id_x_tok},
			caseArmNode,
		},
	}

	for _, test := range tests {
		t.Run(test.name, maybeOutputFUT_endCheck(test.input, data.Just(test.want), maybeParseCaseArm, -1))
	}
}

// rule:
//
//	```
//	case arms = case arm | "(", {"\n"}, case arm, {{"\n"}, case arm}, {"\n"}, ")" ;
//	```
//
// The tests are pretty small here b/c `parseCaseArms` just wraps a function more thoroughly tested
// in other tests.
//
// The point of this test is just to make sure the `maybeParseCaseArm` (no `s`) integrates with
// the aforementioned function
func TestParseCaseArms(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  caseArms
	}{
		{
			"single case arm",
			// x => x
			[]api.Token{id_x_tok, thickArrow, id_x_tok},
			data.EConstruct[caseArms](caseArmNode),
		},
		{
			"single case arm enclosed",
			// (x => x)
			[]api.Token{lparen, id_x_tok, thickArrow, id_x_tok, rparen},
			data.EConstruct[caseArms](caseArmNode),
		},
		{
			"case arms",
			/*
				(
					x => x
					x => x
				)
			*/
			[]api.Token{lparen, newline, id_x_tok, thickArrow, id_x_tok, newline, id_x_tok, thickArrow, id_x_tok, newline, rparen},
			data.EConstruct[caseArms](caseArmNode, caseArmNode),
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseCaseArms, -1))
	}
}

// rule:
//
//	```
//	case expr = "case", {"\n"}, pattern, {"\n"}, "of", {"\n"}, case arms ;
//	```
func TestParseMaybeCaseExpr(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  data.Maybe[caseExpr]
	}{
		{
			"case expr - 000",
			// case x of x => x
			[]api.Token{caseTok, id_x_tok, of, id_x_tok, thickArrow, id_x_tok},
			data.Just(caseExprNode),
		},
		{
			"case expr - 001",
			/*
				case x of
				 	x => x
			*/
			[]api.Token{caseTok, id_x_tok, of, newline, id_x_tok, thickArrow, id_x_tok},
			data.Just(caseExprNode),
		},
		{
			"case expr - 010",
			/*
				case x
					of x => x
			*/
			[]api.Token{caseTok, id_x_tok, newline, of, id_x_tok, thickArrow, id_x_tok},
			data.Just(caseExprNode),
		},
		{
			"case expr - 011",
			/*
				case x
					of
						x => x
			*/
			[]api.Token{caseTok, id_x_tok, newline, of, newline, id_x_tok, thickArrow, id_x_tok},
			data.Just(caseExprNode),
		},
		{
			"case expr - 100",
			/*
				case
					x of x => x
			*/
			[]api.Token{caseTok, newline, id_x_tok, of, id_x_tok, thickArrow, id_x_tok},
			data.Just(caseExprNode),
		},
		{
			"case expr - 101",
			/*
				case
					x of
						x => x
			*/
			[]api.Token{caseTok, newline, id_x_tok, of, newline, id_x_tok, thickArrow, id_x_tok},
			data.Just(caseExprNode),
		},
		{
			"case expr - 110",
			/*
				case
					x
					of x => x
			*/
			[]api.Token{caseTok, newline, id_x_tok, newline, of, id_x_tok, thickArrow, id_x_tok},
			data.Just(caseExprNode),
		},
		{
			"case expr - 111",
			/*
				case
					x
					of
						x => x
			*/
			[]api.Token{caseTok, id_x_tok, newline, of, newline, id_x_tok, thickArrow, id_x_tok},
			data.Just(caseExprNode),
		},
	}

	for _, test := range tests {
		t.Run(test.name, maybeOutputFUT_endCheck(test.input, test.want, parseMaybeCaseExpr, -1))
	}
}

// rule:
//
//	```
//	let binding =
//		binding group member
//		| "(", {"\n"}, binding group member, {{"\n"}, binding group member}, {"\n"}, ")" ;
//	```
func TestParseLetBinding(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  letBinding
	}{
		{
			"binding group member - binder`",
			[]api.Token{id_x_tok, colonEqual, id_x_tok},
			letBinding_b, // x := x
		},
		{
			"binding group member - typing",
			[]api.Token{id_x_tok, colon, id_x_tok},
			letBinding_t, // x : x
		},
		{
			"binding group member - assigned typing",
			[]api.Token{id_x_tok, colon, id_x_tok, colonEqual, id_x_tok},
			letBinding_a, // x : x := x
		},
		{
			"enclosed binding group member",
			[]api.Token{lparen, id_x_tok, colonEqual, id_x_tok, rparen},
			letBinding_b, // (x := x)
		},
		{
			"enclosed binding group members",
			[]api.Token{lparen, id_x_tok, colonEqual, id_x_tok, newline, id_x_tok, colon, id_x_tok, rparen},
			letBinding_bt, // (x := x\nx : x)
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseLetBinding, -1))
	}
}

func TestParseLetExpr(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  letExpr
	}{
		{
			"let expr - 000",
			[]api.Token{let, id_x_tok, colonEqual, id_x_tok, in, id_x_tok},
			letExprNode,
		},
		{
			"let expr - 001",
			[]api.Token{let, id_x_tok, colonEqual, id_x_tok, newline, in, id_x_tok},
			letExprNode,
		},
		{
			"let expr - 010",
			[]api.Token{let, id_x_tok, newline, colonEqual, id_x_tok, in, id_x_tok},
			letExprNode,
		},
		{
			"let expr - 011",
			[]api.Token{let, id_x_tok, newline, colonEqual, id_x_tok, newline, in, id_x_tok},
			letExprNode,
		},
		{
			"let expr - 100",
			[]api.Token{let, newline, id_x_tok, colonEqual, id_x_tok, in, id_x_tok},
			letExprNode,
		},
		{
			"let expr - 101",
			[]api.Token{let, newline, id_x_tok, colonEqual, id_x_tok, newline, in, id_x_tok},
			letExprNode,
		},
		{
			"let expr - 110",
			[]api.Token{let, newline, id_x_tok, newline, colonEqual, id_x_tok, in, id_x_tok},
			letExprNode,
		},
		{
			"let expr - 111",
			[]api.Token{let, newline, id_x_tok, newline, colonEqual, id_x_tok, newline, in, id_x_tok},
			letExprNode,
		},
	}

	for _, test := range tests {
		t.Run(test.name, maybeOutputFUT_endCheck(test.input, data.Just(test.want), parseMaybeLetExpr, -1))
	}
}

// rule:
//
//	```
//	expr term = [expr atom | "(", {"\n"}, enc expr, {"\n"}, ")" | let expr | case expr] ;
//	```
func TestParseMaybeExprTerm(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  data.Maybe[expr]
	}{
		{
			"expr atom",
			[]api.Token{id_x_tok},
			data.Just[expr](name_x),
		},
		{
			"enclosed expr",
			[]api.Token{lparen, id_x_tok, rparen},
			data.Just(enclosedExpr),
		},
		{
			"let expr",
			[]api.Token{let, id_x_tok, colonEqual, id_x_tok, in, id_x_tok},
			data.Just[expr](letExprNode),
		},
		{
			"case expr",
			[]api.Token{caseTok, id_x_tok, of, id_x_tok, thickArrow, id_x_tok},
			data.Just[expr](caseExprNode),
		},
		{
			"nothing",
			[]api.Token{},
			data.Nothing[expr](token.EndOfTokens.Make()),
		},
	}

	for _, test := range tests {
		t.Run(test.name, maybeOutputFUT_endCheck(test.input, test.want, parseMaybeExprTerm, -1))
	}
}

func TestMaybeEnclosedExpr(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  data.Maybe[expr]
	}{
		{
			"expr - 00",
			[]api.Token{lparen, id_x_tok, rparen},
			data.Just[expr](enclosedExpr),
		},
		{
			"expr - 01",
			[]api.Token{lparen, id_x_tok, newline, rparen},
			data.Just[expr](enclosedExpr),
		},
		{
			"expr - 10",
			[]api.Token{lparen, newline, id_x_tok, rparen},
			data.Just[expr](enclosedExpr),
		},
		{
			"expr - 11",
			[]api.Token{lparen, newline, id_x_tok, newline, rparen},
			data.Just[expr](enclosedExpr),
		},
	}

	for _, test := range tests {
		t.Run(test.name, maybeOutputFUT_endCheck(test.input, test.want, parseMaybeEnclosedExpr, -1))
	}
}

// rule:
//
//	```
//	expr atom = pattern atom | lambda abstraction ;
//	```
func TestParseExprAtom(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  exprAtom
	}{
		{
			"pattern atom",
			[]api.Token{id_x_tok},
			exprAtomNode,
		},
		{
			"lambda abstraction",
			[]api.Token{backslash, id_x_tok, thickArrow, id_x_tok},
			data.Inr[patternAtom](lambdaAbs1),
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseExprAtom, -1))
	}
}

func TestParseLambdaAbstraction(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  lambdaAbstraction
	}{
		{
			name: "no newlines, single binder - 000",
			// \x => x
			input: []api.Token{backslash, id_x_tok, thickArrow, id_x_tok},
			want:  lambdaAbs1,
		},
		{
			name:  "newlines, single binder - 001",
			input: []api.Token{backslash, id_x_tok, thickArrow, newline, id_x_tok},
			want:  lambdaAbs1,
		},
		{
			name:  "newlines, single binder - 010",
			input: []api.Token{backslash, id_x_tok, newline, thickArrow, id_x_tok},
			want:  lambdaAbs1,
		},
		{
			name:  "newlines, single binder - 011",
			input: []api.Token{backslash, id_x_tok, newline, thickArrow, newline, id_x_tok},
			want:  lambdaAbs1,
		},
		{
			name:  "newlines, single binder - 100",
			input: []api.Token{backslash, newline, id_x_tok, thickArrow, id_x_tok},
			want:  lambdaAbs1,
		},
		{
			name:  "newlines, single binder - 101",
			input: []api.Token{backslash, newline, id_x_tok, thickArrow, newline, id_x_tok},
			want:  lambdaAbs1,
		},
		{
			name:  "newlines, single binder - 110",
			input: []api.Token{backslash, newline, id_x_tok, newline, thickArrow, id_x_tok},
			want:  lambdaAbs1,
		},
		{
			name:  "newlines, single binder - 111",
			input: []api.Token{backslash, newline, id_x_tok, newline, thickArrow, newline, id_x_tok},
			want:  lambdaAbs1,
		},

		// multiple binders
		{
			name: "no newlines, multiple binders - x00xx",
			// \x, x => x
			input: []api.Token{backslash, id_x_tok, comma, id_x_tok, thickArrow, id_x_tok},
			want:  lambdaAbs2,
		},
		{
			name:  "newlines, multiple binders - x01xx",
			input: []api.Token{backslash, id_x_tok, comma, newline, id_x_tok, thickArrow, id_x_tok},
			want:  lambdaAbs2,
		},
		{
			name:  "newlines, multiple binders - x10xx",
			input: []api.Token{backslash, id_x_tok, newline, comma, id_x_tok, thickArrow, id_x_tok},
			want:  lambdaAbs2,
		},
		{
			name:  "newlines, multiple binders - x11xx",
			input: []api.Token{backslash, id_x_tok, newline, comma, newline, id_x_tok, thickArrow, id_x_tok},
			want:  lambdaAbs2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseLambdaAbstraction, -1))
	}
}

// rule:
//
//	```
//	binder = lower ident | upper ident | "(", {"\n"}, enc pattern, {"\n"}, ")" ;
//	```
//
// rule:
//
//	```
//	lambda binder = binder | "_" ;
//	```
func TestMaybeParseBinder_TestMaybeParseLambdaBinder(t *testing.T) {
	// binder tests
	var binder_tests = []struct {
		name  string
		input []api.Token
		want  binder
	}{
		{
			name:  "lower ident",
			input: []api.Token{id_x_tok},
			want:  lowerBinder,
		},
		{
			name:  "upper ident",
			input: []api.Token{id_MyId_tok},
			want:  data.Inl[pattern](upperId),
		},
		{
			name:  "enclosed pattern newline left-right-flanked",
			input: []api.Token{lparen, newline, id_x_tok, newline, rparen},
			want:  data.Inr[ident](encPattern),
		},
		{
			name:  "enclosed pattern newline right-flanked",
			input: []api.Token{lparen, id_x_tok, newline, rparen},
			want:  data.Inr[ident](encPattern),
		},
		{
			name:  "enclosed pattern newline left-flanked",
			input: []api.Token{lparen, newline, id_x_tok, rparen},
			want:  data.Inr[ident](encPattern),
		},
	}
	for _, test := range binder_tests {
		t.Run(test.name, maybeOutputFUT_endCheck(test.input, data.Just(test.want), parseMaybeBinder, -1))
	}

	// lambda binder tests
	type lambdaBinderTest struct {
		name  string
		input []api.Token
		want  data.Maybe[lambdaBinder]
	}
	tests := []lambdaBinderTest{}
	for _, test := range binder_tests {
		tests = append(tests, lambdaBinderTest{
			test.name,
			test.input,
			data.Just(data.EInl[lambdaBinder](test.want)),
		})
	}
	// add underscore test
	tests = append(tests,
		lambdaBinderTest{
			"wildcard",
			[]api.Token{token.Underscore.Make()},
			data.Just(data.EInr[lambdaBinder](wildcardNode)),
		},
		lambdaBinderTest{
			"nothing",
			[]api.Token{},
			data.Nothing[lambdaBinder](token.EndOfTokens.Make()),
		},
	)
	for _, test := range tests {
		t.Run(test.name, maybeOutputFUT_endCheck(test.input, test.want, parseMaybeLambdaBinder, -1))
	}
}
