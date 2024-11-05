package helpers

import "net/http"

func SendError(w http.ResponseWriter, errI error) {
	err, ok := errI.(*CustomError)
	if !ok {
		http.Error(w, errI.Error(), http.StatusInternalServerError)
	} else {
		http.Error(w, err.Message, err.Code)
	}
}
