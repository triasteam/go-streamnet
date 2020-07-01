package store

type StorageProvider interface {
	Init(path string) error
	Save(key []byte, value []byte) error
	SaveValue(value []byte) ([]byte, error)
	Get(key []byte) ([]byte, error)

	/*
		isAvailable() bool
		shutdown()
		delete(key types.Hash, value []byte)
		update(key types.Hash, value []byte)
		exists(key types.Hash, value []byte) bool
		saveBatch() bool
		deleteBatch()
		mayExist() bool
	*/

}
