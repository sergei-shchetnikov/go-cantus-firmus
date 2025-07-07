// Package rules provides validation functions for contrapuntal rules
// in cantus firmus. The package operates on abstract representations
// of cantus firmus as sequences of diatonic intervals between consecutive notes.
// The package implements contrapuntal rules as validation functions that check
// whether a given sequence of intervals violates specific composition guidelines.
// These functions are designed to be integrated with cantus firmus generation
// algorithms (like in package cantusgen) to ensure generated melodies adhere
// to traditional counterpoint rules.
package rules

import "go-cantus-firmus/internal/utils"

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

// LimitDirectionalMotion enforces contrapuntal rules regarding consecutive melodic motion:
// 1. Prohibits more than four consecutive intervals in the same direction (ascending/descending)
// 2. Restricts the cumulative melodic span in one direction to a sixth
func LimitDirectionalMotion(currentSlice []int) bool {
	n := len(currentSlice)
	if n == 0 {
		return true
	}

	currentSign := sign(currentSlice[0])
	currentSum := currentSlice[0]
	count := 1

	for i := 1; i < n; i++ {
		s := sign(currentSlice[i])

		if s == currentSign {
			count++
			currentSum += currentSlice[i]

			// Check for five consecutive same-sign numbers
			if count >= 5 {
				return false
			}

			// Check if absolute sum exceeds 5
			if utils.Abs(currentSum) > 5 {
				return false
			}
		} else {
			// Reset for new sign
			currentSign = s
			currentSum = currentSlice[i]
			count = 1
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
// appears more than 3 times in the cantus firmus. Works with partial slices during generation.
//
// Returns:
//   - false if any note repeats more than 3 times (rule violated)
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

		if sumCounts[currentSum] > 3 {
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
// Detects the following repeating patterns:
//   - Immediate repetitions:
//   - 2-note patterns: ..., a, b, a, b, ...
//   - 3-note patterns: ..., a, b, c, a, b, c, ...
//   - Patterns with separators between repetitions:
//   - 3-note patterns with 1 separator: ..., a, b, c, X, a, b, c, ...
//   - 3-note patterns with 2 separators: ..., a, b, c, X, Y, a, b, c, ...
//   - 3-note patterns with 3 separators: ..., a, b, c, X, Y, Z, a, b, c, ...
//
// where X, Y, Z can be any single pitch value.
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

	// Check for 2-note patterns (a,b,a,b)
	for i := 0; i <= n-4; i++ {
		a, b := partialSums[i], partialSums[i+1]
		if partialSums[i+2] == a && partialSums[i+3] == b {
			return false
		}
	}

	// Check for 3-note patterns (a,b,c,a,b,c)
	for i := 0; i <= n-6; i++ {
		a, b, c := partialSums[i], partialSums[i+1], partialSums[i+2]
		if partialSums[i+3] == a && partialSums[i+4] == b && partialSums[i+5] == c {
			return false
		}
	}

	// Check for 3-note patterns with 1 separator (a,b,c,x,a,b,c)
	for i := 0; i <= n-7; i++ {
		a, b, c := partialSums[i], partialSums[i+1], partialSums[i+2]
		if partialSums[i+4] == a && partialSums[i+5] == b && partialSums[i+6] == c {
			return false
		}
	}

	// Check for 3-note patterns with 2 separators (a,b,c,x,y,a,b,c)
	for i := 0; i <= n-8; i++ {
		a, b, c := partialSums[i], partialSums[i+1], partialSums[i+2]
		if partialSums[i+5] == a && partialSums[i+6] == b && partialSums[i+7] == c {
			return false
		}
	}

	// Check for 3-note patterns with 3 separators (a,b,c,x,y,z,a,b,c)
	for i := 0; i <= n-9; i++ {
		a, b, c := partialSums[i], partialSums[i+1], partialSums[i+2]
		if partialSums[i+6] == a && partialSums[i+7] == b && partialSums[i+8] == c {
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
	absLast := utils.Abs(last)

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
		if sign(prev) == -sign(last) && utils.Abs(prev) >= 2 {
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
		if prev1 < 0 && prev2 < 0 && (utils.Abs(prev1) >= 2 || utils.Abs(prev2) >= 2) {
			return true
		}
		fallthrough
	case n >= 2:
		prev := intervals[n-2]
		return prev < 0 && utils.Abs(prev) >= 3
	default:
		return false
	}
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
		if utils.Abs(intervals[i]) > 2 {
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
		absLeap := utils.Abs(leap)
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
	return (sign(leap) == -sign(next1) && utils.Abs(next1) >= 2) ||
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
		return (next1 < 0 && utils.Abs(next1) >= 3) ||
			(sign(leap) == -sign(next1) && sign(leap) == -sign(next2))
	}

	// Case with at least three elements after leap
	next3 := intervals[3]
	return (next1 < 0 && utils.Abs(next1) >= 3) ||
		(next1 < 0 && next2 < 0 && (utils.Abs(next1)+utils.Abs(next2)) >= 3) ||
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

// NoNoteRepetitionAfterLeap checks that there is no immediate note repetition after
// two consecutive leaps of equal magnitude in opposite directions.
// Returns false if two consecutive intervals:
//   - have absolute value > 1 (leaps)
//   - have opposite signs
//   - have equal absolute values
func NoNoteRepetitionAfterLeap(intervals []int) bool {
	if len(intervals) < 2 {
		return true
	}

	for i := 0; i < len(intervals)-1; i++ {
		current := intervals[i]
		next := intervals[i+1]

		if utils.Abs(current) > 1 &&
			utils.Abs(next) > 1 &&
			sign(current) == -sign(next) &&
			utils.Abs(current) == utils.Abs(next) {
			return false
		}
	}

	return true
}

// NoRepeatingExtremes checks that there are no two identical adjacent peaks or valleys
// in the cantus firmus. It builds a sequence of partial sums (representing note heights),
// then identifies local extrema (excluding the first and last notes), and checks for
// the pattern ..., a, b, a, ... in the extrema sequence.
//
// Returns:
//   - false if the pattern is found (rule violated)
//   - true otherwise (rule satisfied)
func NoRepeatingExtremes(intervals []int) bool {
	if len(intervals) < 3 {
		return true
	}

	// Build partial sums (note heights)
	partialSums := make([]int, len(intervals)+1)
	partialSums[0] = 0
	for i, interval := range intervals {
		partialSums[i+1] = partialSums[i] + interval
	}

	// Find all local extrema (excluding first and last notes)
	extrema := make([]int, 0)
	for i := 1; i < len(partialSums)-1; i++ {
		prev := partialSums[i-1]
		current := partialSums[i]
		next := partialSums[i+1]

		// Check for peak (current > neighbors) or valley (current < neighbors)
		if (current > prev && current > next) || (current < prev && current < next) {
			extrema = append(extrema, current)
		}
	}

	// Check for the pattern a, b, a in extrema
	for i := 0; i < len(extrema)-2; i++ {
		if extrema[i] == extrema[i+2] {
			return false
		}
	}

	return true
}

// AvoidSeventhBetweenExtrema checks that no adjacent extrema (peaks/valleys, including first/last notes)
// in the cantus firmus are separated by a seventh (interval of 6).
// Returns:
//   - false if any adjacent extrema differ by a seventh (rule violated)
//   - true otherwise (rule satisfied)
func AvoidSeventhBetweenExtrema(intervals []int) bool {
	if len(intervals) < 1 {
		return true
	}

	// Build partial sums (note heights)
	partialSums := make([]int, len(intervals)+1)
	partialSums[0] = 0
	for i, interval := range intervals {
		partialSums[i+1] = partialSums[i] + interval
	}

	// Find all local extrema, including first and last notes
	extrema := make([]int, 0)
	extrema = append(extrema, partialSums[0]) // Include first note

	for i := 1; i < len(partialSums)-1; i++ {
		prev := partialSums[i-1]
		current := partialSums[i]
		next := partialSums[i+1]

		// Check for peak or valley
		if (current > prev && current > next) || (current < prev && current < next) {
			extrema = append(extrema, current)
		}
	}

	extrema = append(extrema, partialSums[len(partialSums)-1]) // Include last note

	for i := 0; i < len(extrema)-1; i++ {
		if utils.Abs(extrema[i]-extrema[i+1]) == 6 {
			return false
		}
	}

	return true
}

// MinDirectionChanges checks that the melody changes direction (ascending/descending)
// at least twice in the complete interval sequence.
// Returns:
//   - false if there are less than 2 direction changes
//   - true otherwise
func MinDirectionChanges(intervals []int) bool {
	if len(intervals) < 3 {
		return false // Need at least 3 intervals to have 2 direction changes
	}

	directionChanges := 0
	prevSign := sign(intervals[0])

	for i := 1; i < len(intervals); i++ {
		currentSign := sign(intervals[i])
		if currentSign != prevSign {
			directionChanges++
			prevSign = currentSign
		}
	}

	return directionChanges >= 2
}

// ValidateClimax checks the climax rules for the cantus firmus:
// - If all heights are >= 0 (relative to starting note), there should be exactly one maximum
// - If all heights are <= 0, there should be exactly one minimum
// - If heights are both positive and negative, there should be exactly one maximum and one minimum
func ValidateClimax(intervals []int) bool {
	if len(intervals) == 0 {
		return true
	}

	// Build partial sums (note heights relative to starting note)
	partialSums := make([]int, len(intervals)+1)
	partialSums[0] = 0
	for i, interval := range intervals {
		partialSums[i+1] = partialSums[i] + interval
	}

	allPositive := true
	allNegative := true
	for _, sum := range partialSums {
		if sum < 0 {
			allPositive = false
		}
		if sum > 0 {
			allNegative = false
		}
	}

	if allPositive {
		return countMaxima(partialSums) == 1
	}

	if allNegative {
		return countMinima(partialSums) == 1
	}

	return countMaxima(partialSums) == 1 && countMinima(partialSums) == 1
}

// countMaxima counts how many times the maximum value appears in the slice
func countMaxima(sums []int) int {
	if len(sums) == 0 {
		return 0
	}

	maxVal := sums[0]
	count := 0

	for _, val := range sums {
		if val > maxVal {
			maxVal = val
			count = 1
		} else if val == maxVal {
			count++
		}
	}

	return count
}

// countMinima counts how many times the minimum value appears in the slice
func countMinima(sums []int) int {
	if len(sums) == 0 {
		return 0
	}

	minVal := sums[0]
	count := 0

	for _, val := range sums {
		if val < minVal {
			minVal = val
			count = 1
		} else if val == minVal {
			count++
		}
	}

	return count
}

// NoSequences prohibits any repeating patterns in the cantus firmus.
// The logic is divided into three parts:
// a) If the interval slice contains a pattern ..., a, b, a, b, a, ... where a != b, return false immediately
// b) If there are consecutive three-element patterns separated by one element (..., a, b, c, d, a, b, c, ...), return false
// c) Search for other three-element patterns containing leaps and check for their repetition
func NoSequences(intervals []int) bool {
	// Part a) Check for alternating a, b, a, b, a pattern
	if hasAlternatingPattern(intervals) {
		return false
	}

	// Part b) Check for consecutive patterns with one-element separation
	if hasConsecutivePatterns(intervals) {
		return false
	}

	// Part c) Check for repeating leap patterns
	return !hasRepeatingLeapPatterns(intervals)
}

// hasAlternatingPattern checks for the presence of a, b, a, b, a pattern where a != b
func hasAlternatingPattern(intervals []int) bool {
	if len(intervals) < 5 {
		return false
	}

	for i := 0; i <= len(intervals)-5; i++ {
		a := intervals[i]
		b := intervals[i+1]
		if a == b {
			continue // a and b must be different
		}

		// Check for a, b, a, b, a pattern
		if intervals[i+2] == a && intervals[i+3] == b && intervals[i+4] == a {
			return true
		}
	}

	return false
}

// hasConsecutivePatterns checks for consecutive three-element patterns separated by one element
func hasConsecutivePatterns(intervals []int) bool {
	if len(intervals) < 7 {
		return false
	}

	for i := 0; i <= len(intervals)-7; i++ {
		a1 := intervals[i]
		b1 := intervals[i+1]
		c1 := intervals[i+2]
		_ = intervals[i+3] // separator element (any value)
		a2 := intervals[i+4]
		b2 := intervals[i+5]
		c2 := intervals[i+6]

		// Check if the two triplets match with one element in between
		if a1 == a2 && b1 == b2 && c1 == c2 {
			return true
		}
	}

	return false
}

// hasRepeatingLeapPatterns checks for repeating patterns containing leaps
func hasRepeatingLeapPatterns(intervals []int) bool {
	if len(intervals) < 3 {
		return false
	}

	leaps := map[int]bool{-4: true, -3: true, -2: true, 2: true, 3: true, 4: true, 5: true}

	// Collect all potential patterns (triplets) with at least one leap and not all equal
	patterns := make([][3]int, 0)
	patternIndices := make([]int, 0) // Store starting indices of patterns

	for i := 0; i <= len(intervals)-3; i++ {
		a := intervals[i]
		b := intervals[i+1]
		c := intervals[i+2]

		// Verify at least two elements are different and at least one is a leap
		if (a != b || b != c) && (leaps[a] || leaps[b] || leaps[c]) {
			patterns = append(patterns, [3]int{a, b, c})
			patternIndices = append(patternIndices, i)
		}
	}

	// Look for repeating patterns
	for i := 0; i < len(patterns); i++ {
		for j := i + 1; j < len(patterns); j++ {
			if patterns[i] == patterns[j] {
				// Found repeating pattern
				start1 := patternIndices[i]
				start2 := patternIndices[j]

				// Extend first pattern up to the start of the second
				extended1 := extendPattern(intervals, start1, start2)
				if len(extended1) == 0 {
					continue // Not enough elements for comparison
				}

				// Extend second pattern to match the length of the first
				extended2 := extendPattern(intervals, start2, start2+len(extended1))
				if len(extended2) != len(extended1) {
					continue // Not enough elements for comparison
				}

				// Compare the extended patterns
				if equalSlices(extended1, extended2) {
					return true
				}
			}
		}
	}

	return false
}

// extendPattern extends the pattern from start to end using elements from intervals
func extendPattern(intervals []int, start, end int) []int {
	if end > len(intervals) {
		return nil // Not enough elements to extend
	}
	return intervals[start:end]
}

// equalSlices checks if two slices are equal
func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// AvoidSeventhNinthBetweenExtremes checks that there are no seventh (6) or ninth (8) intervals
// between the tonic and extreme notes (highest, lowest).
// This function should only be applied to complete interval slices.
//
// Logic:
//   - Builds partial sums slice (notes relative to tonic)
//   - Checks:
//     1. If maximum is multiple of 6 or 8 - false (seventh/ninth between tonic and highest note)
//     2. If minimum is multiple of 6 or 8 - false (seventh/ninth between tonic and lowest note)
//     3. If difference between max and min is multiple of 6 or 8 - false (seventh/ninth between extremes)
//   - Returns true in all other cases
func AvoidSeventhNinthBetweenExtremes(intervals []int) bool {
	if len(intervals) == 0 {
		return true
	}

	// Build partial sums (notes relative to tonic)
	partialSums := make([]int, len(intervals)+1)
	partialSums[0] = 0
	for i, interval := range intervals {
		partialSums[i+1] = partialSums[i] + interval
	}

	// Find maximum and minimum
	maxSum := partialSums[0]
	minSum := partialSums[0]
	for _, sum := range partialSums {
		if sum > maxSum {
			maxSum = sum
		}
		if sum < minSum {
			minSum = sum
		}
	}

	// Check conditions
	if isMultipleOfSixOrEight(maxSum) {
		return false
	}
	if isMultipleOfSixOrEight(minSum) {
		return false
	}
	if isMultipleOfSixOrEight(maxSum - minSum) {
		return false
	}

	return true
}

// isMultipleOfSixOrEight checks if a number is multiple of 6 or 8 (absolute value)
func isMultipleOfSixOrEight(x int) bool {
	absX := utils.Abs(x)
	return absX == 6 || absX == 8 || absX == 12 || absX == 16 // etc. for larger octaves
}

// ValidateLeadingTone checks the rules for the introductory tone in a partial interval slice.
// Returns true if all rules are satisfied, false otherwise.
func ValidateLeadingTone(intervals []int) bool {
	if len(intervals) == 0 {
		return true
	}

	// Build a slice of partial sums (notes relative to the starting note)
	partialSums := make([]int, len(intervals)+1)
	partialSums[0] = 0
	for i, interval := range intervals {
		partialSums[i+1] = partialSums[i] + interval
	}

	// Check each partial sum against the introductory tone rules
	for i, sum := range partialSums {
		switch sum {
		case -1:
			if !isValidMinusOne(partialSums, i) {
				return false
			}
		case 6:
			if !isValidSix(partialSums, i) {
				return false
			}
		case -8:
			if !isValidMinusEight(partialSums, i) {
				return false
			}
		}
	}

	return true
}

// isValidMinusOne checks valid configurations for the number -1
func isValidMinusOne(sums []int, index int) bool {
	// Check configuration ..., 0, -1, 0, ...
	if index > 0 && index < len(sums)-1 {
		if sums[index-1] == 0 && sums[index+1] == 0 {
			return true
		}
	}

	// Check configuration ..., -2, -1, 0, ...
	if index > 0 && index < len(sums)-1 {
		if sums[index-1] == -2 && sums[index+1] == 0 {
			return true
		}
	}

	// Check configuration ..., 0, -1, -2, ...
	if index > 0 && index < len(sums)-1 {
		if sums[index-1] == 0 && sums[index+1] == -2 {
			return true
		}
	}

	return false
}

// isValidSix checks valid configurations for the number 6
func isValidSix(sums []int, index int) bool {
	// Check configuration ..., 7, 6, 7, ...
	if index > 0 && index < len(sums)-1 {
		if sums[index-1] == 7 && sums[index+1] == 7 {
			return true
		}
	}

	// Check configuration ..., 5, 6, 7, ...
	if index > 0 && index < len(sums)-1 {
		if sums[index-1] == 5 && sums[index+1] == 7 {
			return true
		}
	}

	// Check configuration ..., 7, 6, 5, ...
	if index > 0 && index < len(sums)-1 {
		if sums[index-1] == 7 && sums[index+1] == 5 {
			return true
		}
	}

	return false
}

// isValidMinusEight checks valid configurations for the number -8
func isValidMinusEight(sums []int, index int) bool {
	// Check configuration ..., -7, -8, -7, ...
	if index > 0 && index < len(sums)-1 {
		if sums[index-1] == -7 && sums[index+1] == -7 {
			return true
		}
	}

	// Check configuration ..., -9, -8, -7, ...
	if index > 0 && index < len(sums)-1 {
		if sums[index-1] == -9 && sums[index+1] == -7 {
			return true
		}
	}

	// Check configuration ..., -7, -8, -9, ...
	if index > 0 && index < len(sums)-1 {
		if sums[index-1] == -7 && sums[index+1] == -9 {
			return true
		}
	}

	return false
}

// NoCloseLargeLeaps checks that there are no two leaps (absolute value > 2)
// separated by a single step (any interval). This prevents patterns like leap-step-leap
// or leap-leap-leap
// Returns:
//   - false if two leaps are found with one interval between them (rule violated)
//   - true otherwise (rule satisfied)
func NoCloseLargeLeaps(intervals []int) bool {
	if len(intervals) < 3 {
		return true
	}

	for i := 0; i <= len(intervals)-3; i++ {
		first := intervals[i]
		second := intervals[i+2]

		if utils.Abs(first) > 2 && utils.Abs(second) > 2 {
			return false
		}
	}

	return true
}

// NoMoreThanTwoConsecutiveThirds checks that there are no more than two consecutive intervals
// with absolute value equal to 2 in the interval sequence.
// Returns:
//   - false if three or more consecutive intervals with absolute value 2 are found (rule violated)
//   - true otherwise (rule satisfied)
func NoMoreThanTwoConsecutiveThirds(intervals []int) bool {
	if len(intervals) < 3 {
		return true
	}

	consecutiveTwos := 0
	for _, interval := range intervals {
		if utils.Abs(interval) == 2 {
			consecutiveTwos++
			if consecutiveTwos > 2 {
				return false
			}
		} else {
			consecutiveTwos = 0
		}
	}

	return true
}
