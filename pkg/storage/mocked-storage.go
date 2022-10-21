package storage

import (
	"errors"
)

type mShardedMap []*Shard

func NewMockedShardedMap(nShard, maxLgt int) mShardedMap {
	shards := make([]*Shard, nShard)

	return shards
}

func (ms mShardedMap) Set(key string, value interface{}) error {
	if key == "error" {
		return errors.New("an errr...")
	}
	return nil
}

func (ms mShardedMap) Get(key string) (interface{}, error) {
	if key == "key" {
		return "value", nil
	}
	return nil, errors.New("an error")
}

func (ms mShardedMap) Delete(key string) error {
	if key == "err" {
		return errors.New("an err")
	}
	return nil
}

func (ms mShardedMap) Keys() []string {
	var str []string

	return str
}
