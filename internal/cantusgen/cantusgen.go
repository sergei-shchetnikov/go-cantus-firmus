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

	// Determine the required number of elements from 'steps' and 'leaps'
	// Account for the last 2 elements already being 'steps'
	requiredSteps := int(float64(n) * 0.7)
	requiredLeaps := n - requiredSteps

	// If requiredSteps is less than 2, it means even with the last two steps
	// we won't be able to reach 70% 'steps' if n is very small.
	// Or if requiredLeaps is negative, which is also an error.
	if requiredSteps < 2 {
		return nil
	}
	if requiredLeaps < 0 {
		return nil
	}

	// Number of 'steps' and 'leaps' for the first (n-2) positions
	stepsForPrefix := requiredSteps - 2
	leapsForPrefix := requiredLeaps

	if stepsForPrefix < 0 {
		return nil
	}

	// Recursive function to generate the prefix
	var generatePrefix func(currentIndex int, currentSlice []int, currentSum int, currentStepsCount int, currentLeapsCount int)
	generatePrefix = func(currentIndex int, currentSlice []int, currentSum int, currentStepsCount int, currentLeapsCount int) {
		// Base case: prefix is filled
		if currentIndex == n-2 {
			// Now try adding the last two elements from 'steps'
			for _, end1Val := range steps {
				for _, end2Val := range steps {
					finalSlice := make([]int, n)
					copy(finalSlice, currentSlice)
					finalSlice[n-2] = end1Val
					finalSlice[n-1] = end2Val

					// !!! Important: Final check for NoFiveOfSameSign on the complete slice !!!
					// Although partial checks are done during generation,
					// this ensures the very end of the slice (which was not part of the prefix)
					// also adheres to the rule if it forms a consecutive sequence.
					if !rules.NoFiveOfSameSign(finalSlice) {
						continue // Skip this finalSlice if it violates the rule
					}

					totalSum := currentSum + end1Val + end2Val
					if totalSum == 0 {
						result = append(result, finalSlice)
					}
				}
			}
			return
		}

		// Attempt to add an element from 'steps'
		if currentStepsCount < stepsForPrefix {
			for _, val := range steps {
				nextSlice := append(currentSlice, val)
				// No need to check for 5 here, as steps does not contain 5

				// !!! Key step: Intermediate check for NoFiveOfSameSign !!!
				if !rules.NoFiveOfSameSign(nextSlice) { //
					continue // Skip this branch if the condition is violated
				}

				generatePrefix(currentIndex+1, nextSlice, currentSum+val, currentStepsCount+1, currentLeapsCount)
			}
		}

		// Attempt to add an element from 'leaps'
		if currentLeapsCount < leapsForPrefix {
			for _, val := range leaps {
				// If it's the first element (currentIndex == 0) and the value is 5, skip it
				if currentIndex == 0 && val == 5 {
					continue
				}
				nextSlice := append(currentSlice, val)

				// !!! Key step: Intermediate check for NoFiveOfSameSign !!!
				if !rules.NoFiveOfSameSign(nextSlice) { //
					continue // Skip this branch if the condition is violated
				}

				generatePrefix(currentIndex+1, nextSlice, currentSum+val, currentStepsCount, currentLeapsCount+1)
			}
		}
	}

	// Start the generation process
	generatePrefix(0, []int{}, 0, 0, 0)

	return result
}
