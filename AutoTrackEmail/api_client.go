package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var ServerURL string

// set ServerURL
func SetServerURL(url string) {
	ServerURL = url
	log.Printf("Server URL set to: %s", ServerURL)
}

// send job data to backend API
func sendJobToAPI(job Job) {
	jsonData, err := json.Marshal(job)
	if err != nil {
		log.Println("Error marshailing job data:", err)
		return
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	//create request
	url := fmt.Sprintf("%s/jobs", ServerURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Error creating req:", err)
	}

	//Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "JobEmailWatcher/1.0")

	//make request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending job to API:", err)
		return
	}
	defer resp.Body.Close()

	//handle response
	switch resp.StatusCode {
	case http.StatusOK:
		log.Printf("✓ Job saved successfully: %s at %s", job.Title, job.Company)
	case http.StatusCreated:
		log.Printf("✓ Job saved successfully: %s at %s", job.Title, job.Company)
	case http.StatusConflict:
		log.Printf("⚠ Job already exists: %s at %s", job.Title, job.Company)
	case http.StatusBadRequest:
		log.Printf("✗ Bad request for job: %s at %s", job.Title, job.Company)
	case http.StatusInternalServerError:
		log.Printf("✗ Server error when saving job: %s at %s", job.Title, job.Company)
	default:
		log.Printf("✗ Failed to save job: %s (status: %s)", job.Title, resp.Status)
	}
}
