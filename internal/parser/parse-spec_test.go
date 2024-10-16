//go:build test
// +build test

package parser

import (
	"testing"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/common/data"
)

// rule:
//
//	```
//	constrainer = upper ident, pattern | "(", {"\n"}, enc constrainer {"\n"}, ")" ;
//	enc constrainer = upper ident, {"\n"}, pattern ;
//	```
func TestParseConstrainer(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  constrainer
	}{
		{
			"single",
			[]api.Token{id_MyId_tok, id_x_tok},
			constrainerNode,
		},
		{
			"enclosed - 000",
			[]api.Token{lparen, id_MyId_tok, id_x_tok, rparen},
			constrainerNode,
		},
		{
			"enclosed - 001",
			[]api.Token{lparen, id_MyId_tok, id_x_tok, newline, rparen},
			constrainerNode,
		},
		{
			"enclosed - 010",
			[]api.Token{lparen, id_MyId_tok, newline, id_x_tok, rparen},
			constrainerNode,
		},
		{
			"enclosed - 011",
			[]api.Token{lparen, id_MyId_tok, newline, id_x_tok, newline, rparen},
			constrainerNode,
		},
		{
			"enclosed - 100",
			[]api.Token{lparen, newline, id_MyId_tok, id_x_tok, rparen},
			constrainerNode,
		},
		{
			"enclosed - 101",
			[]api.Token{lparen, newline, id_MyId_tok, id_x_tok, newline, rparen},
			constrainerNode,
		},
		{
			"enclosed - 110",
			[]api.Token{lparen, newline, id_MyId_tok, newline, id_x_tok, rparen},
			constrainerNode,
		},
		{
			"enclosed - 111",
			[]api.Token{lparen, newline, id_MyId_tok, newline, id_x_tok, newline, rparen},
			constrainerNode,
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseConstrainer, -1))
	}
}

func TestParseRequiringClause(t *testing.T) {

}

func TestParseSpecBody(t *testing.T) {

}

// rule:
//
//	```
//	spec def = "spec", {"\n"}, spec head, [{"\n"}, spec dependency], {"\n"}, "where", {"\n"}, spec body, [{"\n"}, requiring clause] ;
//	```
func TestParseSpecDef(t *testing.T) {
	var (
		depNoReq = makeSpecDef(specHeadNode, data.Just(dependencyNode), specDefBodyNode, data.Nothing[requiringClause]())
		depReq   = makeSpecDef(specHeadNode, data.Just(dependencyNode), specDefBodyNode, data.Just(requiringNode))
		noDepReq = makeSpecDef(specHeadNode, data.Nothing[pattern](), specDefBodyNode, data.Just(requiringNode))
	)
	tests := []struct {
		name  string
		input []api.Token
		want  specDef
		end int
	}{
		// no dependency, no requiring clause

		{
			"no dependency, no requiring clause - 000",
			[]api.Token{spec, id_MyId_tok, id_x_tok, where, id_x_tok, colon, id_x_tok},
			specDefNode, -1,
		},
		{
			"no dependency, no requiring clause - 001",
			[]api.Token{spec, id_MyId_tok, id_x_tok, where, newline, id_x_tok, colon, id_x_tok},
			specDefNode, -1,
		},
		{
			"no dependency, no requiring clause - 010",
			[]api.Token{spec, id_MyId_tok, id_x_tok, newline, where, id_x_tok, colon, id_x_tok},
			specDefNode, -1,
		},
		{
			"no dependency, no requiring clause - 100",
			[]api.Token{spec, newline, id_MyId_tok, id_x_tok, where, id_x_tok, colon, id_x_tok},
			specDefNode, -1,
		},
		{
			"no dependency, no requiring clause - ends correctly",
			[]api.Token{spec, id_MyId_tok, id_x_tok, where, id_x_tok, colon, id_x_tok, newline},
			//                                                                       ^-- should end here
			specDefNode, -2,
		},

		// dependency, no requiring clause

		{
			"dependency, no requiring clause",
			// spec MyId x from x where x : x
			[]api.Token{spec, id_MyId_tok, id_x_tok, from, id_x_tok, where, id_x_tok, colon, id_x_tok},
			depNoReq, -1,
		},
		{
			"dependency, no requiring clause - newline",
			[]api.Token{spec, id_MyId_tok, id_x_tok, newline, from, id_x_tok, where, id_x_tok, colon, id_x_tok},
			depNoReq, -1,
		},
		{
			"dependency, no requiring clause - ends correctly",
			// spec MyId x from x where x : x
			[]api.Token{spec, id_MyId_tok, id_x_tok, from, id_x_tok, where, id_x_tok, colon, id_x_tok, newline},
			//                                                                     should end here --^
			depNoReq, -2,
		},

		// no dependency, requiring clause

		{
			"no dependency, requiring clause",
			// spec MyId x where x : x requiring x = x
			[]api.Token{spec, id_MyId_tok, id_x_tok, where, id_x_tok, colon, id_x_tok, requiring, id_x_tok, equal, id_x_tok},
			noDepReq, -1,
		},
		{
			"no dependency, requiring clause - newline",
			[]api.Token{spec, id_MyId_tok, id_x_tok, where, id_x_tok, colon, id_x_tok, newline, requiring, id_x_tok, equal, id_x_tok},
			noDepReq, -1,
		},
		{
			"no dependency, requiring clause - ends correctly",
			[]api.Token{spec, id_MyId_tok, id_x_tok, where, id_x_tok, colon, id_x_tok, requiring, id_x_tok, equal, id_x_tok, newline},
			//                                                                                           should end here --^ 
			noDepReq, -2,
		},

		// dependency, requiring clause

		{
			"dependency, requiring clause",
			// spec MyId x from x where x : x requiring x = x
			[]api.Token{spec, id_MyId_tok, id_x_tok, from, id_x_tok, where, id_x_tok, colon, id_x_tok, requiring, id_x_tok, equal, id_x_tok},
			depReq, -1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseSpecDef, test.end))
	}
}

func TestParseSpecInst(t *testing.T) {

}

func TestParseSpecMemberGroup(t *testing.T) {

}

// rule:
//
//	```
//	spec head = [constraint, {"\n"}, "=>", {"\n"}], constrainer ;
//	```
func TestParseSpecHead(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  specHead
	}{
		{
			"no constraint",
			[]api.Token{id_MyId_tok, id_x_tok},
			specHeadNode,
		},
		{
			"with constraint - 00",
			[]api.Token{id_MyId_tok, id_x_tok, thickArrow, id_MyId_tok, id_x_tok},
			specHeadConstrNode,
		},
		{
			"with constraint - 01",
			[]api.Token{id_MyId_tok, id_x_tok, thickArrow, newline, id_MyId_tok, id_x_tok},
			specHeadConstrNode,
		},
		{
			"with constraint - 10",
			[]api.Token{id_MyId_tok, id_x_tok, newline, thickArrow, id_MyId_tok, id_x_tok},
			specHeadConstrNode,
		},
		{
			"with constraint - 11",
			[]api.Token{id_MyId_tok, id_x_tok, newline, thickArrow, newline, id_MyId_tok, id_x_tok},
			specHeadConstrNode,
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseSpecHead, -1))
	}
}

func TestParseSpecDependency(t *testing.T) {

}

func TestParseSpecInstTarget(t *testing.T) {

}

func TestParseSpecInstWhereClause(t *testing.T) {

}

func TestParseUpperIdSequence(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  data.List[upperIdent]
		end   int
	}{
		{
			"empty",
			[]api.Token{},
			data.Nil[upperIdent](),
			0,
		},
		{
			"single, no comma",
			[]api.Token{id_MyId_tok},
			//         ^ end
			data.Nil[upperIdent](),
			0,
		},
		{
			"single, with comma",
			// MyId,
			[]api.Token{id_MyId_tok, comma},
			//                            ^ end
			data.Makes(MyId_as_upper),
			2,
		},
		{
			"multiple, no trailing comma",
			// MyId, MyId
			[]api.Token{id_MyId_tok, comma, id_MyId_tok},
			//                            ^ end
			data.Makes(MyId_as_upper),
			2,
		},
		{
			"multiple, with trailing comma",
			// MyId, MyId,
			[]api.Token{id_MyId_tok, comma, id_MyId_tok, comma},
			//                                                ^ end
			data.Makes(MyId_as_upper, MyId_as_upper),
			4,
		},
		{
			"multiple, trailing with constraint tail",
			// MyId, MyId, MyId x
			[]api.Token{id_MyId_tok, comma, id_MyId_tok, comma, id_MyId_tok, id_x_tok},
			//                                                ^ end
			data.Makes(MyId_as_upper, MyId_as_upper),
			4,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := initTestParser(test.input)
			got := parseUpperIdSequence(p)
			if !equals(got, test.want) {
				t.Errorf("parseUpperIdSequence() = expected \n%v\n, got \n%v\n", sprintTree(test.want), sprintTree(got))
			}

			if p.tokenCounter != test.end {
				t.Errorf("after parseUpperIdSequence(): expected (*ParserState).tokenCounter=%d, but got (*ParserState).tokenCounter=%d", test.end, p.tokenCounter)
			}
		})
	}
}
