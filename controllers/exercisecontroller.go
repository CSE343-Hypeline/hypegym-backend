package controllers

import (
	"hypegym-backend/initializers"
	"hypegym-backend/models"
	"net/http"

	"github.com/dranikpg/dto-mapper"
	"github.com/gin-gonic/gin"
)

func GetExercises(context *gin.Context) {
	var exercises []models.Exercise
	var response []models.ExerciseResponseDto = nil
	if result := initializers.DB.Find(&exercises); result.Error != nil {
		context.AbortWithError(http.StatusNotFound, result.Error)
		return
	}
	dto.Map(&response, exercises)
	context.JSON(http.StatusOK, &response)
}

func GetExercise(context *gin.Context) {
	id := context.Param("id")
	var exercise models.Exercise
	if result := initializers.DB.First(&exercise, id); result.Error != nil {
		context.AbortWithError(http.StatusNotFound, result.Error)
		return
	}
	context.JSON(http.StatusOK, &exercise)
}
