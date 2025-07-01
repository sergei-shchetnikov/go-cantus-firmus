package cantusgen

import (
	"go-cantus-firmus/internal/rules"
)

var steps = []int{-1, 1}
var leaps = []int{-4, -3, -2, 2, 3, 4, 5}

// Validation functions that can be checked on partial slices during generation
var cantusValidators = []rules.ValidationFunc{
	rules.NoBeginWithFive,
	rules.NoExcessiveNoteRepetition,
	rules.NoFiveOfSameSign,
	rules.NoRangeExceedsDecima,
	rules.NoRepeatingPatterns,
	rules.PreparedLeaps,
	rules.ValidateLeapResolution,
	rules.NoTripleAlternatingNote,
	rules.NoNoteRepetitionAfterLeap,
	rules.NoRepeatingExtremes,
	rules.AvoidSeventhBetweenExtrema,
	rules.NoTwoNoteSequences,
}

// Validation functions that require complete slices (length n) to evaluate
var completeCantusValidators = []rules.ValidationFunc{
	rules.MinDirectionChanges,
	rules.ValidateClimax,
}

// GenerateCantus generates a set of integer slices of length n,
// satisfying specific contrapuntal and structural conditions:
//   - The sum of all intervals in the complete slice equals 0 (returns to starting pitch)
//   - The slice always ends with two step motions (values from {-1, 1})
//   - All slices adhere to both partial (cantusValidators) and complete (completeCantusValidators) rules
//
// Parameters:
//   - n: the number of intervals between adjacent pairs of notes in cantus firmus
//   - allowedLeaps: slice of integers specifying allowed number of leaps (e.g. []int{2,3,4})
//
// The function uses recursive backtracking with these optimization strategies:
//   - Early pruning of invalid partial melodies using cantusValidators
//   - Final validation of complete melodies using completeCantusValidators
func GenerateCantus(n int, allowedLeaps []int) [][]int {
	if n < 2 {
		return nil
	}

	var result [][]int

	// Convert allowedLeaps to a map for faster lookup
	leapCounts := make(map[int]bool)
	for _, count := range allowedLeaps {
		if count >= 0 && count <= n-2 { // -2 because last two must be steps
			leapCounts[count] = true
		}
	}

	if len(leapCounts) == 0 {
		return nil
	}

	var generatePrefix func(currentIndex int, currentSlice []int, currentSum int, currentLeapsCount int)
	generatePrefix = func(currentIndex int, currentSlice []int, currentSum int, currentLeapsCount int) {
		// Validate partial melody against partial rules
		if !rules.AllRules(currentSlice, cantusValidators) {
			return
		}

		// When we reach the position where we need to add the final two steps
		if currentIndex == n-2 {
			// Check if current leaps count is in allowed counts
			if !leapCounts[currentLeapsCount] {
				return
			}

			for _, end1Val := range steps {
				for _, end2Val := range steps {
					finalSlice := make([]int, n)
					copy(finalSlice, currentSlice)
					finalSlice[n-2] = end1Val
					finalSlice[n-1] = end2Val

					// Validate complete melody against all rule sets
					if !rules.AllRules(finalSlice, cantusValidators) {
						continue
					}

					totalSum := currentSum + end1Val + end2Val
					if totalSum == 0 {
						// Final check for complete melody-specific rules
						if rules.AllRules(finalSlice, completeCantusValidators) {
							result = append(result, finalSlice)
						}
					}
				}
			}
			return
		}

		// Try adding a step (if we can still have steps)
		if (n - 2 - currentLeapsCount) > 0 { // -2 for final two steps
			for _, val := range steps {
				nextSlice := append(currentSlice, val)
				generatePrefix(currentIndex+1, nextSlice, currentSum+val, currentLeapsCount)
			}
		}

		// Try adding a leap (if we haven't exceeded allowed leaps)
		if currentLeapsCount < maxKey(leapCounts) {
			for _, val := range leaps {
				nextSlice := append(currentSlice, val)
				generatePrefix(currentIndex+1, nextSlice, currentSum+val, currentLeapsCount+1)
			}
		}
	}

	// Start generation with empty slice
	generatePrefix(0, []int{}, 0, 0)

	return result
}

// Helper function to get maximum key from leapCounts map
func maxKey(m map[int]bool) int {
	max := 0
	for k := range m {
		if k > max {
			max = k
		}
	}
	return max
}
