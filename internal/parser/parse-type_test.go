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
//	type =
//		["forall", {"\n"}, forall binders, {"\n"}, "in", {"\n"}], type tail
//		| "(", {"\n"}, enc type, {"\n"}, ")" ;
//	```
func TestParseType(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  typ
	}{
		{
			"single",
			[]api.Token{id_x_tok},
			typ_x,
		},
		{
			"enclosed - 00",
			[]api.Token{lparen, id_x_tok, rparen},
			typ_x,
		},
		{
			"enclosed - 01",
			[]api.Token{lparen, id_x_tok, newline, rparen},
			typ_x,
		},
		{
			"enclosed - 10",
			[]api.Token{lparen, newline, id_x_tok, rparen},
			typ_x,
		},
		{
			"enclosed - 11",
			[]api.Token{lparen, newline, id_x_tok, newline, rparen},
			typ_x,
		},
		{
			"forall type - 000",
			[]api.Token{forall, id_x_tok, in, id_x_tok},
			forallTypeNode,
		},
		{
			"forall type - 001",
			[]api.Token{forall, id_x_tok, in, newline, id_x_tok},
			forallTypeNode,
		},
		{
			"forall type - 010",
			[]api.Token{forall, id_x_tok, newline, in, id_x_tok},
			forallTypeNode,
		},
		{
			"forall type - 011",
			[]api.Token{forall, id_x_tok, newline, in, newline, id_x_tok},
			forallTypeNode,
		},
		{
			"forall type - 100",
			[]api.Token{forall, newline, id_x_tok, in, id_x_tok},
			forallTypeNode,
		},
		{
			"forall type - 101",
			[]api.Token{forall, newline, id_x_tok, in, newline, id_x_tok},
			forallTypeNode,
		},
		{
			"forall type - 110",
			[]api.Token{forall, newline, id_x_tok, newline, in, id_x_tok},
			forallTypeNode,
		},
		{
			"forall type - 111",
			[]api.Token{forall, newline, id_x_tok, newline, in, newline, id_x_tok},
			forallTypeNode,
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, ParseType, -1))
	}
}

// rule:
//
//	```
//	type tail = type term, {type term}, [{"\n"}, ("->" | "=>"), {"\n"}, type tail] ;
//	```
func TestParseTypeTail(t *testing.T) {
	tests := []struct {
		name     string
		input    []api.Token
		want     typ
		enclosed bool
	}{
		// application cases
		{
			"single",
			[]api.Token{id_x_tok},
			typ_x,
			false,
		},
		{
			"single - 0",
			[]api.Token{id_x_tok},
			typ_x,
			true,
		},
		{
			"type application",
			[]api.Token{id_x_tok, id_x_tok},
			data.MakeApp[appType](typ_x, typ_x),
			false,
		},
		{
			"type application - 0",
			[]api.Token{id_x_tok, id_x_tok},
			data.MakeApp[appType](typ_x, typ_x),
			true,
		},
		{
			"type application - 1",
			[]api.Token{id_x_tok, newline, id_x_tok},
			data.MakeApp[appType](typ_x, typ_x),
			true,
		},
		// function cases
		{
			"function - 00",
			// x -> x
			[]api.Token{id_x_tok, arrow, id_x_tok},
			makeFunc(typ_x, typ_x),
			false, // NOTE: this is not enclosed
		},
		{
			"function - 01",
			// x -> \n x
			[]api.Token{id_x_tok, arrow, newline, id_x_tok},
			makeFunc(typ_x, typ_x),
			false, // NOTE: this is not enclosed, but shouldn't disallow newlines next to **the arrow**
		},
		{
			"function - 10",
			// x \n -> x
			[]api.Token{id_x_tok, newline, arrow, id_x_tok},
			makeFunc(typ_x, typ_x),
			false,
		},
		{
			"function - 11",
			// x \n -> \n x
			[]api.Token{id_x_tok, newline, arrow, newline, id_x_tok},
			makeFunc(typ_x, typ_x),
			false,
		},
		// function w/ lhs application
		{
			"function w/ lhs application",
			[]api.Token{id_x_tok, id_x_tok, arrow, id_x_tok},
			// x x -> x
			makeFunc(data.MakeApp[appType](typ_x, typ_x), typ_x),
			false,
		},
		{
			"constrained type",
			[]api.Token{id_MyId_tok, id_x_tok, thickArrow, id_x_tok},
			// MyId x => x
			makeUnverifiedConstrainedType(data.MakeApp[appType](typ(name_MyId), typ_x), typ_x),
			false,
		},
	}

	for _, test := range tests {
		fut := fun.Bind1stOf2(parseTypeTail, test.enclosed)
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, fut, -1))
	}
}

// rule:
//
//	```
//	forall binding = "forall", {"\n"}, forall binders, {"\n"}, "in", {"\n"}
//	```
func TestParseForallBinding(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  forallBinders
	}{
		{
			"single - lower",
			[]api.Token{forall, id_x_tok},
			data.EConstruct[forallBinders](lowerId),
		},
		{
			"single - upper",
			[]api.Token{forall, id_MyId_tok},
			data.EConstruct[forallBinders](upperId),
		},
		{
			"multiple - lower * 2",
			[]api.Token{forall, id_x_tok, id_x_tok},
			data.EConstruct[forallBinders](lowerId, lowerId),
		},
		{
			"multiple - lower * upper",
			[]api.Token{forall, id_x_tok, id_MyId_tok},
			data.EConstruct[forallBinders](lowerId, upperId),
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseForallBinding, -1))
	}
}

// rule:
//
//	```
//	type exceptions = "_" | "()" | "="
//	```
func TestParseTypeTermException(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  typ
	}{
		{
			"wildcard",
			[]api.Token{underscoreTok},
			wildcardNode, // x
		},
		{
			"unit",
			[]api.Token{unitTypeTok},
			data.EOne[unitType](unitTypeTok), // ()
		},
		{
			"equal",
			[]api.Token{equal},
			name_eq, // =
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseTypeTermException, -1))
	}
}

// rule:
//
//	```
//	type term =
//		expr root
//		| "_" | "()" | "="
//		| "(", {"\n"}, enc type inner, [{"\n"}, enc typing end], {"\n"}, ")"
//		| "{", {"\n"}, enc type inner, [{"\n"}, enc typing end, [{"\n"}, default expr]], {"\n"}, "}" ;
//	```
func TestMaybeParseTypeTerm(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  data.Maybe[typ]
	}{
		{
			"expr atom",
			[]api.Token{id_x_tok},
			data.Just(typ_x), // x
		},
		{
			"underscore",
			[]api.Token{underscoreTok},
			data.Just[typ](wildcardNode), // _
		},
		{
			"unit",
			[]api.Token{unitTypeTok},
			data.Just[typ](data.EOne[unitType](unitTypeTok)), // ()
		},
		{
			"equal",
			[]api.Token{equal},
			data.Just[typ](name_eq), // =
		},
		{
			"enclosed",
			[]api.Token{lparen, id_x_tok, rparen},
			data.Just[typ](enclosedTypeNode), // ( x )
		},
		{
			"enclosed with typing",
			[]api.Token{lparen, id_x_tok, colon, id_x_tok, rparen},
			data.Just[typ](enclosedTypingNode), // (x : x)
		},
		{
			"enclosed with once modality",
			[]api.Token{lparen, onceTok, id_x_tok, colon, id_x_tok, rparen},
			data.Just[typ](enclosedOnceTypingNode), // (once x : x)
		},
		{
			"enclosed with erase modality",
			[]api.Token{lparen, eraseTok, id_x_tok, colon, id_x_tok, rparen},
			data.Just[typ](enclosedEraseTypingNode), // (erase x : x)
		},
		{
			"enclosed with type sequence",
			[]api.Token{lparen, id_x_tok, comma, id_x_tok, rparen},
			data.Just[typ](enclosedTypeSeqNode), // (x, x)
		},
		{
			"enclosed with typing sequence",
			[]api.Token{lparen, id_x_tok, comma, id_x_tok, colon, id_x_tok, rparen},
			data.Just[typ](enclosedTypingSeqNode), // (x, x : x)
		},
		{
			"enclosed implicit",
			[]api.Token{lbrace, id_x_tok, rbrace},
			data.Just[typ](implicitEnclosedTypeNode), // { x }
		},
		{
			"enclosed implicit sequence",
			[]api.Token{lbrace, id_x_tok, comma, id_x_tok, rbrace},
			data.Just[typ](implicitEnclosedTypeSeqNode), // { x, x }
		},
		{
			"enclosed implicit with typing",
			[]api.Token{lbrace, id_x_tok, colon, id_x_tok, rbrace},
			data.Just[typ](implicitEnclosedTypingNode), // { x : x }
		},
		{
			"enclosed implicit with default",
			[]api.Token{lbrace, id_x_tok, colon, id_x_tok, colonEqual, id_x_tok, rbrace},
			data.Just[typ](implicitEnclosedTypingNode_def), // { x : x := x }
		},
		{
			"enclosed implicit with once modality",
			[]api.Token{lbrace, onceTok, id_x_tok, colon, id_x_tok, rbrace},
			data.Just[typ](implicitEnclosedOnceTypingNode), // { once x : x }
		},
		{
			"enclosed implicit with erase modality",
			[]api.Token{lbrace, eraseTok, id_x_tok, colon, id_x_tok, rbrace},
			data.Just[typ](implicitEnclosedEraseTypingNode), // { erase x : x }
		},
		{
			"enclosed implicit with default sequence",
			[]api.Token{lbrace, id_x_tok, comma, id_x_tok, colon, id_x_tok, colonEqual, id_x_tok, rbrace},
			data.Just[typ](implicitEnclosedTypingSeqNode_def), // { x, x : x := x }
		},
		{
			"enclosed implicit with typing sequence",
			[]api.Token{lbrace, id_x_tok, comma, id_x_tok, colon, id_x_tok, rbrace},
			data.Just[typ](implicitEnclosedTypingSeqNode), // { x, x : x }
		},
	}

	for _, test := range tests {
		t.Run(test.name, maybeOutputFUT_endCheck(test.input, test.want, maybeParseTypeTerm, -1))
	}
}
