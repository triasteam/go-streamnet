package dag

import (
	"container/list"

	"github.com/triasteam/go-streamnet/types"
)

func Update(graph map[types.Hash]types.Set, score map[types.Hash]float64, key types.Hash, weight float64) map[types.Hash]float64 {
	ret := score
	queue := list.New()

	queue.PushBack(key)
	visited := types.NewSet()
	visited.Add(key)

	for queue.Len() != 0 {
		h := queue.Front()
		queue.Remove(h)
		for _, e := range graph[h].List() {
			_, ok1 := graph[e]
			ok2 := visited.Has(e)
			if ok1 && !ok2 {
				queue.PushBack(e)
				visited.Add(e)
			}
		}
		if _, ok := ret[h]; !ok {
			ret[h] = weight
		} else {
			ret[h] = ret[h] + weight
		}

	}
	return ret
}

func UpdateParentScore(parentGraph map[types.Hash]types.Hash, parentScore map[types.Hash]float64, key types.Hash, weight float64) map[types.Hash]float64 {
	ret := parentScore
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
		start = parentGraph[start]
		_, ok = parentGraph[start]
	}

	return ret
}

func ComputeParentScore(parentGraph map[types.Hash]types.Hash, revParentGraph map[types.Hash]types.Set) map[types.Hash]float64 {
	ret := make(map[types.Hash]float64)

	for key, _ := range parentGraph {
		start := key
		visited := types.NewSet()
		for start != nil && !visited.Has(start) {
			if _, ok := ret[start]; ok {
				ret[start] = ret[start] + 1
			} else {
				ret[start] = 1.0
			}
			visited.Add(start)
			start = parentGraph[start]
		}
	}

	return ret
}

func compute(revGraph, graph map[types.Hash]types.Set, genesis types.Hash) map[types.Hash]float64 {
	ret := make(map[types.Hash]float64)
	queue := list.New()

	queue.PushBack(genesis)
	visited := types.NewSet()
	visited.Add(genesis)

	for queue.Len() != 0 {
		h := queue.Front()
		queue.Remove(h)
		if _, ok := revGraph[h]; ok {
			for e := range revGraph[h].List() {
				_, ok1 := revGraph[e]
				_, ok2 := graph[e]
				if (ok1 || ok2) && !visited.Has(e) {
					queue.PushBack(e)
					visited.Add(e)
				}
			}
		}
		ret = Update(graph, ret, h, 1.0)
	}
	return ret
}
