package getter

import (
	"github.com/djedjethai/generation0/pkg/storage"
)

type Getter interface {
	Get(string) (string, error)
}

type getter struct {
	st storage.StorageRepo
}

func NewGetter(s storage.StorageRepo) Getter {
	return &getter{st: s}
}

func (s *getter) Get(key string) (string, error) {
	value, err := s.st.Get(key)
	if err != nil {
		return "", err
	}
	return value, nil
}
