package getter

import (
	"context"

	"github.com/djedjethai/generation0/pkg/storage"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
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
	req *metric.Int64Counter
}

func NewGetter(s storage.ShardedMap, requests *metric.Int64Counter) Getter {
	return &getter{
		st:  s,
		req: requests,
	}
}

func (s *getter) Get(ctx context.Context, key string) (interface{}, error) {

	lb := label.Key("getter").String("get")
	s.req.Add(ctx, 1, lb)

	value, err := s.st.Get(key)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (s *getter) GetKeys(ctx context.Context) []string {

	lb := label.Key("getter").String("getkeys")
	s.req.Add(ctx, 1, lb)

	var keys []string
	keys = s.st.Keys()

	return keys
}
