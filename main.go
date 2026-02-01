package main

import (
	"net/http"
	"log"

	"github.com/muskiteer/Ai-Scam/routes"
)

func main() {
	mux := http.NewServeMux()

	routes.SetupRoutes(mux)
	
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
	
}	