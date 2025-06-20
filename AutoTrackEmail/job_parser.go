package main

import (
	"regexp"
	"strings"
)

// Extract job info from email data
func ParseJobFromEmail(email EmailData) *Job {
	company := extractCompanyFromEmail(email.Subject, email.Body)
	title := extractTitleFromEmail(email.Subject, email.Body)
	status := determineJobStatus(email.Subject)

	// Skip if cant find data
	if company == "" && title == "" {
		return nil
	}

	return &Job{
		Company:     company,
		Title:       title,
		Status:      status,
		AppliedDate: email.Date,
		Notes:       email.Subject,
	}
}

// extract company info
func extractCompanyFromEmail(subject, body string) string {
	//First try subject
	company := extractCompanyFromSubject(subject)
	if company != "Unknown Company" {
		return company
	}

	//then try body
	return extractCompanyFromBody(body)
}

// extract title
func extractTitleFromEmail(subject, body string) string {
	//First try subject
	title := extractTitleFromSubject(subject)
	if title != "Unknown Position" {
		return title
	}

	//then try body
	return extractTitleFromBody(body)
}

// extract company from subject
func extractCompanyFromSubject(subject string) string {
	// LinkedIn-specific patterns first
	linkedInPatterns := []string{
		`(?i)your application was sent to\s+(.+?)(?:\s*$|\s*\|)`,
		`(?i)application sent.*?at\s+(.+?)(?:\s*$|\s*\|)`,
		`(?i)(.+?):\s*your application`,
		`(?i)application.*?to\s+(.+?)(?:\s+for|\s*$)`,
		`(?i)thank you for applying to\s+(.+?)(?:\s+for|\s*$)`,
	}

	for _, pattern := range linkedInPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(subject)
		if len(matches) > 1 {
			company := strings.TrimSpace(matches[1])
			// Clean up common suffixes
			company = strings.TrimSuffix(company, " - LinkedIn")
			company = strings.TrimSuffix(company, " via LinkedIn")
			if len(company) > 1 && len(company) < 100 {
				return company
			}
		}
	}

	// Common patterns for company names in email subjects
	patterns := []string{
		`(?i)from\s+([A-Za-z\s&.]+?)(?:\s+team|\s+careers|\s+hiring|$)`,
		`(?i)at\s+([A-Za-z\s&.]+?)(?:\s+team|\s+careers|\s+-|$)`,
		`(?i)([A-Za-z\s&.]+?)\s+(?:team|careers|hiring|hr)`,
		`(?i)^([A-Za-z\s&.]+?)\s+(?:-|:)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(subject)
		if len(matches) > 1 {
			company := strings.TrimSpace(matches[1])
			if len(company) > 1 && len(company) < 50 {
				return company
			}
		}
	}

	// Fallback: try to extract from common email patterns
	if strings.Contains(strings.ToLower(subject), "thank you for your application") {
		// Look for company name before common phrases
		words := strings.Fields(subject)
		for i, word := range words {
			if strings.ToLower(word) == "team" && i > 0 {
				return words[i-1]
			}
		}
	}

	return "Unknown Company"
}

// extract company from body
func extractCompanyFromBody(body string) string {
	if body == "" {
		return "Unknown Company"
	}

	lines := strings.Split(body, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Look for common patterns in email body
		patterns := []string{
			`(?i)from\s+([A-Za-z\s&.]+?)\s+team`,
			`(?i)at\s+([A-Za-z\s&.]+?)[,.]`,
			`(?i)([A-Za-z\s&.]+?)\s+careers`,
			`(?i)([A-Za-z\s&.]+?)\s+hiring`,
		}

		for _, pattern := range patterns {
			re := regexp.MustCompile(pattern)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				company := strings.TrimSpace(matches[1])
				if len(company) > 1 && len(company) < 50 {
					return company
				}
			}
		}
	}

	return "Unknown Company"
}

// extract title from subject
func extractTitleFromSubject(subject string) string {

	// Common patterns for job titles
	patterns := []string{
		`(?i)for\s+the\s+([^-]+?)(?:\s+position|\s+role)`,
		`(?i)position:\s+([^-\n]+?)(?:\s+at|\s+-|$)`,
		`(?i)role:\s+([^-\n]+?)(?:\s+at|\s+-|$)`,
		`(?i)([A-Za-z\s]+?)\s+position`,
		`(?i)([A-Za-z\s]+?)\s+role`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(subject)
		if len(matches) > 1 {
			title := strings.TrimSpace(matches[1])
			if len(title) > 1 && len(title) < 100 {
				return title
			}
		}
	}

	// Look for common job titles
	commonTitles := []string{
		"software engineer", "software developer", "backend developer", "developer",
		"programmer", "business analyst", "analyst",
		"manager", "director", "consultant", "specialist", "coordinator",
		"associate", "senior", "junior", "intern", "architect",
	}

	subjectLower := strings.ToLower(subject)
	for _, title := range commonTitles {
		if strings.Contains(subjectLower, title) {
			return title
		}
	}

	return "Unknown Position"
}

// extract title from body
func extractTitleFromBody(body string) string {
	if body == "" {
		return "Unknown Position"
	}

	lines := strings.Split(body, "\n")

	// LinkedIn-specific logic: look for job title after company line
	foundCompanyLine := false
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Remove markdown formatting
		cleanLine := strings.ReplaceAll(line, "**", "")
		cleanLine = strings.ReplaceAll(cleanLine, "*", "")
		cleanLine = strings.TrimSpace(cleanLine)

		// Check if this line contains "Your application was sent to"
		if strings.Contains(strings.ToLower(cleanLine), "your application was sent to") {
			foundCompanyLine = true
			continue
		}

		// If we found the company line, the next non-empty line should be the job title
		if foundCompanyLine && cleanLine != "" {
			// Clean up the job title line - remove company name and location info
			title := cleanLine

			// Remove company name and location patterns
			// Look for patterns like "Company Name · Location" or "Company Name - Location"
			parts := strings.Split(title, "·")
			if len(parts) > 1 {
				title = strings.TrimSpace(parts[0])
			}

			parts = strings.Split(title, " - ")
			if len(parts) > 1 && !strings.Contains(strings.ToLower(parts[1]), "developer") {
				title = strings.TrimSpace(parts[0])
			}

			// Remove "Applied on" and date information
			if strings.Contains(strings.ToLower(title), "applied on") {
				continue
			}

			// Validate the title
			if len(title) > 1 && len(title) < 100 && !isLocationOrDate(title) {
				return title
			}
		}
	}

	// Fallback to general body parsing
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Look for job title patterns in body
		patterns := []string{
			`(?i)position:\s+([^,\n]+)`,
			`(?i)role:\s+([^,\n]+)`,
			`(?i)for\s+the\s+([^,\n]+?)\s+position`,
			`(?i)([A-Za-z\s]+?)\s+position`,
		}

		for _, pattern := range patterns {
			re := regexp.MustCompile(pattern)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				title := strings.TrimSpace(matches[1])
				if len(title) > 1 && len(title) < 100 {
					return title
				}
			}
		}
	}

	return "Unknown Position"
}

// get job status from subject
func determineJobStatus(subject string) string {
	subjectLower := strings.ToLower(subject)

	// Check for specific status indicators
	if strings.Contains(subjectLower, "interview") {
		return "interview"
	}
	if strings.Contains(subjectLower, "offer") || strings.Contains(subjectLower, "congratulations") {
		return "offer"
	}
	if strings.Contains(subjectLower, "reject") || strings.Contains(subjectLower, "unfortunately") {
		return "rejected"
	}
	if strings.Contains(subjectLower, "received") || strings.Contains(subjectLower, "thank you") {
		return "applied"
	}

	// Default status
	return "applied"
}

// Helper function to check if a string looks like location or date info
func isLocationOrDate(text string) bool {
	textLower := strings.ToLower(text)

	// Common location indicators
	locationWords := []string{"toronto", "vancouver", "montreal", "calgary", "ottawa", "remote", "hybrid", "on-site", "ontario", "canada", "usa", "bc", "ab", "qc"}
	for _, word := range locationWords {
		if strings.Contains(textLower, word) {
			return true
		}
	}

	// Date patterns
	if strings.Contains(textLower, "applied on") || strings.Contains(textLower, "2024") || strings.Contains(textLower, "2025") {
		return true
	}

	// Check if it contains only location-like patterns (City, Province format)
	if strings.Contains(text, ",") && len(strings.Split(text, ",")) == 2 {
		return true
	}

	return false
}
