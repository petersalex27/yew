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
		lex := Init(StdinSpec)
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
		lex := Init(StdinSpec)
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
		lex := Init(StdinSpec)
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
			source:         token.Token{Value: "annotation_a", Type: token.Affixed},
			expect:         token.Token{Value: "annotation_a", Type: token.Affixed},
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
			expect: token.Token{Value: "annotation", Type: token.At, End: 11},
		},
		{
			source: `@Annotation`,
			expect: token.Token{Value: "Annotation", Type: token.At, End: 11},
		},
		{
			source: `@annotation _ stuff`,
			expect: token.Token{Value: "annotation", Type: token.At, End: 11},
		},
	}

	for _, test := range tests {
		lex := Init(StdinSpec)
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
		lex := Init(StdinSpec)
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
		lex.Source = []byte(test.source)
		lex.PositionRanges = []int{len(test.source)}
		lex.Line = 1

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
		{
			source: `_+_`,
			expect: token.Token{Value: "_+_", Type: token.Affixed, End: 3},
		},
		{
			source: `_>>=_`,
			expect: token.Token{Value: "_>>=_", Type: token.Affixed, End: 5},
		},
		{
			source: `_mod_`,
			expect: token.Token{Value: "_mod_", Type: token.Affixed, End: 5},
		},
		{
			source: `mod`,
			expect: token.Token{Value: "mod", Type: token.Id, End: 3},
		},
		{
			source: `!`,
			expect: token.Token{Value: "!", Type: token.Id, End: 1},
		},
		{
			source: `_!`,
			expect: token.Token{Value: "_!", Type: token.Affixed, End: 2},
		},
		{
			source: `!_`,
			expect: token.Token{Value: "!_", Type: token.Affixed, End: 2},
		},
		{
			source: `if_then_else_`,
			expect: token.Token{Value: "if_then_else_", Type: token.Affixed, End: 13},
		},
	}

	for _, test := range tests {
		lex := Init(StdinSpec)
		lex.Source = []byte(test.source)
		lex.PositionRanges = []int{len(test.source)}
		lex.Line = 1

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
	expect := token.Token{Value: "_", Type: token.Underscore, End: 1}
	lex := Init(StdinSpec)
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
		t.Fatalf("unexpected token (%v): got %v", expect, actual)
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
		lex := Init(StdinSpec)
		lex.Source = test.source
		lex.PositionRanges = []int{len(test.source)}
		lex.Line = 1

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
		source []byte
		pos    []int
		expect token.Token
	}{
		{
			source: []byte(`-**-`),
			pos: []int{4},
			expect: token.Token{Value: "", Type: token.Comment, End: 4},
		},
		{
			source: []byte(`-*****-`),
			pos: []int{7},
			expect: token.Token{Value: "***", Type: token.Comment, End: 7},
		},
		{
			source: []byte(`-*-*-`),
			pos: []int{5},
			expect: token.Token{Value: "-", Type: token.Comment, End: 5},
		},
		{
			source: []byte(`-*--*-`),
			pos: []int{6},
			expect: token.Token{Value: "--", Type: token.Comment, End: 6},
		},
		{
			source: []byte(`-*-**-`),
			pos: []int{6},
			expect: token.Token{Value: "-*", Type: token.Comment, End: 6},
		},
		{
			source: []byte(`-*-**-a`),
			pos: []int{7},
			expect: token.Token{Value: "-*", Type: token.Comment, End: 6},
		},
		{
			source: []byte(`-*a comment*-`),
			pos: []int{13},
			expect: token.Token{Value: "a comment", Type: token.Comment, End: 13},
		},
		{
			source: []byte(`-* a comment *-`),
			pos: []int{15},
			expect: token.Token{Value: " a comment ", Type: token.Comment, End: 15},
		},
		{
			source: []byte("-*\n*-"),
			pos: []int{3,5},
			expect: token.Token{Value: "\n", Type: token.Comment, End: 5},
		},
		{
			source: []byte("-* this\nis a comment*-"),
			pos: []int{8,22},
			expect: token.Token{Value: " this\nis a comment", Type: token.Comment, End: 22},
		},
		{
			source: []byte("-*\nthis\nis\na\ncomment*-"),
			pos: []int{3,9,11,13,22},
			expect: token.Token{Value: "\nthis\nis\na\ncomment", Type: token.Comment, End: 22},
		},
		{
			source: []byte("-*\nthis\nis\na\ncomment\n*-"),
			pos: []int{3,9,11,13,23},
			expect: token.Token{Value: "\nthis\nis\na\ncomment\n", Type: token.Comment, End: 23},
		},
	}

	for _, test := range tests {
		lex := Init(StdinSpec)
		lex.Source = test.source
		lex.PositionRanges = test.pos
		lex.Line = 1

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
