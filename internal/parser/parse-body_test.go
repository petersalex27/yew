//go:build test
// +build test

package parser

import (
	"testing"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/common/data"
)

func TestParseBody(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  body
	}{
		{
			"typing",
			// x : x
			[]api.Token{id_x_tok, colon, id_x_tok},
			body_typing,
		},
		{
			"type def",
			// MyId : x where MyId : x
			[]api.Token{id_MyId_tok, colon, id_x_tok, where, id_MyId_tok, colon, id_x_tok},
			body_typeDef,
		},
		{
			"spec def",
			// spec MyId x where x : x
			[]api.Token{spec, id_MyId_tok, id_x_tok, where, id_x_tok, colon, id_x_tok},
			body_specDef,
		},
		{
			"spec inst",
			// inst MyId x where x = x
			[]api.Token{inst, id_MyId_tok, id_x_tok, where, id_x_tok, equal, id_x_tok},
			body_specInst,
		},
		{
			"type alias",
			// alias MyId = MyId
			[]api.Token{alias, id_MyId_tok, equal, id_MyId_tok},
			body_alias,
		},
		{
			"syntax",
			// syntax `my` x = x
			[]api.Token{syntaxTok, raw_my_tok, id_x_tok, equal, id_x_tok},
			body_syntax,
		},
		{
			"def",
			// x = x
			[]api.Token{id_x_tok, equal, id_x_tok},
			data.EMakes[body](bodyElement(defNode)),
		},
		{
			"annotated",
			// --@test\nx : x
			[]api.Token{annot, newline, id_x_tok, colon, id_x_tok},
			data.EMakes[body](bodyElement(annotTypingNode)),
		},
		{
			"multiple",
			// x : x\nx = x
			[]api.Token{id_x_tok, colon, id_x_tok, newline, id_x_tok, equal, id_x_tok},
			data.EMakes[body](bodyElement(typingNode), bodyElement(defNode)),
		},
	}

	for _, test := range tests {
		fut := func(p Parser) data.Either[data.Ers, body] {
			bd, _ := parseBody(p)
			if es, mBd, isMBd := bd.Break(); !isMBd {
				return data.PassErs[body](es)
			} else if b, ok := mBd.Break(); ok {
				return data.Ok(b)
			} else {
				return data.Fail[body]("could not parse body", p)
			}
		}

		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, fut, -1))
	}
}

func TestParseBodyElement(t *testing.T) {

}

func TestParseConstructorName(t *testing.T) {

}

func TestParseConstructor(t *testing.T) {

}

func TestParseDef(t *testing.T) {

}

func TestParseDefBody(t *testing.T) {

}

func TestParseDerivingBody(t *testing.T) {

}

func TestParseDerivingClause(t *testing.T) {

}

func TestParseMainElement(t *testing.T) {

}

func TestParseSyntax(t *testing.T) {

}

func TestParseSyntaxBindingSymbol(t *testing.T) {
	
}

func TestParseSyntaxRule(t *testing.T) {

}

func TestParseTypeAlias(t *testing.T) {

}

func TestParseTypeConstructor(t *testing.T) {

}

func TestParseTypeDef(t *testing.T) {

}

func TestParseTyping(t *testing.T) {

}

func TestParseVisibleBodyElement(t *testing.T) {

}

func TestParseWhereClause(t *testing.T) {

}

func TestParseWithClause(t *testing.T) {

}

func TestParseWithClauseArms(t *testing.T) {

}

func TestParseWithClauseArm(t *testing.T) {

}
