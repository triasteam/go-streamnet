package types

import "github.com/triasteam/go-streamnet/dag"

// Transaction is a core struct.
type Transaction struct {
<<<<<<< HEAD
	bytes   []byte
	address Hash

	trunk  Hash
	branch Hash

	value     int64
	timestamp int64

	weightMagnitude int64
	nonce           []byte

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

	return s
=======
	bytes                         []byte
	address                       Hash
	trunk                         Hash
	branch                        Hash
	obsoleteTag                   Hash
	value                         int64
	currentIndex                  int64
	lastIndex                     int64
	timestamp                     int64
	tag                           Hash
	attachmentTimestamp           int64
	attachmentTimestampLowerBound int64
	attachmentTimestampUpperBound int64
	height                        int64
	sender                        string
	weightMagnitude               int64
	nonce                         []byte
>>>>>>> 688dc7a... implement type 'Set'
}

func FromHash(dag *dag.Dag, hash Hash) *Transaction {
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

	return s
}
