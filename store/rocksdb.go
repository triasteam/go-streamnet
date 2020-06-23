package store

import (
	"errors"
	"fmt"
	"github.com/tecbot/gorocksdb"
	"log"
	//"strconv"
)

type Storage struct {
	db           *gorocksdb.DB
	readOptions  *gorocksdb.ReadOptions
	writeOptions *gorocksdb.WriteOptions
	//columnFamilies gorocksdb.ColumnFamilyHandles
}

func Init(path string) (*Storage, error) {
	db, err := OpenDB(path)
	if err != nil {
		fmt.Println("Open RocksDB failed: %v!", err)
		return nil, err
	}

	readOptions := gorocksdb.NewDefaultReadOptions()
	readOptions.SetFillCache(true)

	writeOptions := gorocksdb.NewDefaultWriteOptions()
	writeOptions.SetSync(true)

	return &Storage{
		db,
		readOptions,
		writeOptions,
	}, nil
}

func (store *Storage) Save(key []byte, value []byte) error {
	err := store.db.Put(store.writeOptions, key, value)
	if err != nil {
		fmt.Println("Write data to RocksDB failed!")
		return err
	}
}

func (store *Storage) Get(key []byte) []byte {
	slice, err := store.db.Get(store.readOptions, key)
	if err != nil {
		fmt.Println("Get data from RocksDB failed!")
		return nil
	}
	return slice.Data()
}

func OpenDB(path string) (*gorocksdb.DB, error) {
	options := gorocksdb.NewDefaultOptions()
	options.SetCreateIfMissing(true)

	bloomFilter := gorocksdb.NewBloomFilter(10)

	readOptions := gorocksdb.NewDefaultReadOptions()
	readOptions.SetFillCache(false)

	rateLimiter := gorocksdb.NewRateLimiter(10000000, 10000, 10)
	options.SetRateLimiter(rateLimiter)
	options.SetCreateIfMissing(true)
	options.EnableStatistics()
	//options.SetWriteBufferSize(8 * gorocksdb.KB)
	options.SetMaxWriteBufferNumber(3)
	options.SetMaxBackgroundCompactions(10)
	//options.SetCompression(gorocksdb.SnappyCompression)
	options.SetCompactionStyle(gorocksdb.UniversalCompactionStyle)

	options.SetHashSkipListRep(2000000, 4, 4)

	blockBasedTableOptions := gorocksdb.NewDefaultBlockBasedTableOptions()
	//blockBasedTableOptions.SetBlockCache(gorocksdb.NewLRUCache(64 * gorocksdb.KB))
	blockBasedTableOptions.SetFilterPolicy(bloomFilter)
	blockBasedTableOptions.SetBlockSizeDeviation(5)
	blockBasedTableOptions.SetBlockRestartInterval(10)
	//blockBasedTableOptions.SetBlockCacheCompressed(gorocksdb.NewLRUCache(64 * gorocksdb.KB))
	blockBasedTableOptions.SetCacheIndexAndFilterBlocks(true)
	blockBasedTableOptions.SetIndexType(gorocksdb.KHashSearchIndexType)

	options.SetBlockBasedTableFactory(blockBasedTableOptions)
	//log.Println(bloomFilter, readOptions)
	options.SetPrefixExtractor(gorocksdb.NewFixedPrefixTransform(3))

	options.SetAllowConcurrentMemtableWrites(false)
	db, err := gorocksdb.OpenDb(options, path)

	if err != nil {
		log.Fatalln("OPEN DB error", db, err)
		db.Close()
		return nil, errors.New("fail to open db")
	} else {
		log.Println("OPEN DB success", db)
	}
	return db, nil
}
