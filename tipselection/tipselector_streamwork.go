package tipselection

import (
	"github.com/triasteam/go-streamnet/dag"
	"github.com/triasteam/go-streamnet/types"
)

type TipSelectorStreamWork struct {
	dag *dag.Dag
	ep  EntryPoint
	cal RatingCalculator
}

func (tips *TipSelectorStreamWork) Init(dag *dag.Dag) {
	tips.dag = dag

}

func (ts *TipSelectorStreamWork) GetTransactionsToApprove(depth int, reference types.Hash) types.List {
	tips := types.List{}

	// Parental tip
	parentTip := ts.dag.GetLastPivot(ts.dag.GetGenesis())
	tips.Add(parentTip)

	// Reference tip
	entryPoint := ts.ep.GetEntryPoint(ts.dag, depth)

	rating := ts.cal.Calculate(entryPoint)

	var walkValidator WalkerValidatorImpl

	var walker WalkerAlpha
	walker.Init(ts.dag)

	refTip := walker.walk(entryPoint, rating, walkValidator)
	tips.Append(refTip)

	// TODO validate UTXO etc.

	return tips
}
