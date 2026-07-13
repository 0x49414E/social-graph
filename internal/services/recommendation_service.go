package services

import (
	"sort"

	"social-graph/internal/structures"
)

type Recommendation struct {
	User          structures.NodeID
	Score         float64
	MutualFriends int
}

type RecommendationService struct {
	g structures.UnweightedGraph
}

// RecommendFriends suggests up to `limit` new friends for `user`, drawn from the
// friends-of-friends the user isn't already connected to. Each candidate is
// scored by the Jaccard similarity of their friend set with the user's, and ties
// break by NodeID so the output is deterministic. A limit <= 0 means no cap.
func (r *RecommendationService) RecommendFriends(user structures.NodeID, limit int) []Recommendation {
	res := make([]Recommendation, 0)

	userFriends := r.g.Neighbors(user)

	// Never recommend the user themselves or someone they're already friends with.
	excluded := make(map[structures.NodeID]struct{}, len(userFriends)+1)
	excluded[user] = struct{}{}
	for _, f := range userFriends {
		excluded[f] = struct{}{}
	}

	userSet := toSet(userFriends)

	// Walk the friends-of-friends and score each unique candidate once.
	scored := make(map[structures.NodeID]struct{})
	for _, friend := range userFriends {
		for _, candidate := range r.g.Neighbors(friend) {
			if _, skip := excluded[candidate]; skip {
				continue
			}
			if _, seen := scored[candidate]; seen {
				continue
			}
			scored[candidate] = struct{}{}

			// One pass over the friend sets gives us both the similarity score
			// and the mutual-friend count.
			score, mutual := structures.JaccardWithOverlap(userSet, toSet(r.g.Neighbors(candidate)))
			res = append(res, Recommendation{
				User:          candidate,
				Score:         score,
				MutualFriends: mutual,
			})
		}
	}

	// Rank by score descending, breaking ties by NodeID for a stable order.
	sort.Slice(res, func(i, j int) bool {
		if res[i].Score != res[j].Score {
			return res[i].Score > res[j].Score
		}
		return lessNodeID(res[i].User, res[j].User)
	})

	if limit > 0 && len(res) > limit {
		res = res[:limit]
	}
	return res
}

// toSet turns a neighbor slice into a set so it can be fed to Jaccard.
func toSet(ids []structures.NodeID) map[structures.NodeID]struct{} {
	s := make(map[structures.NodeID]struct{}, len(ids))
	for _, id := range ids {
		s[id] = struct{}{}
	}
	return s
}

// lessNodeID gives NodeIDs a stable ordering so tied scores rank the same way on
// every run, instead of following Go's random map iteration.
func lessNodeID(a, b structures.NodeID) bool {
	ab, bb := [16]byte(a), [16]byte(b)
	for i := range ab {
		if ab[i] != bb[i] {
			return ab[i] < bb[i]
		}
	}
	return false
}
