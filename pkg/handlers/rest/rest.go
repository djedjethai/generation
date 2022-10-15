package rest

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/djedjethai/generation0/pkg/deleter"
	"github.com/djedjethai/generation0/pkg/getter"
	"github.com/djedjethai/generation0/pkg/logger"
	"github.com/djedjethai/generation0/pkg/setter"
	"github.com/gorilla/mux"
)

var ErrorNoSuchKey = errors.New("no such key")

func Handler(setSrv setter.Setter, getSrv getter.Getter, delSrv deleter.Deleter, loggerFacade *logger.LoggerFacade) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", helloMuxHandler())
	r.HandleFunc("/v1/{key}", keyValuePutHandler(setSrv, loggerFacade)).Methods("PUT")
	r.HandleFunc("/v1/{key}", keyValueGetHandler(getSrv)).Methods("GET")
	r.HandleFunc("/v1/{key}", keyValueDeleteHandler(delSrv, loggerFacade)).Methods("DELETE")
	r.HandleFunc("/v1/util/keys", keyValueGetKeysHandler(getSrv)).Methods("GET")

	return r
}

func keyValueDeleteHandler(delSrv deleter.Deleter, loggerFacade *logger.LoggerFacade) func(w http.ResponseWriter, r *http.Request) {
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
		loggerFacade.WriteDelete(key)

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

		var result string
		switch value.(type) {
		case string:
			result = value.(string)
		case int64:
			result = fmt.Sprintf("%v", value.(int64))
		case float32:
			result = fmt.Sprintf("%v", value.(float32))
		default:
			result = "Invalid type"
		}
		w.Write([]byte(result))
	}
}

func keyValuePutHandler(setSrv setter.Setter, loggerFacade *logger.LoggerFacade) func(w http.ResponseWriter, r *http.Request) {
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
		loggerFacade.WritePut(key, string(value))

		w.WriteHeader(http.StatusCreated)
	}
}

func keyValueGetKeysHandler(getSrv getter.Getter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		value := getSrv.GetKeys()
		stringByte := strings.Join(value, ",")

		w.Write([]byte(stringByte))
	}
}

func helloMuxHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello gorilla/mux!\n"))
	}
}
