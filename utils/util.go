package util

import "fmt"

type Stringable interface {
	ToString() string
}

type ColorString struct {

}

/*
Fmap takes a list 'ts' and a function 'f' and returns a list that has 'f' applied to each member of
'ts'.

e.g., 
ts := []int{1, 2, 3}
f := func(i int) float64 {
	return float64(i) + 1.5
}
Fmap(ts, f) // returns []float64{2.5, 3.5, 4.5}
*/
func Fmap[T any, K any](ts []T, f func(T) K) []K {
	out := make([]K, len(ts))
	for i, v := range ts {
		out[i] = f(v)
	}
	return out
}

func PrintError(actual any, expected any) {
	fmt.Printf("Actual: %v\nExpected: %v\n", actual, expected)
}

/* 
Takes a list 'ts' and applies a function 'f' to all members that are neighbors in 'ts';
it does so from left to right.
*/
func FoldLeft[T any, K any](ts []T, base K, f func(K, T) K) K {
	for _, v := range ts {
		base = f(base, v)
	}
	return base
}

/* 
Takes a list 'ts' and applies a function 'f' to all members that are neighbors in 'ts';
it does so from right to left.
*/
func FoldRight[T any, K any](ts []T, base K, f func(K, T) K) K {
	for i := len(ts) - 1; i >= 0; i-- {
		base = f(base, ts[i])
	}
	return base
}

/* 
Filter takes the predicate 'f' and applies it to each member of 'ts' (from left to right), 
returning a list of all the members that 'f' returns true for in the order they where discovered.
*/ 
func Filter[T any](ts []T, f func(T) bool) []T {
	var out []T
	for _, t := range ts {
		if f(t) {
			out = append(out, t)
		}
	}
	return out
}

func Tail[T any](xs []T) (T, bool) {
	at := len(xs) - 1
	if at < 0 {
		out := new(T)
		return *out, false
	}
	return xs[at], true
}

func Head[T []any](xs []T) []T {
	if len(xs) <= 1 {
		return xs
	}
	return xs[1:]
}