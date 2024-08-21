// =================================================================================================
// Alex Peters - January 30, 2024
// =================================================================================================
package common

import (
	"testing"
)

func TestNumDigits(t *testing.T) {
	tests := []struct {
		num    uint64
		base   uint8
		expect uint64
	}{
		{num: 0, base: 255, expect: 1},
		{num: 0, base: 16, expect: 1},
		{num: 0, base: 10, expect: 1},
		{num: 0, base: 8, expect: 1},
		{num: 0, base: 7, expect: 1},
		{num: 0, base: 2, expect: 1},
		{num: 0, base: 1, expect: 0},

		{num: 9, base: 255, expect: 1},
		{num: 9, base: 16, expect: 1},
		{num: 9, base: 10, expect: 1},
		{num: 9, base: 8, expect: 2},
		{num: 9, base: 7, expect: 2},
		{num: 9, base: 2, expect: 4},
		{num: 9, base: 1, expect: 9},

		{num: 255, base: 255, expect: 2},
		{num: 16, base: 16, expect: 2},
		{num: 10, base: 10, expect: 2},
		{num: 8, base: 8, expect: 2},
		{num: 7, base: 7, expect: 2},
		{num: 2, base: 2, expect: 2},
		{num: 1, base: 1, expect: 1},

		{num: 18446744073709551615, base: 255, expect: 9},
		{num: 18446744073709551615, base: 16, expect: 16},
		{num: 18446744073709551615, base: 10, expect: 20},
		{num: 18446744073709551615, base: 8, expect: 22},
		{num: 18446744073709551615, base: 7, expect: 23},
		{num: 18446744073709551615, base: 2, expect: 64},
		{num: 18446744073709551615, base: 1, expect: 18446744073709551615},
	}
	for _, test := range tests {
		actual := NumDigits(test.num, test.base)
		if actual != test.expect {
			t.Fatalf("unexpected (%d): got %d", test.expect, actual)
		}
	}
}
