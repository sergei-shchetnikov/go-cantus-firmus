package music

import (
	"errors"
	"fmt"
	"go-cantus-firmus/internal/utils"
)

// Note represents a musical note
//
// Fields:
//   - Step: diatonic step number (0 = C, 1 = D, ..., 6 = B)
//   - Octave: octave number (4 is the middle octave)
//   - Alteration: accidental for the note (-1 = flat, 0 = natural, 1 = sharp)
type Note struct {
	Step       int
	Octave     int
	Alteration int
}

// String returns the string representation of the note in standard musical notation.
// The step numbers are mapped to diatonic note names (0=C, 1=D, ..., 6=B).
// Octave numbers follow scientific pitch notation.
// Alteration affects the note name:
//
//	-1 → flat (represented as "b")
//	 0 → natural (no symbol)
//	 1 → sharp (represented as "#")
//
// Examples:
//   - Note{0, 4, 0}  → "C4" (Middle C)
//   - Note{0, 4, 1}  → "C#4" (C sharp)
//   - Note{1, 4, -1} → "Db4" (D flat)
//   - Note{6, 3, 0}  → "B3" (B below Middle C)
func (n Note) String() string {
	noteNames := []string{"C", "D", "E", "F", "G", "A", "B"}
	alterationSymbol := ""
	switch n.Alteration {
	case 1:
		alterationSymbol = "#"
	case -1:
		alterationSymbol = "b"
	}
	return fmt.Sprintf("%s%s%d", noteNames[n.Step], alterationSymbol, n.Octave)
}

// ParseNote parses a string representation of a musical note into a Note struct.
//
// Examples of valid input:
//   - "C4" (Middle C)
//   - "C#4" (C sharp)
//   - "Db4" (D flat)
//   - "G3" (G below middle C)
//   - "Fb5" (F flat)
//   - "B#2" (B sharp)
//
// Returns:
//   - Note struct if parsing is successful
//   - error if the format is invalid (with specific reason)
func ParseNote(s string) (Note, error) {
	if len(s) < 2 {
		return Note{}, errors.New("invalid note format: string too short")
	}

	// Extract note character (first character)
	noteChar := s[0]
	if noteChar < 'A' || noteChar > 'G' && noteChar < 'a' || noteChar > 'g' {
		return Note{}, fmt.Errorf("invalid note character: %c", noteChar)
	}

	// Convert to uppercase for consistency
	if noteChar >= 'a' && noteChar <= 'g' {
		noteChar -= 'a' - 'A'
	}

	// Determine step
	var step int
	switch noteChar {
	case 'C':
		step = 0
	case 'D':
		step = 1
	case 'E':
		step = 2
	case 'F':
		step = 3
	case 'G':
		step = 4
	case 'A':
		step = 5
	case 'B':
		step = 6
	}

	// Check for alteration
	alteration := 0
	rest := s[1:]
	if len(rest) > 0 {
		switch rest[0] {
		case '#':
			alteration = 1
			rest = rest[1:]
		case 'b':
			alteration = -1
			rest = rest[1:]
		}
	}

	// Parse octave
	if len(rest) == 0 {
		return Note{}, errors.New("invalid note format: missing octave")
	}

	var octave int
	_, err := fmt.Sscanf(rest, "%d", &octave)
	if err != nil {
		return Note{}, fmt.Errorf("invalid octave: %v", err)
	}

	return Note{
		Step:       step,
		Octave:     octave,
		Alteration: alteration,
	}, nil
}

// Transpose transposes a note by the given interval.
func Transpose(n Note, i Interval) Note {
	stepDelta := int(i)
	totalSteps := n.Step + stepDelta
	newStep := Mod7(totalSteps)

	octaveDelta := totalSteps / 7
	if totalSteps < 0 && newStep != 0 {
		octaveDelta -= 1
	}

	return Note{
		Step:   newStep,
		Octave: n.Octave + octaveDelta,
	}
}

// IsLeap determines whether the interval between two notes is a leap (larger than a second).
// Returns true if the interval is larger than a second (i.e., a third or greater).
func IsLeap(n1, n2 Note) bool {
	// Calculate the total step count for each note including octaves
	n1TotalStep := n1.Step + n1.Octave*7
	n2TotalStep := n2.Step + n2.Octave*7

	// Calculate the absolute difference in steps
	stepDiff := utils.Abs((n2TotalStep - n1TotalStep))

	// A leap is any interval larger than a second (step difference > 1)
	return stepDiff > 1
}

// Semitones returns the number of semitones from C0 (0 semitones)
// This can be used to compare note pitches
func (n Note) Semitones() int {
	// Diatonic steps from C to B (C=0, D=2, E=4, F=5, G=7, A=9, B=11)
	diatonicSemitones := []int{0, 2, 4, 5, 7, 9, 11}

	// Calculate semitones from diatonic step and alteration
	semitones := diatonicSemitones[n.Step] + n.Alteration

	// Add octave (12 semitones per octave)
	return semitones + n.Octave*12
}

// Less returns true if this note is lower in pitch than the other note
func (n Note) Less(other Note) bool {
	return n.Semitones() < other.Semitones()
}

// Greater returns true if this note is higher in pitch than the other note
func (n Note) Greater(other Note) bool {
	return n.Semitones() > other.Semitones()
}

// EqualPitch returns true if the notes have the same pitch (same semitones value)
func (n Note) EqualPitch(other Note) bool {
	return n.Semitones() == other.Semitones()
}
