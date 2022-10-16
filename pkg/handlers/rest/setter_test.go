package rest

import (
	"errors"
	"strings"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/djedjethai/generation0/pkg/config"
	dele "github.com/djedjethai/generation0/pkg/deleter"
	"github.com/djedjethai/generation0/pkg/logger"
	sett "github.com/djedjethai/generation0/pkg/setter"
	"github.com/djedjethai/generation0/pkg/storage"
)

func Test_put_should_return_nil_if_value_is_added(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	st := storage.NewShardedMap(1, 3)
	ss := sett.NewSetter(st)
	ds := dele.NewDeleter(st)

	lf, _ := logger.NewLoggerFacade(ss, ds, false, false, config.PostgresDBParams{}, "encryptk")

	mockSetterSrv.EXPECT().Set("key-a", []uint8{118, 97, 108, 117, 101, 45, 97}).Return(nil)

	router.HandleFunc("/v1/{key}", keyValuePutHandler(mockSetterSrv, lf))

	request, _ := http.NewRequest(http.MethodPut, "/v1/key-a", strings.NewReader("value-a"))
	// request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Error("Failed while checking status code")
	}

	if recorder.Body.String() != "" {
		t.Error("Failed while testing the value")
	}
}

func Test_put_should_return_err_service_return_err(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	st := storage.NewShardedMap(1, 3)
	ss := sett.NewSetter(st)
	ds := dele.NewDeleter(st)

	lf, _ := logger.NewLoggerFacade(ss, ds, false, false, config.PostgresDBParams{}, "encryptk")

	mockSetterSrv.EXPECT().Set("key-a", []uint8{118, 97, 108, 117, 101, 45, 97}).Return(errors.New("what ever..."))

	router.HandleFunc("/v1/{key}", keyValuePutHandler(mockSetterSrv, lf))

	request, _ := http.NewRequest(http.MethodPut, "/v1/key-a", strings.NewReader("value-a"))
	// request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed while checking status code")
	}
}
