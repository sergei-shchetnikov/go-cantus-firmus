package music

import (
	"fmt"
)

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
//   - Adds sharps to G and F when the configuration ..., A, G, F, G, A, ... appears
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
		// Configuration ..., A, G, F, G, A, ... (requires at least 5 notes)
		if i >= 2 && i < len(adjusted)-2 {
			prevPrev := adjusted[i-2]
			nextNext := adjusted[i+2]

			if prevPrev.Step == 5 && // A
				prev.Step == 4 && // G
				current.Step == 3 && // F
				next.Step == 4 && // G
				nextNext.Step == 5 { // A

				if prev.Alteration == 0 {
					adjusted[i-1].Alteration = 1 // G
				}
				if current.Alteration == 0 {
					adjusted[i].Alteration = 1 // F
				}
				if next.Alteration == 0 {
					adjusted[i+1].Alteration = 1 // G
				}
			}
		}
	}

	return adjusted
}

// IsNoteSurroundedByLinearMotion determines if the note at the given index `i` in a Realization is surrounded by linear motion.
// This means that notes at `i-1`, `i`, and `i+1` must form a consecutive ascending or descending sequence,
// where the interval between successive notes is a second (stepwise).
//
// Returns:
//   - true if the notes are in linear (stepwise) motion
//   - false otherwise (including when the index `i` is out of bounds for checking surrounding notes)
func IsNoteSurroundedByLinearMotion(r Realization, i int) bool {
	if i <= 0 || i >= len(r)-1 {
		return false // Not enough notes to check for surrounding linear motion
	}

	nPrev := r[i-1]
	nCurrent := r[i]
	nNext := r[i+1]

	// Calculate the total step count for each note including octaves
	nPrevTotalStep := nPrev.Step + nPrev.Octave*7
	nCurrentTotalStep := nCurrent.Step + nCurrent.Octave*7
	nNextTotalStep := nNext.Step + nNext.Octave*7

	// Check if the steps are consecutive and in the same direction
	// Case 1: Ascending linear motion (e.g., C4, D4, E4 or B3, C4, D4)
	if nCurrentTotalStep == nPrevTotalStep+1 && nNextTotalStep == nCurrentTotalStep+1 {
		return true
	}

	// Case 2: Descending linear motion (e.g., E4, D4, C4 or D4, C4, B3)
	if nCurrentTotalStep == nPrevTotalStep-1 && nNextTotalStep == nCurrentTotalStep-1 {
		return true
	}

	return false
}
