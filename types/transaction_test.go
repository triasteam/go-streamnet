package types

import (
	"log"
	"testing"
	"time"
)

func TestString(t *testing.T) {
	tx := Transaction{}
	tx.Timestamp = time.Now()
	tx.DataHash = NilHash
	s, err := tx.String()
	if err != nil {
		log.Fatal("String failed!")
	}
	log.Print(s)
}

func TestInit(t *testing.T) {
	tx := Transaction{}
	parents := List{}
	trunk := Sha256([]byte("StreamNet_Trunk"))
	branch := Sha256([]byte("StreamNet_Branch"))
	parents.Append(trunk)
	parents.Append(branch)
	tx.Init(parents, RandomHash())
	if tx.trunk != trunk || tx.branch != branch {
		log.Fatal("Init failed!")
	}
}
