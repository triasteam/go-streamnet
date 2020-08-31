package tipselection

import (
	"github.com/triasteam/go-streamnet/dag"
	"github.com/triasteam/go-streamnet/types"
)

type CumulativeWeightCalculator struct {
	dag *dag.Dag
}

func (cu CumulativeWeightCalculator) calculate(dag *dag.Dag, entryPoint types.Hash) map[types.Hash]int {
	r := make(map[types.Hash]int)
	return r
}
