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

func (u *UnionFind[T]) Find(x T) T {
	if u.id[x] != x {
		u.id[x] = u.Find(u.id[x]) // path compression
	}
	return u.id[x]
}

func (u *UnionFind[T]) Union(a, b T) {
	rootA := u.Find(a)
	rootB := u.Find(b)

	if rootA == rootB {
		return
	}

	// union by rank: hang the shorter tree under the taller one
	if u.rank[rootA] < u.rank[rootB] {
		rootA, rootB = rootB, rootA
	}
	u.id[rootB] = rootA
	if u.rank[rootA] == u.rank[rootB] {
		u.rank[rootA]++
	}
}
