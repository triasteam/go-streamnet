package service

import (
	"encoding/json"

	"github.com/triasteam/go-streamnet/noderank"
	"github.com/triasteam/go-streamnet/store"
	"github.com/triasteam/go-streamnet/utils/crypto/tmhash"
)

type transService struct{}

func NewTransServer() *transService {
	return &transService{}
}

func (transServ *transService) StoreDagData(key string, val string) string {
	redisStore := store.NewRedisStore() //init redisServer
	hashUtil := tmhash.NewTruncated()   //init hashutil
	hashKey := hashUtil.Sum([]byte(key + val))
	var hashKeyStr = string(hashKey[:]) //get hashkey
	redisStore.Set(hashKeyStr, val)     //store
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
