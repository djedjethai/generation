package storage

import (
	"context"
	"fmt"
	"github.com/djedjethai/generation/internal/models"
	"github.com/djedjethai/generation/internal/observability"
	"testing"
)

// test Put
func TestPut(t *testing.T) {
	ctx := context.Background()

	_ = shardedMap.Set(ctx, "test", "put")

	shard := shardedMap.getShard("test")

	// fmt.Println("grrr: ", shardedMap.Keys())
	if shard.m["test"].val != "put" {
		t.Error("err in store Put() failed")
	}
}

// test Get
func TestGet(t *testing.T) {

	ctx := context.Background()

	dt, _ := shardedMap.Get(ctx, "test")

	if dt != "put" {
		t.Error("err in store Get() failed")
	}
}

// test get all Keys
func TestKeys(t *testing.T) {

	ctx := context.Background()

	keys := shardedMap.Keys(ctx)

	fmt.Println("see keys: ", keys)
	if len(keys) != 1 && keys[0] != "test" {
		t.Error("err in store Keys() failed")
	}
}

// test delete
func TestDelete(t *testing.T) {

	ctx := context.Background()

	_ = shardedMap.Delete(ctx, "test", nil)

	shard := shardedMap.getShard("test")

	if len(shard.m) > 0 {
		t.Error("err in store Delete() failed")
	}
}

// test getKeysValues
func TestGetKeysValues(t *testing.T) {

	ctx := context.Background()

	shardedMap.Set(ctx, "key1", "val1")
	shardedMap.Set(ctx, "key2", "val2")
	shardedMap.Set(ctx, "key3", "val3")

	kv := make(chan models.KeysValues)

	go func() {
		err := shardedMap.KeysValues(ctx, kv)
		if err != nil {
			fmt.Println("err: ", err)
		}
	}()

	var resKeys []string
	var resValues []string
	for v := range kv {
		t.Run(v.Key, func(t *testing.T) {
			resKeys = append(resKeys, v.Key)
			resValues = append(resValues, v.Value)
		})
	}

	if len(resKeys) != 3 || len(resValues) != 3 {
		t.Error("err in store TestGetKeysValues, 3 key-value pairs should be returned")
	}
}

// make sure the fixed size is respected when one shard and many items for this shard
func TestStorageKeepTheSettedSizeWithOneShardAnsManyItemsPerShard(t *testing.T) {

	ctx := context.Background()
	obs := observability.Observability{}

	sm := NewShardedMap(1, 2, obs)
	sm.Set(ctx, "key1", "val1")
	sm.Set(ctx, "key2", "val2")
	sm.Set(ctx, "key3", "val3")
	sm.Set(ctx, "key4", "val4")

	// check the remained element into the dll
	_, t1 := sm.shd[0].m["key1"]
	_, t2 := sm.shd[0].m["key2"]
	_, t3 := sm.shd[0].m["key3"]
	_, t4 := sm.shd[0].m["key4"]
	if t1 || t2 || !t3 || !t4 {
		t.Error("err t3 or/and t4 in store TestStorageKeepTheSettedSizeWithOneShardAnsManyItemsPerShard")
	}

	ks := sm.Keys(ctx)
	if len(ks) != 2 || ks[0] != "key3" || ks[1] != "key4" {
		t.Error("err in store TestStorageKeepTheSettedSizeWithOneShardAnsManyItemsPerShard")
	}
}

// make sure the fixed size is respected when many shards and a single item for each shard
// func TestStorageKeepTheSettedSizeManyShardAnsOneItemPerShard(t *testing.T) {
//
// 	ctx := context.Background()
// 	obs := observability.Observability{}
//
// 	sm := NewShardedMap(2, 1, obs)
// 	sm.Set(ctx, "key1", "val1")
// 	sm.Set(ctx, "key2", "val2")
// 	sm.Set(ctx, "key3", "val3")
// 	sm.Set(ctx, "key4", "val4")
//
// 	ks := sm.Keys(ctx)
// 	if len(ks) != 2 || ks[0] != "key4" || ks[1] != "key3" {
// 		t.Error("err in store TestStorageKeepTheSettedSizeManyShardAnsOneItemPerShard")
// 	}
// }

// make sure a key won't be repeated twice
func TestStorageDoNotStoreTheSameKeyTwice(t *testing.T) {

	ctx := context.Background()
	obs := observability.Observability{}

	sm := NewShardedMap(2, 2, obs)
	sm.Set(ctx, "key1", "val1")
	sm.Set(ctx, "key2", "val2")
	sm.Set(ctx, "key3", "val3")
	sm.Set(ctx, "key3", "val3")
	sm.Set(ctx, "key4", "val4")
	sm.Set(ctx, "key1", "val1")
	sm.Set(ctx, "key2", "val2")
	sm.Set(ctx, "key1", "val1")

	ks := sm.Keys(ctx)
	if len(ks) != 4 {
		var val1 = false
		var val2 = false
		var val3 = false
		var val4 = false
		for _, v := range ks {
			switch v {
			case "val1":
				val1 = true
			case "val2":
				val2 = true
			case "val3":
				val3 = true
			case "val4":
				val4 = true
			}
		}
		if !val1 || !val2 || !val3 || !val4 {
			t.Error("err in store TestStorageDoNotStoreTheSameKeyTwice")
		}
	}
}