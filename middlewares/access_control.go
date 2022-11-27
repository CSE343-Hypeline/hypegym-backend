package middlewares

import (
	"hypegym-backend/models/enums"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

func AccessControl(validRoles []enums.Role) gin.HandlerFunc {
	return func(context *gin.Context) {
		role, exist := context.Get("role")

		if !exist {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Some problem on access control middleware"})
			context.Abort()
			return
		}

		if !slices.Contains(validRoles, role.(enums.Role)) {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			context.Abort()
			return
		}
		context.Next()
	}
}
