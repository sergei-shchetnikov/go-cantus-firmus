package rules

// NoFiveOfSameSign checks that there are no five consecutive numbers
// with the same sign (positive or negative) in the currentSlice.
// In a musical context, this rule helps prevent excessive or monotonous
// stepwise motion or leaps in a single direction (e.g., always ascending or always descending)
// for an extended period, which can lead to a less engaging melodic line.
// Returns false if five consecutive numbers of the same sign are found, otherwise true.
// Works with incomplete slices.
func NoFiveOfSameSign(currentSlice []int) bool {
	n := len(currentSlice)
	if n < 5 {
		return true
	}

	for i := 0; i <= n-5; i++ {
		s1 := sign(currentSlice[i])
		s2 := sign(currentSlice[i+1])
		s3 := sign(currentSlice[i+2])
		s4 := sign(currentSlice[i+3])
		s5 := sign(currentSlice[i+4])

		if s1 == s2 && s2 == s3 && s3 == s4 && s4 == s5 { // Updated condition
			return false
		}
	}
	return true
}

// sign returns the sign of a number:
//
//	1 for positive numbers
//
// -1 for negative numbers
func sign(x int) int {
	if x > 0 {
		return 1
	}
	// Since 0 is not expected, any non-positive number must be negative
	return -1
}
