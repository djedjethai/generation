package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_getter_should_return_a_value_from_key(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	// arrange
	mockedGetterResponse := "value"
	mockGetterSrv.EXPECT().Get("key-a").Return(mockedGetterResponse, nil)

	request, _ := http.NewRequest(http.MethodGet, "/v1/key-a", nil)

	// act
	router.HandleFunc("/v1/{key}", keyValueGetHandler(mockGetterSrv))
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

	// arrange
	mockGetterSrv.EXPECT().Get("key-a").Return(nil, errors.New("what ever..."))

	request, _ := http.NewRequest(http.MethodGet, "/v1/key-a", nil)

	// act
	router.HandleFunc("/v1/{key}", keyValueGetHandler(mockGetterSrv))
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

	mockedGetKeysResponse := []string{"key1, key2"}
	mockGetterSrv.EXPECT().GetKeys().Return(mockedGetKeysResponse)

	request, _ := http.NewRequest(http.MethodGet, "/util/keys", nil)

	router.HandleFunc("/util/keys", keyValueGetKeysHandler(mockGetterSrv))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Error("Failed while checking status code")
	}

	if recorder.Body.String() != "key1, key2" {
		t.Error("Failed while testing the value")
	}
}
