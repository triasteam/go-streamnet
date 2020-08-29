package streamnet

import (
	"github.com/triasteam/go-streamnet/dag"
	"github.com/triasteam/go-streamnet/store"
	"github.com/triasteam/go-streamnet/tipselection"
)

// StreamNet is the biggest structure.
type StreamNet struct {
	Dag   *dag.Dag
	Store store.StorageProvider
	Tips  tipselection.TipSelector
}
