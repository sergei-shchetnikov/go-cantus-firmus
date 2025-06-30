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
}

// Validation functions that require complete slices (length n) to evaluate
// These rules check structural properties that only make sense for complete compositions
var completeCantusValidators = []rules.ValidationFunc{
	rules.MinDirectionChanges,
}

// GenerateCantus generates a set of integer slices of length n,
// satisfying specific contrapuntal and structural conditions:
//   - The sum of all intervals in the complete slice equals 0 (returns to starting pitch)
//   - Approximately 70% of intervals are step motions (values from {-1, 1})
//   - The slice always ends with two step motions (values from {-1, 1})
//   - All slices adhere to both partial (cantusValidators) and complete (completeCantusValidators) rules
//
// The function uses recursive backtracking with these optimization strategies:
//   - Early pruning of invalid partial melodies using cantusValidators
//   - Final validation of complete melodies using completeCantusValidators
//   - Step/leap ratio enforcement throughout generation
func GenerateCantus(n int) [][]int {
	if n < 2 {
		return nil
	}

	var result [][]int

	// Calculate required step/leap distribution (70% steps)
	requiredSteps := int(float64(n) * 0.7)
	requiredLeaps := n - requiredSteps

	if requiredSteps < 2 || requiredLeaps < 0 {
		return nil
	}

	// Reserve last 2 positions for steps (standard cadence)
	stepsForPrefix := requiredSteps - 2
	leapsForPrefix := requiredLeaps

	if stepsForPrefix < 0 {
		return nil
	}

	var generatePrefix func(currentIndex int, currentSlice []int, currentSum int, currentStepsCount int, currentLeapsCount int)
	generatePrefix = func(currentIndex int, currentSlice []int, currentSum int, currentStepsCount int, currentLeapsCount int) {
		// Validate partial melody against partial rules
		if !rules.AllRules(currentSlice, cantusValidators) {
			return
		}

		// When we reach the position where we need to add the final two steps
		if currentIndex == n-2 {
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

		// Continue building melody with steps if we haven't used our step quota
		if currentStepsCount < stepsForPrefix {
			for _, val := range steps {
				nextSlice := append(currentSlice, val)
				generatePrefix(currentIndex+1, nextSlice, currentSum+val, currentStepsCount+1, currentLeapsCount)
			}
		}

		// Continue building melody with leaps if we haven't used our leap quota
		if currentLeapsCount < leapsForPrefix {
			for _, val := range leaps {
				nextSlice := append(currentSlice, val)
				generatePrefix(currentIndex+1, nextSlice, currentSum+val, currentStepsCount, currentLeapsCount+1)
			}
		}
	}

	// Start generation with empty slice
	generatePrefix(0, []int{}, 0, 0, 0)

	return result
}
