package middleware

import "github.com/gin-gonic/gin"

func RequireAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !ctx.GetBool("is_admin") {
			ctx.AbortWithStatusJSON(403, gin.H{"error": "admin only"})
			return
		}
		ctx.Next()
	}
}
