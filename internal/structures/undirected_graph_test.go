package structures

import (
	"sync"
	"testing"

	"github.com/google/uuid"
)

// newID hands back a fresh, unique NodeID.
func newID() NodeID {
	return NodeID(uuid.New())
}

// contains reports whether target shows up in ids.
func contains(ids []NodeID, target NodeID) bool {
	for _, id := range ids {
		if id == target {
			return true
		}
	}
	return false
}

func TestAddEdgeSymmetry(t *testing.T) {
	g := NewUndirectedGraph()
	a, b := newID(), newID()
	mustAddNode(t, g, a, b)

	if err := g.AddEdge(a, b); err != nil {
		t.Fatalf("AddEdge(a, b) unexpected error: %v", err)
	}

	// The graph is undirected, so the edge should show up from both ends.
	if !contains(g.Neighbors(a), b) {
		t.Error("Neighbors(a) should contain b")
	}
	if !contains(g.Neighbors(b), a) {
		t.Error("Neighbors(b) should contain a")
	}
}

func TestAddEdgeSelfLoop(t *testing.T) {
	g := NewUndirectedGraph()
	a := newID()
	mustAddNode(t, g, a)

	if err := g.AddEdge(a, a); err == nil {
		t.Error("AddEdge(a, a) should return a self-loop error")
	}
}

func TestAddEdgeDuplicate(t *testing.T) {
	g := NewUndirectedGraph()
	a, b := newID(), newID()
	mustAddNode(t, g, a, b)

	if err := g.AddEdge(a, b); err != nil {
		t.Fatalf("first AddEdge unexpected error: %v", err)
	}
	if err := g.AddEdge(a, b); err == nil {
		t.Error("duplicate AddEdge(a, b) should return an error")
	}
}

func TestAddEdgeMissingNode(t *testing.T) {
	t.Run("missing src", func(t *testing.T) {
		g := NewUndirectedGraph()
		a, b := newID(), newID()
		mustAddNode(t, g, b) // only b exists
		if err := g.AddEdge(a, b); err == nil {
			t.Error("AddEdge with missing src should return an error")
		}
	})

	t.Run("missing dst", func(t *testing.T) {
		g := NewUndirectedGraph()
		a, b := newID(), newID()
		mustAddNode(t, g, a) // only a exists
		if err := g.AddEdge(a, b); err == nil {
			t.Error("AddEdge with missing dst should return an error")
		}
	})
}

// Removing a node should also wipe every back reference to it, so none of its
// old neighbors are left pointing at a node that no longer exists.
func TestRemoveNodeCleansCrossReferences(t *testing.T) {
	g := NewUndirectedGraph()
	center, a, b := newID(), newID(), newID()
	mustAddNode(t, g, center, a, b)
	mustAddEdge(t, g, center, a)
	mustAddEdge(t, g, center, b)

	if err := g.RemoveNode(center); err != nil {
		t.Fatalf("RemoveNode unexpected error: %v", err)
	}

	// Its old neighbors shouldn't mention it anymore.
	if contains(g.Neighbors(a), center) {
		t.Error("Neighbors(a) still references the removed node")
	}
	if contains(g.Neighbors(b), center) {
		t.Error("Neighbors(b) still references the removed node")
	}

	// Removing it a second time should fail, since it's already gone.
	if err := g.RemoveNode(center); err == nil {
		t.Error("RemoveNode on an already-removed node should return an error")
	}
}

func TestRemoveEdgeBothDirections(t *testing.T) {
	g := NewUndirectedGraph()
	a, b := newID(), newID()
	mustAddNode(t, g, a, b)
	mustAddEdge(t, g, a, b)

	if err := g.RemoveEdge(a, b); err != nil {
		t.Fatalf("RemoveEdge unexpected error: %v", err)
	}

	if contains(g.Neighbors(a), b) {
		t.Error("Neighbors(a) should not contain b after RemoveEdge")
	}
	if contains(g.Neighbors(b), a) {
		t.Error("Neighbors(b) should not contain a after RemoveEdge")
	}
}

// UndirectedGraph claims to be thread-safe, so hammering it from lots of
// goroutines at once shouldn't trip the race detector. Run this with
// `go test -race`. It doesn't assert a specific outcome, since the ordering is
// non-deterministic; its whole job is to surface unsynchronized access.
func TestUndirectedGraphConcurrentAccess(t *testing.T) {
	g := NewUndirectedGraph()

	ids := make([]NodeID, 50)
	for i := range ids {
		ids[i] = newID()
		mustAddNode(t, g, ids[i])
	}

	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			a, b := ids[i%len(ids)], ids[(i+1)%len(ids)]

			// Mix writers and readers so the RWMutex gets exercised both ways.
			// Errors are ignored on purpose: with concurrent writers an edge may
			// already exist, and that's fine here.
			_ = g.AddEdge(a, b)
			_ = g.Neighbors(a)
			_ = g.DegreesOfSeparation(a, b)
			_ = g.RemoveEdge(a, b)
		}(i)
	}
	wg.Wait()
}

// mustAddNode adds each id and fails the test if anything goes wrong.
func mustAddNode(t *testing.T, g *UndirectedGraph, ids ...NodeID) {
	t.Helper()
	for _, id := range ids {
		if err := g.AddNode(id); err != nil {
			t.Fatalf("AddNode(%v) unexpected error: %v", id, err)
		}
	}
}

// mustAddEdge wires up src and dst and fails the test if anything goes wrong.
func mustAddEdge(t *testing.T, g *UndirectedGraph, src, dst NodeID) {
	t.Helper()
	if err := g.AddEdge(src, dst); err != nil {
		t.Fatalf("AddEdge(%v, %v) unexpected error: %v", src, dst, err)
	}
}
