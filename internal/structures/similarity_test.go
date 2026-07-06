package structures

import "testing"

// set builds a set out of the given items so the tests stay readable.
func set[T comparable](items ...T) map[T]struct{} {
	s := make(map[T]struct{}, len(items))
	for _, item := range items {
		s[item] = struct{}{}
	}
	return s
}

func TestJaccard(t *testing.T) {
	tests := []struct {
		name string
		a, b map[int]struct{}
		want float64
	}{
		{
			name: "identical sets",
			a:    set(1, 2, 3),
			b:    set(1, 2, 3),
			want: 1.0,
		},
		{
			name: "disjoint sets",
			a:    set(1, 2, 3),
			b:    set(4, 5, 6),
			want: 0.0,
		},
		{
			name: "partial overlap",
			// they share {2,3} and cover {1,2,3,4} together, so 2/4 = 0.5
			a:    set(1, 2, 3),
			b:    set(2, 3, 4),
			want: 0.5,
		},
		{
			name: "single common element",
			// only 2 is shared and there are 5 distinct values total, so 1/5 = 0.2
			a:    set(1, 2, 3),
			b:    set(2, 4, 5),
			want: 0.2,
		},
		{
			name: "both empty is defined as identical",
			a:    set[int](),
			b:    set[int](),
			want: 1.0,
		},
		{
			name: "one empty set",
			a:    set(1, 2),
			b:    set[int](),
			want: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Jaccard(tt.a, tt.b); got != tt.want {
				t.Errorf("Jaccard() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Jaccard iterates over the smaller set internally, so passing the arguments in
// either order should give back the exact same number.
func TestJaccardIsSymmetric(t *testing.T) {
	a := set(1, 2, 3, 4)
	b := set(3, 4)

	if forward, backward := Jaccard(a, b), Jaccard(b, a); forward != backward {
		t.Errorf("Jaccard not symmetric: Jaccard(a,b)=%v, Jaccard(b,a)=%v", forward, backward)
	}
}
