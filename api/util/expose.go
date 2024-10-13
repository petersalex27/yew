package util

import "strings"

func Expose[T any](expose func(T) string, items ...T) string {
	res := ExposeList(expose, items, " ")
	// ExposeList guarantees that the result will be a string with a left-most '[' and a right-most ']',
	// so the brackets can safely be removed
	return res[1:len(res)-1]
}

// ExposeList returns a string representation of a list of items
//
// this function is guaranteed to return a string with a left-most '[' and a right-most ']' and a 
// possibly empty string between them
func ExposeList[T any](expose func(T) string, items []T, sep string) string {
	if len(items) == 0 {
		return "[]"
	}

	if len(items) == 1 {
		return "[" + expose(items[0]) + "]"
	}

	b := &strings.Builder{}
	b.WriteString(expose(items[0]))
	for _, item := range items[1:] {
		b.WriteString(sep)
		b.WriteString(expose(item))
	}
	return "[" + b.String() + "]"
}