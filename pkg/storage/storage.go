package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"sync"
)

var ErrorNoSuchKey = errors.New("no such key")

type StorageRepo interface {
	Set(string, string) error
	Get(string) (string, error)
	Delete(string) error
}

type Shard struct {
	sync.RWMutex
	m   map[string]*node
	dll dll
}

type ShardedMap []*Shard

func NewShardedMap(nShard, maxLgt int) ShardedMap {
	shards := make([]*Shard, nShard)

	for i := 0; i < nShard; i++ {
		shard := make(map[string]*node)
		shards[i] = &Shard{
			m:   shard,
			dll: NewDll(maxLgt),
		}
	}

	return shards
}

func (m ShardedMap) getShardIndex(key string) int {
	checksum := sha1.Sum([]byte(key)) // Use Sum from "crypto/sha1"
	hash := int(checksum[17])         // Pick an arbitrary byte as the hash
	fmt.Println("see the index(modulo): ", hash%len(m))
	return hash % len(m)
}

// retrieve the shard where the key is stored
func (m ShardedMap) getShard(key string) *Shard {
	index := m.getShardIndex(key)
	return m[index]
}

func (m ShardedMap) Set(key string, value string) error {
	shard := m.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	newN, outN := shard.dll.unshift(key, value)
	if outN != nil {
		delete(shard.m, outN.key)
	}

	shard.m[key] = newN

	return nil
}

func (m ShardedMap) Get(key string) (string, error) {
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
		shard.dll.unshiftNode(ndExist)
	}

	return nd.val, nil
}

func (m ShardedMap) Delete(key string) error {
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
func (m ShardedMap) Keys() []string {
	keys := make([]string, 0) // Create an empty keys slice

	mutex := sync.Mutex{} // Mutex for write safety to keys

	wg := sync.WaitGroup{}
	wg.Add(len(m))

	for _, shard := range m { // Run a goroutine for each slice
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

// type storage struct {
// 	sync.RWMutex
// 	store map[string]*node
// 	dll   dll
// }
//
// func NewStorage(maxLgt int) StorageRepo {
// 	store := make(map[string]*node)
// 	return &storage{
// 		store: store,
// 		dll:   NewDll(maxLgt),
// 	}
// }
//
// func (s *storage) Set(key string, value string) error {
// 	// create node in dll
// 	s.Lock()
// 	newN, outN := s.dll.unshift(key, value)
// 	if outN != nil {
// 		// in case dll poped out the last item
// 		delete(s.store, outN.key)
// 	}
//
// 	// add node to map
// 	s.store[key] = newN
// 	s.Unlock()
//
// 	return nil
// }
//
// func (s *storage) Get(key string) (string, error) {
//
// 	s.RLock()
// 	nd, ok := s.store[key]
// 	s.RUnlock()
// 	if !ok {
// 		return "", ErrorNoSuchKey
// 	}
//
// 	// move the Get node(so nd) to head of dll
// 	s.Lock()
// 	ndExist := s.dll.removeNode(nd)
// 	if ndExist != nil {
// 		// re-unshift ndExist or nd(what ever they point to the same location)
// 		s.dll.unshiftNode(ndExist)
// 	}
// 	s.Unlock()
//
// 	return nd.val, nil
// }
//
// func (s *storage) Delete(key string) error {
//
// 	s.RLock()
// 	nd, ok := s.store[key]
// 	s.RUnlock()
//
// 	if ok {
// 		s.Lock()
// 		_ = s.dll.removeNode(nd)
//
// 		delete(s.store, key)
// 		s.Unlock()
// 	}
//
// 	return nil
// }
