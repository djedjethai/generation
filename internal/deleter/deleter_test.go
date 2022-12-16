package deleter

import (
	"context"
	"github.com/djedjethai/generation/internal/observability"
	"github.com/djedjethai/generation/internal/storage"
	"testing"
)

var deleterMocked deleter

func setup() {
	storage := storage.NewMockedShardedMap(1, 0)

	obs := observability.Observability{}

	deleterMocked = deleter{storage, &obs}
}

func Test_delete_return_nil_if_no_err(t *testing.T) {

	setup()

	ctx := context.Background()

	err := deleterMocked.Delete(ctx, "key")

	if err != nil {
		t.Error("test deleter.Delete() should return nil")
	}

}

func Test_delete_return_error_when_a_problem_occur(t *testing.T) {

	setup()

	ctx := context.Background()

	err := deleterMocked.Delete(ctx, "err")

	if err == nil {
		t.Error("test deleter.Delete() should return an err if the storage return an err")
	}
}
