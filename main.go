package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/muskiteer/Ai-Scam/routes"
)

func main() {
	// Load .env file FIRST before anything else
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	} else {
		log.Println("Loaded .env file successfully")
	}

	// Get PORT from environment, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Setup routes
	mux := http.NewServeMux()
	routes.SetupRoutes(mux)

	// Start server
	addr := ":" + port
	log.Printf("Starting server on %s", addr)
	log.Printf("API Key configured: %v", os.Getenv("API_KEY") != "")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
