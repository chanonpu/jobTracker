package main

type Job struct {
	Company     string `json:"company"` //company name
	Title       string `json:"title"`   //job title
	Status      string `json:"status"`  // applied,interview,offer,rejected
	AppliedDate string `json:"applied_date"`
	Notes       string `json:"notes"`
}
