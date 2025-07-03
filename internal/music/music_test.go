package music

import "testing"

func TestMod7(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{"Number less than 7", 5, 5},
		{"Number greater than -7", -5, 2},
		{"Positive number", 10, 3},
		{"Negative number", -10, 4},
		{"Positive multiple of 7", 14, 0},
		{"Negative multiple of 7", -7, 0},
		{"Zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Mod7(tt.n); got != tt.want {
				t.Errorf("Mod7(%d) = %d, want %d", tt.n, got, tt.want)
			}
		})
	}
}

func TestTranspose(t *testing.T) {
	tests := []struct {
		name     string
		note     Note
		interval Interval
		want     Note
	}{
		{
			name:     "unison",
			note:     Note{Step: 0, Octave: 4}, // C4
			interval: 0,
			want:     Note{Step: 0, Octave: 4}, // C4
		},
		{
			name:     "second up",
			note:     Note{Step: 0, Octave: 4}, // C4
			interval: 1,
			want:     Note{Step: 1, Octave: 4}, // D4
		},
		{
			name:     "third down",
			note:     Note{Step: 2, Octave: 4}, // E4
			interval: -2,
			want:     Note{Step: 0, Octave: 4}, // C4
		},
		{
			name:     "octave up",
			note:     Note{Step: 0, Octave: 4}, // C4
			interval: 7,
			want:     Note{Step: 0, Octave: 5}, // C5
		},
		{
			name:     "octave down",
			note:     Note{Step: 0, Octave: 4}, // C4
			interval: -7,
			want:     Note{Step: 0, Octave: 3}, // C3
		},
		{
			name:     "octave up from D",
			note:     Note{Step: 1, Octave: 4}, // D4
			interval: 7,
			want:     Note{Step: 1, Octave: 5}, // D5
		},
		{
			name:     "octave down from E",
			note:     Note{Step: 2, Octave: 4}, // E4
			interval: -7,
			want:     Note{Step: 2, Octave: 3}, // E3
		},
		{
			name:     "crossing octave boundary up",
			note:     Note{Step: 6, Octave: 4}, // B4
			interval: 1,
			want:     Note{Step: 0, Octave: 5}, // C5
		},
		{
			name:     "crossing octave boundary down",
			note:     Note{Step: 1, Octave: 4}, // D4
			interval: -2,
			want:     Note{Step: 6, Octave: 3}, // B3
		},
		{
			name:     "large interval up",
			note:     Note{Step: 0, Octave: 4}, // C4
			interval: 15,
			want:     Note{Step: 1, Octave: 6}, // D6
		},
		{
			name:     "large interval down",
			note:     Note{Step: 3, Octave: 4}, // F4
			interval: -10,
			want:     Note{Step: 0, Octave: 3}, // C3
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Transpose(tt.note, tt.interval)
			if got != tt.want {
				t.Errorf("Transpose(%v, %v) = %v, want %v", tt.note, tt.interval, got, tt.want)
			}
		})
	}
}

func TestNote_String(t *testing.T) {
	tests := []struct {
		name     string
		note     Note
		expected string
	}{
		{
			name:     "Middle C",
			note:     Note{Step: 0, Octave: 4},
			expected: "C4",
		},
		{
			name:     "B below Middle C",
			note:     Note{Step: 6, Octave: 3},
			expected: "B3",
		},
		{
			name:     "High G",
			note:     Note{Step: 4, Octave: 6},
			expected: "G6",
		},
		{
			name:     "C sharp",
			note:     Note{Step: 0, Octave: 4, Alteration: 1},
			expected: "C#4",
		},
		{
			name:     "D flat",
			note:     Note{Step: 1, Octave: 4, Alteration: -1},
			expected: "Db4",
		},
		{
			name:     "F sharp high octave",
			note:     Note{Step: 3, Octave: 5, Alteration: 1},
			expected: "F#5",
		},
		{
			name:     "G flat low octave",
			note:     Note{Step: 4, Octave: 2, Alteration: -1},
			expected: "Gb2",
		},
		{
			name:     "Natural A after alteration",
			note:     Note{Step: 5, Octave: 4, Alteration: 0},
			expected: "A4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.note.String(); got != tt.expected {
				t.Errorf("Note.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIntervalString(t *testing.T) {
	tests := []struct {
		name     string
		interval Interval
		want     string
	}{
		{"unison", 0, "unison"},

		{"second up", 1, "second up"},
		{"third up", 2, "third up"},
		{"fourth up", 3, "fourth up"},
		{"fifth up", 4, "fifth up"},
		{"sixth up", 5, "sixth up"},
		{"seventh up", 6, "seventh up"},
		{"octave up", 7, "octave up"},

		{"second down", -1, "second down"},
		{"third down", -2, "third down"},
		{"fourth down", -3, "fourth down"},
		{"fifth down", -4, "fifth down"},
		{"sixth down", -5, "sixth down"},
		{"seventh down", -6, "seventh down"},
		{"octave down", -7, "octave down"},

		{"9th up", 8, "9th up"},
		{"10th up", 9, "10th up"},
		{"11th up", 10, "11th up"},
		{"12th up", 11, "12th up"},
		{"13th up", 12, "13th up"},
		{"14th up", 13, "14th up"},
		{"15th up", 14, "15th up"},

		{"9th down", -8, "9th down"},
		{"10th down", -9, "10th down"},
		{"11th down", -10, "11th down"},
		{"12th down", -11, "12th down"},
		{"13th down", -12, "13th down"},
		{"14th down", -13, "14th down"},
		{"15th down", -14, "15th down"},

		{"16th up", 15, "16th up"},
		{"17th up", 16, "17th up"},
		{"21st up", 20, "21st up"},
		{"22nd up", 21, "22nd up"},
		{"23rd up", 22, "23rd up"},
		{"24th up", 23, "24th up"},
		{"31st up", 30, "31st up"},
		{"32nd down", -31, "32nd down"},
		{"33rd down", -32, "33rd down"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.interval.String(); got != tt.want {
				t.Errorf("Interval(%d).String() = %v, want %v", tt.interval, got, tt.want)
			}
		})
	}
}

func TestParseNote(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Note
		wantErr bool
	}{
		// Valid notes without alteration
		{"C4", "C4", Note{0, 4, 0}, false},
		{"D5", "D5", Note{1, 5, 0}, false},
		{"lowercase a4", "a4", Note{5, 4, 0}, false},
		{"E2", "E2", Note{2, 2, 0}, false},

		// Valid notes with alteration
		{"C#4", "C#4", Note{0, 4, 1}, false},
		{"Db5", "Db5", Note{1, 5, -1}, false},
		{"F#3", "F#3", Note{3, 3, 1}, false},
		{"Gb2", "Gb2", Note{4, 2, -1}, false},
		{"lowercase with alteration", "ab4", Note{5, 4, -1}, false},

		// Invalid formats
		{"empty string", "", Note{}, true},
		{"too short", "C", Note{}, true},
		{"invalid note char", "H4", Note{}, true},
		{"invalid alteration", "Cx4", Note{}, true},
		{"missing octave after alteration", "C#", Note{}, true},
		{"invalid octave", "CA", Note{}, true},
		{"double alteration", "C##4", Note{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseNote(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseNote() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
