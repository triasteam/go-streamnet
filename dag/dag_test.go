package dag

import (
	"log"
	"testing"

	"github.com/triasteam/go-streamnet/config"

	"github.com/triasteam/go-streamnet/store"
	"github.com/triasteam/go-streamnet/types"
)

func TestInit(t *testing.T) {
	s := store.Storage{}
	defer s.Close()

	s.Init("/tmp/gorocksdb_test")
	d := Dag{}
	d.Init(&s)

	t.Log(d)
}

func TestAdd(t *testing.T) {
	// storage
	s := store.Storage{}
	defer s.Close()
	s.Init("/tmp/gorocksdb_test")

	// dag
	d := Dag{}
	d.Init(&s)

	// transaction
	trunk := config.GenesisTrunk
	branch := config.GenesisBranch
	l := types.List{}
	l.Append(trunk)
	l.Append(branch)
	tx := types.Transaction{}
	tx.Init(l, types.NewHash([]byte("data")))

	bytes, _ := tx.Bytes()
	key := types.Sha256(bytes)

	// add
	d.Add(key, &tx)

	// test updateGraph
	t.Log("graph: ", d.graph)
	t.Log("parentGraph: ", d.parentGraph)
	t.Log("revGraph: ", d.revGraph)
	t.Log("parentRevGraph: ", d.parentRevGraph)
	t.Log("degrees: ", d.degrees)

	// test updateTopologicalOrder
	t.Log("topOrderStreaming: ", d.topOrderStreaming)
	t.Log("levels: ", d.levels)
	t.Log("totalDepth: ", d.totalDepth)

	// test updateScore
	t.Log("score: ", d.score)
	t.Log("parentScore: ", d.parentScore)

	// test GetGenesis
	g := d.GetGenesis()
	t.Log("Genesis: ", g)

	// test GetPivotalHash
	h := d.GetPivotalHash(0)
	log.Print("Pivotal hash: ", h)
	if key != h {
		t.Fatal("Genesis is wrong!")
	}

	// test GetLastPivot
	last := d.GetLastPivot(g)
	if last != g {
		t.Fatal("LastPivot is wrong!")
	}
	t.Log(last)

	// test Contains
	if !d.Contains(g) {
		t.Fatal("Constains failed!")
	}

	// test GetScore
	score := d.GetScore(key)
	if score != 1 {
		t.Fatal("Score is wrong!")
	}

	// test GetChildren
	children := d.GetChildren(key)
	if children.Len() != 0 {
		t.Fatal("GetChildren is wrong!")
	}
}
