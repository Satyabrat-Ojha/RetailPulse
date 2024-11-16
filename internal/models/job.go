package models

// Job represents a job to be processed, containing job ID, status, and associated visits.
type Job struct {
	JobID  int64
	Status string
	Visits []Visit
}

// Visit represents a visit associated with a job, containing visit ID, store ID, visit time, and image URLs.
type Visit struct {
	StoreID   string   `json:"store_id"`
	ImageURLs []string `json:"image_url"`
	VisitTime string   `json:"visit_time"`
}
