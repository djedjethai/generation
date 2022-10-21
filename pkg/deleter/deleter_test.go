package deleter

import (
	"github.com/djedjethai/generation0/pkg/storage"
	"testing"
)

var deleterMocked deleter

func setup() {
	storage := storage.NewMockedShardedMap(1, 0)

	deleterMocked = deleter{storage}
}

func Test_delete_return_nil_if_no_err(t *testing.T) {

	setup()

	err := deleterMocked.Delete("key")

	if err != nil {
		t.Error("test deleter.Delete() should return nil")
	}

}

func Test_delete_return_error_when_a_problem_occur(t *testing.T) {

	setup()

	err := deleterMocked.Delete("err")

	if err == nil {
		t.Error("test deleter.Delete() should return an err if the storage return an err")
	}
}
