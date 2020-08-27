package types

import "fmt"

// StoreData is the content of input.
type StoreData struct {
	Attester string
	Attestee string
	Score    string
}

// String transfers the StoreData to String.
func (d StoreData) String() string {
	return fmt.Sprintf("Attester: %s, Attestee: %s, Score: %s", d.Attester, d.Attestee, d.Score)
}

// StoreReply is the struct to reply to user.
type StoreReply struct {
	Code int
	Hash string
}

// String transfers the StoreReply to String.
func (r StoreReply) String() string {
	return fmt.Sprintf("Code: %d, Hash: %s", r.Code, r.Hash)
}

// GetReq is the struct of request of 'get'.
type GetReq struct {
	Key string
}

// GetReply is the struct of reply of 'get'.
type GetReply struct {
	Value string
}
