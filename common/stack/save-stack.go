package stack

import (
	"strings"
)

type SaveStack[a any] struct {
	// basically, a stack of stacks
	frames Stack[*Stack[a]]
}

func NewSS[a any]() *SaveStack[a] {
	stack := Stack[*Stack[a]]{}
	stack.grow()
	stack.Push(New[a]())
	return &SaveStack[a]{frames: stack}
}

// returns the stack's current frame data in a new slice and moves down to the next frame
//
// the order of elements in the returned slice is the same as the order in which they were pushed
// with index 0 being the least recently pushed element
//
// See: FlushAll for returning all elements on the stack from all frames
// See: Return for just discarding the current frame and moving down to the next frame
func (s *SaveStack[a]) Flush() []a {
	if frame, ok := s.frames.Pop(); ok {
		out := frame.Flush()
		s.Return() // discard the frame
		return out
	}
	return []a{}
}

// prints the stack's entire contents as follows:
//
//	frame.String() + "\n"
//	+ "-----------------" + "\n"
//	+ frame2.String() + "\n"
//	+ "-----------------" + "\n"
//	+ ...
//	+ "-----------------" + "\n"
//	+ finalFrame.String()
//
// empty frames are represented as "_", as is the stack itself if it is empty
//
// this is to distinguish between items that become the empty string which would appear as "[]" in a
// frame--which would be ambiguous if empty frames were also represented as "[]"
//
// the line "---" + ... + "---" + "\n" has the same length as the longest frame string
func (s *SaveStack[a]) String() string {
	if s.frames.ctr == 0 {
		return "_"
	}

	maxLen := 0
	strs := make([]string, 0, s.frames.ctr)
	for i := s.frames.ctr - 1; i >= 0; i-- {
		str, n := s.frames.data[i].stringStack()
		maxLen = max(maxLen, n)
		strs = append(strs, str)
	}
	return strings.Join(strs, strings.Repeat("-", maxLen)+"\n")
}

// returns all elements on the stack from all frames
func (s *SaveStack[a]) FlushAll() [][]a {
	frames := make([][]a, s.CountFrames())
	for i := 0; s.CountFrames() > 0; i++ {
		frames[i] = s.Flush()
	}
	return frames
}

// discard the current frame and move down to the next frame
func (s *SaveStack[a]) Return() {
	if frame, ok := s.frames.Pop(); ok {
		frame.data = nil // let the GC do its thing
	}
}

// length of the current frame
//
// See: CountFrames for the number of frames on the stack
func (s *SaveStack[a]) Len() int {
	if frame, ok := s.frames.Peek(); ok {
		return frame.Len()
	}
	return 0
}

// number of frames on the stack
func (s *SaveStack[a]) CountFrames() int {
	return s.frames.Len()
}

// push a new element onto the current frame
func (s *SaveStack[a]) Push(e a) {
	frame, ok := s.frames.Peek()
	if !ok {
		// create new frames
		s.frames.grow()
		frame, ok = s.frames.Peek()
		if !ok {
			panic("stack: failed to create new frame")
		}
	}

	// push onto the current frame
	frame.Push(e) // this will update the reference in the stack too
}

// returns the top element of the current frame without removing it
func (s *SaveStack[a]) Peek() (top a, ok bool) {
	var frame *Stack[a]
	frame, ok = s.frames.Peek()
	if ok {
		top, ok = frame.Peek()
	} // else no frame to peek into
	return top, ok
}

// returns the top element of the current frame and removes it
func (s *SaveStack[a]) Pop() (top a, ok bool) {
	var frame *Stack[a]
	frame, ok = s.frames.Peek()
	if ok {
		top, ok = frame.Pop() // this will update the reference in the stack too
	} // else no frame to pop from
	return top, ok
}
