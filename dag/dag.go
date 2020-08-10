package dag

import (
	"github.com/triasteam/go-streamnet/store"
	"github.com/triasteam/go-streamnet/types"
)

// Dag is the most important struct in the whole procedure.
type Dag struct {
	graph       map[types.Hash]types.Set
	parentGraph map[types.Hash]types.Hash

	revGraph       map[types.Hash]types.Set
	parentRevGraph map[types.Hash]types.Set

	score       map[types.Hash]float64
	parentScore map[types.Hash]float64

	degrees           map[types.Hash]int64
	topOrder          map[int64]types.Set
	topOrderStreaming map[int64]types.Set

	subGraph          map[types.Hash]types.Set
	subRevGraph       map[types.Hash]types.Set
	subParentGraph    map[types.Hash]types.Hash
	subParentRevGraph map[types.Hash]types.Set

	levelmap map[types.Hash]int64
	namemap  map[types.Hash]string

	totalDepth int
	store      *store.Storage

	// to use
	pivotChain types.Set
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
		degrees:           make(map[types.Hash]int64),
		topOrder:          make(map[int64]types.Set),
		topOrderStreaming: make(map[int64]types.Set),
		subGraph:          make(map[types.Hash]types.Set),
		subRevGraph:       make(map[types.Hash]types.Set),
		subParentGraph:    make(map[types.Hash]types.Hash),
		subParentRevGraph: make(map[types.Hash]types.Set),
		levelmap:          make(map[types.Hash]int64),
		namemap:           make(map[types.Hash]string),
		totalDepth:        0,
		store:             db,
		pivotChain:        types.NewSet(),
	}
}

// Close will free all the resources.
func (dag *Dag) Close() {

}
