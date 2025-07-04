package rules

import (
	"go-cantus-firmus/internal/music"
)

// IsFreeOfAugmentedDiminished checks if a Realization contains no augmented ("A") or diminished ("d") intervals
// between adjacent notes or notes separated by one note, and also checks for these intervals in extremum notes.
// Returns:
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

	// Find all extremum notes (step 3)
	extremums := make(music.Realization, 0)
	extremums = append(extremums, r[0]) // Always include first note

	for i := 1; i < len(r)-1; i++ {
		prev := r[i-1]
		current := r[i]
		next := r[i+1]

		// Check if current note is an extremum
		if (current.Greater(prev) && current.Greater(next)) ||
			(current.Less(prev) && current.Less(next)) {
			extremums = append(extremums, current)
		}
	}

	extremums = append(extremums, r[len(r)-1]) // Always include last note

	// Check intervals between adjacent extremums (step 4)
	for i := 0; i < len(extremums)-1; i++ {
		quality, err := music.CalculateIntervalQuality(extremums[i], extremums[i+1])
		if err != nil {
			return false
		}
		if quality == "d" || quality == "A" {
			return false
		}
	}

	// Check intervals between extremums separated by one note (step 5)
	for i := 0; i < len(extremums)-2; i++ {
		quality, err := music.CalculateIntervalQuality(extremums[i], extremums[i+2])
		if err != nil {
			return false
		}
		if quality == "d" || quality == "A" {
			return false
		}
	}

	return true
}
