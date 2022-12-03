package middlewares

import (
	"hypegym-backend/auth"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString, err := context.Cookie("Authorization")
		if err != nil {
			context.JSON(401, gin.H{"error": "request does not contain an access token"})
			context.Abort()
			return
		}
		validateErr, claims := auth.ValidateToken(tokenString)
		if validateErr != nil {
			context.JSON(401, gin.H{"error": validateErr.Error()})
			context.Abort()
			return
		}
		context.Set("jwt", claims)
		context.Next()
	}
}
