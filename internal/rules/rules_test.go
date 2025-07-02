package rules

import (
	"testing"
)

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
			name:      "exactly 3 repetitions",
			intervals: []int{2, -1, -1, -1, 1}, // C4 E4 D4 C4 B3 C4
			want:      true,
		},
		{
			name:      "4 repetitions (violation)",
			intervals: []int{3, -1, -2, 1, -1, 4, -4}, // C F E C D C G C
			want:      false,
		},
		{
			name:      "multiple notes with 3 reps",
			intervals: []int{2, -1, -1, 3, -1, 2, -1, -1, -1, -1}, // Fux
			want:      true,
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

func TestPreparedLeaps(t *testing.T) {
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
			name:      "single small step",
			intervals: []int{2},
			want:      true,
		},
		{
			name:      "fourth leap with contrary motion",
			intervals: []int{-2, 3},
			want:      true,
		},
		{
			name:      "fourth leap without contrary motion",
			intervals: []int{1, 3},
			want:      false,
		},
		{
			name:      "fifth leap with single contrary step",
			intervals: []int{-2, 4},
			want:      true,
		},
		{
			name:      "fifth leap with double contrary motion",
			intervals: []int{-1, -1, 4},
			want:      true,
		},
		{
			name:      "fifth leap without preparation",
			intervals: []int{1, 4},
			want:      false,
		},
		{
			name:      "sixth leap with triple descending",
			intervals: []int{-1, -1, -1, 5},
			want:      true,
		},
		{
			name:      "sixth leap with double descending leaps",
			intervals: []int{-3, -2, 5},
			want:      true,
		},
		{
			name:      "sixth leap with single large descending",
			intervals: []int{-4, 5},
			want:      true,
		},
		{
			name:      "sixth leap without preparation",
			intervals: []int{1, 5}, // не подготовлен
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PreparedLeaps(tt.intervals); got != tt.want {
				t.Errorf("PreparedLeaps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateLeapResolution(t *testing.T) {
	tests := []struct {
		name      string
		intervals []int
		want      bool
	}{
		{
			name:      "empty sequence",
			intervals: []int{},
			want:      true,
		},
		{
			name:      "single element",
			intervals: []int{1},
			want:      true,
		},
		{
			name:      "no leaps",
			intervals: []int{1, 2, -1, 2, -2, 1},
			want:      true,
		},
		{
			name:      "single leap at end (no resolution needed)",
			intervals: []int{1, 1, 3},
			want:      true,
		},
		{
			name:      "properly resolved fourth (3)",
			intervals: []int{3, -1},
			want:      true,
		},
		{
			name:      "properly resolved fourth (-3)",
			intervals: []int{-3, 1},
			want:      true,
		},
		{
			name:      "unresolved fourth (3)",
			intervals: []int{3, 1},
			want:      false,
		},
		{
			name:      "properly resolved fifth (4)",
			intervals: []int{4, -2},
			want:      true,
		},
		{
			name:      "properly resolved fifth with two steps (4)",
			intervals: []int{4, -1, -1},
			want:      true,
		},
		{
			name:      "unresolved fifth (4)",
			intervals: []int{4, 1},
			want:      false,
		},
		{
			name:      "properly resolved sixth (5)",
			intervals: []int{5, -3},
			want:      true,
		},
		{
			name:      "properly resolved sixth with multiple steps (5)",
			intervals: []int{5, -1, -2},
			want:      true,
		},
		{
			name:      "unresolved sixth (5)",
			intervals: []int{5, 1},
			want:      false,
		},
		{
			name:      "multiple leaps all resolved",
			intervals: []int{3, -1, 4, -2, 1, 5, -3},
			want:      true,
		},
		{
			name:      "multiple leaps with one unresolved",
			intervals: []int{3, -1, 4, 1, 5, -3},
			want:      false,
		},
		{
			name:      "leap at beginning and end",
			intervals: []int{3, -1, 1, 1, 5},
			want:      true,
		},
		{
			name:      "consecutive leaps with proper resolution",
			intervals: []int{3, -3, 4, -4},
			want:      true,
		},
		{
			name:      "consecutive leaps with improper resolution",
			intervals: []int{3, 3, 4, -4},
			want:      false,
		},
		{
			name:      "complex sequence with mixed leaps",
			intervals: []int{1, 3, -1, 2, 5, -2, -1, 4, -2, 1, -3, 1},
			want:      true,
		},
		{
			name:      "complex sequence with unresolved leap",
			intervals: []int{1, 3, -1, 2, 5, -1, -1, 4, -2, 1, -3, 1},
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateLeapResolution(tt.intervals); got != tt.want {
				t.Errorf("ValidateLeapResolution() = %v, want %v for sequence %v", got, tt.want, tt.intervals)
			}
		})
	}
}

func TestNoTripleAlternatingNote(t *testing.T) {
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
			name:      "too short sequence",
			intervals: []int{1, -1, 1},
			want:      true,
		},
		{
			name:      "no alternating pattern",
			intervals: []int{1, 2, 3, 4, 5},
			want:      true,
		},
		{
			name:      "simple alternating pattern (a, b, a, c, a)",
			intervals: []int{1, -1, 1, -1},
			want:      false,
		},
		{
			name:      "complex alternating pattern",
			intervals: []int{1, 1, 3, -4, 4, 2, -2},
			want:      false,
		},
		{
			name:      "pattern at beginning",
			intervals: []int{1, -1, 4, -4, 1, -1},
			want:      false,
		},
		{
			name:      "pattern at end",
			intervals: []int{2, -1, -1, -1, -2, 2, 1, -1, 1},
			want:      false,
		},
		{
			name:      "pattern in middle",
			intervals: []int{2, 2, 1, -1, -3, 3, -2, -2},
			want:      false,
		},
		{
			name:      "no pattern with same notes but not alternating",
			intervals: []int{1, 1, 1, 1, 1},
			want:      true,
		},
		{
			name:      "long sequence without pattern",
			intervals: []int{2, -1, -1, 3, -1, 2, -1, -1, -1, -1},
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NoTripleAlternatingNote(tt.intervals); got != tt.want {
				t.Errorf("NoTripleAlternatingNote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoNoteRepetitionAfterLeap(t *testing.T) {
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
			name:      "single interval",
			intervals: []int{3},
			want:      true,
		},
		{
			name:      "equal leaps opposite direction",
			intervals: []int{3, -3},
			want:      false,
		},
		{
			name:      "equal leaps same direction",
			intervals: []int{3, 3},
			want:      true,
		},
		{
			name:      "unequal leaps opposite direction",
			intervals: []int{3, -4},
			want:      true,
		},
		{
			name:      "thirds",
			intervals: []int{2, -2},
			want:      false,
		},
		{
			name:      "large equal leaps opposite direction",
			intervals: []int{-5, 5},
			want:      false,
		},
		{
			name:      "multiple intervals no violation",
			intervals: []int{3, 2, -3, 1, 4, -2},
			want:      true,
		},
		{
			name:      "violation in middle",
			intervals: []int{2, 3, -3, 1, 4},
			want:      false,
		},
		{
			name:      "violation at end",
			intervals: []int{1, 2, 4, -4},
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NoNoteRepetitionAfterLeap(tt.intervals); got != tt.want {
				t.Errorf("NoNoteRepetitionAfterLeap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoRepeatingExtremes(t *testing.T) {
	tests := []struct {
		name      string
		intervals []int
		want      bool
	}{
		{
			name:      "empty sequence",
			intervals: []int{},
			want:      true,
		},
		{
			name:      "too short sequence",
			intervals: []int{1, 2},
			want:      true,
		},
		{
			name:      "simple ascending",
			intervals: []int{1, 1, 1, 1},
			want:      true,
		},
		{
			name:      "simple descending",
			intervals: []int{-1, -1, -1, -1},
			want:      true,
		},
		{
			name:      "single peak",
			intervals: []int{1, 1, -2},
			want:      true,
		},
		{
			name:      "single valley",
			intervals: []int{-1, -1, -1, 3},
			want:      true,
		},
		{
			name:      "repeating peaks", // C D E A E D C
			intervals: []int{1, 1, -4, 4, -1, -1},
			want:      false,
		},
		{
			name:      "repeating valleys", // C B D E B C
			intervals: []int{-1, 2, 1, -3, 1},
			want:      false,
		},
		{
			name:      "complex valid cantus firmus by Fus", // D F E D G F A G F E D
			intervals: []int{2, -1, -1, 3, -1, 2, -1, -1, -1, -1},
			want:      true,
		},
		{
			name:      "complex invalid sequence", // D F E D F C D
			intervals: []int{2, -1, -1, 2, -3, 1},
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NoRepeatingExtremes(tt.intervals); got != tt.want {
				t.Errorf("NoRepeatingExtremes() = %v, want %v (intervals: %v)", got, tt.want, tt.intervals)
			}
		})
	}
}

func TestAvoidSeventhBetweenExtrema(t *testing.T) {
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
			name:      "Single interval (no extrema)",
			intervals: []int{2},
			want:      true,
		},
		{
			name:      "seventh between the first and last note",
			intervals: []int{2, 2, 2}, // C E G B
			want:      false,
		},
		{
			name:      "Seventh between first note and first peak",
			intervals: []int{6, -1}, // C B A
			want:      false,
		},
		{
			name:      "Seventh between adjacent extrema",
			intervals: []int{1, -1, 1, 5, -1, -1, -1}, // C D C D B A G F
			want:      false,
		},
		{
			name:      "Seventh between last note and previous extremum",
			intervals: []int{1, 1, -2, 6}, // C D E C B
			want:      false,
		},
		{
			name:      "Valid melody by Fux (no sevenths)",
			intervals: []int{2, -1, -1, 3, -1, 2, -1, -1, -1, -1}, // D F E D G F A G F E D
			want:      true,
		},
		{
			name:      "Large leap but not seventh",
			intervals: []int{7, -7},
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AvoidSeventhBetweenExtrema(tt.intervals); got != tt.want {
				t.Errorf("AvoidSeventhBetweenExtrema() = %v, want %v (intervals: %v)", got, tt.want, tt.intervals)
			}
		})
	}
}

func TestMinDirectionChanges(t *testing.T) {
	tests := []struct {
		name      string
		intervals []int
		want      bool
	}{
		{
			name:      "empty slice",
			intervals: []int{},
			want:      false,
		},
		{
			name:      "single interval",
			intervals: []int{1},
			want:      false,
		},
		{
			name:      "two intervals, same direction",
			intervals: []int{1, 2},
			want:      false,
		},
		{
			name:      "two intervals, different directions",
			intervals: []int{1, -1},
			want:      false, // still only 1 change (need at least 2)
		},
		{
			name:      "three intervals, two changes",
			intervals: []int{1, -1, 1},
			want:      true,
		},
		{
			name:      "four intervals, one change",
			intervals: []int{1, 2, 3, -1},
			want:      false,
		},
		{
			name:      "complex melody with multiple changes",
			intervals: []int{1, -2, 3, -1, 2, -3, 1},
			want:      true,
		},
		{
			name:      "all same direction",
			intervals: []int{1, 1, 1, 1, 1},
			want:      false,
		},
		{
			name:      "alternating directions",
			intervals: []int{1, -1, 1, -1, 1},
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MinDirectionChanges(tt.intervals); got != tt.want {
				t.Errorf("MinDirectionChanges() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: fix cases
func TestValidateClimax(t *testing.T) {
	tests := []struct {
		name      string
		intervals []int
		want      bool
	}{
		{
			name:      "single climax, all heights are >= 0",
			intervals: []int{1, 1, 2, -1, -2, -1}, // C D E G F D C
			want:      true,
		},
		{
			name:      "multiple climaxes, all heights are >= 0",
			intervals: []int{1, 2, -1, 1, -2, -1}, // C D F E F D C
			want:      false,
		},
		{
			name:      "single climax, all heights are <= 0",
			intervals: []int{-1, -1, -1, 1, 1, 1}, // C4 B3 A3 G3 A3 B3 C4
			want:      true,
		},
		{
			name:      "multiple climaxes, all heights are <= 0",
			intervals: []int{-2, -1, 1, -1, 2, 1}, // C4 A3 G3 A3 G3 B3 C4
			want:      false,
		},
		{
			name:      "single max and min",
			intervals: []int{-1, -1, 5, -1, -2}, // C4 B3 A3 F4 E4 C4
			want:      true,
		},
		{
			name:      "Mixed heights - multiple max",
			intervals: []int{-1, 1, 3, -1, 1, -2, -1}, // C4 B3 C4 F4 E4 F4 D4 C4
			want:      false,
		},
		{
			name:      "Mixed heights - multiple min",
			intervals: []int{-2, -2, 1, 1, -2, 3, 1, 1, -1}, // C4 A3 F3 G3 A3 F3 B3 C4 D4 C4
			want:      false,
		},
		{
			name:      "Empty sequence",
			intervals: []int{},
			want:      true,
		},
		{
			name:      "Single interval - positive",
			intervals: []int{1},
			want:      true,
		},
		{
			name:      "Single interval - negative",
			intervals: []int{-1},
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateClimax(tt.intervals)
			if got != tt.want {
				t.Errorf("ValidateClimax() = %v, want %v for sequence %v", got, tt.want, tt.intervals)
			}
		})
	}
}

func TestNoSequences(t *testing.T) {
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
			name:      "Too short for any pattern",
			intervals: []int{1, 2},
			want:      true,
		},
		// Part a) Alternating pattern tests
		{
			name:      "Alternating pattern a,b,a,b,a",
			intervals: []int{1, -1, 1, -1, 1},
			want:      false,
		},
		{
			name:      "Alternating pattern with different values",
			intervals: []int{2, -1, 2, -1, 2},
			want:      false,
		},
		{
			name:      "Not alternating (a == b)",
			intervals: []int{1, 1, 1, 1, 1},
			want:      true,
		},
		// Part b) Consecutive patterns with one-element separation
		{
			name:      "Consecutive patterns with separation",
			intervals: []int{1, 2, 3, 0, 1, 2, 3},
			want:      false,
		},
		{
			name:      "Consecutive patterns with leap separation",
			intervals: []int{4, -2, 3, 5, 4, -2, 3},
			want:      false,
		},
		{
			name:      "Longer sequence with multiple repeats",
			intervals: []int{2, 3, -1, 0, 2, 3, -1, 0, 2, 3, -1},
			want:      false,
		},
		{
			name:      "Longer sequence. No repeating patterns",
			intervals: []int{2, 3, -1, 1, 5, 2, 3, -1, 2, 1, 2, 3, -1},
			want:      true,
		},
		// Part c) Leap pattern repetition tests
		{
			name:      "Repeating leap pattern",
			intervals: []int{3, -2, 4, 1, 3, -2, 4, 1},
			want:      false,
		},
		{
			name:      "Repeating leap pattern with different surrounding",
			intervals: []int{1, 3, -2, 4, -1, -1, 2, 3, -2, 4, -1, -1, 2},
			want:      false,
		},
		{
			name:      "No repeating patterns",
			intervals: []int{1, 2, -1, 3, -2, 4, -3},
			want:      true,
		},
		{
			name:      "Leaps but no repetition",
			intervals: []int{3, -2, 4, 2, -3, 5, -4},
			want:      true,
		},
		// Edge cases
		{
			name:      "Almost alternating but not quite",
			intervals: []int{1, -1, 1, -1, 2},
			want:      true,
		},
		{
			name:      "Almost consecutive pattern but different",
			intervals: []int{1, 2, 3, 1, 2, 4},
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NoSequences(tt.intervals); got != tt.want {
				t.Errorf("NoRepeatingPatterns() = %v, want %v for case %v", got, tt.want, tt.intervals)
			}
		})
	}
}
