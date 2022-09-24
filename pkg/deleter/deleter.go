package deleter

import (
	"github.com/djedjethai/generation0/pkg/storage"
)

type Deleter interface {
	Delete(string) error
}

type deleter struct {
	st storage.StorageRepo
}

func NewDeleter(s storage.StorageRepo) Deleter {
	return &deleter{st: s}
}

func (s *deleter) Delete(key string) error {
	err := s.st.Delete(key)
	if err != nil {
		return err
	}
	return nil
}
