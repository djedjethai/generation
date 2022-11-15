package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/djedjethai/generation0/pkg/config"
	"go.opentelemetry.io/otel"
	"sync"
)

var ErrorNoSuchKey = errors.New("no such key")

type StorageRepo interface {
	Set(context.Context, string, interface{}) error
	Get(context.Context, string) (interface{}, error)
	Keys(context.Context) []string
	Delete(context.Context, string) error
}

type Shard struct {
	sync.RWMutex
	m   map[string]*node
	dll dll
}

// type ShardedMap []*Shard
type ShardedMap struct {
	shd []*Shard
	obs config.Observability
}

func NewShardedMap(nShard, maxLgt int, observ config.Observability) ShardedMap {
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
	fmt.Println("shard Index: ", index)
	return m.shd[index]
}

func (m ShardedMap) Set(ctx context.Context, key string, value interface{}) error {
	if m.obs.IsTracing {
		tr := otel.GetTracerProvider().Tracer(m.obs.ServiceName)
		_, sp := tr.Start(ctx, "StorageSet")
		defer sp.End()
	}

	shard := m.getShard(key)

	// if key already exist, remove it first
	shard.RLock()
	_, ok := shard.m[key]
	shard.RUnlock()

	if ok {
		m.Delete(ctx, key)
	}

	shard.Lock()
	defer shard.Unlock()

	newN, outN, err := shard.dll.unshift(key, value)
	if err != nil {
		return err
	}
	if outN != nil {
		delete(shard.m, outN.key)
	}

	shard.m[key] = newN

	return nil
}

// TODO get per type, check with gRPC if that works...
func (m ShardedMap) Get(ctx context.Context, key string) (interface{}, error) {
	if m.obs.IsTracing {
		tr := otel.GetTracerProvider().Tracer(m.obs.ServiceName)
		_, sp := tr.Start(ctx, "StorageGet")
		defer sp.End()
	}

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

func (m ShardedMap) Delete(ctx context.Context, key string) error {
	if m.obs.IsTracing {
		tr := otel.GetTracerProvider().Tracer(m.obs.ServiceName)
		_, sp := tr.Start(ctx, "StorageDelete")
		defer sp.End()
	}

	shard := m.getShard(key)

	shard.RLock()
	nd, ok := shard.m[key]
	shard.RUnlock()

	shard.Lock()
	defer shard.Unlock()

	if ok {
		_ = shard.dll.removeNode(nd)
		delete(shard.m, key)
	}

	return nil
}

// establish lock(concurrently) on all the table to get all the keys
func (m ShardedMap) Keys(ctx context.Context) []string {
	if m.obs.IsTracing {
		tr := otel.GetTracerProvider().Tracer(m.obs.ServiceName)
		_, sp := tr.Start(ctx, "StorageKeys")
		defer sp.End()
	}

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
