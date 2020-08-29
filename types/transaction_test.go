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
	s := tx.String()
	if s == "" {
		log.Fatal("String failed!")
	}
	log.Print(s)
}
