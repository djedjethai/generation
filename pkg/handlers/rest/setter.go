package rest

import (
	"context"
	"io"
	"net/http"

	"github.com/djedjethai/generation0/pkg/logger"
	"github.com/djedjethai/generation0/pkg/setter"
	"github.com/gorilla/mux"
)

func keyValuePutHandler(setSrv setter.Setter, loggerFacade *logger.LoggerFacade) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		value, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		// create a context
		ctx := context.Background()

		if err != nil {
			http.Error(w,
				err.Error(),
				http.StatusInternalServerError)
			return
		}

		err = setSrv.Set(ctx, key, value)

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
