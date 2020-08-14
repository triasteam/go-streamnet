package tipselection

import (
	"github.com/triasteam/go-streamnet/dag"
	"github.com/triasteam/go-streamnet/types"
)

type TipSelector interface {
	getTransactionsToApprove(depth int, reference *types.Hash)([]types.Hash)
}


func getTransactionsToApprove(d *dag.Dag, depth int, reference Optional<types.Hash> ) types.List {
	tips := types.List{}

	// Parental tip
	parentTip := d.GetLastPivot(d.GetGenesis())
	tips.Add(parentTip)

	// Reference tip
	entryPoint := GetEntryPoint(d, depth)

	UnIterableMap<HashId, Integer> rating = ratingCalculator.calculate(entryPoint);

	WalkValidator walkValidator = new WalkValidatorImpl(tangle, ledgerValidator, milestoneTracker, config);
	if(BaseIotaConfig.getInstance().getWalkValidator().equals("NULL")) {
		walkValidator = new WalkValidatorNull();
	}

	Hash refTip;
	refTip = walker.walk(entryPoint, rating, walkValidator);
	tips.add(refTip);

	// TODO validate UTXO etc.

	return tips;
}