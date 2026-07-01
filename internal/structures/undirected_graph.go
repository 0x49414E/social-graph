package structures

import (
	"fmt"
	"sync"
)

const (
	errNodeNotFound      = "node %s does not exist"
	errEdgeNotFound      = "edge %s does not exist"
	errEdgeAlreadyExists = "edge (%s,%s) already exists"
	errNodeAlreadyExists = "node %s already exists"
	errSelfLoop          = "self loop occurred at node %s"
)

// Edge is a reference to another node in the UndirectedGraph adjacency list.
type Edge struct {
	To NodeID
}

// UndirectedGraph is a thread-safe graph implementation.
// It stores an adjacency list mapping each node to its outgoing edges.
type UndirectedGraph struct {
	mu    sync.RWMutex
	edges map[NodeID]map[NodeID]struct{}
}

// Neighbors returns the IDs of all nodes directly connected to the specfied id.
func (g *UndirectedGraph) Neighbors(id NodeID) []NodeID {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.neighborsUnsafe(id)
}

// DegreesOfSeparation returns the shortest path length between src and dst
// using BFS. Returns -1 if dst is unreachable from src.
func (g *UndirectedGraph) DegreesOfSeparation(src, dst NodeID) int {
	// Take a single read lock so the graph doesn't change mid-BFS.
	g.mu.RLock()
	defer g.mu.RUnlock()

	return bfsShortestPath(src, dst, g.neighborsUnsafe)
}

// AddEdge adds the edge (src, dst) and the edge (dst, src) to the graph.
func (g *UndirectedGraph) AddEdge(src, dst NodeID) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if src == dst {
		return fmt.Errorf(errSelfLoop, src)
	}

	srcEdges, ok1 := g.edges[src]
	_, ok2 := g.edges[dst]
	if !ok1 {
		return fmt.Errorf(errNodeNotFound, src)
	}
	if !ok2 {
		return fmt.Errorf(errNodeNotFound, dst)
	}

	if _, exists := srcEdges[dst]; exists {
		return fmt.Errorf(errEdgeAlreadyExists, src, dst)
	}

	srcEdges[dst] = struct{}{}
	g.edges[dst][src] = struct{}{}

	return nil
}

func (g *UndirectedGraph) AddNode(a NodeID) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.edges[a]; exists {
		return fmt.Errorf(errNodeAlreadyExists, a)
	}

	g.edges[a] = make(map[NodeID]struct{})

	return nil
}

func (g *UndirectedGraph) RemoveNode(id NodeID) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	neighbors, exists := g.edges[id]
	if !exists {
		return fmt.Errorf(errNodeNotFound, id)
	}

	for neighbor := range neighbors {
		delete(g.edges[neighbor], id)
	}

	delete(g.edges, id)

	return nil
}

func (g *UndirectedGraph) RemoveEdge(src, dst NodeID) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if _, exists := g.edges[src]; !exists {
		return fmt.Errorf(errNodeNotFound, src)
	}

	delete(g.edges[src], dst)
	delete(g.edges[dst], src)

	return nil
}

// neighborsUnsafe returns the neighbor IDs of id without acquiring any lock.
// Callers must hold at least a read lock (g.mu.RLock) before calling this.
func (g *UndirectedGraph) neighborsUnsafe(id NodeID) []NodeID {
	edges := g.edges[id]
	result := make([]NodeID, 0, len(edges))
	for i := range edges {
		result = append(result, i)
	}
	return result
}
