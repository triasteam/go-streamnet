package tipselection

import (
	dag2 "github.com/triasteam/go-streamnet/dag"
	"github.com/triasteam/go-streamnet/types"
)

// Milestone methods is the way IOTA seleted. StreamNet won't use it.
type EntryPointMilestone struct {
	dag *dag2.Dag
}

func (ms *EntryPointMilestone) GetEntryPoint(depth int) types.Hash {
	return types.NilHash
}
