package controllers

import (
	"hypegym-backend/initializers"
	"hypegym-backend/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TrainerCreate(user *models.UserRequestDto) {
	trainer := models.Trainer{}
	trainer.UserID = int(user.ID)
	record := initializers.DB.Create(&trainer)

	if record.Error != nil {
		log.Fatal(record.Error)
	}
}

func AssignMembers(context *gin.Context) {
	id := context.Param("id")
	var dto models.TrainerDto
	var member models.Member
	var trainer models.Trainer

	record := initializers.DB.Where("user_id = ?", id).Find(&trainer)
	if record.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	if err := context.ShouldBindJSON(&dto); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	record2 := initializers.DB.Where("user_id = ?", dto.UserID).First(&member)
	if record2.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	// Error checking missing
	initializers.DB.Model(&trainer).Association("Members").Append(&member)
	context.JSON(http.StatusOK, gin.H{"message": "Member assigned successfuly"})
}

func DismissMember(context *gin.Context) {
	id := context.Param("id")
	var dto models.TrainerDto
	var member models.Member
	var trainer models.Trainer

	record := initializers.DB.Where("user_id = ?", id).Find(&trainer)
	if record.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	if err := context.ShouldBindJSON(&dto); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	record2 := initializers.DB.Where("user_id = ?", dto.UserID).First(&member)
	if record2.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	initializers.DB.Model(&trainer).Association("Members").Delete(&member)
	context.JSON(http.StatusOK, gin.H{"message": "Member dismissed successfuly"})
}

func GetMembers(context *gin.Context) {
	id := context.Param("id")
	var trainer models.Trainer
	var members []models.Member

	record := initializers.DB.Where("user_id = ?", id).Find(&trainer)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	initializers.DB.Model(&trainer).Association("Members").Find(&members)
	context.JSON(http.StatusOK, gin.H{"members": members})
}
