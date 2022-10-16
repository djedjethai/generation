package setter

import (
	"github.com/djedjethai/generation0/pkg/storage"
)

//go:generate mockgen -destination=../mocks/setter/mockSetter.go -package=setter github.com/djedjethai/generation0/pkg/setter Setter
type Setter interface {
	Set(string, []byte) error
}

type setter struct {
	// st storage.StorageRepo
	st storage.ShardedMap
}

func NewSetter(s storage.ShardedMap) Setter {
	return &setter{st: s}
}

func (s *setter) Set(key string, value []byte) error {
	err := s.st.Set(key, string(value))
	if err != nil {
		return err
	}
	return nil
}
