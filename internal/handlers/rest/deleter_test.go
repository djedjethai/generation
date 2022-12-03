package rest

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_delete_should_return_nil_if_value_is_deleted(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	ctx := context.Background()

	mockDeleterSrv.EXPECT().Delete(ctx, "key-a").Return(nil)

	router.HandleFunc("/v1/{key}", handler.keyValueDeleteHandler())

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

	ctx := context.Background()

	mockDeleterSrv.EXPECT().Delete(ctx, "key-a").Return(errors.New("what ever..."))

	router.HandleFunc("/v1/{key}", handler.keyValueDeleteHandler())

	request, _ := http.NewRequest(http.MethodDelete, "/v1/key-a", nil)
	// request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed while checking status code")
	}
}
