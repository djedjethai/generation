package storage

import (
	"context"
	"errors"
	"github.com/djedjethai/generation/internal/models"
)

type mShardedMap struct {
	shd []*Shard
}

func NewMockedShardedMap(nShard, maxLgt int) mShardedMap {
	shards := make([]*Shard, nShard)

	return mShardedMap{shards}
}

func (ms mShardedMap) Set(ctx context.Context, key string, value interface{}) error {
	if key == "error" {
		return errors.New("an errr...")
	}
	return nil
}

func (ms mShardedMap) Get(ctx context.Context, key string) (interface{}, error) {
	if key == "key" {
		return "value", nil
	}
	return nil, errors.New("an error")
}

func (ms mShardedMap) Delete(ctx context.Context, key string, shard *Shard) error {
	if key == "err" {
		return errors.New("an err")
	}
	return nil
}

func (ms mShardedMap) Keys(ctx context.Context) []string {
	var str []string

	return str
}

func (ms mShardedMap) KeysValues(ctx context.Context, kv chan models.KeysValues) error {

	kv <- models.KeysValues{"key1", "value1"}
	kv <- models.KeysValues{"key2", "value2"}
	kv <- models.KeysValues{"key3", "value3"}

	close(kv)

	return nil
}
