package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Home(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "HOME PAGE HADİ BAKİMMM"})
}
