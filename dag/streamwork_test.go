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
	value := types.FromHash(genesis)

	A := types.NewHashString("A")
	B := types.NewHashString("B")
	C := types.NewHashString("C")
	D := types.NewHashString("D")
	E := types.NewHashString("E")

	d.Add(genesis, value)

	// FIXME block, transaction has no trunk and branck
	d.Add(A, value)
	d.Add(B, value)
	d.Add(C, value)
	d.Add(D, value)
	d.Add(E, value)

	pivot := d.getPivot(E)
	t.Log(pivot)
	a := assert.New(t)

	a.Equal(pivot, genesis, "pivot is not genesis")
}

// transaction not implement, holding.
// func TestGetTotalOrder(t *testing.T) {
// 	s := store.Storage{}
// 	s.Init("/tmp/gorocksdb_test")
// 	d := Init(&s)

// 	genesis := types.NewHashString("genesis")
// 	value := types.FromHash(genesis)

// 	A := types.NewHashString("A")
// 	B := types.NewHashString("B")
// 	C := types.NewHashString("C")
// 	D := types.NewHashString("D")
// 	E := types.NewHashString("E")

// 	d.Add(genesis, value)

// 	// FIXME block, transaction has no trunk and branck
// 	tranA := types.FromHash(A)
// 	tranA.GetApprovers().Add(genesis)
// 	tranA.GetApprovers().Add(genesis)
// 	d.Add(A, tranA)
// 	tranB := types.FromHash(B)
// 	tranB.GetApprovers().Add(A)
// 	tranB.GetApprovers().Add(A)
// 	d.Add(B, tranB)
// 	tranC := types.FromHash(C)
// 	tranC.GetApprovers().Add(B)
// 	tranC.GetApprovers().Add(B)
// 	d.Add(C, tranC)
// 	tranD := types.FromHash(D)
// 	tranD.GetApprovers().Add(C)
// 	tranD.GetApprovers().Add(C)
// 	d.Add(D, tranD)
// 	tranE := types.FromHash(E)
// 	tranE.GetApprovers().Add(D)
// 	tranE.GetApprovers().Add(D)
// 	d.Add(E, value)

// 	pivot := d.getPivot(E)
// 	t.Log(pivot)
// 	a := assert.New(t)

// 	a.Equal(pivot, genesis, "pivot is not genesis")
// }

// test totalOrder, G <A,B<C,D<E
func TestGetTotalOrder(t *testing.T) {
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

// simple dag: G<A<B<C<D<E
func TestGetTotalOrder2(t *testing.T) {
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
	d.graph[B].Add(A)
	d.graph[B].Add(A)

	d.graph[C] = types.NewSet()
	d.graph[C].Add(B)
	d.graph[C].Add(B)

	d.graph[D] = types.NewSet()
	d.graph[D].Add(C)
	d.graph[D].Add(C)

	d.graph[E] = types.NewSet()
	d.graph[E].Add(D)
	d.graph[E].Add(D)

	// file revGraph
	d.revGraph[X] = types.NewSet()
	d.revGraph[X].Add(genesis)

	d.revGraph[Y] = types.NewSet()
	d.revGraph[Y].Add(genesis)

	d.revGraph[genesis] = types.NewSet()
	d.revGraph[genesis].Add(A)

	d.revGraph[A] = types.NewSet()
	d.revGraph[A].Add(B)

	d.revGraph[B] = types.NewSet()
	d.revGraph[B].Add(C)

	d.revGraph[C] = types.NewSet()
	d.revGraph[C].Add(D)

	d.revGraph[D] = types.NewSet()
	d.revGraph[D].Add(E)
	// fill parentGraph
	d.parentGraph[genesis] = X
	d.parentGraph[A] = genesis
	d.parentGraph[B] = A
	d.parentGraph[C] = B
	d.parentGraph[D] = C
	d.parentGraph[E] = D
	// fill parentRevGraph
	d.parentRevGraph[X] = types.NewSet()
	d.parentRevGraph[X].Add(genesis)

	d.parentRevGraph[genesis] = types.NewSet()
	d.parentRevGraph[genesis].Add(A)

	d.parentRevGraph[A] = types.NewSet()
	d.parentRevGraph[A].Add(B)

	d.parentRevGraph[B] = types.NewSet()
	d.parentRevGraph[B].Add(C)

	d.parentRevGraph[C] = types.NewSet()
	d.parentRevGraph[C].Add(D)

	d.parentRevGraph[D] = types.NewSet()
	d.parentRevGraph[D].Add(E)
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
