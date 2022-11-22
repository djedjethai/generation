package storage

import (
	"context"
	"errors"
	// "fmt"
	"github.com/djedjethai/generation0/pkg/observability"
	"sync"
)

var ErrorNoSuchKey = errors.New("no such key")

type StorageRepo interface {
	Set(context.Context, string, interface{}) error
	Get(context.Context, string) (interface{}, error)
	Keys(context.Context) []string
	Delete(context.Context, string, *Shard) error
}

type Shard struct {
	sync.RWMutex
	m   map[string]*node
	dll dll
}

// TODO idea: improvement: encode key ??
// TODO idea: the key are saved into the node(and in the shardedMap), if remove key from node
// how to know which key has been removed when last element is poped out

// type ShardedMap []*Shard
type ShardedMap struct {
	shd []*Shard
	obs observability.Observability
}

func NewShardedMap(nShard, maxLgt int, observ observability.Observability) ShardedMap {
	shards := make([]*Shard, nShard)

	for i := 0; i < nShard; i++ {
		shard := make(map[string]*node)
		shards[i] = &Shard{
			m:   shard,
			dll: NewDll(maxLgt),
		}
	}

	return ShardedMap{shards, observ}
}

func (m ShardedMap) getShardIndex(key string) int {
	return calculeIndex(key, len(m.shd))
}

// retrieve the shard where the key is stored
func (m ShardedMap) getShard(key string) *Shard {
	index := m.getShardIndex(key)
	// fmt.Println("shard Index: ", index)
	return m.shd[index]
}

func (m ShardedMap) Set(ctx context.Context, key string, value interface{}) error {

	teardown := m.obs.CarryOnTrace(ctx, "StorageSet")
	defer teardown()

	shard := m.getShard(key)

	// if key already exist, remove it first
	shard.RLock()
	_, ok := shard.m[key]
	shard.RUnlock()

	if ok {
		m.obs.Logger.Debug("ShardedMap.Set()", "delete existing key")
		m.Delete(ctx, key, shard)
	}

	shard.Lock()
	defer shard.Unlock()

	newN, outN, err := shard.dll.unshift(key, value)
	if err != nil {
		m.obs.Logger.Error("ShardedMap.Set() failed", err)
		return err
	}
	if outN != nil {
		m.obs.Logger.Debug("ShardedMap.Set()", "delete existing expired queue element")
		// delete the poped node from the shard record
		delete(shard.m, outN.key)
	}

	shard.m[key] = newN

	return nil
}

// TODO get per type, check with gRPC if that works...
func (m ShardedMap) Get(ctx context.Context, key string) (interface{}, error) {

	teardown := m.obs.CarryOnTrace(ctx, "StorageGet")
	defer teardown()

	shard := m.getShard(key)

	shard.RLock()
	nd, ok := shard.m[key]
	shard.RUnlock()
	if !ok {
		return "", ErrorNoSuchKey
	}

	shard.Lock()
	defer shard.Unlock()

	ndExist := shard.dll.removeNode(nd)
	if ndExist != nil {
		m.obs.Logger.Debug("ShardedMap.Get()", "unshift shifted node")
		_, _ = shard.dll.unshiftNode(ndExist)
	}

	if nd.val != "" {
		return nd.val, nil
	} else if nd.valInt != 0 {
		return nd.valInt, nil
	} else if nd.valFloat != 0 {
		return nd.valFloat, nil
	} else {
		return nil, nil
	}
}

func (m ShardedMap) Delete(ctx context.Context, key string, sh *Shard) error {

	teardown := m.obs.CarryOnTrace(ctx, "StorageDelete")
	defer teardown()

	// in the case delete a poped node, we already have the *Shard
	var shard *Shard
	if sh != nil {
		shard = sh
	} else {
		shard = m.getShard(key)
	}

	shard.RLock()
	nd, ok := shard.m[key]
	shard.RUnlock()

	shard.Lock()
	defer shard.Unlock()

	if ok {
		m.obs.Logger.Debug("ShardedMap.Delete()", "delete node")
		_ = shard.dll.removeNode(nd)
		delete(shard.m, key)
	}

	return nil
}

// establish lock(concurrently) on all the table to get all the keys
func (m ShardedMap) Keys(ctx context.Context) []string {

	teardown := m.obs.CarryOnTrace(ctx, "StorageKeys")
	defer teardown()

	keys := make([]string, 0) // Create an empty keys slice

	mutex := sync.Mutex{} // Mutex for write safety to keys

	wg := sync.WaitGroup{}
	wg.Add(len(m.shd))

	for _, shard := range m.shd { // Run a goroutine for each slice
		go func(s *Shard) {
			s.RLock()

			for key := range s.m {
				mutex.Lock()
				keys = append(keys, key)
				mutex.Unlock()
			}

			s.RUnlock()
			wg.Done()

		}(shard)
	}

	wg.Wait()

	return keys
}
