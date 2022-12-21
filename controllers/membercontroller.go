package controllers

import (
	"hypegym-backend/initializers"
	"hypegym-backend/models"
	"log"
)

func MemberCreate(user *models.User) {
	member := models.Member{}
	member.UserID = int(user.ID)
	record := initializers.DB.Omit("TrainerID").Create(&member)

	if record.Error != nil {
		log.Fatal(record.Error)
	}
}
