package rest

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_getter_should_return_a_value_from_key(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	ctx := context.Background()

	// arrange
	mockedGetterResponse := "value"
	mockGetterSrv.EXPECT().Get(ctx, "key-a").Return(mockedGetterResponse, nil)

	request, _ := http.NewRequest(http.MethodGet, "/v1/key-a", nil)

	// act
	router.HandleFunc("/v1/{key}", handler.keyValueGetHandler())
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// assert
	if recorder.Code != http.StatusOK {
		t.Error("Failed while testing the status code")
	}

	if recorder.Body.String() != "value" {
		t.Error("Failed while testing the value")
	}
}

func Test_getter_should_return_an_error(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	ctx := context.Background()

	// arrange
	mockGetterSrv.EXPECT().Get(ctx, "key-a").Return(nil, errors.New("what ever..."))

	request, _ := http.NewRequest(http.MethodGet, "/v1/key-a", nil)

	// act
	router.HandleFunc("/v1/{key}", handler.keyValueGetHandler())
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// assert
	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed while testing the status code")
	}
}

func Test_getkeys_should_return_the_keys(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	ctx := context.Background()

	mockedGetKeysResponse := []string{"key1, key2"}
	mockGetterSrv.EXPECT().GetKeys(ctx).Return(mockedGetKeysResponse)

	request, _ := http.NewRequest(http.MethodGet, "/util/keys", nil)

	router.HandleFunc("/util/keys", handler.keyValueGetKeysHandler())
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Error("Failed while checking status code")
	}

	if recorder.Body.String() != "key1, key2" {
		t.Error("Failed while testing the value")
	}
}
