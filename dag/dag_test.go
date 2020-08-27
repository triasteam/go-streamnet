package dag

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/triasteam/go-streamnet/store"
	"github.com/triasteam/go-streamnet/types"
)

func TestInit(t *testing.T) {
	s := store.Storage{}
	s.Init("/tmp/gorocksdb_test")
	d := Init(&s)

	t.Log(d)
}

func TestSave(t *testing.T) {
	s := store.Storage{}
	s.Init("/tmp/gorocksdb_test")
	d := Init(&s)

	key := types.NewHashString("genesis")
	value := types.FromHash(key)

	d.Add(key, value)
	pivot := d.getPivot(key)
	t.Log(pivot)
	a := assert.New(t)

	a.Equal(pivot, key, "there is only one genesis block")
}

func TestTotalOrder(t *testing.T) {
	s := store.Storage{}
	s.Init("/tmp/gorocksdb_test")
	d := Init(&s)

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
