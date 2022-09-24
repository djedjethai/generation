package rest

import (
	"errors"
	"io"
	"net/http"

	"github.com/djedjethai/generation0/pkg/deleter"
	"github.com/djedjethai/generation0/pkg/getter"
	"github.com/djedjethai/generation0/pkg/logger"
	"github.com/djedjethai/generation0/pkg/setter"
	"github.com/gorilla/mux"
)

var ErrorNoSuchKey = errors.New("no such key")

func Handler(setSrv setter.Setter, getSrv getter.Getter, delSrv deleter.Deleter, logger logger.TransactionLogger, dbLogger logger.TransactionLogger) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", helloMuxHandler())
	r.HandleFunc("/v1/{key}", keyValuePutHandler(setSrv, logger, dbLogger)).Methods("PUT")
	r.HandleFunc("/v1/{key}", keyValueGetHandler(getSrv)).Methods("GET")
	r.HandleFunc("/v1/{key}", keyValueDeleteHandler(delSrv, logger, dbLogger)).Methods("DELETE")

	return r
}

func keyValueDeleteHandler(delSrv deleter.Deleter, logger logger.TransactionLogger, dbLogger logger.TransactionLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		err := delSrv.Delete(key)
		if errors.Is(err, ErrorNoSuchKey) {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logger.WriteDelete(key)
		dbLogger.WriteDelete(key)

		w.WriteHeader(http.StatusOK)
	}
}

func keyValueGetHandler(getSrv getter.Getter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		value, err := getSrv.Get(key)
		if errors.Is(err, ErrorNoSuchKey) {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write([]byte(value))
	}
}

func keyValuePutHandler(setSrv setter.Setter, logger logger.TransactionLogger, dbLogger logger.TransactionLogger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		value, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			http.Error(w,
				err.Error(),
				http.StatusInternalServerError)
			return
		}

		err = setSrv.Set(key, value)

		if err != nil {
			http.Error(w,
				err.Error(),
				http.StatusInternalServerError)
			return
		}
		logger.WritePut(key, string(value))
		dbLogger.WritePut(key, string(value))

		w.WriteHeader(http.StatusCreated)
	}
}

func helloMuxHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello gorilla/mux!\n"))
	}
}
