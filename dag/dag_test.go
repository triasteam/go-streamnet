package dag

import (
	"testing"

	"github.com/triasteam/go-streamnet/store"
)

func TestInit(t *testing.T) {
	s := store.Init("/tmp/gorocksdb_test")
	d := Init(s)
}
