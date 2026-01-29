package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireApprovedUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		status := ctx.GetString("user_status")
		if status != "APPROVED" {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "account not approved",
			})
			return
		}

		ctx.Next()
	}
}
