package tipselection

import (
	"github.com/triasteam/go-streamnet/dag"
	"github.com/triasteam/go-streamnet/types"
)

func GetEntryPoint(d *dag.Dag, depth int) types.Hash {
	streamingGraph := true

	if streamingGraph {
		return d.GetPivotalHash(depth)
	} else {
		d.BuildGraph()
		d.ComputeScore()

		d.GetPivotalHash(depth)
	}
}
