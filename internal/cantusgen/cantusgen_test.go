package cantusgen

import (
	"slices"
	"testing"
)

func TestGenerateCantus_InvalidInput(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		{"n less than 2", 1},
		{"n too small for 70% steps", 2}, // 70% of 2 is 1.4, but we need at least 2 steps at the end
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateCantus(tt.n)
			if result != nil {
				t.Errorf("Expected nil result for n=%d, got %v", tt.n, result)
			}
		})
	}
}

func TestGenerateCantus_ValidInput(t *testing.T) {
	tests := []struct {
		name string
		n    int
	}{
		{"n=5", 5},
		{"n=8", 8},
		{"n=10", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateCantus(tt.n)
			if result == nil {
				t.Fatalf("Expected non-nil result for n=%d", tt.n)
			}

			for _, sequence := range result {
				// Check sequence length
				if len(sequence) != tt.n {
					t.Errorf("Expected sequence length %d, got %d", tt.n, len(sequence))
				}

				// Check sum of elements
				sum := 0
				for _, val := range sequence {
					sum += val
				}
				if sum != 0 {
					t.Errorf("Expected sum 0, got %d for sequence %v", sum, sequence)
				}

				// Verify last two elements are from steps
				if !contains(steps, sequence[tt.n-2]) || !contains(steps, sequence[tt.n-1]) {
					t.Errorf("Last two elements must be from steps, got %v", sequence[tt.n-2:])
				}

				// Check steps and leaps distribution
				stepsCount := 0
				leapsCount := 0
				for _, val := range sequence {
					if contains(steps, val) {
						stepsCount++
					} else if contains(leaps, val) {
						leapsCount++
					} else {
						t.Errorf("Unexpected value %d in sequence %v", val, sequence)
					}
				}

				// Verify approximate 70/30 ratio (with ±1 tolerance)
				expectedSteps := int(float64(tt.n)*0.7 + 0.5) // rounding to nearest integer
				expectedLeaps := tt.n - expectedSteps

				if stepsCount < expectedSteps-1 || stepsCount > expectedSteps+1 {
					t.Errorf("Expected about %d steps (70%%), got %d in sequence %v", expectedSteps, stepsCount, sequence)
				}
				if leapsCount < expectedLeaps-1 || leapsCount > expectedLeaps+1 {
					t.Errorf("Expected about %d leaps (30%%), got %d in sequence %v", expectedLeaps, leapsCount, sequence)
				}
			}
		})
	}
}

// Helper function to check if a value exists in a slice
func contains(slice []int, val int) bool {
	return slices.Contains(slice, val)
}

func TestGenerateCantusContainsSpecificSlice(t *testing.T) {
	n := 10
	target := []int{2, -1, -1, 3, -1, 2, -1, -1, -1, -1} // cantus by Johann Joseph Fux

	result := GenerateCantus(n)

	found := false
	for _, slice := range result {
		if equalSlices(slice, target) {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected slice %v not found in GenerateCantus(%d) result", target, n)
	}
}

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
