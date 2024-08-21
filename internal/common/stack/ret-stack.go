package stack

import (
	"fmt"
	"strings"
)

type SaveStack[T any] struct {
	Stack[T]
	bc          uint // base counter
	returnStack Stack[uint]
}

// only up to 2 arguments are used, one argument required; cap is the starting
// capacity for the stack. returnCap[0] is the starting capacity for the
// number of returns that can be saved. returnCap has a default value of 8
func NewSaveStack[T any](cap uint, returnCap ...uint) *SaveStack[T] {
	out := new(SaveStack[T])
	out.Stack = makeStack[T](cap)
	out.bc = 0
	rCap := uint(8)
	if len(returnCap) > 0 {
		rCap = returnCap[0]
	}
	out.returnStack = makeStack[uint](rCap)
	return out
}

func (stack *SaveStack[T]) Push(elem T) {
	stack.Stack.Push(elem)
}

func elemStringSepHelper[T any](elems []T, encloseL, encloseR, sep string) string {
	var b strings.Builder

	if len(elems) == 0 {
		return ""
	}

	b.WriteString(encloseL + fmt.Sprint(elems[0]) + encloseR)

	for _, elem := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(encloseL)
		if st, ok := any(elem).(fmt.Stringer); ok && st != nil {
			b.WriteString(st.String())
		} else {
			b.WriteString("_?_")
		}
		b.WriteString(encloseR)
	}
	return b.String()
}

func (s *SaveStack[T]) ElemStringSep(encloseL, encloseR, frameSep, sep string) string {
	if s.GetFullCount() == 0 {
		return ""
	}

	var b strings.Builder
	bcCount := s.returnStack.GetCount()
	bcs := s.returnStack.elems[:bcCount]
	if once := bcCount == 0; once {
		return elemStringSepHelper(s.elems[:s.sc], encloseL, encloseR, sep)
	}

	var prev, next uint = 0, bcs[0]
	for i := 0; i < len(bcs); i++ {
		prev = next
		if len(bcs) == i+1 {
			next = s.bc
		} else {
			next = bcs[i+1]
		}

		res := elemStringSepHelper(s.elems[prev:next], encloseL, encloseR, sep)
		b.WriteString(res)
		b.WriteString(frameSep)
	}

	res := elemStringSepHelper(s.elems[next:s.sc], encloseL, encloseR, sep)
	b.WriteString(res)

	return b.String()
}

func (stack *SaveStack[T]) Pop() (elem T, stat StackStatus) {
	if stack.bc == stack.sc {
		stat = Empty
		return
	}
	return stack.Stack.Pop()
}

func (stack *SaveStack[T]) Peek() (elem T, stat StackStatus) {
	if stack.bc == stack.sc {
		stat = Empty
		return
	}
	return stack.Stack.Peek()
}

// number of stack frames
func (stack *SaveStack[T]) GetFrames() uint {
	return stack.returnStack.GetCount() + 1
}

// number of elements in current frame
func (stack *SaveStack[T]) GetCount() uint {
	return stack.GetFullCount() - stack.bc
}

// true iff current frame is empty
func (stack *SaveStack[T]) Empty() bool {
	return stack.GetCount() == 0
}

// true iff all frames are empty
func (stack *SaveStack[T]) FullEmpty() bool {
	return stack.GetFullCount() == 0
}

// number of elems in all frames
func (stack *SaveStack[T]) GetFullCount() uint {
	return stack.Stack.GetCount()
}

func (stack *SaveStack[T]) Status() StackStatus {
	return stack.Stack.Status()
}

// create return point point to return to later
func (stack *SaveStack[T]) Save() {
	save := stack.bc
	stack.returnStack.Push(save)
	stack.bc = stack.sc
}

// Returns base counter to saved base counter but does NOT change the stack
// counter. This function effectively "merges" the top two stack frames into
// a single frame
func (stack *SaveStack[T]) Rebase() (stat StackStatus) {
	if stack.returnStack.GetCount() == 0 { // return stack is empty
		return IllegalReturn
	}
	var previousBaseCounter uint // base counter of saved frame
	previousBaseCounter, stat = stack.returnStack.Pop()
	if stat.IsOk() {
		stack.bc = previousBaseCounter // merge top two frames
	}
	// by way of not setting the stack counter, the previous stack frame
	// incorporates all the elements of the frame above it
	return
}

// Return base counter to saved base counter and return stack counter to 
// previous stack counter. Return fails when nothing is saved.
func (stack *SaveStack[T]) Return() (stat StackStatus) {
	old := stack.bc // where previous stack frame ended
	stat = stack.Rebase() // merge top two frames
	if stat.IsOk() {
		stack.sc = old // remove top (now merged with frame under it) frame
	}
	return
}

// search from top to bottom for first occurrence of element that predicate occurs for
func (s *SaveStack[T]) Search(predicate func(T) bool) (found bool, elem T) {
	for i := int64(s.sc) - 1; i >= int64(s.bc); i-- {
		if found = predicate(s.elems[i]); found {
			elem = s.elems[i]
			return
		}
	}
	return
}
