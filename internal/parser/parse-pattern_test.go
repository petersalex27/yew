//go:build test
// +build test

package parser

import (
	"testing"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/util/fun"
	"github.com/petersalex27/yew/common/data"
)

// rule:
//
//	```
//	pattern = pattern term, {pattern term} ;
//	```
func TestParsePattern(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  pattern
	}{
		{
			"one",
			[]api.Token{id_x_tok},
			patternNode, // x
		},
		{
			"two",
			[]api.Token{id_x_tok, id_x_tok},
			patternAppNode, // x x
		},
		{
			"three",
			[]api.Token{id_x_tok, id_x_tok, id_x_tok},
			patternAppNode2, // x x x
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, ParsePattern, -1))
	}
}

// rule:
//
//	```
//	enc pattern term = "=" | pattern term ;
//	pattern term =
//		pattern atom
//		| "_"
//		| "(", {"\n"}, enc pattern inner, {"\n"}, ")"
//		| "{", {"\n"}, enc pattern inner, {"\n"}, "}" ;
//	enc pattern inner = enc pattern, {{"\n"}, ",", enc pattern}, [{"\n"}, ","] ;
//	```
func TestPatternTerm(t *testing.T) {
	tests := []struct {
		name     string
		input    []api.Token
		want     data.Maybe[pattern]
		enclosed bool
	}{
		{
			"pattern atom",
			[]api.Token{id_x_tok},
			data.Just[pattern](name_x), // x
			false,
		},
		{
			"underscore",
			[]api.Token{underscoreTok},
			data.Just[pattern](wildcardNode), // _
			false,
		},
		{
			"equal",
			[]api.Token{equal},
			data.Just[pattern](name_eq), // =
			true,
		},
		{
			"enclosed tuple",
			[]api.Token{lparen, id_x_tok, comma, id_x_tok, rparen},
			data.Just(encPattern2), // (x, x)
			false,
		},
		{
			"enclosed tuple - trailing",
			[]api.Token{lparen, id_x_tok, comma, id_x_tok, comma, rparen},
			data.Just(encPattern2), // (x, x,)
			false,
		},
		{
			"implicit pattern arg seq",
			[]api.Token{lbrace, id_x_tok, comma, id_x_tok, rbrace},
			data.Just(encPattern2Implicit), // {x, x}
			false,
		},
		{
			"implicit pattern arg seq - trailing",
			[]api.Token{lbrace, id_x_tok, comma, id_x_tok, comma, rbrace},
			data.Just(encPattern2Implicit), // {x, x,}
			false,
		},
		{
			"enclosed - 00",
			[]api.Token{lparen, id_x_tok, rparen},
			data.Just(encPattern), // ( x )
			false,
		},
		{
			"enclosed - 01",
			[]api.Token{lparen, id_x_tok, newline, rparen},
			data.Just(encPattern), // (x\n)
			false,
		},
		{
			"enclosed - 10",
			[]api.Token{lparen, newline, id_x_tok, rparen},
			data.Just(encPattern), // (\nx)
			false,
		},
		{
			"enclosed - 11",
			[]api.Token{lparen, newline, id_x_tok, newline, rparen},
			data.Just(encPattern), // (\nx\n)
			false,
		},
	}

	for _, test := range tests {
		fut := fun.BinBind1st_PairTarget(maybeParsePatternTerm, test.enclosed)
		t.Run(test.name, maybeOutputFUT_endCheck(test.input, test.want, fut, -1))
	}
}

func TestMaybeParseName(t *testing.T) {
	// infix: lower, upper, name, symbol
	// non-infix: lower, upper, name, symbol

	tests := []struct {
		name  string
		input []api.Token
		want  data.Maybe[name]
	}{
		{
			"lower",
			[]api.Token{id_x_tok}, // x
			data.Just(name_x),
		},
		{
			"upper",
			[]api.Token{id_MyId_tok}, // MyId
			data.Just(name_MyId),
		},
		{
			"symbol",
			[]api.Token{id_dollar_tok}, // $
			data.Just(name_dollar),
		},
		{
			"infix",
			[]api.Token{infix_dollar_tok}, // ($)
			data.Just(name_infix_dollar),
		},
	}

	for _, test := range tests {
		p := &ParserState{state: state{tokens: test.input}, ast: nil}
		actual := maybeParseName(p)
		if !equals(actual, test.want) {
			t.Errorf("expected \n%v\n, got \n%v\n", test.want, actual)
		}
	}
}

// rule:
//
//	```
//	pattern atom = literal | name | "[]" | hole ;
//	```
func TestParsePatternAtom(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  patternAtom
	}{
		{
			"literal",
			[]api.Token{integerValTok},
			data.Inl[patternName](literalNode),
		},
		{
			"name",
			[]api.Token{id_x_tok},
			patternAtomNode,
		},
		{
			"empty list",
			[]api.Token{nilListTok},
			data.Inr[literal](nilList),
		},
		{
			"hole",
			[]api.Token{hole_x_tok},
			data.Inr[literal](holePatName),
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parsePatternAtom, -1))
	}
}
