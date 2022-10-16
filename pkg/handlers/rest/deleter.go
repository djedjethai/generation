package rest

import (
	"errors"
	"net/http"

	"github.com/djedjethai/generation0/pkg/deleter"
	"github.com/djedjethai/generation0/pkg/logger"
	"github.com/gorilla/mux"
)

// var ErrorNoSuchKey = errors.New("no such key")

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
