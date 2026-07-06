package structures

import "testing"

func TestDegreesOfSeparationDirectNeighbor(t *testing.T) {
	g := NewUndirectedGraph()
	a, b := newID(), newID()
	mustAddNode(t, g, a, b)
	mustAddEdge(t, g, a, b)

	if got := g.DegreesOfSeparation(a, b); got != 1 {
		t.Errorf("DegreesOfSeparation(a, b) = %d, want 1", got)
	}
}

func TestDegreesOfSeparationMultiHop(t *testing.T) {
	g := NewUndirectedGraph()
	// Path: a - b - c - d
	a, b, c, d := newID(), newID(), newID(), newID()
	mustAddNode(t, g, a, b, c, d)
	mustAddEdge(t, g, a, b)
	mustAddEdge(t, g, b, c)
	mustAddEdge(t, g, c, d)

	if got := g.DegreesOfSeparation(a, d); got != 3 {
		t.Errorf("DegreesOfSeparation(a, d) = %d, want 3", got)
	}
}

func TestDegreesOfSeparationSameNode(t *testing.T) {
	g := NewUndirectedGraph()
	a := newID()
	mustAddNode(t, g, a)

	if got := g.DegreesOfSeparation(a, a); got != 0 {
		t.Errorf("DegreesOfSeparation(a, a) = %d, want 0", got)
	}
}

func TestDegreesOfSeparationDisconnected(t *testing.T) {
	g := NewUndirectedGraph()
	// Two separate components: a - b   and   c - d
	a, b, c, d := newID(), newID(), newID(), newID()
	mustAddNode(t, g, a, b, c, d)
	mustAddEdge(t, g, a, b)
	mustAddEdge(t, g, c, d)

	if got := g.DegreesOfSeparation(a, c); got != -1 {
		t.Errorf("DegreesOfSeparation(a, c) across components = %d, want -1", got)
	}
}

// A node with no edges at all can't reach anyone, so it should be -1 to everyone
// else, while the distance to itself is still 0.
func TestDegreesOfSeparationIsolatedNode(t *testing.T) {
	g := NewUndirectedGraph()
	// isolated has no edges. a and b form a separate little component off to the side.
	isolated, a, b := newID(), newID(), newID()
	mustAddNode(t, g, isolated, a, b)
	mustAddEdge(t, g, a, b)

	if got := g.DegreesOfSeparation(isolated, a); got != -1 {
		t.Errorf("DegreesOfSeparation(isolated, a) = %d, want -1", got)
	}
	if got := g.DegreesOfSeparation(isolated, isolated); got != 0 {
		t.Errorf("DegreesOfSeparation(isolated, isolated) = %d, want 0", got)
	}
}

// A cyclic graph shouldn't trap BFS in an infinite loop, and it should still
// come back with the shortest distance.
func TestDegreesOfSeparationWithCycle(t *testing.T) {
	g := NewUndirectedGraph()
	// Triangle a - b - c - a, plus a tail c - d.
	a, b, c, d := newID(), newID(), newID(), newID()
	mustAddNode(t, g, a, b, c, d)
	mustAddEdge(t, g, a, b)
	mustAddEdge(t, g, b, c)
	mustAddEdge(t, g, c, a)
	mustAddEdge(t, g, c, d)

	// Shortest path from a to d is a, c, d at 2 hops, not a, b, c, d at 3.
	if got := g.DegreesOfSeparation(a, d); got != 2 {
		t.Errorf("DegreesOfSeparation(a, d) in cyclic graph = %d, want 2", got)
	}
}
