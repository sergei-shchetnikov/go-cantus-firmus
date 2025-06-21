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
