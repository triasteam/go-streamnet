package tipselection

import "github.com/triasteam/go-streamnet/types"

type Walker interface {
	Walk(entryPoint types.Hash, ratings map[types.Hash]int, walkValidator WalkValidator) types.Hash
}
