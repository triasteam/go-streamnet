package types

import (
	"encoding/json"
	"fmt"
)

// StoreData is the content of input.
type StoreData struct {
	Attester string
	Attestee string
	Score    string
}

// String transfers the StoreData to String.
func (d StoreData) String() string {
	str, _ := json.Marshal(d)
	return string(str)
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

// QueryNodeReq is req param of QueryNodes
type QueryNodeReq struct {
	Period  uint32 `json:"period"`
	NumRank uint32 `json:"numrank"`
}

// TeeCtx contains vote detail
type TeeCtx struct {
	Attester string  `json:"attester"`
	Attestee string  `json:"attestee"`
	Score    float64 `json:"score"`
	Time     string  `json:"time,omitempty"`
	Nonce    int64   `json:"nonce,omitempty"`
}

// TeeScore contains score of attestee
type TeeScore struct {
	Attestee string  `json:"attestee"`
	Score    float64 `json:"score"`
}

// DataTee contains teescore and teectx
type DataTee struct {
	Teescore []TeeScore `json:"teescore"`
	Teectx   []TeeCtx   `json:"teectx"`
}

// Message return result
type Message struct {
	Code      uint32
	Timestamp int64
	Message   string
	Data      DataTee
}
