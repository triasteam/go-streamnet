package tipselection

import (
	"github.com/triasteam/go-streamnet/dag"
	dag2 "github.com/triasteam/go-streamnet/dag"
	"github.com/triasteam/go-streamnet/types"
)

// Milestone methods is the way IOTA seleted. StreamNet won't use it.
type EntryPointMilestone struct {
	dag *dag2.Dag
}

func (ms *EntryPointMilestone) GetEntryPoint(dag *dag.Dag, depth int) types.Hash {
	return types.NilHash
}
