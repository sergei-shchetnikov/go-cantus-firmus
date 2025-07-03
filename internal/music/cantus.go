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
