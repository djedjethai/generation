package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/djedjethai/generation0/pkg/config"
	dele "github.com/djedjethai/generation0/pkg/deleter"
	"github.com/djedjethai/generation0/pkg/logger"
	sett "github.com/djedjethai/generation0/pkg/setter"
	"github.com/djedjethai/generation0/pkg/storage"
)

func Test_delete_should_return_nil_if_value_is_deleted(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	st := storage.NewShardedMap(1, 3)
	ss := sett.NewSetter(st)
	ds := dele.NewDeleter(st)

	lf, _ := logger.NewLoggerFacade(ss, ds, false, false, config.PostgresDBParams{}, "encryptk")

	mockDeleterSrv.EXPECT().Delete("key-a").Return(nil)

	router.HandleFunc("/v1/{key}", keyValueDeleteHandler(mockDeleterSrv, lf))

	request, _ := http.NewRequest(http.MethodDelete, "/v1/key-a", nil)
	// request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Error("Failed while checking status code")
	}

	if recorder.Body.String() != "" {
		t.Error("Failed while testing the value")
	}
}

func Test_delete_should_return_err_if_delete_service_return_an_err(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	// no storage will return an err
	st := storage.NewShardedMap(1, 3)
	ss := sett.NewSetter(st)
	ds := dele.NewDeleter(st)

	lf, _ := logger.NewLoggerFacade(ss, ds, false, false, config.PostgresDBParams{}, "encryptk")

	mockDeleterSrv.EXPECT().Delete("key-a").Return(errors.New("what ever..."))

	router.HandleFunc("/v1/{key}", keyValueDeleteHandler(mockDeleterSrv, lf))

	request, _ := http.NewRequest(http.MethodDelete, "/v1/key-a", nil)
	// request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed while checking status code")
	}
}
