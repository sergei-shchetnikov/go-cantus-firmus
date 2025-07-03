package music

import "testing"

func TestCantusFirmus_Realize(t *testing.T) {
	tests := []struct {
		name        string
		cf          CantusFirmus
		mode        string
		wantNotes   []string
		wantErr     bool
		errContains string
	}{
		{
			name:      "Major mode - simple ascending",
			cf:        CantusFirmus{0, 1, 2}, // Unison, second up, third up
			mode:      "Major",
			wantNotes: []string{"C4", "C4", "D4", "F4"},
		},
		{
			name:      "Dorian mode - up and down",
			cf:        CantusFirmus{1, -1, 2}, // Second up, second down, third up
			mode:      "Dorian",
			wantNotes: []string{"D4", "E4", "D4", "F4"},
		},
		{
			name:      "Phrygian mode - descending",
			cf:        CantusFirmus{-1, -1, -2}, // Second down, second down, third down
			mode:      "Phrygian",
			wantNotes: []string{"E4", "D4", "C4", "A3"},
		},
		{
			name:      "Lydian mode - octave jump",
			cf:        CantusFirmus{7, -7}, // Octave up, octave down
			mode:      "Lydian",
			wantNotes: []string{"F4", "F5", "F4"},
		},
		{
			name:      "Mixolydian mode - complex",
			cf:        CantusFirmus{1, 2, -3, 4}, // Second up, third up, fourth down, fifth up
			mode:      "Mixolydian",
			wantNotes: []string{"G4", "A4", "C5", "G4", "D5"},
		},
		{
			name:      "Minor mode - simple",
			cf:        CantusFirmus{1, 1, -2}, // Second up, second up, third down
			mode:      "Minor",
			wantNotes: []string{"A4", "B4", "C5", "A4"},
		},
		{
			name:      "Locrian mode - large leap",
			cf:        CantusFirmus{5, -5}, // Sixth up, sixth down
			mode:      "Locrian",
			wantNotes: []string{"B4", "G5", "B4"},
		},
		{
			name:        "Invalid mode",
			cf:          CantusFirmus{1, 2},
			mode:        "Blues",
			wantErr:     true,
			errContains: "unknown mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cf.Realize(tt.mode)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Realize() expected error, got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("Realize() error = %v, want containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("Realize() unexpected error: %v", err)
				return
			}

			if len(got) != len(tt.wantNotes) {
				t.Errorf("Realize() length mismatch: got %d notes, want %d", len(got), len(tt.wantNotes))
				return
			}

			for i, note := range got {
				if note.String() != tt.wantNotes[i] {
					t.Errorf("Realize() note %d = %v, want %v", i, note.String(), tt.wantNotes[i])
				}
			}
		})
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}

func TestAdjustMinorAlterations(t *testing.T) {
	tests := []struct {
		name        string
		input       Realization
		expected    Realization
		description string
	}{
		{
			name: "A_G_A pattern",
			input: Realization{
				{Step: 5, Octave: 4, Alteration: 0}, // A
				{Step: 4, Octave: 4, Alteration: 0}, // G
				{Step: 5, Octave: 4, Alteration: 0}, // A
			},
			expected: Realization{
				{Step: 5, Octave: 4, Alteration: 0},
				{Step: 4, Octave: 4, Alteration: 1}, // G#
				{Step: 5, Octave: 4, Alteration: 0},
			},
			description: "Should add sharp to G in A-G-A pattern",
		},
		{
			name: "F_G_A pattern",
			input: Realization{
				{Step: 3, Octave: 4, Alteration: 0}, // F
				{Step: 4, Octave: 4, Alteration: 0}, // G
				{Step: 5, Octave: 4, Alteration: 0}, // A
			},
			expected: Realization{
				{Step: 3, Octave: 4, Alteration: 1}, // F#
				{Step: 4, Octave: 4, Alteration: 1}, // G#
				{Step: 5, Octave: 4, Alteration: 0},
			},
			description: "Should add sharps to F and G in F-G-A pattern",
		},
		{
			name: "No alteration pattern",
			input: Realization{
				{Step: 5, Octave: 4, Alteration: 0}, // A
				{Step: 3, Octave: 4, Alteration: 0}, // F
				{Step: 5, Octave: 4, Alteration: 0}, // A
			},
			expected: Realization{
				{Step: 5, Octave: 4, Alteration: 0},
				{Step: 3, Octave: 4, Alteration: 0},
				{Step: 5, Octave: 4, Alteration: 0},
			},
			description: "Should not alter notes when no pattern matches",
		},
		{
			name: "Already altered notes",
			input: Realization{
				{Step: 5, Octave: 4, Alteration: 0},  // A
				{Step: 4, Octave: 4, Alteration: -1}, // Gb
				{Step: 5, Octave: 4, Alteration: 0},  // A
			},
			expected: Realization{
				{Step: 5, Octave: 4, Alteration: 0},
				{Step: 4, Octave: 4, Alteration: -1}, // Gb remains unchanged
				{Step: 5, Octave: 4, Alteration: 0},
			},
			description: "Should not alter already altered notes",
		},
		{
			name: "Short sequence",
			input: Realization{
				{Step: 5, Octave: 4, Alteration: 0}, // A
				{Step: 4, Octave: 4, Alteration: 0}, // G
			},
			expected: Realization{
				{Step: 5, Octave: 4, Alteration: 0},
				{Step: 4, Octave: 4, Alteration: 0},
			},
			description: "Should not alter notes in sequences shorter than 3",
		},
		{
			name: "Multiple patterns",
			input: Realization{
				{Step: 5, Octave: 4, Alteration: 0}, // A
				{Step: 4, Octave: 4, Alteration: 0}, // G
				{Step: 5, Octave: 4, Alteration: 0}, // A
				{Step: 3, Octave: 4, Alteration: 0}, // F
				{Step: 4, Octave: 4, Alteration: 0}, // G
				{Step: 5, Octave: 4, Alteration: 0}, // A
			},
			expected: Realization{
				{Step: 5, Octave: 4, Alteration: 0},
				{Step: 4, Octave: 4, Alteration: 1}, // G#
				{Step: 5, Octave: 4, Alteration: 0},
				{Step: 3, Octave: 4, Alteration: 1}, // F#
				{Step: 4, Octave: 4, Alteration: 1}, // G#
				{Step: 5, Octave: 4, Alteration: 0},
			},
			description: "Should handle multiple patterns in sequence",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := adjustMinorAlterations(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("%s: expected length %d, got %d", tt.description, len(tt.expected), len(result))
				return
			}

			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("%s: at index %d expected %v, got %v",
						tt.description, i, tt.expected[i], result[i])
				}
			}
		})
	}
}
