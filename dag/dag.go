package dag

import (
	"sync"

	"github.com/triasteam/go-streamnet/store"
	"github.com/triasteam/go-streamnet/types"
	"github.com/triasteam/go-streamnet/utils"
)

// Dag is the most important struct in the whole procedure.
type Dag struct {
	graph       map[types.Hash]types.Set  // parents of one node
	parentGraph map[types.Hash]types.Hash // trunk parent of one node

	revGraph       map[types.Hash]types.Set // children of one node
	parentRevGraph map[types.Hash]types.Set // children of one trunk node

	degrees map[types.Hash]int64 // every node's degree, for streamwork

	score       map[types.Hash]float64 // every node's score
	parentScore map[types.Hash]float64 // every trunk node's score
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
func (dag *Dag) Init(db *store.Storage) {
	dag.graph = make(map[types.Hash]types.Set)
	dag.parentGraph = make(map[types.Hash]types.Hash)
	dag.revGraph = make(map[types.Hash]types.Set)
	dag.parentRevGraph = make(map[types.Hash]types.Set)
	dag.score = make(map[types.Hash]float64)
	dag.parentScore = make(map[types.Hash]float64)
	dag.freshScore = false
	dag.degrees = make(map[types.Hash]int64)
	dag.topOrder = make(map[int64]types.Set)
	dag.topOrderStreaming = make(map[int64]types.Set)
	dag.subGraph = make(map[types.Hash]types.Set)
	dag.subRevGraph = make(map[types.Hash]types.Set)
	dag.subParentGraph = make(map[types.Hash]types.Hash)
	dag.subParentRevGraph = make(map[types.Hash]types.Set)
	dag.levels = make(map[types.Hash]int64)
	dag.namemap = make(map[types.Hash]string)
	dag.totalDepth = 0
	dag.store = db
	dag.pivotChain = types.List{}
	dag.ancestors = types.NewStack()

	// TODO: init dag with data stored in database if restart.
	dag.load()
}

// Close will free all the resources.
func (d *Dag) Close() {

}

// Add puts one tx hash into dag.
func (d *Dag) Add(key types.Hash, tx *types.Transaction) error {
	trunk := tx.GetTrunkTransactionHash()
	branch := tx.GetBranchTransactionHash()

	d.graphLock.Lock()
	defer d.graphLock.Unlock()

	d.updateGraph(key, trunk, branch)

	d.updateTopologicalOrder(key, trunk, branch)

	d.updateScore(key)

	return nil
}

func (d *Dag) updateGraph(key, trunk, branch types.Hash) {
	// Approve graph
	if _, exist := d.graph[key]; !exist {
		d.graph[key] = types.NewSet()
	}
	d.graph[key].Add(trunk)
	d.graph[key].Add(branch)

	//parentGraph
	d.parentGraph[key] = trunk

	// Approvee graph
	if _, exist := d.revGraph[trunk]; !exist {
		d.revGraph[trunk] = types.NewSet()
	}
	d.revGraph[trunk].Add(key)
	if _, exist := d.revGraph[branch]; !exist {
		d.revGraph[branch] = types.NewSet()
	}
	d.revGraph[branch].Add(key)

	if _, exist := d.parentRevGraph[trunk]; !exist {
		d.parentRevGraph[trunk] = types.NewSet()
	}
	d.parentRevGraph[trunk].Add(key)

	// update degrees
	if _, exist := d.degrees[key]; !exist || d.degrees[key] == 0 {
		d.degrees[key] = 2
	}
	if _, exist := d.degrees[trunk]; !exist {
		d.degrees[trunk] = 0
	}
	if _, exist := d.degrees[branch]; !exist {
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
		level := utils.Min(trunkLevel, branchLevel) + 1
		if _, exist := d.topOrderStreaming[level]; !exist {
			d.topOrderStreaming[level] = types.NewSet()
			d.totalDepth++
		}
		d.topOrderStreaming[level].Add(key)
		d.levels[key] = level
	}
}

func (d *Dag) updateScore(key types.Hash) {
	// todo: use config to choose score algorithm.
	scoreAlg := "CUM_WEIGHT"

	if scoreAlg == "CUM_WEIGHT" {
		CumulateWeight{}.UpdateScore(d.graph, d.score, key, 1)
		CumulateWeight{}.UpdateTrunkScore(d.parentGraph, d.parentScore, key, 1)

	} else if scoreAlg == "KATZ" {
		d.score[key] = 1.0 / (float64(len(d.score)) + 1.0)
		centrality := NewKatz(d.graph, d.revGraph, d.score, 0.5)
		d.score = centrality.Compute()
		CumulateWeight{}.UpdateTrunkScore(d.parentGraph, d.parentScore, key, 1)
	}
	d.freshScore = false
}

// GetPivotalHash returns the pivot hash of dag from genesis. depth is 0, 1, 2...
func (d *Dag) GetPivotalHash(depth int) types.Hash {
	var ret types.Hash
	d.BuildPivotChain()
	if depth == -1 || depth >= d.pivotChain.Length() {
		set := d.topOrderStreaming[1]
		if set.IsEmpty() {
			return types.NilHash
		}
		ret = set.List()[0]
		return ret
	}

	ret = d.pivotChain.Index(d.pivotChain.Length() - depth - 1)
	return ret
}

func (d *Dag) BuildPivotChain() {
	d.graphLock.RLock()
	defer d.graphLock.RUnlock()
	d.pivotChain = d.getPivotChainFrom(d.GetGenesis())
}

// get the pivot chain after the giving node.
func (d *Dag) getPivotChainFrom(start types.Hash) types.List {
	/*d.graphLock.RLock()
	defer d.graphLock.RUnlock()*/

	list := types.List{}

	if _, exist := d.graph[start]; start == types.NilHash || !exist {
		return list
	}

	list.Append(start)

	children, ok := d.parentRevGraph[start]
	for ok && !children.IsEmpty() {
		child := d.getMax(children)
		if child == types.NilHash {
			return list
		}

		list.Append(child)

		start = child
		children, ok = d.parentRevGraph[start]
	}

	return list
}

// todo: use config to get genesis directly; if genesis_forward is implemented, it should be as following.
func (d *Dag) GetGenesis() types.Hash {
	if d.ancestors != nil && !d.ancestors.Empty() {
		return d.ancestors.Peek()
	}

	// find one node whose trunk parent is not in parentGraph
	for key, trunk := range d.parentGraph {
		if _, ok := d.parentGraph[trunk]; !ok {
			return key
		}
	}
	return types.NilHash
}

func (d *Dag) getMax(set types.Set) types.Hash {
	tmpMaxScore := -1.0
	result := types.NilHash

	for _, child := range set.List() {
		if score, exist := d.parentScore[child]; exist {
			if score > tmpMaxScore {
				tmpMaxScore = score
				result = child
			} else if score == tmpMaxScore {
				sStr := result.String()
				blockStr := child.String()
				if sStr < blockStr {
					result = child
				}
			}
		}
	}
	return result
}

func (d *Dag) GetLastPivot(start types.Hash) types.Hash {
	d.graphLock.RLock()
	defer d.graphLock.RUnlock()

	if _, exist := d.graph[start]; start == types.NilHash || !exist {
		return types.NilHash
	}
	children, exist := d.parentRevGraph[start]
	for exist && !children.IsEmpty() {
		child := d.getMax(children)
		if child == types.NilHash {
			return start
		}
		start = child
		children, exist = d.parentRevGraph[start]
	}
	return start
}

func (d *Dag) GetChildren(cur types.Hash) types.Set {
	if children, exist := d.revGraph[cur]; exist {
		return children
	}

	return types.NewSet()
}

func (d *Dag) Contains(key types.Hash) bool {
	_, exist := d.graph[key]
	return exist
}

func (d *Dag) GetScore(key types.Hash) float64 {
	d.graphLock.RLock()
	defer d.graphLock.RUnlock()

	score, exist := d.score[key]
	if exist {
		return score
	} else {
		return 0.0
	}
}
