package stack

import (
	"testing"
)

func runPush(bs []byte) *Stack[byte] {
	s := New[byte]()
	for _, b := range bs {
		s.Push(b)
	}
	return s
}

func TestFlush(t *testing.T) {
	tests := [][]byte{
		{},
		{'a'},
		{'a', 'b', 'c'},
	}

	for _, test := range tests {
		s := runPush(test)
		if got := s.Flush(); len(got) != len(test) {
			t.Errorf("Expected %v, got %v", test, got)
		} else if string(got) != string(test) {
			t.Errorf("Expected %v, got %v", test, got)
		}
	}
}

func TestPeek(t *testing.T) {
	tests := [][]byte{
		{},
		{'a'},
		{'a', 'b', 'c'},
	}

	for _, test := range tests {
		s := runPush(test)

		if c, ok := s.Peek(); ok {
			if c != test[len(test)-1] {
				t.Errorf("Expected %c, got %c", test[len(test)-1], c)
			}
		} else if len(test) > 0 {
			t.Errorf("Expected (%c, true), got (%c, false)", test[len(test)-1], c)
		}
	}
}

func TestPush(t *testing.T) {
	tests := [][]byte{
		{'a'},
		{'a', 'b', 'c'},
	}

	for _, test := range tests {
		s := runPush(test)

		if s.ctr != len(test) {
			t.Errorf("Expected s.ctr=%d, got %d", len(test), len(s.data))
		}

		for i, b := range test {
			if s.data[i] != b {
				t.Errorf("Expected s.data[%d]=%c, got %c", i, b, s.data[i])
			}
		}
	}
}

func TestPop(t *testing.T) {
	tests := [][]byte{
		{},
		{'a'},
		{'a', 'b', 'c'},
	}

	for _, test := range tests {
		s := runPush(test)
		for i := len(test) - 1; i >= 0; i-- {
			if got, ok := s.Pop(); !ok || got != test[i] {
				t.Errorf("Expected (%c, %t), got (%c, %t)", test[i], true, got, ok)
			}
		}

		if c, ok := s.Pop(); ok {
			t.Errorf("Expected (0, false), got (%c, %t)", c, ok)
		}
	}
}
