package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

// var ErrorNoSuchKey = errors.New("no such key")

func (h *Handler) keyValueDeleteHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

		// create a context
		ctx := context.Background()

		err := h.services.Deleter.Delete(ctx, key)
		if errors.Is(err, ErrorNoSuchKey) {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		h.loggerFacade.WriteDelete(key)

		w.WriteHeader(http.StatusOK)
	}
}
