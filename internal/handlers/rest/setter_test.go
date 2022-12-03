package rest

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_put_should_return_nil_if_value_is_added(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	ctx := context.Background()

	mockSetterSrv.EXPECT().Set(ctx, "key-a", []uint8{118, 97, 108, 117, 101, 45, 97}).Return(nil)

	router.HandleFunc("/v1/{key}", handler.keyValueSetHandler())

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

	ctx := context.Background()

	mockSetterSrv.EXPECT().Set(ctx, "key-a", []uint8{118, 97, 108, 117, 101, 45, 97}).Return(errors.New("what ever..."))

	router.HandleFunc("/v1/{key}", handler.keyValueSetHandler())

	request, _ := http.NewRequest(http.MethodPut, "/v1/key-a", strings.NewReader("value-a"))
	// request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed while checking status code")
	}
}
