package main

import (
	"log"
	"net/http"
	"os"
	"retail-pulse/internal/api"
	"retail-pulse/internal/db"

	"github.com/joho/godotenv"
)

func main() {
	// Initialize the database
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	db.InitDB()
	defer db.CloseDB()

	// Set up routes
	router := api.SetupRoutes()
	port := "8080"
	port = os.Getenv("PORT")
	// Start the server
	log.Printf("Server running at http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
