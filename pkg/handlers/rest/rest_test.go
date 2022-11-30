package rest

import (
	"github.com/djedjethai/generation/pkg/mocks/deleter"
	"github.com/djedjethai/generation/pkg/mocks/getter"
	"github.com/djedjethai/generation/pkg/mocks/setter"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"testing"
)

var router *mux.Router
var mockGetterSrv *getter.MockGetter
var mockSetterSrv *setter.MockSetter
var mockDeleterSrv *deleter.MockDeleter

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockGetterSrv = getter.NewMockGetter(ctrl)
	mockSetterSrv = setter.NewMockSetter(ctrl)
	mockDeleterSrv = deleter.NewMockDeleter(ctrl)

	router = mux.NewRouter()

	return func() {
		router = nil
		defer ctrl.Finish()
	}
}
