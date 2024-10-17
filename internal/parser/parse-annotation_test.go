//go:build test
// +build test

package parser

import (
	"strings"
	"testing"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/api/util"
	"github.com/petersalex27/yew/common/data"
)

func Test_assert_annotatable(*testing.T) {
	type annotatable interface{ annotate(data.Maybe[annotations]) }
	var (
		_ annotatable = &def{}
		_ annotatable = &specDef{}
		_ annotatable = &specInst{}
		_ annotatable = &typeDef{}
		_ annotatable = &typeAlias{}
		_ annotatable = &typing{}
		_ annotatable = &syntax{}
		_ annotatable = &typeConstructor{}
	)
	// yippee!
}

func TestMaybeParseAnnotation(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  data.Maybe[annotation]
	}{
		{
			"empty",
			[]api.Token{},
			data.Nothing[annotation](token.EndOfTokens.Make()),
		},
		{
			"one flat annotation",
			[]api.Token{annot},
			data.Just(annotation_flat),
		},
		{
			"one enclosed annotation",
			[]api.Token{annotOpenTok, id_test_tok, rbracket},
			data.Just(annotation_enclosed),
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseAnnotation, -1))
	}
}

func TestParseAnnotations(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  data.Maybe[annotations]
	}{
		{
			"empty",
			[]api.Token{},
			data.Nothing[annotations](token.EndOfTokens.Make()),
		},
		{
			"one flat annotation",
			[]api.Token{annot},
			data.Just(annotationBlock1),
		},
		{
			"one enclosed annotation",
			[]api.Token{annotOpenTok, id_test_tok, rbracket},
			data.Just(data.EConstruct[annotations](annotation_enclosed)),
		},
		{
			// --@test
			// [@test]
			"multiple annotations",
			[]api.Token{annot, newline, annotOpenTok, id_test_tok, rbracket},
			data.Just(annotationBlock2),
		},
	}

	for _, test := range tests {
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, parseAnnotations_, -1))
	}
}

func TestParseOptionalFlatAnnotation(t *testing.T) {

	//err := inl[flatAnnotation](terr)

	res := data.EOne[flatAnnotation](annot)
	tests := []struct {
		name     string
		ts       []api.Token
		expected data.Maybe[flatAnnotation]
	}{
		{"empty", []api.Token{}, data.Nothing[flatAnnotation](res)},
		{"non-annotation", []api.Token{__}, data.Nothing[flatAnnotation](res)},
		{"good", []api.Token{annot}, data.Just(res)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &ParserState{state: state{tokens: test.ts}}

			actual, actualIsSomething := parseOptionalFlatAnnotation(p).Break()
			if !test.expected.IsNothing() != actualIsSomething {
				t.Errorf("expected %v, got %v", test.expected, actual)
			}

			unit, just := test.expected.Break()
			if !just {
				return // nothing to compare
			} else if !equals(actual, unit) {
				t.Errorf("expected %v, got %v", test.expected, actual)
			}
		})
	}
}

func TestParseOptionalEnclosedAnnotation(t *testing.T) {
	tests := []struct {
		name     string
		ts       []api.Token
		expected data.Either[data.Ers, data.Maybe[enclosedAnnotation]]
	}{
		{"empty", []api.Token{}, data.Ok(data.Nothing[enclosedAnnotation](__))},
		{"non-annotation", []api.Token{id_test_tok}, data.Ok(data.Nothing[enclosedAnnotation](__))},
		{
			"no content",
			[]api.Token{
				annotOpenTok, id_test_tok, rbracket,
			},
			data.Ok(data.Just(annotSimple)),
		},
		{
			"some content",
			[]api.Token{annotOpenTok, id_test_tok, id_test_tok, rbracket},
			data.Ok(data.Just(annotSomeContent)),
		},
		{
			"has newlines 1",
			[]api.Token{
				annotOpenTok,
				newline, id_test_tok,
				newline,
				newline, id_test_tok,
				newline,
				rbracket,
			},
			data.Ok(data.Just(annotSomeContent)),
		},
		{
			"has newlines 2",
			[]api.Token{
				annotOpenTok, id_test_tok,
				newline,
				newline, id_test_tok, rbracket,
			},
			data.Ok(data.Just(annotSomeContent)),
		},
		{
			"has inner brackets",
			[]api.Token{
				annotOpenTok, id_test_tok, lbracket, rbracket, rbracket,
			},
			data.Ok(data.Just(annotWithInnerBrackets)),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := &ParserState{state: state{tokens: test.ts}}

			actual := parseOptionalEnclosedAnnotation(p)
			if test.expected.IsLeft() != actual.IsLeft() {
				t.Errorf("expected %v, got %v", test.expected, actual)
			}

			if test.expected.IsLeft() {
				return // nothing to compare
			} else if !equals(actual, test.expected) {
				b := &strings.Builder{}
				util.PrintTree(b, test.expected)
				exp := b.String()
				b.Reset()
				util.PrintTree(b, actual)
				act := b.String()
				t.Errorf("expected \n%s\n, got \n%s\n", exp, act)
			}
		})
	}
}
