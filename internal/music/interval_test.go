package music

import (
	"testing"
)

func TestCalculateIntervalQuality(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		note1    Note
		note2    Note
		expected string
		wantErr  bool
	}{
		// Perfect Intervals
		{"C4 to C4 (Unison)", Note{0, 4, 0}, Note{0, 4, 0}, "P", false},
		{"C4 to F4 (Perfect Fourth)", Note{0, 4, 0}, Note{3, 4, 0}, "P", false},
		{"C4 to G4 (Perfect Fifth)", Note{0, 4, 0}, Note{4, 4, 0}, "P", false},
		{"C4 to C5 (Perfect Octave)", Note{0, 4, 0}, Note{0, 5, 0}, "P", false},

		// Major/Minor Intervals
		{"C4 to D4 (Major Second)", Note{0, 4, 0}, Note{1, 4, 0}, "M", false},
		{"C4 to Db4 (Minor Second)", Note{0, 4, 0}, Note{1, 4, -1}, "m", false},
		{"C4 to E4 (Major Third)", Note{0, 4, 0}, Note{2, 4, 0}, "M", false},
		{"C4 to Eb4 (Minor Third)", Note{0, 4, 0}, Note{2, 4, -1}, "m", false},
		{"C4 to A4 (Major Sixth)", Note{0, 4, 0}, Note{5, 4, 0}, "M", false},
		{"C4 to Ab4 (Minor Sixth)", Note{0, 4, 0}, Note{5, 4, -1}, "m", false},
		{"C4 to B4 (Major Seventh)", Note{0, 4, 0}, Note{6, 4, 0}, "M", false},
		{"C4 to Bb4 (Minor Seventh)", Note{0, 4, 0}, Note{6, 4, -1}, "m", false},

		// Augmented Intervals
		{"C4 to C#4 (Augmented Unison)", Note{0, 4, 0}, Note{0, 4, 1}, "A", false},
		{"C4 to D#4 (Augmented Second)", Note{0, 4, 0}, Note{1, 4, 1}, "A", false},
		{"C4 to F#4 (Augmented Fourth)", Note{0, 4, 0}, Note{3, 4, 1}, "A", false},
		{"C4 to G#4 (Augmented Fifth)", Note{0, 4, 0}, Note{4, 4, 1}, "A", false},
		// Corrected: C4 to Cb4 is an Augmented Unison (1 semitone from 0-semitone perfect unison)
		{"C4 to Cb4 (Augmented Unison)", Note{0, 4, 0}, Note{0, 4, -1}, "A", false},

		// Diminished Intervals
		// Corrected: B4 to F4 is an Augmented Fourth
		{"C4 to Dbb4 (Diminished Second - theoretical)", Note{0, 4, 0}, Note{1, 4, -2}, "d", false},
		{"C4 to Fb4 (Diminished Fourth - theoretical)", Note{0, 4, 0}, Note{3, 4, -1}, "d", false},
		{"B4 to F4 (Augmented Fourth)", Note{6, 4, 0}, Note{3, 4, 0}, "A", false},
		{"C4 to Gb4 (Diminished Fifth)", Note{0, 4, 0}, Note{4, 4, -1}, "d", false},

		// Cross-octave intervals
		{"C4 to D5 (Major Ninth)", Note{0, 4, 0}, Note{1, 5, 0}, "M", false},
		{"C4 to G3 (Perfect Fourth Down)", Note{0, 4, 0}, Note{4, 3, 0}, "P", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateIntervalQuality(tt.note1, tt.note2)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateIntervalQuality() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("CalculateIntervalQuality() got = %v, expected %v", got, tt.expected)
			}
		})
	}
}

// Helper for testing noteToSemitones if needed for debugging
func TestNoteToSemitones(t *testing.T) {
	tests := []struct {
		name     string
		note     Note
		expected int
	}{
		{"C4", Note{0, 4, 0}, 48}, // C0 is 0, C1 is 12, C4 is 4*12 = 48
		{"C#4", Note{0, 4, 1}, 49},
		{"Db4", Note{1, 4, -1}, 49}, // D is 2 semitones from C. Db is 2-1 = 1. C4 is 48. So 48+1=49
		{"D4", Note{1, 4, 0}, 50},
		{"E4", Note{2, 4, 0}, 52},
		{"F4", Note{3, 4, 0}, 53},
		{"G4", Note{4, 4, 0}, 55},
		{"A4", Note{5, 4, 0}, 57},
		{"B4", Note{6, 4, 0}, 59},
		{"C5", Note{0, 5, 0}, 60},
		{"B3", Note{6, 3, 0}, 47}, // B is 11 semitones from C. 3*12 + 11 = 36 + 11 = 47
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := noteToSemitones(tt.note)
			if got != tt.expected {
				t.Errorf("noteToSemitones(%v) got = %d, expected %d", tt.note, got, tt.expected)
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
