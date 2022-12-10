package getter

import (
	"context"
	"fmt"

	"github.com/djedjethai/generation/internal/models"
	"github.com/djedjethai/generation/internal/observability"
	"github.com/djedjethai/generation/internal/storage"
)

// run: go generate ./...

//go:generate mockgen -destination=../mocks/getter/mockGetter.go -package=getter github.com/djedjethai/generation/internal/getter Getter
type Getter interface {
	// Get(context.Context, string) (interface{}, error)
	Get(context.Context, string) interface{}
	GetKeys(context.Context) []string
	GetKeysValues(context.Context, chan models.KeysValues) error
}

type getter struct {
	st  storage.StorageRepo
	obs *observability.Observability
}

func NewGetter(s storage.ShardedMap, observ *observability.Observability) Getter {
	return &getter{
		st:  s,
		obs: observ,
	}
}

// func (s *getter) Get(ctx context.Context, key string) (interface{}, error) {
func (s *getter) Get(ctx context.Context, key string) interface{} {

	s.obs.Logger.Debug("Getter/Get()", "hit func")

	ctx, teardown := s.obs.StartTrace(ctx, "GetterGet")
	defer teardown()

	s.obs.AddMetricsAndSpecificLabel(ctx, "getter", "get")

	value, err := s.st.Get(ctx, key)
	if err != nil {
		s.obs.Logger.Warning("Getter/Get() failed", fmt.Sprintf("%v", err))
		return err
	}

	s.obs.Logger.Debug("Getter/Get()", "executed successfully")
	return value
}

func (s *getter) GetKeys(ctx context.Context) []string {

	s.obs.Logger.Debug("Getter/GetKeys()", "hit func")

	ctx, teardown := s.obs.StartTrace(ctx, "GetterGetKeys")
	defer teardown()

	s.obs.AddMetricsAndSpecificLabel(ctx, "getter", "getkeys")

	var keys []string
	keys = s.st.Keys(ctx)

	s.obs.Logger.Debug("Getter/Get()", "executed successfully")

	return keys
}

func (s *getter) GetKeysValues(ctx context.Context, kv chan models.KeysValues) error {
	return s.st.KeysValues(ctx, kv)
}
