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

func TestNoRangeExceedsDecima(t *testing.T) {
	tests := []struct {
		name      string
		intervals []int
		expected  bool
	}{
		{
			name:      "empty slice",
			intervals: []int{},
			expected:  true,
		},
		{
			name:      "single interval within range",
			intervals: []int{5},
			expected:  true,
		},
		{
			name:      "exact decima range",
			intervals: []int{5, 4},
			expected:  true,
		},
		{
			name:      "exceeds decima",
			intervals: []int{5, 5},
			expected:  false,
		},
		{
			name:      "multiple steps within range",
			intervals: []int{1, 1, 1, 1, 1, -1, -1, -1, -1, -1},
			expected:  true,
		},
		{
			name:      "ascending then descending within range",
			intervals: []int{5, -3, 4, -2},
			expected:  true,
		},
		{
			name:      "exceeds decima with negative intervals",
			intervals: []int{-5, -5},
			expected:  false,
		},
		{
			name:      "boundary case - just below decima",
			intervals: []int{4, 4},
			expected:  true,
		},
		{
			name:      "boundary case - exactly decima",
			intervals: []int{4, 5},
			expected:  true,
		},
		{
			name:      "boundary case - just above decima",
			intervals: []int{5, 5},
			expected:  false,
		},
		{
			name:      "complex pattern within range",
			intervals: []int{3, -1, 4, -2, 2, -3, 1, -4},
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NoRangeExceedsDecima(tt.intervals)
			if got != tt.expected {
				t.Errorf("NoRangeExceedsDecima(%v) = %v, want %v", tt.intervals, got, tt.expected)
			}
		})
	}
}

func TestNoRepeatingPatterns(t *testing.T) {
	tests := []struct {
		name      string
		intervals []int
		want      bool
	}{
		{
			name:      "Empty slice",
			intervals: []int{},
			want:      true,
		},
		{
			name:      "Too short sequence",
			intervals: []int{1, -1},
			want:      true,
		},
		{
			name:      "No patterns - valid sequence",
			intervals: []int{1, 2, -1, 2, -2, 1},
			want:      true,
		},
		{
			name:      "Repeating pattern of length 2 (a,b,a,b)",
			intervals: []int{1, -1, 1, -1}, // Heights: [0, 1, 0, 1, 0]
			want:      false,
		},
		{
			name:      "Repeating pattern of length 3 (a,b,c,a,b,c)",
			intervals: []int{1, 2, -3, 1, 2, -3}, // Heights: [0, 1, 3, 0, 1, 3, 0]
			want:      false,
		},
		{
			name:      "Partial repeating pattern at start",
			intervals: []int{1, -1, 1}, // Heights: [0, 1, 0, 1]
			want:      false,
		},
		{
			name:      "Repeating pattern not at start",
			intervals: []int{2, 1, -1, 1, -1, 1}, // Heights: [0, 2, 3, 2, 3, 2, 3]
			want:      false,
		},
		{
			name:      "Large intervals but no repeating patterns",
			intervals: []int{5, -3, 4, -4, 2, -2},
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NoRepeatingPatterns(tt.intervals); got != tt.want {
				t.Errorf("NoRepeatingPatterns() = %v, want %v (intervals: %v)", got, tt.want, tt.intervals)
			}
		})
	}
}
