package storage

import (
	"fmt"
	"testing"
)

// test Put
func TestPut(t *testing.T) {
	_ = shardedMap.Set("test", "put")

	shard := shardedMap.getShard("test")

	// fmt.Println("grrr: ", shardedMap.Keys())
	if shard.m["test"].val != "put" {
		t.Error("err in store Put() failed")
	}
}

// test Get
func TestGet(t *testing.T) {
	dt, _ := shardedMap.Get("test")

	if dt != "put" {
		t.Error("err in store Get() failed")
	}
}

// test get all Keys
func TestKeys(t *testing.T) {
	keys := shardedMap.Keys()

	fmt.Println("see keys: ", keys)
	if len(keys) != 1 && keys[0] != "test" {
		t.Error("err in store Keys() failed")
	}
}

// test delete
func TestDelete(t *testing.T) {
	_ = shardedMap.Delete("test")

	shard := shardedMap.getShard("test")

	if len(shard.m) > 0 {
		t.Error("err in store Delete() failed")
	}
}

// make sure the fixed size is respected when one shard and many items for this shard
func TestStorageKeepTheSettedSizeWithOneShardAnsManyItemsPerShard(t *testing.T) {
	sm := NewShardedMap(1, 2)
	sm.Set("key1", "val1")
	sm.Set("key2", "val2")
	sm.Set("key3", "val3")
	sm.Set("key4", "val4")

	ks := sm.Keys()
	if len(ks) != 2 || ks[0] != "key3" || ks[1] != "key4" {
		t.Error("err in store TestStorageKeepTheSettedSizeWithOneShardAnsManyItemsPerShard")
	}
}

// make sure the fixed size is respected when many shards and a single item for each shard
func TestStorageKeepTheSettedSizeManyShardAnsOneItemPerShard(t *testing.T) {
	sm := NewShardedMap(2, 1)
	sm.Set("key1", "val1")
	sm.Set("key2", "val2")
	sm.Set("key3", "val3")
	sm.Set("key4", "val4")

	ks := sm.Keys()
	if len(ks) != 2 || ks[0] != "key4" || ks[1] != "key3" {
		t.Error("err in store TestStorageKeepTheSettedSizeManyShardAnsOneItemPerShard")
	}
}

// make sure a key won't be repeated twice
func TestStorageDoNotStoreTheSameKeyTwice(t *testing.T) {
	sm := NewShardedMap(2, 2)
	sm.Set("key1", "val1")
	sm.Set("key2", "val2")
	sm.Set("key3", "val3")
	sm.Set("key3", "val3")
	sm.Set("key4", "val4")
	sm.Set("key1", "val1")
	sm.Set("key2", "val2")
	sm.Set("key1", "val1")

	ks := sm.Keys()
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
