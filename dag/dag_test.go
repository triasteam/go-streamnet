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
	d := Dag{}
	d.Init(&s)

	t.Log(d)
}

func TestSave(t *testing.T) {
	s := store.Storage{}
	s.Init("/tmp/gorocksdb_test")
	d := Dag{}
	d.Init(&s)

	key := types.NewHashString("genesis")
	value := types.FromHash(key)

	d.Add(key, value)
	pivot := d.getPivot(key)
	t.Log(pivot)

	a := assert.New(t)

	a.Equal(pivot, key, "there is only one genesis block")
}
