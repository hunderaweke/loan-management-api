package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		isAdmin, exists := ctx.Get("isAdmin")
		if !exists || !isAdmin.(bool) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
