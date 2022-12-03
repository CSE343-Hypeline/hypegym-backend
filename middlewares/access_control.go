package middlewares

import (
	"hypegym-backend/auth"
	"hypegym-backend/models/enums"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

func AccessControl(validRoles []enums.Role) gin.HandlerFunc {
	return func(context *gin.Context) {
		claim, exist := context.Get("jwt")
		var role enums.Role
		if !exist {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Some problem on access control middleware"})
			context.Abort()
			return
		}

		if claim, ok := claim.(*auth.JWTClaim); ok {
			role = claim.Role
		}

		if !slices.Contains(validRoles, role) {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			context.Abort()
			return
		}
		context.Next()
	}
}
