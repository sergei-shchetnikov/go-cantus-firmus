package main

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

func TestCantusFirmus_ToNotes(t *testing.T) {
	tests := []struct {
		name      string
		intervals CantusFirmus
		startNote Note
		want      []Note
	}{
		{
			name:      "Empty intervals",
			intervals: CantusFirmus{},
			startNote: Note{Step: 0, Octave: 4}, // C4
			want:      []Note{{0, 4, 0}},
		},
		{
			name:      "Simple ascending",
			intervals: CantusFirmus{1, 1, 1},                              // up 2nd, up 2nd, up 2nd
			startNote: Note{0, 4, 0},                                      // C4
			want:      []Note{{0, 4, 0}, {1, 4, 0}, {2, 4, 0}, {3, 4, 0}}, // C4, D4, E4, F4
		},
		{
			name:      "Mixed intervals",
			intervals: CantusFirmus{1, 2, -1},                             // up 2nd, up 3rd, down 2nd
			startNote: Note{0, 4, 0},                                      // C4
			want:      []Note{{0, 4, 0}, {1, 4, 0}, {3, 4, 0}, {2, 4, 0}}, // C4, D4, F4, E4
		},
		{
			name:      "Octave crossing up",
			intervals: CantusFirmus{6, 1},                      // up 7th, up 2nd
			startNote: Note{6, 3, 0},                           // B3
			want:      []Note{{6, 3, 0}, {5, 4, 0}, {6, 4, 0}}, // B3, A4, B4
		},
		{
			name:      "Octave crossing down",
			intervals: CantusFirmus{-6, -1},                    // down 7th, down 2nd
			startNote: Note{0, 5, 0},                           // C5
			want:      []Note{{0, 5, 0}, {1, 4, 0}, {0, 4, 0}}, // C5, D4, C4
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.intervals.ToNotes(tt.startNote)
			if len(got) != len(tt.want) {
				t.Errorf("ToNotes() length = %d, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("ToNotes()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
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
		errText string
	}{
		{
			name:  "valid note C4",
			input: "C4",
			want:  Note{Step: 0, Octave: 4},
		},
		{
			name:  "valid note lowercase a5",
			input: "a5",
			want:  Note{Step: 5, Octave: 5},
		},
		{
			name:  "valid note B3",
			input: "B3",
			want:  Note{Step: 6, Octave: 3},
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
			errText: "invalid note format: string too short",
		},
		{
			name:    "too short",
			input:   "A",
			wantErr: true,
			errText: "invalid note format: string too short",
		},
		{
			name:    "invalid note char",
			input:   "X5",
			wantErr: true,
			errText: "invalid note character: X",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseNote(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseNote(%q) expected error, got nil", tt.input)
				} else if tt.errText != "" && err.Error() != tt.errText {
					t.Errorf("ParseNote(%q) error = %v, wantErr %v", tt.input, err.Error(), tt.errText)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseNote(%q) unexpected error: %v", tt.input, err)
				return
			}

			if got != tt.want {
				t.Errorf("ParseNote(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
