package common

type Number interface {
	comparable
	SignedNumber | UnsignedNumber
}

type SignedNumber interface {
	comparable
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

type UnsignedNumber interface {
	comparable
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// fixes the range of a number within 0 (inclusive) to the max value for that type (inclusive)
//
// in other words, it returns a value for the parameterized type that is in the non-negative subset of that type
//
//	NonNegative(x) = max(x, 0)
func NonNegative[n Number](x n) n {
	return max(x, 0)
}

func MinPositive[n Number](x n, ys ...n) (minimum n) {
	startIndex := 0
	ysLen := len(ys)
	// find the first positive value
	for minimum = NonNegative(x); minimum == 0 && startIndex < ysLen; startIndex++ {
		startIndex++
		minimum = NonNegative(ys[startIndex])
	}

	// find the minimum positive value (or skips and just returns 0 if there are no positive values)
	for _, y := range ys[startIndex:] {
		if y > 0 { // only consider positive values
			minimum = min(minimum, y)
		}
	}
	return minimum
}

// returns the minimum non-negative value of the given numbers. If all the numbers are negative, it returns 0
func MinNonNegative[n Number](x n, ys ...n) n {
	minimum := NonNegative(x)
	for _, y := range ys {
		minimum = min(minimum, NonNegative(y))
	}
	return minimum
}

func Abs[n Number](x n) n {
	if x < 0 {
		return -x
	}
	return x
}

func Sgn[n Number](x n) int {
	if x == 0 {
		return 0
	} else if x < 0 {
		return -1
	}
	return 1
}
