package pair

import (
	"testing"
)

func TestWeakenWith(t *testing.T) {
	tests := []struct {
		name     string
		p1       Data[int, int]
		p2       Data[int, int]
		expected Data[int, int]
	}{
		{
			name:     "both positive",
			p1:       Data[int, int]{3, 5},
			p2:       Data[int, int]{7, 2},
			expected: Data[int, int]{2, 7},
		},
		{
			name:     "one negative in p1",
			p1:       Data[int, int]{-3, 5},
			p2:       Data[int, int]{7, 2},
			expected: Data[int, int]{2, 7},
		},
		{
			name:     "one negative in p2",
			p1:       Data[int, int]{3, 5},
			p2:       Data[int, int]{-7, 2},
			expected: Data[int, int]{3, 5},
		},
		{
			name:     "both negative",
			p1:       Data[int, int]{-3, -5},
			p2:       Data[int, int]{-7, -2},
			expected: Data[int, int]{0, 0},
		},
		{
			name:     "mixed positive and zero",
			p1:       Data[int, int]{0, 5},
			p2:       Data[int, int]{7, 0},
			expected: Data[int, int]{0, 7},
		},
		{
			name: "same values",
			p1:   Data[int, int]{3, 5},
			p2:   Data[int, int]{3, 5},
			expected: Data[int, int]{3, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WeakenRange(tt.p1, tt.p2)
			if result != tt.expected {
				t.Errorf("WeakenWith(%v, %v) = %v; expected %v", tt.p1, tt.p2, result, tt.expected)
			}
		})
	}
}
