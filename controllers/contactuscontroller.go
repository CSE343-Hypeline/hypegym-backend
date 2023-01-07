package controllers

import (
	"hypegym-backend/models"
	"hypegym-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Contact(context *gin.Context) {
	var contactUs models.Contact

	if err := context.ShouldBindJSON(&contactUs); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	services.SendContactUsMail(contactUs)
	context.JSON(http.StatusOK, gin.H{"message": "Mail send successfuly"})
}
