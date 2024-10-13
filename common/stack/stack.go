package stack

import (
	"fmt"
	"math"
	"strings"
)

type Stack[a any] struct {
	cap, ctr int
	data     []a
}

// elemCopy is optional--if more than one argument is provided, only the first argument is used
func (s *Stack[a]) Copy(elemCopy ...func(a) a) *Stack[a] {
	var cpf func(a) a
	if len(elemCopy) > 0 {
		cpf = elemCopy[0]
	} else {
		cpf = func(elem a) a { return elem } // id function
	}

	out := New[a]()
	elem, ok := s.Pop()
	for ; ok; elem, ok = s.Pop() {
		out.Push(cpf(elem))
	}
	return out
}

func New[a any]() *Stack[a] {
	stack := &Stack[a]{}
	stack.grow()
	return stack
}

func addString(b *strings.Builder, str string) int {
	if !strings.ContainsRune(str, '\n') {
		return len(str)
	}

	n := 0
	strs := strings.Split(str, "\n")
	if len(strs) > 0 {
		n, _ = b.WriteString(strs[0])
	}
	for i := 1; i < len(strs); i++ {
		b.WriteString("\n" + strs[i])
		n = max(n, len(strs[i]))
	}
	return n
}

// returns a string representation of the stack's data along with the length of the
// longest line in the string (not including the newline character)
func (s *Stack[a]) stringStack() (string, int) {
	if s.ctr == 0 {
		return "_", 1
	}

	b := strings.Builder{}

	var str string
	var n int = 0
	for i := s.ctr - 1; i >= 1; i-- {
		str = fmt.Sprintf("[%v]", s.data[i])
		n = max(n, addString(&b, str))
		b.WriteByte('\n')
	}

	str = fmt.Sprintf("[%v]", s.data[0])
	n = max(n, addString(&b, str))
	return b.String(), n
}

// return "_" if the stack is empty, otherwise return the stack's data in the format:
//
//	"[" + top as string + "]" + "\n"
//	"[" + next as string + "]" + "\n"
//	...
//	"[" + bottom as string + "]"
func (s *Stack[a]) String() string {
	str, _ := s.stringStack()
	return str
}

// returns the stack's data in a new slice and resets the stack
//
// the order of elements in the returned slice is the same as the order in which they were pushed
// with index 0 being the least recently pushed element
func (s *Stack[a]) Flush() []a {
	if s.ctr == 0 {
		return []a{}
	}

	out := make([]a, s.ctr)
	copy(out, s.data[:s.ctr])
	s.ctr = 0
	return out
}

func (s *Stack[a]) Len() int { return s.ctr }

func (s *Stack[a]) Push(e a) {
	if s.cap == s.ctr {
		s.grow()
	}
	s.data[s.ctr] = e
	s.ctr++
}

func (s *Stack[a]) Peek() (top a, ok bool) {
	if ok = s.ctr > 0; ok {
		top = s.data[s.ctr-1]
	}
	return top, ok
}

func (s *Stack[a]) Pop() (top a, ok bool) {
	if top, ok = s.Peek(); ok {
		s.ctr--
	}
	return top, ok
}

const (
	smallCap int = 8
)

func (s *Stack[a]) grow() {
	if s.cap == 0 {
		s.cap = smallCap
	} else if s.cap < math.MaxInt {
		s.cap *= 2 // this works b/c starting cap is factor of 2
	} else {
		// that's a lotta data ...
		panic("stack: cannot grow beyond max capacity")
	}

	old := s.data
	s.data = make([]a, s.cap)
	if old != nil {
		copy(s.data, old)
	}
}
