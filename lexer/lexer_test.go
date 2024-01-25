// =================================================================================================
// Alex Peters - January 21, 2024
// =================================================================================================
package lexer

import "testing"

func TestWrite(t *testing.T) {
	{
		write := func(lex *Lexer) bool {
			lex.Source = append(lex.Source, []byte("test")...)
			return true
		}
		lex := Lexer{
			write: write,
			Source: []byte("test"),
		}
		expected := 4
		actual := lex.Write()
		if actual != expected {
			t.Fatalf("unexpected return value(%d): got %d", expected, actual)
		}
	}
	
	{
		write := func(*Lexer) bool {
			return false
		}
		lex := Lexer{
			write: write,
			Source: []byte("test"),
		}
		expected := -1
		actual := lex.Write()
		if actual != expected {
			t.Fatalf("unexpected return value: got %d", actual)
		}
	}
}

func TestCurrentSourceChar(t *testing.T) {
	tests := []struct{
		Lexer
		char byte
		ok bool
	}{
		{
			Lexer: Lexer{
				Source: []byte("a"),
				Pos: 0,
				PositionRanges: []int{1},
			},
			char: 'a',
			ok: true,
		},
		{
			Lexer: Lexer{
				Source: []byte("a\nb"),
				Pos: 2,
				PositionRanges: []int{2,3},
			},
			char: 'b',
			ok: true,
		},
		{
			Lexer: Lexer{
				Source: []byte(""),
				Pos: 0,
			},
			char: 0,
			ok: false,
		},
		{
			Lexer: Lexer{
				Source: []byte(""),
				Pos: 1,
			},
			char: 0,
			ok: false,
		},
		{
			Lexer: Lexer{
				Source: []byte("abc"),
				Pos: 0,
				PositionRanges: []int{3},
			},
			char: 'a',
			ok: true,
		},
		{
			Lexer: Lexer{
				Source: []byte("abc"),
				Pos: 1,
				PositionRanges: []int{3},
			},
			char: 'b',
			ok: true,
		},
		{
			Lexer: Lexer{
				Source: []byte("abc"),
				Pos: 2,
				PositionRanges: []int{3},
			},
			char: 'c',
			ok: true,
		},
		{
			Lexer: Lexer{
				Source: []byte("abc"),
				Pos: 3,
				PositionRanges: []int{3},
			},
			char: 0,
			ok: false,
		},
	}

	for _, test := range tests {
		char, ok := test.Lexer.currentSourceChar()
		if test.char != char {
			t.Fatalf("unexpected character: got %c", char)
		}
		if test.ok != ok {
			t.Fatalf("unexpected ok value: got %t", ok)
		}
	}
}

func TestCurrentSourceLine(t *testing.T) {
	tests := []struct{
		Lexer
		line string
		ok bool
	}{
		{
			Lexer: Lexer{
				Line: 0,
			},
			line: "",
			ok: false,
		},
		{
			Lexer: Lexer{
				Source: []byte("a"),
				PositionRanges: []int{1},
				Line: 1,
			},
			line: "a",
			ok: true,
		},
		{
			Lexer: Lexer{
				Source: []byte("a"),
				PositionRanges: []int{1},
				Line: 1,
			},
			line: "a",
			ok: true,
		},
		{
			Lexer: Lexer{
				Source: []byte("a\nb"),
				PositionRanges: []int{2,3},
				Line: 2,
			},
			line: "b",
			ok: true,
		},
		{
			Lexer: Lexer{
				Source: []byte(""),
				Line: 1,
			},
			line: "",
			ok: false,
		},
		{
			Lexer: Lexer{
				Source: []byte("abc"),
				PositionRanges: []int{3},
				Line: 1,
			},
			line: "abc",
			ok: true,
		},
		{
			Lexer: Lexer{
				Source: []byte("abc"),
				PositionRanges: []int{3},
				Line: 1,
			},
			line: "abc",
			ok: true,
		},
		{
			Lexer: Lexer{
				Source: []byte("abc"),
				PositionRanges: []int{3},
				Line: 2,
			},
			line: "",
			ok: false,
		},
	}

	for _, test := range tests {
		line, ok := test.Lexer.currentSourceLine()
		if test.line != line {
			t.Fatalf("unexpected string (Lexer=%v, line=\"%s\", ok=%t): got \"%s\"", (&test.Lexer), test.line, test.ok, line)
		}
		if test.ok != ok {
			t.Fatalf("unexpected ok value: got %t", ok)
		}
	}
}