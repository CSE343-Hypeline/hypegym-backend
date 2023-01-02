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

	// create activities table for the gym
	if err := createGymActivitiesTable(gym.ID); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusCreated, gin.H{})
}

func createGymActivitiesTable(gymID uint) error {
	tableName := fmt.Sprintf("gym_activities_%d", gymID)
	return initializers.DB.Table(tableName).AutoMigrate(&models.UserActivity{})
}

func GetAttendance(c *gin.Context) {
	gymID, err := strconv.Atoi(c.Param("gymID"))
	if err != nil {
		c.JSON(400, gin.H{"message": "invalid gym ID"})
		return
	}

	path := c.Request.URL.String()
	now := time.Now()
	currentDay := now.Day()
	currentMonth := int(now.Month())
	currentYear := now.Year()
	countMale := 0
	countFemale := 0
	countOther := 0
	tableName := fmt.Sprintf("gym_activities_%d", gymID)
	var queryString string
	var activities []models.UserActivity

	if strings.Contains(path, "month") {
		queryString = "MONTH(check_in_at) = ? AND YEAR(check_in_at) = ?"
		initializers.DB.Table(tableName).Where(queryString, currentMonth, currentYear).Find(&activities)
	} else {
		queryString = "DAY(check_in_at) = ? AND MONTH(check_in_at) = ? AND YEAR(check_in_at) = ?"
		initializers.DB.Table(tableName).Where(queryString, currentDay, currentMonth, currentYear).Find(&activities)
	}

	for _, activity := range activities {
		if activity.Gender == enums.Male {
			countMale++
		} else if activity.Gender == enums.Female {
			countFemale++
		} else {
			countOther++
		}
	}

	c.JSON(200, gin.H{
		"attendance_count_male":   countMale,
		"attendance_count_female": countFemale,
		"attendance_count_other":  countOther,
	})
}

func GetAllOnlines(c *gin.Context) {
	gymID, err := strconv.Atoi(c.Param("gymID"))
	if err != nil {
		c.JSON(400, gin.H{"message": "invalid gym ID"})
		return
	}

	key := fmt.Sprintf("online-users-%d", gymID)
	onlineMemberIDs, err := initializers.RDB.SMembers(initializers.CTX, key).Result()
	if err != nil {
		c.JSON(500, gin.H{"message": "failed to retrieve online users"})
		return
	}

	c.JSON(200, gin.H{
		"online_user_ids": onlineMemberIDs,
	})
}
