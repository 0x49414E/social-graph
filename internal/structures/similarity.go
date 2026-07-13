package structures

// Jaccard returns the Jaccard similarity of a and b: |a ∩ b| / |a ∪ b|.
// Two empty sets are defined as identical (1).
func Jaccard[T comparable](a, b map[T]struct{}) float64 {
	score, _ := JaccardWithOverlap(a, b)
	return score
}

// JaccardWithOverlap returns the Jaccard similarity of a and b together with the
// size of their intersection (the overlap).
func JaccardWithOverlap[T comparable](a, b map[T]struct{}) (float64, int) {
	if len(a) == 0 && len(b) == 0 {
		return 1, 0
	}

	// Iterate over the smaller set
	if len(a) > len(b) {
		a, b = b, a
	}

	intersection := 0
	for x := range a {
		if _, ok := b[x]; ok {
			intersection++
		}
	}

	union := len(a) + len(b) - intersection
	return float64(intersection) / float64(union), intersection
}
