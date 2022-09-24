package storage

import (
	"testing"
)

func TestPut(t *testing.T) {
	_ = storageT.Set("test", "put")

	if storeT["test"] != "put" {
		t.Error("err in store Put() failed")
	}
}

func TestGet(t *testing.T) {
	dt, _ := storageT.Get("test")

	if dt != "put" {
		t.Error("err in store Get() failed")
	}
}

func TestDelete(t *testing.T) {
	_ = storageT.Delete("test")

	if len(storeT) > 0 {
		t.Error("err in store Delete() failed")
	}
}
