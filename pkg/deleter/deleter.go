package deleter

import (
	"context"

	"github.com/djedjethai/generation0/pkg/config"
	"github.com/djedjethai/generation0/pkg/storage"
	"go.opentelemetry.io/otel/label"
	// "go.opentelemetry.io/otel/metric"
)

//go:generate mockgen -destination=../mocks/deleter/mockDeleter.go -package=deleter github.com/djedjethai/generation0/pkg/deleter Deleter
type Deleter interface {
	Delete(context.Context, string) error
}

type deleter struct {
	st  storage.StorageRepo
	obs config.Observability
}

func NewDeleter(s storage.ShardedMap, observ config.Observability) Deleter {
	// run the query: golru_requests_total{deleter="delete"}
	lb := label.Key("deleter").String("delete")
	observ.Labels = append(observ.Labels, lb)
	return &deleter{
		st:  s,
		obs: observ,
	}
}

func (s *deleter) Delete(ctx context.Context, key string) error {
	if s.obs.IsTracing {
		ctx1, sp := s.obs.Tracer.Start(context.Background(), "DeleterDelete")
		defer sp.End()

		ctx = ctx1
	}

	if s.obs.IsMetrics {
		s.obs.Requests.Add(ctx, 1, s.obs.Labels...)
	}

	err := s.st.Delete(ctx, key)
	if err != nil {
		return err
	}
	return nil
}
