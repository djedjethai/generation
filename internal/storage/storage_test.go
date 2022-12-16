package storage

import (
	"context"
	"github.com/djedjethai/generation/internal/models"
	"github.com/djedjethai/generation/internal/observability"
	"testing"
)

// test Put
func Test_put(t *testing.T) {
	ctx := context.Background()

	_ = shardedMap.Set(ctx, "test", "put")

	shard := shardedMap.getShard("test")

	if shard.m["test"].val != "put" {
		t.Error("err in store Put() failed")
	}
}

// test Get
func Test_get(t *testing.T) {

	ctx := context.Background()

	dt, _ := shardedMap.Get(ctx, "test")

	if dt != "put" {
		t.Error("err in store Get() failed")
	}
}

// test get all Keys
func Test_keys(t *testing.T) {

	ctx := context.Background()

	keys := shardedMap.Keys(ctx)

	if len(keys) != 1 && keys[0] != "test" {
		t.Error("err in store Keys() failed")
	}
}

// test delete
func Test_delete(t *testing.T) {

	ctx := context.Background()

	_ = shardedMap.Delete(ctx, "test", nil)

	shard := shardedMap.getShard("test")

	if len(shard.m) > 0 {
		t.Error("err in store Delete() failed")
	}
}

// test getKeysValues
func Test_get_keys_values(t *testing.T) {

	ctx := context.Background()

	shardedMap.Set(ctx, "key1", "value1")
	shardedMap.Set(ctx, "key2", "value2")
	shardedMap.Set(ctx, "key3", "value3")

	kv := make(chan models.KeysValues, 4)

	err := shardedMap.KeysValues(ctx, kv)
	if err != nil {
		t.Error("err in store TestGetKeysValues, .KeysValues() return an err: ", err)
	}

	var res = make(map[string]string)
	for v := range kv {
		t.Run(v.Key, func(t *testing.T) {
			res[v.Key] = v.Value
		})
	}

	if res["key1"] != "value1" || res["key2"] != "value2" || res["key3"] != "value3" {
		t.Error("err in store TestGetKeysValues, 3 key-value pairs should be returned")
	}
}

// make sure the fixed size is respected when one shard and many items for this shard
func Test_storage_keep_the_setted_size_with_one_shard_and_many_items_per_shard(t *testing.T) {

	ctx := context.Background()
	obs := observability.Observability{}

	sm := NewShardedMap(1, 3, &obs)
	sm.Set(ctx, "key1", "val1")
	sm.Set(ctx, "key2", "val2")
	sm.Set(ctx, "key3", "val3")
	sm.Set(ctx, "key4", "val4")

	// check the remained element into the dll
	_, t1 := sm.shd[0].m["key1"]
	_, t2 := sm.shd[0].m["key2"]
	_, t3 := sm.shd[0].m["key3"]
	_, t4 := sm.shd[0].m["key4"]
	if t1 || !t2 || !t3 || !t4 {
		t.Error("err t3 or/and t4 in store TestStorageKeepTheSettedSizeWithOneShardAnsManyItemsPerShard")
	}

	head := sm.shd[0].dll.head.val
	middle := sm.shd[0].dll.head.next.val
	tail := sm.shd[0].dll.tail.val

	if head != "val4" || middle != "val3" || tail != "val2" {
		t.Error("err in store TestStorageKeepTheSettedSizeWithOneShardAnsManyItemsPerShard")
	}
}

// make sure the fixed size is respected when one shard and many items for this shard
func Test_storage_delete_and_unshift_item_when_item_already_exist(t *testing.T) {

	ctx := context.Background()
	obs := observability.Observability{}

	sm := NewShardedMap(1, 3, &obs)
	sm.Set(ctx, "key1", "val1")
	sm.Set(ctx, "key2", "val2")
	sm.Set(ctx, "key3", "val3")
	sm.Set(ctx, "key4", "val4")
	sm.Set(ctx, "key3", "val3")

	// check the remained element into the dll
	_, t1 := sm.shd[0].m["key1"]
	_, t2 := sm.shd[0].m["key2"]
	_, t3 := sm.shd[0].m["key3"]
	_, t4 := sm.shd[0].m["key4"]
	if t1 || !t2 || !t3 || !t4 {
		t.Error("err t3 or/and t4 in store TestStorageKeepTheSettedSizeWithOneShardAnsManyItemsPerShard")
	}

	head := sm.shd[0].dll.head.val
	middle := sm.shd[0].dll.head.next.val
	tail := sm.shd[0].dll.tail.val

	if head != "val3" || middle != "val4" || tail != "val2" {
		t.Error("err in store TestStorageKeepTheSettedSizeWithOneShardAnsManyItemsPerShard")
	}
}

// make sure the last item is removed(from dll and map) from the list(if over storage size)
func Test_item_has_been_properly_removed_when_outbound_the_storage_size(t *testing.T) {

	ctx := context.Background()
	obs := observability.Observability{}

	sm := NewShardedMap(1, 2, &obs)
	sm.Set(ctx, "key1", "val1")
	sm.Set(ctx, "key3", "val3")
	sm.Set(ctx, "key3", "val3")
	sm.Set(ctx, "key4", "val4")

	head := sm.shd[0].dll.head.val
	headNext := sm.shd[0].dll.head.next.val
	tail := sm.shd[0].dll.tail.val
	tailPrev := sm.shd[0].dll.tail.prev.val
	_, map1 := sm.shd[0].m["key1"]
	_, map2 := sm.shd[0].m["key2"]
	_, map3 := sm.shd[0].m["key3"]
	_, map4 := sm.shd[0].m["key4"]
	if head != "val4" ||
		tail != "val3" ||
		headNext != tail ||
		head != tailPrev ||
		map1 || map2 || !map3 || !map4 {
		t.Error("err in store test ItemHasBeenProperlyRemovedWhenOutboud1")
	}

	sm.Set(ctx, "key3", "val3")
	head = sm.shd[0].dll.head.val
	headNext = sm.shd[0].dll.head.next.val
	tail = sm.shd[0].dll.tail.val
	tailPrev = sm.shd[0].dll.tail.prev.val
	_, map3 = sm.shd[0].m["key3"]
	_, map4 = sm.shd[0].m["key4"]
	if head != "val3" ||
		tail != "val4" ||
		headNext != tail ||
		head != tailPrev ||
		!map3 || !map4 {
		t.Error("err in store test ItemHasBeenProperlyRemovedWhenOutboud2")
	}

	sm.Set(ctx, "key3", "val3")
	head = sm.shd[0].dll.head.val
	headNext = sm.shd[0].dll.head.next.val
	tail = sm.shd[0].dll.tail.val
	tailPrev = sm.shd[0].dll.tail.prev.val
	_, map3 = sm.shd[0].m["key3"]
	_, map4 = sm.shd[0].m["key4"]
	if head != "val3" ||
		tail != "val4" ||
		headNext != tail ||
		head != tailPrev ||
		!map3 || !map4 {
		t.Error("err in store test ItemHasBeenProperlyRemovedWhenOutboud3")
	}

}

// make sure a key won't be repeated twice
func Test_storage_do_not_store_the_same_key_twice(t *testing.T) {

	ctx := context.Background()
	obs := observability.Observability{}

	sm := NewShardedMap(2, 2, &obs)
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
