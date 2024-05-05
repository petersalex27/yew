// =================================================================================================
// Alex Peters - February 17, 2024
// =================================================================================================
package common

// a pair (Left, Right): (T, U)
type Pair[T, U any] struct {
	Left  T
	Right U
}

// given values a and b, returns pair (a, b)
func Join[T, U any](a T, b U) Pair[T, U] { return Pair[T, U]{a, b} }

// given a pair p
//
//	p = (a, b)
//
// returns a, b
func (pair Pair[T, U]) Split() (left T, right U) { return pair.Left, pair.Right }

// given a pair p
//
//	p = (a, b)
//
// and a value v
//
//	Append(p, v) = (p, v) = ((a, b), v)
func Append[T, U, V any](pair Pair[T, U], v V) Pair[Pair[T, U], V] {
	return Pair[Pair[T, U], V]{pair, v}
}