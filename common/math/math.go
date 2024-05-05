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


// returns the absolute value of some integer type
func Abs[T ~int](a T) T {
	if a > 0 {
		return a
	}
	return -a
}

// rounds up to the nearest power of two
func PowerOfTwoCeil(n uint) uint {
	n = n - 1
	n = n | (n >> 1)
	n = n | (n >> 2)
	n = n | (n >> 4)
	n = n | (n >> 8)
	n = n | (n >> 16)
	return n + 1
}