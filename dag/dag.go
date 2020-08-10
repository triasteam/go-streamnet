package dag

const HASH_LEN = 128

type Hash [HASH_LEN]byte

type DagService interface {
	add(Hash) error
	totalOrder() []Hash
	getPeriod(int) []Hash
}

type Dag struct {
	graph       map[Hash][]Hash
	parentGraph map[Hash]Hash

	revGraph       map[Hash][]Hash
	parentRevGraph map[Hash][]Hash

	score       map[Hash]float64
	parentScore map[Hash]float64

	degrees           map[Hash]int64
	topOrder          map[int64][]Hash
	topOrderStreaming map[int64][]Hash

	subGraph          map[Hash][]Hash
	subRevGraph       map[Hash][]Hash
	subParentGraph    map[Hash]Hash
	subParentRevGraph map[Hash][]Hash

	levelmap map[Hash]int64
	namemap  map[Hash]string

	totalDepth int
	store      Store
	// to use
	pivotChain []Hash
}
