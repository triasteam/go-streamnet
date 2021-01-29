// Copyright 2017 The GoReporter Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package dag provides graph storage.
package dag

import (
	"sync"

	"github.com/triasteam/go-streamnet/store"
	"github.com/triasteam/go-streamnet/types"
	"github.com/triasteam/go-streamnet/utils"
)

// Dag is the most important struct in the whole procedure.
// memory storage of graph, output a total order of nodes with conflux algorithm
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
func (dag *Dag) Close() {

}

// Add puts one tx hash into dag.
func (dag *Dag) Add(key types.Hash, tx *types.Transaction) error {
	trunk := tx.GetTrunkTransactionHash()
	branch := tx.GetBranchTransactionHash()

	dag.graphLock.Lock()
	defer dag.graphLock.Unlock()

	dag.updateGraph(key, trunk, branch)

	dag.updateTopologicalOrder(key, trunk, branch)

	dag.updateScore(key)

	return nil
}

func (dag *Dag) updateGraph(key, trunk, branch types.Hash) {
	// Approve graph
	if _, exist := dag.graph[key]; !exist {
		dag.graph[key] = types.NewSet()
	}
	dag.graph[key].Add(trunk)
	dag.graph[key].Add(branch)

	//parentGraph
	dag.parentGraph[key] = trunk

	// Approvee graph
	if _, exist := dag.revGraph[trunk]; !exist {
		dag.revGraph[trunk] = types.NewSet()
	}
	dag.revGraph[trunk].Add(key)
	if _, exist := dag.revGraph[branch]; !exist {
		dag.revGraph[branch] = types.NewSet()
	}
	dag.revGraph[branch].Add(key)

	if _, exist := dag.parentRevGraph[trunk]; !exist {
		dag.parentRevGraph[trunk] = types.NewSet()
	}
	dag.parentRevGraph[trunk].Add(key)

	// update degrees
	if _, exist := dag.degrees[key]; !exist || dag.degrees[key] == 0 {
		dag.degrees[key] = 2
	}
	if _, exist := dag.degrees[trunk]; !exist {
		dag.degrees[trunk] = 0
	}
	if _, exist := dag.degrees[branch]; !exist {
		dag.degrees[branch] = 0
	}
}

func (dag *Dag) updateTopologicalOrder(key, trunk, branch types.Hash) {
	if len(dag.topOrderStreaming) == 0 {
		dag.topOrderStreaming[1] = types.NewSet()
		dag.topOrderStreaming[1].Add(key)
		dag.levels[key] = 1
		dag.topOrderStreaming[0] = types.NewSet()
		dag.topOrderStreaming[0].Add(trunk)
		dag.topOrderStreaming[0].Add(branch)
		dag.totalDepth = 1
		return
	} else {
		// TODO: check trunk or branch exist !!!!!!!!!
		// Or we won't call Add if trunk or branch not exist!!!!
		trunkLevel := dag.levels[trunk]
		branchLevel := dag.levels[branch]
		level := utils.Min(trunkLevel, branchLevel) + 1
		if _, exist := dag.topOrderStreaming[level]; !exist {
			dag.topOrderStreaming[level] = types.NewSet()
			dag.totalDepth++
		}
		dag.topOrderStreaming[level].Add(key)
		dag.levels[key] = level
	}
}

func (dag *Dag) updateScore(key types.Hash) {
	// todo: use config to choose score algorithm.
	scoreAlg := "CUM_WEIGHT"

	if scoreAlg == "CUM_WEIGHT" {
		CumulateWeight{}.UpdateScore(dag.graph, dag.score, key, 1)
		CumulateWeight{}.UpdateTrunkScore(dag.parentGraph, dag.parentScore, key, 1)

	} else if scoreAlg == "KATZ" {
		dag.score[key] = 1.0 / (float64(len(dag.score)) + 1.0)
		centrality := NewKatz(dag.graph, dag.revGraph, dag.score, 0.5)
		dag.score = centrality.Compute()
		CumulateWeight{}.UpdateTrunkScore(dag.parentGraph, dag.parentScore, key, 1)
	}
	dag.freshScore = false
}

// GetPivotalHash returns the pivot hash of dag from genesis. depth is 0, 1, 2...
func (dag *Dag) GetPivotalHash(depth int) types.Hash {
	var ret types.Hash
	dag.BuildPivotChain()
	if depth == -1 || depth >= dag.pivotChain.Length() {
		set := dag.topOrderStreaming[1]
		if set.IsEmpty() {
			return types.NilHash
		}
		ret = set.List()[0]
		return ret
	}

	ret = dag.pivotChain.Index(dag.pivotChain.Length() - depth - 1)
	return ret
}

func (dag *Dag) BuildPivotChain() {
	dag.graphLock.RLock()
	defer dag.graphLock.RUnlock()
	dag.pivotChain = dag.getPivotChainFrom(dag.GetGenesis())
}

// get the pivot chain after the giving node.
func (dag *Dag) getPivotChainFrom(start types.Hash) types.List {
	/*dag.graphLock.RLock()
	defer dag.graphLock.RUnlock()*/

	list := types.List{}

	if _, exist := dag.graph[start]; start == types.NilHash || !exist {
		return list
	}

	list.Append(start)

	children, ok := dag.parentRevGraph[start]
	for ok && !children.IsEmpty() {
		child := dag.getMax(&children)
		if child == types.NilHash {
			return list
		}

		list.Append(child)

		start = child
		children, ok = dag.parentRevGraph[start]
	}

	return list
}

// GetGenesis todo: use config to get genesis directly; if genesis_forward is implemented, it should be as following.
func (dag *Dag) GetGenesis() types.Hash {
	if dag.ancestors != nil && !dag.ancestors.Empty() {
		return dag.ancestors.Peek()
	}

	// find one node whose trunk parent is not in parentGraph
	for key, trunk := range dag.parentGraph {
		if _, ok := dag.parentGraph[trunk]; !ok {
			return key
		}
	}
	return types.NilHash
}

func (dag *Dag) getMax(set *types.Set) types.Hash {
	tmpMaxScore := -1.0
	result := types.NilHash

	for _, child := range set.List() {
		if score, exist := dag.parentScore[child]; exist {
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

// GetLastPivot ...
func (dag *Dag) GetLastPivot(start types.Hash) types.Hash {
	dag.graphLock.RLock()
	defer dag.graphLock.RUnlock()

	if _, exist := dag.graph[start]; start == types.NilHash || !exist {
		return types.NilHash
	}
	children, exist := dag.parentRevGraph[start]
	for exist && !children.IsEmpty() {
		child := dag.getMax(&children)
		if child == types.NilHash {
			return start
		}
		start = child
		children, exist = dag.parentRevGraph[start]
	}
	return start
}

// GetChildren returns nodes whoes parent is input param `cur`,
// if `cur` has no nodes referenced empty set will return
func (dag *Dag) GetChildren(cur types.Hash) *types.Set {
	if children, exist := dag.revGraph[cur]; exist {
		return &children
	}

	s := types.NewSet()
	return &s
}

// Contains ...
func (dag *Dag) Contains(key types.Hash) bool {
	_, exist := dag.graph[key]
	return exist
}

// GetScore not calculate, get from score map
func (dag *Dag) GetScore(key types.Hash) float64 {
	dag.graphLock.RLock()
	defer dag.graphLock.RUnlock()

	score, exist := dag.score[key]
	if exist {
		return score
	}
	return 0.0
}
