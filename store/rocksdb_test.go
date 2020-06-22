package store

import (
	"log"
	"strconv"
)

const (
	DB_PATH = "/tmp/gorocksdb"
)

func test() {
	db, err := OpenDB(DB_PATH)
	if err != nil {
		log.Println("fail to open db,", nil, db)
	}

	readOptions := gorocksdb.NewDefaultReadOptions()
	readOptions.SetFillCache(true)

	writeOptions := gorocksdb.NewDefaultWriteOptions()
	writeOptions.SetSync(true)

	for i := 0; i < 10000; i++ {
		keyStr := "aa" + strconv.Itoa(i)
		var key []byte = []byte(keyStr)
		db.Put(writeOptions, key, key)
		log.Println(i, keyStr)
		slice, err2 := db.Get(readOptions, key)
		if err2 != nil {
			log.Println("获取数据异常：", key, err2)
			continue
		}
		log.Println("获取数据：", slice.Size(), string(slice.Data()))
	}

	//defer readOptions.Destroy()
	//defer writeOptions.Destroy()
}
