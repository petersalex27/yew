// =================================================================================================
// Alex Peters - February 16, 2024
// =================================================================================================
package common

import "fmt"

type Set[T any] struct {
	internal map[string]T
}

func Create[T any](elems ...T) Set[T] {
	set := Set[T]{internal: make(map[string]T, len(elems))}
	for _, elem := range elems {
		set.MutAdd(elem)
	}
	return set
}

func str[T any](elem T) string {
	return fmt.Sprint(elem)
}

func (s Set[T]) Contains(elem T) bool {
	_, exists := s.internal[str(elem)]
	return exists
}

// Takes the set difference of S and R. S is not changed
//
//	S \ R = {x : x in S & x !in R}
//
// NOTE: this considers two elements to be the same when fmt.Sprintf returns the same string for
// both
func Difference[T, U any](S Set[T], R Set[U]) Set[T] {
	result := Set[T]{internal: make(map[string]T, len(S.internal))}
	// copy result <- S
	for k, v := range S.internal {
		result.internal[k] = v
	}

	return MutDifference(result, R)
}

// Takes the set difference of S and R. S **IS** changed
//
//	S \ R = {x : x in S /\ x !in R}
//
// NOTE: this considers two elements to be the same when fmt.Sprintf returns the same string for
// both
func MutDifference[T, U any](S Set[T], R Set[U]) Set[T] {
	// take set difference
	for k, _ := range R.internal {
		if _, found := S.internal[k]; found {
			delete(S.internal, k)
		}
	}
	return S
}

// Takes set intersection of S and R. S is **NOT** changed
//
//	S U R = {x : x in S \/ x in R}
//
// NOTE: if S and R "share" an element, the element of S is preferred
func Intersection[T, U any](S Set[T], R Set[U]) Set[T] {
	result := Set[T]{internal: make(map[string]T)}
	for k, v := range S.internal {
		if _, found := R.internal[k]; found {
			result.internal[k] = v
		}
	}

	return result
}

// Takes the set union of S and {s}. S is not changed
//
//	S U {s} = {x : x in S \/ x == s}
//
// NOTE: this considers two elements to be the same when fmt.Sprintf returns the same string for
// both
//
// NOTE: if s already exists in S, uses the existing s instead of argument `s`
func (S Set[T]) Add(s T) Set[T] {
	key := str(s)
	if _, found := S.internal[key]; found {
		return S
	}


	result := Set[T]{internal: make(map[string]T, len(S.internal)+1)}
	// copy result <- S
	for k, v := range S.internal {
		result.internal[k] = v
	}
	result.internal[key] = s
	return result
}

// Takes the set union of S and {s}. S **IS** changed
//
//	S U {s} = {x : x in S \/ x == s}
//
// NOTE: this considers two elements to be the same when fmt.Sprintf returns the same string for
// both
//
// NOTE: if s already exists in S, uses the existing s instead of argument `s`
func (S Set[T]) MutAdd(s T) {
	key := str(s)
	if _, found := S.internal[key]; found {
		return
	}
	S.internal[key] = s
}

// Takes the set difference of S and {s}. S is not changed
//
//	S \ {s} = {x : x in S & x != s}
//
// NOTE: this considers two elements to be the same when fmt.Sprintf returns the same string for
// both
func (S Set[T]) Remove(s T) Set[T] {
	key := str(s)
	if _, found := S.internal[key]; !found {
		return S
	}


	result := Set[T]{internal: make(map[string]T, len(S.internal)-1)}
	// copy result <- S
	for k, v := range S.internal {
		if k == key {
			continue
		}
		result.internal[k] = v
	}

	return result
}

// Takes the set difference of S and {s}. S **IS** changed
//
//	S \ {s} = {x : x in S & x != s}
//
// NOTE: this considers two elements to be the same when fmt.Sprintf returns the same string for
// both
func (S Set[T]) MutRemove(s T) {
	delete(S.internal, str(s))
}

// Takes set union of S and R. S is **NOT** changed
//
//	S U R = {x : x in S \/ x in R}
//
// NOTE: if S and R "share" an element, the element of S is preferred
//
// NOTE: panics if U cannot be cast to T
func Union[T, U any](S Set[T], R Set[U]) Set[T] {
	result := Set[T]{internal: make(map[string]T, len(S.internal))}
	// copy result <- S
	for k, v := range S.internal {
		result.internal[k] = v
	}

	return MutUnion(result, R)
}

// Takes set union of S and R. S **IS** changed
//
//	S U R = {x : x in S \/ x in R}
//
// NOTE: if S and R "share" an element, the element of S is preferred
//
// NOTE: panics if U cannot be cast to T
func MutUnion[T, U any](S Set[T], R Set[U]) Set[T] {
	for k, v := range R.internal {
		if _, found := S.internal[k]; !found {
			S.internal[k] = any(v).(T)
		}
	}
	return S
}

