package stack

import (
	"testing"
)

func TestPush(t *testing.T) {
	tests := []struct{
		expect string
	}{
		{""},
		{"a"},
		{"hello, world"},
		{"this is a string with length greater than 32"},
	}

	for _, test := range tests {
		stack := NewStack[byte](32)
		// push entire expected string
		for _, b := range test.expect {
			stack.Push(byte(b))
		}
		
		// check stack counter
		expectedCounter := uint(len(test.expect))
		if stack.sc != expectedCounter {
			t.Fatalf("unexpected stack counter (%d): got %d", expectedCounter, stack.sc)
		}
		// check stack elements
		actual := string(stack.elems[:stack.sc])
		if test.expect != actual {
			t.Fatalf("unexpected stack element (\"%s\"): got \"%s\"", test.expect, actual)
		}
	}
}

func TestPop_empty(t *testing.T) {
	stack := NewStack[byte](8)
	_, stat := stack.Pop()
	if !stat.IsEmpty() {
		t.Fatalf("expected stack to return empty status: got %v", stat)
	}
}

func TestPop(t *testing.T) {
	tests := []struct{
		put []byte
	}{
		{[]byte("a")},
		{[]byte("hello, world")},
		{[]byte("this is a string with length greater than 32")},
	}

	for _, test := range tests {
		stack := NewStack[byte](uint(len(test.put)))
		// push entire expected string
		for i, b := range test.put {
			stack.elems[i] = byte(b)
		}
		stack.sc = uint(len(test.put))

		for j := range test.put {
			actual, stat := stack.Pop()
			if stat.NotOk() {
				t.Fatalf("expected pop status to be Ok: got %v", stat)
			}
			expect := test.put[len(test.put)-1-j]
			if actual != expect {
				t.Fatalf("unexpected element popped (%v): got %v", expect, actual)
			}
		}
		
		// check stack counter
		expectedCounter := uint(0)
		if stack.sc != expectedCounter {
			t.Fatalf("unexpected stack counter (%v): got %v", expectedCounter, stack.sc)
		}
	}
}

func TestPeek(t *testing.T) {
	tests := []struct{
		put []byte
		stat StackStatus
	}{
		{[]byte(""), Empty},
		{[]byte("a"), Ok},
		{[]byte("hello, world"), Ok},
		{[]byte("this is a string with length greater than 32"), Ok},
	}

	for _, test := range tests {
		stack := NewStack[byte](uint(len(test.put)))
		// push entire expected string
		for i, b := range test.put {
			stack.elems[i] = byte(b)
		}
		stack.sc = uint(len(test.put))

		actual, stat := stack.Peek()
		if !stat.Is(test.stat) {
			t.Fatalf("unexpected peek status (%v): got %v", test.stat, stat)
		}
		
		if len(test.put) > 0 {
			expect := test.put[len(test.put)-1]
			if actual != expect {
				t.Fatalf("unexpected element peeked (%v): got %v", expect, actual)
			}
		}
		
		// check stack counter
		expectedCounter := uint(len(test.put))
		if stack.sc != expectedCounter {
			t.Fatalf("unexpected stack counter (%v): got %v", expectedCounter, stack.sc)
		}
	}
}

func TestGetCount(t *testing.T) {
	tests := []struct{
		sc uint
	}{
		{0},
		{1},
		{10000},
	}

	for _, test := range tests {
		stack := NewStack[byte](8)
		stack.sc = test.sc
		actual := stack.GetCount()
		if actual != test.sc {
			t.Fatalf("unexpected stack counter (%v): got %v", test.sc, actual)
		}
	}
}

func TestStatus(t *testing.T) {
	tests := []struct{
		sc uint
		stat StackStatus
	}{
		{0, Empty},
		{1, Ok},
		{10000, Ok},
	}

	for _, test := range tests {
		stack := NewStack[byte](8)
		stack.sc = test.sc
		actual := stack.Status()
		if actual != test.stat {
			t.Fatalf("unexpected status (%v): got %v", test.stat, actual)
		}
	}
}