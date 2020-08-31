package tipselection

import (
	"github.com/triasteam/go-streamnet/dag"
	"github.com/triasteam/go-streamnet/types"
)

type EntryPointKatz struct {
}

func (kz *EntryPointKatz) Init() {

}

func (kz *EntryPointKatz) GetEntryPoint(dag *dag.Dag, depth int) types.Hash {
	// todo:
	streamingGraph := true

	if streamingGraph {
		return dag.GetPivotalHash(depth)
	} else {
		dag.BuildGraph()
		dag.ComputeScore()

		return dag.GetPivotalHash(depth)
	}
}
