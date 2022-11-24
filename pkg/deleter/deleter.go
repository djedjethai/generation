package deleter

import (
	"context"
	"github.com/djedjethai/generation0/pkg/observability"
	"github.com/djedjethai/generation0/pkg/storage"
	"go.opentelemetry.io/otel/label"
)

//go:generate mockgen -destination=../mocks/deleter/mockDeleter.go -package=deleter github.com/djedjethai/generation0/pkg/deleter Deleter
type Deleter interface {
	Delete(context.Context, string) error
}

type deleter struct {
	st  storage.StorageRepo
	obs observability.Observability
}

func NewDeleter(s storage.ShardedMap, observ observability.Observability) Deleter {
	// run the query: golru_requests_total{deleter="delete"}
	lb := label.Key("deleter").String("delete")
	observ.Labels = append(observ.Labels, lb)
	return &deleter{
		st:  s,
		obs: observ,
	}
}

func (s *deleter) Delete(ctx context.Context, key string) error {

	s.obs.Logger.Debug("Deleter/Delete()", "hit func")

	ctx, teardown := s.obs.StartTrace(ctx, "DeleterDelete")
	defer teardown()

	s.obs.AddMetrics(ctx)

	err := s.st.Delete(ctx, key, nil)
	if err != nil {
		s.obs.Logger.Error("Deleter/Delete() failed", err)
		return err
	}

	s.obs.Logger.Debug("Deleter/Delete()", "executed successfully")
	return nil
}
