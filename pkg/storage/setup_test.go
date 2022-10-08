package storage

import (
	"os"
	"testing"
)

var shardedMap ShardedMap

// var storeT = make(map[string]string)

func TestMain(m *testing.M) {

	shardedMap = NewShardedMap(3, 10)

	os.Exit(m.Run())
}
