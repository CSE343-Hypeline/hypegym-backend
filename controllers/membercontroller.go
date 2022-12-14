package controllers

import (
	"hypegym-backend/initializers"
	"hypegym-backend/models"
	"log"
	"net/http"

	"github.com/dranikpg/dto-mapper"
	"github.com/gin-gonic/gin"
)

func MemberCreate(user *models.User) {
	member := models.Member{}
	member.UserID = int(user.ID)
	record := initializers.DB.Omit("trainer_id").Create(&member)

	if record.Error != nil {
		log.Fatal(record.Error)
	}
}

func GetTrainerOf(context *gin.Context) {
	id := context.Param("id")
	var member models.Member
	var trainer models.Trainer
	var trainerDto models.TrainerDto

	record := initializers.DB.Where("user_id = ?", id).Find(&member)
	if record.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	record2 := initializers.DB.Where("user_id = ?", member.TrainerID).Find(&trainer)
	if record2.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": record2.Error.Error()})
		context.Abort()
		return
	}

	dto.Map(&trainerDto, trainer)
	context.JSON(http.StatusOK, &trainerDto)
}

func AssignProgram(context *gin.Context) {
	id := context.Param("id")
	var member models.Member
	var program models.Program

	if err := context.ShouldBindJSON(&program); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	programRecord := initializers.DB.Create(&program)
	if programRecord.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": programRecord.Error.Error()})
		context.Abort()
		return
	}

	record := initializers.DB.Where("user_id = ?", id).Find(&member)
	if record.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	if err := initializers.DB.Model(&member).Association("Programs").Append(&program); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Assigned successfuly"})
}

func AssignPrograms(context *gin.Context) {
	id := context.Param("id")
	var member models.Member
	var programs []models.Program

	if err := context.ShouldBindJSON(&programs); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	programRecord := initializers.DB.Create(&programs)
	if programRecord.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": programRecord.Error.Error()})
		context.Abort()
		return
	}

	record := initializers.DB.Where("user_id = ?", id).Find(&member)
	if record.Error != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	if err := initializers.DB.Model(&member).Association("Programs").Append(&programs); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Assigned successfuly"})
}

func DismissProgram(context *gin.Context) {
	id := context.Param("id")
	var dto models.ProgramRequestDto
	var program models.Program
	var member models.Member

	record := initializers.DB.Where("user_id = ?", id).Find(&member)
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

	if result := initializers.DB.First(&program, dto.ProgramID); result.Error != nil {
		context.AbortWithError(http.StatusNotFound, result.Error)
		return
	}

	initializers.DB.Model(&member).Association("Programs").Delete(&program)
	context.JSON(http.StatusOK, gin.H{"message": "Dismissed successfuly"})
}

func GetPrograms(context *gin.Context) {
	id := context.Param("id")
	var member models.Member
	var programs []models.Program
	var response []models.ProgramResponseDto = nil

	record := initializers.DB.Where("user_id = ?", id).Find(&member)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}
	initializers.DB.Model(&member).Association("Programs").Find(&programs)
	dto.Map(&response, programs)
	exercises := make([]models.Exercise, len(response))

	for i := 0; i < len(response); i++ {
		initializers.DB.First(&exercises[i], response[i].ExerciseID)
		response[i].Exercise = exercises[i]
	}

	context.JSON(http.StatusOK, gin.H{"programs": response})
}
