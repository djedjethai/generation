package deleter

import (
	"context"

	"github.com/djedjethai/generation0/pkg/config"
	"github.com/djedjethai/generation0/pkg/storage"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
)

//go:generate mockgen -destination=../mocks/deleter/mockDeleter.go -package=deleter github.com/djedjethai/generation0/pkg/deleter Deleter
type Deleter interface {
	Delete(context.Context, string) error
}

type deleter struct {
	st  storage.StorageRepo
	lbl []label.KeyValue
	req *metric.Int64Counter
	trc config.Tracer
}

func NewDeleter(s storage.ShardedMap, labels []label.KeyValue, request *metric.Int64Counter, tracer config.Tracer) Deleter {
	// run the query: golru_requests_total{deleter="delete"}
	lb := label.Key("deleter").String("delete")
	labels = append(labels, lb)
	return &deleter{
		st:  s,
		lbl: labels,
		req: request,
		trc: tracer,
	}
}

func (s *deleter) Delete(ctx context.Context, key string) error {
	ctx, sp := s.trc.Start(context.Background(), "DeleterDelete")
	defer sp.End()

	s.req.Add(ctx, 1, s.lbl...)

	err := s.st.Delete(ctx, key)
	if err != nil {
		return err
	}
	return nil
}
