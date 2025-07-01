package cantusgen

import (
	"slices"
	"testing"
)

func TestGenerateCantus_InvalidInput(t *testing.T) {
	tests := []struct {
		name         string
		n            int
		allowedLeaps []int
	}{
		{"n less than 2", 1, []int{2, 3}},
		{"no allowed leaps", 10, []int{}},
		{"all leap counts too large", 10, []int{9, 10}}, // only 8 possible leaps max (n-2)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateCantus(tt.n, tt.allowedLeaps)
			if result != nil {
				t.Errorf("Expected nil result for n=%d, allowedLeaps=%v, got %v", tt.n, tt.allowedLeaps, result)
			}
		})
	}
}

func TestGenerateCantus_ValidInput(t *testing.T) {
	tests := []struct {
		name         string
		n            int
		allowedLeaps []int
	}{
		{"n=5 with 1-2 leaps", 5, []int{1, 2}},
		{"n=8 with 2-3 leaps", 8, []int{2, 3}},
		{"n=10 with 2-4 leaps", 10, []int{2, 3, 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateCantus(tt.n, tt.allowedLeaps)
			if result == nil {
				t.Fatalf("Expected non-nil result for n=%d, allowedLeaps=%v", tt.n, tt.allowedLeaps)
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

				// Count steps and leaps (excluding last two steps)
				stepsCount := 0
				leapsCount := 0
				for i := 0; i < tt.n-2; i++ {
					if contains(steps, sequence[i]) {
						stepsCount++
					} else if contains(leaps, sequence[i]) {
						leapsCount++
					} else {
						t.Errorf("Unexpected value %d in sequence %v", sequence[i], sequence)
					}
				}

				// Verify leap count is in allowedLeaps
				if !contains(tt.allowedLeaps, leapsCount) {
					t.Errorf("Leap count %d not in allowedLeaps %v for sequence %v", leapsCount, tt.allowedLeaps, sequence)
				}

				// Verify last two elements are steps
				if !contains(steps, sequence[tt.n-2]) || !contains(steps, sequence[tt.n-1]) {
					t.Errorf("Last two elements must be steps, got %v", sequence[tt.n-2:])
				}
			}
		})
	}
}

func TestGenerateCantusContainsSpecificSlice(t *testing.T) {
	n := 10
	allowedLeaps := []int{2, 3, 4}
	targets := [][]int{
		{2, -1, -1, 3, -1, 2, -1, -1, -1, -1}, // cantus by Johann Joseph Fux
		{1, 2, -1, 1, 1, 1, -1, -2, -1, -1},   // Schenker
	}

	result := GenerateCantus(n, allowedLeaps)

	for _, target := range targets {
		found := false
		for _, slice := range result {
			if equalSlices(slice, target) {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Expected slice %v not found in GenerateCantus(%d, %v) result", target, n, allowedLeaps)
		}
	}
}

// Helper function to check if a value exists in a slice
func contains(slice []int, val int) bool {
	return slices.Contains(slice, val)
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
