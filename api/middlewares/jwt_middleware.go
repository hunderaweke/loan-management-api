package middlewares

import (
	"loan-management/pkg/infrastructures"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is missing"})
			ctx.Abort()
			return
		}
		header := strings.Split(authHeader, " ")
		if strings.ToLower(header[0]) != "bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}
		tokenString := header[1]
		claims, err := infrastructures.ValidateJWTToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.Set("userID", claims.UserID)
		ctx.Set("email", claims.Email)
		ctx.Set("isAdmin", claims.IsAdmin)

		ctx.Next()
	}
}
