package tipselection

import (
	"github.com/triasteam/go-streamnet/dag"
	"github.com/triasteam/go-streamnet/types"
)

type EntryPointKatz struct {
	dag *dag.Dag
}

func (kz *EntryPointKatz) Init(dag *dag.Dag) {
	kz.dag = dag
}

func (kz *EntryPointKatz) GetEntryPoint(depth int) types.Hash {
	streamingGraph := true

	if streamingGraph {
		return kz.dag.GetPivotalHash(depth)
	} else {
		kz.dag.BuildGraph()
		kz.dag.ComputeScore()

		return kz.dag.GetPivotalHash(depth)
	}
}
