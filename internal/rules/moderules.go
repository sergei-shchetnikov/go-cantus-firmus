package rules

import "go-cantus-firmus/internal/music"

// IsFreeOfAugmentedDiminished checks a Realization for specific conditions related to augmented or diminished intervals.
func IsFreeOfAugmentedDiminished(r music.Realization) bool {
	return rule1(r) && rule2(r)
}

// rule1 checks every pair of notes n1 and n2 within a distance of 2 or fewer other notes
// (i.e., indices differ by 3 or less), if the interval between n1 and n2 is augmented ("A")
// or diminished ("d"), then at least one of n1 or n2 must be surrounded by linear motion
// (as determined by IsNoteSurroundedByLinearMotion). If this condition is not met for any pair,
// the function immediately returns false. If all such pairs satisfy the condition, it returns true.
func rule1(r music.Realization) bool {
	for i := range r {
		for j := i + 1; j < len(r); j++ {
			if j-i <= 3 {
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

	for i := 0; i < len(r)-1; i++ {
		for j := i + 1; j < len(r); j++ {
			subsequence := r[i : j+1]

			isCurrentSubsequenceMonotonic, _ := isStrictlyMonotonic(subsequence)

			if isCurrentSubsequenceMonotonic {
				canExtendLeft := false
				if i > 0 {
					potentialLeftSubsequence := r[i-1 : j+1]
					isLeftMonotonic, _ := isStrictlyMonotonic(potentialLeftSubsequence)
					canExtendLeft = isLeftMonotonic
				}

				canExtendRight := false
				if j < len(r)-1 {
					potentialRightSubsequence := r[i : j+2]
					isRightMonotonic, _ := isStrictlyMonotonic(potentialRightSubsequence)
					canExtendRight = isRightMonotonic
				}

				isMaximalByLeftExtension := (i == 0) || !canExtendLeft
				isMaximalByRightExtension := (j == len(r)-1) || !canExtendRight

				if isMaximalByLeftExtension && isMaximalByRightExtension {
					firstNote := subsequence[0]
					lastNote := subsequence[len(subsequence)-1]

					quality, err := music.CalculateIntervalQuality(firstNote, lastNote)
					if err != nil {
						return false
					}

					if quality == "A" || quality == "d" {
						return false
					}
				}
			}
		}
	}
	return true
}

// isStrictlyMonotonic is a helper function to check if a Realization subsequence is strictly monotonic.
// A sequence is strictly monotonic if its notes consistently move in one direction (all ascending or all descending),
// and no two adjacent notes have the same pitch.
// It returns true if the subsequence is strictly ascending or strictly descending, along with its direction (true for ascending, false for descending).
// Returns false and nil direction otherwise.
func isStrictlyMonotonic(subsequence music.Realization) (bool, *bool) {
	if len(subsequence) < 2 {
		return false, nil
	}

	var direction *bool = nil

	for k := 0; k < len(subsequence)-1; k++ {
		n1 := subsequence[k]
		n2 := subsequence[k+1]

		isCurrentMovementAscending := n2.Greater(n1)
		isSamePitch := n2.EqualPitch(n1)

		if isSamePitch {
			return false, nil
		}

		if direction == nil {
			b := isCurrentMovementAscending
			direction = &b
		} else if *direction != isCurrentMovementAscending {
			return false, nil
		}
	}
	return true, direction
}
