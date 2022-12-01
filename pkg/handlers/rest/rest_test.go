package rest

import (
	"testing"

	"github.com/djedjethai/generation/pkg/config"
	"github.com/djedjethai/generation/pkg/logger"
	// "github.com/djedjethai/generation/pkg/mocks/config"
	"github.com/djedjethai/generation/pkg/mocks/deleter"
	"github.com/djedjethai/generation/pkg/mocks/getter"
	"github.com/djedjethai/generation/pkg/mocks/setter"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

var router *mux.Router

var mockGetterSrv *getter.MockGetter
var mockSetterSrv *setter.MockSetter
var mockDeleterSrv *deleter.MockDeleter

var handler *Handler

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockGetterSrv = getter.NewMockGetter(ctrl)
	mockSetterSrv = setter.NewMockSetter(ctrl)
	mockDeleterSrv = deleter.NewMockDeleter(ctrl)

	lf, _ := logger.NewLoggerFacade(mockSetterSrv, mockDeleterSrv, false, config.PostgresDBParams{})

	svc := config.Services{mockSetterSrv, mockGetterSrv, mockDeleterSrv}

	handler = NewHandler(&svc, lf)

	// router = mux.NewRouter()
	router = handler.Multiplex().(*mux.Router)

	return func() {
		router = nil
		defer ctrl.Finish()
	}
}
