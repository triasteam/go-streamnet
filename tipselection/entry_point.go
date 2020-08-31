package tipselection

import (
	"github.com/triasteam/go-streamnet/dag"
	"github.com/triasteam/go-streamnet/types"
)

type EntryPoint interface {
	GetEntryPoint(dag *dag.Dag, depth int) types.Hash
}
