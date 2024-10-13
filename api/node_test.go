package api

import (
	"testing"
)

type mockPositioned struct {
	Position
}

func TestWeakenRangeFrom(t *testing.T) {
	tests := []struct {
		name     string
		x        mockPositioned
		ys       []mockPositioned
		expected Position
	}{
		{
			name:     "Single position",
			x:        mockPositioned{MakePosition(1, 5)},
			ys:       []mockPositioned{},
			expected: MakePosition(1, 5),
		},
		{
			name:     "Single position, negative",
			x:        mockPositioned{MakePosition(-1, -5)},
			ys:       []mockPositioned{},
			expected: ZeroPosition(),
		},
		{
			name:     "Single position, zero",
			x:        mockPositioned{ZeroPosition()},
			ys:       []mockPositioned{},
			expected: ZeroPosition(),
		},
		{
			name:     "Two positions",
			x:        mockPositioned{MakePosition(1, 5)},
			ys:       []mockPositioned{{MakePosition(3, 7)}},
			expected: MakePosition(1, 7),
		},
		{
			name:     "Multiple positions",
			x:        mockPositioned{MakePosition(1, 5)},
			ys:       []mockPositioned{{MakePosition(3, 7)}, {MakePosition(0, 4)}},
			expected: MakePosition(1, 7),
		},
		{
			name:     "Overlapping positions",
			x:        mockPositioned{MakePosition(1, 5)},
			ys:       []mockPositioned{{MakePosition(2, 4)}, {MakePosition(3, 6)}},
			expected: MakePosition(1, 6),
		},
		{
			name:     "Negative positions",
			x:        mockPositioned{MakePosition(-1, 5)},
			ys:       []mockPositioned{{MakePosition(3, 4)}, {MakePosition(-3, 4)}},
			expected: MakePosition(3, 5),
		},
		{
			name:     "Zero as minimum",
			x:        mockPositioned{MakePosition(0, 3)},
			ys:       []mockPositioned{{MakePosition(3, 0)}, {MakePosition(0, 3)}},
			expected: MakePosition(0, 3),
		},
		{
			name:     "Zero as maximum",
			x:        mockPositioned{ZeroPosition()},
			ys:       []mockPositioned{{ZeroPosition()}, {ZeroPosition()}},
			expected: ZeroPosition(),
		},
		{
			name:     "Zero as minimum and maximum, negative input only",
			x:        mockPositioned{MakePosition(-49, -113)},
			ys:       []mockPositioned{{MakePosition(-193, -1)}, {MakePosition(-88, -97)}},
			expected: ZeroPosition(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WeakenRangeOver(tt.x, tt.ys...)
			if result != tt.expected {
				t.Errorf("WeakenRangeFrom() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
