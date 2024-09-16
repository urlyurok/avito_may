package controller

import (
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte("ok"))
	if err != nil {
		http.Error(w, "server is not ready", http.StatusInternalServerError)
	}
}
