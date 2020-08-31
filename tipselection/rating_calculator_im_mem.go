package tipselection

import (
	dag "github.com/triasteam/go-streamnet/dag"
	"github.com/triasteam/go-streamnet/types"
)

type CumulativeWeightMemCalculator struct {
}

func (cal *CumulativeWeightMemCalculator) Init() {
}

func (cu *CumulativeWeightMemCalculator) Calculate(dag *dag.Dag, entryPoint types.Hash) map[types.Hash]int {
	ret := make(map[types.Hash]int)

	visited := types.NewSet()
	queue := types.List{}
	queue.Append(entryPoint)

	for !queue.IsEmpty() {
		cur := queue.RemoveAtIndex(0)
		for _, child := range dag.GetChildren(cur).List() {
			if dag.Contains(child) && !visited.Has(child) {
				queue.Append(child)
				visited.Add(child)
			}
		}
		ret[cur] = int(dag.GetScore(cur))
	}

	return ret
}
