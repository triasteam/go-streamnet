package tipselection

import (
	"container/list"
	"math"
	"math/rand"

	"github.com/triasteam/go-streamnet/dag"
	"github.com/triasteam/go-streamnet/types"
)

const ALPHA = 0.001

type WalkerAlpha struct {
	dag *dag.Dag
}

func (walker *WalkerAlpha) Init(d *dag.Dag) {
	walker.dag = d
}

func (walker *WalkerAlpha) walk(entryPoint types.Hash, ratings map[types.Hash]int, walkValidator WalkValidator) types.Hash {
	if !walkValidator.IsValid(entryPoint) {
		return types.NewHash(nil)
	}

	var nextStep types.Hash
	traversedTails := types.List{}
	traversedTails.Append(entryPoint)

	//Walk
	ok := true
	for ok {
		nextStep = walker.selectApprover(traversedTails.GetLast(), ratings, walkValidator)
		if nextStep != types.NewHash(nil) {
			traversedTails.Append(nextStep)
		}
		ok = nextStep == types.NewHash(nil)
	}

	return traversedTails.GetLast()
}

func (walker *WalkerAlpha) selectApprover(tailHash types.Hash, ratings map[types.Hash]int, walkValidator WalkValidator) types.Hash {
	approvers := types.NewSet()

	approvers1 := walker.dag.GetChild(tailHash)
	approvers.AddAll(approvers1)

	return findNextValidTail(ratings, approvers, walkValidator)
}

func findNextValidTail(ratings map[types.Hash]int, approvers types.Set, walkValidator WalkValidator) types.Hash {
	nextTailHash := types.NilHash

	//select next tail to step to
	for nextTailHash == types.NewHash(nil) {
		nextTxHash := selects(ratings, approvers)
		if nextTxHash == types.NewHash(nil) {
			//no existing approver = tip
			return types.NewHash(nil)
		}

		nextTailHash = findTailIfValid(nextTxHash, walkValidator)
		approvers.Remove(nextTxHash)
		//if next tail is not valid, re-select while removing it from approvers set
	}

	return nextTailHash
}

func selects(ratings map[types.Hash]int, approversSet types.Set) types.Hash {
	//filter based on tangle state when starting the walk
	approvers := types.List{}
	for _, hash := range approversSet.List() {
		if _, ok := ratings[hash]; ok {
			approvers.Append(hash)
		}
	}

	//After filtering, if no approvers are available, it's a tip.
	if approvers.Length() == 0 {
		return types.NewHash(nil)
	}

	//calculate the probabilities
	walkRatings := list.New()
	for i := 0; i < approvers.Length(); i++ {
		hash := approvers.Index(i)
		v := ratings[hash]
		walkRatings.PushBack(v)
	}

	var maxRating = 0
	for e := walkRatings.Front(); e != nil; e = e.Next() {
		if maxRating < e.Value.(int) {
			maxRating = e.Value.(int)
		}
	}

	//transition probability function (normalize ratings based on Hmax)
	var normalizedWalkRatings = list.New()
	for e := walkRatings.Front(); e != nil; e = e.Next() {
		normalizedWalkRatings.PushBack(e.Value.(int) - maxRating)
	}

	var weights = list.New()
	for e := normalizedWalkRatings.Front(); e != nil; e = e.Next() {
		weights.PushBack(math.Exp(ALPHA * e.Value.(float64)))
	}

	//select the next transaction
	var weightsSum = 0.0
	for e := weights.Front(); e != nil; e = e.Next() {
		weightsSum += e.Value.(float64)
	}

	target := rand.Float64() * weightsSum

	var approverIndex int
	for e := weights.Front(); e != nil; e = e.Next() {
		target = target - e.Value.(float64)
		if target <= 0 {
			break
		}
		approverIndex++
	}

	return approvers.Index(approverIndex)
}

func findTailIfValid(transactionHash types.Hash, validator WalkValidator) types.Hash {
	/*tailHash := tailFinder.findTail(transactionHash)
	if tailHash.isPresent() && validator.isValid(tailHash.get()) {
		return tailHash
	}
	return Optional.empty()*/
	return types.NilHash
}