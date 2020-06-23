package tipselection

import "github.com/triasteam/go-streamnet/types"

type TipSelector interface {
	getTransactionsToApprove(depth int, reference *types.Hash)([]types.Hash)
}


