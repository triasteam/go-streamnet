package tipselection

import "github.com/triasteam/StreamNet-go/types"

type TipSelector interface {
	getTransactionsToApprove(depth int, reference *types.Hash)([]types.Hash)
}


