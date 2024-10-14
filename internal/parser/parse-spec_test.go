//go:build test
// +build test

package parser

import (
	"testing"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/common/data"
)

func TestParseConstrainer(t *testing.T) {
	
}

func TestParseRequiringClause(t *testing.T) {
	
}

func TestParseSpecBody(t *testing.T) {
	
}

func TestParseSpecDef(t *testing.T) {
	
}

func TestParseSpecInst(t *testing.T) {

}

func TestParseSpecMemberGroup(t *testing.T) {
	
}

func TestParseSpecHead(t *testing.T) {
	
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
		want data.List[upperIdent]
		end int
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