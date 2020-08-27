package tipselection

import "github.com/triasteam/go-streamnet/types"

type TailFinder interface {
	findTail(hash types.Hash) types.Hash
}
