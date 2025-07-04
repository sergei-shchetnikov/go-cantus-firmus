package rules

import (
	"go-cantus-firmus/internal/music"
)

// IsFreeOfAugmentedDiminished checks if a Realization contains no augmented ("A") or diminished ("d") intervals
// between adjacent notes or notes separated by one note. Returns:
//   - true if all checked intervals are perfect, major, or minor
//   - false if any augmented/diminished interval is found
//   - false if interval calculation fails (invalid notes)
func IsFreeOfAugmentedDiminished(r music.Realization) bool {
	if len(r) < 2 {
		return true // No intervals to check
	}

	// Check intervals between adjacent notes (step 1)
	for i := 0; i < len(r)-1; i++ {
		quality, err := music.CalculateIntervalQuality(r[i], r[i+1])
		if err != nil {
			return false
		}
		if quality == "d" || quality == "A" {
			return false
		}
	}

	// Check intervals between notes separated by one note (step 2)
	for i := 0; i < len(r)-2; i++ {
		quality, err := music.CalculateIntervalQuality(r[i], r[i+2])
		if err != nil {
			return false
		}
		if quality == "d" || quality == "A" {
			return false
		}
	}

	return true
}
