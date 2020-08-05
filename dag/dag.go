package dag

import (
	"github.com/triasteam/go-streamnet/store"
	"github.com/triasteam/go-streamnet/types"
)

// Dag is the most important struct in the whole procedure.
type Dag struct {
	graph       map[types.Hash][]types.Hash
	parentGraph map[types.Hash]types.Hash

	revGraph       map[types.Hash][]types.Hash
	parentRevGraph map[types.Hash][]types.Hash

	score       map[types.Hash]float64
	parentScore map[types.Hash]float64

	degrees           map[types.Hash]int64
	topOrder          map[int64][]types.Hash
	topOrderStreaming map[int64][]types.Hash

	subGraph          map[types.Hash][]types.Hash
	subRevGraph       map[types.Hash][]types.Hash
	subParentGraph    map[types.Hash]types.Hash
	subParentRevGraph map[types.Hash][]types.Hash

	levelmap map[types.Hash]int64
	namemap  map[types.Hash]string

	totalDepth int
	store      store.Storage
	// to use
	pivotChain []types.Hash
}

// Init return a new Dag struct.
func Init() *Dag {

}

// Close will free all the resources.
func (dag *Dag) Close() {

}
