
# Retail Pulse Image Processing Service

The Retail Pulse Job Processing Service is a backend service developed in Go to handle image processing tasks for retail stores. The service allows users to submit jobs containing multiple images from different store visits. It calculates the perimeter of each image and simulates GPU processing with random sleep times. The service is designed to handle large volumes of jobs concurrently, utilising background workers to process images asynchronously. The status of each job can be queried to track progress or retrieve error details if a failure occurs. The project supports both local and Docker-based environments for easy setup, and it is designed for scalability with efficient concurrency management.




## Documentation

[Documentation](https://docs.google.com/document/d/117_rjXCOkD-YcO42Y0MdvhypiowwgrXKLwtLoV5RxlU/edit?usp=sharing)


## API Reference

#### Submit  Job

```http
  POST /api/items

    {
        "count": 2,
        "visits": [
            {
                "store_id": "S00339218",
                "image_url": ["https://example.com/image1.jpg", "https://example.com/image2.jpg"],
                "visit_time": "2024-11-17T12:00:00Z"
            },
            {
            "store_id": "S01408764",
            "image_url": ["https://example.com/image3.jpg"],
            "visit_time": "2024-11-17T14:00:00Z"
            }
        ]
    }

```

#### Get Job Status

```http
  GET /api/status?jobid=1
```


## Installation

### Using Docker

```bash
  git clone https://github.com/Satyabrat-Ojha/RetailPulse.git
  cd RetailPulse 
  docker build -t retail-pulse-app .
  docker run --rm -p 8080:8080 retail-pulse-app
```
### Without Docker

```bash
  git clone https://github.com/Satyabrat-Ojha/RetailPulse.git
  cd RetailPulse 
  go run main.go
```
The service will be available at http://localhost:8080.
