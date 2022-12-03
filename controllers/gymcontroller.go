package controllers

import (
	"hypegym-backend/initializers"
	"hypegym-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GymCreate(context *gin.Context) {

	var gym models.Gym
	if err := context.ShouldBindJSON(&gym); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	record := initializers.DB.Create(&gym)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, gin.H{})
}
