//go:build test
// +build test

package parser

import (
	"testing"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/util/fun"
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


func TestParseConstructorNameErrors(t *testing.T) {
	tests := []struct {
		name string
		input api.Token
		want string
	}{
		{
			"error - lower ident",
			id_x_tok,
			IllegalLowercaseConstructorName,
		},
		{
			"error - lower infix ident",
			infix_x_tok,
			IllegalLowercaseConstructorName,
		},
		{
			"error - method name",
			method_run_tok,
			IllegalMethodTypeConstructor,
		},
		{
			"error - non-name",
			alias,
			ExpectedTypeConstructorName,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := constructorName_Error(test.input)
			if res != test.want {
				t.Errorf("failed: got %q, want %q", res, test.want)
			}
		})
	}
}

// rule:
//
//	```
//	constructor name = infix upper ident | upper ident | symbol | infix symbol ;
//	```
func TestParseConstructorName(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  data.Maybe[name]
	}{
		{
			"infix upper ident",
			[]api.Token{infix_MyId_tok},
			data.Just(name_infix_MyId),
		},
		{
			"upper ident",
			[]api.Token{id_MyId_tok},
			data.Just(name_MyId),
		},
		{
			"symbol",
			[]api.Token{id_dollar_tok},
			data.Just(name_dollar),
		},
		{
			"infix symbol",
			[]api.Token{infix_dollar_tok},
			data.Just(name_infix_dollar),
		},
	}

	for _, test := range tests {
		t.Run(test.name, maybeOutputFUT_endCheck(test.input, test.want, maybeParseConstructorName, -1))
	}
}

// rule: 
//
//	```
//	type constructor = constructor name, {{"\n"}, ",", {"\n"}, constructor name}, {"\n"}, ":", {"\n"}, type ;
//	```
func TestParseConstructor(t *testing.T) {
	tests := []struct {
		name string
		input []api.Token
		want data.NonEmpty[typeConstructor]
	}{
		{
			"single - 00",
			[]api.Token{id_MyId_tok, colon, id_x_tok},
			singleConsNode,
		},
		{
			"single - 01",
			[]api.Token{id_MyId_tok, colon, newline, id_x_tok},
			singleConsNode,
		},
		{
			"single - 10",
			[]api.Token{id_MyId_tok, newline, colon, id_x_tok},
			singleConsNode,
		},
		{
			"single - 11",
			[]api.Token{id_MyId_tok, newline, colon, newline, id_x_tok},
			singleConsNode,
		},
		{
			"multiple - 00",
			[]api.Token{id_MyId_tok, comma, id_MyId_tok, colon, id_x_tok},
			multiConsNode,
		},
		{
			"multiple - 01",
			[]api.Token{id_MyId_tok, comma, newline, id_MyId_tok, colon, id_x_tok},
			multiConsNode,
		},
		{
			"multiple - 10",
			[]api.Token{id_MyId_tok, newline, comma, id_MyId_tok, colon, id_x_tok},
			multiConsNode,
		},
		{
			"multiple - 11",
			[]api.Token{id_MyId_tok, newline, comma, newline, id_MyId_tok, colon, id_x_tok},
			multiConsNode,
		},
		{
			"single - trailing comma",
			[]api.Token{id_MyId_tok, comma, colon, id_x_tok},
			data.Construct(makeCons(name_MyId, typ_x)),
		},
		{
			"multiple - trailing comma",
			[]api.Token{id_MyId_tok, comma, id_MyId_tok, comma, colon, id_x_tok},
			multiConsNode,
		},
	}

	for _, test := range tests {
		fut := fun.Bind1stOf2(parseTypeConstructor, data.Nothing[annotations]())
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, fut , -1))
	}
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

func TestParseTypeDefBody(t *testing.T) {
	tests := []struct {
		name string
		input []api.Token
		want typeDefBody
	}{
		{
			"impossible",
			[]api.Token{impossibleTok},
			data.Inr[data.NonEmpty[typeConstructor]](impossibleNode),
		},
		{
			"non-grouped",
			[]api.Token{id_MyId_tok, colon, id_x_tok},
			data.Inl[impossible](singleConsNode),
		},
		{
			"enclosed - 00",
			[]api.Token{lparen, id_MyId_tok, colon, id_x_tok, rparen},
			data.Inl[impossible](singleConsNode),
		},
		{
			"enclosed - 01",
			[]api.Token{lparen, id_MyId_tok, colon, id_x_tok, newline, rparen},
			data.Inl[impossible](singleConsNode),
		},
		{
			"enclosed - 10",
			[]api.Token{lparen, newline, id_MyId_tok, colon, id_x_tok, rparen},
			data.Inl[impossible](singleConsNode),
		},
		{
			"enclosed - 11",
			[]api.Token{lparen, newline, id_MyId_tok, colon, id_x_tok, newline, rparen},
			data.Inl[impossible](singleConsNode),
		},
		{
			"enclosed - multiple",
			[]api.Token{lparen, id_MyId_tok, colon, id_x_tok, newline, id_MyId_tok, colon, id_x_tok, rparen},
			data.Inl[impossible](multiConsNode),
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseTypeDefBody, -1))
	}
}

// not much to test here, just make sure the name parse, colon parse, and type parse are 
// correctly sequenced to allow for newlines in appropriate places
func TestParseTyping(t *testing.T) {
	tests := []struct {
		name string
		input []api.Token
		want typing
	}{
		{
			"00",
			[]api.Token{id_x_tok, colon, id_x_tok},
			typingNode,
		},
		{
			"01",
			[]api.Token{id_x_tok, colon, newline, id_x_tok},
			typingNode,
		},
		{
			"10",
			[]api.Token{id_x_tok, newline, colon, id_x_tok},
			typingNode,
		},
		{
			"11",
			[]api.Token{id_x_tok, newline, colon, newline, id_x_tok},
			typingNode,
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseTypeSig, -1))
	}
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
