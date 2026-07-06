package structures

type UnionFind[T comparable] struct {
	id   map[T]T
	rank map[T]int
}

func NewUnionFind[T comparable]() *UnionFind[T] {
	return &UnionFind[T]{
		id:   make(map[T]T),
		rank: make(map[T]int),
	}
}

// MakeSet registers x as a new singleton set if it isn't known yet.
func (u *UnionFind[T]) MakeSet(x T) {
	if _, exists := u.id[x]; !exists {
		u.id[x] = x
		u.rank[x] = 0
	}
}

// Find returns the representative root of Xs set, applying path compression
// along the way. If x was never registered it is lazily initialized as its own
// singleton set.
func (u *UnionFind[T]) Find(x T) T {
	if _, known := u.id[x]; !known {
		u.MakeSet(x)
	}
	if u.id[x] != x {
		u.id[x] = u.Find(u.id[x]) // path compression
	}
	return u.id[x]
}

// Union merges the sets containing a and b using union by rank: the shorter
// tree is hung under the taller one to keep the structure shallow. Elements
// not yet registered are lazily initialized by Find. It is a no-op if a and b
// already share a root.
func (u *UnionFind[T]) Union(a, b T) {
	rootA := u.Find(a)
	rootB := u.Find(b)

	if rootA == rootB {
		return
	}

	if u.rank[rootA] < u.rank[rootB] {
		rootA, rootB = rootB, rootA
	}
	u.id[rootB] = rootA
	if u.rank[rootA] == u.rank[rootB] {
		u.rank[rootA]++
	}
}
