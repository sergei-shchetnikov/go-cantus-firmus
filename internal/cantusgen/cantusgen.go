package cantusgen

import (
	"go-cantus-firmus/internal/rules"
)

var steps = []int{-1, 1}
var leaps = []int{-4, -3, -2, 2, 3, 4, 5}

// Define the set of validation functions
var cantusValidators = []rules.ValidationFunc{
	rules.NoBeginWithFive,
	rules.NoExcessiveNoteRepetition,
	rules.NoFiveOfSameSign,
	rules.NoRangeExceedsDecima,
	rules.NoRepeatingPatterns,
	rules.PreparedLeaps,
	rules.ValidateLeapResolution,
}

// GenerateCantus generates a set of integer slices of length n,
// satisfying specific contrapuntal and structural conditions:
//   - The sum of all intervals in the complete slice equals 0 (returns to starting pitch)
//   - Approximately 70% of intervals are step motions (values from {-1, 1})
//   - The slice always ends with two step motions (values from {-1, 1})
//   - All slices adhere to contrapuntal rules defined in cantusValidators
//
// The function uses a recursive backtracking approach, pruning invalid partial
// melodies early based on the `cantusValidators` to efficiently generate
// only valid cantus firmus melodies.
// The resulting slice of intervals of length 'n' serves as the basis
// for the implementation of cantus firmus of length 'n+1' in specific notes.
func GenerateCantus(n int) [][]int {
	if n < 2 {
		return nil
	}

	var result [][]int

	requiredSteps := int(float64(n) * 0.7)
	requiredLeaps := n - requiredSteps

	if requiredSteps < 2 || requiredLeaps < 0 {
		return nil
	}

	stepsForPrefix := requiredSteps - 2
	leapsForPrefix := requiredLeaps

	if stepsForPrefix < 0 {
		return nil
	}

	var generatePrefix func(currentIndex int, currentSlice []int, currentSum int, currentStepsCount int, currentLeapsCount int)
	generatePrefix = func(currentIndex int, currentSlice []int, currentSum int, currentStepsCount int, currentLeapsCount int) {
		if !rules.AllRules(currentSlice, cantusValidators) {
			return
		}

		if currentIndex == n-2 {
			for _, end1Val := range steps {
				for _, end2Val := range steps {
					finalSlice := make([]int, n)
					copy(finalSlice, currentSlice)
					finalSlice[n-2] = end1Val
					finalSlice[n-1] = end2Val

					if !rules.AllRules(finalSlice, cantusValidators) {
						continue
					}

					totalSum := currentSum + end1Val + end2Val
					if totalSum == 0 {
						result = append(result, finalSlice)
					}
				}
			}
			return
		}

		if currentStepsCount < stepsForPrefix {
			for _, val := range steps {
				nextSlice := append(currentSlice, val)
				generatePrefix(currentIndex+1, nextSlice, currentSum+val, currentStepsCount+1, currentLeapsCount)
			}
		}

		if currentLeapsCount < leapsForPrefix {
			for _, val := range leaps {
				nextSlice := append(currentSlice, val)
				generatePrefix(currentIndex+1, nextSlice, currentSum+val, currentStepsCount, currentLeapsCount+1)
			}
		}
	}

	generatePrefix(0, []int{}, 0, 0, 0)

	return result
}
