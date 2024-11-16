package api

import (
	"retail-pulse/internal/api/handlers"

	"github.com/gorilla/mux"
)

// SetupRoutes initializes all API routes
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Job routes
	router.HandleFunc("/api/submit/", handlers.SubmitJob).Methods("POST")
	router.HandleFunc("/api/status", handlers.GetJobStatus).Methods("GET")

	return router
}
