package tipselection

import (
	"github.com/triasteam/go-streamnet/dag"
	"github.com/triasteam/go-streamnet/types"
)

type TipSelectorStreamWork struct {
	dag    *dag.Dag
	ep     EntryPoint
	cal    RatingCalculator
	walker Walker
}

func (tips *TipSelectorStreamWork) Init(dag *dag.Dag) {
	tips.dag = dag

	// todo: using config to choose entrypoint selector.
	ep := EntryPointKatz{}
	ep.Init(dag)
	tips.ep = &ep

	// todo: using config to choose rating calculator.
	cal := CumulativeWeightMemCalculator{}
	cal.Init(dag)
	tips.cal = &cal

	// todo: using config to choose tail finder.
	tf := TailFinderImpl{}
	tf.Init(dag)

	// todo: using config to choose walker.
	walker := WalkerAlpha{}
	walker.Init(dag, &tf)
	tips.walker = &walker
}

func (ts *TipSelectorStreamWork) GetTransactionsToApprove(depth int, reference types.Hash) types.List {
	tips := types.List{}

	// Parental tip
	parentTip := ts.dag.GetLastPivot(ts.dag.GetGenesis())
	tips.Add(parentTip)

	// Reference tip
	entryPoint := ts.ep.GetEntryPoint(depth)

	rating := ts.cal.Calculate(entryPoint)

	var walkValidator WalkerValidatorImpl

	refTip := ts.walker.Walk(entryPoint, rating, &walkValidator)
	tips.Append(refTip)

	// TODO validate UTXO etc.

	return tips
}
