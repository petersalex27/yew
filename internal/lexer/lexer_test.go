// =================================================================================================
// Alex Peters - January 21, 2024
// =================================================================================================
package lexer

import (
	"testing"

	"github.com/petersalex27/yew/source"
)

func TestWrite(t *testing.T) {
	{
		write := func(lex *Lexer) bool {
			lex.Source = append(lex.Source, []byte("test")...)
			return true
		}
		lex := Lexer{
			write:      write,
			SourceCode: source.SourceCode{Source: []byte("test")},
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
			write:      write,
			SourceCode: source.SourceCode{Source: []byte("test")},
		}
		expected := -1
		actual := lex.Write()
		if actual != expected {
			t.Fatalf("unexpected return value: got %d", actual)
		}
	}
}

func TestCurrentSourceChar(t *testing.T) {
	tests := []struct {
		Lexer
		char byte
		ok   bool
	}{
		{
			Lexer: Lexer{
				SourceCode: source.SourceCode{
					Source:         []byte("a"),
					PositionRanges: []int{1},
				},
				Pos: 0,
			},
			char: 'a',
			ok:   true,
		},
		{
			Lexer: Lexer{
				SourceCode: source.SourceCode{
					Source:         []byte("a\nb"),
					PositionRanges: []int{2, 3},
				},
				Pos: 2,
			},
			char: 'b',
			ok:   true,
		},
		{
			Lexer: Lexer{
				SourceCode: source.SourceCode{Source: []byte("")},
				Pos:        0,
			},
			char: 0,
			ok:   false,
		},
		{
			Lexer: Lexer{
				SourceCode: source.SourceCode{Source: []byte("")},
				Pos:        1,
			},
			char: 0,
			ok:   false,
		},
		{
			Lexer: Lexer{
				SourceCode: source.SourceCode{
					Source:         []byte("abc"),
					PositionRanges: []int{3},
				},
				Pos: 0,
			},
			char: 'a',
			ok:   true,
		},
		{
			Lexer: Lexer{
				SourceCode: source.SourceCode{
					Source:         []byte("abc"),
					PositionRanges: []int{3},
				},
				Pos: 1,
			},
			char: 'b',
			ok:   true,
		},
		{
			Lexer: Lexer{
				SourceCode: source.SourceCode{
					Source:         []byte("abc"),
					PositionRanges: []int{3},
				},
				Pos: 2,
			},
			char: 'c',
			ok:   true,
		},
		{
			Lexer: Lexer{
				SourceCode: source.SourceCode{
					Source:         []byte("abc"),
					PositionRanges: []int{3},
				},
				Pos: 3,
			},
			char: 0,
			ok:   false,
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
	tests := []struct {
		Lexer
		line string
		ok   bool
	}{
		{
			Lexer: Lexer{
				Line: 0,
			},
			line: "",
			ok:   false,
		},
		{
			Lexer: Lexer{
				SourceCode: source.SourceCode{
					Source:         []byte("a"),
					PositionRanges: []int{1},
				},
				Line: 1,
			},
			line: "a",
			ok:   true,
		},
		{
			Lexer: Lexer{
				SourceCode: source.SourceCode{
					Source:         []byte("a"),
					PositionRanges: []int{1},
				},
				Line: 1,
			},
			line: "a",
			ok:   true,
		},
		{
			Lexer: Lexer{
				SourceCode: source.SourceCode{
					Source:         []byte("a\nb"),
					PositionRanges: []int{2, 3},
				},
				Line: 2,
			},
			line: "b",
			ok:   true,
		},
		{
			Lexer: Lexer{
				SourceCode: source.SourceCode{Source: []byte("")},
				Line:       1,
			},
			line: "",
			ok:   false,
		},
		{
			Lexer: Lexer{
				SourceCode: source.SourceCode{
					Source:         []byte("abc"),
					PositionRanges: []int{3},
				},
				Line: 1,
			},
			line: "abc",
			ok:   true,
		},
		{
			Lexer: Lexer{
				SourceCode: source.SourceCode{
					Source:         []byte("abc"),
					PositionRanges: []int{3},
				},
				Line: 1,
			},
			line: "abc",
			ok:   true,
		},
		{
			Lexer: Lexer{
				SourceCode: source.SourceCode{
					Source:         []byte("abc"),
					PositionRanges: []int{3},
				},
				Line: 2,
			},
			line: "",
			ok:   false,
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
