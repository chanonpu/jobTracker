package models

type Job struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Company     string `json:"company"` //company name
	Title       string `json:"title"`   //job title
	Status      string `json:"status"`  // applied,interview,offer,rejected
	AppliedDate string `json:"appled_date"`
	Notes       string `json:"notes"`
}
