package tipselection

import "github.com/triasteam/go-streamnet/types"

type WalkValidator interface {
	IsValid(hash types.Hash) bool
}
