package parser

import (
	"testing"

	"github.com/petersalex27/yew/token"
)

func TestFindMeaningfulIndent(t *testing.T) {
	// no meaningful indent
	// meaningful indent is first indent
	// meaningful indent is not first indent

	tests := []struct {
		tokens  []token.Token
		indent  token.Token
		success bool
	}{
		{
			tokens:  []token.Token{},
			indent:  endOfTokensToken(),
			success: false,
		},
		{
			tokens:  []token.Token{token.Id.MakeValued("x")},
			indent:  token.Id.MakeValued("x"),
			success: false,
		},
		{
			tokens:  []token.Token{token.Indent.MakeValued(" ")},
			indent:  token.Indent.MakeValued(" "),
			success: true,
		},
		{
			tokens:  []token.Token{token.Indent.MakeValued(" "), token.Id.MakeValued("x")},
			indent:  token.Indent.MakeValued(" "),
			success: true,
		},
		{
			tokens:  []token.Token{token.Indent.MakeValued(" "), token.Newline.Make(), token.Indent.MakeValued("  ")},
			indent:  token.Indent.MakeValued("  "),
			success: true,
		},
		{
			// indent should follow newline
			tokens:  []token.Token{token.Indent.MakeValued(" "), token.Newline.Make(), token.Id.MakeValued("x")},
			indent:  token.Indent.MakeValued(" "),
			success: false,
		},
		{
			// indent should follow newline
			//
			//	1 | ⎵--my comment
			//  2 | x
			tokens: []token.Token{
				token.Indent.MakeValued(" "),
				token.Comment.MakeValued("my comment"),
				token.Newline.Make(),
				token.Id.MakeValued("x"),
			},
			indent:  token.Indent.MakeValued(" "),
			success: false,
		},
		{
			//	1 | ⎵-*my comment*- x
			tokens: []token.Token{
				token.Indent.MakeValued(" "),
				token.Comment.MakeValued("my comment"),
				token.Id.MakeValued("x"),
			},
			indent:  token.Indent.MakeValued(" "),
			success: true,
		},
		{
			//  1 | ⎵-*my comment*-
			//  2 | ⎵⎵x
			tokens: []token.Token{
				token.Indent.MakeValued(" "),
				token.Comment.MakeValued("my comment"),
				token.Newline.Make(),
				token.Indent.MakeValued("  "),
				token.Id.MakeValued("x"),
			},
			indent:  token.Indent.MakeValued("  "),
			success: true,
		},
		{
			//  1 | ⎵-*my comment*-
			//  2 | ⎵
			//	3 | ⎵
			//  4 | ⎵⎵x
			tokens: []token.Token{
				token.Indent.MakeValued(" "),
				token.Comment.MakeValued("my comment"),
				token.Newline.Make(),
				token.Indent.MakeValued(" "),
				token.Newline.Make(),
				token.Indent.MakeValued(" "),
				token.Newline.Make(),
				token.Indent.MakeValued("  "),
				token.Id.MakeValued("x"),
			},
			indent:  token.Indent.MakeValued("  "),
			success: true,
		},
		{
			tokens:  []token.Token{token.Indent.MakeValued(" "), token.Newline.Make(), token.Indent.MakeValued("  ")},
			indent:  token.Indent.MakeValued("  "),
			success: true,
		},
	}

	for _, test := range tests {
		p := &tokenInfo{tokens: test.tokens}
		actualIndent, actualSuccess := p.findMeaningfulIndent()
		if actualSuccess != test.success {
			t.Fatalf("unexpected success (%v): got %t", test, actualSuccess)
		}
		if !test.indent.Equals(actualIndent) {
			t.Fatalf("unexpected token (%v): got %v", test, actualIndent)
		}
	}
}
