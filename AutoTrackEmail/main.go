package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// load environment variable
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	email := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")
	serverURL := os.Getenv("SERVER_URL")

	if email == "" || password == "" {
		log.Fatal("Please set Email and password in .env file")
	}

	if serverURL == "" {
		serverURL = "http://localhost:8080" // set as fallback
		log.Println("Server URL not set, using default server")
	}

	// Set Server URL for the api client
	SetServerURL(serverURL)

	//Connect to email
	log.Println("Connecting to GMail...")
	client, err := ConnectToGmail(email, password)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer client.Logout() // logout at the end

	// fetch unread emails
	log.Println("Fetching unread emails...")
	emails, err := FetchUnreadEmails(client)
	if err != nil {
		log.Fatal("Failed to fetch emails:", err)
	}

	if len(emails) == 0 {
		log.Println("No unread job-related emails found")
		return
	}

	//post each email to db
	for _, email := range emails {
		//extract job data from email
		job := ParseJobFromEmail(email)

		if job != nil {
			sendJobToAPI(*job)
		}
	}

	log.Println("Email process completed")
}
