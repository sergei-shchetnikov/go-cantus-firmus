package main

import "testing"

func TestMod7(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{"Number less than 7", 5, 5},
		{"Number greater than -7", -5, 2},
		{"Positive number", 10, 3},
		{"Negative number", -10, 4},
		{"Positive multiple of 7", 14, 0},
		{"Negative multiple of 7", -7, 0},
		{"Zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Mod7(tt.n); got != tt.want {
				t.Errorf("Mod7(%d) = %d, want %d", tt.n, got, tt.want)
			}
		})
	}
}
