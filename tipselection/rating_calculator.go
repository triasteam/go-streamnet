package tipselection

import "github.com/triasteam/go-streamnet/types"

type RatingCalculator interface {
	Calculate(entryPoint types.Hash) map[types.Hash]int
}
