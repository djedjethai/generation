package getter

import (
	"github.com/djedjethai/generation0/pkg/storage"
)

// run: go generate ./...
//
//go:generate mockgen -destination=../mocks/getter/mockGetter.go -package=getter github.com/djedjethai/generation0/pkg/getter Getter
type Getter interface {
	Get(string) (interface{}, error)
	GetKeys() []string
}

type getter struct {
	st storage.StorageRepo
}

func NewGetter(s storage.ShardedMap) Getter {
	return &getter{st: s}
}

func (s *getter) Get(key string) (interface{}, error) {
	value, err := s.st.Get(key)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (s *getter) GetKeys() []string {
	var keys []string
	keys = s.st.Keys()

	return keys
}
