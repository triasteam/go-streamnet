package tipselection

import (
	"github.com/triasteam/go-streamnet/dag"
	"github.com/triasteam/go-streamnet/types"
)

type RatingCalculator interface {
	Calculate(dag *dag.Dag, entryPoint types.Hash) map[types.Hash]int
}
