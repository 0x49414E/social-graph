package structures

import "testing"

func TestUnionFindBasic(t *testing.T) {
	u := NewUnionFind[int]()
	u.MakeSet(1)
	u.MakeSet(2)

	// They haven't been unioned yet, so they should sit in different sets.
	if u.Find(1) == u.Find(2) {
		t.Fatal("1 and 2 should be in different sets before Union")
	}

	u.Union(1, 2)

	if u.Find(1) != u.Find(2) {
		t.Error("1 and 2 should share the same root after Union")
	}
}

// Union everything into one long chain. Afterwards every element should resolve
// to the same root, and Find should have flattened the tree so each id points
// straight at that root.
func TestUnionFindPathCompression(t *testing.T) {
	u := NewUnionFind[int]()
	for i := 1; i <= 5; i++ {
		u.MakeSet(i)
	}

	u.Union(1, 2)
	u.Union(2, 3)
	u.Union(3, 4)
	u.Union(4, 5)

	root := u.Find(1)
	for i := 2; i <= 5; i++ {
		if got := u.Find(i); got != root {
			t.Errorf("Find(%d) = %v, want common root %v", i, got, root)
		}
	}

	// Now that Find has touched every element, each id should point right at the
	// root in a single hop.
	for i := 1; i <= 5; i++ {
		if u.id[i] != root {
			t.Errorf("path not compressed: id[%d] = %v, want %v", i, u.id[i], root)
		}
	}
}

// Unioning two elements that are already together should be a no-op and leave
// the structure untouched.
func TestUnionFindSameSet(t *testing.T) {
	u := NewUnionFind[int]()
	u.MakeSet(1)
	u.MakeSet(2)
	u.Union(1, 2)

	root := u.Find(1)

	u.Union(1, 2) // already together
	u.Union(2, 1) // reversed order, still together

	if u.Find(1) != root || u.Find(2) != root {
		t.Error("Union of elements already in the same set changed their root")
	}
}

// Calling Find on something that was never registered should treat it as its own
// singleton and hand back the element itself, not the zero value.
func TestUnionFindFindWithoutMakeSet(t *testing.T) {
	u := NewUnionFind[int]()

	if got := u.Find(42); got != 42 {
		t.Errorf("Find(42) without MakeSet = %v, want 42 (lazily its own root)", got)
	}
}

// A set whose root is the zero value used to swallow any unregistered element.
// An unknown element should never end up sharing a root with a real zero-value
// set, since that would be a silent false positive.
func TestUnionFindLazyInitNoZeroValueCollision(t *testing.T) {
	u := NewUnionFind[int]()
	u.MakeSet(0) // a real set that happens to have the zero value as its root

	if u.Find(999) == u.Find(0) {
		t.Error("unregistered element collided with the zero-value set")
	}
}

// Union resolves its operands through Find, so it should register unknown
// elements on its own instead of choking on them.
func TestUnionFindUnionWithoutMakeSet(t *testing.T) {
	u := NewUnionFind[int]()

	u.Union(7, 8) // neither was MakeSet'd

	if u.Find(7) != u.Find(8) {
		t.Error("Union of unregistered elements should place them in the same set")
	}
}
