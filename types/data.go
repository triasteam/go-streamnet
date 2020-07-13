package types

import "fmt"

type StoreData struct {
	Attester string
	Attestee string
	Score    string
}

func (d StoreData) String() string {
	return fmt.Sprintf("Attester: %s, Attestee: %s, Score: %s", d.Attester, d.Attestee, d.Score)
}

type StoreReply struct {
	Code int
	Hash string
}

func (r StoreReply) String() string {
	return fmt.Sprintf("Code: %d, Hash: %s", r.Code, r.Hash)
}

type GetReq struct {
	Key string
}

type GetReply struct {
	Value string
}
