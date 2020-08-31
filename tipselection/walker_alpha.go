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
	dag        *dag.Dag
	tailFinder TailFinder
}

func (walker *WalkerAlpha) Init(d *dag.Dag, tf TailFinder) {
	walker.dag = d
	walker.tailFinder = tf
}

func (walker *WalkerAlpha) Walk(entryPoint types.Hash, ratings map[types.Hash]int, walkValidator WalkValidator) types.Hash {
	if !walkValidator.IsValid(entryPoint) {
		return types.NilHash
	}

	var nextStep types.Hash
	traversedTails := types.List{}
	traversedTails.Append(entryPoint)

	//Walk
	for {
		nextStep = walker.findNextValidTail(traversedTails.GetLast(), ratings, walkValidator)
		if nextStep == types.NilHash {
			break
		}
		traversedTails.Append(nextStep)
	}

	return traversedTails.GetLast()
}

func (walker *WalkerAlpha) findNextValidTail(tailHash types.Hash, ratings map[types.Hash]int, walkValidator WalkValidator) types.Hash {
	nextTailHash := types.NilHash

	approvers := walker.dag.GetChildren(tailHash)

	//select next tail to step to
	for nextTailHash == types.NilHash {
		nextTxHash := selects(ratings, approvers)
		if nextTxHash == types.NilHash {
			//no existing approver = tip
			return types.NilHash
		}

		nextTailHash = walker.findTailIfValid(nextTxHash, walkValidator)
		approvers.Remove(nextTxHash)
		//if next tail is not valid, re-select while removing it from approvers set
	}

	return nextTailHash
}

func selects(ratings map[types.Hash]int, approversSet types.Set) types.Hash {
	//filter based on tangle state when starting the walk
	approvers := types.List{}
	for _, hash := range approversSet.List() {
		if _, exist := ratings[hash]; exist {
			approvers.Append(hash)
		}
	}

	//After filtering, if no approvers are available, it's a tip.
	if approvers.Length() == 0 {
		return types.NilHash
	}

	//calculate the probabilities
	walkRatings := list.New()
	maxRating := 0
	for i := 0; i < approvers.Length(); i++ {
		hash := approvers.Index(i)
		v := ratings[hash]
		walkRatings.PushBack(v)

		if maxRating < v {
			maxRating = v
		}
	}

	//transition probability function (normalize ratings based on Hmax)
	var normalizedWalkRatings = list.New()
	for e := walkRatings.Front(); e != nil; e = e.Next() {
		normalizedWalkRatings.PushBack(e.Value.(int) - maxRating)
	}

	var weights = list.New()
	for e := normalizedWalkRatings.Front(); e != nil; e = e.Next() {
		weights.PushBack(math.Exp(ALPHA * float64(e.Value.(int))))
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

func (walker *WalkerAlpha) findTailIfValid(transactionHash types.Hash, validator WalkValidator) types.Hash {
	tailHash := walker.tailFinder.findTail(transactionHash)
	if tailHash != types.NilHash && validator.IsValid(tailHash) {
		return tailHash
	}
	return types.NilHash
}
