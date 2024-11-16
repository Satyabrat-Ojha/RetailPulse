package handlers

import (
	"encoding/json"
	"net/http"
	"retail-pulse/internal/db"
	"retail-pulse/internal/models"
	"retail-pulse/internal/worker"
)

// JobSubmission represents the structure of a submitted job
type JobSubmission struct {
	Count  int            `json:"count"`
	Visits []models.Visit `json:"visits"`
}

// SubmitJob handles job submissions
func SubmitJob(w http.ResponseWriter, r *http.Request) {
	var job JobSubmission

	// Parse the request body
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, `{"error": "invalid payload"}`, http.StatusBadRequest)
		return
	}

	// Validate count
	if len(job.Visits) != job.Count {
		http.Error(w, `{"error": "count mismatch"}`, http.StatusBadRequest)
		return
	}

	// Insert job into DB
	jobID, err := db.InsertJob("created")
	if err != nil {
		http.Error(w, `{"error": "failed to insert job"}`, http.StatusInternalServerError)
		return
	}

	// Insert visits into DB
	for _, visit := range job.Visits {
		// Convert image URLs to JSON format
		imageURLsJSON, err := json.Marshal(visit.ImageURLs)
		if err != nil {
			http.Error(w, `{"error": "failed to marshal image URLs"}`, http.StatusInternalServerError)
			return
		}

		// Insert visit into DB
		_, err = db.InsertVisit(jobID, visit.StoreID, visit.VisitTime, string(imageURLsJSON))
		if err != nil {
			http.Error(w, `{"error": "failed to insert visit"}`, http.StatusInternalServerError)
			return
		}
	}

	go worker.Worker(models.Job{
		JobID:  jobID,
		Status: "processing",
		Visits: job.Visits,
	})

	// Return the job ID
	response := map[string]interface{}{
		"job_id": jobID,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
