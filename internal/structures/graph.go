package structures

import "github.com/google/uuid"

// NodeID is simply a UUID, used for unique nodes in a graph.
type NodeID uuid.UUID

// Helper interfaces

type Number interface {
	~int | ~int64 | ~uint | ~uint64 | ~float32 | ~float64
}

// UnweightedGraph represents a graph structure exposing read access to node connections.
// Implementations may or may not be safe for concurrent use.
type UnweightedGraph interface {
	GraphDistance
	GraphReader
	GraphWriter
}

type GraphWriter interface {
	// AddEdge returns an error if src/dst are non-existent in the graph.
	AddEdge(src, dst NodeID) error

	// AddNode returns an error if the node already exists in the graph.
	AddNode(node NodeID) error

	// RemoveEdge returns an error if the edge doesn't exist.
	RemoveEdge(src, dst NodeID) error

	// RemoveNode returns an error if the node doesn't exist.
	RemoveNode(node NodeID) error
}

type GraphReader interface {
	// Neighbors returns the IDs of all nodes directly connected to id.
	// Returns an empty slice if id has no neighbors.
	Neighbors(id NodeID) []NodeID
}

type GraphDistance interface {
	// DegreesOfSeparation returns the length of the shortest path between
	// src and dst, in number of hops. Returns -1 if dst is unreachable
	// from src.
	DegreesOfSeparation(src, dst NodeID) int
}

// -------------- Weighted Graph interface --------------

type WeightedGraphWriter[T Number] interface {
	// GraphWriter should take a default weight for the implementations
	// of AddEdge.
	GraphWriter

	// AddWeightedEdge returns an error if the edge (src, dst, weight) already exists
	// in the graph.
	AddWeightedEdge(src, dst NodeID, weight T) error
}

type WeightedGraphReader[T Number] interface {
	GraphReader

	// NeighborWeight returns the weight of the edge (src, dst).
	// Returns false if the edge doesn't exist.
	NeighborWeight(src, dst NodeID) (T, bool)
}

type WeightedGraphDistance[T Number] interface {
	GraphDistance

	// ShortestWeightedPath returns the shortest path from src to dst.
	// Returns -1 (in the numeric type specified) if the path is non-existent.
	ShortestWeightedPath(src, dst NodeID) T
}

type WeightedGraph[T Number] interface {
	WeightedGraphWriter[T]
	WeightedGraphReader[T]
	WeightedGraphDistance[T]
}
