package music

import (
	"fmt"
)

// Mod7 returns the non-negative remainder of division of n by 7.
// The result will be in the range [0, 6] for any integer input.
func Mod7(n int) int {
	result := n % 7
	if result < 0 {
		result += 7
	}
	return result
}

// Interval represents a musical interval using Taneyev's digital notation system.
//
// This practical numbering system, derived from Sergei Taneyev's "Convertible Counterpoint",
// represents intervals as signed integers where:
//   - 0: Unison
//   - Positive integers: Ascending intervals
//   - 1: Second up
//   - 2: Third up
//   - ...
//   - Negative integers: Descending intervals
//   - -1: Second down
//   - -2: Third down
//   - ...
//
// Note: The interval value represents the scale degree difference and is abstracted from
// the absolute pitch magnitude (e.g., minor/major seconds are both represented as 1).
type Interval int

// String returns the human-readable representation of the interval.
func (i Interval) String() string {
	absVal := int(i)
	if absVal < 0 {
		absVal = -absVal
	}

	direction := "up"
	if int(i) < 0 {
		direction = "down"
	}

	switch absVal {
	case 0:
		return "unison"
	case 1:
		return fmt.Sprintf("second %s", direction)
	case 2:
		return fmt.Sprintf("third %s", direction)
	case 3:
		return fmt.Sprintf("fourth %s", direction)
	case 4:
		return fmt.Sprintf("fifth %s", direction)
	case 5:
		return fmt.Sprintf("sixth %s", direction)
	case 6:
		return fmt.Sprintf("seventh %s", direction)
	case 7:
		return fmt.Sprintf("octave %s", direction)
	}

	intervalNum := absVal + 1 // Convert to musical interval number
	var suffix string

	switch intervalNum % 10 {
	case 1:
		suffix = "st"
	case 2:
		suffix = "nd"
	case 3:
		suffix = "rd"
	default:
		suffix = "th"
	}

	if intervalNum >= 11 && intervalNum <= 13 {
		suffix = "th"
	}

	return fmt.Sprintf("%d%s %s", intervalNum, suffix, direction)
}

// CantusFirmus represents a melodic contour abstracted from rhythm, meter, key, or specific pitches.
// It captures only the sequence of diatonic intervals between consecutive notes, serving as the foundation
// for later elaboration into a complete melody by applying tonality, mode, and other musical parameters.
//
// Example: [third up, second down, second down] → "D4, F4, E4, D4" (if starting from D4).
type CantusFirmus []Interval

// Realize generates a concrete musical realization of the CantusFirmus in the specified mode.
// The first note will be the tonic of the mode (C for Major, D for Dorian, E for Phrygian,
// F for Lydian, G for Mixolydian, A for Minor, B for Locrian),
// and subsequent notes will follow the intervals of the CantusFirmus.
func (cf CantusFirmus) Realize(mode string) (Realization, error) {
	var startingNote Note
	switch mode {
	case "Major":
		startingNote = Note{Step: 0, Octave: 4} // C4
	case "Dorian":
		startingNote = Note{Step: 1, Octave: 4} // D4
	case "Phrygian":
		startingNote = Note{Step: 2, Octave: 4} // E4
	case "Lydian":
		startingNote = Note{Step: 3, Octave: 4} // F4
	case "Mixolydian":
		startingNote = Note{Step: 4, Octave: 4} // G4
	case "Minor":
		startingNote = Note{Step: 5, Octave: 4} // A4
	case "Locrian":
		startingNote = Note{Step: 6, Octave: 4} // B4
	default:
		return nil, fmt.Errorf("unknown mode: %s", mode)
	}

	realization := Realization{startingNote}

	currentNote := startingNote
	for _, interval := range cf {
		currentNote = Transpose(currentNote, interval)
		realization = append(realization, currentNote)
	}

	// Apply alteration rules for minor mode
	if mode == "Minor" {
		realization = adjustMinorAlterations(realization)
	}

	return realization, nil
}

// Realization represents a concrete musical realization of a CantusFirmus as a sequence of notes.
// It transforms the abstract interval sequence of a CantusFirmus into actual pitches,
// preserving the melodic contour while making the pitches explicit.
type Realization []Note

// adjustMinorAlterations adds necessary alteration marks to a Realization in minor mode.
//
// Rules:
//   - Adds sharp to G when the configuration ..., A, G, A, ... appears
//   - Adds sharps to both F and G when the configuration ..., F, G, A, ... appears
//   - In all other cases, no sharps are added
func adjustMinorAlterations(realization Realization) Realization {
	if len(realization) < 3 {
		return realization // Not enough notes to analyze configurations
	}

	adjusted := make(Realization, len(realization))
	copy(adjusted, realization)

	for i := 1; i < len(adjusted)-1; i++ {
		prev := adjusted[i-1]
		current := adjusted[i]
		next := adjusted[i+1]

		// Configuration ..., A, G, A, ...
		if prev.Step == 5 && current.Step == 4 && next.Step == 5 {
			if current.Alteration == 0 { // Only add sharp if the note hasn't been altered yet
				adjusted[i].Alteration = 1
			}
		}

		// Configuration ..., F, G, A, ...
		if prev.Step == 3 && current.Step == 4 && next.Step == 5 {
			if prev.Alteration == 0 {
				adjusted[i-1].Alteration = 1
			}
			if current.Alteration == 0 {
				adjusted[i].Alteration = 1
			}
		}
	}

	return adjusted
}
