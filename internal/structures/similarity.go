package structures

func Jaccard[T comparable](a, b map[T]struct{}) float64 {
	if len(a) == 0 && len(b) == 0 {
		return 1
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
	return float64(intersection) / float64(union)
}
