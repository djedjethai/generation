package getter

import (
	"context"
	"fmt"

	"github.com/djedjethai/generation0/pkg/observability"
	"github.com/djedjethai/generation0/pkg/storage"
)

// run: go generate ./...
//
//go:generate mockgen -destination=../mocks/getter/mockGetter.go -package=getter github.com/djedjethai/generation0/pkg/getter Getter
type Getter interface {
	Get(context.Context, string) (interface{}, error)
	GetKeys(context.Context) []string
}

type getter struct {
	st  storage.StorageRepo
	obs observability.Observability
}

func NewGetter(s storage.ShardedMap, observ observability.Observability) Getter {
	return &getter{
		st:  s,
		obs: observ,
	}
}

func (s *getter) Get(ctx context.Context, key string) (interface{}, error) {

	s.obs.Logger.Debug("Getter/Get()", "hit func")

	ctx, teardown := s.obs.StartTrace(ctx, "GetterGet")
	defer teardown()

	s.obs.AddMetricsAndSpecificLabel(ctx, "getter", "get")

	value, err := s.st.Get(ctx, key)
	if err != nil {
		s.obs.Logger.Warning("Getter/Get() failed", fmt.Sprintf("%v", err))
		return "", err
	}

	s.obs.Logger.Debug("Getter/Get()", "executed successfully")
	return value, nil
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
