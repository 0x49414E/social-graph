package structures

// bfsShortestPath is a pure BFS helper.
// Callers are responsible for ensuring `neighborsOf` is safe to call
// under their own concurrency model.
func bfsShortestPath(src, dst NodeID, neighborsOf func(NodeID) []NodeID) int {
	q := NewQueue[NodeID]()
	q.Push(src)

	dist := make(map[NodeID]int)
	dist[src] = 0

	for q.Len() != 0 {
		v, _ := q.Pop()

		if v == dst {
			return dist[v]
		}

		for _, neighbor := range neighborsOf(v) {
			if _, seen := dist[neighbor]; !seen {
				dist[neighbor] = dist[v] + 1
				q.Push(neighbor)
			}
		}
	}

	return -1
}
