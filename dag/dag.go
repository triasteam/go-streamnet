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

	pivotChain types.Set

	graphLock sync.RWMutex
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
		pivotChain:        types.NewSet(),
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
	scoreAlg := "KATZ"

	if scoreAlg == "CUM_WEIGHT" {
		Update(d.graph, d.score, key, 1)
		UpdateParentScore(d.parentGraph, d.parentScore, key, 1)

	} else if scoreAlg == "KATZ" {
		d.score[key] = 1.0 / (float64(len(d.score)) + 1.0)
		centrality := NewKatz(d.graph, d.revGraph, 0.5)
		centrality.SetScore(d.score)
		d.score = centrality.Compute()
		d.parentScore = UpdateParentScore(d.parentGraph, d.parentScore, key, 1.0)
	}
	d.freshScore = false
}
