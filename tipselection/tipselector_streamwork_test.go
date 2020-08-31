package tipselection

import (
	"testing"

	"github.com/triasteam/go-streamnet/dag"
	"github.com/triasteam/go-streamnet/store"
	"github.com/triasteam/go-streamnet/types"
)

func TestGetTransactionsToApprove(t *testing.T) {
	// storage
	s := store.Storage{}
	defer s.Close()
	s.Init("/tmp/gorocksdb_test")

	// dag
	d := dag.Dag{}
	d.Init(&s)

	// tipselection
	ts := TipSelectorStreamWork{}
	ts.Init(&d)

	// transaction 1
	tx1 := types.Transaction{}
	tips := ts.GetTransactionsToApprove(16, types.NilHash)
	t.Log(tips.Length())
	for i := 0; i < tips.Length(); i++ {
		t.Log(tips.Index(i))
	}

	tx1.Init(tips, types.RandomHash())
	bytes, _ := tx1.Bytes()
	key := types.Sha256(bytes)
	t.Log("tx1 hash = ", key)
	d.Add(key, &tx1)

	// transaction 2
	tx2 := types.Transaction{}
	tips = ts.GetTransactionsToApprove(16, types.NilHash)
	t.Log(tips.Length())
	for i := 0; i < tips.Length(); i++ {
		t.Log(tips.Index(i))
	}
	tx2.Init(tips, types.RandomHash())

	bytes, _ = tx2.Bytes()
	key = types.Sha256(bytes)
	t.Log("tx2 hash = ", key)
	d.Add(key, &tx2)
}
