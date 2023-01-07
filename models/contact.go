package models

type Contact struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name" `
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
}
