package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/djedjethai/generation0/pkg/mocks/getter"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

var router *mux.Router
var mockGetterSrv *getter.MockGetter

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockGetterSrv = getter.NewMockGetter(ctrl)

	router = mux.NewRouter()

	router.HandleFunc("/v1/{key}", keyValueGetHandler(mockGetterSrv))

	return func() {
		router = nil
		defer ctrl.Finish()
	}
}

func Test_getter_should_return_a_value_from_key(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	// arrange
	mockedGetterResponse := "value"
	mockGetterSrv.EXPECT().Get("key-a").Return(mockedGetterResponse, nil)

	request, _ := http.NewRequest(http.MethodGet, "/v1/key-a", nil)

	// act
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
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// assert
	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed while testing the status code")
	}
}
