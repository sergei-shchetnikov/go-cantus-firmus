package utils

import (
	"testing"
	"time"
)

func TestAbs(t *testing.T) {
	tests := []struct {
		name  string
		input int
		want  int
	}{
		{"positive number", 5, 5},
		{"negative number", -3, 3},
		{"zero", 0, 0},
		{"max positive", 1<<31 - 1, 1<<31 - 1},
		{"min negative", -1 << 31, 1 << 31},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Abs(tt.input); got != tt.want {
				t.Errorf("Abs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelectRandomItems(t *testing.T) {
	// Log the random seed for potential debugging
	randSeed := time.Now().UnixNano()
	t.Logf("Using random seed: %d", randSeed)

	tests := []struct {
		name     string
		items    []int
		count    int
		expected int // expected number of items in the result
	}{
		{
			name:     "empty slice",
			items:    []int{},
			count:    3,
			expected: 0,
		},
		{
			name:     "zero count",
			items:    []int{1, 2, 3},
			count:    0,
			expected: 0,
		},
		{
			name:     "negative count",
			items:    []int{1, 2, 3},
			count:    -1,
			expected: 0,
		},
		{
			name:     "count equals slice length",
			items:    []int{1, 2, 3, 4, 5},
			count:    5,
			expected: 5,
		},
		{
			name:     "count greater than slice length",
			items:    []int{1, 2, 3},
			count:    5,
			expected: 3,
		},
		{
			name:     "normal case",
			items:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			count:    3,
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SelectRandomItems(tt.items, tt.count)

			if len(result) != tt.expected {
				t.Errorf("expected %d items, got %d", tt.expected, len(result))
			}

			if len(result) > 0 {
				itemSet := make(map[int]bool)
				for _, item := range tt.items {
					itemSet[item] = true
				}

				for _, item := range result {
					if !itemSet[item] {
						t.Errorf("result contains item %v not in original slice", item)
					}
				}
			}
		})
	}

	t.Run("probability distribution", func(t *testing.T) {
		const (
			iterations = 10000
			sliceSize  = 10
			selectSize = 3
		)

		items := make([]int, sliceSize)
		for i := 0; i < sliceSize; i++ {
			items[i] = i + 1
		}

		counts := make(map[int]int)
		for i := 0; i < iterations; i++ {
			result := SelectRandomItems(items, selectSize)
			for _, item := range result {
				counts[item]++
			}
		}

		expectedProbability := float64(selectSize) / float64(sliceSize)
		allowedVariance := 0.02

		for item, count := range counts {
			actualProbability := float64(count) / float64(iterations)
			if abs(actualProbability-expectedProbability) > allowedVariance {
				t.Errorf("item %d was selected with probability %.4f, expected %.4f (±%.4f)",
					item, actualProbability, expectedProbability, allowedVariance)
			}
		}
	})
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
