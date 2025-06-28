package rules

import "testing"

func TestNoBeginWithFive(t *testing.T) {
	tests := []struct {
		name      string
		intervals []int
		want      bool
	}{
		{
			name:      "empty slice",
			intervals: []int{},
			want:      true,
		},
		{
			name:      "starts with 5",
			intervals: []int{5, 1, 2},
			want:      false,
		},
		{
			name:      "starts with 1",
			intervals: []int{1, 5, 2},
			want:      true,
		},
		{
			name:      "starts with other leap",
			intervals: []int{4, 5, 2},
			want:      true,
		},
		{
			name:      "single element 5",
			intervals: []int{5},
			want:      false,
		},
		{
			name:      "single element not 5",
			intervals: []int{3},
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NoBeginWithFive(tt.intervals); got != tt.want {
				t.Errorf("NoBeginWithFive() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestNoExcessiveNoteRepetition(t *testing.T) {
	tests := []struct {
		name      string
		intervals []int
		want      bool
	}{
		{
			name:      "empty slice",
			intervals: []int{},
			want:      true,
		},
		{
			name:      "no repetitions",
			intervals: []int{1, -1, 2, -2},
			want:      true,
		},
		{
			name:      "exactly 4 repetitions",
			intervals: []int{1, -1, 1, -1, 1, -1},
			want:      true,
		},
		{
			name:      "5 repetitions (violation)",
			intervals: []int{1, -1, 1, -1, 1, -1, 1, -1, 1, -1},
			want:      false,
		},
		{
			name:      "multiple notes with 4 reps",
			intervals: []int{1, 1, -2, 2, -2, 2, -1, 2, -1},
			want:      true,
		},
		{
			name:      "one note with 5 reps",
			intervals: []int{0, 0, 0, 0, 0},
			want:      false,
		},
		{
			name:      "complex pattern valid",
			intervals: []int{2, -1, 1, -2, 3, -3, 1, -1, 2, -1},
			want:      true,
		},
		{
			name:      "complex pattern invalid",
			intervals: []int{2, -2, 2, -2, 2, -2, 2, -2, 2, -2, 2},
			want:      false,
		},
		{
			name:      "partial slice valid",
			intervals: []int{1, -1, 1, -1},
			want:      true,
		},
		{
			name:      "partial slice invalid",
			intervals: []int{0, 0, 0, 0, 0},
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NoExcessiveNoteRepetition(tt.intervals); got != tt.want {
				t.Errorf("NoExcessiveNoteRepetition() = %v, want %v", got, tt.want)
			}
		})
	}
}
