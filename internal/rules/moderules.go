package rules

import "go-cantus-firmus/internal/music"

// IsFreeOfAugmentedDiminished checks a Realization for specific conditions related to augmented or diminished intervals.
func IsFreeOfAugmentedDiminished(r music.Realization) bool {
	return rule1(r) && rule2(r)
}

// rule1 checks every pair of notes n1 and n2 within a distance of 3 or fewer other notes
// (i.e., indices differ by 4 or less), if the interval between n1 and n2 is augmented ("A")
// or diminished ("d"), then at least one of n1 or n2 must be surrounded by linear motion
// (as determined by IsNoteSurroundedByLinearMotion). If this condition is not met for any pair,
// the function immediately returns false. If all such pairs satisfy the condition, it returns true.
func rule1(r music.Realization) bool {
	for i := range r {
		for j := i + 1; j < len(r); j++ {
			if j-i <= 4 {
				n1 := r[i]
				n2 := r[j]

				quality, err := music.CalculateIntervalQuality(n1, n2)
				if err != nil {
					return false
				}

				if quality == "A" || quality == "d" {
					n1LinearMotion := music.IsNoteSurroundedByLinearMotion(r, i)
					n2LinearMotion := music.IsNoteSurroundedByLinearMotion(r, j)

					if !(n1LinearMotion || n2LinearMotion) {
						return false
					}
				}
			}
		}
	}
	return true
}

// rule2 checks a Realization for ascending or descending sequences of notes.
// It returns false immediately if the interval between the first and last notes
// of such a sequence is augmented ("A") or diminished ("d").
// This rule considers any change in step as part of the sequence as long as the direction is maintained.
func rule2(r music.Realization) bool {
	if len(r) < 2 {
		return true
	}

	for i := 0; i < len(r); i++ {
		currentNote := r[i]
		var sequenceDirection *bool = nil
		sequenceStartIdx := i

		for j := i + 1; j < len(r); j++ {
			nextNote := r[j]

			isCurrentMovementAscending := nextNote.Greater(currentNote)
			isCurrentMovementDescending := nextNote.Less(currentNote)

			if !isCurrentMovementAscending && !isCurrentMovementDescending {
				if j-1 > sequenceStartIdx {
					n1 := r[sequenceStartIdx]
					n2 := r[j-1]

					quality, err := music.CalculateIntervalQuality(n1, n2)
					if err != nil {
						return false
					}

					if quality == "A" || quality == "d" {
						return false
					}
				}
				i = j - 1
				break
			}

			if sequenceDirection == nil {
				b := isCurrentMovementAscending
				sequenceDirection = &b
			} else if *sequenceDirection != isCurrentMovementAscending {
				if j-1 > sequenceStartIdx {
					n1 := r[sequenceStartIdx]
					n2 := r[j-1]

					quality, err := music.CalculateIntervalQuality(n1, n2)
					if err != nil {
						return false
					}

					if quality == "A" || quality == "d" {
						return false
					}
				}
				i = j - 1
				break
			}

			if j == len(r)-1 {
				n1 := r[sequenceStartIdx]
				n2 := r[j]

				quality, err := music.CalculateIntervalQuality(n1, n2)
				if err != nil {
					return false
				}

				if quality == "A" || quality == "d" {
					return false
				}
				i = j
			}
			currentNote = nextNote
		}
	}
	return true
}
