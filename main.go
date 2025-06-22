package main

import "fmt"

// Project: go-cantus-firmus
// Created: 2025-06-21

func main() {

}

// Mod7 returns the non-negative remainder of division of n by 7.
// The result will be in the range [0, 6] for any integer input.
func Mod7(n int) int {
	result := n % 7
	if result < 0 {
		result += 7
	}
	return result
}

// Note represents a musical note in diatonic notation.
//
// Fields:
//   - Step: diatonic step number (0 = C, 1 = D, ..., 6 = B)
//   - Octave: octave number (4 is the middle octave)
type Note struct {
	Step   int
	Octave int
}

// String returns the string representation of the note in the format "C4".
// The step numbers are mapped to diatonic note names (0=C, 1=D, ..., 6=B).
// Octave numbers follow the scientific pitch notation standard.
// Examples:
//   - Note{0, 4} → "C4" (Middle C)
//   - Note{6, 3} → "B3" (B below Middle C)
//   - Note{0, 5} → "C5" (C one octave above Middle C)
func (n Note) String() string {
	noteNames := []string{"C", "D", "E", "F", "G", "A", "B"}
	return fmt.Sprintf("%s%d", noteNames[n.Step], n.Octave)
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
