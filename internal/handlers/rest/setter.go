package rest

import (
	"context"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func (h *Handler) keyValueSetHandler() func(w http.ResponseWriter, r *http.Request) {
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

		err = h.services.Setter.Set(ctx, key, value)

		if err != nil {
			http.Error(w,
				err.Error(),
				http.StatusInternalServerError)
			return
		}
		h.loggerFacade.WriteSet(key, string(value))

		w.WriteHeader(http.StatusCreated)
	}
}
