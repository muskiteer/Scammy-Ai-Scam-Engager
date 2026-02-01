package main

import (
	"net/http"
	"log"

	"github.com/muskiteer/Ai-Scam/routes"
	"github.com/joho/godotenv"
)

func main() {
	mux := http.NewServeMux()

	routes.SetupRoutes(mux)
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env file")
	}
	
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}	