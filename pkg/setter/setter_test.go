package setter

import (
	"context"
	"github.com/djedjethai/generation/pkg/observability"
	"github.com/djedjethai/generation/pkg/storage"
	"testing"
)

var setterMocked setter

func setup() {
	obs := observability.Observability{}

	storage := storage.NewMockedShardedMap(1, 0)

	setterMocked = setter{storage, obs}

}

func Test_set_return_no_error_when_all_is_good(t *testing.T) {

	setup()

	ctx := context.Background()

	// create a moked service
	err := setterMocked.Set(ctx, "key", []byte("value"))

	if err != nil {
		t.Error("test setter.Set() should not return an err when all is good")
	}
}

func Test_set_return_error_when_a_problem_occur(t *testing.T) {

	setup()

	ctx := context.Background()

	// create a moked service
	err := setterMocked.Set(ctx, "error", []byte("value"))

	if err == nil {
		t.Error("test setter.Set() should return an err if the storage return an err")
	}
}
