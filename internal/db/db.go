package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

var DB *sql.DB

// InitDB initializes the SQLite database
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "./retail_pulse.db") // Use "sqlite" for modernc.org/sqlite
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Create tables if they do not exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS jobs (
		job_id INTEGER PRIMARY KEY AUTOINCREMENT,
		status TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS visits (
		visit_id INTEGER PRIMARY KEY AUTOINCREMENT,
		job_id INTEGER NOT NULL,
		store_id TEXT NOT NULL,
		visit_time TEXT NOT NULL,
		image_urls TEXT NOT NULL,
		FOREIGN KEY (job_id) REFERENCES jobs(job_id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS errors (
		error_id INTEGER PRIMARY KEY AUTOINCREMENT,
		job_id INTEGER,
		store_id TEXT,
		error_message TEXT,
		FOREIGN KEY (job_id) REFERENCES jobs(job_id)
	);
	`

	_, err = DB.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}
	log.Println("Database initialized successfully")
}

// CloseDB closes the database connection
func CloseDB() {
	if err := DB.Close(); err != nil {
		log.Fatalf("Error closing database: %v", err)
	}
}

// InsertJob inserts a new job into the job table
func InsertJob(status string) (int64, error) {
	stmt, err := DB.Prepare("INSERT INTO jobs(status) VALUES (?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(status)
	if err != nil {
		return 0, err
	}

	jobID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return jobID, nil
}

// InsertVisit inserts a new visit into the visits table
func InsertVisit(jobID int64, storeID, visitTime, imageURLs string) (int64, error) {
	stmt, err := DB.Prepare("INSERT INTO visits(job_id, store_id, visit_time, image_urls) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(jobID, storeID, visitTime, imageURLs)
	if err != nil {
		return 0, err
	}

	visitID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return visitID, nil
}

// GetJob retrieves a job by its ID
func GetJobStatus(jobID int64) (string, error) {
	var status string
	err := DB.QueryRow("SELECT status FROM jobs WHERE job_id = ?", jobID).Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}

// InsertError inserts a new error record into the errors table
func InsertError(jobID int64, storeID, errorMessage string) (int64, error) {
	stmt, err := DB.Prepare("INSERT INTO errors(job_id, store_id, error_message) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(jobID, storeID, errorMessage)
	if err != nil {
		return 0, err
	}

	errorID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return errorID, nil
}

// GetErrorsByJobID retrieves all errors associated with a specific job ID
func GetErrorsByJobID(jobID int64) ([]map[string]interface{}, error) {
	rows, err := DB.Query("SELECT error_id, store_id, error_message FROM errors WHERE job_id = ?", jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var errors []map[string]interface{}
	for rows.Next() {
		var errorID int64
		var storeID, errorMessage string

		err = rows.Scan(&errorID, &storeID, &errorMessage)
		if err != nil {
			return nil, err
		}

		errorRecord := map[string]interface{}{
			// "error_id":      errorID,
			"store_id":      storeID,
			"error_message": errorMessage,
		}
		errors = append(errors, errorRecord)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return errors, nil
}

// UpdateJobStatus updates the status of a job in the jobs table by job ID
func UpdateJobStatus(jobID int64, newStatus string) error {
	stmt, err := DB.Prepare("UPDATE jobs SET status = ? WHERE job_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(newStatus, jobID)
	if err != nil {
		return err
	}

	return nil
}
