package setter

import (
	"context"
	"github.com/djedjethai/generation/pkg/observability"
	"github.com/djedjethai/generation/pkg/storage"
	"go.opentelemetry.io/otel/label"
)

//go:generate mockgen -destination=../mocks/setter/mockSetter.go -package=setter github.com/djedjethai/generation/pkg/setter Setter
type Setter interface {
	Set(context.Context, string, []byte) error
}

type setter struct {
	st  storage.StorageRepo
	obs observability.Observability
}

func NewSetter(s storage.ShardedMap, observ observability.Observability) Setter {
	lb := label.Key("setter").String("set")
	observ.Labels = append(observ.Labels, lb)
	return &setter{
		st:  s,
		obs: observ,
	}
}

func (s *setter) Set(ctx context.Context, key string, value []byte) error {

	s.obs.Logger.Debug("Setter/Set()", "hit func")

	// exemple of logs types
	// s.obs.Logger.Error("test alert", errors.New("my error"))
	// s.obs.Logger.Warning("test alert", "domainouuuooo")
	// s.obs.Logger.Info("test alert", "domainouuuooo")
	// s.obs.Logger.Alert("test alert", "setter", "Set", "test errrrrrr")

	ctx, teardown := s.obs.StartTrace(ctx, "SetterSet")
	defer teardown()

	s.obs.AddMetrics(ctx)

	err := s.st.Set(ctx, key, string(value))
	if err != nil {
		s.obs.Logger.Error("Setter/Set() failed", err)
		return err
	}

	s.obs.Logger.Debug("Setter/Set()", "executed successfully")
	return nil

}
