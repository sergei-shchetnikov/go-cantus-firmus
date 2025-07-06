package utils

import "testing"

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
