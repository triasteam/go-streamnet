package tipselection

import (
	"github.com/triasteam/go-streamnet/dag"
	"github.com/triasteam/go-streamnet/types"
)

type TailFinderImpl struct {
	dag *dag.Dag
}

func (tf *TailFinderImpl) findTail(hash types.Hash) types.Hash {
	/*tx := types.FromHash(tf.dag, hash)

	approvees := tx.GetApprovers()
	foundApprovee := false
	for _, approvee := range approvees.List() {
		nextTx := types.fromHash(tf.dag, approvee)
		if nextTx.getCurrentIndex() == index && bundleHash.equals(nextTx.getBundleHash()) {
			tx = nextTx
			foundApprovee = true
			break
		}
	}
	if !foundApprovee {
		break
	}

	if tx.getCurrentIndex() == 0 {
		return Optional.of(tx.getHash())
	}*/
	return types.NilHash
}
