package storage

import (
	"os"
	"testing"
)

var storageT StorageRepo
var storeT = make(map[string]string)

func TestMain(m *testing.M) {

	storageT = NewShardedMap(3, 10)

	os.Exit(m.Run())
}
