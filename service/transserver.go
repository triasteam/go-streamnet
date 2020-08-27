package service

import (
	"encoding/json"

	"github.com/triasteam/go-streamnet/noderank"
	"github.com/triasteam/go-streamnet/store"
	"github.com/triasteam/go-streamnet/types"
)

type transService struct{}

func NewTransServer() *transService {
	return &transService{}
}

func (transServ *transService) StoreDagData(val string) string {
	redisStore := store.NewRedisStore() //init redisServerstreamstreamstrea
	hashKey := types.Sha256([]byte(val))
	var hashKeyStr = hashKey.String() //get hashkey
	redisStore.Set(hashKeyStr, val)   //store
	return hashKeyStr
}

func (transServ *transService) GetNodeRank(hashes []string, duration int, period int, numrank int) ([]noderank.Teescore, []noderank.Teectx, error) {
	redisStore := store.NewRedisStore()
	hashBlocks := make([]string, len(hashes))
	var block string

	for i, hashStr := range hashes {
		block = redisStore.Get(hashStr)
		hashBlocks[i] = block
	}
	blockJson, err := json.Marshal(hashBlocks)
	if err != nil {
		return nil, nil, err
	}

	noderankRequest := &noderank.GetRankRequest{
		Blocks:   string(blockJson),
		Duration: duration,
	}
	teescore, teectx, err := noderank.GetRank(noderankRequest, period, numrank)
	return teescore, teectx, nil
}
