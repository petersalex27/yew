// =================================================================================================
// Alex Peters - January 22, 2024
// =================================================================================================

package lexer

import (
	"testing"

	"github.com/petersalex27/yew/errors"
	"github.com/petersalex27/yew/source"
	"github.com/petersalex27/yew/token"
)

func TestAnalyzeNumber(t *testing.T) {
	tests := []struct {
		source string
		pos    []int
		expect token.Token
	}{
		{
			source: "1",
			expect: token.Token{Value: "1", Type: token.IntValue, End: 1},
		},
		{
			source: "0x1",
			expect: token.Token{Value: "0x1", Type: token.IntValue, End: 3},
		},
		{
			source: "0xa",
			expect: token.Token{Value: "0xa", Type: token.IntValue, End: 3},
		},
		{
			source: "0Xa",
			expect: token.Token{Value: "0Xa", Type: token.IntValue, End: 3},
		},
		{
			source: "0o1",
			expect: token.Token{Value: "0o1", Type: token.IntValue, End: 3},
		},
		{
			source: "0O1",
			expect: token.Token{Value: "0O1", Type: token.IntValue, End: 3},
		},
		{
			source: "0b1",
			expect: token.Token{Value: "0b1", Type: token.IntValue, End: 3},
		},
		{
			source: "0B1",
			expect: token.Token{Value: "0B1", Type: token.IntValue, End: 3},
		},
		{
			source: "1.0",
			expect: token.Token{Value: "1.0", Type: token.FloatValue, End: 3},
		},
		{
			source: "1e1",
			expect: token.Token{Value: "1e1", Type: token.FloatValue, End: 3},
		},
		{
			source: "1E1",
			expect: token.Token{Value: "1E1", Type: token.FloatValue, End: 3},
		},
		{
			source: "1e+1",
			expect: token.Token{Value: "1e+1", Type: token.FloatValue, End: 4},
		},
		{
			source: "1e-1",
			expect: token.Token{Value: "1e-1", Type: token.FloatValue, End: 4},
		},
		{
			source: "1.0e1",
			expect: token.Token{Value: "1.0e1", Type: token.FloatValue, End: 5},
		},
		{
			source: "0_1",
			expect: token.Token{Value: "1", Type: token.IntValue, End: 3},
		},
		{
			source: "00_1",
			expect: token.Token{Value: "1", Type: token.IntValue, End: 4},
		},
		{
			source: "11__1",
			expect: token.Token{Value: "111", Type: token.IntValue, End: 5},
		},
		{
			source: "0x1_1",
			expect: token.Token{Value: "0x11", Type: token.IntValue, End: 5},
		},
		{
			source: "0o1_1",
			expect: token.Token{Value: "0o11", Type: token.IntValue, End: 5},
		},
		{
			source: "0b1_1",
			expect: token.Token{Value: "0b11", Type: token.IntValue, End: 5},
		},
	}

	for _, test := range tests {
		lex := Init(source.StdinSpec)
		lex.Source = []byte(test.source)
		lex.PositionRanges = []int{len(test.source)}
		lex.Line = 1

		ok, eof := lex.analyzeNumber()
		if !ok {
			t.Fatalf("could not analyze number")
		}
		if eof {
			t.Fatalf("unexpected end of file")
		}

		if len(lex.Tokens) != 1 {
			t.Fatalf("unexpected token stream: got %v", lex.Tokens)
		}

		actual := lex.Tokens[0]
		if actual != test.expect {
			t.Fatalf("unexpected token (%v): got %v", test.expect, actual)
		}
	}
}

func TestAnalyzeChar(t *testing.T) {
	tests := []struct {
		source string
		expect token.Token
	}{
		{
			source: `'a'`,
			expect: token.Token{Value: "a", Type: token.CharValue, End: 3},
		},
		{
			source: `' '`,
			expect: token.Token{Value: " ", Type: token.CharValue, End: 3},
		},
		{
			source: `'@'`,
			expect: token.Token{Value: "@", Type: token.CharValue, End: 3},
		},
		{
			source: `'\n'`,
			expect: token.Token{Value: "\n", Type: token.CharValue, End: 4},
		},
		{
			source: `'\t'`,
			expect: token.Token{Value: "\t", Type: token.CharValue, End: 4},
		},
		{
			source: `'\a'`,
			expect: token.Token{Value: "\a", Type: token.CharValue, End: 4},
		},
		{
			source: `'\b'`,
			expect: token.Token{Value: "\b", Type: token.CharValue, End: 4},
		},
		{
			source: `'\v'`,
			expect: token.Token{Value: "\v", Type: token.CharValue, End: 4},
		},
		{
			source: `'\f'`,
			expect: token.Token{Value: "\f", Type: token.CharValue, End: 4},
		},
		{
			source: `'\r'`,
			expect: token.Token{Value: "\r", Type: token.CharValue, End: 4},
		},
		{
			source: `'\''`,
			expect: token.Token{Value: "'", Type: token.CharValue, End: 4},
		},
		{
			source: `'\\'`,
			expect: token.Token{Value: "\\", Type: token.CharValue, End: 4},
		},
	}

	for _, test := range tests {
		lex := Init(source.StdinSpec)
		lex.Source = []byte(test.source)
		lex.PositionRanges = []int{len(test.source)}
		lex.Line = 1

		ok, eof := lex.analyzeChar()
		if !ok {
			t.Fatalf("could not analyze char")
		}
		if eof {
			t.Fatalf("unexpected end of file")
		}

		if len(lex.Tokens) != 1 {
			t.Fatalf("unexpected token stream: got %v", lex.Tokens)
		}

		actual := lex.Tokens[0]
		if actual != test.expect {
			t.Fatalf("unexpected token (%v): got %v", test.expect, actual)
		}
	}
}

func TestAnalyzeString(t *testing.T) {
	tests := []struct {
		source string
		expect token.Token
	}{
		{
			source: `""`,
			expect: token.Token{Value: "", Type: token.StringValue, End: 2},
		},
		{
			source: `" "`,
			expect: token.Token{Value: " ", Type: token.StringValue, End: 3},
		},
		{
			source: `"--"`,
			expect: token.Token{Value: "--", Type: token.StringValue, End: 4},
		},
		{
			source: `"this is a string"`,
			expect: token.Token{Value: "this is a string", Type: token.StringValue, End: 18},
		},
		{
			source: `"\n\t\a\b\v\f\r\"\\"`,
			expect: token.Token{Value: "\n\t\a\b\v\f\r\"\\", Type: token.StringValue, End: 20},
		},
	}

	for _, test := range tests {
		lex := Init(source.StdinSpec)
		lex.Source = []byte(test.source)
		lex.PositionRanges = []int{len(test.source)}
		lex.Line = 1

		ok, eof := lex.analyzeString()
		if !ok {
			t.Fatalf("could not analyze string")
		}
		if eof {
			t.Fatalf("unexpected end of file")
		}

		if len(lex.Tokens) != 1 {
			t.Fatalf("unexpected token stream: got %v", lex.Tokens)
		}

		actual := lex.Tokens[0]
		if actual != test.expect {
			t.Fatalf("unexpected token (%v): got %v", test.expect, actual)
		}
	}
}

func TestAnalyzeAnnotation(t *testing.T) {
	tests := []struct {
		source string
		expect token.Token
	}{
		{
			source: `-- @annotation`,
			expect: token.Token{Value: "annotation", Type: token.At, Start: 4, End: 14},
		},
		{
			source: `-- @Annotation`,
			expect: token.Token{Value: "Annotation", Type: token.At, Start: 4, End: 14},
		},
		{
			source: `-- @annotation _ stuff`,
			expect: token.Token{Value: "annotation", Type: token.At, Start: 4, End: 14},
		},
	}

	for _, test := range tests {
		lex := Init(source.StdinSpec)
		lex.Source = []byte(test.source)
		lex.PositionRanges = []int{len(test.source)}
		lex.Line = 1

		ok, eof := lex.analyzeAnnotation()
		if !ok {
			t.Fatalf("could not analyze annotation")
		}
		if eof {
			t.Fatalf("unexpected end of file")
		}

		if len(lex.Tokens) != 1 {
			t.Fatalf("unexpected token stream: got %v", lex.Tokens)
		}

		actual := lex.Tokens[0]
		if actual != test.expect {
			t.Fatalf("unexpected token (%v): got %v", test.expect, actual)
		}
	}
}

func TestAnalyzeStandalone(t *testing.T) {
	tests := []struct {
		source string
		expect token.Token
	}{
		{
			source: `(`,
			expect: token.Token{Value: "(", Type: token.LeftParen, End: 1},
		},
		{
			source: `[`,
			expect: token.Token{Value: "[", Type: token.LeftBracket, End: 1},
		},
		{
			source: `{`,
			expect: token.Token{Value: "{", Type: token.LeftBrace, End: 1},
		},
		{
			source: `,`,
			expect: token.Token{Value: ",", Type: token.Comma, End: 1},
		},
		{
			source: `)`,
			expect: token.Token{Value: ")", Type: token.RightParen, End: 1},
		},
		{
			source: `]`,
			expect: token.Token{Value: "]", Type: token.RightBracket, End: 1},
		},
		{
			source: `}`,
			expect: token.Token{Value: "}", Type: token.RightBrace, End: 1},
		},
	}

	for _, test := range tests {
		lex := Init(source.StdinSpec)
		lex.Source = []byte(test.source)
		lex.PositionRanges = []int{len(test.source)}
		lex.Line = 1

		ok, eof := lex.analyzeStandalone()
		if !ok {
			t.Fatalf("could not analyze standalone symbol")
		}
		if eof {
			t.Fatalf("unexpected end of file")
		}

		if len(lex.Tokens) != 1 {
			t.Fatalf("unexpected token stream: got %v", lex.Tokens)
		}

		actual := lex.Tokens[0]
		if actual != test.expect {
			t.Fatalf("unexpected token (%v): got %v", test.expect, actual)
		}
	}
}

func TestMatchKeyword(t *testing.T) {
	tests := []struct {
		source string
		expect token.Type
	}{
		{
			source: `:`,
			expect: token.Colon,
		},
		{
			source: `=`,
			expect: token.Equal,
		},
		{
			source: `|`,
			expect: token.Bar,
		},
		{
			source: `=>`,
			expect: token.ThickArrow,
		},
		{
			source: `->`,
			expect: token.Arrow,
		},
		{
			source: `:=`,
			expect: token.ColonEqual,
		},
		{
			source: `\`,
			expect: token.Backslash,
		},
		{
			source: `..`,
			expect: token.DotDot,
		},
		{
			source: `.`,
			expect: token.Dot,
		},
		{
			source: `alias`,
			expect: token.Alias,
		},
		{
			source: `derives`,
			expect: token.Derives,
		},
		{
			source: `with`,
			expect: token.With,
		},
		{
			source: `let`,
			expect: token.Let,
		},
		{
			source: `import`,
			expect: token.Import,
		},
		{
			source: `in`,
			expect: token.In,
		},
		{
			source: `module`,
			expect: token.Module,
		},
		{
			source: `use`,
			expect: token.Use,
		},
		{
			source: `spec`,
			expect: token.Spec,
		},
		{
			source: `where`,
			expect: token.Where,
		},
		{
			source: `with`,
			expect: token.With,
		},
		{
			source: `as`,
			expect: token.As,
		},
		{
			source: `of`,
			expect: token.Of,
		},
		{
			source: `erase`,
			expect: token.Erase,
		},
		{
			source: `once`,
			expect: token.Once,
		},
		// identifiers
		{
			source: `not`,
			expect: token.Id,
		},
		{
			source: `/`,
			expect: token.Id,
		},
	}

	for _, test := range tests {
		lex := Init(source.StdinSpec)
		lex.Source = []byte(test.source)
		lex.PositionRanges = []int{len(test.source)}
		lex.Line = 1

		actual := matchKeyword(test.source, token.Id)
		if actual != test.expect {
			t.Fatalf("token type (%v): got %v", test.expect, actual)
		}
	}
}

func TestCharNum(t *testing.T) {
	const src string = "test.source"
	lex := Init(source.StdinSpec)
	lex.Source = []byte(src)
	lex.PositionRanges = []int{len(src)}
	lex.Line = 1

	for i := 0; i < len(src); i++ {
		actual, eof := lex.charNumber()
		if actual != i+1 {
			t.Fatalf("expected %d: got %d", i+1, actual)
		}
		if eof {
			t.Fatalf("unexpected end of file")
		}
		lex.Pos++
	}

	_, eof := lex.charNumber()
	if !eof {
		t.Fatalf("expected end of file")
	}
}

func TestAnalyzeIdentifier(t *testing.T) {
	tests := []struct {
		source string
		expect token.Token
	}{
		{
			source: `+`,
			expect: token.Token{Value: "+", Type: token.Id, End: 1},
		},
		{
			source: `+{`,
			expect: token.Token{Value: "+", Type: token.Id, End: 1},
		},
		{
			source: `+=`,
			expect: token.Token{Value: "+=", Type: token.Id, End: 2},
		},
		// {
		// 	source: `(+)`,
		// 	expect: token.Token{Value: "(+)", Type: token.Affixed, End: 3},
		// },
		// {
		// 	source: `(>>=)`,
		// 	expect: token.Token{Value: "(>>=)", Type: token.Affixed, End: 5},
		// },
		// {
		// 	source: `(mod)`,
		// 	expect: token.Token{Value: "(mod)", Type: token.Affixed, End: 5},
		// },
		{
			source: `mod`,
			expect: token.Token{Value: "mod", Type: token.Id, End: 3},
		},
		{
			source: `?mod`,
			expect: token.Token{Value: "?mod", Type: token.Hole, End: 4},
		},
		{
			source: `!`,
			expect: token.Token{Value: "!", Type: token.Id, End: 1},
		},
		// {
		// 	source: `(!)`,
		// 	expect: token.Token{Value: "(!)", Type: token.Affixed, End: 3},
		// },
		{
			source: `a`,
			expect: token.Token{Value: "a", Type: token.Id, End: 1},
		},
		{
			source: `a'`,
			expect: token.Token{Value: "a'", Type: token.Id, End: 2},
		},
		{
			source: `a''`,
			expect: token.Token{Value: "a''", Type: token.Id, End: 3},
		},
		{
			source: `a1`,
			expect: token.Token{Value: "a1", Type: token.Id, End: 2},
		},
		{
			source: `a12`,
			expect: token.Token{Value: "a12", Type: token.Id, End: 3},
		},
		{
			source: `a1'`,
			expect: token.Token{Value: "a1'", Type: token.Id, End: 3},
		},
	}

	for _, test := range tests {
		lex := Init(source.StdinSpec)
		lex.Source = []byte(test.source)
		lex.PositionRanges = []int{len(test.source)}
		lex.Line = 1

		ok, eof := lex.analyzeSymbol()
		if errors.PrintErrors(lex.messages) != 0 {
			t.Fatalf("failed with above errors")
		}
		if !ok {
			t.Fatalf("could not analyze identifier")
		}
		if eof {
			t.Fatalf("unexpected end of file")
		}

		if len(lex.Tokens) != 1 {
			t.Fatalf("unexpected token stream: got %v", lex.Tokens)
		}

		actual := lex.Tokens[0]
		if actual != test.expect {
			t.Fatalf("unexpected token (%v): got %v", test.expect.Debug(), actual.Debug())
		}
	}
}

func TestAnalyzeUnderscore(t *testing.T) {
	expect := token.Token{Value: "_", Type: token.Underscore, End: 1}
	lex := Init(source.StdinSpec)
	lex.Source = []byte(`_`)
	lex.Line = 1

	ok, eof := lex.analyzeUnderscore()
	if !ok {
		t.Fatalf("could not analyze underscore symbol")
	}
	if eof {
		t.Fatalf("unexpected end of file")
	}

	if len(lex.Tokens) != 1 {
		t.Fatalf("unexpected token stream: got %v", lex.Tokens)
	}

	actual := lex.Tokens[0]
	if actual != expect {
		t.Fatalf("unexpected token (%v): got %v", expect.Debug(), actual.Debug())
	}
}

func TestAnalyzeSingleLineComment(t *testing.T) {
	tests := []struct {
		source []byte
		expect token.Token
	}{
		{
			source: []byte(`--`),
			expect: token.Token{Value: "", Type: token.Comment, End: 2},
		},
		{
			source: []byte(`-- `),
			expect: token.Token{Value: " ", Type: token.Comment, End: 3},
		},
		{
			source: []byte(`--comment`),
			expect: token.Token{Value: "comment", Type: token.Comment, End: 9},
		},
		{
			source: []byte(`-- comment`),
			expect: token.Token{Value: " comment", Type: token.Comment, End: 10},
		},
		{
			source: []byte("--\tcomment"),
			expect: token.Token{Value: "\tcomment", Type: token.Comment, End: 10},
		},
		{
			source: []byte("--comment "),
			expect: token.Token{Value: "comment ", Type: token.Comment, End: 10},
		},
		{
			source: []byte("-- a comment"),
			expect: token.Token{Value: " a comment", Type: token.Comment, End: 12},
		},
		{
			source: []byte("-------"),
			expect: token.Token{Value: "-----", Type: token.Comment, End: 7},
		},
	}

	for _, test := range tests {
		lex := Init(source.StdinSpec)
		lex.Source = test.source
		lex.PositionRanges = []int{len(test.source)}
		lex.Line = 1
		lex.SetKeepComments(true)

		ok, eof := lex.analyzeComment()
		if !ok {
			t.Fatalf("could not analyze comment")
		}
		if eof {
			t.Fatalf("unexpected end of file")
		}

		if len(lex.Tokens) != 1 {
			t.Fatalf("unexpected token stream: got %v", lex.Tokens)
		}

		actual := lex.Tokens[0]
		if actual != test.expect {
			t.Fatalf("unexpected token (%v): got %v", test.expect.Debug(), actual.Debug())
		}
	}
}
