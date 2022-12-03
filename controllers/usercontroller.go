package controllers

import (
	"hypegym-backend/auth"
	"hypegym-backend/initializers"
	"hypegym-backend/models"
	"hypegym-backend/models/enums"
	"net/http"
	"strings"

	"github.com/dranikpg/dto-mapper"
	"github.com/gin-gonic/gin"
)

type UserLoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponseDto struct {
	ID    uint   `json:"ID"`
	Email string `json:"email"`
	Role  string `json:"role"`
	GymID uint   `json:"gym_id"`
}

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
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"message": user.Role + " created",
	})
}

func UserGet(context *gin.Context) {
	id := context.Param("id")
	response := UserResponseDto{}
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
	var response []UserResponseDto = nil
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
	var response []UserResponseDto = nil
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
	context.JSON(http.StatusCreated, gin.H{
		"message": user.Email + " deleted",
	})
}

func UserLogin(context *gin.Context) {
	var request UserLoginDto
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

	tokenString, err := auth.GenerateJWT(user.Email, user.Role)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.SetSameSite(http.SameSiteLaxMode)
	context.SetCookie("Authorization", tokenString, 3600*34, "", "", false, true)

	context.JSON(http.StatusOK, gin.H{})
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
