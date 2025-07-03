package music

import "testing"

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

// Additional tests for edge cases
func TestParseNoteEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Note
	}{
		{"B# highest octave", "B#9", Note{6, 9, 1}},
		{"Cb lowest octave", "Cb0", Note{0, 0, -1}},
		{"mixed case with alteration", "Ab3", Note{5, 3, -1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseNote(tt.input)
			if err != nil {
				t.Errorf("ParseNote() unexpected error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("ParseNote() = %v, want %v", got, tt.want)
			}
		})
	}
}
