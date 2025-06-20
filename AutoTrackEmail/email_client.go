package main

import (
	"io"
	"log"
	"mime"
	"mime/multipart"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
)

type EmailData struct {
	Subject string
	Date    string
	Body    string
}

// Connect to Gmail IMAP server
func ConnectToGmail(email, password string) (*client.Client, error) {
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		return nil, err
	}

	if err = c.Login(email, password); err != nil {
		c.Logout()
		return nil, err
	}

	log.Println("Successfully logged in to Gmail")
	return c, nil
}

// FetchUnreadEmails retrieves unread job-related emails

func FetchUnreadEmails(c *client.Client) ([]EmailData, error) {
	// Select inbox
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return nil, err
	}

	if mbox.Messages == 0 {
		log.Println("No messages in inbox")
		return []EmailData{}, nil
	}

	// Search for unread messages
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{"\\Seen"}       // get unread message
	criteria.Since = time.Now().Add(-24 * time.Hour) // get only last 24 hour message
	seqNums, err := c.Search(criteria)
	if err != nil {
		return nil, err
	}

	if len(seqNums) == 0 {
		log.Println("No unread emails found")
		return []EmailData{}, nil
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(seqNums...)

	//Fetch messages
	section := &imap.BodySectionName{}
	messages := make(chan *imap.Message, 10)
	err = c.Fetch(seqSet, []imap.FetchItem{imap.FetchEnvelope, section.FetchItem()}, messages)
	if err != nil {
		log.Fatal(err)
	}

	var emails []EmailData

	//process messages
	for msg := range messages {
		r := msg.GetBody(section)
		//if no body skip it
		if r == nil {
			log.Println("Empty body, skipping...")
			continue
		}

		mr, err := message.Read(r)
		if err != nil {
			log.Println("Message parse error:", err)
			continue
		}

		header := mr.Header
		subject := header.Get("Subject")
		dateStr := header.Get("Date")

		//parse the date string
		date, err := parseEmailDate(dateStr)
		if err != nil {
			log.Println("Error parsing date:", err)
			continue
		}

		//Filter job-related emails
		if isJobRelatedEmail(subject) {
			log.Printf("Found job email: %s (%s)\n", subject, date.Format("2006-01-02"))

			//extract email body
			body := extractEmailBody(mr)

			emails = append(emails, EmailData{
				Subject: subject,
				Date:    date.Format("2006-01-02"),
				Body:    body,
			})
		}
	}

	return emails, nil
}

// check if email is job-related by subject
func isJobRelatedEmail(subject string) bool {
	subject = strings.ToLower(subject)

	jobKeywords := []string{"application", "thank you for apply", "next step"} // will add more

	for _, keyword := range jobKeywords {
		if strings.Contains(subject, keyword) {
			return true
		}
	}

	return false
}

// parse Gmail date format (if change provider need to be changed)
// Gmail uses RFC 2822 format: "Mon, 02 Jan 2006 15:04:05 -0700"
func parseEmailDate(dateStr string) (time.Time, error) {
	return time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", dateStr)
}

// extract body from email message
func extractEmailBody(mr *message.Entity) string {
	// check if it's a multipart message
	if mr.Header.Get("Content-Type") != "" {
		mediaType, params, err := mime.ParseMediaType(mr.Header.Get("Content-Type"))
		if err == nil && strings.HasPrefix(mediaType, "multipart/") {
			return extractMultipartBody(mr, params["boundary"])
		}
	}

	//single part message
	body, err := io.ReadAll(mr.Body)
	if err != nil {
		log.Println("Error reading email body:", err)
		return ""
	}

	return string(body)

}

// multipart email message
func extractMultipartBody(mr *message.Entity, boundary string) string {
	multipartReader := multipart.NewReader(mr.Body, boundary)

	var textBody, htmlBody string

	for {
		part, err := multipartReader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Error reading multipart:", err)
			break
		}

		contentType := part.Header.Get("Content-Type")
		mediaType, _, _ := mime.ParseMediaType(contentType)

		body, err := io.ReadAll(part)
		if err != nil {
			continue
		}

		switch mediaType {
		case "text/plain":
			textBody = string(body)
		case "text/html":
			if textBody == "" { //only use HTML if no plain text
				htmlBody = stripHTMLTags(string(body))
			}
		}

		part.Close()
	}

	// Prefer plain text over HTML
	if textBody != "" {
		return textBody
	}
	return htmlBody
}

// removes basic HTML tags (simple implementation)
func stripHTMLTags(html string) string {
	// very basic HTML stripper
	lines := strings.Split(html, "\n")
	var result []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Remove common HTML tags
		line = strings.ReplaceAll(line, "<br>", "\n")
		line = strings.ReplaceAll(line, "<br/>", "\n")
		line = strings.ReplaceAll(line, "<br />", "\n")
		line = strings.ReplaceAll(line, "<p>", "")
		line = strings.ReplaceAll(line, "</p>", "\n")
		line = strings.ReplaceAll(line, "<div>", "")
		line = strings.ReplaceAll(line, "</div>", "\n")

		// Basic tag removal (not comprehensive)
		for strings.Contains(line, "<") && strings.Contains(line, ">") {
			start := strings.Index(line, "<")
			end := strings.Index(line[start:], ">")
			if end == -1 {
				break
			}
			line = line[:start] + line[start+end+1:]
		}

		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}
