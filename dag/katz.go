package dag

import (
	"math"

	"github.com/triasteam/go-streamnet/types"
)

const DEFAULT_ALPHA = 1.0
const DEFAULT_BETA = 1.0
const MAX_ITERATIONS = 100
const EPSILON = 1e-6

type KatzCentrality struct {
	alpha float64
	beta  float64

	network     map[types.Hash]types.Set
	revNetwork  map[types.Hash]types.Set
	allVertices types.Set
	score       map[types.Hash]float64
}

func NewKatz(network, revNetwork map[types.Hash]types.Set, alpha float64) *KatzCentrality {
	kz := KatzCentrality{}

	kz.network = network
	kz.alpha = alpha
	kz.beta = DEFAULT_BETA

	kz.allVertices = types.NewSet()

	for k, _ := range kz.network {
		kz.allVertices.Add(k)
		for _, v := range kz.network[k].List() {
			if _, ok := network[v]; !ok {
				kz.allVertices.Add(v)
			}
		}
	}
	if len(revNetwork) == 0 {
		kz.revNetwork = make(map[types.Hash]types.Set)
		for k, _ := range network {
			for _, v := range network[k].List() {
				if _, ok := kz.revNetwork[v]; !ok {
					kz.revNetwork[v] = types.NewSet()
				}
				kz.revNetwork[v].Add(k)
			}
		}
	} else {
		kz.revNetwork = revNetwork
	}

	return &kz
}

func (kz *KatzCentrality) SetScore(score map[types.Hash]float64) {
	kz.score = score
}

func (kz *KatzCentrality) Compute() map[types.Hash]float64 {
	size := float64(kz.allVertices.Len())
	centrality := make(map[types.Hash]float64)
	old := make(map[types.Hash]float64)

	if len(kz.score) == 0 {
		for _, v := range kz.allVertices.List() {
			centrality[v] = 1.0 / size
			old[v] = 1.0 / size
		}
	} else {
		scale := float64(len(kz.score)-1) / float64(len(kz.score))
		for _, v := range kz.allVertices.List() {
			centrality[v] = kz.score[v] * scale
			old[v] = kz.score[v] * scale
		}
	}
	if size <= 1 {
		return centrality
	}

	// Power iteration: O(k(n+m))
	// The value of norm converges to the dominant eigenvalue, and the vector 'centrality' to an associated eigenvector
	// ref. http://en.wikipedia.org/wiki/Power_iteration
	change := math.MaxFloat64
	for iteration := 0; (iteration < MAX_ITERATIONS) && (change > EPSILON); iteration++ {
		tmp := old
		old = centrality
		centrality = tmp
		sum2 := 0.0
		for _, v := range kz.allVertices.List() {
			centrality[v] = 0.0
			if _, ok := kz.revNetwork[v]; !ok {
				continue
			}
			for _, u := range kz.revNetwork[v].List() {
				centrality[v] = centrality[v] + old[u]
			}
			centrality[v] = kz.alpha*centrality[v] + kz.beta/size
			sum2 += centrality[v] * centrality[v]
		}
		// Normalization
		norm := math.Sqrt(sum2)
		change = 0
		for _, v := range kz.allVertices.List() {
			centrality[v] = centrality[v] / norm
			if math.Abs(centrality[v]-old[v]) > change {
				change = math.Abs(centrality[v] - old[v])
			}
		}
	}
	return centrality
}
