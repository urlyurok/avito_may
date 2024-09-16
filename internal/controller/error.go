package controller

import (
	log "log/slog"
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, message string, code int) {
	http.Error(w, "{ \"message\": \""+message+"\" }", code)
	return
}

func HandleRequestError(w http.ResponseWriter, err error) {
	if err != nil {
		ErrorResponse(w, "invalid params", http.StatusBadRequest)
		log.Debug("err: " + err.Error())
		return
	}
}
