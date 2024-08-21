package stack

import (
	"testing"
)

func TestPush_ret(t *testing.T) {
	tests := []struct {
		expect string
	}{
		{""},
		{"a"},
		{"hello, world"},
		{"this is a string with length greater than 32"},
	}

	for _, test := range tests {
		stack := NewSaveStack[byte](32)
		// push entire expected string
		for _, b := range test.expect {
			stack.Push(byte(b))
		}

		// check stack counter
		expectedCounter := uint(len(test.expect))
		if stack.sc != expectedCounter {
			t.Fatalf("unexpected stack counter (%v): got %v", expectedCounter, stack.sc)
		}
		// check stack elements
		actual := string(stack.elems[:stack.sc])
		if test.expect != actual {
			t.Fatalf("unexpected stack element (\"%s\"): got \"%s\"", test.expect, actual)
		}
	}
}

func TestPop_empty_ret(t *testing.T) {
	stack := NewSaveStack[byte](8)
	_, stat := stack.Pop()
	if !stat.IsEmpty() {
		t.Fatalf("expected stack to return empty status: got %v", stat)
	}
}

func TestPop_ret(t *testing.T) {
	tests := []struct {
		put []byte
	}{
		{[]byte("a")},
		{[]byte("hello, world")},
		{[]byte("this is a string with length greater than 32")},
	}

	for _, test := range tests {
		stack := NewSaveStack[byte](uint(len(test.put)))
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

func TestPeek_ret(t *testing.T) {
	tests := []struct {
		put  []byte
		stat StackStatus
	}{
		{[]byte(""), Empty},
		{[]byte("a"), Ok},
		{[]byte("hello, world"), Ok},
		{[]byte("this is a string with length greater than 32"), Ok},
	}

	for _, test := range tests {
		stack := NewSaveStack[byte](uint(len(test.put)))
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

func TestGetCount_ret(t *testing.T) {
	tests := []struct {
		sc uint
	}{
		{0},
		{1},
		{10000},
	}

	for _, test := range tests {
		stack := NewSaveStack[byte](8)
		stack.sc = test.sc
		actual := stack.GetCount()
		if actual != test.sc {
			t.Fatalf("unexpected stack counter (%v): got %v", test.sc, actual)
		}
	}
}

func TestStatus_ret(t *testing.T) {
	tests := []struct {
		sc   uint
		stat StackStatus
	}{
		{0, Empty},
		{1, Ok},
		{10000, Ok},
	}

	for _, test := range tests {
		stack := NewSaveStack[byte](8)
		stack.sc = test.sc
		actual := stack.Status()
		if actual != test.stat {
			t.Fatalf("unexpected status (%v): got %v", test.stat, actual)
		}
	}
}

func TestSave(t *testing.T) {
	tests := []struct {
		beforeSave string
		afterSave  string
		peek       rune
		stat       StackStatus
		total      string
	}{
		{"", "", 0, Empty, ""},
		{"", "a", 'a', Ok, "a"},
		{"a", "", 0, Empty, "a"},
		{"123", "456", '6', Ok, "123456"},
	}

	for _, test := range tests {
		stack := NewSaveStack[byte](32, 4)

		// push
		for _, b := range []byte(test.beforeSave) {
			stack.Push(b)
		}

		// now save and check
		stack.Save()
		expectBC := uint(len(test.beforeSave))
		actualBC := stack.bc
		if expectBC != actualBC {
			t.Fatalf("unexpected base counter (%v): got %v", expectBC, actualBC)
		}

		// push after save
		for _, b := range []byte(test.afterSave) {
			stack.Push(b)
		}
		// check that bc hasn't changed
		if expectBC != stack.bc {
			t.Fatalf("base counter changed (%v): changed to %v", expectBC, stack.bc)
		}
		actual, stat := stack.Peek()
		if !stat.Is(test.stat) {
			t.Fatalf("unexpected peek status (%v): got %v", test.stat, actual)
		}

		if stat.IsOk() {
			expect := test.afterSave[len(test.afterSave)-1]
			if actual != expect {
				t.Fatalf("unexpected peek result (%v): got %v", expect, actual)
			}
		}

		actualTotal := string(stack.elems[:stack.sc])
		if actualTotal != test.total {
			t.Fatalf("unexpected total result (%v): got %v", test.total, actualTotal)
		}
	}
}

func TestReturn(t *testing.T) {
	tests := []struct {
		afterReturn  string
		returnStatus StackStatus
		peek         rune
		stat         StackStatus
		total        string
	}{
		{"", Ok, 0, Empty, ""},
		{"", Ok, 0, Empty, "a"},
		{"a", Ok, 'a', Ok, "a"},
		{"123", Ok, '3', Ok, "123456"},
	}

	for _, test := range tests {
		stack := NewSaveStack[byte](32, 4)
		// place all elems
		for i := range test.total {
			stack.elems[i] = test.total[i]
		}

		// set counter
		stack.sc = uint(len(test.total))
		// set return point
		stack.returnStack.Push(0)
		// set base counter
		ret := uint(len(test.afterReturn))
		stack.bc = ret

		// now return
		retStat := stack.Return()
		if !retStat.Is(test.returnStatus) {
			t.Fatalf("unexpected return status (%v): got %v", test.returnStatus, retStat)
		}

		expectBC := uint(0)
		actualBC := stack.bc
		if expectBC != actualBC {
			t.Fatalf("unexpected base ")
		}

		actual, stat := stack.Peek()
		if !stat.Is(test.stat) {
			t.Fatalf("unexpected peek status (%v): got %v", test.stat, stat)
		}

		if stat.IsOk() {
			expect := byte(test.peek)
			if actual != expect {
				t.Fatalf("unexpected element peeked (%v): got %v", expect, actual)
			}
		}

		expectSC, actualSC := uint(len(test.afterReturn)), stack.sc
		if expectSC != actualSC {
			t.Fatalf("unexpected stack counter post-return (%v): got %v", expectSC, actualBC)
		}

		actualResult := string(stack.elems[:stack.sc])
		if actualResult != test.afterReturn {
			t.Fatalf("unexpected post-return stack (%v): got %v", test.afterReturn, actualResult)
		}
	}
}

func TestRebase(t *testing.T) {
	tests := []struct {
		savedFrames   []string
		rebaseStatus  StackStatus
		peek          rune
		stat          StackStatus
		expectedFrame string
	}{
		{[]string{""}, Ok, 0, Empty, ""},
		{[]string{""}, Ok, 'a', Ok, "a"},
		{[]string{"a"}, Ok, 'a', Ok, "a"},
		{[]string{"123"}, Ok, '6', Ok, "123456"},
		{[]string{"123", "45"}, Ok, '6', Ok, "123456"},
	}

	for _, test := range tests {
		stack := NewSaveStack[byte](32, 4)
		// place all elems
		for i := range test.expectedFrame {
			stack.elems[i] = test.expectedFrame[i]
		}

		// set counter
		stack.sc = uint(len(test.expectedFrame))

		// for each frame, set return point
		previousTop := uint(0)
		for _, frame := range test.savedFrames {
			// save return point
			stack.returnStack.Push(previousTop)
			previousTop = uint(len(frame)) // stack counter for frame
		}

		// set base counter for current frame (base counter is previous frame's 
		// stack counter)
		ret := previousTop
		stack.bc = ret

		// test rebase for each saved frame
		numFrames := len(test.savedFrames)
		for j, frameIndex := 0, numFrames - 1; j < len(test.savedFrames); j, frameIndex = j + 1, frameIndex - 1 {
			// now rebase
			retStat := stack.Rebase()
			if !retStat.Is(test.rebaseStatus) {
				t.Fatalf("rebase status (%v): got %v", test.rebaseStatus, retStat)
			}

			// new base counter should be frame saved before frame's stack counter 
			// (else 0)
			var expectBC uint
			if frameIndex == 0 {
				expectBC = 0
			} else {
				expectBC = uint(len(test.savedFrames[frameIndex-1]))
			}
			actualBC := stack.bc
			if expectBC != actualBC {
				t.Fatalf("base counter (%v): got %v", expectBC, actualBC)
			}

			actual, stat := stack.Peek() // stack counter should never change, so top should never change
			if !stat.Is(test.stat) {
				t.Fatalf("peek status (%v): %v", test.stat, stat)
			}

			if stat.IsOk() {
				expect := byte(test.peek)
				if actual != expect {
					t.Fatalf("element peeked (%v): got %v", expect, actual)
				}
			}

			// stack counter should never change
			expectSC, actualSC := uint(len(test.expectedFrame)), stack.sc
			if expectSC != actualSC {
				t.Fatalf("stack counter after rebase (%v): got %v", expectSC, actualSC)
			}

			actualResult := string(stack.elems[stack.bc:stack.sc])
			expectedResult := test.expectedFrame[expectBC:]
			if actualResult != expectedResult {
				t.Fatalf("after rebase result (%v): got %v", expectedResult, actualResult)
			}
		}
	}
}
