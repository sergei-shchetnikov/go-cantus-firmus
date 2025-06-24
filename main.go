package main

import (
	"errors"
	"fmt"
)

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
// Alterations are not supported.
//
// Examples of valid input:
//   - "C4" (Middle C)
//   - "G3" (G below middle C)
//
// Returns:
//   - Note struct if parsing is successful
//   - error if the format is invalid (with specific reason)
func ParseNote(s string) (Note, error) {
	if len(s) < 2 {
		return Note{}, errors.New("invalid note format: string too short")
	}

	noteChar := s[0]
	octaveStr := s[1:]

	var step int
	switch noteChar {
	case 'C', 'c':
		step = 0
	case 'D', 'd':
		step = 1
	case 'E', 'e':
		step = 2
	case 'F', 'f':
		step = 3
	case 'G', 'g':
		step = 4
	case 'A', 'a':
		step = 5
	case 'B', 'b':
		step = 6
	default:
		return Note{}, fmt.Errorf("invalid note character: %c", noteChar)
	}

	var octave int
	_, err := fmt.Sscanf(octaveStr, "%d", &octave)
	if err != nil {
		return Note{}, fmt.Errorf("invalid octave: %v", err)
	}

	return Note{Step: step, Octave: octave}, nil
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

// CantusFirmus represents a melodic contour abstracted from rhythm, meter, key, or specific pitches.
// It captures only the sequence of diatonic intervals between consecutive notes, serving as the foundation
// for later elaboration into a complete melody by applying tonality, mode, and other musical parameters.
//
// Example: [third up, second down, second down] → "D4, F4, E4, D4" (if starting from D4).
type CantusFirmus []Interval
