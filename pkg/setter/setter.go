package setter

import (
	"context"

	"github.com/djedjethai/generation0/pkg/config"
	"github.com/djedjethai/generation0/pkg/storage"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
)

//go:generate mockgen -destination=../mocks/setter/mockSetter.go -package=setter github.com/djedjethai/generation0/pkg/setter Setter
type Setter interface {
	Set(context.Context, string, []byte) error
}

type setter struct {
	st  storage.StorageRepo
	lbl []label.KeyValue
	req *metric.Int64Counter
	trc config.Tracer
	// st storage.ShardedMap
}

func NewSetter(s storage.ShardedMap, labels []label.KeyValue, requests *metric.Int64Counter, tracer config.Tracer) Setter {
	lb := label.Key("setter").String("set")
	labels = append(labels, lb)
	return &setter{
		st:  s,
		lbl: labels,
		req: requests,
		trc: tracer,
	}
}

func (s *setter) Set(ctx context.Context, key string, value []byte) error {
	ctx, sp := s.trc.Start(context.Background(), "SetterSet")
	defer sp.End()

	s.req.Add(ctx, 1, s.lbl...)

	err := s.st.Set(ctx, key, string(value))
	if err != nil {
		return err
	}
	return nil
}
