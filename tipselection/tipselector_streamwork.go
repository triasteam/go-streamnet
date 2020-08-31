package tipselection

import (
	"github.com/triasteam/go-streamnet/config"
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
	ep.Init()
	tips.ep = &ep

	// todo: using config to choose rating calculator.
	cal := CumulativeWeightMemCalculator{}
	cal.Init()
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
	trunkTip := ts.dag.GetLastPivot(ts.dag.GetGenesis())
	tips.Add(trunkTip)

	// Reference tip
	entryPoint := ts.ep.GetEntryPoint(ts.dag, depth)

	rating := ts.cal.Calculate(ts.dag, entryPoint)

	var walkValidator WalkerValidatorImpl

	branchTip := ts.walker.Walk(entryPoint, rating, &walkValidator)
	tips.Append(branchTip)

	// set genesis if trunk and branch are nil.
	if tips.Index(0) == types.NilHash || tips.Index(1) == types.NilHash {
		// Using genesis.
		tips = types.List{}
		tips.Append(config.GenesisTrunk)
		tips.Append(config.GenesisBranch)
	}

	return tips
}
