package stack

import (
	"fmt"
	"strings"

	"github.com/petersalex27/yew/common/math"
)

// interface for usual stack operations (i.e., push, pop, peek, and number of
// elements) and an additional operation to query various states
type StackType[T any] interface {
	Push(T)
	Pop() (T, StackStatus)
	Peek() (T, StackStatus)
	GetCount() uint
	Status() StackStatus
}

// stacks with the ability to pop and push multiple elements 
type StackPlus[T any] interface {
	StackType[T]
	Clear(uint)
	MultiPush(...T)
	MultiPop(uint) ([]T, StackStatus)
}

// stacks that can create stack frames
type ReturnableStack[T any] interface {
	StackType[T]
	GetFullCount() uint
	Save()
	Return() StackStatus
}

type Stack[T any] struct {
	st    uint // stack top/capacity
	sc    uint // stack counter
	elems []T  // elements
}

// wraps fmt.Sprint(s.elems)
func (s *Stack[T]) ElemString() string { return fmt.Sprint(s.elems[:s.sc]) }

func (s *Stack[T]) ElemStringSep(encloseL, encloseR string, sep string) string {
	var b strings.Builder
	if s.GetCount() == 0 {
		return ""
	}

	b.WriteString(encloseL + fmt.Sprint(s.elems[0]) + encloseR)

	for _, elem := range s.elems[1:s.sc] {
		b.WriteString(sep)
		b.WriteString(encloseL)
		b.WriteString(fmt.Sprint(elem))
		b.WriteString(encloseR)
	}
	return b.String()
}

// stack with fixed size (i.e., stack cannot request more capacity)
type StaticStack[T any] Stack[T]

// initializes a stack w/ an initial capacity of `cap` rounded up to the 
// nearest power of two
func makeStack[T any](cap uint) (out Stack[T]) {
	if cap == 0 {
		cap = 8
	} else {
		// make capacity a power of 2
		cap = math.PowerOfTwoCeil(cap)
	}
	out.st, out.sc = cap, 0
	out.elems = make([]T, cap)
	return out
}

// Return a pointer to a new, initialized stack with room for at least `cap`
// elements.
func NewStack[T any](cap uint) *Stack[T] {
	out := new(Stack[T])
	*out = makeStack[T](cap)
	return out
}

// returns a copy of elements on stack
func (s *Stack[T]) getElemsCopy(newCap uint) []T {
	if newCap < uint(cap(s.elems)) {
		newCap = uint(cap(s.elems))
	}

	newElems := make([]T, newCap)
	for i := uint(0); i < s.sc; i++ {
		newElems[i] = s.elems[i]
	}
	return newElems
}

// grow stack size
func (s *Stack[T]) grow() StackStatus {
	if (s.st << 1) < s.st { // check for overflow
		return Overflow
	}

	newCap := s.st << 1
	// don't want to use copy here because copy will take the len of elems, but
	// len might be greater than s.sc (s.sc is the effective length of the
	// stack)
	s.elems = s.getElemsCopy(newCap)
	s.st = newCap
	return Ok
}

// true iff stack cannot have more elements pushed onto it
func (s *Stack[T]) full() bool { return s.st == s.sc }

func (s *Stack[T]) MultiPush(elems ...T) {
	for _, elem := range elems {
		s.Push(elem)
	}
}

// MultiCheck returns the top `n` elements of the stack (or an IllegalOperation 
// if there are fewer than `n` elements). Elements are returned in 
// reverse-popped order (see "IMPORTANT" below), e.g.,
//	 [w, x, y, z].MultiCheck(3) == [x, y, z]
//
// IMPORTANT: note that MultiCheck does NOT return a slice in the order 
// (left-to-right) that elemnts would be popped! In fact, it does the opposite.
//
// NOTE: elements are not removed from the stack
func (s *Stack[T]) MultiCheck(n int) (elems []T, stat StackStatus) {
	if s.sc < uint(n) {
		return nil, IllegalOperation
	}

	elems = make([]T, n)
	for i := n - 1; i >= 0; i-- {
		elems[n-1-i] = s.unsafePeekOffset(uint(i))
	}
	return elems, Ok
}

// MultiPop pops and returns the top `n` elements of the stack (or an 
// IllegalOperation if there are fewer than `n` elements). Elements are 
// returned in they order (left-to-right in return value) that they are popped,
// e.g.,
//		[w, x, y, z].MultiPop(3) == [z, y, x]
// To return the elements in the order they appear on the stack, use 
// 		(*Stack[T]) MultiCheck(int)
func (s *Stack[T]) MultiPop(n uint) (elems []T, stat StackStatus) {
	if s.sc < n {
		stat = IllegalOperation
		return
	}

	elems = make([]T, n)
	for i := uint(0); i < n; i++ {
		elems[i] = s.unsafePeekOffset(i)
	}
	s.sc = s.sc - n
	stat = Ok
	return
}

// Puts elem onto top of stack
func (s *Stack[T]) Push(elem T) {
	if s.full() {
		if s.grow().IsOverflow() {
			panic("stack overflow")
		}
	}

	s.elems[s.sc] = elem
	s.sc++
}

// removes min(n, s.GetCount()) elements from the stack
func (s *Stack[T]) Clear(n uint) {
	n = math.Min(uint(s.GetCount()), n)
	s.sc = s.sc - n
}

// returns the top element of the stack if it exists, else the Empty status
func (s *Stack[T]) Peek() (elem T, stat StackStatus) {
	if s.Empty() {
		stat = Empty
	} else {
		stat = Ok
		elem = s.elems[s.sc-1]
	}
	return
}

// unchecked peek at stack counter - 1 - n
func (s *Stack[T]) unsafePeekOffset(n uint) T {
	return s.elems[s.sc-1-n]
}

// removes the top element of the stack and returns it if stack is not empty,
// else Empty status is returned
func (s *Stack[T]) Pop() (T, StackStatus) {
	elem, stat := s.Peek()
	if stat.IsOk() {
		s.sc--
	}
	return elem, stat
}

// returns capicty of stack
func (s *Stack[T]) GetCapacity() uint {
	return s.st
}

// returns number of elements in stack
func (s *Stack[T]) GetCount() uint {
	return s.sc
}

// returns condition of stack (Ok or Empty)
func (s *Stack[T]) Status() StackStatus {
	if s.Empty() {
		return Empty
	}
	return Ok
}

// returns true iff stack has no elements
func (s *Stack[T]) Empty() bool {
	return s.GetCount() == 0
}

// searches from top to bottom for first thing predicate returns true for
func (s *Stack[T]) Search(predicate func(T) bool) (found bool, elem T) {
	for i := int64(s.sc)-1; i >= 0; i-- {
		if found = predicate(s.elems[i]); found {
			elem = s.elems[i]
			return
		}
	}
	return
}