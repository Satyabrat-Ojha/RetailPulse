package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"retail-pulse/internal/db"
	"strconv"
)

// GetJobStatus handles job status retrieval
func GetJobStatus(w http.ResponseWriter, r *http.Request) {
	// Extract job ID from query parameters
	jobid := r.URL.Query().Get("jobid")
	if jobid == "" {
		http.Error(w, `{"error": "job ID is required"}`, http.StatusBadRequest)
		return
	}

	// Convert job ID to integer
	jobidInt, err := strconv.Atoi(jobid)
	if err != nil {
		log.Printf("error converting job id to integer: %s", err.Error())
		http.Error(w, `{"error": "invalid job ID"}`, http.StatusBadRequest)
		return
	}

	// Get current status of the job
	currentStatus, err := db.GetJobStatus(int64(jobidInt))
	if err != nil {
		log.Printf("error getting job status: %s", err.Error())
		http.Error(w, `{"error": "unable to retrieve job status"}`, http.StatusBadRequest)
		return
	}

	// Prepare the response
	response := map[string]interface{}{
		"status": currentStatus,
		"job_id": jobid,
	}

	// If the job failed, fetch errors and include them in the response
	if currentStatus == "failed" {
		errors, err := db.GetErrorsByJobID(int64(jobidInt))
		if err != nil {
			log.Printf("error fetching job errors: %s", err.Error())
			http.Error(w, `{"error": "unable to retrieve job errors"}`, http.StatusInternalServerError)
			return
		}
		response["error"] = errors // Assuming `errors` is in the required format
	}

	// Send the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
