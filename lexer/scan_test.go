// =================================================================================================
// Alex Peters - January 22, 2024
// =================================================================================================

package lexer

import (
	"testing"

	"github.com/petersalex27/yew/token"
)

func TestAnalyzeNumber(t *testing.T) {
	tests := []struct {
		source string
		expect token.Token
	}{
		{
			source: "1",
			expect: token.Token{Value: "1", Type: token.IntValue, Start: 1, End: 2, Line: 1},
		},
		{
			source: "0x1",
			expect: token.Token{Value: "0x1", Type: token.IntValue, Start: 1, End: 4, Line: 1},
		},
		{
			source: "0xa",
			expect: token.Token{Value: "0xa", Type: token.IntValue, Start: 1, End: 4, Line: 1},
		},
		{
			source: "0Xa",
			expect: token.Token{Value: "0Xa", Type: token.IntValue, Start: 1, End: 4, Line: 1},
		},
		{
			source: "0o1",
			expect: token.Token{Value: "0o1", Type: token.IntValue, Start: 1, End: 4, Line: 1},
		},
		{
			source: "0O1",
			expect: token.Token{Value: "0O1", Type: token.IntValue, Start: 1, End: 4, Line: 1},
		},
		{
			source: "0b1",
			expect: token.Token{Value: "0b1", Type: token.IntValue, Start: 1, End: 4, Line: 1},
		},
		{
			source: "0B1",
			expect: token.Token{Value: "0B1", Type: token.IntValue, Start: 1, End: 4, Line: 1},
		},
		{
			source: "1.0",
			expect: token.Token{Value: "1.0", Type: token.FloatValue, Start: 1, End: 4, Line: 1},
		},
		{
			source: "1e1",
			expect: token.Token{Value: "1e1", Type: token.FloatValue, Start: 1, End: 4, Line: 1},
		},
		{
			source: "1E1",
			expect: token.Token{Value: "1E1", Type: token.FloatValue, Start: 1, End: 4, Line: 1},
		},
		{
			source: "1e+1",
			expect: token.Token{Value: "1e+1", Type: token.FloatValue, Start: 1, End: 5, Line: 1},
		},
		{
			source: "1e-1",
			expect: token.Token{Value: "1e-1", Type: token.FloatValue, Start: 1, End: 5, Line: 1},
		},
		{
			source: "1.0e1",
			expect: token.Token{Value: "1.0e1", Type: token.FloatValue, Start: 1, End: 6, Line: 1},
		},
		{
			source: "0_1",
			expect: token.Token{Value: "1", Type: token.IntValue, Start: 1, End: 4, Line: 1},
		},
		{
			source: "00_1",
			expect: token.Token{Value: "1", Type: token.IntValue, Start: 1, End: 5, Line: 1},
		},
		{
			source: "11__1",
			expect: token.Token{Value: "111", Type: token.IntValue, Start: 1, End: 6, Line: 1},
		},
		{
			source: "0x1_1",
			expect: token.Token{Value: "0x11", Type: token.IntValue, Start: 1, End: 6, Line: 1},
		},
		{
			source: "0o1_1",
			expect: token.Token{Value: "0o11", Type: token.IntValue, Start: 1, End: 6, Line: 1},
		},
		{
			source: "0b1_1",
			expect: token.Token{Value: "0b11", Type: token.IntValue, Start: 1, End: 6, Line: 1},
		},
	}

	for _, test := range tests {
		lex := Init(StdinSpec)
		lex.Source = []string{test.source}
		lex.Line, lex.Char = 1, 1

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
			expect: token.Token{Value: "a", Type: token.CharValue, Start: 1, End: 4, Line: 1},
		},
		{
			source: `' '`,
			expect: token.Token{Value: " ", Type: token.CharValue, Start: 1, End: 4, Line: 1},
		},
		{
			source: `'@'`,
			expect: token.Token{Value: "@", Type: token.CharValue, Start: 1, End: 4, Line: 1},
		},
		{
			source: `'\n'`,
			expect: token.Token{Value: "\n", Type: token.CharValue, Start: 1, End: 5, Line: 1},
		},
		{
			source: `'\t'`,
			expect: token.Token{Value: "\t", Type: token.CharValue, Start: 1, End: 5, Line: 1},
		},
		{
			source: `'\a'`,
			expect: token.Token{Value: "\a", Type: token.CharValue, Start: 1, End: 5, Line: 1},
		},
		{
			source: `'\b'`,
			expect: token.Token{Value: "\b", Type: token.CharValue, Start: 1, End: 5, Line: 1},
		},
		{
			source: `'\v'`,
			expect: token.Token{Value: "\v", Type: token.CharValue, Start: 1, End: 5, Line: 1},
		},
		{
			source: `'\f'`,
			expect: token.Token{Value: "\f", Type: token.CharValue, Start: 1, End: 5, Line: 1},
		},
		{
			source: `'\r'`,
			expect: token.Token{Value: "\r", Type: token.CharValue, Start: 1, End: 5, Line: 1},
		},
		{
			source: `'\''`,
			expect: token.Token{Value: "'", Type: token.CharValue, Start: 1, End: 5, Line: 1},
		},
		{
			source: `'\\'`,
			expect: token.Token{Value: "\\", Type: token.CharValue, Start: 1, End: 5, Line: 1},
		},
	}

	for _, test := range tests {
		lex := Init(StdinSpec)
		lex.Source = []string{test.source}
		lex.Line, lex.Char = 1, 1

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
			expect: token.Token{Value: "", Type: token.StringValue, Start: 1, End: 3, Line: 1},
		},
		{
			source: `" "`,
			expect: token.Token{Value: " ", Type: token.StringValue, Start: 1, End: 4, Line: 1},
		},
		{
			source: `"--"`,
			expect: token.Token{Value: "--", Type: token.StringValue, Start: 1, End: 5, Line: 1},
		},
		{
			source: `"this is a string"`,
			expect: token.Token{Value: "this is a string", Type: token.StringValue, Start: 1, End: 19, Line: 1},
		},
		{
			source: `"\n\t\a\b\v\f\r\"\\"`,
			expect: token.Token{Value: "\n\t\a\b\v\f\r\"\\", Type: token.StringValue, Start: 1, End: 21, Line: 1},
		},
	}

	for _, test := range tests {
		lex := Init(StdinSpec)
		lex.Source = []string{test.source}
		lex.Line, lex.Char = 1, 1

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

func TestFixAnnotation(t *testing.T) {
	tests := []struct {
		source         token.Token
		expect         token.Token
		expectErrorMsg string
	}{
		{
			source: token.Token{Value: "annotation", Type: token.Id},
			expect: token.Token{Value: "annotation", Type: token.At},
		},
		{
			source: token.Token{Value: "Annotation", Type: token.CapId},
			expect: token.Token{Value: "Annotation", Type: token.At},
		},
		{
			source: token.Token{Value: "annotation_a", Type: token.Affixed},
			expect: token.Token{Value: "annotation_a", Type: token.Affixed},
			expectErrorMsg: InvalidAnnotation,
		},
	}

	for _, test := range tests {
		actual, actualErrorMsg := fixAnnotation(test.source)
		if actual != test.expect {
			t.Fatalf("unexpected token (%v): got %v", test.expect, actual)
		}
		if actualErrorMsg != test.expectErrorMsg {
			t.Fatalf("unexpected error message (\"%s\"): got \"%s\"", test.expectErrorMsg, actualErrorMsg)
		}
	}
}

func TestAnalyzeAnnotation(t *testing.T) {
	tests := []struct {
		source string
		expect token.Token
	}{
		{
			source: `@annotation`,
			expect: token.Token{Value: "annotation", Type: token.At, Line: 1, Start: 1, End: 12},
		},
		{
			source: `@Annotation`,
			expect: token.Token{Value: "Annotation", Type: token.At, Line: 1, Start: 1, End: 12},
		},
		{
			source: `@annotation _ stuff`,
			expect: token.Token{Value: "annotation", Type: token.At, Line: 1, Start: 1, End: 12},
		},
	}

	for _, test := range tests {
		lex := Init(StdinSpec)
		lex.Source = []string{test.source}
		lex.Line, lex.Char = 1, 1

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
			expect: token.Token{Value: "(", Type: token.LeftParen, Start: 1, End: 2, Line: 1},
		},
		{
			source: `[`,
			expect: token.Token{Value: "[", Type: token.LeftBracket, Start: 1, End: 2, Line: 1},
		},
		{
			source: `{`,
			expect: token.Token{Value: "{", Type: token.LeftBrace, Start: 1, End: 2, Line: 1},
		},
		{
			source: `,`,
			expect: token.Token{Value: ",", Type: token.Comma, Start: 1, End: 2, Line: 1},
		},
		{
			source: `)`,
			expect: token.Token{Value: ")", Type: token.RightParen, Start: 1, End: 2, Line: 1},
		},
		{
			source: `]`,
			expect: token.Token{Value: "]", Type: token.RightBracket, Start: 1, End: 2, Line: 1},
		},
		{
			source: `}`,
			expect: token.Token{Value: "}", Type: token.RightBrace, Start: 1, End: 2, Line: 1},
		},
	}

	for _, test := range tests {
		lex := Init(StdinSpec)
		lex.Source = []string{test.source}
		lex.Line, lex.Char = 1, 1

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
			source: `trait`,
			expect: token.Trait,
		},
		{
			source: `where`,
			expect: token.Where,
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
		lex := Init(StdinSpec)
		lex.Source = []string{test.source}
		lex.Line, lex.Char = 1, 1

		actual := matchKeyword(test.source, token.Id)
		if actual != test.expect {
			t.Fatalf("token type (%v): got %v", test.expect, actual)
		}
	}
}

func TestMatchId(t *testing.T) {
	tests := []struct {
		line           string
		expectMatched  string
		expectErrorMsg string
		expectIllegal  bool
	}{
		{
			line:           `_+__`,
			expectMatched:  "_+_",
			expectErrorMsg: "",
		},
		{
			line:           `_+_`,
			expectMatched:  "_+_",
			expectErrorMsg: "",
		},
		{
			line:           `if_then_else_`,
			expectMatched:  "if_then_else_",
			expectErrorMsg: "",
		},
		{
			line:           `if_?_||_`,
			expectMatched:  "if_",
			expectErrorMsg: "",
		},
		{
			line:           `_?_or_`,
			expectMatched:  "_?_",
			expectErrorMsg: "",
		},
		{
			line:           `._+_`,
			expectMatched:  "",
			expectErrorMsg: InvalidCharacter,
		},
		{
			line:           `a`,
			expectMatched:  "a",
			expectErrorMsg: "",
		},
		{
			line:           `+`,
			expectMatched:  "+",
			expectErrorMsg: "",
		},
		{
			line:           ``,
			expectMatched:  "",
			expectErrorMsg: "",
			expectIllegal:  true,
		},
	}

	for _, test := range tests {
		actualMatched, actualErrorMsg, actualIllegal := matchId(test.line)
		if test.expectMatched != actualMatched {
			t.Errorf("unexpected matched (\"%s\"): got \"%s\"", test.expectMatched, actualMatched)
		}
		if test.expectErrorMsg != actualErrorMsg {
			t.Errorf("unexpected error message (\"%s\"): got \"%s\"", test.expectErrorMsg, actualErrorMsg)
		}
		if test.expectIllegal != actualIllegal {
			t.Errorf("unexpected error message (%v): got %v", test.expectIllegal, actualIllegal)
		}

		if t.Failed() {
			t.FailNow()
		}
	}
}

func TestCheckAffixed(t *testing.T) {
	tests := []struct {
		line   string
		id     string
		expect string
	}{
		{
			line:   `+__`,
			id:     `+_`,
			expect: InvalidAffixId,
		},
		{
			line:   `_+__`,
			id:     `_+_`,
			expect: InvalidAffixId,
		},
		{
			line:   `_+___`,
			id:     `_+_`,
			expect: InvalidAffixId,
		},
		{
			line:   `abc d`,
			id:     `abc`,
			expect: "",
		},
		{
			line:   `_+_: Num n => n -> n -> n`,
			id:     `_+_`,
			expect: "",
		},
		{
			line:   `_+_ _`,
			id:     `_+_`,
			expect: "",
		},
		{
			line:   `_+_`,
			id:     `_+_`,
			expect: "",
		},
	}

	for _, test := range tests {
		actual := checkAffixed(test.line, test.id)
		if test.expect != actual {
			t.Fatalf("unexpected error message (\"%s\"): got \"%s\"", test.expect, actual)
		}
	}
}

func TestAnalyzeIdentifier(t *testing.T) {
	tests := []struct {
		source string
		expect token.Token
	}{
		{
			source: `+`,
			expect: token.Token{Value: "+", Type: token.Id, Start: 1, End: 2, Line: 1},
		},
		{
			source: `+{`,
			expect: token.Token{Value: "+", Type: token.Id, Start: 1, End: 2, Line: 1},
		},
		{
			source: `+=`,
			expect: token.Token{Value: "+=", Type: token.Id, Start: 1, End: 3, Line: 1},
		},
		{
			source: `_+_`,
			expect: token.Token{Value: "_+_", Type: token.Affixed, Start: 1, End: 4, Line: 1},
		},
		{
			source: `_>>=_`,
			expect: token.Token{Value: "_>>=_", Type: token.Affixed, Start: 1, End: 6, Line: 1},
		},
		{
			source: `_mod_`,
			expect: token.Token{Value: "_mod_", Type: token.Affixed, Start: 1, End: 6, Line: 1},
		},
		{
			source: `mod`,
			expect: token.Token{Value: "mod", Type: token.Id, Start: 1, End: 4, Line: 1},
		},
		{
			source: `!`,
			expect: token.Token{Value: "!", Type: token.Id, Start: 1, End: 2, Line: 1},
		},
		{
			source: `_!`,
			expect: token.Token{Value: "_!", Type: token.Affixed, Start: 1, End: 3, Line: 1},
		},
		{
			source: `!_`,
			expect: token.Token{Value: "!_", Type: token.Affixed, Start: 1, End: 3, Line: 1},
		},
		{
			source: `if_then_else_`,
			expect: token.Token{Value: "if_then_else_", Type: token.Affixed, Start: 1, End: 14, Line: 1},
		},
	}

	for _, test := range tests {
		lex := Init(StdinSpec)
		lex.Source = []string{test.source}
		lex.Line, lex.Char = 1, 1

		ok, eof := lex.analyzeIdentifier()
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
			t.Fatalf("unexpected token (%v): got %v", test.expect, actual)
		}
	}
}

func TestAnalyzeUnderscore(t *testing.T) {
	expect := token.Token{Value: "_", Type: token.Underscore, Start: 1, End: 2, Line: 1}
	lex := Init(StdinSpec)
	lex.Source = []string{`_`}
	lex.Line, lex.Char = 1, 1

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
		t.Fatalf("unexpected token (%v): got %v", expect, actual)
	}
}

func TestAnalyzeSingleLineComment(t *testing.T) {
	tests := []struct {
		source []string
		expect token.Token
	}{
		{
			source: []string{`--`},
			expect: token.Token{Value: "", Type: token.Comment, Start: 1, End: 3, Line: 1},
		},
		{
			source: []string{`-- `},
			expect: token.Token{Value: " ", Type: token.Comment, Start: 1, End: 4, Line: 1},
		},
		{
			source: []string{`--comment`},
			expect: token.Token{Value: "comment", Type: token.Comment, Start: 1, End: 10, Line: 1},
		},
		{
			source: []string{`-- comment`},
			expect: token.Token{Value: " comment", Type: token.Comment, Start: 1, End: 11, Line: 1},
		},
		{
			source: []string{"--\tcomment"},
			expect: token.Token{Value: "\tcomment", Type: token.Comment, Start: 1, End: 11, Line: 1},
		},
		{
			source: []string{"--comment "},
			expect: token.Token{Value: "comment ", Type: token.Comment, Start: 1, End: 11, Line: 1},
		},
		{
			source: []string{"-- a comment"},
			expect: token.Token{Value: " a comment", Type: token.Comment, Start: 1, End: 13, Line: 1},
		},
		{
			source: []string{"-------"},
			expect: token.Token{Value: "-----", Type: token.Comment, Start: 1, End: 8, Line: 1},
		},
	}

	for _, test := range tests {
		lex := Init(StdinSpec)
		lex.Source = test.source
		lex.Line, lex.Char = 1, 1

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
			t.Fatalf("unexpected token (%v): got %v", test.expect, actual)
		}
	}
}

func TestAnalyzeMultiLineComment(t *testing.T) {
	tests := []struct {
		source []string
		expect token.Token
	}{
		{
			source: []string{`-**-`},
			expect: token.Token{Value: "", Type: token.Comment, Start: 1, End: 5, Line: 1},
		},
		{
			source: []string{`-*****-`},
			expect: token.Token{Value: "***", Type: token.Comment, Start: 1, End: 8, Line: 1},
		},
		{
			source: []string{`-*-*-`},
			expect: token.Token{Value: "-", Type: token.Comment, Start: 1, End: 6, Line: 1},
		},
		{
			source: []string{`-*--*-`},
			expect: token.Token{Value: "--", Type: token.Comment, Start: 1, End: 7, Line: 1},
		},
		{
			source: []string{`-*-**-`},
			expect: token.Token{Value: "-*", Type: token.Comment, Start: 1, End: 7, Line: 1},
		},
		{
			source: []string{`-*-**-a`},
			expect: token.Token{Value: "-*", Type: token.Comment, Start: 1, End: 7, Line: 1},
		},
		{
			source: []string{`-*a comment*-`},
			expect: token.Token{Value: "a comment", Type: token.Comment, Start: 1, End: 14, Line: 1},
		},
		{
			source: []string{`-* a comment *-`},
			expect: token.Token{Value: " a comment ", Type: token.Comment, Start: 1, End: 16, Line: 1},
		},
		{
			source: []string{
				"-*\n",
				`*-`,
			},
			expect: token.Token{Value: "\n", Type: token.Comment, Start: 1, End: 3, Line: 2},
		},
		{
			source: []string{
				"-* this\n",
				`is a comment*-`,
			},
			expect: token.Token{Value: " this\nis a comment", Type: token.Comment, Start: 1, End: 15, Line: 2},
		},
		{
			source: []string{
				"-*\n",
				"this\n",
				"is\n",
				"a\n",
				"comment*-",
			},
			expect: token.Token{Value: "\nthis\nis\na\ncomment", Type: token.Comment, Start: 1, End: 10, Line: 5},
		},
		{
			source: []string{
				"-*\n",
				"this\n",
				"is\n",
				"a\n",
				"comment\n",
				"*-",
			},
			expect: token.Token{Value: "\nthis\nis\na\ncomment\n", Type: token.Comment, Start: 1, End: 3, Line: 6},
		},
	}

	for _, test := range tests {
		lex := Init(StdinSpec)
		lex.Source = test.source
		lex.Line, lex.Char = 1, 1

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
			t.Fatalf("unexpected token (%v): got %v", test.expect, actual)
		}
	}
}
