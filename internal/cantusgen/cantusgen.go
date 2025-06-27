package cantusgen

import "fmt"

var steps = []int{-1, 1}
var leaps = []int{-4, -3, -2, 2, 3, 4, 5}

// GenerateCantus generates a set of integer slices of length n,
// satisfying the following conditions:
// - Each slice contains approximately 70% elements from 'steps' and 30% from 'leaps'.
// - Each slice ends with two numbers from 'steps'.
// - The sum of all elements in each slice must be 0.
// - All possible permutations with repetitions that meet the conditions are generated.
func GenerateCantus(n int) [][]int {
	if n < 2 {
		fmt.Println("Slice length (n) must be at least 2, as the last two elements must be from 'steps'.")
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
		fmt.Printf("Cannot satisfy the 70%% 'steps' condition for n=%d. At least 2 'steps' are required for the end. Increase n or reconsider the percentage.\n", n)
		return nil
	}
	if requiredLeaps < 0 {
		fmt.Printf("Cannot satisfy the 30%% 'leaps' condition for n=%d. Increase n or reconsider the percentage.\n", n)
		return nil
	}

	// Number of 'steps' and 'leaps' for the first (n-2) positions
	stepsForPrefix := requiredSteps - 2
	leapsForPrefix := requiredLeaps

	if stepsForPrefix < 0 {
		fmt.Printf("Insufficient 'steps' for the prefix. 'n' might be too small or the 'steps' percentage too high. n=%d, stepsForPrefix=%d\n", n, stepsForPrefix)
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
				generatePrefix(currentIndex+1, append(currentSlice, val), currentSum+val, currentStepsCount+1, currentLeapsCount)
			}
		}

		// Attempt to add an element from 'leaps'
		if currentLeapsCount < leapsForPrefix {
			for _, val := range leaps {
				generatePrefix(currentIndex+1, append(currentSlice, val), currentSum+val, currentStepsCount, currentLeapsCount+1)
			}
		}
	}

	// Start the generation process
	generatePrefix(0, []int{}, 0, 0, 0)

	return result
}
