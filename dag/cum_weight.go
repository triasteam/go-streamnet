package dag

import (
	"log"

	"github.com/triasteam/go-streamnet/types"
)

type CumulateWeight struct{}

// Update adds wight to score in a DFS(深度优先遍历) way
func (cw CumulateWeight) UpdateScore(graph map[types.Hash]types.Set, score map[types.Hash]float64, key types.Hash, weight float64) {
	queue := types.List{}
	queue.Append(key)

	visited := types.NewSet()
	visited.Add(key)

	for queue.Length() != 0 {
		cur := queue.RemoveAtIndex(0)
		for _, parent := range graph[cur].List() {
			_, exist := graph[parent]
			visit := visited.Has(parent)
			if exist && !visit {
				queue.Append(parent)
				visited.Add(parent)
			}
		}
		if _, exist := score[cur]; !exist {
			score[cur] = 0.0
		}
		score[cur] = score[cur] + weight
	}
}

// UpdateParentScore adds weight to score in a linked list way.
func (cw CumulateWeight) UpdateTrunkScore(parentGraph map[types.Hash]types.Hash, parentScore map[types.Hash]float64, key types.Hash, weight float64) {
	start := key
	visited := types.NewSet()

	_, ok := parentGraph[start]
	for ok {
		if visited.Has(start) {
			log.Printf("Error!!! Circle exist: %s\n", start)
			break
		} else {
			visited.Add(start)
		}

		if _, exist := parentScore[start]; !exist {
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
		cw.UpdateScore(graph, ret, h, 1.0)
	}
	return ret
}
