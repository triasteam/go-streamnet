package streamwork

import (
	"fmt"
	"math"

	"github.com/phf/go-queue/queue"
)

// StreamWork is the interface of dag
type StreamWork interface {
	getTotalOrder()
}

// Dag is the implement of StreamWork
// blocks (or transaction) is stored as a graph,
// a graph would be convert to a chain to output
type Dag struct {
	graph          map[string][]string
	revGraph       map[string][]string
	parentGraph    map[string]string
	revParentGraph map[string][]string

	parentScore map[string]float64

	lvlMap map[string]int64

	genesis string
}

// IfCovered test if a node is son of ancestor
func (d *Dag) IfCovered(block string, ancestor string, covered []string) bool {

	if d.revGraph[block] == nil {
		return false
	}

	if block == ancestor {
		return true
	}

	visited := []string{}
	fmt.Println("visited", visited)

	queue := &queue.Queue{}
	queue.Init()

	queue.PushBack(block)
	visited = append(visited, block)

	for queue.Len() > 0 {
		if h, ok := queue.PopFront().(string); ok {
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
func (d *Dag) DiffSet(block string, parent string, covered []string) []string {
	if d.graph[block] == nil {
		return []string{}
	}

	ret := []string{}
	queue := &queue.Queue{}
	queue.Init()

	queue.PushBack(block)

	for queue.Len() > 0 {
		if h, ok := queue.PopFront().(string); ok {
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
func (d *Dag) GetMax(start string) string {
	tmpMaxScore := float64(-1)
	s := ""
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

func (d *Dag) getPivot(start string) string {
	if start == "" || d.graph[start] == nil {
		return ""
	}

	for d.revParentGraph[start] != nil {
		s := d.GetMax(start)
		if s == "" {
			return start
		}
		start = s
	}
	return start
}

// buildSubGraph returns ...
func (d *Dag) buildSubGraph(blocks []string) map[string][]string {
	subMap := make(map[string][]string)
	for _, h := range blocks {
		s := d.graph[h]
		ss := []string{}

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
func (d *Dag) StreamWork(block string) []string {
	list := []string{}
	covered := []string{}
	if block == "" || d.graph[block] == nil {
		return list
	}

	for {
		parent := d.parentGraph[block]
		subTopOrder := []string{}
		diff := d.DiffSet(block, parent, covered)
		for len(diff) != 0 {
			subGraph := d.buildSubGraph(diff)
			noBeforeInTmpGraph := []string{}
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

		if d.parentGraph[block] == "" {
			break
		}
	}
	return list
}

// GetTotalOrder returns total order of the graph
func (d *Dag) GetTotalOrder() []string {
	pivot := d.getPivot(d.genesis)
	return d.StreamWork(pivot)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// StringToFloat make a string to float
func StringToFloat(s string) float64 {
	result := float64(0)
	for i := range s {
		result += float64(s[i]) * math.Pow10(i)
	}

	return result
}

// RemoveAll first cp a empty target string then append if not contained by another slice
func RemoveAll(obj []string, toRemove []string) []string {
	rtn := []string{}
	for _, s := range obj {
		if !contains(toRemove, s) {
			rtn = append(rtn, s)
		}
	}
	return rtn
}

// SortByLvl return sorted slice by lvl
func SortByLvl(obj []string, lvlMap map[string]int64) []string {
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
