package dag

import (
	"github.com/triasteam/go-streamnet/types"
)

type CumulateWeight struct{}

// Update adds wight to score in a DFS(深度优先遍历) way
func (cw CumulateWeight) Update(graph map[types.Hash]types.Set, score map[types.Hash]float64, key types.Hash, weight float64) {
	queue := types.List{}

	queue.Append(key)
	visited := types.NewSet()
	visited.Add(key)

	for queue.Length() != 0 {
		h := queue.RemoveAtIndex(0)
		for _, e := range graph[h].List() {
			_, ok1 := graph[e]
			ok2 := visited.Has(e)
			if ok1 && !ok2 {
				queue.Append(e)
				visited.Add(e)
			}
		}
		if _, ok := score[h]; !ok {
			score[h] = 0.0
		}
		score[h] = score[h] + weight
	}
}

// UpdateParentScore adds weight to score in a linked list way.
func (cw CumulateWeight) UpdateParentScore(parentGraph map[types.Hash]types.Hash, parentScore map[types.Hash]float64, key types.Hash, weight float64) {
	start := key
	visited := types.NewSet()

	_, ok := parentGraph[start]
	for ok {
		if visited.Has(start) {
			//log.error("Circle exist: " + start)
			break
		} else {
			visited.Add(start)
		}

		if _, o := parentScore[start]; !o {
			parentScore[start] = 0.0
		}
		parentScore[start] = parentScore[start] + weight
		start, ok = parentGraph[start]
	}
}

func (cw CumulateWeight) ComputeParentScore(parentGraph map[types.Hash]types.Hash, revParentGraph map[types.Hash]types.Set) map[types.Hash]float64 {
	ret := make(map[types.Hash]float64)

	for key, _ := range parentGraph {
		start := key
		visited := types.NewSet()
		ok := true
		for ok && !visited.Has(start) {
			if _, o := ret[start]; o {
				ret[start] = ret[start] + 1
			} else {
				ret[start] = 1.0
			}
			visited.Add(start)
			start, ok = parentGraph[start]
		}
	}

	return ret
}

func (cw CumulateWeight) Compute(revGraph, graph map[types.Hash]types.Set, genesis types.Hash) map[types.Hash]float64 {
	ret := make(map[types.Hash]float64)
	queue := types.List{}

	queue.Append(genesis)
	visited := types.NewSet()
	visited.Add(genesis)

	for queue.Length() != 0 {
		h := queue.RemoveAtIndex(0)
		if _, ok := revGraph[h]; ok {
			for _, e := range revGraph[h].List() {
				_, ok1 := revGraph[e]
				_, ok2 := graph[e]
				if (ok1 || ok2) && !visited.Has(e) {
					queue.Append(e)
					visited.Add(e)
				}
			}
		}
		cw.Update(graph, ret, h, 1.0)
	}
	return ret
}
