//go:build test
// +build test

package parser

import (
	"testing"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/common/data"
)

func TestParseYewSource(t *testing.T) {
	tests := []struct {
		name  string
		input []api.Token
		want  yewSource
	}{
		{
			"source - 0000 (empty footer only)",
			[]api.Token{},
			makeYewSource(
				//data.Nothing[meta](),
				data.Nothing[header](),
				data.Nothing[body](),
				data.Nothing[annotations](),
			),
		},
		{
			"source - 0001 (footer only)",
			[]api.Token{annot},
			makeYewSource(
				//data.Nothing[meta](),
				data.Nothing[header](),
				data.Nothing[body](),
				data.Just(data.EConstruct[annotations](annotation_flat)),
			),
		},
		{
			"source - 0010 (body only)",
			[]api.Token{id_x_tok, colon, id_x_tok},
			makeYewSource(
				//data.Nothing[meta](),
				data.Nothing[header](),
				data.Just(body_typing),
				data.Nothing[annotations](),
			),
		},
		{
			"source - 0011 (body and footer)",
			[]api.Token{id_x_tok, colon, id_x_tok, newline, annot},
			makeYewSource(
				//data.Nothing[meta](),
				data.Nothing[header](),
				data.Just(body_typing),
				data.Just(data.EConstruct[annotations](annotation_flat)),
			),
		},
		{
			"source - 0100 (header only)",
			[]api.Token{moduleTok, id_x_tok},
			makeYewSource(
				data.Just(data.EMakePair[header](data.Just(module_x), data.Nil[importStatement]())),
				data.Nothing[body](),
				data.Nothing[annotations](),
			),
		},
		{
			"source - 0101 (header and footer)",
			[]api.Token{moduleTok, id_x_tok, newline, annot},
			makeYewSource(
				//data.Nothing[meta](),
				data.Just(data.EMakePair[header](data.Just(module_x), data.Nil[importStatement]())),
				data.Nothing[body](),
				data.Just(data.EConstruct[annotations](annotation_flat)),
			),
		},
		{
			"source - 0110 (header and body)",
			[]api.Token{moduleTok, id_x_tok, newline, id_x_tok, colon, id_x_tok},
			makeYewSource(
				//data.Nothing[meta](),
				data.Just(data.EMakePair[header](data.Just(module_x), data.Nil[importStatement]())),
				data.Just(body_typing),
				footer{data.Nothing[annotations]()},
			),
		},
		{
			"source - 0111 (header, body, and footer)",
			[]api.Token{moduleTok, id_x_tok, newline, id_x_tok, colon, id_x_tok, newline, annot},
			makeYewSource(
				data.Just(data.EMakePair[header](data.Just(module_x), data.Nil[importStatement]())),
				data.Just(body_typing),
				data.Just(data.EConstruct[annotations](annotation_flat)),
			),
		},

		// NOTE: this shows there are really just three *syntactic* sections of a yew source file; however,
		// once name analysis is complete, the annotations will be divided b/w the module and the header.

		// "source - 1000 (meta only)" is not possible, as the footer will take the annotations
		// 		- would be equivalent to "source - 0001 (footer only)"

		// "source - 1001 (meta and footer)" is not possible, as the footer will take the annotations
		// 		- would be equivalent-ish to "source - 0001 (footer only)"--with > 1 annotation though

		// "source - 1010 (meta and body)" is not possible, as the body will take the annotations
		// 		- would be equivalent to "source - 0010 (body only)"--with the body element annotated

		// "source - 1011 (meta, body, and footer)" is not possible, as the body will take the annotations
		//		- would be equivalent to "source - 0011 (body and footer)"--with the body element annotated

		{
			"source - 1100 (meta and header)",
			[]api.Token{annot, newline, moduleTok, id_x_tok},
			makeYewSource(
				data.Just(data.EMakePair[header](data.Just(module_annot_x), data.Nil[importStatement]())),
				data.Nothing[body](),
				data.Nothing[annotations](),
			),
		},
		{
			"source - 1101 (meta, header, and footer)",
			[]api.Token{annot, newline, moduleTok, id_x_tok, newline, annot},
			makeYewSource(
				data.Just(data.EMakePair[header](data.Just(module_annot_x), data.Nil[importStatement]())),
				data.Nothing[body](),
				data.Just(data.EConstruct[annotations](annotation_flat)),
			),
		},
		{
			"source - 1110 (meta, header, and body)",
			[]api.Token{annot, newline, moduleTok, id_x_tok, newline, id_x_tok, colon, id_x_tok},
			makeYewSource(
				data.Just(data.EMakePair[header](data.Just(module_annot_x), data.Nil[importStatement]())),
				data.Just(body_typing),
				data.Nothing[annotations](),
			),
		},
		{
			"source - 1111 (meta, header, body, and footer--all source section present)",
			[]api.Token{annot, newline, moduleTok, id_x_tok, newline, id_x_tok, colon, id_x_tok, newline, annot},
			makeYewSource(
				data.Just(data.EMakePair[header](data.Just(module_annot_x), data.Nil[importStatement]())),
				data.Just(body_typing),
				data.Just(data.EConstruct[annotations](annotation_flat)),
			),
		},
	}

	for _, test := range tests {
		fut := func(p parser) data.Either[data.Ers, yewSource] {
			p = parseYewSource(p)
			if ps, ok := p.(*ParserState); ok {
				return data.Ok(ps.ast)
			} else {
				return data.Fail[yewSource]("(test failure) could not parse yew source", p)
			}
		}
		t.Run(test.name, resultOutputFUT_endCheck(test.input, test.want, fut, -1))
	}
}