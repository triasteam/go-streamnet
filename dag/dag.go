package dag

import (
	"sync"

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

		d.pivotChain.Add(s)

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
