package common

// returns min {a, b} or a if a==b
func Min[T ~int](a, b T) T {
	if a > b {
		return b
	}
	return a
}

// returns max {a, b} or a if a==b
func Max[T ~int](a, b T) T {
	if a < b {
		return b
	}
	return a
}