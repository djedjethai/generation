package storage

import (
	"errors"
	"sync"
)

var ErrorNoSuchKey = errors.New("no such key")

type StorageRepo interface {
	Set(string, string) error
	Get(string) (string, error)
	Delete(string) error
}

type storage struct {
	sync.RWMutex
	store map[string]*node
	dll   dll
}

func NewStorage(maxLgt int) StorageRepo {
	store := make(map[string]*node)
	return &storage{
		store: store,
		dll:   NewDll(maxLgt),
	}
}

func (s *storage) Set(key string, value string) error {
	// create node in dll
	s.Lock()
	newN, outN := s.dll.unshift(key, value)
	if outN != nil {
		// in case dll poped out the last item
		delete(s.store, outN.key)
	}

	// add node to map
	s.store[key] = newN
	s.Unlock()

	return nil
}

func (s *storage) Get(key string) (string, error) {

	s.RLock()
	nd, ok := s.store[key]
	s.RUnlock()
	if !ok {
		return "", ErrorNoSuchKey
	}

	// move the Get node(so nd) to head of dll
	s.Lock()
	ndExist := s.dll.removeNode(nd)
	if ndExist != nil {
		// re-unshift ndExist or nd(what ever they point to the same location)
		s.dll.unshiftNode(ndExist)
	}
	s.Unlock()

	return nd.val, nil
}

func (s *storage) Delete(key string) error {

	s.RLock()
	nd, ok := s.store[key]
	s.RUnlock()

	if ok {
		s.Lock()
		_ = s.dll.removeNode(nd)

		delete(s.store, key)
		s.Unlock()
	}

	return nil
}
