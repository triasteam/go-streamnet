package types

import (
	"encoding/json"
	"log"
	"time"
)

// Transaction is a core struct.
type Transaction struct {
	Address Hash `json:"address"`

	DataHash Hash `json:"DataHash"`

	Trunk  Hash `json:"trunk"`
	Branch Hash `json:"branch"`

	Timestamp time.Time `json:"timestamp"`

	WeightMagnitude int64  `json:"weightMagnitude"`
	Nonce           []byte `json:"nonce"`

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

func (tx *Transaction) Init(parents List, value Hash) {
	tx.Trunk = parents.Index(0)
	tx.Branch = parents.Index(1)

	// timestamp
	tx.Timestamp = time.Now()

	tx.DataHash = value
}

func FromHash(hash Hash) *Transaction {
	var tx Transaction

	return &tx
}

// GetTrunkTransactionHash returns the trunk hash.
func (tx *Transaction) GetTrunkTransactionHash() Hash {
	if tx == nil {
		return NilHash
	}
	return tx.Trunk
}

// GetBranchTransactionHash returns the branch hash.
func (tx *Transaction) GetBranchTransactionHash() Hash {
	if tx == nil {
		return NilHash
	}
	return tx.Branch
}

/*func (tx *Transaction) GetApprovers() Set {
	s := NewSet()
	s.Add(tx.Trunk)
	s.Add(tx.Branch)
	return s
}*/

func (tx *Transaction) Bytes() ([]byte, error) {
	b, err := json.Marshal(tx)
	if err != nil {
		log.Print("Json failed!")
		return nil, err
	}
	return b, nil
}

func (tx *Transaction) String() (string, error) {
	b, err := tx.Bytes()
	if err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}

func TransactionFromBytes(b []byte) *Transaction {
	tx := Transaction{}
	err := json.Unmarshal(b, &tx)
	if err != nil {
		log.Print("UnJson failed!")
		return nil
	}
	return &tx
}
