package worker

import (
	"fmt"
	"image"
	_ "image/jpeg" // Required for decoding JPEG images
	_ "image/png"  // Required for decoding PNG images
	"log"
	"math/rand"
	"net/http"
	"retail-pulse/internal/db"
	"retail-pulse/internal/models"
	"time"
)

// ProcessJob processes each image in the job, calculating perimeter and simulating processing delay
func Worker(job models.Job) {
	log.Printf("Starting job %d", job.JobID)
	failed := 0
	// Loop over each visit in the job
	for _, visit := range job.Visits {
		log.Printf("Processing visit for store ID %s", visit.StoreID)

		// Process each image URL in the visit
		for _, imageURL := range visit.ImageURLs {
			// Mock downloading the image
			log.Printf("Downloading image from URL: %s", imageURL)

			// Calculate the perimeter (mock calculation with random dimensions)
			perimeter, err := CalculateImagePerimeter(imageURL)
			if err != nil {
				db.InsertError(job.JobID, visit.StoreID, err.Error())
				failed++
				break
			}

			// Sleep for a random time to simulate GPU processing
			sleepTime := time.Duration(rand.Intn(300)+100) * time.Millisecond
			time.Sleep(sleepTime)

			// Log the result (this can also be stored in a database or result file)
			log.Printf("Processed image %s with perimeter: %d", imageURL, perimeter)
		}
	}
	if failed == 0 {
		db.UpdateJobStatus(job.JobID, "completed")
		log.Printf("Completed job %d", job.JobID)
	} else {
		db.UpdateJobStatus(job.JobID, "failed")
	}
}

// CalculateImagePerimeter downloads an image from a URL, decodes it, and calculates its perimeter
func CalculateImagePerimeter(url string) (int, error) {
	// Send HTTP GET request to download the image
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to download image: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response is successful
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to download image: status code %d", resp.StatusCode)
	}

	// Decode the image to get its dimensions
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to decode image: %v", err)
	}

	// Get image dimensions
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Calculate perimeter
	perimeter := 2 * (width + height)
	return perimeter, nil
}
