package service

import (
	"github.com/triasteam/go-streamnet/store"
	"github.com/triasteam/go-streamnet/utils/crypto/tmhash"
)

type storeService struct{}

func NewStoreserver() *storeService {
	return &storeService{}
}

func (storeServ *storeService) StoreDagData(key string, val string) string {
	redisStore := store.NewRedisStore() //init redisServer
	hashUtil := tmhash.NewTruncated()   //init hashutil
	hashKey := hashUtil.Sum([]byte(key + val))
	var hashKeyStr = string(hashKey[:]) //get hashkey
	redisStore.Set(hashKeyStr, val)     //store
	return hashKeyStr
}
