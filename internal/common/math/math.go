package math

import "cmp"

func Min[T cmp.Ordered](a, b T) T {
	if a > b {
		return b
	}
	return a
}

// returns max {a, b} or a if a==b
func Max[T cmp.Ordered](a, b T) T {
	if a < b {
		return b
	}
	return a
}

type WholeNumber interface {
	~int | ~uint
}

// returns the absolute value of some integer type
func Abs[N WholeNumber](a N) N {
	if a > 0 {
		return a
	}
	return -a
}

// rounds up to the nearest power of two
func PowerOfTwoCeil[N WholeNumber](n N) N {
	n = n - 1
	n = n | (n >> 1)
	n = n | (n >> 2)
	n = n | (n >> 4)
	n = n | (n >> 8)
	n = n | (n >> 16)
	return n + 1
}

// finds the lowest and highest non-zero values in a list of numbers
//
// returns 0, 0 if no arguments are provided (or all arguments are zero)
func LowHighNon0[N WholeNumber](ns ...N) (low, high N) {
	if len(ns) == 0 {
		return 0, 0
	}

	low = ns[0]
	high = ns[0]
	for _, n := range ns {
		if n != 0 {
			if n < low {
				low = n
			}
			if n > high {
				high = n
			}
		}
	}
	return low, high
}