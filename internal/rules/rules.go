// Package rules provides validation functions for contrapuntal rules
// in cantus firmus. The package operates on abstract representations
// of cantus firmus as sequences of diatonic intervals between consecutive notes.
// The package implements contrapuntal rules as validation functions that check
// whether a given sequence of intervals violates specific composition guidelines.
// These functions are designed to be integrated with cantus firmus generation
// algorithms (like in package cantusgen) to ensure generated melodies adhere
// to traditional counterpoint rules.
package rules

// ValidationFunc defines the type for a validation function.
type ValidationFunc func(s []int) bool

// AllRules checks a slice against a given set of validation functions.
// It returns false if any function returns false, true otherwise.
func AllRules(s []int, validators []ValidationFunc) bool {
	for _, validate := range validators {
		if !validate(s) {
			return false
		}
	}
	return true
}

// NoBeginWithFive checks that the interval sequence doesn't start with 5.
// Returns false if the first interval is 5, true otherwise.
func NoBeginWithFive(intervals []int) bool {
	if len(intervals) > 0 && intervals[0] == 5 {
		return false
	}
	return true
}

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

		if s1 == s2 && s2 == s3 && s3 == s4 && s4 == s5 {
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

// NoExcessiveNoteRepetition checks that no single note (as represented by cumulative interval sum)
// appears more than 4 times in the cantus firmus. Works with partial slices during generation.
//
// Returns:
//   - false if any note repeats more than 4 times (rule violated)
//   - true otherwise (rule satisfied)
func NoExcessiveNoteRepetition(intervals []int) bool {
	if len(intervals) == 0 {
		return true
	}

	sumCounts := make(map[int]int)
	currentSum := 0
	sumCounts[currentSum] = 1 // Count the starting note

	for _, interval := range intervals {
		currentSum += interval
		sumCounts[currentSum]++

		if sumCounts[currentSum] > 4 {
			return false
		}
	}

	return true
}

// NoRangeExceedsDecima checks that the range of the cantus firmus (difference between
// highest and lowest notes) does not exceed a decima (9 in interval notation).
// Works with partial slices during generation.
//
// Returns:
//   - false if the range exceeds 9 (rule violated)
//   - true otherwise (rule satisfied)
func NoRangeExceedsDecima(intervals []int) bool {
	if len(intervals) == 0 {
		return true
	}

	currentSum := 0
	minSum := 0
	maxSum := 0

	for _, interval := range intervals {
		currentSum += interval
		if currentSum < minSum {
			minSum = currentSum
		}
		if currentSum > maxSum {
			maxSum = currentSum
		}

		if maxSum-minSum > 9 {
			return false
		}
	}

	return true
}

// NoRepeatingPatterns checks that the cantus firmus doesn't contain repeating pitch patterns
// by examining the sequence of partial sums (note heights relative to the starting note).
// Detects patterns like ..., a, b, a, b, ... or ..., a, b, c, a, b, c, ...
// Works with partial slices during generation.
//
// Returns:
//   - false if any repeating pitch pattern is found (rule violated)
//   - true otherwise (rule satisfied)
func NoRepeatingPatterns(intervals []int) bool {
	if len(intervals) < 3 {
		return true
	}

	partialSums := make([]int, len(intervals)+1)
	partialSums[0] = 0
	for i, interval := range intervals {
		partialSums[i+1] = partialSums[i] + interval
	}

	n := len(partialSums)

	for i := 0; i <= n-4; i++ {
		a, b := partialSums[i], partialSums[i+1]
		if partialSums[i+2] == a && partialSums[i+3] == b {
			return false
		}
	}

	for i := 0; i <= n-6; i++ {
		a, b, c := partialSums[i], partialSums[i+1], partialSums[i+2]
		if partialSums[i+3] == a && partialSums[i+4] == b && partialSums[i+5] == c {
			return false
		}
	}

	return true
}
