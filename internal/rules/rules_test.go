package rules

import "testing"

func TestNoFiveOfSameSign(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected bool
	}{
		{
			name:     "Empty slice",
			input:    []int{},
			expected: true, // Not enough elements to violate (min 5)
		},
		{
			name:     "Slice less than 5 elements (e.g., 4 elements)",
			input:    []int{1, 2, -3, 4},
			expected: true, // Not enough elements to violate (min 5)
		},
		{
			name:     "Exactly 5 elements, no violation (mixed signs)",
			input:    []int{1, -2, 3, -4, 5},
			expected: true,
		},
		{
			name:     "Exactly 5 elements, violation (all positive)",
			input:    []int{1, 2, 3, 4, 5},
			expected: false,
		},
		{
			name:     "Exactly 5 elements, violation (all negative)",
			input:    []int{-1, -2, -3, -4, -5},
			expected: false,
		},
		{
			name:     "More than 5 elements, violation in middle (positive)",
			input:    []int{1, -1, 2, 3, 4, 5, 6, -7}, // 2,3,4,5,6 is a violation
			expected: false,
		},
		{
			name:     "More than 5 elements, violation at end (negative)",
			input:    []int{1, 2, -3, -4, -5, -6, -7}, // -3,-4,-5,-6,-7 is a violation
			expected: false,
		},
		{
			name:     "More than 5 elements, no violation (alternating signs)",
			input:    []int{1, -1, 2, -2, 3, -3, 4, -4, 5, -5},
			expected: true,
		},
		{
			name:     "Four positive, then negative (no violation)",
			input:    []int{1, 2, 3, 4, -5, 6},
			expected: true,
		},
		{
			name:     "Four negative, then positive (no violation)",
			input:    []int{-1, -2, -3, -4, 5, -6},
			expected: true,
		},
		{
			name:     "Long slice, no violation",
			input:    []int{1, 2, 3, 4, -1, 2, 3, 4, -1, 2, 3, 4, -1},
			expected: true,
		},
		{
			name:     "Long slice, violation at start",
			input:    []int{1, 2, 3, 4, 5, -1, 2, -3, 4}, // 1,2,3,4,5 is a violation
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NoFiveOfSameSign(tt.input)
			if got != tt.expected {
				t.Errorf("NoFourOfSameSign(%v) = %v; expected %v", tt.input, got, tt.expected)
			}
		})
	}
}
