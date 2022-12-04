package getter

import (
	"context"
	"reflect"
	"testing"

	"github.com/djedjethai/generation/internal/observability"
	"github.com/djedjethai/generation/internal/storage"
)

var getterMocked getter

func setup() {

	obs := observability.Observability{}

	str := storage.NewMockedShardedMap(1, 0)

	getterMocked = getter{str, obs}
}

func Test_get_return_a_value_and_no_error(t *testing.T) {

	setup()

	ctx := context.Background()

	// create a moked service
	value, err := getterMocked.Get(ctx, "key")

	if value != "value" {
		t.Error("test getter.Get() should return a value")
	}

	if err != nil {
		t.Error("test getter.Get() should not return an err when all is good")
	}
}

func Test_get_return_error_when_a_problem_occur(t *testing.T) {

	setup()

	ctx := context.Background()

	// create a moked service
	_, err := getterMocked.Get(ctx, "invalidKey")

	if err == nil {
		t.Error("test getter.Get() should return an err if the storage return an err")
	}
}

func Test_keys_return_an_array_of_keys(t *testing.T) {

	setup()

	ctx := context.Background()

	// create a moked service
	keys := getterMocked.GetKeys(ctx)

	rt := reflect.TypeOf(keys)

	if rt.Kind() != reflect.Slice {
		t.Error("test getter.GetKeys() should return slice of strings")
	}

	if rt.Elem().String() != "string" {
		t.Error("test getter.GetKeys() should return strings")
	}
}
