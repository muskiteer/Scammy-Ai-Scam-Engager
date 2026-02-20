package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/muskiteer/Ai-Scam/middleware"
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
	muxWithLogging := middleware.Logging(mux)

	// Configure server with very long timeouts for Render free tier
	// This prevents the server from killing long-running requests
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      muxWithLogging,
		ReadTimeout:  0,                 // No timeout - wait indefinitely for request
		WriteTimeout: 0,                 // No timeout - wait indefinitely for response
		IdleTimeout:  600 * time.Second, // 10 minutes idle before closing connection
	}

	log.Printf("Starting server on %s", server.Addr)
	log.Printf("API Key configured: %v", os.Getenv("API_KEY") != "")
	log.Printf("Timeouts: Read=NONE, Write=NONE (waits for request to complete)")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
