package types

import (
	"encoding/json"
	"log"
	"time"
)

// Transaction is a core struct.
type Transaction struct {
	address Hash

	DataHash Hash

	trunk  Hash
	branch Hash

	Timestamp time.Time

	weightMagnitude int64
	nonce           []byte

	//bytes   []byte
	//value     int64
	//currentIndex                  int64
	//lastIndex                     int64
	//attachmentTimestamp           int64
	//attachmentTimestampLowerBound int64
	//attachmentTimestampUpperBound int64
	//obsoleteTag Hash
	//tag       Hash
	//sender          string
	//height int64
}

func (tx *Transaction) Init(parents List) {
	tx.trunk = parents.Index(0)
	tx.branch = parents.Index(1)
}

func FromHash(hash Hash) *Transaction {
	var tx Transaction

	return &tx
}

// GetTrunkTransactionHash returns the trunk hash.
func (tx *Transaction) GetTrunkTransactionHash() Hash {
	if tx == nil {
		return NewHash(nil)
	}
	return tx.trunk
}

// GetBranchTransactionHash returns the branch hash.
func (tx *Transaction) GetBranchTransactionHash() Hash {
	if tx == nil {
		return NewHash(nil)
	}
	return tx.branch
}

func (tx *Transaction) GetApprovers() Set {
	s := NewSet()
	/*s.Add(tx.trunk)
	s.Add(tx.branch)*/
	return s
}

func (tx *Transaction) String() string {
	b, err := json.Marshal(tx)
	if err != nil {
		log.Print("Json failed!")
		return ""
	}
	return string(b)
}
