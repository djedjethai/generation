package storage

import (
	"os"
	"testing"
)

var storageT StorageRepo
var storeT = make(map[string]string)

func TestMain(m *testing.M) {

	storageT = NewStorage(storeT)

	os.Exit(m.Run())
}
