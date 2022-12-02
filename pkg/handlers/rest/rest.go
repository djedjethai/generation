package rest

import (
	"errors"
	"github.com/djedjethai/generation/pkg/config"
	"github.com/djedjethai/generation/pkg/logger"
	"github.com/gorilla/mux"
	"net/http"
)

var ErrorNoSuchKey = errors.New("no such key")

type Handler struct {
	services     *config.Services
	loggerFacade *logger.LoggerFacade
}

func NewHandler(svc *config.Services, lf *logger.LoggerFacade) *Handler {
	return &Handler{
		services:     svc,
		loggerFacade: lf,
	}
}

func (h *Handler) Multiplex() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", h.helloMuxHandler())
	r.HandleFunc("/v1/{key}", h.keyValueSetHandler()).Methods("PUT")
	r.HandleFunc("/v1/{key}", h.keyValueGetHandler()).Methods("GET")
	r.HandleFunc("/v1/{key}", h.keyValueDeleteHandler()).Methods("DELETE")
	r.HandleFunc("/v1/util/keys", h.keyValueGetKeysHandler()).Methods("GET")

	return r
}

func (h *Handler) helloMuxHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello gorilla/mux!\n"))
	}
}
