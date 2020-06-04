package store

import "github.com/triasteam/StreamNet-go/types"

type StorageProvider interface {
	init()
	isAvailable() bool
	shutdown()
    save(key types.Hash, value []byte) bool
	delete(key types.Hash, value []byte)
	update(key types.Hash, value []byte)
    exists(key types.Hash, value []byte) bool
	saveBatch() bool
    deleteBatch()
	get()
	mayExist() bool
}
