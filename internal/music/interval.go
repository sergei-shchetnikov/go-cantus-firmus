package music

import (
	"fmt"
	"math"
)

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

// noteToSemitones converts a Note to its absolute semitone value,
// assuming C0 as the reference point (0 semitones).
// This function helps in calculating the exact pitch distance between notes.
func noteToSemitones(n Note) int {
	// Semitone values for each step relative to C in an octave
	// C, D, E, F, G, A, B
	// 0, 2, 4, 5, 7, 9, 11
	stepSemitones := []int{0, 2, 4, 5, 7, 9, 11}

	// Calculate the base semitones from the step and octave
	// Each octave has 12 semitones
	semitones := n.Octave*12 + stepSemitones[n.Step]

	// Add alteration
	semitones += n.Alteration

	return semitones
}

// CalculateIntervalQuality determines the quality of the interval between two notes.
// It returns "P" for perfect, "A" for augmented, "M" for major, or "m" for minor.
// The order of notes (n1, n2) determines whether the interval is ascending or descending,
// but for quality determination, the absolute semitone difference is used first.
func CalculateIntervalQuality(n1, n2 Note) (string, error) {
	// Calculate the absolute semitone difference between the two notes
	semitoneDiffAbs := int(math.Abs(float64(noteToSemitones(n2) - noteToSemitones(n1))))

	// Calculate the diatonic numerical interval (e.g., 1st, 2nd, 3rd, etc.)
	// This accounts for octave differences and descending intervals correctly.
	// We need to determine the numerical interval based on the *span* of steps, not just Mod7.
	// Example: C4 to C5 is an 8th (octave), not a 1st.
	// C4 to G3: StepDiff (G-C) = 4. OctaveDiff (3-4) = -1. This is not simply Mod7.

	// First, normalize notes to the same octave for step calculation to avoid negative steps.
	// Or, calculate the total number of steps, including octave changes.

	// Calculate the total step distance (including octave shifts)
	// This effectively treats C, D, E, F, G, A, B as a continuous scale.
	// C4 = 0, D4 = 1, E4 = 2, F4 = 3, G4 = 4, A4 = 5, B4 = 6
	// C5 = 7, D5 = 8, etc.
	n1TotalStep := n1.Step + n1.Octave*7
	n2TotalStep := n2.Step + n2.Octave*7

	// The raw diatonic step difference
	rawStepDiff := int(math.Abs(float64(n2TotalStep - n1TotalStep)))

	// The numerical interval is rawStepDiff + 1 (unison is 1, second is 2, etc.)
	numericalInterval := rawStepDiff + 1

	// Define the semitone counts for standard perfect and major/minor intervals
	// This map stores the semitone values for perfect (P) and major (M) intervals
	// Key: numerical interval (1 for unison, 2 for second, etc.)
	// Value: [perfect_semitones, major_semitones] (perfect for P intervals, major for M/m intervals)
	// These values are for ascending intervals.
	standardSemitones := map[int][2]int{
		1:  {0, 0},  // Unison (P1)
		2:  {0, 2},  // Second (M2)
		3:  {0, 4},  // Third (M3)
		4:  {5, 0},  // Fourth (P4)
		5:  {7, 0},  // Fifth (P5)
		6:  {0, 9},  // Sixth (M6)
		7:  {0, 11}, // Seventh (M7)
		8:  {12, 0}, // Octave (P8)
		9:  {0, 14}, // Major Ninth (Octave + Major Second)
		10: {0, 16}, // Major Tenth (Octave + Major Third)
		11: {17, 0}, // Perfect Eleventh (Octave + Perfect Fourth)
		12: {19, 0}, // Perfect Twelfth (Octave + Perfect Fifth)
		13: {0, 21}, // Major Thirteenth (Octave + Major Sixth)
		14: {0, 23}, // Major Fourteenth (Octave + Major Seventh)
		15: {24, 0}, // Perfect Fifteenth (Double Octave)
	}

	// Get the expected semitones for the numerical interval
	expected, ok := standardSemitones[numericalInterval]
	if !ok {
		// For very large intervals, we can calculate based on octave equivalency
		// or return an error if not explicitly defined.
		// For now, let's allow larger intervals by calculating their octave equivalent.
		// A 9th is a 2nd + octave, a 10th is a 3rd + octave etc.
		// The quality is derived from the simple interval.

		// If the numerical interval is not in the map, assume it's a compound interval.
		// Calculate the equivalent simple interval and its octave offset.
		simpleNumericalInterval := Mod7(numericalInterval-1) + 1 // (numericalInterval - 1) is the number of steps
		octaveOffset := (numericalInterval - 1) / 7

		if simpleNumericalInterval == 0 { // This can happen if numericalInterval is 7, 14 etc.
			simpleNumericalInterval = 7 // Treat as 7th if numericalInterval is multiple of 7
			if numericalInterval == 1 {
				simpleNumericalInterval = 1
			} // Unison
			if numericalInterval == 8 {
				simpleNumericalInterval = 8
			} // Octave
		}

		expected, ok = standardSemitones[simpleNumericalInterval]
		if !ok {
			return "", fmt.Errorf("unsupported numerical interval beyond typical compound range: %d", numericalInterval)
		}

		// Adjust expected semitones for compound intervals
		if simpleNumericalInterval != 1 && simpleNumericalInterval != 8 { // Not unison or octave
			expected[0] += octaveOffset * 12
			expected[1] += octaveOffset * 12
		} else { // For unison/octave, only add 12 semitones per octave
			expected[0] = expected[0] + octaveOffset*12
			expected[1] = expected[1] + octaveOffset*12 // This will be 0 for perfect, so it's fine.
		}

	}

	// Determine the quality based on semitone difference
	// Perfect intervals (Unison, 4th, 5th, Octave, and their compounds)
	if numericalInterval == 1 || numericalInterval == 4 || numericalInterval == 5 || numericalInterval == 8 ||
		numericalInterval == 11 || numericalInterval == 12 || numericalInterval == 15 { // Include compound perfect intervals

		// Also check if the simple interval is perfect
		isPerfectFamily := false
		simpleNum := Mod7(numericalInterval-1) + 1
		if simpleNum == 1 || simpleNum == 4 || simpleNum == 5 { // Unison, 4th, 5th
			isPerfectFamily = true
		}
		if numericalInterval == 8 || numericalInterval == 15 { // Octave, Double Octave
			isPerfectFamily = true
		}

		if isPerfectFamily {
			if semitoneDiffAbs == expected[0] {
				return "P", nil // Perfect
			} else if semitoneDiffAbs > expected[0] {
				return "A", nil // Augmented
			} else {
				return "d", nil // Diminished
			}
		}
	}

	// Major/Minor intervals (2nd, 3rd, 6th, 7th, and their compounds)
	// For these, we use expected[1] (Major semitones) as the reference
	if semitoneDiffAbs == expected[1] {
		return "M", nil // Major
	} else if semitoneDiffAbs == expected[1]-1 {
		return "m", nil // Minor
	} else if semitoneDiffAbs > expected[1] {
		return "A", nil // Augmented
	} else {
		// If semitoneDiffAbs is less than expected[1]-1 (i.e., more than 1 semitone less than major)
		// it's diminished (e.g., diminished 3rd is 2 semitones, Major 3rd is 4, Minor 3rd is 3)
		return "d", nil // Diminished
	}
}
