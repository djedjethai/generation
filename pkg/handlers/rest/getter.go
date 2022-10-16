package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/djedjethai/generation0/pkg/getter"
	"github.com/gorilla/mux"
)

// var ErrorNoSuchKey = errors.New("no such key")

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

func keyValueGetKeysHandler(getSrv getter.Getter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		value := getSrv.GetKeys()
		stringByte := strings.Join(value, ",")

		w.Write([]byte(stringByte))
	}
}
