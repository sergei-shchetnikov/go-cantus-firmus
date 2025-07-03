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
				{Step: 1, Octave: 4}, // D4 (M2)
				{Step: 2, Octave: 4}, // E4 (M2)
				{Step: 0, Octave: 5}, // C5 (m6)
			},
			expected: true,
		},
		{
			name: "Contains augmented interval",
			input: music.Realization{
				{Step: 0, Octave: 4},                // C4
				{Step: 3, Octave: 4, Alteration: 1}, // F#4 (A4)
			},
			expected: false,
		},
		{
			name: "Contains diminished interval",
			input: music.Realization{
				{Step: 3, Octave: 4},                 // F4
				{Step: 0, Octave: 5, Alteration: -1}, // Cb5 (d5)
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
				{Step: 1, Octave: 4},                // D4 (M2)
				{Step: 3, Octave: 4, Alteration: 1}, // F#4 (A4)
				{Step: 0, Octave: 5},                // C5 (P5)
			},
			expected: false,
		},
		{
			name: "All minor intervals",
			input: music.Realization{
				{Step: 5, Octave: 4}, // A4
				{Step: 0, Octave: 5}, // C5 (m3)
				{Step: 3, Octave: 5}, // F5 (m3)
			},
			expected: true,
		},
		{
			name: "All perfect intervals",
			input: music.Realization{
				{Step: 0, Octave: 4}, // C4
				{Step: 4, Octave: 4}, // G4 (P5)
				{Step: 0, Octave: 5}, // C5 (P4)
			},
			expected: true,
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
