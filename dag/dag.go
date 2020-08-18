package dag

import (
	"fmt"
	"math"
	"sync"

	"github.com/phf/go-queue/queue"
	"github.com/triasteam/go-streamnet/store"
	"github.com/triasteam/go-streamnet/types"
	"github.com/triasteam/go-streamnet/utils"
)

// Dag is the most important struct in the whole procedure.
type Dag struct {
	graph       map[types.Hash]types.Set
	parentGraph map[types.Hash]types.Hash

	revGraph       map[types.Hash]types.Set
	parentRevGraph map[types.Hash]types.Set

	degrees map[types.Hash]int64

	score       map[types.Hash]float64
	parentScore map[types.Hash]float64
	freshScore  bool

	topOrder          map[int64]types.Set
	topOrderStreaming map[int64]types.Set
	totalDepth        int

	subGraph          map[types.Hash]types.Set
	subRevGraph       map[types.Hash]types.Set
	subParentGraph    map[types.Hash]types.Hash
	subParentRevGraph map[types.Hash]types.Set

	levels  map[types.Hash]int64
	namemap map[types.Hash]string

	store *store.Storage

	pivotChain types.List

	graphLock sync.RWMutex

	ancestors *types.Stack

	genesis types.Hash
}

// Init return a new Dag struct.
func Init(db *store.Storage) *Dag {
	return &Dag{
		graph:             make(map[types.Hash]types.Set),
		parentGraph:       make(map[types.Hash]types.Hash),
		revGraph:          make(map[types.Hash]types.Set),
		parentRevGraph:    make(map[types.Hash]types.Set),
		score:             make(map[types.Hash]float64),
		parentScore:       make(map[types.Hash]float64),
		freshScore:        false,
		degrees:           make(map[types.Hash]int64),
		topOrder:          make(map[int64]types.Set),
		topOrderStreaming: make(map[int64]types.Set),
		subGraph:          make(map[types.Hash]types.Set),
		subRevGraph:       make(map[types.Hash]types.Set),
		subParentGraph:    make(map[types.Hash]types.Hash),
		subParentRevGraph: make(map[types.Hash]types.Set),
		levels:            make(map[types.Hash]int64),
		namemap:           make(map[types.Hash]string),
		totalDepth:        0,
		store:             db,
		pivotChain:        types.List{},
		ancestors:         types.NewStack(),
	}
}

// Close will free all the resources.
func (d *Dag) Close() {

}

// Add puts one tx hash into dag.
func (d *Dag) Add(key types.Hash, value *types.Transaction) bool {
	trunk := value.GetTrunkTransactionHash()
	branch := value.GetBranchTransactionHash()

	d.graphLock.Lock()
	defer d.graphLock.Unlock()

	d.updateGraph(key, trunk, branch)

	d.updateTopologicalOrder(key, trunk, branch)

	d.updateScore(key)

	return true
}

func (d *Dag) updateGraph(key, trunk, branch types.Hash) {
	// Approve graph
	if _, ok := d.graph[key]; !ok {
		d.graph[key] = types.NewSet()
	}
	d.graph[key].Add(trunk)
	d.graph[key].Add(branch)

	//parentGraph
	d.parentGraph[key] = trunk

	// Approvee graph
	if _, ok := d.revGraph[trunk]; !ok {
		d.revGraph[trunk] = types.NewSet()
	}
	d.revGraph[trunk].Add(key)
	if _, ok := d.revGraph[branch]; !ok {
		d.revGraph[branch] = types.NewSet()
	}
	d.revGraph[branch].Add(key)

	if _, ok := d.parentRevGraph[trunk]; !ok {
		d.parentRevGraph[trunk] = types.NewSet()
	}
	d.parentRevGraph[trunk].Add(key)

	// update degrees
	if _, ok := d.degrees[key]; !ok || d.degrees[key] == 0 {
		d.degrees[key] = 2
	}

	if _, ok := d.degrees[trunk]; !ok {
		d.degrees[trunk] = 0
	}
	if _, ok := d.degrees[branch]; !ok {
		d.degrees[branch] = 0
	}
}

func (d *Dag) updateTopologicalOrder(key, trunk, branch types.Hash) {
	if len(d.topOrderStreaming) == 0 {
		d.topOrderStreaming[1] = types.NewSet()
		d.topOrderStreaming[1].Add(key)
		d.levels[key] = 1
		d.topOrderStreaming[0] = types.NewSet()
		d.topOrderStreaming[0].Add(trunk)
		d.topOrderStreaming[0].Add(branch)
		d.totalDepth = 1
		return
	} else {
		// TODO: check trunk or branch exist !!!!!!!!!
		// Or we won't call Add if trunk or branch not exist!!!!
		trunkLevel := d.levels[trunk]
		branchLevel := d.levels[branch]
		lvl := utils.Min(trunkLevel, branchLevel) + 1
		if _, ok := d.topOrderStreaming[lvl]; !ok {
			d.topOrderStreaming[lvl] = types.NewSet()
			d.totalDepth++
		}
		d.topOrderStreaming[lvl].Add(key)
		d.levels[key] = lvl
	}
}

func (d *Dag) updateScore(key types.Hash) {
	scoreAlg := "CUM_WEIGHT"

	if scoreAlg == "CUM_WEIGHT" {
		CumulateWeight{}.Update(d.graph, d.score, key, 1)
		CumulateWeight{}.UpdateParentScore(d.parentGraph, d.parentScore, key, 1)

	} else if scoreAlg == "KATZ" {
		d.score[key] = 1.0 / (float64(len(d.score)) + 1.0)
		centrality := NewKatz(d.graph, d.revGraph, d.score, 0.5)
		d.score = centrality.Compute()
		CumulateWeight{}.UpdateParentScore(d.parentGraph, d.parentScore, key, 1)
	}
	d.freshScore = false
}

func computeToplogicalOrder() {

}

func (d *Dag) GetPivotalHash(depth int) types.Hash {
	var ret types.Hash
	d.buildPivotChain()
	if depth == -1 || depth >= d.pivotChain.Length() {
		set := d.topOrderStreaming[1]
		if set.IsEmpty() {
			return types.NewHash(nil)
		}
		ret = set.List()[0]
		return ret
	}

	ret = d.pivotChain.Index(d.pivotChain.Length() - depth - 1)
	return ret
}

func (d *Dag) buildPivotChain() {
	d.pivotChain = d.PivotChain(d.GetGenesis())
}

func (d *Dag) PivotChain(start types.Hash) types.List {
	d.graphLock.RLock()
	defer d.graphLock.Unlock()

	list := types.List{}

	if _, ok := d.graph[start]; start == types.NewHash(nil) || !ok {
		return list
	}

	list.Append(start)

	v, ok := d.parentRevGraph[start]
	for ok && !v.IsEmpty() {
		s := d.getMax(v)
		if s == types.NewHash(nil) {
			return list
		}

		d.pivotChain.Append(s)

		start = s
		v, ok = d.parentRevGraph[start]
	}

	return list
}

func (d *Dag) GetGenesis() types.Hash {
	if d.ancestors != nil && !d.ancestors.Empty() {
		return d.ancestors.Peek()
	}

	for key, v := range d.parentGraph {
		if _, ok := d.parentGraph[v]; !ok {
			return key
		}
	}
	return types.NewHash(nil)
}

func (d *Dag) getMax(set types.Set) types.Hash {
	tmpMaxScore := -1.0
	s := types.NewHash(nil)

	for _, block := range set.List() {
		if v, ok := d.parentScore[block]; ok {
			if v > tmpMaxScore {
				tmpMaxScore = v
				s = block
			} else if v == tmpMaxScore {
				sStr := s.String()
				blockStr := block.String()
				if sStr < blockStr {
					s = block
				}
			}
		}
	}
	return s
}

func (d *Dag) GetLastPivot(start types.Hash) types.Hash {
	d.graphLock.RLock()

	if _, ok := d.graph[start]; start == types.NewHash(nil) || !ok {
		return types.NewHash(nil)
	}
	v, ok := d.parentRevGraph[start]
	for ok && !v.IsEmpty() {
		s := d.getMax(v)
		if s == types.NewHash(nil) {
			return start
		}
		start = s
		v, ok = d.parentRevGraph[start]
	}
	return start

}

func (d *Dag) BuildGraph() {
}

func (d *Dag) ComputeScore() {
}

func (d *Dag) GetChild(block types.Hash) types.Set {
	if v, ok := d.revGraph[block]; ok {
		return v
	}

	return types.NewSet()
}

func (d *Dag) Contains(key types.Hash) bool {
	_, ok := d.graph[key]
	return ok
}

func (d *Dag) GetScore(key types.Hash) float64 {
	d.graphLock.RLock()
	defer d.graphLock.Unlock()

	score, ok := d.score[key]
	if ok {
		return score
	} else {
		return 0.0
	}
}

// IfCovered test if a node is son of ancestor
func (d *Dag) IfCovered(block types.Hash, ancestor types.Hash, covered []types.Hash) bool {

	if _, ok := d.revGraph[block]; !ok {
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
			if set, ok := d.revGraph[h]; ok {
				for _, e := range set.List() {
					if e == ancestor {
						return true
					} else {
						if _, ok := d.revGraph[e]; ok && !contains(visited, e) && !contains(covered, e) {
							queue.PushBack(e)
							visited = append(visited, e)
						}
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
	if _, ok := d.graph[block]; !ok {
		return []types.Hash{}
	}

	ret := []types.Hash{}
	queue := &queue.Queue{}
	queue.Init()

	queue.PushBack(block)

	for queue.Len() > 0 {
		if h, ok := queue.PopFront().(types.Hash); ok {
			if set, ok := d.graph[h]; ok {
				for _, e := range set.List() {
					if _, ok := d.graph[e]; ok && !contains(ret, e) && !d.IfCovered(e, parent, covered) {
						queue.PushBack(e)
					}
				}
				ret = append(ret, h)
			}
		}
	}
	return ret
}

// GetMax returns ...
func (d *Dag) GetMax(start types.Hash) types.Hash {
	tmpMaxScore := float64(-1)
	s := types.Hash{}
	if set, ok := d.parentRevGraph[start]; ok {
		for _, block := range set.List() {
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
	}
	return s
}

func (d *Dag) getPivot(start types.Hash) types.Hash {
	if _, ok := d.graph[start]; !ok || &start == nil {
		return types.Hash{}
	}

	set, ok := d.parentRevGraph[start]
	for ok && !set.IsEmpty() {
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
		s, ok := d.graph[h]
		if !ok {
			continue
		}

		ss := []types.Hash{}

		for _, hh := range s.List() {
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
	if set, ok := d.graph[block]; !ok || set.IsEmpty() || &block == nil {
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
				if d.levels[s] != 0 {
					d.levels[s] = math.MaxInt64
				}
			}
			// TODO 2 按lvl从小到大排序，level一样按照字符串float大小排序
			SortByLvl(noBeforeInTmpGraph, d.levels)
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
