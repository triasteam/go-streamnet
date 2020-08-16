package streamwork

import (
	"fmt"
	"math"

	"github.com/phf/go-queue/queue"
	"github.com/triasteam/go-streamnet/types"
)

// StreamWork is the interface of dag
type StreamWork interface {
	getTotalOrder()
}

// Dag is the implement of StreamWork
// blocks (or transaction) is stored as a graph,
// a graph would be convert to a chain to output
type Dag struct {
	graph          map[types.Hash][]types.Hash
	revGraph       map[types.Hash][]types.Hash
	parentGraph    map[types.Hash]types.Hash
	revParentGraph map[types.Hash][]types.Hash

	parentScore map[types.Hash]float64

	lvlMap map[types.Hash]int64

	genesis types.Hash
}

// IfCovered test if a node is son of ancestor
func (d *Dag) IfCovered(block types.Hash, ancestor types.Hash, covered []types.Hash) bool {

	if d.revGraph[block] == nil {
		return false
	}

	if block == ancestor {
		return true
	}

	visited := []types.Hash{}
	fmt.Println("visited", visited)

	queue := &queue.Queue{}
	queue.Init()

	queue.PushBack(block)
	visited = append(visited, block)

	for queue.Len() > 0 {
		if h, ok := queue.PopFront().(types.Hash); ok {
			for _, e := range d.revGraph[h] {
				if e == ancestor {
					return true
				} else {
					if d.revGraph[e] != nil && !contains(visited, e) && !contains(covered, e) {
						queue.PushBack(e)
						visited = append(visited, e)
					}
				}
			}
		} else {

		}
	}
	return false
}

// DiffSet is ...
func (d *Dag) DiffSet(block types.Hash, parent types.Hash, covered []types.Hash) []types.Hash {
	if d.graph[block] == nil {
		return []types.Hash{}
	}

	ret := []types.Hash{}
	queue := &queue.Queue{}
	queue.Init()

	queue.PushBack(block)

	for queue.Len() > 0 {
		if h, ok := queue.PopFront().(types.Hash); ok {
			for _, e := range d.graph[h] {
				if d.graph[e] != nil && !contains(ret, e) && !d.IfCovered(e, parent, covered) {
					queue.PushBack(e)
				}
			}
			ret = append(ret, h)
		}
	}
	return ret
}

// GetMax returns ...
func (d *Dag) GetMax(start types.Hash) types.Hash {
	tmpMaxScore := float64(-1)
	s := types.Hash{}
	for _, block := range d.revParentGraph[start] {
		if d.parentScore[block] != 0 {
			if d.parentScore[block] > tmpMaxScore {
				tmpMaxScore = d.parentScore[block]
				s = block
			} else if d.parentScore[block] == tmpMaxScore {
				//按特定顺序，此处以block ascii 码大小
				if StringToFloat(s) > StringToFloat(block) {
					s = block
				}
			}
		}
	}

	return s
}

func (d *Dag) getPivot(start types.Hash) types.Hash {
	if &start == nil || d.graph[start] == nil {
		return types.Hash{}
	}

	for d.revParentGraph[start] != nil {
		s := d.GetMax(start)
		if &s == nil {
			return start
		}
		start = s
	}
	return start
}

// buildSubGraph returns ...
func (d *Dag) buildSubGraph(blocks []types.Hash) map[types.Hash][]types.Hash {
	subMap := make(map[types.Hash][]types.Hash)
	for _, h := range blocks {
		s := d.graph[h]
		ss := []types.Hash{}

		for _, hh := range s {
			if contains(blocks, hh) {
				ss = append(ss, hh)
			}
		}
		subMap[h] = ss
	}
	return subMap
}

// StreamWork returns ...
func (d *Dag) StreamWork(block types.Hash) []types.Hash {
	list := []types.Hash{}
	covered := []types.Hash{}
	if &block == nil || d.graph[block] == nil {
		return list
	}

	for {
		parent := d.parentGraph[block]
		subTopOrder := []types.Hash{}
		diff := d.DiffSet(block, parent, covered)
		for len(diff) != 0 {
			subGraph := d.buildSubGraph(diff)
			noBeforeInTmpGraph := []types.Hash{}
			for k, v := range subGraph {
				if len(v) != 0 {
					continue
				}
				noBeforeInTmpGraph = append(noBeforeInTmpGraph, k)
			}
			// TODO init lvlMap
			for _, s := range noBeforeInTmpGraph {
				if d.lvlMap[s] != 0 {
					d.lvlMap[s] = math.MaxInt64
				}
			}
			// TODO 2 按lvl从小到大排序，level一样按照字符串float大小排序
			SortByLvl(noBeforeInTmpGraph, d.lvlMap)
			subTopOrder = append(subTopOrder, noBeforeInTmpGraph...)
			diff = RemoveAll(diff, noBeforeInTmpGraph)
		}
		list = append(list, subTopOrder...)
		covered = append(covered, subTopOrder...)
		block = d.parentGraph[block]

		if len(d.parentGraph[block]) < 1 {
			break
		}
	}
	return list
}

// GetTotalOrder returns total order of the graph
func (d *Dag) GetTotalOrder() []types.Hash {
	pivot := d.getPivot(d.genesis)
	return d.StreamWork(pivot)
}

func contains(s []types.Hash, e types.Hash) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// StringToFloat make a types.Hash to float
func StringToFloat(s types.Hash) float64 {
	result := float64(0)
	for i := range s {
		result += float64(s[i]) * math.Pow10(i)
	}

	return result
}

// RemoveAll first cp a empty target types.Hash then append if not contained by another slice
func RemoveAll(obj []types.Hash, toRemove []types.Hash) []types.Hash {
	rtn := []types.Hash{}
	for _, s := range obj {
		if !contains(toRemove, s) {
			rtn = append(rtn, s)
		}
	}
	return rtn
}

// SortByLvl return sorted slice by lvl
func SortByLvl(obj []types.Hash, lvlMap map[types.Hash]int64) []types.Hash {
	length := len(obj)
	for i := 0; i < length; i++ {
		for j := i + 1; j < length; j++ {
			one := obj[i]
			ano := obj[j]
			oneLevel := lvlMap[one]
			anoLevel := lvlMap[ano]
			if oneLevel < anoLevel {
				tmp := one
				obj[i] = ano
				obj[j] = tmp
			} else if oneLevel == anoLevel {
				if StringToFloat(one) < StringToFloat(ano) {
					tmp := one
					obj[i] = ano
					obj[j] = tmp
				}
			}
		}
	}
	return obj
}
