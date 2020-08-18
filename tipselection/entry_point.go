package tipselection

import (
	"github.com/triasteam/go-streamnet/dag"
	"github.com/triasteam/go-streamnet/types"
)

type EntryPoint interface {
	GetEntryPoint(d *dag.Dag, depth int) types.Hash
}
