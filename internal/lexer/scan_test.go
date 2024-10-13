// =================================================================================================
// Alex Peters - January 22, 2024
// =================================================================================================

package lexer

import (
	"testing"

	"github.com/petersalex27/yew/api/token"
	"github.com/petersalex27/yew/api/util"
)

func TestAnalyzeNumber(t *testing.T) {
	tests := []struct {
		source string
		pos    []int
		expect token.Token
	}{
		{
			source: "1",
			expect: token.Token{Value: "1", Typ: token.IntValue, End: 1},
		},
		{
			source: "0x1",
			expect: token.Token{Value: "0x1", Typ: token.IntValue, End: 3},
		},
		{
			source: "0xa",
			expect: token.Token{Value: "0xa", Typ: token.IntValue, End: 3},
		},
		{
			source: "0Xa",
			expect: token.Token{Value: "0Xa", Typ: token.IntValue, End: 3},
		},
		{
			source: "0o1",
			expect: token.Token{Value: "0o1", Typ: token.IntValue, End: 3},
		},
		{
			source: "0O1",
			expect: token.Token{Value: "0O1", Typ: token.IntValue, End: 3},
		},
		{
			source: "0b1",
			expect: token.Token{Value: "0b1", Typ: token.IntValue, End: 3},
		},
		{
			source: "0B1",
			expect: token.Token{Value: "0B1", Typ: token.IntValue, End: 3},
		},
		{
			source: "1.0",
			expect: token.Token{Value: "1.0", Typ: token.FloatValue, End: 3},
		},
		{
			source: "1e1",
			expect: token.Token{Value: "1e1", Typ: token.FloatValue, End: 3},
		},
		{
			source: "1E1",
			expect: token.Token{Value: "1E1", Typ: token.FloatValue, End: 3},
		},
		{
			source: "1e+1",
			expect: token.Token{Value: "1e+1", Typ: token.FloatValue, End: 4},
		},
		{
			source: "1e-1",
			expect: token.Token{Value: "1e-1", Typ: token.FloatValue, End: 4},
		},
		{
			source: "1.0e1",
			expect: token.Token{Value: "1.0e1", Typ: token.FloatValue, End: 5},
		},
		{
			source: "0_1",
			expect: token.Token{Value: "1", Typ: token.IntValue, End: 3},
		},
		{
			source: "00_1",
			expect: token.Token{Value: "1", Typ: token.IntValue, End: 4},
		},
		{
			source: "11_1",
			expect: token.Token{Value: "111", Typ: token.IntValue, End: 4},
		},
		{
			source: "0x1_1",
			expect: token.Token{Value: "0x11", Typ: token.IntValue, End: 5},
		},
		{
			source: "0o1_1",
			expect: token.Token{Value: "0o11", Typ: token.IntValue, End: 5},
		},
		{
			source: "0b1_1",
			expect: token.Token{Value: "0b11", Typ: token.IntValue, End: 5},
		},
	}

	for _, test := range tests {
		lex := Init(util.StringSource(test.source))

		actual := lex.number()
		if !actual.Equals(test.expect) {
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
			expect: token.Token{Value: "a", Typ: token.CharValue, End: 3},
		},
		{
			source: `' '`,
			expect: token.Token{Value: " ", Typ: token.CharValue, End: 3},
		},
		{
			source: `'@'`,
			expect: token.Token{Value: "@", Typ: token.CharValue, End: 3},
		},
		{
			source: `'\n'`,
			expect: token.Token{Value: "\n", Typ: token.CharValue, End: 4},
		},
		{
			source: `'\t'`,
			expect: token.Token{Value: "\t", Typ: token.CharValue, End: 4},
		},
		{
			source: `'\a'`,
			expect: token.Token{Value: "\a", Typ: token.CharValue, End: 4},
		},
		{
			source: `'\b'`,
			expect: token.Token{Value: "\b", Typ: token.CharValue, End: 4},
		},
		{
			source: `'\v'`,
			expect: token.Token{Value: "\v", Typ: token.CharValue, End: 4},
		},
		{
			source: `'\f'`,
			expect: token.Token{Value: "\f", Typ: token.CharValue, End: 4},
		},
		{
			source: `'\r'`,
			expect: token.Token{Value: "\r", Typ: token.CharValue, End: 4},
		},
		{
			source: `'\''`,
			expect: token.Token{Value: "'", Typ: token.CharValue, End: 4},
		},
		{
			source: `'\\'`,
			expect: token.Token{Value: "\\", Typ: token.CharValue, End: 4},
		},
	}

	for _, test := range tests {
		lex := Init(util.StringSource(test.source))

		actual := lex.char()
		if !actual.Equals(test.expect) {
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
			expect: token.Token{Value: "", Typ: token.StringValue, End: 2},
		},
		{
			source: `" "`,
			expect: token.Token{Value: " ", Typ: token.StringValue, End: 3},
		},
		{
			source: `"--"`,
			expect: token.Token{Value: "--", Typ: token.StringValue, End: 4},
		},
		{
			source: `"this is a string"`,
			expect: token.Token{Value: "this is a string", Typ: token.StringValue, End: 18},
		},
		{
			source: `"\n\t\a\b\v\f\r\"\\"`,
			expect: token.Token{Value: "\n\t\a\b\v\f\r\"\\", Typ: token.StringValue, End: 20},
		},

		// TODO: prompted strings

		// import paths
		{
			source: `"a/b"`,
			expect: token.Token{Value: "a/b", Typ: token.ImportPath, End: 5},
		},
		{
			source: `"a/b/c"`,
			expect: token.Token{Value: "a/b/c", Typ: token.ImportPath, End: 7},
		},
		// almost, but not quite, import paths
		{
			source: `"a/"`,
			expect: token.Token{Value: "a/", Typ: token.StringValue, End: 4},
		},
		{
			source: `"/a"`,
			expect: token.Token{Value: "/a", Typ: token.StringValue, End: 4},
		},
		{
			source: `"a/b."`,
			expect: token.Token{Value: "a/b.", Typ: token.StringValue, End: 6},
		},
		{
			source: `"a/b/c.d"`,
			expect: token.Token{Value: "a/b/c.d", Typ: token.StringValue, End: 9},
		},
		{
			source: `"a/b/c.d/e"`,
			expect: token.Token{Value: "a/b/c.d/e", Typ: token.StringValue, End: 11},
		},
	}

	for _, test := range tests {
		lex := Init(util.StringSource(test.source))

		actual := lex.stringLiteral()
		if !actual.Equals(test.expect) {
			t.Fatalf("unexpected token (%v): got %v", test.expect.Debug(), actual.Debug())
		}
	}
}

func TestAnalyzeAnnotation(t *testing.T) {
	tests := []struct {
		source string
		expect token.Token
	}{
		{
			source: `--@annotation`,
			expect: token.Token{Value: "annotation", Typ: token.FlatAnnotation, Start: 0, End: 13},
		},
		{
			source: `-- @annotation`,
			expect: token.Token{Value: "annotation", Typ: token.FlatAnnotation, Start: 0, End: 14},
		},
		{
			source: `-- @ annotation`,
			expect: token.Token{Value: "annotation", Typ: token.FlatAnnotation, Start: 0, End: 15},
		},
		{
			source: "-- @\tannotation",
			expect: token.Token{Value: "annotation", Typ: token.FlatAnnotation, Start: 0, End: 15},
		},
		{
			source: "-- @\t annotation",
			expect: token.Token{Value: "annotation", Typ: token.FlatAnnotation, Start: 0, End: 16},
		},
		{
			source: `-- @Annotation`,
			expect: token.Token{Value: "Annotation", Typ: token.FlatAnnotation, Start: 0, End: 14},
		},
		{
			source: `-- @annotation _ stuff`,
			expect: token.Token{Value: "annotation _ stuff", Typ: token.FlatAnnotation, Start: 0, End: 22},
		},
		{
			source: `-- @ annotation _ stuff`,
			expect: token.Token{Value: "annotation _ stuff", Typ: token.FlatAnnotation, Start: 0, End: 23},
		},
	}

	for _, test := range tests {
		lex := Init(util.StringSource(test.source))

		actual := lex.analyzeComment()
		if !actual.Equals(test.expect) {
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
			source: `()`,
			expect: token.Token{Value: `()`, Typ: token.EmptyParenEnclosure, End: 2},
		},
		{
			source: `[]`,
			expect: token.Token{Value: `[]`, Typ: token.EmptyBracketEnclosure, End: 2},
		},
		{
			source: `[@`,
			expect: token.Token{Value: `[@`, Typ: token.LeftBracketAt, End: 2},
		},
		{
			source: `(`,
			expect: token.Token{Value: "(", Typ: token.LeftParen, End: 1},
		},
		{
			source: `[`,
			expect: token.Token{Value: "[", Typ: token.LeftBracket, End: 1},
		},
		{
			source: `{`,
			expect: token.Token{Value: "{", Typ: token.LeftBrace, End: 1},
		},
		{
			source: `,`,
			expect: token.Token{Value: ",", Typ: token.Comma, End: 1},
		},
		{
			source: `)`,
			expect: token.Token{Value: ")", Typ: token.RightParen, End: 1},
		},
		{
			source: `]`,
			expect: token.Token{Value: "]", Typ: token.RightBracket, End: 1},
		},
		{
			source: `}`,
			expect: token.Token{Value: "}", Typ: token.RightBrace, End: 1},
		},
	}

	for _, test := range tests {
		lex := Init(util.StringSource(test.source))

		actual := lex.standalone(func(*Lexer) token.Token {
			t.Fatalf("standalone called fallback, but expected actual standalone")
			return token.Token{}
		})
		if !actual.Equals(test.expect) {
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
			source: `()`,
			expect: token.EmptyParenEnclosure,
		},
		{
			source: `[]`,
			expect: token.EmptyBracketEnclosure,
		},
		{
			source: `[@`,
			expect: token.LeftBracketAt,
		},
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
			source: `deriving`,
			expect: token.Deriving,
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
			source: `using`,
			expect: token.Using,
		},
		{
			source: `spec`,
			expect: token.Spec,
		},
		{
			source: `inst`,
			expect: token.Inst,
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
			source: `from`,
			expect: token.From,
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
		{
			source: `impossible`,
			expect: token.Impossible,
		},
		{
			source: `requiring`,
			expect: token.Requiring,
		},
		{
			source: `forall`,
			expect: token.Forall,
		},
		{
			source: `public`,
			expect: token.Public,
		},
		{
			source: `open`,
			expect: token.Open,
		},
		{
			source: `auto`,
			expect: token.Auto,
		},
		{
			source: `inst`,
			expect: token.Inst,
		},
		{
			source: `syntax`,
			expect: token.Syntax,
		},
		{
			source: `requiring`,
			expect: token.Requiring,
		},
		// reserved, but unused
		{
			source: `pattern`,
			expect: token.Pattern,
		},
		{
			source: `term`,
			expect: token.Term,
		},
		{
			source: `ref`,
			expect: token.Ref,
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
		actual := matchKeyword(test.source, token.Id)
		if actual != test.expect {
			t.Fatalf("token type (%v): got %v", test.expect, actual)
		}
	}
}

func TestCharNum(t *testing.T) {
	const src string = "test.source"
	lex := Init(util.StringSource(src))

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
			expect: token.Token{Value: "+", Typ: token.Id, End: 1},
		},
		{
			source: `+{`,
			expect: token.Token{Value: "+", Typ: token.Id, End: 1},
		},
		{
			source: `++{`,
			expect: token.Token{Value: "++", Typ: token.Id, End: 2},
		},
		{
			source: `+=`,
			expect: token.Token{Value: "+=", Typ: token.Id, End: 2},
		},
		{
			source: `mod`,
			expect: token.Token{Value: "mod", Typ: token.Id, End: 3},
		},
		{
			source: `?mod`,
			expect: token.Token{Value: "?mod", Typ: token.Hole, End: 4},
		},
		{
			source: `?Mod`,
			expect: token.Token{Value: "?Mod", Typ: token.Hole, End: 4},
		},
		{
			source: `!`,
			expect: token.Token{Value: "!", Typ: token.Id, End: 1},
		},
		{
			source: `a`,
			expect: token.Token{Value: "a", Typ: token.Id, End: 1},
		},
		{
			source: `a'`,
			expect: token.Token{Value: "a'", Typ: token.Id, End: 2},
		},
		{
			source: `a''`,
			expect: token.Token{Value: "a''", Typ: token.Id, End: 3},
		},
		{
			source: `a1`,
			expect: token.Token{Value: "a1", Typ: token.Id, End: 2},
		},
		{
			source: `a12`,
			expect: token.Token{Value: "a12", Typ: token.Id, End: 3},
		},
		{
			source: `a1'`,
			expect: token.Token{Value: "a1'", Typ: token.Id, End: 3},
		},
	}

	for _, test := range tests {
		lex := Init(util.StringSource(test.source))

		actual := lex.symbol()
		if !actual.Equals(test.expect) {
			t.Fatalf("unexpected token (%v): got %v", test.expect.Debug(), actual.Debug())
		}
	}
}

func TestAnalyzeInfix(t *testing.T) {
	tests := []struct {
		source string
		expect token.Token
	}{
		{
			source: `(+)`,
			expect: token.Token{Value: "+", Typ: token.Infix, End: 3},
		},
		{
			source: `(>>=)`,
			expect: token.Token{Value: ">>=", Typ: token.Infix, End: 5},
		},
		{
			source: `(mod)`,
			expect: token.Token{Value: "mod", Typ: token.Infix, End: 5},
		},
		{
			source: `(List)`,
			expect: token.Token{Value: "List", Typ: token.Infix, End: 6},
		},
		// allow even keywords to be enclosed in infix parens
		{
			source: `(->)`,
			expect: token.Token{Value: "->", Typ: token.Infix, End: 4},
		},
		// method symbols
		{
			source: `(.method)`,
			expect: token.Token{Value: "method", Typ: token.MethodSymbol, End: 9},
		},
		{
			source: `(.Method)`,
			expect: token.Token{Value: "Method", Typ: token.MethodSymbol, End: 9},
		},
		{
			source: `(.?)`,
			expect: token.Token{Value: "?", Typ: token.MethodSymbol, End: 4},
		},
	}

	for _, test := range tests {
		lex := Init(util.StringSource(test.source))

		actual := lex.infix()
		if !actual.Equals(test.expect) {
			t.Fatalf("unexpected token (%v): got %v", test.expect.Debug(), actual.Debug())
		}
	}
}

func TestAnalyzeUnderscore(t *testing.T) {
	expect := token.Token{Value: "_", Typ: token.Underscore, End: 1}
	lex := Init(util.StringSource("_"))

	actual := lex.underscore()
	if !actual.Equals(expect) {
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
			expect: token.Token{Value: "", Typ: token.Comment, End: 2},
		},
		{
			source: []byte(`-- `),
			expect: token.Token{Value: " ", Typ: token.Comment, End: 3},
		},
		{
			source: []byte(`--comment`),
			expect: token.Token{Value: "comment", Typ: token.Comment, End: 9},
		},
		{
			source: []byte(`-- comment`),
			expect: token.Token{Value: " comment", Typ: token.Comment, End: 10},
		},
		{
			source: []byte("--\tcomment"),
			expect: token.Token{Value: "\tcomment", Typ: token.Comment, End: 10},
		},
		{
			source: []byte("--comment "),
			expect: token.Token{Value: "comment ", Typ: token.Comment, End: 10},
		},
		{
			source: []byte("-- a comment"),
			expect: token.Token{Value: " a comment", Typ: token.Comment, End: 12},
		},
		{
			source: []byte("-------"),
			expect: token.Token{Value: "-----", Typ: token.Comment, End: 7},
		},
	}

	for _, test := range tests {
		lex := Init(util.BytesSource(test.source))
		lex.SetKeepComments(true)

		actual := lex.analyzeComment()
		if !actual.Equals(test.expect) {
			t.Fatalf("unexpected token (%v): got %v", test.expect.Debug(), actual.Debug())
		}
	}
}
func TestAnalyzeRawString(t *testing.T) {
	tests := []struct {
		source string
		expect token.Token
	}{
		{
			source: "``",
			expect: token.Token{Value: "", Typ: token.RawStringValue, End: 2},
		},
		{
			source: "`raw string`",
			expect: token.Token{Value: "raw string", Typ: token.RawStringValue, End: 12},
		},
		{
			source: "`with\nnew\nlines`",
			expect: token.Token{Value: "with\nnew\nlines", Typ: token.RawStringValue, End: 16},
		},
		{
			source: "`with\ttabs`",
			expect: token.Token{Value: "with\ttabs", Typ: token.RawStringValue, End: 11},
		},
		{
			source: "`with special chars !@#$%^&*()`",
			expect: token.Token{Value: "with special chars !@#$%^&*()", Typ: token.RawStringValue, End: 31},
		},
	}

	for _, test := range tests {
		lex := Init(util.StringSource(test.source))

		actual := lex.rawStringLiteral()
		if !actual.Equals(test.expect) {
			t.Fatalf("unexpected token (`%v`): got `%v`", test.expect, actual)
		}
	}
}