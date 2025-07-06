package utils

import (
	"math/rand"
)

// Abs returns the absolute value of an integer
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// SelectRandomItems selects 'count' random items from a slice using reservoir sampling algorithm
func SelectRandomItems[T any](items []T, count int) []T {
	if count <= 0 || len(items) == 0 {
		return nil
	}
	if count >= len(items) {
		result := make([]T, len(items))
		copy(result, items)
		return result
	}

	result := make([]T, count)
	copy(result, items[:count])

	for i := count; i < len(items); i++ {
		j := rand.Intn(i + 1)
		if j < count {
			result[j] = items[i]
		}
	}

	return result
}
