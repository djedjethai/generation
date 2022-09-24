package setter

import (
	"github.com/djedjethai/generation0/pkg/storage"
)

type Setter interface {
	Set(string, []byte) error
}

type setter struct {
	st storage.StorageRepo
}

func NewSetter(s storage.StorageRepo) Setter {
	return &setter{st: s}
}

func (s *setter) Set(key string, value []byte) error {
	err := s.st.Set(key, string(value))
	if err != nil {
		return err
	}
	return nil
}
