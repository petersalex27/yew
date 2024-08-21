// =================================================================================================
// Alex Peters - January 30, 2024
//
// Counting various things
// =================================================================================================
package common

// counts number of digits that would be in the number `num` if it was a number string written in
// base `base`.
//
// given a base > 1, returns an integer `digits`: 0 < `digits` < 65
//
// if given a base == 1, `digits` == `num`
//
// panics with "illegal argument: base cannot be zero" if `base` is 0
func NumDigits(num uint64, base uint8) (digits uint64) {
	if base == 0 {
		panic("illegal argument: base cannot be zero")
	} else if base == 1 {
		return num
	}

	digits = 1
	num = num / uint64(base)
	for base := uint64(base); num > 0; num = num / base {
		digits++
	}
	return digits
}
