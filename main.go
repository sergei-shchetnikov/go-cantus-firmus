package main

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
type Interval int
