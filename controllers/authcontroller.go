package controllers

import (
	"hypegym-backend/auth"
	"hypegym-backend/initializers"
	"hypegym-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(context *gin.Context) {
	var request models.UserLoginDto
	var user models.User
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}
	record := initializers.DB.Where("email = ?", request.Email).First(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	credentialError := user.CheckPassword(request.Password)
	if credentialError != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		context.Abort()
		return
	}

	tokenString, err := auth.GenerateJWT(user.Email, user.Role, user.GymID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.SetSameSite(http.SameSiteLaxMode)
	context.SetCookie("Authorization", tokenString, 3600*34, "", "", false, true)

	context.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Logout(context *gin.Context) {
	tokenString := "un4uthorized"
	context.SetSameSite(http.SameSiteLaxMode)
	context.SetCookie("Authorization", tokenString, 3600*34, "", "", false, true)
	context.JSON(http.StatusOK, gin.H{"message": "Susscessfuly logged out"})
}

func Me(context *gin.Context) {
	claim, exist := context.Get("jwt")

	if !exist {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Some problem on access control middleware"})
		context.Abort()
		return
	}
	context.JSON(http.StatusOK, claim)
}
