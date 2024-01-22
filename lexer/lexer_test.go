// =================================================================================================
// Alex Peters - January 21, 2024
// =================================================================================================
package lexer

import "testing"

func TestWrite(t *testing.T) {
	{
		write := func(lex *Lexer) bool {
			lex.Source = append(lex.Source, "test")
			return true
		}
		lex := Lexer{
			write: write,
			Source: []string{"test"},
		}
		expected := 1
		actual := lex.Write()
		if actual != expected {
			t.Fatalf("unexpected return value: got %d", actual)
		}
	}
	
	{
		write := func(*Lexer) bool {
			return false
		}
		lex := Lexer{
			write: write,
			Source: []string{"test"},
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
				Source: []string{"a"},
				Line: 1,
				Char: 1,
			},
			char: 'a',
			ok: true,
		},
		{
			Lexer: Lexer{
				Source: []string{"a", "b"},
				Line: 2,
				Char: 1,
			},
			char: 'b',
			ok: true,
		},
		{
			Lexer: Lexer{
				Source: []string{},
				Line: 1,
				Char: 1,
			},
			char: 0,
			ok: false,
		},
		{
			Lexer: Lexer{
				Source: []string{"abc"},
				Line: 1,
				Char: 1,
			},
			char: 'a',
			ok: true,
		},
		{
			Lexer: Lexer{
				Source: []string{"abc"},
				Line: 1,
				Char: 2,
			},
			char: 'b',
			ok: true,
		},
		{
			Lexer: Lexer{
				Source: []string{"abc"},
				Line: 1,
				Char: 3,
			},
			char: 'c',
			ok: true,
		},
		{
			Lexer: Lexer{
				Source: []string{"abc"},
				Line: 1,
				Char: 4,
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
				Char: 1,
			},
			line: "",
			ok: false,
		},
		{
			Lexer: Lexer{
				Source: []string{"a"},
				Line: 1,
				Char: 0,
			},
			line: "a",
			ok: true,
		},
		{
			Lexer: Lexer{
				Source: []string{"a"},
				Line: 1,
				Char: 1,
			},
			line: "a",
			ok: true,
		},
		{
			Lexer: Lexer{
				Source: []string{"a", "b"},
				Line: 2,
				Char: 1,
			},
			line: "b",
			ok: true,
		},
		{
			Lexer: Lexer{
				Source: []string{},
				Line: 1,
				Char: 1,
			},
			line: "",
			ok: false,
		},
		{
			Lexer: Lexer{
				Source: []string{"abc"},
				Line: 1,
				Char: 3,
			},
			line: "abc",
			ok: true,
		},
		{
			Lexer: Lexer{
				Source: []string{"abc"},
				Line: 1,
				Char: 4,
			},
			line: "abc",
			ok: true,
		},
		{
			Lexer: Lexer{
				Source: []string{"abc"},
				Line: 2,
				Char: 1,
			},
			line: "",
			ok: false,
		},
	}

	for _, test := range tests {
		line, ok := test.Lexer.currentSourceLine()
		if test.line != line {
			t.Fatalf("unexpected string (Lexer=%v, line=\"%s\", ok=%t): got %s", (&test.Lexer), test.line, test.ok, line)
		}
		if test.ok != ok {
			t.Fatalf("unexpected ok value: got %t", ok)
		}
	}
}