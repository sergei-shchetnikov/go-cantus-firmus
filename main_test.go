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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.note.String(); got != tt.expected {
				t.Errorf("Note.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}
