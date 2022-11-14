package getter

import (
	"context"

	"github.com/djedjethai/generation0/pkg/config"
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
	trc config.Tracer
}

func NewGetter(s storage.ShardedMap, requests *metric.Int64Counter, tracer config.Tracer) Getter {
	return &getter{
		st:  s,
		req: requests,
		trc: tracer,
	}
}

func (s *getter) Get(ctx context.Context, key string) (interface{}, error) {
	ctx, sp := s.trc.Start(context.Background(), "GetterGet")
	defer sp.End()

	lb := label.Key("getter").String("get")
	s.req.Add(ctx, 1, lb)

	value, err := s.st.Get(ctx, key)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (s *getter) GetKeys(ctx context.Context) []string {
	ctx, sp := s.trc.Start(context.Background(), "GetterGetkeys")
	defer sp.End()

	lb := label.Key("getter").String("getkeys")
	s.req.Add(ctx, 1, lb)

	var keys []string
	keys = s.st.Keys(ctx)

	return keys
}
