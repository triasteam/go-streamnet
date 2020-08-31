package types

import (
	"log"
	"testing"
	"time"
)

func TestTransaction_Init(t *testing.T) {
	tx := Transaction{}
	parents := List{}
	trunk := Sha256([]byte("StreamNet_Trunk"))
	branch := Sha256([]byte("StreamNet_Branch"))
	parents.Append(trunk)
	parents.Append(branch)
	tx.Init(parents, RandomHash())
	if tx.Trunk != trunk || tx.Branch != branch {
		log.Fatal("Init failed!")
	}
}

func TestTransaction_String(t *testing.T) {
	tx := Transaction{}
	tx.Timestamp = time.Now()
	tx.DataHash = NilHash
	tx.Trunk = Sha256([]byte("StreamNet_Trunk"))
	tx.Branch = Sha256([]byte("StreamNet_Branch"))
	s, err := tx.String()
	if err != nil {
		log.Fatal("String failed!")
	}
	log.Print(s)
}

func TestTransactionFromBytes(t *testing.T) {
	tx := Transaction{}
	parents := List{}
	trunk := Sha256([]byte("StreamNet_Trunk"))
	branch := Sha256([]byte("StreamNet_Branch"))
	parents.Append(trunk)
	parents.Append(branch)
	tx.Init(parents, RandomHash())

	b, err := tx.Bytes()
	if err != nil {
		log.Fatal("Bytes failed!")
	}

	tx1 := TransactionFromBytes(b)

	if tx1.DataHash != tx.DataHash ||
		tx1.Trunk != tx.Trunk ||
		tx1.Branch != tx.Branch ||
		!tx1.Timestamp.Equal(tx.Timestamp) {
		t.Fatal("Marshal and Unmarshal failed!")
	}
}
