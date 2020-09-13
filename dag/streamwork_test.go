package dag

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/triasteam/go-streamnet/store"
	"github.com/triasteam/go-streamnet/types"
)

// no mock
func TestGetPivot(t *testing.T) {
	s := store.Storage{}
	// s.Init("/tmp/gorocksdb_test")
	d := Dag{}
	d.Init(&s)

	genesis := types.NewHashString("genesis")

	A := types.NewHashString("A")
	B := types.NewHashString("B")
	C := types.NewHashString("C")
	D := types.NewHashString("D")
	E := types.NewHashString("E")

	valueGenesis := types.Transaction{}
	refGenesis := types.List{}
	refGenesis.Add(types.NewHashString("X"))
	refGenesis.Add(types.NewHashString("Y"))
	valueGenesis.Init(refGenesis, types.NewHashString("data of genesis"))

	valueA := types.Transaction{}
	refA := &types.List{}
	refA.Add(genesis)
	refA.Add(genesis)
	valueA.Init(*refA, types.NewHashString("dataA"))

	valueB := types.Transaction{}
	refB := types.List{}
	refB.Add(A)
	refB.Add(A)
	valueB.Init(refB, types.NewHashString("dataB"))

	valueC := types.Transaction{}
	refC := types.List{}
	refC.Add(B)
	refC.Add(B)
	valueC.Init(refC, types.NewHashString("dataC"))

	valueD := types.Transaction{}
	refD := types.List{}
	refD.Add(C)
	refD.Add(C)
	valueD.Init(refD, types.NewHashString("dataD"))

	valueE := types.Transaction{}
	refE := types.List{}
	refE.Add(D)
	refE.Add(D)
	valueE.Init(refE, types.NewHashString("dataE"))

	d.Add(genesis, &valueGenesis)

	// FIXME block, transaction has no trunk and branck
	d.Add(A, &valueA)
	d.Add(B, &valueB)
	d.Add(C, &valueC)
	d.Add(D, &valueD)
	d.Add(E, &valueE)

	pivot := d.getPivot(E)
	t.Log(pivot)
	a := assert.New(t)

	a.Equal(pivot, E, "pivot is not E")
}

// test totalOrder, G <A,B<C,D<E
func TestGetTotalOrderWithoutTransaction(t *testing.T) {
	// mock graph
	s := store.Storage{}
	s.Init("/tmp/gorocksdb_test")
	d := Dag{}
	d.Init(&s)

	X := types.NewHashString("X")
	Y := types.NewHashString("Y")
	genesis := types.NewHashString("genesis")
	A := types.NewHashString("A")
	B := types.NewHashString("B")
	C := types.NewHashString("C")
	D := types.NewHashString("D")
	E := types.NewHashString("E")

	alias := map[types.Hash]string{}
	alias[genesis] = "genesis"
	alias[A] = "A"
	alias[B] = "B"
	alias[C] = "C"
	alias[D] = "D"
	alias[E] = "E"

	// set genesis
	d.genesis = genesis
	// fill graph
	d.graph[genesis] = types.NewSet()
	d.graph[genesis].Add(X)
	d.graph[genesis].Add(Y)

	d.graph[A] = types.NewSet()
	d.graph[A].Add(genesis)
	d.graph[A].Add(genesis)

	d.graph[B] = types.NewSet()
	d.graph[B].Add(genesis)
	d.graph[B].Add(genesis)

	d.graph[C] = types.NewSet()
	d.graph[C].Add(A)
	d.graph[C].Add(B)

	d.graph[D] = types.NewSet()
	d.graph[D].Add(B)
	d.graph[D].Add(B)

	d.graph[E] = types.NewSet()
	d.graph[E].Add(C)
	d.graph[E].Add(D)

	// file revGraph
	d.revGraph[X] = types.NewSet()
	d.revGraph[X].Add(genesis)

	d.revGraph[Y] = types.NewSet()
	d.revGraph[Y].Add(genesis)

	d.revGraph[genesis] = types.NewSet()
	d.revGraph[genesis].Add(A)
	d.revGraph[genesis].Add(B)

	d.revGraph[A] = types.NewSet()
	d.revGraph[A].Add(C)

	d.revGraph[B] = types.NewSet()
	d.revGraph[B].Add(C)
	d.revGraph[B].Add(D)

	d.revGraph[C] = types.NewSet()
	d.revGraph[C].Add(E)

	d.revGraph[D] = types.NewSet()
	d.revGraph[D].Add(E)
	// fill parentGraph
	d.parentGraph[genesis] = X
	d.parentGraph[A] = genesis
	d.parentGraph[B] = genesis
	d.parentGraph[C] = A
	d.parentGraph[D] = B
	d.parentGraph[E] = C
	// fill parentRevGraph
	d.parentRevGraph[X] = types.NewSet()
	d.parentRevGraph[X].Add(genesis)

	d.parentRevGraph[genesis] = types.NewSet()
	d.parentRevGraph[genesis].Add(A)
	d.parentRevGraph[genesis].Add(B)

	d.parentRevGraph[A] = types.NewSet()
	d.parentRevGraph[A].Add(C)

	d.parentRevGraph[B] = types.NewSet()
	d.parentRevGraph[B].Add(D)

	d.parentRevGraph[C] = types.NewSet()
	d.parentRevGraph[C].Add(E)

	// d.parentRevGraph[D] = types.NewSet()
	// d.parentRevGraph[D].Add(E)
	// file score
	d.score[genesis] = 6
	d.score[A] = 5
	d.score[B] = 4
	d.score[C] = 3
	d.score[D] = 2
	d.score[E] = 1

	d.parentScore[genesis] = 6.0
	d.parentScore[A] = 5
	d.parentScore[B] = 4
	d.parentScore[C] = 3
	d.parentScore[D] = 2
	d.parentScore[E] = 1

	to := d.GetTotalOrder()
	t.Log("total order is : ", to)
	for idx, v := range to {
		t.Log(idx)
		t.Log(alias[v])
	}
}

// Test getTotalOrder with transaction
func TestGetTotalOrder(t *testing.T) {
	s := store.Storage{}
	// s.Init("/tmp/gorocksdb_test")
	d := Dag{}
	d.Init(&s)

	genesis := types.NewHashString("genesis")

	A := types.NewHashString("A")
	B := types.NewHashString("B")
	C := types.NewHashString("C")
	D := types.NewHashString("D")
	E := types.NewHashString("E")

	valueGenesis := types.Transaction{}
	refGenesis := types.List{}
	refGenesis.Add(types.NewHashString("X"))
	refGenesis.Add(types.NewHashString("Y"))
	valueGenesis.Init(refGenesis, types.NewHashString("data of genesis"))

	valueA := types.Transaction{}
	refA := &types.List{}
	refA.Add(genesis)
	refA.Add(genesis)
	valueA.Init(*refA, types.NewHashString("dataA"))

	valueB := types.Transaction{}
	refB := types.List{}
	refB.Add(A)
	refB.Add(A)
	valueB.Init(refB, types.NewHashString("dataB"))

	valueC := types.Transaction{}
	refC := types.List{}
	refC.Add(B)
	refC.Add(B)
	valueC.Init(refC, types.NewHashString("dataC"))

	valueD := types.Transaction{}
	refD := types.List{}
	refD.Add(C)
	refD.Add(C)
	valueD.Init(refD, types.NewHashString("dataD"))

	valueE := types.Transaction{}
	refE := types.List{}
	refE.Add(D)
	refE.Add(D)
	valueE.Init(refE, types.NewHashString("dataE"))

	d.Add(genesis, &valueGenesis)

	// FIXME block, transaction has no trunk and branck
	d.Add(A, &valueA)
	d.Add(B, &valueB)
	d.Add(C, &valueC)
	d.Add(D, &valueD)
	d.Add(E, &valueE)

	totalOrder := d.GetTotalOrder()
	t.Log(totalOrder)

	alias := map[types.Hash]string{}
	alias[genesis] = "genesis"
	alias[A] = "A"
	alias[B] = "B"
	alias[C] = "C"
	alias[D] = "D"
	alias[E] = "E"

	assertList := make([]string, 0)
	for _, h := range totalOrder {
		assertList = append(assertList, alias[h])
	}

	a := assert.New(t)

	a.Equal([]string{"E", "D", "C", "B", "A", "genesis"}, assertList, "total order is not expected.")
}

// TestGetTotalOrderWithSpecialSeq test GetTotalOrder with special sequence
// genesis < A, B < C, D < E
// 问题：
//    1. List 元素顺序问题，java版的List元素索引按add顺序从小到大，现在是反过来的。
//    2. 得分计算问题，非pivot链的tip节点不会被纳入totalOrder。
func TestGetTotalOrderWithSpecialSeq(t *testing.T) {
	s := store.Storage{}
	// s.Init("/tmp/gorocksdb_test")
	d := Dag{}
	d.Init(&s)

	genesis := types.NewHashString("genesis")

	A := types.NewHashString("A")
	B := types.NewHashString("B")
	C := types.NewHashString("C")
	D := types.NewHashString("D")
	E := types.NewHashString("E")

	valueGenesis := types.Transaction{}
	refGenesis := types.List{}
	refGenesis.Add(types.NewHashString("X"))
	refGenesis.Add(types.NewHashString("Y"))
	valueGenesis.Init(refGenesis, types.NewHashString("data of genesis"))

	valueA := types.Transaction{}
	refA := &types.List{}
	refA.Add(genesis)
	refA.Add(genesis)
	valueA.Init(*refA, types.NewHashString("dataA"))

	valueB := types.Transaction{}
	refB := types.List{}
	refB.Add(genesis)
	refB.Add(genesis)
	valueB.Init(refB, types.NewHashString("dataB"))

	valueC := types.Transaction{}
	refC := types.List{}
	refC.Add(A)
	refC.Add(B)
	valueC.Init(refC, types.NewHashString("dataC"))

	valueD := types.Transaction{}
	refD := types.List{}
	refD.Add(A)
	refD.Add(B)
	valueD.Init(refD, types.NewHashString("dataD"))

	valueE := types.Transaction{}
	refE := types.List{}
	refE.Add(B)
	refE.Add(D)
	valueE.Init(refE, types.NewHashString("dataE"))

	d.Add(genesis, &valueGenesis)

	// FIXME block, transaction has no trunk and branck
	d.Add(A, &valueA)
	d.Add(B, &valueB)
	d.Add(C, &valueC)
	d.Add(D, &valueD)
	d.Add(E, &valueE)

	totalOrder := d.GetTotalOrder()
	t.Log(totalOrder)

	alias := map[types.Hash]string{}
	alias[genesis] = "genesis"
	alias[A] = "A"
	alias[B] = "B"
	alias[C] = "C"
	alias[D] = "D"
	alias[E] = "E"

	assertList := make([]string, 0)
	for _, h := range totalOrder {
		assertList = append(assertList, alias[h])
	}

	a := assert.New(t)

	a.Equal([]string{"E", "D", "A", "genesis"}, assertList, "total order is not expected.")
}

// List 的实现可能有坑，原有的逻辑下，先放进去的元素索引是0，然后依次递增。现有逻辑下先放进的元素索引是最大的。
func TestList(t *testing.T) {
	A := types.NewHashString("A")
	B := types.NewHashString("B")
	C := types.NewHashString("C")
	D := types.NewHashString("D")

	list := types.List{}
	list.Add(A)
	list.Add(B)
	list.Add(C)
	list.Add(D)

	a := assert.New(t)

	a.Equal(A, list.Index(3))
}
