package rules

import (
	"go-cantus-firmus/internal/music"
)

// IsFreeOfAugmentedDiminished checks if a Realization contains no augmented ("A") or diminished ("d") intervals
// between adjacent notes. Returns:
//   - true if all intervals are perfect, major, or minor
//   - false if any augmented/diminished interval is found
//   - false if interval calculation fails (invalid notes)
func IsFreeOfAugmentedDiminished(r music.Realization) bool {
	if len(r) < 2 {
		return true // No intervals to check
	}

	for i := 0; i < len(r)-1; i++ {
		quality, err := music.CalculateIntervalQuality(r[i], r[i+1])
		if err != nil {
			// Treat calculation errors as invalid intervals
			return false
		}

		if quality == "d" || quality == "A" {
			return false
		}
	}

	return true
}
