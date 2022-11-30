package rest

import (
	"errors"
	"net/http"

	"github.com/djedjethai/generation/pkg/deleter"
	"github.com/djedjethai/generation/pkg/getter"
	"github.com/djedjethai/generation/pkg/logger"
	"github.com/djedjethai/generation/pkg/setter"
	"github.com/gorilla/mux"
)

var ErrorNoSuchKey = errors.New("no such key")

func Handler(setSrv setter.Setter, getSrv getter.Getter, delSrv deleter.Deleter, loggerFacade *logger.LoggerFacade) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", helloMuxHandler())
	r.HandleFunc("/v1/{key}", keyValueSetHandler(setSrv, loggerFacade)).Methods("PUT")
	r.HandleFunc("/v1/{key}", keyValueGetHandler(getSrv)).Methods("GET")
	r.HandleFunc("/v1/{key}", keyValueDeleteHandler(delSrv, loggerFacade)).Methods("DELETE")
	r.HandleFunc("/v1/util/keys", keyValueGetKeysHandler(getSrv)).Methods("GET")

	return r
}

func helloMuxHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello gorilla/mux!\n"))
	}
}
