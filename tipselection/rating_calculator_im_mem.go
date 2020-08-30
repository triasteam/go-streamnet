package tipselection

import (
	dag "github.com/triasteam/go-streamnet/dag"
	"github.com/triasteam/go-streamnet/types"
)

type CumulativeWeightMemCalculator struct {
	dag *dag.Dag
}

func (cal *CumulativeWeightMemCalculator) Init(dag *dag.Dag) {
	cal.dag = dag
}

func (cu *CumulativeWeightMemCalculator) Calculate(entryPoint types.Hash) map[types.Hash]int {
	r := make(map[types.Hash]int)

	visited := types.NewSet()
	queue := types.List{}
	queue.Append(entryPoint)
	var h types.Hash
	for !queue.IsEmpty() {
		h = queue.RemoveAtIndex(0)
		for _, e := range cu.dag.GetChildren(h).List() {
			if cu.dag.Contains(e) && !visited.Has(e) {
				queue.Append(e)
				visited.Add(e)
			}
		}
		r[h] = int(cu.dag.GetScore(h))
	}

	return r
}
