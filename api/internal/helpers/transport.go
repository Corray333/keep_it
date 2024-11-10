package helpers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func SendError(w http.ResponseWriter, errI error) {
	err, ok := errI.(*CustomError)
	if !ok {
		http.Error(w, errI.Error(), http.StatusInternalServerError)
	} else {
		http.Error(w, err.Error(), err.Code)
	}
}

func GetIntQueryParam(r *http.Request, key string) (int, error) {
	val := r.URL.Query().Get(key)
	if val == "" {
		return 0, nil
	}
	return strconv.Atoi(val)
}

func GetIntUrlParam(r *http.Request, key string) (int, error) {
	val := chi.URLParam(r, key)
	if val == "" {
		return 0, nil
	}
	return strconv.Atoi(val)
}
