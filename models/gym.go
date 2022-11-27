package models

import (
	"gorm.io/gorm"
)

type Gym struct {
	gorm.Model
	Address     string `json:"address"`
	Email       string `json:"email"`
	Name        string `json:"name" gorm:"unique"`
	PhoneNumber string `json:"phone_number" gorm:"unique"`
}
