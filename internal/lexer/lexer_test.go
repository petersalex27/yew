// =================================================================================================
// Alex Peters - January 21, 2024
// =================================================================================================
package lexer

import (
	"testing"

	"github.com/petersalex27/yew/api"
	"github.com/petersalex27/yew/api/util"
)

var _ api.Scanner = (*Lexer)(nil)

func initWithPos(src api.Source, pos int) *Lexer {
	lex := Init(src)
	lex.Pos = pos
	return lex
}

func initWithLine(src api.Source, line int) *Lexer {
	lex := Init(src)
	lex.Line = line
	return lex
}

func TestCurrentSourceChar(t *testing.T) {
	tests := []struct {
		name string
		*Lexer
		char byte
		ok   bool
	}{
		{
			name:  "`a` @ 0",
			Lexer: initWithPos(util.StringSource("a"), 0),
			char:  'a',
			ok:    true,
		},
		{
			name:  "`a\\nb` @ 2",
			Lexer: initWithPos(util.StringSource("a\nb"), 2),
			char:  'b',
			ok:    true,
		},
		{
			name:  "empty source @ 0",
			Lexer: initWithPos(util.StringSource(""), 0),
			char:  0,
			ok:    false,
		},
		{
			name:  "empty source @ 1",
			Lexer: initWithPos(util.StringSource(""), 1),
			char:  0,
			ok:    false,
		},
		{
			name:  "`abc` @ 0",
			Lexer: initWithPos(util.StringSource("abc"), 0),
			char:  'a',
			ok:    true,
		},
		{
			name:  "`abc` @ 1",
			Lexer: initWithPos(util.StringSource("abc"), 1),
			char:  'b',
			ok:    true,
		},
		{
			name:  "`abc` @ 2",
			Lexer: initWithPos(util.StringSource("abc"), 2),
			char:  'c',
			ok:    true,
		},
		{
			name:  "`abc` @ 3",
			Lexer: initWithPos(util.StringSource("abc"), 3),
			char:  0,
			ok:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			char, ok := test.Lexer.currentSourceChar()
			if test.char != char {
				t.Fatalf("unexpected character: got %c", char)
			}
			if test.ok != ok {
				t.Fatalf("unexpected ok value: got %t", ok)
			}
		})
	}
}

func TestCurrentSourceLine(t *testing.T) {
	tests := []struct {
		name string
		*Lexer
		line string
		ok   bool
	}{
		{
			name:  "empty source",
			Lexer: initWithLine(util.EmptySource(), 0),
			line:  "",
			ok:    false,
		},
		{
			// this would tokenize as just <EOF>
			name:  "empty source @ line 1",
			Lexer: initWithLine(util.EmptySource(), 1),
			line:  "",
			ok:    true,
		},
		{
			name:  "source with single line",
			Lexer: initWithLine(util.StringSource("a"), 1),
			line:  "a",
			ok:    true,
		},
		{
			name:  "source with two lines @ line 2",
			Lexer: initWithLine(util.StringSource("a\nb"), 2),
			line:  "b",
			ok:    true,
		},
		{
			name:  "source with one line multiple characters @ line 1",
			Lexer: initWithLine(util.StringSource("abc"), 1),
			line:  "abc",
			ok:    true,
		},
		{
			name:  "source with one line multiple characters @ line 2",
			Lexer: initWithLine(util.StringSource("abc"), 2),
			line:  "",
			ok:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			line, ok := test.Lexer.currentSourceLine()
			if test.ok != ok {
				t.Fatalf("unexpected ok value: got %t", ok)
			}
			if test.line != line {
				t.Fatalf("unexpected string (Lexer=%v, line=\"%s\", ok=%t): got \"%s\"", test.Lexer, test.line, test.ok, line)
			}
		})
	}
}
