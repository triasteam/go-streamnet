package tipselection

import (
	"github.com/triasteam/go-streamnet/dag"
	"github.com/triasteam/go-streamnet/types"
)

type CumulativeWeightWithEdgeCalculator struct {
	dag *dag.Dag
}

func (cu CumulativeWeightWithEdgeCalculator) calculate(entryPoint types.Hash) map[types.Hash]int {
	r := make(map[types.Hash]int)
	return r
}
