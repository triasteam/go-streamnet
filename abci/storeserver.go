package abci

import (
	"github.com/triasteam/go-streamnet/store"
	"github.com/triasteam/go-streamnet/utils/crypto/tmhash"
)

type storeserver struct{}

func NewStoreserver() *storeserver{
	return &storeserver{};
}

func (storeServ * storeserver) StoreDagData(key string, val string) string{
	redisStore := store.NewRedisStore(); //init redisServer
	hashUtil := tmhash.NewTruncated();  //init hashutil
	hashKey := hashUtil.Sum([]byte(key+val));
	var hashKeyStr = string(hashKey[:]) //get hashkey
	redisStore.Set(hashKeyStr,val); //store
	return hashKeyStr;
}