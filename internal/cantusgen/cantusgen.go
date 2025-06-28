package cantusgen

import (
	"go-cantus-firmus/internal/rules"
)

var steps = []int{-1, 1}
var leaps = []int{-4, -3, -2, 2, 3, 4, 5}

// GenerateCantus generates a set of integer slices of length n,
// satisfying the following conditions:
// - Each slice contains approximately 70% elements from 'steps' and 30% from 'leaps'.
// - Each slice ends with two numbers from 'steps'.
// - The sum of all elements in each slice must be 0.
// - All possible permutations with repetitions that meet the conditions are generated.
// - The first element of each slice must not be 5.
// - No more than five consecutive numbers can have the same sign (enforced by rules.NoFiveOfSameSign).
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
		// Check note repetition rule for current partial slice
		if !rules.NoExcessiveNoteRepetition(currentSlice) {
			return
		}

		if currentIndex == n-2 {
			for _, end1Val := range steps {
				for _, end2Val := range steps {
					finalSlice := make([]int, n)
					copy(finalSlice, currentSlice)
					finalSlice[n-2] = end1Val
					finalSlice[n-1] = end2Val

					if !rules.NoFiveOfSameSign(finalSlice) || !rules.NoExcessiveNoteRepetition(finalSlice) {
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

				if !rules.NoFiveOfSameSign(nextSlice) {
					continue
				}

				generatePrefix(currentIndex+1, nextSlice, currentSum+val, currentStepsCount+1, currentLeapsCount)
			}
		}

		if currentLeapsCount < leapsForPrefix {
			for _, val := range leaps {
				if currentIndex == 0 && val == 5 {
					continue
				}
				nextSlice := append(currentSlice, val)

				if !rules.NoFiveOfSameSign(nextSlice) {
					continue
				}

				generatePrefix(currentIndex+1, nextSlice, currentSum+val, currentStepsCount, currentLeapsCount+1)
			}
		}
	}

	generatePrefix(0, []int{}, 0, 0, 0)

	return result
}
