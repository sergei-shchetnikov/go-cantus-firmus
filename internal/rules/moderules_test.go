package rules

import (
	"go-cantus-firmus/internal/music"
	"testing"
)

func TestIsFreeOfAugmentedDiminished(t *testing.T) {
	tests := []struct {
		name     string
		input    music.Realization
		expected bool
	}{
		{
			name: "All perfect/major/minor intervals",
			input: music.Realization{
				{Step: 0, Octave: 4}, // C4
				{Step: 1, Octave: 4}, // D4
				{Step: 2, Octave: 4}, // E4
				{Step: 0, Octave: 5}, // C5
			},
			expected: true,
		},
		{
			name: "Contains augmented interval",
			input: music.Realization{
				{Step: 0, Octave: 4},                // C4
				{Step: 3, Octave: 4, Alteration: 1}, // F#4
			},
			expected: false,
		},
		{
			name: "Contains diminished interval",
			input: music.Realization{
				{Step: 3, Octave: 4},                 // F4
				{Step: 0, Octave: 5, Alteration: -1}, // Cb5
			},
			expected: false,
		},
		{
			name: "Single note (no intervals)",
			input: music.Realization{
				{Step: 0, Octave: 4},
			},
			expected: true,
		},
		{
			name:     "Empty realization",
			input:    music.Realization{},
			expected: true,
		},
		{
			name: "Mixed valid and invalid intervals",
			input: music.Realization{
				{Step: 0, Octave: 4},                // C4
				{Step: 1, Octave: 4},                // D4
				{Step: 3, Octave: 4, Alteration: 1}, // F#4
				{Step: 0, Octave: 5},                // C5
			},
			expected: false,
		},
		{
			name: "All minor intervals",
			input: music.Realization{
				{Step: 5, Octave: 4}, // A4
				{Step: 0, Octave: 5}, // C5
				{Step: 3, Octave: 5}, // F5
			},
			expected: true,
		},
		{
			name: "All perfect intervals",
			input: music.Realization{
				{Step: 0, Octave: 4}, // C4
				{Step: 4, Octave: 4}, // G4
				{Step: 0, Octave: 5}, // C5
			},
			expected: true,
		},
		{
			name: "Valid step-2 intervals",
			input: music.Realization{
				{Step: 0, Octave: 4}, // C4
				{Step: 1, Octave: 4}, // D4
				{Step: 3, Octave: 4}, // F4
			},
			expected: true,
		},
		{
			name: "Invalid augmented step-2 interval",
			input: music.Realization{
				{Step: 0, Octave: 4},                // C4
				{Step: 1, Octave: 4},                // D4
				{Step: 3, Octave: 4, Alteration: 1}, // F#4
			},
			expected: false,
		},
		{
			name: "Invalid diminished step-2 interval",
			input: music.Realization{
				{Step: 3, Octave: 4},                 // F4
				{Step: 4, Octave: 4},                 // G4
				{Step: 0, Octave: 5, Alteration: -1}, // Cb5
			},
			expected: false,
		},
		{
			name: "Valid step-2 intervals with larger sequence",
			input: music.Realization{
				{Step: 0, Octave: 4}, // C4
				{Step: 1, Octave: 4}, // D4
				{Step: 3, Octave: 4}, // F4
				{Step: 4, Octave: 4}, // G4
				{Step: 5, Octave: 4}, // A4
			},
			expected: true,
		},
		{
			name: "Invalid step-2 interval in middle of sequence",
			input: music.Realization{
				{Step: 0, Octave: 4},                // C4
				{Step: 1, Octave: 4},                // D4
				{Step: 3, Octave: 4},                // F4
				{Step: 4, Octave: 4, Alteration: 1}, // G#4
				{Step: 5, Octave: 4},                // A4
			},
			expected: false,
		},
		{
			name: "Valid extremum intervals",
			input: music.Realization{
				{Step: 0, Octave: 4}, // C4 (extremum)
				{Step: 1, Octave: 4}, // D4
				{Step: 2, Octave: 4}, // E4 (extremum)
				{Step: 1, Octave: 4}, // D4
				{Step: 0, Octave: 4}, // C4 (extremum)
			},
			expected: true,
		},
		{
			name: "Invalid augmented interval between extremums",
			input: music.Realization{
				{Step: 0, Octave: 4},                // C4 (extremum)
				{Step: 1, Octave: 4},                // D4
				{Step: 3, Octave: 4, Alteration: 1}, // F#4 (extremum)
			},
			expected: false,
		},
		{
			name: "Invalid diminished interval between extremums",
			input: music.Realization{
				{Step: 3, Octave: 4},                 // F4 (extremum)
				{Step: 2, Octave: 4},                 // E4
				{Step: 0, Octave: 5, Alteration: -1}, // Cb5 (extremum)
			},
			expected: false,
		},
		{
			name: "Invalid augmented step-2 interval between extremums",
			input: music.Realization{
				{Step: 0, Octave: 4},                // C4 (extremum)
				{Step: 1, Octave: 4},                // D4
				{Step: 2, Octave: 4},                // E4
				{Step: 3, Octave: 4, Alteration: 1}, // F#4 (extremum)
			},
			expected: false,
		},
		{
			name: "Complex valid case with extremums",
			input: music.Realization{
				{Step: 0, Octave: 4}, // C4 (extremum)
				{Step: 1, Octave: 4}, // D4
				{Step: 2, Octave: 4}, // E4 (extremum)
				{Step: 1, Octave: 4}, // D4
				{Step: 0, Octave: 4}, // C4 (extremum)
				{Step: 4, Octave: 4}, // G4 (extremum)
				{Step: 3, Octave: 4}, // F4
				{Step: 2, Octave: 4}, // E4 (extremum)
			},
			expected: true,
		},
		{
			name: "Complex invalid case with extremums",
			input: music.Realization{
				{Step: 0, Octave: 4},                 // C4 (extremum)
				{Step: 1, Octave: 4},                 // D4
				{Step: 3, Octave: 4, Alteration: 1},  // F#4 (extremum)
				{Step: 2, Octave: 4},                 // E4
				{Step: 0, Octave: 5, Alteration: -1}, // Cb5 (extremum)
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsFreeOfAugmentedDiminished(tt.input)
			if result != tt.expected {
				t.Errorf("IsFreeOfAugmentedDiminished() = %v, want %v for case %q", result, tt.expected, tt.name)
			}
		})
	}
}
