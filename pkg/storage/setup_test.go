package storage

import (
	"github.com/djedjethai/generation/pkg/observability"
	"os"
	"testing"
)

var shardedMap ShardedMap

// var storeT = make(map[string]string)

func TestMain(m *testing.M) {

	obs := observability.Observability{}

	shardedMap = NewShardedMap(3, 10, obs)

	os.Exit(m.Run())
}
