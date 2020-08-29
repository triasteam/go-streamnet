package tipselection

import (
	"github.com/triasteam/go-streamnet/types"
)

type EntryPoint interface {
	GetEntryPoint(depth int) types.Hash
}
