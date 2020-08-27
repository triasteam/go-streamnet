package tipselection

import "github.com/triasteam/go-streamnet/types"

type WalkerValidatorImpl struct {
}

func (w WalkerValidatorImpl) IsValid(hash types.Hash) bool {

	/*if (belowMaxDepth(hash(),
			milestoneTracker.latestSolidSubtangleMilestoneIndex - config.getMaxDepth())) {
		log.debug("Validation failed: {} is below max depth", hash);
		return false;
	}*/
	return true
}

/*private boolean belowMaxDepth(Hash tip, int lowerAllowedSnapshotIndex) throws Exception {
	//if tip is confirmed stop
	if (TransactionViewModel.fromHash(tangle, tip).snapshotIndex() >= lowerAllowedSnapshotIndex) {
		return false;
	}
	//if tip unconfirmed, check if any referenced tx is confirmed below maxDepth
	Queue<Hash> nonAnalyzedTransactions = new LinkedList<>(Collections.singleton(tip));
	Set<Hash> analyzedTransactions = new HashSet<>();
	Hash hash;
	final int maxAnalyzedTransactions = config.getBelowMaxDepthTransactionLimit();
	while ((hash = nonAnalyzedTransactions.poll()) != null) {
		if (analyzedTransactions.size() == maxAnalyzedTransactions) {
			log.debug("failed below max depth because of exceeding max threshold of {} analyzed transactions",
					maxAnalyzedTransactions);
			return true;
		}

		if (analyzedTransactions.add(hash)) {
			TransactionViewModel transaction = TransactionViewModel.fromHash(tangle, hash);
			if ((transaction.snapshotIndex() != 0 || Objects.equals(Hash.NULL_HASH, transaction.getHash()))
					&& transaction.snapshotIndex() < lowerAllowedSnapshotIndex) {
				log.debug("failed below max depth because of reaching a tx below the allowed snapshot index {}",
						lowerAllowedSnapshotIndex);
				return true;
			}
			if (transaction.snapshotIndex() == 0) {
				if (!maxDepthOkMemoization.contains(hash)) {
					nonAnalyzedTransactions.offer(transaction.getTrunkTransactionHash());
					nonAnalyzedTransactions.offer(transaction.getBranchTransactionHash());
				}
			}
		}
	}
	maxDepthOkMemoization.add(tip);
	return false;
}
}*/
