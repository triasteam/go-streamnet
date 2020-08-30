package tipselection

import (
	"github.com/triasteam/go-streamnet/dag"
	"github.com/triasteam/go-streamnet/types"
)

type Calculator struct {
	dag *dag.Dag
}

func (c *Calculator) Calculate(entryPoint types.Hash) map[types.Hash]int {

	ret := make(map[types.Hash]int)

	visited := types.NewSet()
	queue := types.List{}
	queue.Append(entryPoint)

	for !queue.IsEmpty() {
		h := queue.RemoveAtIndex(0)
		for _, e := range c.dag.GetChildren(h).List() {
			if c.dag.Contains(e) && !visited.Has(e) {
				queue.Append(e)
				visited.Add(e)
			}
		}
		ret[h] = int(c.dag.GetScore(h))
	}

	return ret
}
