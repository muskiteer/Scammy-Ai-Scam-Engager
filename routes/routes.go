package routes 

import (
	"net/http"
	"github.com/muskiteer/Ai-Scam/handler"
)

func SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", handler.HealthCheck)
	mux.HandleFunc("/api/engage", handler.StartConvo)
}
