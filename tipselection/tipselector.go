package tipselection

import (
	"github.com/triasteam/go-streamnet/types"
)

type TipSelector interface {
	GetTransactionsToApprove(depth int, reference types.Hash) types.List
}
