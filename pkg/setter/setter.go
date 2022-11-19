package setter

import (
	"context"
	"errors"
	"github.com/djedjethai/generation0/pkg/config"
	"github.com/djedjethai/generation0/pkg/storage"
	"go.opentelemetry.io/otel/label"
	// "go.uber.org/zap"
)

//go:generate mockgen -destination=../mocks/setter/mockSetter.go -package=setter github.com/djedjethai/generation0/pkg/setter Setter
type Setter interface {
	Set(context.Context, string, []byte) error
}

type setter struct {
	st  storage.StorageRepo
	obs config.Observability
	// st storage.ShardedMap
}

func NewSetter(s storage.ShardedMap, observ config.Observability) Setter {
	lb := label.Key("setter").String("set")
	observ.Labels = append(observ.Labels, lb)
	return &setter{
		st:  s,
		obs: observ,
	}
}

func (s *setter) Set(ctx context.Context, key string, value []byte) error {

	// zap.S().Errorw(
	// 	"Testing zap, in Set",
	// 	"in set",
	// 	"setter",
	// )
	// s.obs.Logger.Debug("test debug", "setter", "set")
	s.obs.Logger.Error("test alert", errors.New("my error"))
	s.obs.Logger.Warning("test alert", "domainouuuooo")
	s.obs.Logger.Info("test alert", "domainouuuooo")
	s.obs.Logger.Debug("test alert", "domainouuuooo")
	// s.obs.Logger.Alert("test alert", "setter", "Set", "test errrrrrr")

	if s.obs.IsTracing {
		ctx1, sp := s.obs.Tracer.Start(context.Background(), "SetterSet")
		defer sp.End()

		ctx = ctx1
	}

	if s.obs.IsMetrics {
		s.obs.Requests.Add(ctx, 1, s.obs.Labels...)
	}

	err := s.st.Set(ctx, key, string(value))
	if err != nil {
		return err
	}
	return nil

}
