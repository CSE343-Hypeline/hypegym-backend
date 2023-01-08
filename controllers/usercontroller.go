package controllers

import (
	"fmt"
	"hypegym-backend/initializers"
	"hypegym-backend/models"
	"hypegym-backend/models/enums"
	"net/http"
	"strconv"
	"strings"
	"time"

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

	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	createSubEntry(&user)
	context.JSON(http.StatusCreated, gin.H{
		"ID": user.ID,
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

func UserUpdate(context *gin.Context) {
	id := context.Param("id")
	var user models.User
	var dto models.UserUpdateeDto
	if result := initializers.DB.First(&user, id); result.Error != nil {
		context.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	if err := context.ShouldBindJSON(&dto); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	user.Name = dto.Name
	user.Email = dto.Email
	user.PhoneNumber = dto.PhoneNumber
	user.Gender = dto.Gender
	user.Address = dto.Address

	initializers.DB.Save(&user)
	context.JSON(http.StatusOK, gin.H{
		"message": "Updated successfuly",
	})
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

func CheckIn(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(400, gin.H{"message": "invalid user ID"})
		return
	}

	gymID, err := strconv.Atoi(c.Param("gymID"))
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	var gender string
	initializers.DB.Table("users").Select("gender").Where("id = ?", userID).Scan(&gender)
	tableName := fmt.Sprintf("gym_activities_%d", gymID)
	var latestActivity models.UserActivity
	err = initializers.DB.Table(tableName).Where("user_id = ?", userID).Order("check_in_at desc").First(&latestActivity).Error
	if err == nil && latestActivity.CheckOutAt == nil {
		c.JSON(400, gin.H{"message": "user must check out before checking in again"})
		return
	}

	checkInTime := time.Now()
	activity := models.UserActivity{
		UserID:    userID,
		Gender:    enums.Gender(gender),
		CheckInAt: checkInTime,
	}
	initializers.DB.Table(tableName).Create(&activity)
	key := fmt.Sprintf("online-users-%d", gymID)
	initializers.RDB.SAdd(initializers.CTX, key, userID)

	c.JSON(200, gin.H{"message": "user checked in"})
}

func CheckOut(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		c.JSON(400, gin.H{"message": "invalid user ID"})
		return
	}
	gymID, err := strconv.Atoi(c.Param("gymID"))
	if err != nil {
		c.JSON(400, gin.H{"message": "invalid gym ID"})
		return
	}

	checkOutTime := time.Now()
	tableName := fmt.Sprintf("gym_activities_%d", gymID)
	result := initializers.DB.Table(tableName).Where("user_id = ? AND check_out_at IS NULL", userID).UpdateColumn("check_out_at", checkOutTime)
	if result.Error != nil {
		c.JSON(500, gin.H{"message": result.Error.Error()})
		return
	}

	key := fmt.Sprintf("online-users-%d", gymID)
	initializers.RDB.SRem(initializers.CTX, key, userID)

	c.JSON(200, gin.H{"message": "user checked out"})
}
