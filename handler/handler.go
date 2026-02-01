package handler

import (
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func StartConvo(w http.ResponseWriter, r *http.Request) {
	
}
	