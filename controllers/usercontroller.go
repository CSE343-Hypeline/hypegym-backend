package controllers

import (
	"hypegym-backend/initializers"
	"hypegym-backend/models"
	"hypegym-backend/models/enums"
	"net/http"
	"strings"

	"github.com/dranikpg/dto-mapper"
	"github.com/gin-gonic/gin"
)

func UserCreate(context *gin.Context) {
	var user models.User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	if err := user.HashPassword(user.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	record := initializers.DB.Create(&user)
	createSubEntry(&user)

	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": user.Role + " created",
	})
}

func createSubEntry(user *models.User) {
	if user.Role == enums.MEMBER {
		MemberCreate(user)
	} else if user.Role == enums.PT {
		TrainerCreate(user)
	}
}

func UserGet(context *gin.Context) {
	id := context.Param("id")
	response := models.UserResponseDto{}
	var user models.User
	if result := initializers.DB.First(&user, id); result.Error != nil {
		context.AbortWithError(http.StatusNotFound, result.Error)
		return
	}
	dto.Map(&response, user)
	context.JSON(http.StatusOK, &response)
}

func UserGetAll(context *gin.Context) {
	gymID := context.Param("gymID")
	var response []models.UserResponseDto = nil
	var users []models.User
	if result := initializers.DB.Find(&users, "gym_id = ?", gymID); result.Error != nil {
		context.AbortWithError(http.StatusNotFound, result.Error)
		return
	}
	dto.Map(&response, users)
	context.JSON(http.StatusOK, &response)
}

func UserGetAllByRole(context *gin.Context) {
	gymID := context.Param("gymID")
	path := context.Request.URL.String()
	var response []models.UserResponseDto = nil
	var users []models.User
	var role enums.Role
	if strings.Contains(path, "members") {
		role = enums.MEMBER
	} else if strings.Contains(path, "pts") {
		role = enums.PT
	}
	if result := initializers.DB.Where("role = ?", role).Find(&users, "gym_id = ?", gymID); result.Error != nil {
		context.AbortWithError(http.StatusNotFound, result.Error)
		return
	}
	dto.Map(&response, users)
	context.JSON(http.StatusOK, &response)
}

func UserDelete(context *gin.Context) {
	id := context.Param("id")
	var user models.User
	if result := initializers.DB.First(&user, id); result.Error != nil {
		context.AbortWithError(http.StatusNotFound, result.Error)
		return
	}
	initializers.DB.Delete(&user)
	context.JSON(http.StatusOK, gin.H{
		"message": user.Email + " deleted",
	})
}
