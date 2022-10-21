package deleter

import (
	"github.com/djedjethai/generation0/pkg/storage"
)

//go:generate mockgen -destination=../mocks/deleter/mockDeleter.go -package=deleter github.com/djedjethai/generation0/pkg/deleter Deleter
type Deleter interface {
	Delete(string) error
}

type deleter struct {
	st storage.StorageRepo
}

func NewDeleter(s storage.ShardedMap) Deleter {
	return &deleter{st: s}
}

func (s *deleter) Delete(key string) error {
	err := s.st.Delete(key)
	if err != nil {
		return err
	}
	return nil
}
