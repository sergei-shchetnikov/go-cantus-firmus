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

// PreparedLeaps checks if leaps are properly prepared by contrary motion according to counterpoint rules
func PreparedLeaps(intervals []int) bool {
	n := len(intervals)
	if n <= 1 {
		return true
	}

	last := intervals[n-1]
	absLast := abs(last)

	// Small steps don't need preparation
	if absLast <= 2 {
		return true
	}

	// Dispatch validation based on leap size
	switch absLast {
	case 3:
		return validateFourthLeap(intervals)
	case 4:
		return validateFifthLeap(intervals)
	case 5:
		return validateSixthLeap(intervals)
	default:
		return true
	}
}

// validateFourthLeap handles preparation for leaps of 3 or -3 (fourth)
func validateFourthLeap(intervals []int) bool {
	n := len(intervals)
	last := intervals[n-1]
	return n >= 2 && sign(intervals[n-2]) == -sign(last)
}

// validateFifthhLeap handles preparation for leaps of 4 or -4 (fifth)
func validateFifthLeap(intervals []int) bool {
	n := len(intervals)
	last := intervals[n-1]

	if n >= 2 {
		prev := intervals[n-2]
		if sign(prev) == -sign(last) && abs(prev) >= 2 {
			return true
		}
	}
	if n >= 3 {
		prev1 := intervals[n-2]
		prev2 := intervals[n-3]
		return sign(prev1) == -sign(last) && sign(prev2) == -sign(last)
	}
	return false
}

// validateSixthLeap handles preparation for leap of 5 (sixth)
func validateSixthLeap(intervals []int) bool {
	n := len(intervals)
	last := intervals[n-1]

	if last != 5 {
		return false
	}

	switch {
	case n >= 4:
		prev1 := intervals[n-2]
		prev2 := intervals[n-3]
		prev3 := intervals[n-4]
		if prev1 < 0 && prev2 < 0 && prev3 < 0 {
			return true
		}
		fallthrough
	case n >= 3:
		prev1 := intervals[n-2]
		prev2 := intervals[n-3]
		if prev1 < 0 && prev2 < 0 && (abs(prev1) >= 2 || abs(prev2) >= 2) {
			return true
		}
		fallthrough
	case n >= 2:
		prev := intervals[n-2]
		return prev < 0 && abs(prev) >= 3
	default:
		return false
	}
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// ValidateLeapResolution checks if all leaps (intervals with absolute value > 2)
// are properly resolved according to counterpoint rules.
// It works with partial slices during generation.
//
// Returns:
//   - true if all leaps are properly resolved or if the slice contains no leaps
//   - false if any leap resolution violates the rules
func ValidateLeapResolution(intervals []int) bool {
	n := len(intervals)
	if n <= 1 {
		return true
	}

	// Find all leaps (absolute value > 2)
	leapIndices := make([]int, 0)
	for i := 0; i < n; i++ {
		if abs(intervals[i]) > 2 {
			leapIndices = append(leapIndices, i)
		}
	}

	// No leaps found
	if len(leapIndices) == 0 {
		return true
	}

	// Check resolution for each leap
	for _, leapIndex := range leapIndices {
		// Skip if the leap is at the end (no resolution needed yet)
		if leapIndex == n-1 {
			continue
		}

		leap := intervals[leapIndex]
		absLeap := abs(leap)
		leapSlice := intervals[leapIndex:]

		// Dispatch validation based on leap size
		var resolved bool
		switch absLeap {
		case 3:
			resolved = validateFourthLeapResolution(leapSlice)
		case 4:
			resolved = validateFifthLeapResolution(leapSlice)
		case 5:
			resolved = validateSixthLeapResolution(leapSlice)
		default:
			resolved = true // Larger leaps are handled elsewhere or not considered
		}

		if !resolved {
			return false
		}
	}

	return true
}

// validateFourthLeapResolution handles resolution for leaps of 3 or -3 (fourth)
func validateFourthLeapResolution(intervals []int) bool {
	if len(intervals) < 2 {
		return true
	}
	return sign(intervals[0]) == -sign(intervals[1])
}

// validateFifthLeapResolution handles resolution for leaps of 4 or -4 (fifth)
func validateFifthLeapResolution(intervals []int) bool {
	n := len(intervals)
	if n < 2 {
		return true
	}

	leap := intervals[0]
	next1 := intervals[1]

	// Case with exactly one element after leap
	if n == 2 {
		return sign(leap) == -sign(next1)
	}

	// Case with at least two elements after leap
	next2 := intervals[2]
	return (sign(leap) == -sign(next1) && abs(next1) >= 2) ||
		(sign(leap) == -sign(next1) && sign(leap) == -sign(next2))
}

// validateSixthLeapResolution handles resolution for leap of 5 (sixth)
func validateSixthLeapResolution(intervals []int) bool {
	n := len(intervals)
	if n < 2 {
		return true
	}

	leap := intervals[0]
	if leap != 5 {
		return false
	}

	next1 := intervals[1]

	// Case with exactly one element after leap
	if n == 2 {
		return sign(leap) == -sign(next1)
	}

	next2 := intervals[2]

	// Case with exactly two elements after leap
	if n == 3 {
		return (next1 < 0 && abs(next1) >= 3) ||
			(sign(leap) == -sign(next1) && sign(leap) == -sign(next2))
	}

	// Case with at least three elements after leap
	next3 := intervals[3]
	return (next1 < 0 && abs(next1) >= 3) ||
		(next1 < 0 && next2 < 0 && (abs(next1)+abs(next2)) >= 3) ||
		(next1 < 0 && next2 < 0 && next3 < 0)
}

// NoTripleAlternatingNote checks that no note repeats three times in an alternating pattern (a, b, a, c, a).
// Works with partial slices during generation.
//
// Returns:
//   - false if the pattern is found (rule violated)
//   - true otherwise (rule satisfied)
func NoTripleAlternatingNote(intervals []int) bool {
	if len(intervals) < 4 {
		return true
	}

	// Compute partial sums (note heights relative to the starting note)
	partialSums := make([]int, len(intervals)+1)
	partialSums[0] = 0
	for i, interval := range intervals {
		partialSums[i+1] = partialSums[i] + interval
	}

	// Check for the pattern a, b, a, c, a
	for i := 0; i <= len(partialSums)-5; i++ {
		a := partialSums[i]

		// Check if the same note appears at positions i, i+2, and i+4
		if partialSums[i+2] == a && partialSums[i+4] == a {
			return false
		}

	}

	return true
}
