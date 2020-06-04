package store

import "github.com/tecbot/gorocksdb"

type Storage struct {
	db gorocksdb.DB
	columnFamilies gorocksdb.ColumnFamilyHandles
}
