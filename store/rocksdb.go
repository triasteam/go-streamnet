package store

import (
	"errors"
	"fmt"
	"log"

	"github.com/tecbot/gorocksdb"
	"github.com/triasteam/go-streamnet/utils/crypto/tmhash"
	//"strconv"
)

type Storage struct {
	db           *gorocksdb.DB
	readOptions  *gorocksdb.ReadOptions
	writeOptions *gorocksdb.WriteOptions
	//columnFamilies gorocksdb.ColumnFamilyHandles
}

func (store *Storage) Init(path string) error {
	db, err := OpenDB(path)
	if err != nil {
		fmt.Printf("Open RocksDB failed: %v!\n", err)
		return err
	}

	readOptions := gorocksdb.NewDefaultReadOptions()
	readOptions.SetFillCache(true)

	writeOptions := gorocksdb.NewDefaultWriteOptions()
	writeOptions.SetSync(true)

	store.db = db
	store.readOptions = readOptions
	store.writeOptions = writeOptions

	return nil
}

func (store *Storage) Save(key []byte, value []byte) error {
	err := store.db.Put(store.writeOptions, key, value)
	if err != nil {
		fmt.Println("Write data to RocksDB failed!")
		return err
	}

	return nil
}

func (store *Storage) SaveValue(value []byte) ([]byte, error) {
	key := tmhash.Sum(value)
	return key, store.Save(key, value)
}

func (store *Storage) Get(key []byte) ([]byte, error) {
	slice, err := store.db.Get(store.readOptions, key)
	if err != nil {
		fmt.Println("Get data from RocksDB failed!")
		return nil, err
	}
	return slice.Data(), nil
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

func (store *Storage) CloseDB() {
	store.db.Close()
}
